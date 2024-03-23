package storage

import (
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) CheckMetric() error {
	switch m.MType {
	case constants.GaugeMetricType:
		{
			if m.Value == nil {
				return constants.ErrGaugeValue
			}
			if m.Delta != nil {
				return constants.ErrGaugeDelta
			}
		}
	case constants.CounterMetricType:
		{
			if m.Delta == nil {
				return constants.ErrCounterDelta
			}
			if m.Value != nil {
				return constants.ErrCounterValue
			}
		}
	default:
		return constants.ErrIncorectMetricType
	}
	return nil
}

func (m *Metrics) UpdateMetric(newMetric Metrics, setCounterDelta bool) error {
	err := newMetric.CheckMetric()
	if err != nil {
		return err
	}
	if m.MType != newMetric.MType {
		return constants.ErrMetricTypeDoNotMatch
	}

	switch newMetric.MType {
	case constants.CounterMetricType:
		{
			d := *(newMetric.Delta)
			if setCounterDelta {
				*m.Delta = d
			} else {
				oldDelta := *m.Delta
				newDelta := d + oldDelta
				*m.Delta = newDelta
			}
		}
	case constants.GaugeMetricType:
		{
			v := *(newMetric.Value)
			*m.Value = v
		}
	default:
		return constants.ErrIncorectMetricType
	}
	return nil
}

func (m *Metrics) Clone() Metrics {
	if m.MType == constants.GaugeMetricType {
		value := *m.Value
		return Metrics{
			ID:    m.ID,
			MType: m.MType,
			Value: &value,
		}
	} else {
		delta := *m.Delta
		return Metrics{
			ID:    m.ID,
			MType: m.MType,
			Delta: &delta,
		}
	}
}
