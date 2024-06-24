package main

import (
	"encoding/json"
	"flag"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
)

type configJSON struct {
	Address        string        `json:"address,omitempty"`
	ReportInterval time.Duration `json:"report_interval,omitempty"`
	PollInterval   time.Duration `json:"poll_interval,omitempty"`
	CryptoKey      string        `json:"crypto_key,omitempty"`
}

type configValues struct {
	Address        string
	ReportInterval int
	PollInterval   int
	EncryptionKey  string
	RateLimit      int
	CryptoKey      string
}

func parseFlags() *config.Config {
	configFlags := configValues{
		Address:        *flag.String("a", "localhost:8080", "Server endpoint address server:port"),
		ReportInterval: *flag.Int("r", 0, "Report Interval in seconds"),
		PollInterval:   *flag.Int("p", 0, "Pool Interval in seconds"),
		EncryptionKey:  *flag.String("k", "", "encryption key"),
		RateLimit:      *flag.Int("l", 0, "rate limit"),
		CryptoKey:      *flag.String("crypto-key", "public.pem", "file with server public key"),
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

	rateLimitEnv, isExist := os.LookupEnv("RATE_LIMIT")
	if isExist {
		rateLimitlInt, err := strconv.Atoi(rateLimitEnv)
		if err == nil {
			configEnv.RateLimit = rateLimitlInt
		}
	}

	reportIntervalEnv, isExist := os.LookupEnv("REPORT_INTERVAL")
	if isExist {
		reportIntervalInt, err := strconv.Atoi(reportIntervalEnv)
		if err == nil {
			configEnv.ReportInterval = reportIntervalInt
		}
	}

	pollIntervalEnv, isExist := os.LookupEnv("POLL_INTERVAL")
	if isExist {
		pollIntervalInt, err := strconv.Atoi(pollIntervalEnv)
		if err == nil {
			configEnv.PollInterval = pollIntervalInt
		}
	}

	encryptionKeyEnv, isExist := os.LookupEnv("KEY")
	if isExist {
		configEnv.EncryptionKey = encryptionKeyEnv
	}

	serverPubKeyEnv, isExist := os.LookupEnv("CRYPTO_KEY")
	if isExist {
		configEnv.CryptoKey = serverPubKeyEnv
	}
	return fillСonfig(cj, configFlags, configEnv)
}

func fillСonfig(cj configJSON, cf configValues, ce configValues) *config.Config {
	conf := config.NewConfig()

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

	if ce.EncryptionKey != "" {
		conf.SetEncryptionKey([]byte(ce.EncryptionKey))
	} else {
		conf.SetEncryptionKey([]byte(cf.EncryptionKey))
	}

	if ce.ReportInterval > 0 {
		conf.SetReportInterval(time.Duration(ce.ReportInterval) * time.Second)
	} else if cf.ReportInterval > 0 {
		conf.SetReportInterval(time.Duration(cf.ReportInterval) * time.Second)
	} else if cj.ReportInterval > 0 {
		conf.SetReportInterval(cj.ReportInterval)
	}

	if ce.PollInterval > 0 {
		conf.SetPollInterval(time.Duration(ce.ReportInterval) * time.Second)
	} else if cf.PollInterval > 0 {
		conf.SetPollInterval(time.Duration(cf.ReportInterval) * time.Second)
	} else if cj.PollInterval > 0 {
		conf.SetPollInterval(cj.ReportInterval)
	}

	isSet := true
	if ce.CryptoKey != "" {
		isSet = conf.ReadServerPubicKey(ce.CryptoKey) == nil
	}
	if !isSet && cf.CryptoKey != "" {
		isSet = conf.ReadServerPubicKey(cf.CryptoKey) == nil
	}
	if !isSet && cj.CryptoKey != "" {
		_ = conf.ReadServerPubicKey(cj.CryptoKey)
	}

	if ce.RateLimit > 0 {
		conf.SetRateLimit(ce.RateLimit)
	} else if cf.RateLimit > 0 {
		conf.SetRateLimit(cf.RateLimit)
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
