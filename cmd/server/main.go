package main

import (
    "net/http"
	"strings"
	"strconv"
	
)

type MemStorage struct {
	gauge map[string]float64
	counter map[string]int64
}

func (ms *MemStorage) ChangeGaugeValue(name string, newValue float64) {
	ms.gauge[name] = newValue

}

func (ms *MemStorage) ChangeCounterValue(name string, newValue int64) {
	_, isExist := ms.counter[name] 
	if isExist {
		ms.counter[name] += newValue
	} else {
		ms.counter[name] = newValue
	}
}

func checkDataAndUpdateGauge(metricName string,mectricValueStr string, ms *MemStorage) int {
	metricValue, err := strconv.ParseFloat(mectricValueStr,64) 
	if err == nil {
		ms.ChangeGaugeValue(metricName,metricValue)
		return http.StatusOK
	} else {
		return http.StatusBadRequest
	}
}

func checkDataAndUpdateCounter(metricName string,mectricValueStr string, ms *MemStorage) int {
	metricValue, err := strconv.ParseInt(mectricValueStr, 10,64) 
	if err == nil {
		ms.ChangeCounterValue(metricName,metricValue)
		return http.StatusOK
	} else {
		return http.StatusBadRequest
	}
}

func checkDataAndUpdateMetric(metricType string, metricName string,mectricValueStr string, ms *MemStorage) int {
	if metricType == "gauge" {
		return checkDataAndUpdateGauge(metricName,mectricValueStr,ms)
	} else { // if "counter" == metricType
		return checkDataAndUpdateCounter(metricName,mectricValueStr,ms)
	}
}

func updateMetrics(r *http.Request, ms *MemStorage) int {
	if http.MethodPost == r.Method {
		newMetricsInfo := strings.Split(r.URL.String(), "/")[2:] 
		if len(newMetricsInfo) < 3 {
			return http.StatusNotFound
			
		}
		metricType := strings.ToLower(newMetricsInfo[0])
		metricName := strings.ToLower(newMetricsInfo[1])
		mectricValueStr := strings.ToLower(newMetricsInfo[2])

		if (metricType == "gauge" || metricType == "counter") {
			return checkDataAndUpdateMetric(metricType,metricName,mectricValueStr, ms)
		} else {
			return http.StatusBadRequest
		}
	} else {
		return http.StatusMethodNotAllowed
	}
	
}

func main() {
	var ms MemStorage
	ms.gauge = make(map[string]float64)
	ms.counter = make(map[string]int64)
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		statusCode := updateMetrics(r,&ms)
		w.WriteHeader(statusCode)
	})

    err := http.ListenAndServe(`:8080`, mux)
    if err != nil {
        panic(err)
    }
}