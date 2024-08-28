package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	_ "go.mongodb.org/mongo-driver/bson"
)

type dbAuth struct {
	user     string
	password string
	uri      string
}

func openMongoConnection(dba dbAuth) (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI(dba.uri).SetAuth(options.Credential{Username: dba.user, Password: dba.password})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func closeMongoConnection(client *mongo.Client) func() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}
}
