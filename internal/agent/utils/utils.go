package utils

import (
	"runtime"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
	"net/http"
	"log"
	"time"
	"math/rand"
)

func sendAllMetricsUpdates(ms *storage.MemStorage, c *config.Config){
	for metrName, metrValue := range ms.ListAllGaugeMetrics() {
		//send post for gauge metrics
		sendOneMetricUpdate(ms,c,"gauge",metrName,fmt.Sprint(metrValue))
	}
	for metrName, metrValue := range ms.ListAllCounterMetric() {
		//send post for counter metrics
		sendOneMetricUpdate(ms,c,"counter",metrName,fmt.Sprint(metrValue))
	}
	ms.SetCounterMetric("PollCount",0)
}

func sendOneMetricUpdate(ms *storage.MemStorage, c *config.Config, metrType string, metrName string, metrValue string){
	urlStr := fmt.Sprintf("%s://%s:%d/update/%s/%s/%s",c.Scheme(),c.SrvAddr(),c.SrvPort(),metrType,metrName,metrValue)
	client := &http.Client{}
    r, _ := http.NewRequest("POST", urlStr, nil)
    r.Header.Add("Content-Type", "text/plain")

    resp, err := client.Do(r)

	_ = resp
	if err != nil {
		log.Fatal(err)
	}
}

func StartMetricsMonitor(ms *storage.MemStorage, c *config.Config){
	iterCount := c.ReportInterval()/c.PollInterval()
	for{
		for i := 0; i < iterCount; i++ {
			collectNewMetrics(ms)
			time.Sleep(time.Duration(c.PollInterval()) * time.Second)
		}
		sendAllMetricsUpdates(ms,c)
	}
}

func collectNewMetrics(ms *storage.MemStorage)  {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	ms.UpdateGaugeMetric("Alloc",float64(rtm.Alloc))
	ms.UpdateGaugeMetric("BuckHashSys",float64(rtm.BuckHashSys))
	ms.UpdateGaugeMetric("Frees",float64(rtm.Frees))
	ms.UpdateGaugeMetric("GCCPUFraction",float64(rtm.GCCPUFraction))
	ms.UpdateGaugeMetric("GCSys",float64(rtm.GCSys))
	ms.UpdateGaugeMetric("HeapAlloc",float64(rtm.HeapAlloc))
	ms.UpdateGaugeMetric("HeapIdle",float64(rtm.HeapIdle))
	ms.UpdateGaugeMetric("HeapInuse",float64(rtm.HeapInuse))
	ms.UpdateGaugeMetric("HeapObjects",float64(rtm.HeapObjects))
	ms.UpdateGaugeMetric("HeapReleased",float64(rtm.HeapReleased))
	ms.UpdateGaugeMetric("HeapSys",float64(rtm.HeapSys))
	ms.UpdateGaugeMetric("LastGC",float64(rtm.LastGC))
	ms.UpdateGaugeMetric("Lookups",float64(rtm.Lookups))
	ms.UpdateGaugeMetric("MCacheInuse",float64(rtm.MCacheInuse))
	ms.UpdateGaugeMetric("MCacheSys",float64(rtm.MCacheSys))
	ms.UpdateGaugeMetric("MSpanInuse",float64(rtm.MSpanInuse))
	ms.UpdateGaugeMetric("MSpanSys",float64(rtm.MSpanSys))
	ms.UpdateGaugeMetric("Mallocs",float64(rtm.Mallocs))
	ms.UpdateGaugeMetric("NextGC",float64(rtm.NextGC))
	ms.UpdateGaugeMetric("NumForcedGC",float64(rtm.NumForcedGC))
	ms.UpdateGaugeMetric("NumGC",float64(rtm.NumGC))
	ms.UpdateGaugeMetric("OtherSys",float64(rtm.OtherSys))
	ms.UpdateGaugeMetric("PauseTotalNs",float64(rtm.PauseTotalNs))
	ms.UpdateGaugeMetric("StackInuse",float64(rtm.StackInuse))
	ms.UpdateGaugeMetric("StackSys",float64(rtm.StackSys))
	ms.UpdateGaugeMetric("Sys",float64(rtm.Sys))
	ms.UpdateGaugeMetric("TotalAlloc",float64(rtm.TotalAlloc))

	ms.UpdateCounterMetric("PollCount",1)
	
	ms.UpdateGaugeMetric("RandomValue",rand.Float64())

}