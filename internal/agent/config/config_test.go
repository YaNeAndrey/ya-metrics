package config

import (
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "First test. Create client config",
			want: &Config{
				enableTLS:      false,
				srvAddr:        "localhost",
				srvPort:        8080,
				pollInterval:   time.Duration(2) * time.Second,
				reportInterval: time.Duration(10) * time.Second,
				rateLimit:      1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Scheme(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want string
	}{
		{
			name: "First test. Get Scheme",
			c:    NewConfig(),
			want: "http",
		},
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
		{
			name: "First test. Get server hostname",
			c:    NewConfig(),
			want: "localhost",
		},
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
		{
			name: "First test. Get server port",
			c:    NewConfig(),
			want: 8080,
		},
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
		want time.Duration
	}{
		{
			name: "First test. Get pool interval",
			c:    NewConfig(),
			want: time.Duration(2) * time.Second,
		},
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
		want time.Duration
	}{
		{
			name: "First test. Get report interval",
			c:    NewConfig(),
			want: time.Duration(10) * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.ReportInterval(); got != tt.want {
				t.Errorf("Config.ReportInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetTLS(t *testing.T) {
	type args struct {
		enableTLS bool
	}
	tests := []struct {
		name string
		c    *Config
		args args
	}{
		{
			name: "First test. Set TLS",
			c:    NewConfig(),
			args: args{
				enableTLS: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetTLS(tt.args.enableTLS)
			assert.Equal(t, tt.args.enableTLS, tt.c.enableTLS)
		})
	}
}

func TestConfig_SetSrvAddr(t *testing.T) {
	type args struct {
		srvAddr string
	}
	tests := []struct {
		name string
		c    *Config
		args args
	}{
		{
			name: "First test. Set hostname",
			c:    NewConfig(),
			args: args{
				srvAddr: "1.1.1.1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetSrvAddr(tt.args.srvAddr)
			assert.Equal(t, tt.args.srvAddr, tt.c.srvAddr)
		})
	}
}

func TestConfig_SetSrvPort(t *testing.T) {
	type args struct {
		srvPort int
	}
	tests := []struct {
		name          string
		c             *Config
		args          args
		wantErr       bool
		expectedError error
	}{
		{
			name: "First test. Set server port",
			c:    NewConfig(),
			args: args{
				srvPort: 80,
			},
			wantErr: false,
		},
		{
			name: "Second test. Trying to set server port with incorrect value",
			c:    NewConfig(),
			args: args{
				srvPort: -1,
			},
			wantErr:       true,
			expectedError: constants.ErrIncorrectPortNumber,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.SetSrvPort(tt.args.srvPort)
			if tt.wantErr {
				assert.Equal(t, err, tt.expectedError, "Error should be: %v, got: %v", tt.expectedError, err)
			} else {
				assert.Equal(t, tt.args.srvPort, tt.c.srvPort)
			}
		})
	}
}

func TestConfig_SetPollInterval(t *testing.T) {
	type args struct {
		pollInterval time.Duration
	}
	tests := []struct {
		name             string
		c                *Config
		args             args
		wantErr          bool
		expectedErrorMsg string
	}{
		{
			name: "First test. Set poll interval",
			c:    NewConfig(),
			args: args{
				pollInterval: 10 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "Second test. Trying to set poll interval with incorrect value",
			c:    NewConfig(),
			args: args{
				pollInterval: -10,
			},
			wantErr:          true,
			expectedErrorMsg: "pollInterval must be greater than 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.SetPollInterval(tt.args.pollInterval)
			if tt.wantErr {
				assert.EqualErrorf(t, err, tt.expectedErrorMsg, "Error should be: %v, got: %v", tt.expectedErrorMsg, err)
			} else {
				assert.Equal(t, tt.args.pollInterval, tt.c.pollInterval)
			}
		})
	}
}

func TestConfig_SetReportInterval(t *testing.T) {
	type args struct {
		reportInterval time.Duration
	}
	tests := []struct {
		name             string
		c                *Config
		args             args
		wantErr          bool
		expectedErrorMsg string
	}{
		{
			name: "First test. Set report interval",
			c:    NewConfig(),
			args: args{
				reportInterval: 80 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "Second test. Trying to set report interval with incorrect value",
			c:    NewConfig(),
			args: args{
				reportInterval: -1,
			},
			wantErr:          true,
			expectedErrorMsg: "reportInterval must be greater than 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.SetReportInterval(tt.args.reportInterval)
			if tt.wantErr {
				assert.EqualErrorf(t, err, tt.expectedErrorMsg, "Error should be: %v, got: %v", tt.expectedErrorMsg, err)
			} else {
				assert.Equal(t, tt.args.reportInterval, tt.c.reportInterval)
			}
		})
	}
}

func TestConfig_GetHostnameWithScheme(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want string
	}{
		{
			name: "First test. Get hostaneme with scheme",
			c:    NewConfig(),
			want: "http://localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetHostnameWithScheme(); got != tt.want {
				t.Errorf("Config.GetHostnameWithScheme() = %v, want %v", got, tt.want)
			}
		})
	}
}
