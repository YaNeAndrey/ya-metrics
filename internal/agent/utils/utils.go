package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
)

func sendAllMetricsUpdates(st *storage.StorageRepo, c *config.Config) {
	for _, metr := range (*st).GetAllMetrics() {
		err := sendOneMetricUpdate(c, metr)
		if err != nil {
			log.Println(err)
		}
	}
	defaultPollInterval := int64(0)
	err := (*st).UpdateMetric(storage.Metrics{ID: "PollCount", MType: constants.CounterMetricType, Delta: &defaultPollInterval}, true)
	if err != nil {
		log.Println(err)
		return
	}
}

func sendOneMetricUpdate(c *config.Config, metric storage.Metrics) error {
	serverAddr := c.GetHostnameWithScheme()

	urlStr, err := url.JoinPath(serverAddr, "update")
	if err != nil {
		log.Println(err)
		return err
	}

	jsonDate, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonDate)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bodyReader)
	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)

	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

/*
func sendOneMetricUpdate(c *config.Config, metrType string, metrName string, metrValue string) error {
	hostname := c.GetHostnameWithScheme()
	urlStr, err := url.JoinPath(hostname, "update", metrType, metrName, metrValue)
	if err != nil {
		log.Println(err)
		return err
	}

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, nil)
	r.Header.Add("Content-Type", "text/plain")

	resp, err := client.Do(r)

	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
*/

func StartMetricsMonitor(st *storage.StorageRepo, c *config.Config) {
	iterCount := int(c.ReportInterval() / c.PollInterval())
	for {
		for i := 0; i < iterCount; i++ {
			collectNewMetrics(st)
			time.Sleep(c.PollInterval())
		}
		sendAllMetricsUpdates(st, c)
	}
}

func collectNewMetrics(st *storage.StorageRepo) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	gaugeMetrics := map[string]float64{
		"Alloc":         float64(rtm.Alloc),
		"BuckHashSys":   float64(rtm.BuckHashSys),
		"Frees":         float64(rtm.Frees),
		"GCCPUFraction": float64(rtm.GCCPUFraction),
		"GCSys":         float64(rtm.GCSys),
		"HeapAlloc":     float64(rtm.HeapAlloc),
		"HeapIdle":      float64(rtm.HeapIdle),
		"HeapInuse":     float64(rtm.HeapInuse),
		"HeapObjects":   float64(rtm.HeapObjects),
		"HeapReleased":  float64(rtm.HeapReleased),
		"HeapSys":       float64(rtm.HeapSys),
		"LastGC":        float64(rtm.LastGC),
		"Lookups":       float64(rtm.Lookups),
		"MCacheInuse":   float64(rtm.MCacheInuse),
		"MCacheSys":     float64(rtm.MCacheSys),
		"MSpanInuse":    float64(rtm.MSpanInuse),
		"MSpanSys":      float64(rtm.MSpanSys),
		"Mallocs":       float64(rtm.Mallocs),
		"NextGC":        float64(rtm.NextGC),
		"NumForcedGC":   float64(rtm.NumForcedGC),
		"NumGC":         float64(rtm.NumGC),
		"OtherSys":      float64(rtm.OtherSys),
		"PauseTotalNs":  float64(rtm.PauseTotalNs),
		"StackInuse":    float64(rtm.StackInuse),
		"StackSys":      float64(rtm.StackSys),
		"Sys":           float64(rtm.Sys),
		"TotalAlloc":    float64(rtm.TotalAlloc),
		"RandomValue":   rand.Float64(),
	}

	for metricName, metricValue := range gaugeMetrics {
		newMetric := storage.Metrics{
			ID:    metricName,
			MType: constants.GaugeMetricType,
			Value: &metricValue,
		}
		//log.Println(newMetric)
		err := (*st).UpdateMetric(newMetric, false)
		if err != nil {
			log.Println(err)
			return
		}
	}
	pollInterval := int64(1)
	err := (*st).UpdateMetric(storage.Metrics{ID: "PollCount", MType: constants.CounterMetricType, Delta: &pollInterval}, false)
	if err != nil {
		log.Println(err)
		return
	}
}
