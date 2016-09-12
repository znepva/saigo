package main

import (
	"html/template"
	"net/http"
)

var homeT *template.Template

var userData = make(map[string]int)

func home(w http.ResponseWriter, r *http.Request) {
	err := homeT.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	userData[username]++

	err := homeT.Execute(w, &userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	var err error
	homeT = template.Must(template.ParseFiles("solution/home.html"))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/home", home)
	http.HandleFunc("/signup", signup)
	http.ListenAndServe(":8080", nil)
}

