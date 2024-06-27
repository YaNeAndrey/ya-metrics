package main

import (
	"github.com/YaNeAndrey/ya-metrics/internal/agent/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {

	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)

	log.SetReportCaller(true)

	conf := parseFlags()
	log.Printf((*conf).String())
	go http.ListenAndServe(":8001", nil)

	utils.StartMetricsMonitorWithWorkers(conf)

}
