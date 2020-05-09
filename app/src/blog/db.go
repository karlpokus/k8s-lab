package main

import (
	"fmt"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDBClient(name, host, port string) (*mongo.Client, error) {
	opts := options.Client()
	opts.SetAppName(name)
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return mongo.Connect(ctx, opts)
}
