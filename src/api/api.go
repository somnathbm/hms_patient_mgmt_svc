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

	"hms_patient_mgmt_svc/models"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const appName = "hms-pmgmt-svc"

var (
	tracer = otel.Tracer(appName)
	logger = otelslog.NewLogger(appName)
)

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
	appRouter.Use(otelgin.Middleware("hms-pm-mgmt-svc"))

	// initialize custom metrics
	allMetrics := metrics.GetAllGaugeMetrics()
	allCounterMetrics := metrics.GetAllCounterMetrics()

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
		allMetrics["PatientCountMetric"].Record(ctx, 1, metric.WithAttributes(patientCountAttr))

		c.JSON(http.StatusOK, gin.H{
			"data": result,
		})
	})

	// patient lookup using phone number
	appRouter.GET("/pm/patients/:phone", func(c *gin.Context) {
		// fire off the tracer
		ctx, span := tracer.Start(c.Request.Context(), c.Request.RequestURI)
		defer span.End()

		phone_num := c.Param("phone")
		result, err := db.GetPatientInfoByPhone(phone_num)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":  "not found",
				"error": err.Error(),
			})
			return
		}

		// set log if the DB operation succeeds
		logger.InfoContext(ctx, "patient.info", "pm-logger", result)
		// metrics - no needed

		c.JSON(http.StatusOK, gin.H{
			"data": result,
		})
	})

	appRouter.POST("/pm/patients", func(c *gin.Context) {
		var patientData models.PatientInfo
		// fire off the tracer
		ctx, span := tracer.Start(c.Request.Context(), c.Request.RequestURI)
		defer span.End()

		// ↴ this validates the payload
		if err := c.ShouldBindJSON(&patientData); err != nil {
			log.Fatalln("content parsing error", err.Error())
			logger.ErrorContext(ctx, "patient.create.error", "pm-logger", err.Error())
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

				// set log
				logger.InfoContext(ctx, "patient.create.ok", "pm-logger", true)
				// increment the metrics
				if patientData.Medical_info.Department == "emg" {
					allCounterMetrics["TotalEmergencyPatientCounterMetric"].Add(ctx, 1)
				}

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
