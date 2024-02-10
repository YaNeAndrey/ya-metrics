package handlers

import (
    "net/http"
	"strconv"
	"strings"
	"html/template"
	"fmt"

	
	"github.com/go-chi/chi/v5"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
)


const tplStr = `<table>
    <thead>
        <tr>
            <th>Metric Name</th>
            <th>Metric Value</th>
        </tr>
    </thead>
    <tbody>
        {{range $key, $value := . }}
			<tr>
                <td>{{ $key }}</td>
                <td>{{ $value }}</td>
            </tr>
        {{ end }}
    </tbody>
</table>`


func HandleGetRoot(w http.ResponseWriter, r *http.Request,ms *storage.MemStorage){
	
	bufMetricMap := make(map[string]string)

	for key, value := range(ms.ListAllGaugeMetrics()){
		bufMetricMap[key] = fmt.Sprintf("%v",value)
	}
	for key, value := range(ms.ListAllCounterMetrics()){
		bufMetricMap[key] = fmt.Sprintf("%v",value)
	}

	tpl, err := template.New("table").Parse(tplStr)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
		return
    }

    err = tpl.Execute(w, bufMetricMap)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
		return
    }

	w.Header().Set("Content-Type", "text/html")
}

func HandleGetMetricValue(w http.ResponseWriter, r *http.Request,ms *storage.MemStorage){
	metricType := strings.ToLower(chi.URLParam(r, "type"))
	metricName := chi.URLParam(r, "name")

	if metricType == "" || metricName == ""{
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(""))
		return
	}


	body := ""
	statusCode := http.StatusBadRequest

	switch metricType {
	case "gauge": 
		body,statusCode = getGaugeMetricValue(metricName,ms)
	case "counter":{
		body,statusCode = getCounterMetricValue(metricName,ms)
	}
	default:
		body,statusCode = "", http.StatusNotFound
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(body))

	
}

func HandlePostUpdateMetricValue(w http.ResponseWriter, r *http.Request,ms *storage.MemStorage){
	metricType := strings.ToLower(chi.URLParam(r, "type"))
	metricName := chi.URLParam(r, "name")
	metricValueStr := chi.URLParam(r, "value")

	if metricType == "" || metricName == "" || metricValueStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	
	statusCode := updateMetric(metricType, metricName,metricValueStr,ms)
	w.WriteHeader(statusCode)
}


func getGaugeMetricValue(metricName string, ms *storage.MemStorage) (string, int) {
	allGaugeMetrics := ms.ListAllGaugeMetrics()

	value, isExist := allGaugeMetrics[metricName]
	if isExist {
		valueStr := strconv.FormatFloat(value, 'f', 0, 64)
		return valueStr, http.StatusOK
	} else {
		return "", http.StatusNotFound
	}
}

func getCounterMetricValue(metricName string, ms *storage.MemStorage) (string, int) {
	allCounterMetrics := ms.ListAllCounterMetrics()
	value,isExist := allCounterMetrics[metricName]
	if isExist {
		valueStr := strconv.FormatInt(value, 10)
			return valueStr, http.StatusOK
	} else {
		return "", http.StatusNotFound
	}
}

func updateMetric(metricType string, metricName string,metricValueStr string, ms *storage.MemStorage) int {
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