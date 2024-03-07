package storagedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"github.com/YaNeAndrey/ya-metrics/internal/server/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type StorageDB struct {
	connectionString string
}

func InitStorageDB(connectionString string) (*StorageDB, error) {
	db, err := utils.TryToOpenDBConnection(connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	myContext := context.TODO()

	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS gauge ( id serial PRIMARY KEY, name VARCHAR (50) UNIQUE NOT NULL, value double precision NOT NULL);")
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS counter ( id serial PRIMARY KEY, name VARCHAR (50) UNIQUE NOT NULL, delta bigint NOT NULL);")
	if err != nil {
		return nil, err
	}

	var resStorage StorageDB
	resStorage.connectionString = connectionString
	return &resStorage, nil
}

func (st *StorageDB) UpdateMetric(newMetric storage.Metrics, setCounterDelta bool) error {
	db, err := utils.TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	myContext := context.TODO()

	metricID, err := GetMetricIDFromBD(myContext, db, newMetric)
	if err != nil {
		return err
	}
	switch newMetric.MType {
	case constants.GaugeMetricType:
		{
			if metricID == 0 {
				_, err = InsertGaugeMetric(myContext, db, newMetric)
				if err != nil {
					return err
				}
			} else {
				_, err = db.ExecContext(myContext, "UPDATE gauge SET value = $1 WHERE id = $2;", newMetric.Value, metricID)
				if err != nil {
					return err
				}
			}
		}
	case constants.CounterMetricType:
		{
			if metricID == 0 {
				_, err = InsertCounterMetric(myContext, db, newMetric)
				if err != nil {
					return err
				}
			} else {
				if setCounterDelta {
					_, err = db.ExecContext(myContext, "UPDATE counter SET delta = $1 WHERE id = $2;", newMetric.Delta, metricID)

				} else {
					_, err = db.ExecContext(myContext, "UPDATE counter SET delta = delta+$1 WHERE id = $2;", newMetric.Delta, metricID)
				}
				if err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("incorrect metric type")
	}
	return nil
}

func (st *StorageDB) GetAllMetrics() ([]storage.Metrics, error) {
	db, err := utils.TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	myContext := context.TODO()
	var resultMetrics []storage.Metrics
	bufMetricSlice, err := SelectAllGaugeMetrics(myContext, db)
	if err != nil {
		return nil, err
	}
	resultMetrics = append(resultMetrics, bufMetricSlice...)

	bufMetricSlice, err = SelectAllCounterMetrics(myContext, db)
	if err != nil {
		return nil, err
	}
	resultMetrics = append(resultMetrics, bufMetricSlice...)

	return resultMetrics, nil
}

func (st *StorageDB) GetMetricByNameAndType(metricName string, metricType string) (storage.Metrics, error) {
	db, err := utils.TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return storage.Metrics{}, err
	}
	defer db.Close()

	myContext := context.TODO()

	var resultMetric storage.Metrics
	switch metricType {
	case constants.GaugeMetricType:
		{
			row := db.QueryRowContext(myContext, "SELECT value FROM gauge WHERE name = $1 LIMIT 1;", metricName)
			var value float64
			err := row.Scan(&value)
			if err != nil {
				return storage.Metrics{}, err
			}
			resultMetric = storage.Metrics{ID: metricName, MType: metricType, Value: &value}
		}
	case constants.CounterMetricType:
		{
			row := db.QueryRowContext(myContext, "SELECT value FROM counter WHERE name = $1 LIMIT 1;", metricName)
			var delta int64
			err := row.Scan(&delta)
			if err != nil {
				return storage.Metrics{}, err
			}
			resultMetric = storage.Metrics{ID: metricName, MType: metricType, Delta: &delta}
		}
	default:
		return storage.Metrics{}, errors.New("incorrect metric type")
	}
	return resultMetric, nil
}

func GetMetricIDFromBD(myContext context.Context, db *sql.DB, metric storage.Metrics) (int, error) {
	selectReq := fmt.Sprintf("SELECT id FROM %s WHERE name = $1 LIMIT 1;", metric.MType)

	row := db.QueryRowContext(myContext, selectReq, metric.ID)

	var id int
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
}

func InsertCounterMetric(myContext context.Context, db *sql.DB, metric storage.Metrics) (int, error) {
	row := db.QueryRowContext(myContext, "INSERT INTO counter (name, delta) VALUES ($1, $2) RETURNING id;", metric.ID, metric.Delta)
	var id int

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func InsertGaugeMetric(myContext context.Context, db *sql.DB, metric storage.Metrics) (int, error) {
	row := db.QueryRowContext(myContext, "INSERT INTO gauge (name, value) VALUES ($1, $2) RETURNING id", metric.ID, metric.Value)
	var id int

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SelectAllGaugeMetrics(myContext context.Context, db *sql.DB) ([]storage.Metrics, error) {
	rows, err := db.QueryContext(myContext, "SELECT name,value FROM gauge")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []storage.Metrics

	for rows.Next() {
		bufMetric := storage.Metrics{MType: constants.GaugeMetricType, Value: new(float64)}
		err = rows.Scan(&bufMetric.ID, bufMetric.Value)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, bufMetric.Clone())
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func SelectAllCounterMetrics(myContext context.Context, db *sql.DB) ([]storage.Metrics, error) {
	rows, err := db.QueryContext(myContext, "SELECT name,delta FROM counter")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []storage.Metrics

	for rows.Next() {
		bufMetric := storage.Metrics{MType: constants.CounterMetricType, Delta: new(int64)}
		err = rows.Scan(&bufMetric.ID, bufMetric.Delta)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, bufMetric.Clone())
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
