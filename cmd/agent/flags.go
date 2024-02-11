package main

import (
    "flag"
	"strings"
	"errors"
	"strconv"

	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
)

func parseFlags()  *config.Config {
	conf := config.NewConfig()

	flag.Func("a", "server Address and Port", func(flagValue string) error {	
		hp := strings.Split(flagValue, ":")
		if len(hp) != 2 {
			return errors.New("need address in a form host:port")
		}
		port, err := strconv.Atoi(hp[1])
		if err != nil{
			return err
		}

		conf.SetSrvAddr(hp[0])
		conf.SetSrvPort(port)
		return nil
	})

	conf.SetReportInterval(int(*flag.Uint("r",10,"reportInterval")))
	conf.SetPollInterval(int(*flag.Uint("p",2,"pollInterval")))


    flag.Parse()
	return conf
}