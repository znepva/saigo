package main

import (
	"html/template"
	"net/http"
	"sync"
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

var homeT *template.Template

// make map concurrency safe
var counter = struct{
	sync.RWMutex
	m map[string]int
}{m: make(map[string]int)}


func loadData(dataFile string, hash map[string]int) error {
	// skip if file does not exists
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return nil
	}

	// open file
	file, err := os.Open(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// read each line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Fields(scanner.Text())

		//convert to int
		counter, err := strconv.Atoi(data[1])
		if err != nil {
			return err
		}
		hash[data[0]] = counter
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}


// render home page
func home(w http.ResponseWriter, r *http.Request) {
	// render view
	counter.RLock()
	err := homeT.Execute(w, &counter.m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	counter.RUnlock()
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
	counter.RUnlock()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func init() {
	filename := "solution/data.txt"

	err := loadData(filename, counter.m)
	if err != nil {
		fmt.Println("Something went wrong")
	}
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

