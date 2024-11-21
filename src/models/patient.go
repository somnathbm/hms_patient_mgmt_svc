package models

type PatientBasicInfo struct {
	Name    string `json:"name" bson:"name" binding:"required"`
	Sex     string `json:"sex" bson:"sex" binding:"required"`
	Age     int    `json:"age" bson:"age" binding:"required"`
	Phone   string `json:"phone" bson:"phone" binding:"required"`
	Email   string `json:"email" bson:"email" binding:"required"`
	Address string `json:"address" bson:"address" binding:"required"`
}

type PatientMedicalInfo struct {
	PatientId       string `json:"patientId" bson:"patientId,omitempty"`
	Illness_primary string `json:"illness_primary" bson:"illness_primary" binding:"required"`
	Department      string `json:"department" bson:"department" binding:"required"`
	Serious         bool   `json:"serious" bson:"serious,omitempty"`
}

type PatientInfo struct {
	// ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Basic_info   PatientBasicInfo   `json:"basic_info" bson:"basic_info" binding:"required"`
	Medical_info PatientMedicalInfo `json:"medical_info" bson:"medical_info" binding:"required"`
}
