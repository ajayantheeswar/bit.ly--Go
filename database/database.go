package database

import (
	"context"
	"time"

	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client - H
var Client *mongo.Client

var Users , Visits , Links *mongo.Collection

//ConnectDatabase - Connect to Remote Database
func ConnectDatabase (){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@cluster0-shard-00-00-krkj7.mongodb.net:27017,cluster0-shard-00-01-krkj7.mongodb.net:27017,cluster0-shard-00-02-krkj7.mongodb.net:27017/bitly?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin&retryWrites=true&w=majority"))

	Client = client
	
	err = client.Ping(ctx, readpref.Primary())

	if err != nil { 
		panic(err) 
	} else {
		Users = client.Database("bitly").Collection("users")
		Visits = client.Database("bitly").Collection("visits")
		Links = client.Database("bitly").Collection("links")
	}
}

