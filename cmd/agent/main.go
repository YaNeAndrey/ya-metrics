package main

import (
	"github.com/YaNeAndrey/ya-metrics/internal/agent/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"log"
)

func main() {
	testMetrics := []storage.Metrics{}

	var st storage.StorageRepo
	st = storage.NewMemStorageJSON(testMetrics)
	cnfg := parseFlags()
	log.Println(cnfg.Scheme())
	log.Println(cnfg.SrvAddr())
	log.Println(cnfg.SrvPort())
	log.Println(cnfg.PollInterval())
	log.Println(cnfg.ReportInterval())
	utils.StartMetricsMonitor(&st, cnfg)
}
