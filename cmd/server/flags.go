package main

import (
	"encoding/json"
	"flag"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/YaNeAndrey/ya-metrics/internal/server/config"
)

type configJSON struct {
	Address        string
	Restore        bool
	Store_interval time.Duration
	Store_file     string
	Database_dsn   string
	Crypto_key     string
}

type configValues struct {
	Address        string
	Restore        bool
	Store_interval int
	Store_file     string
	Database_dsn   string
	EncryptionKey  string
	Crypto_key     string
}

func parseFlags() *config.Config {

	configFlags := configValues{
		Address:        *flag.String("a", "localhost:8080", "Server endpoint address server:port"),
		Restore:        *flag.Bool("r", true, "Restore old metrics? (true or false)"),
		Store_interval: *flag.Int("i", 300, "Store Interval in seconds"),
		Store_file:     *flag.String("f", ".\\tmp\\metrics-db.json", "File storage path (.json)"),
		Database_dsn:   *flag.String("d", "", "dbConnectionString in Postgres format: postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]"),
		EncryptionKey:  *flag.String("k", "", "encryption key"),
		Crypto_key:     *flag.String("crypto-key", "", "file with server private key"),
	}
	configFilePath := *flag.String("c", "", "config file")
	flag.Parse()

	configFileEnv, isExist := os.LookupEnv("CONFIG")
	if isExist {
		configFilePath = configFileEnv
	}

	file, _ := os.ReadFile(configFilePath)
	cj := configJSON{}

	_ = json.Unmarshal([]byte(file), &cj)

	configEnv := configValues{}

	srvEndpointEnv, isExist := os.LookupEnv("ADDRESS")
	if isExist {
		configEnv.Address = srvEndpointEnv
	}

	storeIntervalEnv, isExist := os.LookupEnv("STORE_INTERVAL")
	if isExist {
		storeIntervalInt, err := strconv.Atoi(storeIntervalEnv)
		if err == nil {
			configEnv.Store_interval = storeIntervalInt
		}
	}

	dbConnectionStringEnv, isExist := os.LookupEnv("DATABASE_DSN")
	if isExist {
		configEnv.Database_dsn = dbConnectionStringEnv
	}

	fileStoragePathEnv, isExist := os.LookupEnv("FILE_STORAGE_PATH")
	if isExist {
		configEnv.Store_file = fileStoragePathEnv
	}

	restoreMetricsEnv, isExist := os.LookupEnv("RESTORE")
	if isExist {
		restoreMetricsBool, err := strconv.ParseBool(restoreMetricsEnv)
		if err == nil {
			configEnv.Restore = restoreMetricsBool
		}
	}

	encryptionKeyEnv, isExist := os.LookupEnv("KEY")
	if !isExist {
		configEnv.EncryptionKey = encryptionKeyEnv
	}

	serverPrivKeyEnv, isExist := os.LookupEnv("CRYPTO_KEY")
	if isExist {
		configEnv.Crypto_key = serverPrivKeyEnv
	}

	return fillСonfig(cj, configFlags, configEnv)
}

func fillСonfig(cj configJSON, cf configValues, ce configValues) *config.Config {
	conf := config.NewConfig()

	if ce.Restore || cf.Restore || cj.Restore {
		conf.SetRestoreMetrics(true)
	}

	if ce.EncryptionKey != "" {
		conf.SetEncryptionKey([]byte(ce.EncryptionKey))
	} else {
		conf.SetEncryptionKey([]byte(cf.EncryptionKey))
	}

	if ce.Address != "" {
		if checkEndpoint(ce.Address) == nil {
			conf.SetSrvAddr(ce.Address)
		}
	} else if cf.Address != "" {
		if checkEndpoint(cf.Address) == nil {
			conf.SetSrvAddr(cf.Address)
		}
	} else if cj.Address != "" {
		if checkEndpoint(cj.Address) == nil {
			conf.SetSrvAddr(cj.Address)
		}
	}

	if ce.Store_interval > 0 {
		conf.SetStoreInterval(time.Duration(ce.Store_interval) * time.Second)
	} else if cf.Store_interval > 0 {
		conf.SetStoreInterval(time.Duration(cf.Store_interval) * time.Second)
	} else if cj.Store_interval != 0 {
		conf.SetStoreInterval(cj.Store_interval)
	}

	isSet := true
	if ce.Store_file != "" {
		isSet = conf.SetFileStoragePath(ce.Store_file) == nil
	}
	if !isSet && cf.Store_file != "" {
		isSet = conf.SetFileStoragePath(cf.Store_file) == nil
	}
	if !isSet && cj.Store_file != "" {
		_ = conf.SetFileStoragePath(cj.Store_file)
	}

	isSet = true
	if ce.Crypto_key != "" {
		isSet = conf.ReadPrivateKey(ce.Crypto_key) == nil
	}
	if !isSet && cf.Crypto_key != "" {
		isSet = conf.ReadPrivateKey(cf.Crypto_key) == nil
	}
	if !isSet && cj.Crypto_key != "" {
		_ = conf.ReadPrivateKey(cj.Crypto_key)
	}

	isSet = true
	if ce.Database_dsn != "" {
		isSet = conf.SetDBConnectionString(ce.Database_dsn) == nil
	}
	if !isSet && cf.Database_dsn != "" {
		isSet = conf.SetDBConnectionString(cf.Database_dsn) == nil
	}
	if !isSet && cj.Database_dsn != "" {
		_ = conf.SetDBConnectionString(cj.Database_dsn)
	}
	return conf
}

func checkEndpoint(endpointStr string) error {
	hp := strings.Split(endpointStr, ":")
	if len(hp) != 2 {
		return constants.ErrIncorrectEndpointFormat
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	if port < 65535 && port > 0 {
		return constants.ErrIncorrectPortNumber
	}
	return nil
}
