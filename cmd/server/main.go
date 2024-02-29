package main

import (
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/server/router"
	"github.com/YaNeAndrey/ya-metrics/internal/server/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"log"
	"net/http"
)

func main() {
	conf := parseFlags()
	//floatNum := 6.142434
	//intNum := int64(123456)
	testMetrics := []storage.Metrics{}

	log.Println(*conf)
	st := storage.StorageRepo(storage.NewMemStorageJSON(testMetrics))

	err := utils.ReadMetricsFromFile(*conf, &st)
	if err != nil {
		log.Println(err)
	}
	r := router.InitRouter(*conf, &st)

	//Send Ctrl+C for good exit
	/*c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err := utils.SaveAllMetricsToFile(*conf, &st)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}()
	*/
	go utils.SaveMetricsByTime(*conf, &st)

	defer utils.SaveAllMetricsToFile(*conf, &st)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", conf.SrvAddr(), conf.SrvPort()), r)
	if err != nil {
		panic(err)
	}
}
