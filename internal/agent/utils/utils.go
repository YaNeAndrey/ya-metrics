package utils

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"

	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"

	"github.com/shirou/gopsutil/v3/mem"
)

func sendAllMetricsUpdates(st *storage.StorageRepo, c *config.Config) {
	client := http.Client{}
	myContext := context.TODO()
	metrics, err := (*st).GetAllMetrics(myContext)
	if err != nil {
		log.Println(err)
		return
	}
	err = sendAllMetricsInOneRequest(c, metrics, &client)
	if err != nil {
		log.Println(err)
	}

	/*for _, metr := range metrics {
		err = sendOneMetricUpdate(c, metr, &client)
		if err != nil {
			log.Println(err)
		}
	}
	*/
	defaultPollInterval := int64(0)
	err = (*st).UpdateOneMetric(myContext, storage.Metrics{ID: "PollCount", MType: constants.CounterMetricType, Delta: &defaultPollInterval}, true)
	if err != nil {
		log.Println(err)
		return
	}
}

func sendAllMetricsInOneRequest(c *config.Config, metrics []storage.Metrics, client *http.Client) error {
	serverAddr := c.GetHostnameWithScheme()

	urlStr, err := url.JoinPath(serverAddr, "updates/")
	if err != nil {
		log.Println(err)
		return err
	}

	jsonDate, err := json.Marshal(metrics)
	if err != nil {
		log.Println(err)
		return err
	}
	compressedDate, err := Compress(jsonDate)
	if err != nil {
		log.Println(err)
		return err
	}

	bodyReader := bytes.NewReader(compressedDate)

	//client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bodyReader)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Encoding", "gzip")

	if c.EncryptionKey() != nil {
		hashSHA256 := generateSignature(c.EncryptionKey(), compressedDate)
		r.Header.Add("HashSHA256", base64.URLEncoding.EncodeToString(hashSHA256))
	}

	err = retry.Retry(
		func(attempt uint) error {
			resp, err := client.Do(r)
			if err != nil {
				return err
			}
			err = resp.Body.Close()
			if err != nil {
				return err
			}
			return nil
		},
		strategy.Limit(4),
		strategy.Backoff(backoff.Incremental(-1*time.Second, 2*time.Second)),
	)

	//resp, err := client.Do(r)

	if err != nil {
		return err
	}

	return nil
}

func sendOneMetricUpdate(c *config.Config, metric storage.Metrics, client *http.Client) error {
	serverAddr := c.GetHostnameWithScheme()

	urlStr, err := url.JoinPath(serverAddr, "update")
	if err != nil {
		log.Println(err)
		return err
	}

	jsonDate, err := json.Marshal(metric)
	if err != nil {
		log.Println(err)
		return err
	}
	compressedDate, err := Compress(jsonDate)
	if err != nil {
		log.Println(err)
		return err
	}

	bodyReader := bytes.NewReader(compressedDate)

	r, err := http.NewRequest("POST", urlStr, bodyReader)
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Encoding", "gzip")

	if c.EncryptionKey() != nil {
		hashSHA256 := generateSignature(c.EncryptionKey(), compressedDate)
		r.Header.Add("HashSHA256", base64.URLEncoding.EncodeToString(hashSHA256))
	}

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

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return b.Bytes(), nil
}

/*
func StartMetricsMonitor(st *storage.StorageRepo, c *config.Config) {
	iterCount := int(c.ReportInterval() / c.PollInterval())
	for {
		for i := 0; i < iterCount; i++ {
			collectNewMetrics(st)
			time.Sleep(c.PollInterval())
		}
		sendAllMetricsUpdates(st, c)
	}
}*/

func collectDefaultMetrics(metricsCh chan<- storage.Metrics, c *config.Config) {
	for {
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
			value := metricValue
			newMetric := storage.Metrics{
				ID:    metricName,
				MType: constants.GaugeMetricType,
				Value: &value,
			}
			metricsCh <- newMetric
		}
		pollInterval := int64(1)

		metricsCh <- storage.Metrics{ID: "PollCount", MType: constants.CounterMetricType, Delta: &pollInterval}
		time.Sleep(c.PollInterval())
	}
}

func collectAdditionalMetrics(metricsCh chan<- storage.Metrics, c *config.Config) {
	for {
		v, _ := mem.VirtualMemory()

		gaugeMetrics := map[string]float64{
			"TotalMemory":     float64(v.Total),
			"FreeMemory":      float64(v.Free),
			"CPUutilization1": 1.1,
		}

		for metricName, metricValue := range gaugeMetrics {
			value := metricValue
			newMetric := storage.Metrics{
				ID:    metricName,
				MType: constants.GaugeMetricType,
				Value: &value,
			}

			metricsCh <- newMetric
		}
		time.Sleep(c.PollInterval())
	}
}

func generateSignature(key []byte, date []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(date)
	return h.Sum(nil)
}

func worker(c *config.Config, jobs <-chan storage.Metrics, client *http.Client) {
	for j := range jobs {
		err := sendOneMetricUpdate(c, j, client)
		log.Println(j)
		if err != nil {
			continue
		}
	}
}

func StartMetricsMonitorWithWorkers(c *config.Config) {
	numJobs := 32 // 29 (old metrics) + 3 (new metrics)
	jobs := make(chan storage.Metrics, numJobs)
	//results := make(chan string, numJobs)

	defer close(jobs)
	client := http.Client{}
	go collectDefaultMetrics(jobs, c)
	go collectAdditionalMetrics(jobs, c)

	var wg sync.WaitGroup
	wg.Add(c.RateLimit())
	for w := 1; w <= c.RateLimit(); w++ {
		go worker(c, jobs, &client)
	}
	wg.Wait()
}
