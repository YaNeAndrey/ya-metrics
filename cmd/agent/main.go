package main

import (
	"github.com/YaNeAndrey/ya-metrics/internal/agent/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/storage/storagejson"
	"log"
)

func main() {
	testMetrics := []storage.Metrics{}

	st := storage.StorageRepo(storagejson.NewMemStorageJSON(testMetrics))

	cnfg := parseFlags()
	log.Printf("Config { Server: %s://%s:%d; Poll interval: %s; Report interval: %s", cnfg.Scheme(), cnfg.SrvAddr(), cnfg.SrvPort(), cnfg.PollInterval(), cnfg.ReportInterval())
	utils.StartMetricsMonitor(&st, cnfg)
}
