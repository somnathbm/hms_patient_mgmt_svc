package api

import (
	"fmt"
	"log"
	"net/http"

	"hms_patient_mgmt_svc/db"

	"github.com/gin-gonic/gin"
	// "github.com/penglongli/gin-metrics/ginmetrics"

	"hms_patient_mgmt_svc/models"
)

// func RunAppServer(appMonitor *ginmetrics.Monitor) *gin.Engine {
func RunAppServer() *gin.Engine {
	appRouter := gin.Default()

	// for service liveness check
	appRouter.GET("/pm/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy!",
		})
	})

	// get all patients
	appRouter.GET("/pm/patients", func(c *gin.Context) {
		result, err := db.GetAllPatients()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":  "not found",
				"error": err.Error(),
			})
			return
		}
		// patientNum := len(result)
		// appMonitor.GetMetric("hms_patient_mgmt_patients_total").Add([]string{}, float64(patientNum))
		c.JSON(http.StatusOK, gin.H{
			"data": result,
		})
		return
	})

	// patient lookup using phone number
	appRouter.GET("/pm/patients/:phone", func(c *gin.Context) {
		phone_num := c.Param("phone")
		result, err := db.GetPatientInfoByPhone(phone_num)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":  "not found",
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": result,
		})
		return
	})

	appRouter.POST("/pm/patients", func(c *gin.Context) {
		var patientData models.PatientInfo

		// ↴ this validates the payload
		if err := c.ShouldBindJSON(&patientData); err != nil {
			log.Fatalln("content parsing error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		} else if patientData.Basic_info.Name != "" || patientData.Basic_info.Phone != "" {
			// ↴ lookup by phone number first, to get patient ID
			// ↴ If so, return the whole patient info. Otherwise, register the data and return the whole patient info
			result, _ := db.GetPatientInfoByPhone(patientData.Basic_info.Phone)
			if result != nil {
				// ↴ patient ID found. return as it is
				c.JSON(http.StatusOK, gin.H{
					"data": result,
				})
				return
			} else {
				// ↴ proceed with registering the patient
				insertResult, error := db.CreateNewPatient(patientData)
				if error != nil {
					fmt.Println("Insert operation failed")
					c.JSON(http.StatusBadRequest, gin.H{
						"data":  nil,
						"error": error.Error(),
					})
					return
				}
				// ↴ else insert the patient id into the patient info and return the data back
				patientData.Medical_info.PatientId = *insertResult

				// increment the metrics
				// appMonitor.GetMetric("hms_patient_mgmt_patients_total").Inc([]string{})

				c.JSON(http.StatusOK, gin.H{
					"data": patientData,
				})
				return
			}
		} else {
			fmt.Println("Invalid data")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid data",
			})
			return
		}
	})

	appRouter.Run()

	return appRouter
}
