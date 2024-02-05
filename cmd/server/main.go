package main

import (
    "net/http"
	"strings"
	"github.com/YaNeAndrey/ya-metrics/internal/server/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/server/updater"
)

func main() {
	var ms storage.MemStorage
	ms.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost == r.Method {
			newMetricsInfo := strings.Split(r.URL.String(), "/")[2:] 
			if len(newMetricsInfo) < 3 {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			metricType := strings.ToLower(newMetricsInfo[0])
			metricName := strings.ToLower(newMetricsInfo[1])
			metricValueStr := strings.ToLower(newMetricsInfo[2])


			statusCode := updater.UpdateMetrics(metricType, metricName,metricValueStr,&ms)
			w.WriteHeader(statusCode)
		}else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

    err := http.ListenAndServe(`:8080`, mux)
    if err != nil {
        panic(err)
    }
}