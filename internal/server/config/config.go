package config

import (
	"errors"
	"os"
	"path"
	"strings"
	"time"
)

type Config struct {
	srvAddr         string
	srvPort         int
	storeInterval   time.Duration
	fileStoragePath string
	restoreMetrics  bool
}

func NewConfig() *Config {
	var c Config
	c.srvAddr = "localhost"
	c.srvPort = 8080
	c.storeInterval = time.Duration(300) * time.Second
	c.fileStoragePath = path.Join(".", "tmp", "metrics-db.json")
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
	return errors.New("SrvPort must be in [1:65535]")
}

func (c *Config) SetStoreInterval(storeInterval int) error {
	if storeInterval > -1 {
		c.storeInterval = time.Duration(storeInterval) * time.Second
	}
	return errors.New("StoreInterval must be greater then -1")
}

func (c *Config) SetFileStoragePath(fileStoragePath string) error {
	_, err := os.Stat(fileStoragePath)
	if err != nil {
		if os.IsNotExist(err) {
			separatedPath := strings.Split(fileStoragePath, string(os.PathSeparator))
			dirPath := strings.Join(separatedPath[0:len(separatedPath)-1], string(os.PathSeparator))
			err = os.MkdirAll(dirPath, 0777)
			if err != nil {
				return err
			}
			_, err := os.Create(fileStoragePath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	c.fileStoragePath = fileStoragePath
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

func (c *Config) RestoreMetrics() bool {
	return c.restoreMetrics
}
