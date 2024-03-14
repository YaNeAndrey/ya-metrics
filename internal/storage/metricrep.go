package storage

import "context"

type StorageRepo interface {
	UpdateOneMetric(c context.Context, newMetric Metrics, setCounterDelta bool) error
	UpdateMultipleMetrics(c context.Context, newMetrics []Metrics) error
	GetAllMetrics(c context.Context) ([]Metrics, error)
	GetMetricByNameAndType(c context.Context, metricName string, metricType string) (Metrics, error)
	/*
		UpdateGaugeMetric(name string, newValue float64)
		UpdateCounterMetric(name string, newValue int64)
		SetCounterMetric(name string, newValue int64)
		ListAllGaugeMetrics() map[string]float64
		ListAllCounterMetrics() map[string]int64
	*/
}
