package main

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))

func index(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)

	if err != nil {
		http.Error(w, "Could not load tempalte", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", index)

	http.ListenAndServe(":8080", nil)
}
