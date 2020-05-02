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

type store struct {
	Posts []post
}

func ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
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

func getPosts(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply(w, s.Posts)
	}
}

func getPost(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := httprouter.ParamsFromContext(r.Context()).ByName("title")
		for _, p := range s.Posts {
			if p.Title == title {
				reply(w, p)
				return
			}
		}
		http.Error(w, "post not found", 404)
	}
}

func addPost(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Malformed json body", 400)
			return
		}
		s.Posts = append(s.Posts, p)
		w.Write([]byte("ok"))
	}
}

func main() {
	s := &store{
		Posts: []post{
			{Title: "title-1", Body: "lorem ipsum 1"},
			{Title: "title-2", Body: "lorem ipsum 2"},
		},
	}
	router := httprouter.New()
	router.Handler("GET", "/ping", ping())
	router.Handler("GET", "/posts", getPosts(s))
	router.Handler("GET", "/post/:title", getPost(s))
	router.Handler("POST", "/post", addPost(s))
	log.Fatal(http.ListenAndServe("0.0.0.0:9052", router))
}
