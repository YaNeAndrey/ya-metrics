package main

import (
    "flag"
	"strings"
	"errors"
	"strconv"
	"os"

	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
)

func parseEndpoint(endpointStr string) (string, int, error){
	hp := strings.Split(endpointStr, ":")
	if len(hp) != 2 {
		return "",0,errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil{
		return "",0,err
	}
	return hp[0],port,nil
}

func parseFlags()  *config.Config {
	conf := config.NewConfig()

	srvEndpoit:= flag.String("a", "localhost:8080","Server endpoint address server:port" )
	flag.Parse()


	srvEndpointEnv, isExist := os.LookupEnv("ADDRESS")
	srvAddr,srvPort,err := "",0,errors.New("")
	if !isExist {
		srvAddr,srvPort,err = parseEndpoint(*srvEndpoit)
	} else {
		srvAddr,srvPort,err = parseEndpoint(srvEndpointEnv)
	}

	if err == nil {
		conf.SetSrvAddr(srvAddr)
		conf.SetSrvPort(srvPort)
	}

	return conf
}