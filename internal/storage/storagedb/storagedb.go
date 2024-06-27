package storagedb

import (
	"context"
	"database/sql"
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

func (st *StorageDB) UpdateOneMetric(ctx context.Context, newMetric storage.Metrics, setCounterDelta bool) error {
	db, err := utils.TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	metricID, err := GetMetricIDFromBD(ctx, db, newMetric)
	if err != nil {
		return err
	}
	switch newMetric.MType {
	case constants.GaugeMetricType:
		{
			if metricID == 0 { // if metric not found in DB
				_, err = InsertGaugeMetric(ctx, db, newMetric)
				if err != nil {
					return err
				}
			} else {
				_, err = db.ExecContext(ctx, "UPDATE gauge SET value = $1 WHERE id = $2;", newMetric.Value, metricID)
				if err != nil {
					return err
				}
			}
		}
	case constants.CounterMetricType:
		{
			if metricID == 0 { // if metric not found in DB
				_, err = InsertCounterMetric(ctx, db, newMetric)
				if err != nil {
					return err
				}
			} else {
				if setCounterDelta { // set metric: delta = newdelta
					_, err = db.ExecContext(ctx, "UPDATE counter SET delta = $1 WHERE id = $2;", newMetric.Delta, metricID)

				} else { // update metric: delta = oldDelta + newDelta
					_, err = db.ExecContext(ctx, "UPDATE counter SET delta = delta+$1 WHERE id = $2;", newMetric.Delta, metricID)
				}
				if err != nil {
					return err
				}
			}
		}
	default:
		return constants.ErrIncorectMetricType
	}
	return nil
}

func (st *StorageDB) GetAllMetrics(ctx context.Context) ([]storage.Metrics, error) {
	db, err := utils.TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var resultMetrics []storage.Metrics
	bufMetricSlice, err := SelectAllGaugeMetrics(ctx, db)
	if err != nil {
		return nil, err
	}
	resultMetrics = append(resultMetrics, bufMetricSlice...)

	bufMetricSlice, err = SelectAllCounterMetrics(ctx, db)
	if err != nil {
		return nil, err
	}
	resultMetrics = append(resultMetrics, bufMetricSlice...)

	return resultMetrics, nil
}

func (st *StorageDB) GetMetricByNameAndType(ctx context.Context, metricName string, metricType string) (resultMetric *storage.Metrics, err error) {
	db, err := utils.TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	//var resultMetric storage.Metrics
	switch metricType {
	case constants.GaugeMetricType:
		{
			row := db.QueryRowContext(ctx, "SELECT value FROM gauge WHERE name = $1 LIMIT 1;", metricName)
			var value float64
			err = row.Scan(&value)
			if err != nil {
				return resultMetric, err
			}
			resultMetric = &storage.Metrics{ID: metricName, MType: metricType, Value: &value}
		}
	case constants.CounterMetricType:
		{
			row := db.QueryRowContext(ctx, "SELECT delta FROM counter WHERE name = $1 LIMIT 1;", metricName)
			var delta int64
			err = row.Scan(&delta)
			if err != nil {
				return resultMetric, err
			}
			resultMetric = &storage.Metrics{ID: metricName, MType: metricType, Delta: &delta}
		}
	default:
		return nil, constants.ErrIncorectMetricType
	}
	return resultMetric, nil
}

func (st *StorageDB) UpdateMultipleMetrics(ctx context.Context, newMetrics []storage.Metrics) error {
	db, err := utils.TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, metric := range newMetrics {
		metricID, err := GetMetricIDFromBD(ctx, db, metric)
		if err != nil {
			return err
		}
		switch metric.MType {
		case constants.GaugeMetricType:
			{
				_, err = tx.ExecContext(ctx, "INSERT INTO gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = excluded.value;", metric.ID, metric.Value)
				if err != nil {
					return err
				}
			}
		case constants.CounterMetricType:
			{
				if metricID == 0 {
					_, err = tx.ExecContext(ctx, "INSERT INTO counter (name, delta) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET delta = counter.delta + $2", metric.ID, metric.Delta)

					if err != nil {
						return err

					}
				} else {
					_, err = tx.ExecContext(ctx, "UPDATE counter SET delta = delta+$1 WHERE id = $2;", metric.Delta, metricID)
					if err != nil {
						return err
					}
				}
			}
		default:
			return constants.ErrIncorectMetricType
		}
	}

	return tx.Commit()
}

func GetMetricIDFromBD(ctx context.Context, db *sql.DB, metric storage.Metrics) (int, error) {
	selectQuery := fmt.Sprintf("SELECT id FROM %s WHERE name = $1 LIMIT 1;", metric.MType)

	row := db.QueryRowContext(ctx, selectQuery, metric.ID)

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

func InsertCounterMetric(ctx context.Context, db *sql.DB, metric storage.Metrics) (int, error) {
	row := db.QueryRowContext(ctx, "INSERT INTO counter (name, delta) VALUES ($1, $2) RETURNING id;", metric.ID, metric.Delta)
	var id int

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func InsertGaugeMetric(ctx context.Context, db *sql.DB, metric storage.Metrics) (int, error) {
	row := db.QueryRowContext(ctx, "INSERT INTO gauge (name, value) VALUES ($1, $2) RETURNING id", metric.ID, metric.Value)
	var id int

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SelectAllGaugeMetrics(ctx context.Context, db *sql.DB) (metrics []storage.Metrics, err error) {
	rows, err := db.QueryContext(ctx, "SELECT name,value FROM gauge")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func SelectAllCounterMetrics(ctx context.Context, db *sql.DB) (metrics []storage.Metrics, err error) {
	rows, err := db.QueryContext(ctx, "SELECT name,delta FROM counter")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
