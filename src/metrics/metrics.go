package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const appName = "hms-pm"

var meter = otel.Meter(appName)

func patientCountMetric() (metric.Int64Counter, error) {
	patientCnt, patientCntErr := meter.Int64Counter("patient.count",
		metric.WithDescription("Total number of patients"),
	)
	if patientCntErr != nil {
		return nil, patientCntErr
	}
	return patientCnt, nil
}

func GetAllMetrics() map[string]metric.Int64Counter {
	allMetricsMap := make(map[string]metric.Int64Counter)
	patientCountMetric, pcountMetricErr := patientCountMetric()
	if pcountMetricErr != nil {
		panic(pcountMetricErr)
	}
	allMetricsMap["PatientCountMetric"] = patientCountMetric
	return allMetricsMap
}

// import (
// 	"github.com/gin-gonic/gin"
// 	// "github.com/penglongli/gin-metrics/ginmetrics"
// )

// func GetMetricsInstance() *ginmetrics.Monitor {
// 	monitor := ginmetrics.GetMonitor()
// 	return monitor
// }

// func RunMetricsServer(appMonitor *ginmetrics.Monitor) *gin.Engine {
// 	metricRouter := gin.Default()

// 	appMonitor.SetMetricPath("/metrics")

// 	// // add the metrics
// 	appMonitor.AddMetric(TotalPatientsMetric())
// 	appMonitor.AddMetric(TotalEmergencyPatients())
// 	appMonitor.AddMetric(TotalIPDPatients())
// 	appMonitor.AddMetric(TotalOPDPatients())

// 	// appMonitor.UseWithoutExposingEndpoint(appRouter)
// 	appMonitor.Expose(metricRouter)

// 	return metricRouter
// }

// // Total patients
// func TotalPatientsMetric() *ginmetrics.Metric {
// 	patientMetric := &ginmetrics.Metric{
// 		Type:        ginmetrics.Counter,
// 		Name:        "hms_patient_mgmt_patients_total",
// 		Description: "Number of total patients",
// 		// Labels:      []string{"patients_total"},
// 	}
// 	return patientMetric
// }

// // Total emergency patients
// func TotalEmergencyPatients() *ginmetrics.Metric {
// 	patientMetric := &ginmetrics.Metric{
// 		Type:        ginmetrics.Counter,
// 		Name:        "hms_patient_mgmt_patients_emg",
// 		Description: "Number of total emergency patients",
// 	}
// 	return patientMetric
// }

// // Total IPD patients
// func TotalIPDPatients() *ginmetrics.Metric {
// 	patientMetric := &ginmetrics.Metric{
// 		Type:        ginmetrics.Counter,
// 		Name:        "hms_patient_mgmt_patients_ipd",
// 		Description: "Number of total IPD (In-Patient Department) patients",
// 	}
// 	return patientMetric
// }

// // Total OPD patients
// func TotalOPDPatients() *ginmetrics.Metric {
// 	patientMetric := &ginmetrics.Metric{
// 		Type:        ginmetrics.Counter,
// 		Name:        "hms_patient_mgmt_patients_opd",
// 		Description: "Number of total OPD (Out-Patient Department) patients",
// 	}
// 	return patientMetric
// }
