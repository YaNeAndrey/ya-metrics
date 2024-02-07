package storage

type MemStorage struct {
	gaugeMetrics map[string]float64
	counterMetrics map[string]int64
}

func NewMemStorage() *MemStorage {
	var ms MemStorage
	ms.gaugeMetrics = make(map[string]float64)
	ms.counterMetrics = make(map[string]int64)
	return &ms
}

func (ms *MemStorage) UpdateGaugeMetric(name string, newValue float64) {
	ms.gaugeMetrics[name] = newValue
}

func (ms *MemStorage) UpdateCounterMetric(name string, newValue int64) {
	_, isExist := ms.counterMetrics[name] 
	if isExist {
		ms.counterMetrics[name] += newValue
	} else {
		ms.SetCounterMetric(name,newValue)
	}
}

func (ms *MemStorage) ListAllGaugeMetrics() map[string]float64{
	return ms.gaugeMetrics
}

func (ms *MemStorage) ListAllCounterMetric()map[string]int64{
	return ms.counterMetrics
}

func (ms *MemStorage) SetCounterMetric(name string, newValue int64) {
		ms.counterMetrics[name] = newValue
}