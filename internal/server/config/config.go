package config


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

func (c *Config) SetSrvPort(srvPort int) {
	c.srvPort = srvPort
}

func (c *Config) SrvAddr() string{
	return c.srvAddr
}

func (c *Config) SrvPort() int{
	return c.srvPort
}