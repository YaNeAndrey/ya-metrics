package main

import (
    "net/http"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/server/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {

	ms := storage.NewMemStorage()

	ms.UpdateGaugeMetric("firstGauge", 123.25)
	ms.UpdateGaugeMetric("SecondGauge", 2.1)
	ms.UpdateCounterMetric("CounterMetric", 444)

	
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			handlers.HandleGetRoot(rw,r,ms)
		}) // rootHandle - return all metrics in html table (Change Content type to html)

		r.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}",func(rw http.ResponseWriter, r *http.Request) {
				handlers.HandleGetMetricValue(rw,r,ms)
			}) // HandleGetMetricValue - return value for metric {name}
		})

		r.Route("/update", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", func(rw http.ResponseWriter, r *http.Request) {
				handlers.HandlePostUpdateMetricValue(rw,r,ms)
			})
		})
	}) 

	err:=http.ListenAndServe(":8080", r)
	if err != nil {
        panic(err)
    }


/*
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost == r.Method {
			handlers.HandleUpdateMetrics(w,r,&ms)
		}else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

    err := http.ListenAndServe(`:8080`, mux)
    if err != nil {
        panic(err)
    }
	*/
}