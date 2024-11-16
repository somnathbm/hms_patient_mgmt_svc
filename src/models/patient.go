package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PatientInfo struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Basic_info   primitive.M        `json:"basic_info" bson:"basic_info"`
	Medical_info primitive.M        `json:"medical_info" bson:"medical_info"`
}
