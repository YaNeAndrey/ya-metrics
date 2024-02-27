package storage

import "errors"

type MemStorageJSON struct {
	allMetrics []Metrics
}

func NewMemStorageJSON(allMetrics []Metrics) *MemStorageJSON {
	return &MemStorageJSON{allMetrics: allMetrics}
}

func (ms *MemStorageJSON) UpdateMetric(newMetric Metrics, setCounterDelta bool) error {
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

func (ms *MemStorageJSON) GetAllMetrics() []Metrics {
	return ms.allMetrics
}

func (ms *MemStorageJSON) GetAllGaugeMetrics() []Metrics {
	return nil
}

func (ms *MemStorageJSON) GetAllCounterMetrics() []Metrics {
	return nil
}

func (ms *MemStorageJSON) GetMetricByNameAndType(metricName string, metricType string) (Metrics, error) {
	for _, metr := range ms.GetAllMetrics() {
		if metr.MType == metricType && metr.ID == metricName {
			return metr, nil
		}
	}
	return Metrics{}, errors.New("Metric not found")
}
