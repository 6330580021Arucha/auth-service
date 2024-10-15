package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	_ "os"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserDB struct {
	DB *mongo.Client
}

func ConnectMongo() *UserDB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user_db_uri := os.Getenv("USER_DB_URI")
	clientOptions := options.Client().ApplyURI(user_db_uri)

	// connect database
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// ping database client
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("===== Connected to MongoDB! =====")
	return &UserDB{
		DB: client,
	}
}

func DisconnectMongo(client *mongo.Client) {
	// Disconnect the client
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Disconnected from MongoDB")
}
