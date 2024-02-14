package utils

import (
	"testing"

	
    //"github.com/stretchr/testify/assert"
	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
)

func Test_sendAllMetricsUpdates(t *testing.T) {
	type args struct {
		ms *storage.MemStorage
		c  *config.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "First test",
			args: args{
				ms: storage.NewMemStorage(),
				c: config.NewConfig(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendAllMetricsUpdates(tt.args.ms, tt.args.c)
		})
	}
}

func Test_sendOneMetricUpdate(t *testing.T) {
	type args struct {
		c         *config.Config
		metrType  string
		metrName  string
		metrValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		/*{
			name: "First test. With error",
			args: args {
				c: config.NewConfig(),
				metrType: "gauge",
				metrName: "SomeMetric",
				metrValue: "333.3",
			},
			wantErr: false,
		},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sendOneMetricUpdate(tt.args.c, tt.args.metrType, tt.args.metrName, tt.args.metrValue); (err != nil) != tt.wantErr {
				t.Errorf("sendOneMetricUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
			//err := sendOneMetricUpdate(tt.args.c, tt.args.metrType, tt.args.metrName, tt.args.metrValue)
			//if tt.wantErr {
			//	assert.EqualErrorf(t, err, "expectedErrorMsg", "Error should be: %v, got: %v", expectedErrorMsg, err)
			//}
		})
	}
}

func TestStartMetricsMonitor(t *testing.T) {
	type args struct {
		ms *storage.MemStorage
		c  *config.Config
	}
	tests := []struct {
		name string
		args args
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartMetricsMonitor(tt.args.ms, tt.args.c)
		})
	}
}


func Test_collectNewMetrics(t *testing.T) {
	type args struct {
		ms *storage.MemStorage
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "First test",
			args: args{
				ms: storage.NewMemStorage(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectNewMetrics(tt.args.ms)
		})
	}
}
