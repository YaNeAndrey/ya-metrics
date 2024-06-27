// Пакет utils содержит методы для сбора и отправки метрик.
package utils

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	log "github.com/sirupsen/logrus"
	mrand "math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"

	"github.com/shirou/gopsutil/v3/mem"
)

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
	bodyDate, err := Compress(jsonDate)
	if err != nil {
		log.Println(err)
		return err
	}

	if c.ServerPubKey() != nil {
		bodyDate, err = rsa.EncryptPKCS1v15(rand.Reader, c.ServerPubKey(), bodyDate)
		if err != nil {
			return err
		}
	}

	bodyReader := bytes.NewReader(bodyDate)

	r, err := http.NewRequest("POST", urlStr, bodyReader)
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Encoding", "gzip")

	if c.EncryptionKey() != nil {
		hashSHA256 := generateSignature(c.EncryptionKey(), bodyDate)
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

func collectDefaultMetrics(ctx context.Context, metricsCh chan<- storage.Metrics, c *config.Config) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(c.PollInterval()):

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
				"RandomValue":   mrand.Float64(),
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
			//	time.Sleep(c.PollInterval())
		}
	}
}

func collectAdditionalMetrics(ctx context.Context, metricsCh chan<- storage.Metrics, c *config.Config) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(c.PollInterval()):

			v, _ := mem.VirtualMemory()
			percentage, err := cpu.Percent(0, true)
			CPUusage := float64(0)
			if err == nil {
				for _, value := range percentage {
					CPUusage += value
				}
			}

			gaugeMetrics := map[string]float64{
				"TotalMemory":     float64(v.Total),
				"FreeMemory":      float64(v.Free),
				"CPUutilization1": CPUusage,
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
		}
	}
}

func generateSignature(key []byte, date []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(date)
	return h.Sum(nil)
}

func worker(ctx context.Context, c *config.Config, jobs <-chan storage.Metrics, client *http.Client, wg *sync.WaitGroup) {

	defer wg.Done()

	for /*j := range jobs*/ {
		select {
		case <-ctx.Done():
			return

		case j := <-jobs:
			err := sendOneMetricUpdate(c, j, client)
			if err != nil {
				continue
			}
		}
	}
}

// StartMetricsMonitorWithWorkers - запускает мониторинг метрик.
func StartMetricsMonitorWithWorkers(c *config.Config) {
	numJobs := 32 // 29 (old metrics) + 3 (new metrics)
	jobs := make(chan storage.Metrics, numJobs)

	ctx, cancel := context.WithCancel(context.Background())
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	idleConnsClosed := make(chan struct{})

	client := http.Client{}
	go collectDefaultMetrics(ctx, jobs, c)
	go collectAdditionalMetrics(ctx, jobs, c)

	var wg sync.WaitGroup
	wg.Add(c.RateLimit())
	for w := 1; w <= c.RateLimit(); w++ {
		go worker(ctx, c, jobs, &client, &wg)
	}

	go func() {
		<-exit
		cancel()
		wg.Wait()
		close(jobs)
		close(idleConnsClosed)
	}()

	<-idleConnsClosed
}
