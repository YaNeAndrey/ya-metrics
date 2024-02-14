package main

import (
    "net/http"
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/server/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	conf:= parseFlags()

 
	ms := storage.NewMemStorage()
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			handlers.HandleGetRoot(rw,r,ms)
		})
		
		r.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}",func(rw http.ResponseWriter, r *http.Request) {
				handlers.HandleGetMetricValue(rw,r,ms)
			})
		})

		r.Route("/update", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", func(rw http.ResponseWriter, r *http.Request) {
				handlers.HandlePostUpdateMetricValue(rw,r,ms)
			})
		})
	}) 
	
	err:=http.ListenAndServe(fmt.Sprintf("%s:%d",conf.SrvAddr(),conf.SrvPort()), r)
	if err != nil {
        panic(err)
    }
}