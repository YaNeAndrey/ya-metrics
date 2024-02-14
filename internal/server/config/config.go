package config

import (
	"errors"
)
type Config struct {
	srvAddr string
	srvPort int
}

func NewConfig()(*Config) {
	var c Config
	c.srvAddr = "localhost"
	c.srvPort = 8080
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

func (c *Config) SrvAddr() string{
	return c.srvAddr
}

func (c *Config) SrvPort() int{
	return c.srvPort
}