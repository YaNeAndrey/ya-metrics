package main

import (
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/agent/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
	"fmt"
)

func main() {
	ms := storage.NewMemStorage()
	cnfg := config.NewConfig()

	fmt.Println(cnfg.Scheme())
	fmt.Println(cnfg.SrvAddr())
	fmt.Println(cnfg.SrvPort())
	fmt.Println(cnfg.PollInterval())
	fmt.Println(cnfg.ReportInterval())
	utils.StartMetricsMonitor(ms,cnfg)
	//cnfg.Init("localhost", 8080,2,10)
}
