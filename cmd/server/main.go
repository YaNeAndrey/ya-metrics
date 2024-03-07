package main

import (
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/server/router"
	"github.com/YaNeAndrey/ya-metrics/internal/server/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/storage/storagedb"
	"github.com/YaNeAndrey/ya-metrics/internal/storage/storagejson"
	"log"
	"net/http"
)

func main() {
	conf := parseFlags()
	testMetrics := []storage.Metrics{}

	log.Println(*conf)
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
			log.Println("From main: " + err.Error())
		}

		if conf.StoreInterval() != 0 {
			go utils.SaveMetricsByTime(conf.FileStoragePath(), conf.StoreInterval(), &st)
		}
		defer utils.SaveAllMetricsToFile(conf.FileStoragePath(), &st)
	}

	r := router.InitRouter(*conf, &st)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", conf.SrvAddr(), conf.SrvPort()), r)
	if err != nil {
		panic(err)
	}
}
