package main

import (
	"hms_patient_mgmt_svc/api"
	// "hms_patient_mgmt_svc/metrics"
	// "github.com/penglongli/gin-metrics/ginmetrics"
	// "github.com/joho/godotenv"
)

func main() {
	// var appMonitor *ginmetrics.Monitor = metrics.GetMetricsInstance()
	// err := godotenv.Load()
	// if err != nil {

	// }

	// go func(applicationMonitor *ginmetrics.Monitor) {
	// 	metricRouter := metrics.RunMetricsServer(applicationMonitor)
	// 	metricRouter.Run(":8081")
	// }(appMonitor)

	// api.RunAppServer(appMonitor)
	api.RunAppServer()

}
