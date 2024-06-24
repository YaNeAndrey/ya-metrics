package main

import (
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/server/router"
	"github.com/YaNeAndrey/ya-metrics/internal/server/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/storage/storagedb"
	"github.com/YaNeAndrey/ya-metrics/internal/storage/storagejson"
	log "github.com/sirupsen/logrus"
	"net/http"

	_ "net/http/pprof" // подключаем пакет pprof
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
	testMetrics := []storage.Metrics{}

	var st storage.StorageRepo
	var err error
	if conf.DBConnectionString() != "" {
		st, err = storagedb.InitStorageDB(conf.DBConnectionString())
		if err != nil {
			log.Println(err)
		}
	}

	if st == nil {
		st = storage.StorageRepo(storagejson.NewMemStorageJSON(testMetrics))

		err = utils.ReadMetricsFromFile(conf.FileStoragePath(), &st)
		if err != nil {
			log.Println(err.Error())
		}

		if conf.StoreInterval() != 0 {
			go utils.SaveMetricsByTime(conf.FileStoragePath(), conf.StoreInterval(), &st)
		}
		defer utils.SaveAllMetricsToFile(conf.FileStoragePath(), &st)
	}
	log.Printf(conf.String())
	r := router.InitRouter(*conf, &st)
	err = http.ListenAndServe(fmt.Sprintf("%s", conf.SrvAddr()), r)
	if err != nil {
		panic(err)
	}
}
