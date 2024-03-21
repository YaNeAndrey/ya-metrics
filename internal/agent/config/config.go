package config

import (
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"time"
)

type Config struct {
	enableTLS      bool
	srvAddr        string
	srvPort        int
	pollInterval   time.Duration //in seconds
	reportInterval time.Duration //in seconds
	encryptionKey  []byte
}

func NewConfig() *Config {
	var c Config
	c.SetTLS(false)
	c.SetSrvAddr("localhost")
	c.SetSrvPort(8080)
	c.SetPollInterval(2)
	c.SetReportInterval(10)
	c.encryptionKey = nil
	return &c
}

func (c *Config) Scheme() string {
	scheme := "http"
	if c.enableTLS {
		scheme = "https"
	}
	return scheme
}

func (c *Config) SrvAddr() string {
	return c.srvAddr
}

func (c *Config) EncryptionKey() []byte {
	return c.encryptionKey
}

func (c *Config) SrvPort() int {
	return c.srvPort
}

func (c *Config) PollInterval() time.Duration {
	return c.pollInterval
}

func (c *Config) ReportInterval() time.Duration {
	return c.reportInterval
}

func (c *Config) SetTLS(enableTLS bool) {
	c.enableTLS = enableTLS
}

func (c *Config) SetSrvAddr(srvAddr string) {
	c.srvAddr = srvAddr
}

func (c *Config) SetEncryptionKey(encryptionKey []byte) {
	if len(encryptionKey) != 16 {
		return
	}
	c.encryptionKey = encryptionKey
}

func (c *Config) SetSrvPort(srvPort int) error {
	if srvPort < 65535 && srvPort > 0 {
		c.srvPort = srvPort
		return nil
	}
	return constants.ErrIncorrectPortNumber
}

func (c *Config) SetPollInterval(pollInterval int) error {
	if pollInterval > 0 {
		c.pollInterval = time.Duration(pollInterval) * time.Second
		return nil
	}
	return constants.ErrIncorrectPollInterval
}

func (c *Config) SetReportInterval(reportInterval int) error {
	if reportInterval > 0 {
		c.reportInterval = time.Duration(reportInterval) * time.Second
		return nil
	}
	return constants.ErrIncorrectReportInterval
}

func (c *Config) GetHostnameWithScheme() string {
	return fmt.Sprintf("%s://%s:%d", c.Scheme(), c.SrvAddr(), c.SrvPort())
}

func (c *Config) String() string {
	return fmt.Sprintf("Agent config: { Server: %s://%s:%d; Poll interval: %s; Report interval: %s} ", c.Scheme(), c.SrvAddr(), c.SrvPort(), c.PollInterval(), c.ReportInterval())
}
