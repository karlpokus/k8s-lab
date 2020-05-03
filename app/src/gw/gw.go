package main

import (
	"log"
	"net/http"
	"context"
	"time"
	"io"
	"os"

	"github.com/julienschmidt/httprouter"
)

var client = http.Client{
	Timeout: 3 * time.Second, // tcp ttl
}

type post struct {
	Title, Body string
}

func ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

func remoteCall(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		defer req.Body.Close()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second) // http ttl
	defer cancel()
	return client.Do(req.WithContext(ctx))
}

func reply(w http.ResponseWriter, res *http.Response, err error, url string) {
	if err != nil {
		log.Printf("remote call to %s failed: %s", url, err)
		http.Error(w, "server err", 500)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		http.Error(w, res.Status, res.StatusCode)
		return
	}
	io.Copy(w, res.Body)
}

func auth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pwd, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "missing auth", 400)
			return
		}
		url := "http://user:9051/user"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("request construction of %s failed: %s", url, err)
			http.Error(w, "server error", 500)
			return
		}
		req.SetBasicAuth(user, pwd)
		res, err := remoteCall(req)
		if err != nil {
			log.Printf("remote call to %s failed: %s", url, err)
			http.Error(w, "server err", 500)
			return
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			http.Error(w, res.Status, res.StatusCode)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func getPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := "http://blog:9052/posts"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("request construction of %s failed: %s", url, err)
			http.Error(w, "server error", 500)
			return
		}
		res, err := remoteCall(req)
		reply(w, res, err, url)
	}
}

func getPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := "http://blog:9052/post/"
		title := httprouter.ParamsFromContext(r.Context()).ByName("title")
		if title == "" {
			http.Error(w, "missing title arg", 400)
			return
		}
		url = url + title
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("request construction of %s failed: %s", url, err)
			http.Error(w, "server error", 500)
			return
		}
		res, err := remoteCall(req)
		reply(w, res, err, url)
	}
}

func addPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := "http://blog:9052/post"
		req, err := http.NewRequest("POST", url, r.Body)
		if err != nil {
			log.Printf("request construction of %s failed: %s", url, err)
			http.Error(w, "server error", 500)
			return
		}
		res, err := remoteCall(req)
		reply(w, res, err, url)
	}
}

func logRequest(doLog string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if doLog == "yes" {
			log.Println(r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	}
}

func main() {
	router := httprouter.New()
	router.Handler("GET", "/ping", ping())
	router.Handler("GET", "/posts", auth(getPosts()))
	router.Handler("GET", "/post/:title", auth(getPost()))
	router.Handler("POST", "/post", auth(addPost()))
	log.Fatal(http.ListenAndServe("0.0.0.0:9050", logRequest(os.Getenv("LOG_REQUESTS"), router)))
}
