package api

import (
	"net/http"

	"hms_patient_mgmt_svc/db"

	"github.com/gin-gonic/gin"
)

func Api() {
	server := gin.Default()

	server.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy!",
		})
	})

	server.GET("/patients/:phone", func(c *gin.Context) {
		phone_num := c.Param("phone")
		result := db.GetPatientInfoByPhone(phone_num)
		c.JSON(http.StatusOK, gin.H{
			"data": result,
		})
	})

	server.Run()
}
