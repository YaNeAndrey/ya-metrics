package storage

import "context"

// StorageRepo - интерфейс хранилища метрик
type StorageRepo interface {
	UpdateOneMetric(ctx context.Context, newMetric Metrics, setCounterDelta bool) error
	UpdateMultipleMetrics(ctx context.Context, newMetrics []Metrics) error
	GetAllMetrics(ctx context.Context) ([]Metrics, error)
	GetMetricByNameAndType(ctx context.Context, metricName string, metricType string) (*Metrics, error)
	/*
		UpdateGaugeMetric(name string, newValue float64)
		UpdateCounterMetric(name string, newValue int64)
		SetCounterMetric(name string, newValue int64)
		ListAllGaugeMetrics() map[string]float64
		ListAllCounterMetrics() map[string]int64
	*/
}
