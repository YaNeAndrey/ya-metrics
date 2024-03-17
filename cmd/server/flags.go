package main

import (
	"flag"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"

	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
)

func parseEndpoint(endpointStr string) (string, int, error) {
	hp := strings.Split(endpointStr, ":")
	if len(hp) != 2 {
		return "", 0, constants.ErrIncorrectEndpointFormat
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return "", 0, err
	}
	return hp[0], port, nil
}

func parseFlags() *config.Config {
	conf := config.NewConfig()

	srvEndpoit := flag.String("a", "localhost:8080", "Server endpoint address server:port")
	storeInterval := flag.Uint("i", 300, "Store Interval in seconds")
	fileStoragePath := flag.String("f", ".\\tmp\\metrics-db.json", "File storage path (.json)")
	dbConnectionString := flag.String("d", "", "dbConnectionString in Postgres format: postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]")
	restoreMetrics := flag.Bool("r", true, "Restore old metrics? (true or false)")
	flag.Parse()

	srvEndpointEnv, isExist := os.LookupEnv("ADDRESS")
	srvAddr, srvPort := "", 0
	var err error
	if !isExist {
		srvAddr, srvPort, err = parseEndpoint(*srvEndpoit)
	} else {
		srvAddr, srvPort, err = parseEndpoint(srvEndpointEnv)
	}

	if err == nil {
		conf.SetSrvAddr(srvAddr)
		err := conf.SetSrvPort(srvPort)
		if err != nil {
			log.Println(err)
		}
	}

	storeIntervalEnv, isExist := os.LookupEnv("STORE_INTERVAL")
	if isExist {
		storeIntervalInt, err := strconv.Atoi(storeIntervalEnv)
		if err == nil {
			err := conf.SetStoreInterval(storeIntervalInt)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		err := conf.SetStoreInterval(int(*storeInterval))
		if err != nil {
			log.Println(err)
		}
	}

	dbConnectionStringEnv, isExist := os.LookupEnv("DATABASE_DSN")
	if isExist {
		err := conf.SetDBConnectionString(dbConnectionStringEnv)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := conf.SetDBConnectionString(*dbConnectionString)
		if err != nil {
			log.Println(err)
		}
	}

	fileStoragePathEnv, isExist := os.LookupEnv("FILE_STORAGE_PATH")
	if isExist {
		err := conf.SetFileStoragePath(fileStoragePathEnv)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := conf.SetFileStoragePath(*fileStoragePath)
		if err != nil {
			log.Println(err)
		}
	}

	restoreMetricsEnv, isExist := os.LookupEnv("RESTORE")
	if isExist {
		restoreMetricsBool, err := strconv.ParseBool(restoreMetricsEnv)
		if err == nil {
			conf.SetRestoreMetrics(restoreMetricsBool)
		}
	} else {
		conf.SetRestoreMetrics(*restoreMetrics)
	}

	return conf
}
