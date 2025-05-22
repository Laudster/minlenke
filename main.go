package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/registrer.html", "templates/login.html", "templates/rediger.html", "templates/skins/skin1.html", "templates/skins/skin2.html", "templates/skins/skin3.html"))

var db *sql.DB

type Link struct {
	Title string
	Link  string
}

func index(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)

	err = templates.ExecuteTemplate(w, "index.html", user.Name)

	if err != nil {
		http.Error(w, "Klarte ikke laste inn side", http.StatusInternalServerError)
		return
	}
}

func rediger(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	room, _ := getRoom(user.Name)

	links := formatLinks(room.Links)

	data := map[string]any{
		"Room":  room,
		"Links": links,
	}

	err = templates.ExecuteTemplate(w, "rediger.html", data)

	if err != nil {
		http.Error(w, "Klarte ikke laste inn side", http.StatusInternalServerError)
		return
	}
}

func besøksside(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("navn")

	room, err := getRoom(username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	links := formatLinks(room.Links)

	user, _ := getUser(r)

	data := map[string]any{
		"Room":  room,
		"Image": base64.StdEncoding.EncodeToString(room.Image),
		"Links": links,
		"Admin": strings.ToUpper(username) == strings.ToUpper(user.Name),
	}

	err = templates.ExecuteTemplate(w, "skin"+room.Style+".html", data)

	if err != nil {
		http.Error(w, "Klarte ikke laste inn side "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveTheme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	user, _ := getUser(r)

	if user.Name == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := csrfCheck(r, user.Csrf)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("update rooms set style = $1 where user_id = $2", r.FormValue("theme"), user.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	return
}

func saveBody(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	user, _ := getUser(r)

	if user.Name != r.FormValue("user") || user.Name == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	err := csrfCheck(r, user.Csrf)

	if err != nil {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("update rooms set body = $1 where user_id = $2", r.FormValue("body"), user.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func saveLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	user, _ := getUser(r)

	if user.Name != r.FormValue("user") || user.Name == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := csrfCheck(r, user.Csrf)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	amount := r.FormValue("amount")

	amountInt, _ := strconv.Atoi(amount)

	links := ""

	for i := 0; i < amountInt; i++ {
		titleKey := fmt.Sprintf("Title%d", i)
		links += "#" + r.FormValue(titleKey)

		linkKey := fmt.Sprintf("Link%d", i)
		links += ";" + r.FormValue(linkKey)
	}

	_, err = db.Exec("update rooms set links = $1 where user_id = $2", links, user.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func saveImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	user, _ := getUser(r)

	if user.Name != r.FormValue("user") || user.Name == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := csrfCheck(r, user.Csrf)

	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	r.ParseMultipartForm(5 << 20)

	file, _, err := r.FormFile("image")

	if err != nil {
		http.Error(w, "Image uplaod errror: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	imageData, _ := io.ReadAll(file)

	_, err = db.Exec("update rooms set image = $1 where user_id = $2", imageData, user.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
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

		http.Redirect(w, r, "/rediger", http.StatusFound)
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

	http.HandleFunc("/save-theme", saveTheme)
	http.HandleFunc("/save-image", saveImage)
	http.HandleFunc("/save-links", saveLinks)
	http.HandleFunc("/save-body", saveBody)
	http.HandleFunc("/registrer", registrer)
	http.HandleFunc("/logginn", loggin)
	http.HandleFunc("/loggut", logout)
	http.HandleFunc("/{navn}", besøksside)
	http.HandleFunc("/rediger", rediger)
	http.HandleFunc("/", index)

	http.ListenAndServe(":8080", nil)
}
