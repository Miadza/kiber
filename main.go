package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:123123123@tcp(127.0.0.1:3306)/computer_club")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

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

	http.Redirect(w, r, "/?registered=true", http.StatusSeeOther)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	registered := r.URL.Query().Get("registered") == "true"

	tmpl := `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Document</title>
    <link rel="stylesheet" href="./css/style.css" />
  </head>
  <body>
    <header class="header">
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
              <a id="reg" href="registration.html">Войти</a>
            {{end}}
          </div>
        </div>
      </div>
    </header>
    <main class="main">
      <div class="background-main">
        <img class="img-main" src="./image/GEO-TERRAIN4-20 1.svg" alt="" />
      </div>
      <div class="text-main">
        <h1 class="main-text-h1">Разные зоны с разными условиями</h1>
      </div>
      <div class="cards">
        <div class="cardone">
          <p class="other-us">Игровой ПК и переферия</p>
          <p class="text-nab-card">Стандарт</p>
          <img class="img-card-purple" src="./image/card1.svg" alt="" />
          <p class="zone">Общая зона</p>
          <p class="last-text-card-purple">Подробнее</p>
        </div>
        <div class="cardone">
          <p class="other-us">Игровой ПК и переферия</p>
          <p class="text-nab-card">Стандарт</p>
          <img class="img-card-aqua" src="./image/card2.svg" alt="" />
          <p class="zone">Общая зона</p>
          <p class="last-text-card-aqua">Подробнее</p>
        </div>
        <div class="cardone">
          <p class="other-us">Игровой ПК и переферия</p>
          <p class="text-nab-card">Стандарт</p>
          <img class="img-card-blue" src="./image/card3.svg" alt="" />
          <p class="zone">Общая зона</p>
          <p class="last-text-card-blue">Подробнее</p>
        </div>
      </div>
      <div class="background-stocks"></div>
      <div class="stocks">
        <h1>Акции и скидки на любой вкус</h1>
      </div>
      <div class="promokod">
        <div class="promone">
            <p class="promon">Промокод Cyber</p>
            <p class="promonn">Киберфикация</p>
            <p class="promonnn">Киберфикация</p>
          </div>
          <div class="promtwo">
            <p class="promtw">Промокод Cyber</p>
            <p class="promtww">Киберфикация</p>
          </div>
          <div class="prothree">
            <p class="promth">Промокод Cyber</p>
            <p class="promothh">Киберфикация</p>
          </div>
        </div>
      </div>
    </main>
    <footer class="footer">
      <div class="background-main">
        <img class="img-main" src="./image/GEO-TERRAIN4-20 1.svg" alt="" />
      </div>
      <div class="cards">
        <div class="card-search">
          <p class="how-search">Как нас найти</p>
        <div class="adr-tel">
          <div class="adrr">
            <p class="adres">Адрес: </p> <div class="adr"><p>Пушкинская 120/2</p></div>
          </div>
          <div class="tell">
            <p class="telephone">Телефон: <div class="tel"> <p> +7 977 320 88 88</p></div></p>
          </div>
        </div>
          <div class="button-br">
            <a href="#">Забронировать!</a>
          </div>
        </div>
        <div class="">
          <img src="./image/map.svg" alt="">
        </div>
      </div>
    </footer>
  </body>
</html>`

	t := template.Must(template.New("index").Parse(tmpl))
	t.Execute(w, struct{ Registered bool }{Registered: registered})
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/register", registerHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
