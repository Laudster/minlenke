package main

import (
	"database/sql"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/registrer.html"))

var db *sql.DB

func index(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)

	if err != nil {
		http.Error(w, "Klarte ikke laste inn side", http.StatusInternalServerError)
		return
	}
}

func registrer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := templates.ExecuteTemplate(w, "registrer.html", nil)

		if err != nil {
			http.Error(w, "Klarte ikke laste inn side", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		http.Error(w, "Midlertid ikke registrer post", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	db, _ = createDB()

	http.HandleFunc("/registrer", registrer)
	http.HandleFunc("/", index)

	http.ListenAndServe(":8080", nil)
}
