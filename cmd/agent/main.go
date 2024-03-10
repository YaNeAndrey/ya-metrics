package main

import (
	"github.com/YaNeAndrey/ya-metrics/internal/agent/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/storage/storagejson"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)
	testMetrics := []storage.Metrics{}

	st := storage.StorageRepo(storagejson.NewMemStorageJSON(testMetrics))

	conf := parseFlags()
	log.Printf((*conf).String())
	utils.StartMetricsMonitor(&st, conf)
}
