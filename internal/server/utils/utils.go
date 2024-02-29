package utils

import (
	"encoding/json"
	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"log"
	"os"
	"time"
)

func SaveAllMetricsToFile(c config.Config, st *storage.StorageRepo) error {
	metricSlice := (*st).GetAllMetrics()
	metricsInBytes, err := json.Marshal(metricSlice)
	if err != nil {
		return err
	}

	metricFile, err := os.OpenFile(c.FileStoragePath(), os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer metricFile.Close()

	_, err = metricFile.Write(metricsInBytes)
	if err != nil {
		return err
	}
	return nil
}

func ReadMetricsFromFile(c config.Config, st *storage.StorageRepo) error {
	if c.RestoreMetrics() {
		data, err := os.ReadFile(c.FileStoragePath())
		if err != nil {
			return err
		}
		var metricsSlice []storage.Metrics
		err = json.Unmarshal(data, &metricsSlice)
		if err != nil {
			return err
		}
		newStorage := storage.StorageRepo(storage.NewMemStorageJSON(metricsSlice))
		*st = newStorage
	}
	return nil
}

func SaveMetricsByTime(c config.Config, st *storage.StorageRepo) {
	if c.StoreInterval() == 0 {
		return
	}
	for {
		err := SaveAllMetricsToFile(c, st)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(c.StoreInterval())
	}
}
