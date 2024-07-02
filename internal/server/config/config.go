package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/YaNeAndrey/ya-metrics/internal/server/utils"
)

// Config хранит информацию о конфигурации сервера.
type Config struct {
	srvAddr            string
	storeInterval      time.Duration
	fileStoragePath    string
	dbConnectionString string
	restoreMetrics     bool
	encryptionKey      []byte

	serverPrivKey *rsa.PrivateKey
	trustedSubnet *net.IPNet
}

func NewConfig() *Config {
	var c Config
	c.srvAddr = "localhost:8080"
	c.storeInterval = time.Duration(300) * time.Second
	c.fileStoragePath = path.Join("tmp", "metrics-db.json")
	c.dbConnectionString = ""
	c.restoreMetrics = true
	c.encryptionKey = nil
	return &c
}

func (c *Config) SetSrvAddr(srvAddr string) {
	c.srvAddr = srvAddr
}

func (c *Config) SetEncryptionKey(encryptionKey []byte) {
	if len(encryptionKey) != 16 {
		return
	}
	c.encryptionKey = encryptionKey
}

func (c *Config) SetStoreInterval(storeInterval time.Duration) {
	c.storeInterval = storeInterval
}

func (c *Config) SetFileStoragePath(fileStoragePath string) error {
	err := utils.CheckAndCreateFile(fileStoragePath)
	if err != nil {
		return err
	}
	c.fileStoragePath = fileStoragePath
	return nil
}

func (c *Config) SetDBConnectionString(dbConnectionString string) error {
	db, err := utils.TryToOpenDBConnection(dbConnectionString)
	if err != nil {
		return err
	}
	db.Close()
	c.dbConnectionString = dbConnectionString
	return nil
}

func (c *Config) SetRestoreMetrics(restoreMetrics bool) {
	c.restoreMetrics = restoreMetrics
}

func (c *Config) SetTrustedSubnet(mask string) error {
	_, i, err := net.ParseCIDR(mask)
	if err != nil {
		return err
	}
	c.trustedSubnet = i
	return nil
}

func (c *Config) ReadPrivateKey(filePath string) error {
	privateKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.Println(err)
		return err
	}
	c.serverPrivKey = privateKey
	return nil
}

func (c *Config) SrvAddr() string {
	return c.srvAddr
}

func (c *Config) EncryptionKey() []byte {
	return c.encryptionKey
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

func (c *Config) ServerPrivKey() *rsa.PrivateKey {
	return c.serverPrivKey
}

func (c *Config) TrustedSubnet() *net.IPNet {
	return c.trustedSubnet
}

func (c *Config) String() string {
	if c.dbConnectionString != "" {
		return fmt.Sprintf("Server config: { EndPoint: %s; Store interval: %s; DB connection string: %s} ", c.SrvAddr(), c.StoreInterval(), c.DBConnectionString())

	}
	return fmt.Sprintf("Server config: { EndPoint: %s; Store interval: %s; File storage: %s;} ", c.SrvAddr(), c.StoreInterval(), c.FileStoragePath())
}

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
