package db

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"hms_patient_mgmt_svc/models"
	"hms_patient_mgmt_svc/utils"
)

// Get mongo collection for 1 minute connection pool
func get_db_collection() (*mongo.Collection, *mongo.Client) {
	_, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var base64_encoded_data = map[string]string{
		"APP_DB_URI":      os.Getenv("APP_DB_URI"),
		"APP_DB_NAME":     os.Getenv("APP_DB_NAME"),
		"COLLECTION_NAME": os.Getenv("COLLECTION_NAME"),
	}

	base64_decoded_data, base64_err := utils.DecodeBase64(base64_encoded_data)
	if base64_err != nil {
		// fmt.Println("base64 decode error", base64_err.Error())
	}

	client, err := mongo.Connect(options.Client().ApplyURI(base64_decoded_data["APP_DB_URI"]))
	if err != nil {
		panic(err)
	}
	collection := client.Database(base64_decoded_data["APP_DB_NAME"]).Collection(base64_decoded_data["COLLECTION_NAME"])
	return collection, client
}

// Get patient information by phone number
func GetPatientInfoByPhone(phone_num string) (primitive.M, error) {
	collection, client := get_db_collection()
	if collection == nil || client == nil {
		// panic("unable to proceed operation with mongo")
		return nil, errors.New("could not get db")
	}
	var result bson.M
	err := collection.FindOne(context.TODO(), bson.M{"basic_info.phone": phone_num}, options.FindOne().SetProjection(bson.M{"_id": 0})).Decode(&result)
	if err != nil && err == mongo.ErrNoDocuments {
		// handle error
		return nil, err
	}
	client.Disconnect(context.Background())
	return result, nil
}

// Create new patient that includes a unique patient ID
func CreateNewPatient(patient_info models.PatientInfo) (*string, error) {
	collection, client := get_db_collection()
	if collection == nil || client == nil {
		// panic("unable to proceed operation with mongo")
		return nil, errors.New("could not get db")
	}

	// generate a unique patient ID
	patientId := uuid.New().String()

	// insert the patientId into the patient medical info struct
	patient_info.Medical_info.PatientId = patientId

	// insert the data
	_, err := collection.InsertOne(context.TODO(), patient_info)
	if err != nil {
		// fmt.Println(err.Error())
		client.Disconnect(context.Background())
		return nil, err
	}
	// fmt.Println(result.InsertedID)
	// close the connection
	client.Disconnect(context.Background())

	// return the data back
	return &patientId, nil
}
