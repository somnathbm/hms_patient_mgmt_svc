package api

import (
	"fmt"
	"log"
	"net/http"

	"hms_patient_mgmt_svc/db"

	"github.com/gin-gonic/gin"

	"github.com/penglongli/gin-metrics/ginmetrics"

	"hms_patient_mgmt_svc/models"
)

func Api() {
	server := gin.Default()

	// get global monitor object
	monitor := ginmetrics.GetMonitor()

	// set metric path
	monitor.SetMetricPath("/metrics")

	// use the monitor
	monitor.Use(server)

	// for service liveness check
	server.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy!",
		})
	})

	// patient lookup using phone number
	server.GET("/patients/:phone", func(c *gin.Context) {
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

	server.POST("/patients", func(c *gin.Context) {
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

	server.Run()
}
