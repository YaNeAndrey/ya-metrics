package config

type Config struct {
	scheme string
	srvAddr string
	srvPort int
	pollInterval int //in seconds
	reportInterval int //in seconds
}

func NewConfig()(*Config) {
	var c Config
	c.scheme = "http"
	c.srvAddr = "localhost"
	c.srvPort = 8080
	c.pollInterval = 2
	c.reportInterval = 10
	return &c
}

func (c *Config) SetAllFields(scheme string,srvAddr string,srvPort int,pollInterval int,reportInterval int){
	c.scheme = scheme
	c.srvAddr = srvAddr
	c.srvPort = srvPort
	c.pollInterval = pollInterval
	c.reportInterval = reportInterval
}

func (c *Config) Scheme() string{
	return c.scheme
}

func (c *Config) SrvAddr() string{
	return c.srvAddr
}

func (c *Config) SrvPort() int{
	return c.srvPort
}

func (c *Config) PollInterval() int{
	return c.pollInterval
}

func (c *Config) ReportInterval() int{
	return c.reportInterval
}

func (c *Config) SetSrvAddr(srvAddr string) {
	c.srvAddr = srvAddr
}

func (c *Config) SetSrvPort(srvPort int) {
	c.srvPort = srvPort
}

func (c *Config) SetPollInterval(pollInterval int) {
	c.pollInterval = pollInterval
}

func (c *Config) SetReportInterval(reportInterval int) {
	c.reportInterval = reportInterval
}