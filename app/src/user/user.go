package main

import (
	"log"
	"net/http"

	"go.elastic.co/apm/module/apmhttprouter"
)

func getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pwd, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "missing auth", 400)
			return
		}
		if user != "bob" || pwd != "dylan" {
			http.Error(w, "forbidden", 403)
			return
		}
		w.Write([]byte("ok"))
	}
}

func ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

func main() {
	router := apmhttprouter.New() // wraps httprouter
	router.Handler("GET", "/ping", ping())
	router.Handler("GET", "/user", getUser())
	log.Fatal(http.ListenAndServe("0.0.0.0:9051", router))
}
