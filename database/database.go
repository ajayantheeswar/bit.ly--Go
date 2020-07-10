package database

import(
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"	
)

//Client - mongodb-client
var Client *mongo.Client

//ConnectDatabase - Connect to Remote Database
func ConnectDatabase (){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	
	defer cancel()
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
	"mongodb+srv://admin:<admin>@cluster0-krkj7.mongodb.net/<bitly>?retryWrites=true&w=majority",
	))

	Client = client
	
	if err != nil { 
		panic(err) 
	}
}

