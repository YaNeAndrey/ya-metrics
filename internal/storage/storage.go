package storage


type MemStorage struct {
	gauge map[string]float64
	counter map[string]int64
}

func (ms *MemStorage) Init() {
	ms.gauge = make(map[string]float64)
	ms.counter = make(map[string]int64)
}

func (ms *MemStorage) ChangeGaugeValue(name string, newValue float64) {
	ms.gauge[name] = newValue

}

func (ms *MemStorage) ChangeCounterValue(name string, newValue int64) {
	_, isExist := ms.counter[name] 
	if isExist {
		ms.counter[name] += newValue
	} else {
		ms.counter[name] = newValue
	}
}