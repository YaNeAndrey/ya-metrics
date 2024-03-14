package storagejson

import (
	"context"
	"errors"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
)

type MemStorageJSON struct {
	allMetrics []storage.Metrics
}

func NewMemStorageJSON(allMetrics []storage.Metrics) *MemStorageJSON {
	return &MemStorageJSON{allMetrics: allMetrics}
}

func (ms *MemStorageJSON) UpdateOneMetric(c context.Context, newMetric storage.Metrics, setCounterDelta bool) error {
	err := newMetric.CheckMetric()
	if err != nil {
		return err
	}

	for _, metric := range ms.allMetrics {
		if metric.ID == newMetric.ID {
			return metric.UpdateMetric(newMetric, setCounterDelta)
		}
	}
	ms.allMetrics = append(ms.allMetrics, newMetric.Clone())
	return nil
}

func (ms *MemStorageJSON) GetAllMetrics(c context.Context) ([]storage.Metrics, error) {
	return ms.allMetrics, nil
}

func (ms *MemStorageJSON) GetMetricByNameAndType(c context.Context, metricName string, metricType string) (storage.Metrics, error) {
	metrics, err := ms.GetAllMetrics(c)
	if err != nil {
		return storage.Metrics{}, err
	}
	for _, metr := range metrics {
		if metr.MType == metricType && metr.ID == metricName {
			return metr, nil
		}
	}
	return storage.Metrics{}, errors.New("metric not found")
}

func (ms *MemStorageJSON) UpdateMultipleMetrics(c context.Context, newMetric []storage.Metrics) error {
	for _, metric := range newMetric {
		err := ms.UpdateOneMetric(c, metric, false)
		if err != nil {
			continue
		}
	}
	return nil
}
