package main

import (
	"log"
	"os"

	"github.com/karlpokus/srv"
)

var Stderr = log.New(os.Stderr, "gw ", log.Ldate|log.Ltime)

func conf(s *srv.Server) error {
	s.Router = Routes(os.Getenv("LOG_REQUESTS"))
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
	Stderr.Println("main exited")
}
