package storage

import (
	"reflect"
	"testing"
)

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_UpdateGaugeMetric(t *testing.T) {
	type args struct {
		name     string
		newValue float64
	}
	tests := []struct {
		name string
		ms   *MemStorage
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.UpdateGaugeMetric(tt.args.name, tt.args.newValue)
		})
	}
}

func TestMemStorage_UpdateCounterMetric(t *testing.T) {
	type args struct {
		name     string
		newValue int64
	}
	tests := []struct {
		name string
		ms   *MemStorage
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.UpdateCounterMetric(tt.args.name, tt.args.newValue)
		})
	}
}

func TestMemStorage_ListAllGaugeMetrics(t *testing.T) {
	tests := []struct {
		name string
		ms   *MemStorage
		want map[string]float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ms.ListAllGaugeMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.ListAllGaugeMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_ListAllCounterMetric(t *testing.T) {
	tests := []struct {
		name string
		ms   *MemStorage
		want map[string]int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ms.ListAllCounterMetric(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.ListAllCounterMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_SetCounterMetric(t *testing.T) {
	type args struct {
		name     string
		newValue int64
	}
	tests := []struct {
		name string
		ms   *MemStorage
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.SetCounterMetric(tt.args.name, tt.args.newValue)
		})
	}
}
