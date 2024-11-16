package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Get mongo collection for 1 minute connection pool
func get_db_collection() (*mongo.Collection, *mongo.Client) {
	_, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(os.Getenv("APP_DB_URI")))
	if err != nil {
		panic(err)
	}
	collection := client.Database(os.Getenv("APP_DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))
	return collection, client
}

// Get patient information by phone number
func GetPatientInfoByPhone(phone_num string) primitive.M {
	collection, client := get_db_collection()
	if collection == nil || client == nil {
		panic("unable to proceed operation with mongo")
	}
	var result bson.M
	err := collection.FindOne(context.TODO(), bson.M{"basic_info.phone": phone_num}, options.FindOne().SetProjection(bson.M{"_id": 0})).Decode(&result)
	if err != nil {
		// handle error
		panic(err)
	}
	client.Disconnect(context.Background())
	return result
}
