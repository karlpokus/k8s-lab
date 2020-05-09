package main

import (
	"log"
	"net/http"
	"encoding/json"
	"context"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type post struct {
	Title, Body string
}

/*type store interface {
	GetAll(*[]post) error
	GetOne(string) error
	AddOne()
}*/

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

func getPosts(db *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts := db.Database("blog").Collection("posts")
		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()
		cur, err := posts.Find(ctx, bson.M{})
		if err != nil {
			log.Printf("unable to find posts: %s", err)
			http.Error(w, "server err", 500)
			return
		}
		defer cur.Close(ctx)
		var res []post // TODO: dirty paging by capping this slice
		err = cur.All(ctx, &res) // TODO: try avoid conversion and send raw json instead
		if err != nil {
			log.Printf("cursor all err: %s", err)
			http.Error(w, "server err", 500)
			return
		}
		reply(w, res)
	}
}

func getPost(db *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := httprouter.ParamsFromContext(r.Context()).ByName("title")
		posts := db.Database("blog").Collection("posts")
		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()
		var p post
		err := posts.FindOne(ctx, bson.M{"title": title}).Decode(&p)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "post not found", 404)
				return
			}
			log.Printf("Undefined error for title %s: %s", title, err)
			http.Error(w, "server err", 500)
			return
		}
		reply(w, &p)
	}
}

func addPost(db *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts := db.Database("blog").Collection("posts")
		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()
		var p post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil { // TODO: again, just insert the json without conversion
			http.Error(w, "Malformed json body", 400)
			return
		}
		_, err := posts.InsertOne(ctx, p)
		if err != nil {
			log.Printf("unable to insert post: %s", err)
			http.Error(w, "server err", 500)
			return
		}
		w.Write([]byte("ok"))
	}
}

func main() {
	db, err := NewDBClient("blog", "mongo", "27017")
	if err != nil {
		log.Fatal(err)
	}
	router := httprouter.New()
	router.Handler("GET", "/ping", ping())
	router.Handler("GET", "/posts", getPosts(db))
	router.Handler("GET", "/post/:title", getPost(db))
	router.Handler("POST", "/post", addPost(db))
	log.Fatal(http.ListenAndServe("0.0.0.0:9052", router))
}
