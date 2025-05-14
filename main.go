package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/registrer.html", "templates/login.html", "templates/skins/skin1.html", "templates/skins/skin2.html", "templates/skins/skin3.html"))

var db *sql.DB

func index(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)

	err = templates.ExecuteTemplate(w, "index.html", user.Name)

	if err != nil {
		http.Error(w, "Klarte ikke laste inn side", http.StatusInternalServerError)
		return
	}
}

func besøksside(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("navn")

	room, _ := getRoom(username)

	style := fmt.Sprintf("%d", room.Style)

	err := templates.ExecuteTemplate(w, "skin"+style+".html", room)

	if err != nil {
		http.Error(w, "Klarte ikke laste inn side "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func registrer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		errorMessage := r.URL.Query().Get("error")

		err := templates.ExecuteTemplate(w, "registrer.html", errorMessage)

		if err != nil {
			http.Error(w, "Klarte ikke laste inn side", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		username := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		err := registerate(username, password, email)

		if err != nil {
			redirectURL := fmt.Sprintf("/register?error=%s", url.QueryEscape(err.Error()))
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		sessionToken, _ := generateToken(32)
		csrfToken, _ := generateToken(32)

		err = loggingIn(sessionToken, csrfToken, username, password, w)

		if err != nil {
			redirectURL := fmt.Sprintf("/register?error=%s", url.QueryEscape(err.Error()))
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)

		return
	}
}

func loggin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		errorMessage := r.URL.Query().Get("error")

		err := templates.ExecuteTemplate(w, "login.html", errorMessage)

		if err != nil {
			http.Error(w, "Klarte ikke laste inn side", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		username := r.FormValue("name")
		password := r.FormValue("password")

		sessionToken, _ := generateToken(32)
		csrfToken, _ := generateToken(32)

		err := loggingIn(sessionToken, csrfToken, username, password, w)

		if err != nil {
			redirectURL := fmt.Sprintf("/logginn?error=%s", url.QueryEscape(err.Error()))
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)

		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	user, err := getUser(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = csrfCheck(r, user.Csrf)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = loggingOut(w, user.Name)

	if err != nil {
		http.Error(w, "Error logging out "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)

	return
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	db, _ = createDB()

	http.HandleFunc("/registrer", registrer)
	http.HandleFunc("/logginn", loggin)
	http.HandleFunc("/loggut", logout)
	http.HandleFunc("/{navn}", besøksside)
	http.HandleFunc("/", index)

	http.ListenAndServe(":8080", nil)
}
