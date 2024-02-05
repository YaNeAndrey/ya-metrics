package updater

import (
    "net/http"
	"strconv"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
)

func UpdateMetrics(metricType string, metricName string,metricValueStr string, ms *storage.MemStorage) int {
	switch metricType {
	case "gauge":
		return checkDataAndUpdateGauge(metricName,metricValueStr,ms)
	case "counter":
		return checkDataAndUpdateCounter(metricName,metricValueStr,ms)
	default:
		return http.StatusBadRequest
	}
}

func checkDataAndUpdateGauge(metricName string,metricValueStr string, ms *storage.MemStorage) int {
	metricValue, err := strconv.ParseFloat(metricValueStr,64) 
	if err == nil {
		ms.ChangeGaugeValue(metricName,metricValue)
		return http.StatusOK
	} else {
		return http.StatusBadRequest
	}
}

func checkDataAndUpdateCounter(metricName string,metricValueStr string, ms *storage.MemStorage) int {
	metricValue, err := strconv.ParseInt(metricValueStr, 10,64) 
	if err == nil {
		ms.ChangeCounterValue(metricName,metricValue)
		return http.StatusOK
	} else {
		return http.StatusBadRequest
	}
}