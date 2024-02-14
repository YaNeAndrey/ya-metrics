package storage

type StorageRepo interface {
	UpdateGaugeMetric(name string, newValue float64)
	UpdateCounterMetric(name string, newValue int64)
	SetCounterMetric(name string, newValue int64)
	ListAllGaugeMetrics() map[string]float64
	ListAllCounterMetrics() map[string]int64
}