package storagejson

import (
	"context"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
)

type MemStorageJSON struct {
	allMetrics []storage.Metrics
}

func NewMemStorageJSON(allMetrics []storage.Metrics) *MemStorageJSON {
	return &MemStorageJSON{allMetrics: allMetrics}
}

func (ms *MemStorageJSON) UpdateOneMetric(ctx context.Context, newMetric storage.Metrics, setCounterDelta bool) error {
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

func (ms *MemStorageJSON) GetAllMetrics(ctx context.Context) ([]storage.Metrics, error) {
	return ms.allMetrics, nil
}

func (ms *MemStorageJSON) GetMetricByNameAndType(ctx context.Context, metricName string, metricType string) (*storage.Metrics, error) {
	metrics, err := ms.GetAllMetrics(ctx)
	if err != nil {
		return nil, err
	}
	for _, metr := range metrics {
		if metr.MType == metricType && metr.ID == metricName {
			return &metr, nil
		}
	}
	return nil, constants.ErrIncorectMetricType
}

func (ms *MemStorageJSON) UpdateMultipleMetrics(ctx context.Context, newMetric []storage.Metrics) error {
	for _, metric := range newMetric {
		err := ms.UpdateOneMetric(ctx, metric, false)
		if err != nil {
			continue
		}
	}
	return nil
}
