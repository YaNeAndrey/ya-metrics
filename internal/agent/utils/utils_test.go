package utils

import (
	"testing"

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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendAllMetricsUpdates(tt.args.ms, tt.args.c)
		})
	}
}

func Test_sendOneMetricUpdate(t *testing.T) {
	type args struct {
		ms        *storage.MemStorage
		c         *config.Config
		metrType  string
		metrName  string
		metrValue string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendOneMetricUpdate(tt.args.ms, tt.args.c, tt.args.metrType, tt.args.metrName, tt.args.metrValue)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectNewMetrics(tt.args.ms)
		})
	}
}
