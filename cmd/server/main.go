package main

import (
	"fmt"
	"net/http"

	"github.com/YaNeAndrey/ya-metrics/internal/server/router"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
)

func main() {
	conf := parseFlags()
	//floatNum := 6.142434
	//intNum := int64(123456)
	testMetrics := []storage.Metrics{
		/*{
			ID:    "GaugeMetric",
			MType: constants.GaugeMetricType,
			Delta: nil,
			Value: &floatNum,
		},
		{
			ID:    "CounterMetric",
			MType: constants.CounterMetricType,
			Delta: &intNum,
			Value: nil,
		},*/
	}
	st := storage.StorageRepo(storage.NewMemStorageJSON(testMetrics))
	r := router.InitRouter(&st)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.SrvAddr(), conf.SrvPort()), r)
	if err != nil {
		panic(err)
	}
}
