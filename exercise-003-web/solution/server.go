package main

import (
	"html/template"
	"net/http"
	"sync"
)

var homeT *template.Template

// make map concurrency safe
var counter = struct{
	sync.RWMutex
	m map[string]int
}{m: make(map[string]int)}

// render home page
func home(w http.ResponseWriter, r *http.Request) {
	err := homeT.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// process post request
func signup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")

	// update user counter
	counter.Lock()
	counter.m[username]++
	counter.Unlock()

	// render view
	counter.RLock()
	err := homeT.Execute(w, &counter.m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	counter.RUnlock()
}

func main() {
	// initialize home template
	var err error
	homeT = template.Must(template.ParseFiles("solution/home.html"))
	if err != nil {
		panic(err)
	}

	// routes
	http.HandleFunc("/home", home)
	http.HandleFunc("/signup", signup)

	// start server
	http.ListenAndServe(":8080", nil)
}

