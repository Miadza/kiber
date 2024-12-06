package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	// _ "github.com/go-sql-driver/mysql"git
)

var db *sql.DB

// Инициализация базы данных
func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:123123123@tcp(127.0.0.1:3306)/computer_club")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

// Статические файлы
func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "."+r.URL.Path)
}

// Обработчик регистрации
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	// После регистрации перенаправляем на главную с параметром registered=true
	http.Redirect(w, r, "/?registered=true", http.StatusSeeOther)
}

// Обработчик главной страницы
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, зарегистрирован ли пользователь через параметр "registered"
	registered := r.URL.Query().Get("registered") == "true"

	// Путь к шаблону главной страницы
	tmpl := `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Document</title>
    <link rel="stylesheet" href="/css/style.css" />
  </head>
  <body>
    <header class="header">
      <div class="background">
        <img class="img1" src="/image/chair.svg" alt="" />
        <img class="img2" src="/image/back.svg" alt="" />
      </div>
      <div class="container">
        <div class="navigation">
          <div class="nav">
            <a href="#">Главная</a>
            <a href="#">Зоны</a>
            <a href="#">Цены</a>
            <a href="#">Фото</a>
            <a href="#">Акции</a>
            <a href="#">Контакты</a>
          </div>
          <div class="registration">
            {{if .Registered}}
              <span>Профиль</span>
            {{else}}
              <a id="reg" href="/registration">Войти</a>
            {{end}}
          </div>
        </div>
        <div class="text-header">
          <div class="logo-text">
            <img src="/image/head.svg" alt="" />
          </div>
          <div class="text-rostov">
            <h1>Кибертека</h1>
            <p>Ростов-на-Дону</p>
          </div>
        </div>
        <div class="text-button">
          <p>Работаем круглосуточно</p>
          <p>8 928 136 37 02</p>
          <div class="button-header">
            <a href="#">Забронировать!</a>
          </div>
        </div>
      </div>
    </header>
    <main class="main">
      <div class="background-main">
        <img class="img-main" src="/image/GEO-TERRAIN4-20 1.svg" alt="" />
      </div>
      <div class="text-main">
        <h1 class="main-text-h1">Разные зоны с разными условиями</h1>
      </div>
      <div class="cards">
        <div class="cardone">
          <p class="other-us">Игровой ПК и переферия</p>
          <p class="text-nab-card">Стандарт</p>
          <img class="img-card-purple" src="/image/card1.svg" alt="" />
          <p class="zone">Общая зона</p>
          <p class="last-text-card-purple">Подробнее</p>
        </div>
        <div class="cardone">
          <p class="other-us">Игровой ПК и переферия</p>
          <p class="text-nab-card">Стандарт</p>
          <img class="img-card-aqua" src="/image/card2.svg" alt="" />
          <p class="zone">Общая зона</p>
          <p class="last-text-card-aqua">Подробнее</p>
        </div>
        <div class="cardone">
          <p class="other-us">Игровой ПК и переферия</p>
          <p class="text-nab-card">Стандарт</p>
          <img class="img-card-blue" src="/image/card3.svg" alt="" />
          <p class="zone">Общая зона</p>
          <p class="last-text-card-blue">Подробнее</p>
        </div>
      </div>
    </main>
    <footer class="footer">
      <div class="background-main">
        <img class="img-main" src="/image/GEO-TERRAIN4-20 1.svg" alt="" />
      </div>
    </footer>
  </body>
</html>`

	// Обрабатываем шаблон с данными

	t := template.Must(template.New("index").Parse(tmpl))
	t.Execute(w, struct{ Registered bool }{Registered: registered})
}

func main() {
	initDB()
	defer db.Close()

	// Обработка статических файлов (CSS, изображения и т.д.)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	http.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("./static/image"))))

	// Обработчики
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/registration", registerHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
