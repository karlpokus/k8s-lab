package main

import (
	"log"
	"net/http"
	"encoding/json"

	"github.com/julienschmidt/httprouter"
)

type post struct {
	Title, Body string
}

func ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}
func blogs(posts []post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Println("unable to encode posts json")
			http.Error(w, "server err", 500)
			return
		}
	}
}

func main() {
	posts := []post{
		{Title: "title 1", Body: "lorem ipsum 1"},
		{Title: "title 2", Body: "lorem ipsum 2"},
	}
	router := httprouter.New()
	router.Handler("GET", "/ping", ping())
	router.Handler("GET", "/blogs", blogs(posts))
	log.Fatal(http.ListenAndServe("0.0.0.0:9052", router))
}
