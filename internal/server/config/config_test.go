package config

import (
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// +++
func TestConfig_DBConnectionString(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want string
	}{
		{
			name: "First test. Get DB connection string",
			c:    NewConfig(),
			want: "some/db:connectrion@string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.dbConnectionString = "some/db:connectrion@string"
			if got := tt.c.DBConnectionString(); got != tt.want {
				t.Errorf("Config.DBConnectionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// +++
func TestConfig_FileStoragePath(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want string
	}{
		{
			name: "First test. Get file storage path",
			c:    NewConfig(),
			want: path.Join("tmp", "metrics-db.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.FileStoragePath(); got != tt.want {
				t.Errorf("Config.FileStoragePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ++
func TestConfig_RestoreMetrics(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want bool
	}{
		{
			name: "First test. Get DB restore metric flag",
			c:    NewConfig(),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.RestoreMetrics(); got != tt.want {
				t.Errorf("Config.FileStoragePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_SetDBConnectionString(t *testing.T) {
	type fields struct {
		srvAddr            string
		srvPort            int
		storeInterval      time.Duration
		fileStoragePath    string
		dbConnectionString string
		restoreMetrics     bool
	}
	type args struct {
		dbConnectionString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				srvAddr:            tt.fields.srvAddr,
				srvPort:            tt.fields.srvPort,
				storeInterval:      tt.fields.storeInterval,
				fileStoragePath:    tt.fields.fileStoragePath,
				dbConnectionString: tt.fields.dbConnectionString,
				restoreMetrics:     tt.fields.restoreMetrics,
			}
			if err := c.SetDBConnectionString(tt.args.dbConnectionString); (err != nil) != tt.wantErr {
				t.Errorf("SetDBConnectionString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ++
func TestConfig_SetFileStoragePath(t *testing.T) {
	type args struct {
		fileStoragePath string
	}
	tests := []struct {
		name    string
		c       *Config
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SetFileStoragePath(tt.args.fileStoragePath); (err != nil) != tt.wantErr {
				t.Errorf("SetFileStoragePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// +++
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

// +++
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

// ++
func TestConfig_SetStoreInterval(t *testing.T) {
	type args struct {
		storeInterval int
	}
	tests := []struct {
		name             string
		c                *Config
		args             args
		wantErr          bool
		expectedErrorMsg string
	}{
		{
			name: "First test. Set store interval",
			c:    NewConfig(),
			args: args{
				storeInterval: 2,
			},
			wantErr: false,
		},
		{
			name: "Second test. Trying to set server port with incorrect value",
			c:    NewConfig(),
			args: args{
				storeInterval: -1,
			},
			wantErr:          true,
			expectedErrorMsg: "StoreInterval must be greater then -1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.SetStoreInterval(tt.args.storeInterval)
			if tt.wantErr {
				assert.EqualErrorf(t, err, tt.expectedErrorMsg, "Error should be: %v, got: %v", tt.expectedErrorMsg, err)
			}
		})
	}
}

// ++
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

// ++
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

// +++
func TestConfig_StoreInterval(t *testing.T) {
	tests := []struct {
		name string
		c    *Config
		want time.Duration
	}{
		{
			name: "First test. Get server port",
			c:    NewConfig(),
			want: time.Duration(300) * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.StoreInterval(); got != tt.want {
				t.Errorf("StoreInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_String(t *testing.T) {
	type fields struct {
		srvAddr            string
		srvPort            int
		storeInterval      time.Duration
		fileStoragePath    string
		dbConnectionString string
		restoreMetrics     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				srvAddr:            tt.fields.srvAddr,
				srvPort:            tt.fields.srvPort,
				storeInterval:      tt.fields.storeInterval,
				fileStoragePath:    tt.fields.fileStoragePath,
				dbConnectionString: tt.fields.dbConnectionString,
				restoreMetrics:     tt.fields.restoreMetrics,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ++
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
