package main

import (
	"html/template"
	"net/http"
)

var userData = make(map[string]int)

var homeT = template.Must(template.ParseFiles("solution/home.html"))

func home(w http.ResponseWriter, r *http.Request) {
	homeT.Execute(w, nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	userData[username]++
	homeT.Execute(w, &userData)
}

func main() {
	http.HandleFunc("/home", home)
	http.HandleFunc("/signup", signup)
	http.ListenAndServe(":8080", nil)
}

