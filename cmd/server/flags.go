package main

import (
    "flag"
	"strings"
	"errors"
	"strconv"

	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
)

func parseFlags() *config.Config {

	conf := config.NewConfig()

	flag.Func("a", "server Address and Port", func(flagValue string) error {	
		hp := strings.Split(flagValue, ":")
		if len(hp) != 2 {
			return errors.New("Need address in a form host:port")
		}
		port, err := strconv.Atoi(hp[1])
		if err != nil{
			return err
		}

		conf.SetSrvAddr(hp[0])
		conf.SetSrvPort(port)
		return nil
	})

    flag.Parse()
	return conf
}