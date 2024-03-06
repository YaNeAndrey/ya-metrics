package main

import (
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
	"github.com/YaNeAndrey/ya-metrics/internal/server/router"
	"github.com/YaNeAndrey/ya-metrics/internal/server/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"log"
	"net/http"
)

func main() {
	conf := parseFlags()
	testMetrics := []storage.Metrics{}

	log.Println(*conf)
	st := storage.StorageRepo(storage.NewMemStorageJSON(testMetrics))

	err := utils.ReadMetricsFromFile(*conf, &st)

	if err != nil {
		log.Println("From main: " + err.Error())
	}
	r := router.InitRouter(*conf, &st)

	err = config.CheckAndCreateFile(conf.FileStoragePath())

	if err != nil {
		log.Println("From main: " + err.Error())
	}

	go utils.SaveMetricsByTime(*conf, &st)

	defer utils.SaveAllMetricsToFile(*conf, &st)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", conf.SrvAddr(), conf.SrvPort()), r)
	if err != nil {
		panic(err)
	}
}
