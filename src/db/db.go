package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"hms_patient_mgmt_svc/models"
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

// Get all patients
func GetAllPatients() ([]models.PatientInfo, error) {
	var results []models.PatientInfo
	var decodedResult []models.PatientInfo

	collection, client := get_db_collection()
	if collection == nil || client == nil {
		// panic("unable to proceed operation with mongo")
		return nil, errors.New("could not get db")
	}

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil && err == mongo.ErrNoDocuments {
		// handle error
		return nil, err
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		fmt.Println("ERROR Ocurred!!")
		panic(err)
	}
	// Prints the results of the find operation as structs
	for _, result := range results {
		cursor.Decode(&result)
		decodedResult = append(decodedResult, result)
	}
	client.Disconnect(context.Background())
	return decodedResult, nil
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
