package main

import (
	"log"
	"net/http"
	"context"
	"time"
	"io"

	"github.com/julienschmidt/httprouter"
)

var client = http.Client{
	Timeout: 3 * time.Second, // tcp ttl
}

func remoteCall(req *http.Request) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second) // http ttl
	defer cancel()
	return client.Do(req.WithContext(ctx))
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

func blogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := "http://blog:9052/blogs"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("request construction of %s failed: %s", url, err)
			http.Error(w, "server error", 500)
			return
		}
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
		io.Copy(w, res.Body)
	}
}

func ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

func main() {
	router := httprouter.New()
	router.Handler("GET", "/ping", ping())
	router.Handler("GET", "/blogs", auth(blogs()))
	//router.Handler("GET", "/blog/:title", auth())
	log.Fatal(http.ListenAndServe("0.0.0.0:9050", router))
}
