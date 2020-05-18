package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttprouter"
	"gw/internal/remote"
)

func Routes(c remote.Client, dolog string) http.Handler {
	router := apmhttprouter.New() // wraps httprouter
	router.Handler("GET", "/ping", ping())
	router.Handler("GET", "/posts", logRequest(dolog, auth(c, getPosts(c))))
	router.Handler("GET", "/post/:title", logRequest(dolog, auth(c, getPost(c))))
	router.Handler("POST", "/post", logRequest(dolog, auth(c, addPost(c))))
	return router
}

func ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

func logRequest(doLog string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if doLog == "yes" {
			Stderr.Println(r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	}
}

func auth(c remote.Client, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pwd, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "missing auth", 400)
			return
		}
		url := "http://user:9051/user"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			errPipe(w, r, fmt.Errorf("request construction of %s failed: %s", url, err))
			return
		}
		req.SetBasicAuth(user, pwd)
		res, err := call(c, req, r.Context())
		if err != nil {
			errPipe(w, r, fmt.Errorf("remote call to %s failed: %s", url, err))
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

func getPosts(c remote.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := "http://blog:9052/posts"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			errPipe(w, r, fmt.Errorf("request construction of %s failed: %s", url, err))
			return
		}
		res, err := call(c, req, r.Context())
		reply(w, r, res, err, url)
	}
}

func getPost(c remote.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := "http://blog:9052/post/"
		title := httprouter.ParamsFromContext(r.Context()).ByName("title") // TODO: sanitize title
		if title == "" {
			http.Error(w, "missing title arg", 400)
			return
		}
		url = url + title
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			errPipe(w, r, fmt.Errorf("request construction of %s failed: %s", url, err))
			return
		}
		res, err := call(c, req, r.Context())
		reply(w, r, res, err, url)
	}
}

func addPost(c remote.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := "http://blog:9052/post"
		req, err := http.NewRequest("POST", url, r.Body)
		if err != nil {
			errPipe(w, r, fmt.Errorf("request construction of %s failed: %s", url, err))
			return
		}
		res, err := call(c, req, r.Context())
		reply(w, r, res, err, url)
	}
}

func errPipe(w http.ResponseWriter, r *http.Request, err error) {
	Stderr.Println(err)
	apm.CaptureError(r.Context(), err).Send()
	http.Error(w, "server err", 500)
}

func call(c remote.Client, req *http.Request, parent context.Context) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(parent, 3*time.Second) // http ttl
	defer cancel()
	return c.Do(req.WithContext(ctx))
}

func reply(w http.ResponseWriter, r *http.Request, res *http.Response, err error, url string) {
	if err != nil {
		errPipe(w, r, fmt.Errorf("remote call to %s failed: %s", url, err))
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		http.Error(w, res.Status, res.StatusCode)
		return
	}
	io.Copy(w, res.Body)
}
