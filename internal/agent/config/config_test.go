package config

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetAllFields(t *testing.T) {
	type args struct {
		scheme         string
		srvAddr        string
		srvPort        int
		pollInterval   int
		reportInterval int
	}
	tests := []struct {
		name string
		c    *Config
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetAllFields(tt.args.scheme, tt.args.srvAddr, tt.args.srvPort, tt.args.pollInterval, tt.args.reportInterval)
		})
	}
}

func TestConfig_Scheme(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Scheme(); got != tt.want {
				t.Errorf("Config.Scheme() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SrvAddr(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.SrvAddr(); got != tt.want {
				t.Errorf("Config.SrvAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SrvPort(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.SrvPort(); got != tt.want {
				t.Errorf("Config.SrvPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_PollInterval(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.PollInterval(); got != tt.want {
				t.Errorf("Config.PollInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_ReportInterval(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.ReportInterval(); got != tt.want {
				t.Errorf("Config.ReportInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
