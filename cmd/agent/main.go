package main

import (
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/agent/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
)

func main() {
	ms := storage.NewMemStorage()
	cnfg := config.NewConfig()
	utils.StartMetricsMonitor(ms,cnfg)
	//cnfg.Init("localhost", 8080,2,10)
}
