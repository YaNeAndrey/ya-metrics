package storagejson

import (
	"context"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"reflect"
	"testing"
)

func TestMemStorageJSON_GetAllMetrics(t *testing.T) {
	floatValue := float64(124.2345)
	intValue := int64(124)

	type fields struct {
		allMetrics []storage.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		want   []storage.Metrics
	}{
		{
			name: "First test",
			fields: fields{
				allMetrics: []storage.Metrics{
					{
						ID:    "GaugeMetric",
						MType: constants.GaugeMetricType,
						Value: &floatValue,
					},
					{
						ID:    "CounterMetric",
						MType: constants.CounterMetricType,
						Delta: &intValue,
					},
				},
			},
			want: []storage.Metrics{
				{
					ID:    "GaugeMetric",
					MType: constants.GaugeMetricType,
					Value: &floatValue,
				},
				{
					ID:    "CounterMetric",
					MType: constants.CounterMetricType,
					Delta: &intValue,
				},
			},
		},
	}
	myContext := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorageJSON{
				allMetrics: tt.fields.allMetrics,
			}
			metrics, _ := ms.GetAllMetrics(myContext)
			if !reflect.DeepEqual(metrics, tt.want) {
				t.Errorf("GetAllMetrics() = %v, want %v", metrics, tt.want)
			}
		})
	}
}

func TestMemStorageJSON_GetMetricByNameAndType(t *testing.T) {
	type fields struct {
		allMetrics []storage.Metrics
	}
	type args struct {
		metricName string
		metricType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    storage.Metrics
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	myContext := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorageJSON{
				allMetrics: tt.fields.allMetrics,
			}
			got, err := ms.GetMetricByNameAndType(myContext, tt.args.metricName, tt.args.metricType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetricByNameAndType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetricByNameAndType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorageJSON_UpdateMetric(t *testing.T) {
	type fields struct {
		allMetrics []storage.Metrics
	}
	type args struct {
		newMetric       storage.Metrics
		setCounterDelta bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	myContext := context.TODO()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorageJSON{
				allMetrics: tt.fields.allMetrics,
			}
			if err := ms.UpdateOneMetric(myContext, tt.args.newMetric, tt.args.setCounterDelta); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetrics_CheckMetric(t *testing.T) {
	floatValue := float64(124.2345)
	intValue := int64(124)
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "First test. Without error",
			fields: fields{
				ID:    "GaugeMetric",
				MType: constants.GaugeMetricType,
				Value: &floatValue,
			},
			wantErr: false,
		},
		{
			name: "Second test. With error",
			fields: fields{
				ID:    "GaugeMetric",
				MType: constants.GaugeMetricType,
				Value: &floatValue,
				Delta: &intValue,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &storage.Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			if err := m.CheckMetric(); (err != nil) != tt.wantErr {
				t.Errorf("CheckMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetrics_Clone(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Metrics
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &storage.Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			if got := m.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ++++
func TestMetrics_UpdateMetric(t *testing.T) {
	oldFloatValue := float64(124.2345)
	newFloatValue := float64(555.1223)
	oldIntValue := int64(124)
	newIntValue := int64(124 + 124)
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	type args struct {
		newMetric       storage.Metrics
		setCounterDelta bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    storage.Metrics
		wantErr bool
	}{
		{
			name: "First test. Correct date",
			fields: fields{
				ID:    "GaugeMetric",
				MType: constants.GaugeMetricType,
				Value: &oldFloatValue,
			},
			args: args{
				newMetric: storage.Metrics{
					ID:    "GaugeMetric",
					MType: constants.GaugeMetricType,
					Value: &newFloatValue,
				},
				setCounterDelta: false,
			},
			wantErr: false,
			want: storage.Metrics{
				ID:    "GaugeMetric",
				MType: constants.GaugeMetricType,
				Value: &newFloatValue,
			},
		},
		{
			name: "Second test. Correct date. Update counter",
			fields: fields{
				ID:    "GaugeMetric",
				MType: constants.CounterMetricType,
				Delta: &oldIntValue,
			},
			args: args{
				newMetric: storage.Metrics{
					ID:    "GaugeMetric",
					MType: constants.CounterMetricType,
					Delta: &oldIntValue,
				},
				setCounterDelta: false,
			},
			wantErr: false,
			want: storage.Metrics{
				ID:    "GaugeMetric",
				MType: constants.CounterMetricType,
				Delta: &newIntValue,
			},
		},
		{
			name: "Third test. Correct date. Set counter",
			fields: fields{
				ID:    "GaugeMetric",
				MType: constants.CounterMetricType,
				Delta: &oldIntValue,
			},
			args: args{
				newMetric: storage.Metrics{
					ID:    "GaugeMetric",
					MType: constants.CounterMetricType,
					Delta: &newIntValue,
				},
				setCounterDelta: true,
			},
			wantErr: false,
			want: storage.Metrics{
				ID:    "GaugeMetric",
				MType: constants.CounterMetricType,
				Delta: &newIntValue,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &storage.Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			if err := m.UpdateMetric(tt.args.newMetric, tt.args.setCounterDelta); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
			_ = reflect.DeepEqual(m, tt.want)
		})
	}
}

func TestNewMemStorageJSON(t *testing.T) {
	type args struct {
		allMetrics []storage.Metrics
	}
	tests := []struct {
		name string
		args args
		want *MemStorageJSON
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemStorageJSON(tt.args.allMetrics); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemStorageJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
