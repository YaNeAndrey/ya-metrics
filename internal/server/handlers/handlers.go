package handlers

import (
    "net/http"
	"strconv"
	"github.com/YaNeAndrey/ya-metrics/internal/server/storage"
)

func HandleUpdateMetrics(w http.ResponseWriter, r *http.Request,ms *storage.MemStorage){
	if http.MethodPost == r.Method {
		newMetricsInfo := strings.Split(r.URL.String(), "/")[2:] 
		if len(newMetricsInfo) < 3 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metricType := strings.ToLower(newMetricsInfo[0])
		metricName := strings.ToLower(newMetricsInfo[1])
		metricValueStr := strings.ToLower(newMetricsInfo[2])
		
		statusCode := updateMetrics(metricType, metricName,metricValueStr,ms)
		w.WriteHeader(statusCode)
	}else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func updateMetrics(metricType string, metricName string,metricValueStr string, ms *storage.MemStorage) int {
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