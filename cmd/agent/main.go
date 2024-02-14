package main

import (
	"log"
	
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/agent/utils"
)

func main() {
	ms := storage.NewMemStorage()
	cnfg := parseFlags()

	log.Println(cnfg.Scheme())
	log.Println(cnfg.SrvAddr())
	log.Println(cnfg.SrvPort())
	log.Println(cnfg.PollInterval())
	log.Println(cnfg.ReportInterval())


	utils.StartMetricsMonitor(ms,cnfg)
}
