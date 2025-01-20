package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const appName = "hms-pm-mgmt-svc"

var meter = otel.Meter(appName)

// total patient gauge metric
func totalPatientsGaugeMetric() (metric.Int64Gauge, error) {
	patientCnt, patientCntErr := meter.Int64Gauge("patients.total",
		metric.WithDescription("Total number of patients"),
	)
	if patientCntErr != nil {
		return nil, patientCntErr
	}
	return patientCnt, nil
}

// Total emergency patients - gauge
func totalEmergencyPatientsGaugeMetric() (metric.Int64Gauge, error) {
	patientCnt, patientCntErr := meter.Int64Gauge("patients.emergency",
		metric.WithDescription("Total emergency patients"),
	)
	if patientCntErr != nil {
		return nil, patientCntErr
	}
	return patientCnt, nil
}

// Total emergency patients - counter
func totalEmergencyPatientsCounterMetric() (metric.Int64Counter, error) {
	patientCnt, patientCntErr := meter.Int64Counter("patients.emergency",
		metric.WithDescription("Total emergency patients"),
	)
	if patientCntErr != nil {
		return nil, patientCntErr
	}
	return patientCnt, nil
}

// Total IPD patients
func totalIPDPatientsGaugeMetric() (metric.Int64Gauge, error) {
	patientCnt, patientCntErr := meter.Int64Gauge("patients.ipd",
		metric.WithDescription("Total IPD patients"),
	)
	if patientCntErr != nil {
		return nil, patientCntErr
	}
	return patientCnt, nil
}

// Total OPD patients
func totalOPDPatientsGaugeMetric() (metric.Int64Gauge, error) {
	patientCnt, patientCntErr := meter.Int64Gauge("patients.opd",
		metric.WithDescription("Total OPD patients"),
	)
	if patientCntErr != nil {
		return nil, patientCntErr
	}
	return patientCnt, nil
}

func GetAllGaugeMetrics() map[string]metric.Int64Gauge {
	allMetricsMap := make(map[string]metric.Int64Gauge)
	totalPatientsGaugeMetric, totalPatientsGaugeMetricErr := totalPatientsGaugeMetric()
	if totalPatientsGaugeMetricErr != nil {
		panic(totalPatientsGaugeMetricErr)
	}
	totalEmergencyPatientsGaugeMetric, totalEmergencyPatientsGaugeMetricErr := totalEmergencyPatientsGaugeMetric()
	if totalEmergencyPatientsGaugeMetricErr != nil {
		panic(totalEmergencyPatientsGaugeMetricErr)
	}
	totalIPDPatientsGaugeMetric, totalIPDPatientsGaugeMetricErr := totalIPDPatientsGaugeMetric()
	if totalIPDPatientsGaugeMetricErr != nil {
		panic(totalIPDPatientsGaugeMetricErr)
	}
	totalOPDPatientsGaugeMetric, totalOPDPatientsGaugeMetricErr := totalOPDPatientsGaugeMetric()
	if totalOPDPatientsGaugeMetricErr != nil {
		panic(totalOPDPatientsGaugeMetricErr)
	}
	allMetricsMap["TotalPatientsGaugeMetric"] = totalPatientsGaugeMetric
	allMetricsMap["TotalEmergencyPatientGaugeMetric"] = totalEmergencyPatientsGaugeMetric
	allMetricsMap["totalIPDPatientsGaugeMetric"] = totalIPDPatientsGaugeMetric
	allMetricsMap["totalOPDPatientsGaugeMetric"] = totalOPDPatientsGaugeMetric

	return allMetricsMap
}

func GetAllCounterMetrics() map[string]metric.Int64Counter {
	allMetricsMap := make(map[string]metric.Int64Counter)
	totalEmergencyPatientsCounterMetric, totalEmergencyPatientsCounterMetricErr := totalEmergencyPatientsCounterMetric()
	if totalEmergencyPatientsCounterMetricErr != nil {
		panic(totalEmergencyPatientsCounterMetricErr)
	}
	// totalIPDPatientsGaugeMetric, totalIPDPatientsGaugeMetricErr := totalIPDPatientsGaugeMetric()
	// if totalIPDPatientsGaugeMetricErr != nil {
	// 	panic(totalIPDPatientsGaugeMetricErr)
	// }
	// totalOPDPatientsGaugeMetric, totalOPDPatientsGaugeMetricErr := totalOPDPatientsGaugeMetric()
	// if totalOPDPatientsGaugeMetricErr != nil {
	// 	panic(totalOPDPatientsGaugeMetricErr)
	// }
	// allMetricsMap["TotalPatientsGaugeMetric"] = totalPatientsGaugeMetric
	// allMetricsMap["TotalEmergencyPatientGaugeMetric"] = totalEmergencyPatientsGaugeMetric
	allMetricsMap["TotalEmergencyPatientCounterMetric"] = totalEmergencyPatientsCounterMetric
	// allMetricsMap["totalIPDPatientsGaugeMetric"] = totalIPDPatientsGaugeMetric
	// allMetricsMap["totalOPDPatientsGaugeMetric"] = totalOPDPatientsGaugeMetric

	return allMetricsMap
}
