package config

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	srvAddr            string
	srvPort            int
	storeInterval      time.Duration
	fileStoragePath    string
	dbConnectionString string
	restoreMetrics     bool
}

func NewConfig() *Config {
	var c Config
	c.srvAddr = "localhost"
	c.srvPort = 8080
	c.storeInterval = time.Duration(300) * time.Second
	c.fileStoragePath = path.Join("tmp", "metrics-db.json")
	c.dbConnectionString = ""
	c.restoreMetrics = true
	return &c
}

func (c *Config) SetSrvAddr(srvAddr string) {
	c.srvAddr = srvAddr
}

func (c *Config) SetSrvPort(srvPort int) error {
	if srvPort < 65535 && srvPort > 0 {
		c.srvPort = srvPort
		return nil
	}
	return errors.New("SrvPort must be in [1:65535]")
}

func (c *Config) SetStoreInterval(storeInterval int) error {
	if storeInterval > -1 {
		c.storeInterval = time.Duration(storeInterval) * time.Second
		return nil
	}
	return errors.New("StoreInterval must be greater then -1")
}

func (c *Config) SetFileStoragePath(fileStoragePath string) error {
	err := CheckAndCreateFile(fileStoragePath)
	if err != nil {
		return err
	}
	c.fileStoragePath = fileStoragePath
	return nil
}

func (c *Config) SetDBConnectionString(dbConnectionString string) error {
	err := CheckDBConnection(dbConnectionString)
	if err != nil {
		return err
	}
	c.dbConnectionString = dbConnectionString
	return nil
}

func (c *Config) SetRestoreMetrics(restoreMetrics bool) {
	c.restoreMetrics = restoreMetrics
}

func (c *Config) SrvAddr() string {
	return c.srvAddr
}

func (c *Config) SrvPort() int {
	return c.srvPort
}

func (c *Config) StoreInterval() time.Duration {
	return c.storeInterval
}

func (c *Config) FileStoragePath() string {
	return c.fileStoragePath
}

func (c *Config) DBConnectionString() string {
	return c.dbConnectionString
}

func (c *Config) RestoreMetrics() bool {
	return c.restoreMetrics
}

func CheckAndCreateFile(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(filePath)
			separatedPath := strings.Split(filePath, string(os.PathSeparator))
			log.Println(separatedPath)
			dirPath := strings.Join(separatedPath[0:len(separatedPath)-1], string(os.PathSeparator))
			err = os.MkdirAll(dirPath, 0666)
			if err != nil {
				return err
			}
			_, err := os.Create(filePath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func CheckDBConnection(dbConnectionString string) error {
	db, err := sql.Open("pgx", dbConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
