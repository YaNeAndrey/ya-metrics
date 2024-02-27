package storage

type StorageRepo interface {
	UpdateMetric(newMetric Metrics, setCounterDelta bool) error
	GetAllMetrics() []Metrics
	GetAllGaugeMetrics() []Metrics
	GetAllCounterMetrics() []Metrics
	GetMetricByNameAndType(metricName string, metricType string) (Metrics, error)
	/*
		UpdateGaugeMetric(name string, newValue float64)
		UpdateCounterMetric(name string, newValue int64)
		SetCounterMetric(name string, newValue int64)
		ListAllGaugeMetrics() map[string]float64
		ListAllCounterMetrics() map[string]int64
	*/
}
