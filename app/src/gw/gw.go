package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/karlpokus/srv"
	"go.elastic.co/apm/module/apmhttp"
)

var Stderr = log.New(os.Stderr, "gw ", log.Ldate|log.Ltime)

func conf(s *srv.Server) error {
	client := apmhttp.WrapClient(&http.Client{
		Timeout: 3 * time.Second, // tcp ttl
	})
	s.Router = Routes(client, os.Getenv("LOG_REQUESTS"))
	s.Logger = Stderr
	s.Host = "0.0.0.0"
	s.Port = "9050"
	return nil
}

func main() {
	s, err := srv.New(conf)
	if err != nil {
		Stderr.Fatal(err)
	}
	err = s.Start()
	if err != nil {
		Stderr.Fatal(err)
	}
	Stderr.Println("main uxited")
}
