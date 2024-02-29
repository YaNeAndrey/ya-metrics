package utils

import (
	"encoding/json"
	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"log"
	"os"
	"strings"
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
		err := CheckAndCreateFile(c.FileStoragePath())
		if err != nil {
			return err
		}
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

func CheckAndCreateFile(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(filePath)
			separatedPath := strings.Split(filePath, string(os.PathSeparator))
			log.Println(separatedPath)
			dirPath := strings.Join(separatedPath[0:len(separatedPath)-1], string(os.PathSeparator))
			err = os.MkdirAll(dirPath, 0666)
			if err != nil {
				return err
			}
			_, err := os.Create(filePath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
