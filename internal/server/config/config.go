package config

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	StorageTypeMemory = iota
	StorageTypeDB
)

type config struct {
	servAddr string

	storageType int

	pgConnString string

	storeInterval time.Duration
	storeFile     string
	restore       bool

	hashKey string

	logLevel string
}

func New(args ...string) (*config, error) {
	log.Debug().Strs("args", args).Msg("config.New started")
	cfg := &config{}
	defer func() {
		log.Debug().Str("config", cfg.String()).Msg("config.New ended")
	}()

	cfg.servAddr = "127.0.0.1:8080"
	cfg.storeInterval = time.Second * 300
	cfg.restore = true
	cfg.logLevel = "debug"

	for _, arg := range args {
		switch arg {
		case WithFlag:
			cfg.parseFromOsArgs()
		case WithEnv:
			if err := cfg.parseFromEnv(); err != nil {
				return nil, err
			}
		}
	}

	if cfg.pgConnString != "" {
		cfg.storageType = StorageTypeDB
	} else {
		cfg.storageType = StorageTypeMemory
	}

	return cfg, nil
}

func (c *config) ServAddr() string {
	return c.servAddr
}

func (c *config) StorageType() int {
	return c.storageType
}

func (c *config) PgConnString() string {
	return c.pgConnString
}

func (c *config) StoreInterval() time.Duration {
	return c.storeInterval
}

func (c *config) StoreFile() string {
	return c.storeFile
}

func (c *config) Restore() bool {
	return c.restore
}

func (c *config) HashKey() string {
	return c.hashKey
}

func (c *config) LogLevel() string {
	return c.logLevel
}

func (c *config) String() string {
	restoreStr := "false"
	if c.restore {
		restoreStr = "true"
	}
	return "ServAddr: " + c.servAddr +
		", StorageType: " + strconv.Itoa(c.storageType) +
		", PgConnString: " + c.pgConnString +
		", StoreInterval: " + c.storeInterval.String() +
		", StoreFile: " + c.storeFile +
		", Restore: " + restoreStr
}
