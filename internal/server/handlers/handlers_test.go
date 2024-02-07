package handlers

import (
	"net/http"
	"testing"

	"github.com/YaNeAndrey/ya-metrics/internal/storage"
)

func TestHandleUpdateMetrics(t *testing.T) {
	type args struct {
		w  http.ResponseWriter
		r  *http.Request
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
			HandleUpdateMetrics(tt.args.w, tt.args.r, tt.args.ms)
		})
	}
}

func Test_updateMetrics(t *testing.T) {
	type args struct {
		metricType     string
		metricName     string
		metricValueStr string
		ms             *storage.MemStorage
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateMetrics(tt.args.metricType, tt.args.metricName, tt.args.metricValueStr, tt.args.ms); got != tt.want {
				t.Errorf("updateMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkDataAndUpdateGauge(t *testing.T) {
	type args struct {
		metricName     string
		metricValueStr string
		ms             *storage.MemStorage
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkDataAndUpdateGauge(tt.args.metricName, tt.args.metricValueStr, tt.args.ms); got != tt.want {
				t.Errorf("checkDataAndUpdateGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkDataAndUpdateCounter(t *testing.T) {
	type args struct {
		metricName     string
		metricValueStr string
		ms             *storage.MemStorage
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkDataAndUpdateCounter(tt.args.metricName, tt.args.metricValueStr, tt.args.ms); got != tt.want {
				t.Errorf("checkDataAndUpdateCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}
