package config

import (
	"path"
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
			name: "First test. Create server config",
			want: &Config{
				srvAddr:         "localhost",
				srvPort:         8080,
				storeInterval:   time.Duration(300) * time.Second,
				fileStoragePath: path.Join(".", "tmp", "metrics-db.json"),
				restoreMetrics:  true,
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
			name: "First test. Set server hostname",
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
		name             string
		c                *Config
		args             args
		wantErr          bool
		expectedErrorMsg string
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
			wantErr:          true,
			expectedErrorMsg: "SrvPort must be in [1:65535]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SetSrvPort(tt.args.srvPort); (err != nil) != tt.wantErr {
				t.Errorf("Config.SetSrvPort() error = %v, wantErr %v", err, tt.wantErr)
			}

			err := tt.c.SetSrvPort(tt.args.srvPort)
			if tt.wantErr {
				assert.EqualErrorf(t, err, tt.expectedErrorMsg, "Error should be: %v, got: %v", tt.expectedErrorMsg, err)
			} else {
				assert.Equal(t, tt.args.srvPort, tt.c.srvPort)
			}
		})
	}
}
