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
func getPosts(posts []post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply(w, posts)
	}
}

func reply(w http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("unable to encode %v", data)
		http.Error(w, "server err", 500)
		return
	}
}

func getPost(posts []post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := httprouter.ParamsFromContext(r.Context()).ByName("title")
		for _, p := range posts {
			if p.Title == title {
				reply(w, p)
				return
			}
		}
		http.Error(w, "post not found", 404)
	}
}

func main() {
	posts := []post{
		{Title: "title-1", Body: "lorem ipsum 1"},
		{Title: "title-2", Body: "lorem ipsum 2"},
	}
	router := httprouter.New()
	router.Handler("GET", "/ping", ping())
	router.Handler("GET", "/posts", getPosts(posts))
	router.Handler("GET", "/post/:title", getPost(posts))
	log.Fatal(http.ListenAndServe("0.0.0.0:9052", router))
}
