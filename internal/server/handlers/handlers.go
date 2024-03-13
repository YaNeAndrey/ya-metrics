package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
	"github.com/YaNeAndrey/ya-metrics/internal/server/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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

func HandleGetReadMetrics(w http.ResponseWriter, _ *http.Request, st *storage.StorageRepo) {
	bufMetricMap := make(map[string]string)
	w.Header().Set("Content-Type", "text/html")
	myContext := context.TODO()
	metrics, err := (*st).GetAllMetrics(myContext)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, metr := range metrics {
		if metr.MType == constants.GaugeMetricType {
			bufMetricMap[metr.ID] = fmt.Sprintf("%v", *metr.Value)
		} else {
			bufMetricMap[metr.ID] = fmt.Sprintf("%v", *metr.Delta)
		}
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

	w.WriteHeader(http.StatusOK)
}

func HandleGetPing(c config.Config, w http.ResponseWriter, _ *http.Request) {
	db, err := utils.TryToOpenDBConnection(c.DBConnectionString())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	db.Close()
	w.WriteHeader(http.StatusOK)
}

func HandleGetReadOneMetric(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	metricType := strings.ToLower(chi.URLParam(r, "type"))
	metricName := chi.URLParam(r, "name")

	body := ""

	myContext := context.TODO()
	metricInStorage, err := (*st).GetMetricByNameAndType(myContext, metricName, metricType)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if metricInStorage.MType == constants.GaugeMetricType {
		body = fmt.Sprintf("%v", *metricInStorage.Value)
	} else {
		body = fmt.Sprintf("%v", *metricInStorage.Delta)
	}
	_, err = w.Write([]byte(body))

	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func HandlePostReadOneMetricJSON(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}

	var newMetric storage.Metrics
	err := json.NewDecoder(r.Body).Decode(&newMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	myContext := context.TODO()
	metricInStorage, err := (*st).GetMetricByNameAndType(myContext, newMetric.ID, newMetric.MType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	body, err := json.Marshal(metricInStorage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandlePostUpdateOneMetric(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	metricType := strings.ToLower(chi.URLParam(r, "type"))
	metricName := chi.URLParam(r, "name")
	metricValueStr := chi.URLParam(r, "value")
	if metricType == "" || metricName == "" || metricValueStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	statusCode := updateMetric(metricType, metricName, metricValueStr, st)
	w.WriteHeader(statusCode)
}

func HandlePostUpdateOneMetricJSON(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}

	var newMetric storage.Metrics
	err := json.NewDecoder(r.Body).Decode(&newMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	myContext := context.TODO()
	err = (*st).UpdateOneMetric(myContext, newMetric, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := json.Marshal(newMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandlePostUpdateMultipleMetricsJSON(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}

	var newMetrics []storage.Metrics

	err := json.NewDecoder(r.Body).Decode(&newMetrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(newMetrics) == 0 {
		http.Error(w, "array cannot be empty", http.StatusBadRequest)
		return
	}
	myContext := context.TODO()
	err = (*st).UpdateMultipleMetrics(myContext, newMetrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandlePostMetricValueJSON(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}

	var newMetric storage.Metrics
	err := json.NewDecoder(r.Body).Decode(&newMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	myContext := context.TODO()
	metricInStorage, err := (*st).GetMetricByNameAndType(myContext, newMetric.ID, newMetric.MType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	body, err := json.Marshal(metricInStorage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateMetric(metricType string, metricName string, metricValueStr string, st *storage.StorageRepo) int {
	newMetric := storage.Metrics{
		ID:    metricName,
		MType: metricType,
		Delta: nil,
		Value: nil,
	}

	var floatValue float64
	var intValue int64
	var err error
	switch metricType {
	case constants.GaugeMetricType:
		{
			floatValue, err = strconv.ParseFloat(metricValueStr, 64)
			if err != nil {
				return http.StatusBadRequest
			}
			newMetric.Value = &floatValue
		}
	case constants.CounterMetricType:
		{
			intValue, err = strconv.ParseInt(metricValueStr, 10, 64)
			if err != nil {
				return http.StatusBadRequest
			}
			newMetric.Delta = &intValue
		}
	default:
		return http.StatusBadRequest
	}
	myContext := context.TODO()
	err = (*st).UpdateOneMetric(myContext, newMetric, false)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
