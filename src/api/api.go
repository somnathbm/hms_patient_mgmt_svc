package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"hms_patient_mgmt_svc/db"
	"hms_patient_mgmt_svc/metrics"

	"github.com/gin-gonic/gin"
	// "github.com/penglongli/gin-metrics/ginmetrics"

	"hms_patient_mgmt_svc/models"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const appName = "hms-pm"

var (
	tracer = otel.Tracer(appName)
	meter  = otel.Meter(appName)
	logger = otelslog.NewLogger(appName)
	// patientCountMetric metric.Int64Counter
)

// func RunAppServer(appMonitor *ginmetrics.Monitor) *gin.Engine {
func RunAppServer() {
	// Handle SIGINT (CTRL + C) properly
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Setup Opentelemetry
	otelShutdown, err := SetupOTelSDK(ctx)
	if err != nil {
		// return here
	}
	// Handle shutdown properly so that nothing leaks
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	appRouter := gin.Default()
	appRouter.Use(otelgin.Middleware("hms-pm-svc"))

	// initialize custom metrics
	allMetrics := metrics.GetAllMetrics()

	// for service liveness check
	appRouter.GET("/pm/healthy", func(c *gin.Context) {
		// fire off the tracer
		ctx, span := tracer.Start(c.Request.Context(), "/pm/healthy", trace.WithAttributes(attribute.String("message", "OK!!")))
		defer span.End()

		// set log
		logger.InfoContext(ctx, "service is healthy", "pm-logger", true)

		// no metrics
		// send off the response
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
		patientNum := len(result)

		// fire off the tracer
		ctx, span := tracer.Start(c.Request.Context(), "/pm/patients")
		defer span.End()

		// set log
		logger.InfoContext(ctx, "patient count", "pm-logger", patientNum)

		// set metrics
		patientCountAttr := attribute.Int("patient.total", patientNum)
		span.SetAttributes(patientCountAttr)
		allMetrics["PatientCountMetric"].Add(ctx, 1, metric.WithAttributes(patientCountAttr))

		// appMonitor.GetMetric("hms_patient_mgmt_patients_total").Add([]string{}, float64(patientNum))
		c.JSON(http.StatusOK, gin.H{
			"data": result,
		})
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

	// Wait for interruption.
	select {
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}
}
