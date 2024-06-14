package main

import (
	"github.com/YaNeAndrey/ya-metrics/internal/agent/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	log.SetReportCaller(true)
	//testMetrics := []storage.Metrics{}

	//st := storage.StorageRepo(storagejson.NewMemStorageJSON(testMetrics))

	conf := parseFlags()
	log.Printf((*conf).String())
	go http.ListenAndServe(":8001", nil)
	utils.StartMetricsMonitorWithWorkers(conf)
	//utils.StartMetricsMonitor(&st, conf)
}
