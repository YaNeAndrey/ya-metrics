package config

import (
	"fmt"
	"time"
	"errors"
)


type Config struct {
	enableTLS bool
	srvAddr string
	srvPort int
	pollInterval time.Duration //in seconds
	reportInterval time.Duration //in seconds
}

func NewConfig()(*Config) {
	var c Config
	c.SetTLS(false)
	c.SetSrvAddr("localhost")
	c.SetSrvPort(8080)
	c.SetPollInterval(2)
	c.SetReportInterval(10)
	return &c
}

func (c *Config) Scheme() string{
	scheme := "http"
	if c.enableTLS {
		scheme = "https"
	}
	return scheme
}

func (c *Config) SrvAddr() string{
	return c.srvAddr
}

func (c *Config) SrvPort() int{
	return c.srvPort
}

func (c *Config) PollInterval() time.Duration{
	return c.pollInterval
}

func (c *Config) ReportInterval() time.Duration{
	return c.reportInterval
}

func (c *Config) SetTLS(enableTLS bool) {
	c.enableTLS = enableTLS
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

func (c *Config) SetPollInterval(pollInterval int) error {
	if pollInterval > 0 {
		c.pollInterval = time.Duration(pollInterval)* time.Second
		return nil
	}
	return errors.New("pollInterval must be greater than 0") 
}

func (c *Config) SetReportInterval(reportInterval int) error {
	if reportInterval > 0 {
		c.reportInterval = time.Duration(reportInterval)* time.Second
		return nil
	}
	return errors.New("reportInterval must be greater than 0") 
}


func (c *Config) GetHostnameWithScheme() string {
	return fmt.Sprintf("%s://%s:%d",c.Scheme(),c.SrvAddr(),c.SrvPort())
}