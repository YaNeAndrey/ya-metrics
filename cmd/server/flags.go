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
	Address       string        `json:"address,omitempty"`
	Restore       bool          `json:"restore,omitempty"`
	StoreInterval time.Duration `json:"store_interval,omitempty"`
	StoreFile     string        `json:"store_file,omitempty"`
	DatabaseDSN   string        `json:"database_dsn,omitempty"`
	CryptoKey     string        `json:"crypto_key,omitempty"`
}

type configValues struct {
	Address       string
	Restore       bool
	StoreInterval int
	StoreFile     string
	DatabaseDSN   string
	EncryptionKey string
	CryptoKey     string
}

func parseFlags() *config.Config {

	configFlags := configValues{
		Address:       *flag.String("a", "", "Server endpoint address server:port"),
		Restore:       *flag.Bool("r", true, "Restore old metrics? (true or false)"),
		StoreInterval: *flag.Int("i", 300, "Store Interval in seconds"),
		StoreFile:     *flag.String("f", ".\\tmp\\metrics-db.json", "File storage path (.json)"),
		DatabaseDSN:   *flag.String("d", "", "dbConnectionString in Postgres format: postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]"),
		EncryptionKey: *flag.String("k", "", "encryption key"),
		CryptoKey:     *flag.String("crypto-key", "", "file with server private key"),
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
			configEnv.StoreInterval = storeIntervalInt
		}
	}

	dbConnectionStringEnv, isExist := os.LookupEnv("DATABASE_DSN")
	if isExist {
		configEnv.DatabaseDSN = dbConnectionStringEnv
	}

	fileStoragePathEnv, isExist := os.LookupEnv("FILE_STORAGE_PATH")
	if isExist {
		configEnv.StoreFile = fileStoragePathEnv
	}

	restoreMetricsEnv, isExist := os.LookupEnv("RESTORE")
	if isExist {
		restoreMetricsBool, err := strconv.ParseBool(restoreMetricsEnv)
		if err == nil {
			configEnv.Restore = restoreMetricsBool
		}
	}

	encryptionKeyEnv, isExist := os.LookupEnv("KEY")
	if isExist {
		configEnv.EncryptionKey = encryptionKeyEnv
	}

	serverPrivKeyEnv, isExist := os.LookupEnv("CRYPTO_KEY")
	if isExist {
		configEnv.CryptoKey = serverPrivKeyEnv
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

	if ce.StoreInterval > 0 {
		conf.SetStoreInterval(time.Duration(ce.StoreInterval) * time.Second)
	} else if cf.StoreInterval > 0 {
		conf.SetStoreInterval(time.Duration(cf.StoreInterval) * time.Second)
	} else if cj.StoreInterval != 0 {
		conf.SetStoreInterval(cj.StoreInterval)
	}

	isSet := true
	if ce.StoreFile != "" {
		isSet = conf.SetFileStoragePath(ce.StoreFile) == nil
	}
	if !isSet && cf.StoreFile != "" {
		isSet = conf.SetFileStoragePath(cf.StoreFile) == nil
	}
	if !isSet && cj.StoreFile != "" {
		_ = conf.SetFileStoragePath(cj.StoreFile)
	}

	isSet = true
	if ce.CryptoKey != "" {
		isSet = conf.ReadPrivateKey(ce.CryptoKey) == nil
	}
	if !isSet && cf.CryptoKey != "" {
		isSet = conf.ReadPrivateKey(cf.CryptoKey) == nil
	}
	if !isSet && cj.CryptoKey != "" {
		_ = conf.ReadPrivateKey(cj.CryptoKey)
	}

	isSet = true
	if ce.DatabaseDSN != "" {
		isSet = conf.SetDBConnectionString(ce.DatabaseDSN) == nil
	}
	if !isSet && cf.DatabaseDSN != "" {
		isSet = conf.SetDBConnectionString(cf.DatabaseDSN) == nil
	}
	if !isSet && cj.DatabaseDSN != "" {
		_ = conf.SetDBConnectionString(cj.DatabaseDSN)
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
