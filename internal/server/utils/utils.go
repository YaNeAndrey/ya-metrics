package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/storage/storagejson"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func SaveAllMetricsToFile(filePath string, st *storage.StorageRepo) error {
	myContext := context.TODO()
	metricSlice, err := (*st).GetAllMetrics(myContext)
	if err != nil {
		return err
	}
	metricsInBytes, err := json.Marshal(metricSlice)
	if err != nil {
		return err
	}

	metricFile, err := os.OpenFile(filePath, os.O_TRUNC|os.O_RDWR, 0666)
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

func ReadMetricsFromFile(filePath string, st *storage.StorageRepo) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var metricsSlice []storage.Metrics
	err = json.Unmarshal(data, &metricsSlice)
	if err != nil {
		return err
	}
	newStorage := storage.StorageRepo(storagejson.NewMemStorageJSON(metricsSlice))
	*st = newStorage
	return nil
}

func SaveMetricsByTime(filePath string, storeInterval time.Duration, st *storage.StorageRepo) {
	if storeInterval == 0 {
		return
	}
	for {
		err := SaveAllMetricsToFile(filePath, st)
		if err != nil {
			log.Println("SaveMetricsByTime" + err.Error())
		}
		time.Sleep(storeInterval)
	}
}

func CheckAndCreateFile(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			separatedPath := strings.Split(filePath, string(os.PathSeparator))
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

func TryToOpenDBConnection(dbConnectionString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbConnectionString)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var bufError error
	err = retry.Retry(
		func(attempt uint) error {
			err = db.PingContext(ctx)
			if err != nil {
				if pgerrcode.IsConnectionException(err.Error()) {
					return err
				}
				if err, ok := err.(*pq.Error); ok {
					if err.Code == pgerrcode.UniqueViolation {
						return err
					}
				}
				bufError = err
			}
			return nil
		},
		strategy.Limit(4),
		strategy.Backoff(backoff.Incremental(-1*time.Second, 2*time.Second)),
	)

	if bufError != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
