package main

import (
    "flag"
	"strings"
	"errors"
	"strconv"
	"os"
	//"log"
	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
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
	if port < 1 || port > 65535{
		return "",0,errors.New("port accepts values from the range [1:65535]")
	}
	return hp[0],port,nil
}

func parseFlags()  *config.Config {
	conf := config.NewConfig()

	srvEndpoit:= flag.String("a", "localhost:8080","Server endpoint address server:port" )
	reportInterval := flag.Uint("r",10,"Report Interval in seconds")
	pollInterval := flag.Uint("p",2,"Pool Interval in seconds")
    
	flag.Parse()


	srvEndpointEnv, isExist := os.LookupEnv("ADDRESS")
	srvAddr,srvPort := "",0
	var err error
	if !isExist {
		srvAddr,srvPort,err = parseEndpoint(*srvEndpoit)
	} else {
		srvAddr,srvPort,err = parseEndpoint(srvEndpointEnv)
	}

	if err == nil {
		conf.SetSrvAddr(srvAddr)
		conf.SetSrvPort(srvPort)
	}

	reportIntervalEnv, isExist := os.LookupEnv("REPORT_INTERVAL")
	if isExist {
		reportIntervalInt,err := strconv.Atoi(reportIntervalEnv)
		if err == nil {
			conf.SetReportInterval(reportIntervalInt)
		}
	} else {
		conf.SetReportInterval(int(*reportInterval))
	}
	pollIntervalEnv, isExist := os.LookupEnv("POLL_INTERVAL")
	if isExist {
		pollIntervalInt,err := strconv.Atoi(pollIntervalEnv)
		if err == nil {
			conf.SetPollInterval(pollIntervalInt)
		}
	} else {
		conf.SetPollInterval(int(*pollInterval))
	}
	return conf
}