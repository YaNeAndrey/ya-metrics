package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"

	"github.com/go-chi/chi/v5"
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

func HandleGetRoot(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	bufMetricMap := make(map[string]string)

	for _, metr := range (*st).GetAllMetrics() {
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

	w.Header().Set("Content-Type", "text/html")
}

func HandleGetMetricValue(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	metricType := strings.ToLower(chi.URLParam(r, "type"))
	metricName := chi.URLParam(r, "name")

	body := ""

	metricInStorage, err := (*st).GetMetricByNameAndType(metricName, metricType)
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

	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte(body))
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

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
	metricInStorage, err := (*st).GetMetricByNameAndType(newMetric.ID, newMetric.MType)
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

func HandlePostUpdateMetricValue(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
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

func HandlePostUpdateMetricValueJSON(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Incorrect Content-Type. application/json required", http.StatusBadRequest)
	}

	var newMetric storage.Metrics
	err := json.NewDecoder(r.Body).Decode(&newMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = (*st).UpdateMetric(newMetric, false)
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
}

func updateMetric(metricType string, metricName string, metricValueStr string, st *storage.StorageRepo) int {
	newMetric := storage.Metrics{
		ID:    metricName,
		MType: metricType,
		Delta: nil,
		Value: nil,
	}

	switch metricType {
	case constants.GaugeMetricType:
		{
			metricValue, err := strconv.ParseFloat(metricValueStr, 64)
			if err != nil {
				return http.StatusBadRequest
			}
			newMetric.Value = &metricValue
		}
	case constants.CounterMetricType:
		{
			metricValue, err := strconv.ParseInt(metricValueStr, 10, 64)
			if err != nil {
				return http.StatusBadRequest
			}
			newMetric.Delta = &metricValue
		}
	default:
		return http.StatusBadRequest
	}
	err := (*st).UpdateMetric(newMetric, false)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
