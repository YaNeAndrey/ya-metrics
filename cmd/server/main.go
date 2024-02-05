package main

import (
    "net/http"
	"strings"
	"github.com/YaNeAndrey/ya-metrics/internal/server/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/server/handlers"
)

func main() {
	var ms storage.MemStorage
	ms.Init()

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
}