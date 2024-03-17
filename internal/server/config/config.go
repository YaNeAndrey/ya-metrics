package config

import (
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"path"
	"time"

	"github.com/YaNeAndrey/ya-metrics/internal/server/utils"
)

type Config struct {
	srvAddr            string
	srvPort            int
	storeInterval      time.Duration
	fileStoragePath    string
	dbConnectionString string
	restoreMetrics     bool
}

func NewConfig() *Config {
	var c Config
	c.srvAddr = "localhost"
	c.srvPort = 8080
	c.storeInterval = time.Duration(300) * time.Second
	c.fileStoragePath = path.Join("tmp", "metrics-db.json")
	c.dbConnectionString = ""
	c.restoreMetrics = true
	return &c
}

func (c *Config) SetSrvAddr(srvAddr string) {
	c.srvAddr = srvAddr
}

func (c *Config) SetSrvPort(srvPort int) error {
	if srvPort < 65535 && srvPort > 0 {
		c.srvPort = srvPort
		return nil
	}
	return constants.ErrIncorrectPortNumber
}

func (c *Config) SetStoreInterval(storeInterval int) error {
	if storeInterval > -1 {
		c.storeInterval = time.Duration(storeInterval) * time.Second
		return nil
	}
	return constants.ErrIncorrectStoreInterval
}

func (c *Config) SetFileStoragePath(fileStoragePath string) error {
	err := utils.CheckAndCreateFile(fileStoragePath)
	if err != nil {
		return err
	}
	c.fileStoragePath = fileStoragePath
	return nil
}

func (c *Config) SetDBConnectionString(dbConnectionString string) error {
	db, err := utils.TryToOpenDBConnection(dbConnectionString)
	if err != nil {
		return err
	}
	db.Close()
	c.dbConnectionString = dbConnectionString
	return nil
}

func (c *Config) SetRestoreMetrics(restoreMetrics bool) {
	c.restoreMetrics = restoreMetrics
}

func (c *Config) SrvAddr() string {
	return c.srvAddr
}

func (c *Config) SrvPort() int {
	return c.srvPort
}

func (c *Config) StoreInterval() time.Duration {
	return c.storeInterval
}

func (c *Config) FileStoragePath() string {
	return c.fileStoragePath
}

func (c *Config) DBConnectionString() string {
	return c.dbConnectionString
}

func (c *Config) RestoreMetrics() bool {
	return c.restoreMetrics
}

func (c *Config) String() string {
	if c.dbConnectionString != "" {
		return fmt.Sprintf("Server config: { EndPoint: %s:%d; Store interval: %s; DB connection string: %s} ", c.SrvAddr(), c.SrvPort(), c.StoreInterval(), c.DBConnectionString())

	}
	return fmt.Sprintf("Server config: { EndPoint: %s:%d; Store interval: %s; File storage: %s;} ", c.SrvAddr(), c.SrvPort(), c.StoreInterval(), c.FileStoragePath())
}
