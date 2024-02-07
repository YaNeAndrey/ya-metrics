package handlers

import (
    "net/http"
	"strconv"
	"strings"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
)

func HandleUpdateMetrics(w http.ResponseWriter, r *http.Request,ms *storage.MemStorage){
	//Content-Type: text/plain
	
	newMetricsInfo := strings.Split(r.URL.String(), "/")[2:] 
	if len(newMetricsInfo) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
/*
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
*/
	metricType := newMetricsInfo[0]
	metricName := newMetricsInfo[1]
	metricValueStr := newMetricsInfo[2]
	statusCode := updateMetrics(metricType, metricName,metricValueStr,ms)
	w.WriteHeader(statusCode)
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
		ms.UpdateGaugeMetric(metricName,metricValue)
		return http.StatusOK
	} else {
		return http.StatusBadRequest
	}
}

func checkDataAndUpdateCounter(metricName string,metricValueStr string, ms *storage.MemStorage) int {
	metricValue, err := strconv.ParseInt(metricValueStr, 10,64) 
	if err == nil {
		ms.UpdateCounterMetric(metricName,metricValue)
		return http.StatusOK
	} else {
		return http.StatusBadRequest
	}
}