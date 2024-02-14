package storage

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
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
		{
		name: "First test. Create Gauge metric",
		ms: NewMemStorage(),
		args: args{
			name: "SomeGaugeMetric",
			newValue: 345.12424,
		},
	},
	{
		name: "Second test. Update Gauge metric",
		ms: &MemStorage{
			gaugeMetrics: map[string]float64{
				"SomeGaugeMetric": 1,
			},
			counterMetrics: make(map[string]int64),
		},
		args: args{
			name: "SomeGaugeMetric",
			newValue: 345.12424,
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.UpdateGaugeMetric(tt.args.name, tt.args.newValue)
			assert.Equal(t,tt.ms.ListAllGaugeMetrics()[tt.args.name],tt.args.newValue)
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
		{
			name: "First test. Create counter metric",
			ms: NewMemStorage(),
			args: args{
				name: "SomeCounterMetric",
				newValue: 10,
			},
		},
		{
			name: "Second test. Update cauge metric",
			ms: &MemStorage{
				gaugeMetrics: make(map[string]float64),
				counterMetrics: map[string]int64{
					"SomeCounterMetric": 5,
				},
			},
			args: args{
				name: "SomeCounterMetric",
				newValue: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lastCounterValue := tt.ms.ListAllCounterMetrics()[tt.args.name]
			tt.ms.UpdateCounterMetric(tt.args.name, tt.args.newValue)
			assert.Equal(t,tt.ms.ListAllCounterMetrics()[tt.args.name],tt.args.newValue + lastCounterValue)
		})
	}
}

func TestMemStorage_ListAllGaugeMetrics(t *testing.T) {
	tests := []struct {
		name string
		ms   *MemStorage
		want map[string]float64
	}{
		{
			name: "First test. Get empty list",
			ms: NewMemStorage(),
			want: make(map[string]float64),
		},
		{
			name: "Second test. Get not empty list",
			ms: &MemStorage{
				gaugeMetrics: map[string]float64{
					"SomeGaugeMetric": 5,
				},
				counterMetrics: make(map[string]int64),
			},
			want: map[string]float64{
					"SomeGaugeMetric": 5,
				},
		},
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
		{
			name: "First test. Get empty list",
			ms: NewMemStorage(),
			want: make(map[string]int64),
		},
		{
			name: "Second test. Get not empty list",
			ms: &MemStorage{
				gaugeMetrics: make(map[string]float64),
				counterMetrics: map[string]int64{
					"SomeCounterMetric": 5,
				},
			},
			want: map[string]int64{
					"SomeCounterMetric": 5,
				},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ms.ListAllCounterMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.ListAllCounterMetrics() = %v, want %v", got, tt.want)
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
		{
			name: "First test. Set counter metric",
			ms: NewMemStorage(),
			args: args{
				name: "SomeMetric",
				newValue: 6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.UpdateCounterMetric(tt.args.name, tt.args.newValue)
			assert.Equal(t,tt.ms.ListAllCounterMetrics()[tt.args.name],tt.args.newValue)
		})
	}
}
