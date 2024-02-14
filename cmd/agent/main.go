package main

import (
	"fmt"

	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/agent/utils"
	//"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
)

func main() {
	ms := storage.NewMemStorage()
	cnfg := parseFlags()

	fmt.Println(cnfg.Scheme())
	fmt.Println(cnfg.SrvAddr())
	fmt.Println(cnfg.SrvPort())
	fmt.Println(cnfg.PollInterval())
	fmt.Println(cnfg.ReportInterval())


	utils.StartMetricsMonitor(ms,cnfg)


	//cnfg.Init("localhost", 8080,2,10)
}
