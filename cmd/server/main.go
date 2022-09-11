package main

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/api"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/storage"
)

func main() {
	log.Debug().Str("application", "server").Msg("main started")
	defer log.Debug().Str("application", "server").Msg("main ended")

	cfgOptions := []string{config.WithFlag, config.WithEnv}
	newCfg, err := config.New(cfgOptions...)
	if err != nil {
		log.Fatal().Err(err).
			Strs("config options", cfgOptions).
			Msg("unable to create new config")
	}

	logLevel, err := zerolog.ParseLevel(newCfg.LogLevel())
	if err != nil {
		log.Fatal().Err(err).
			Strs("config options", cfgOptions).
			Msg("unable to parse log level")
	}
	zerolog.SetGlobalLevel(logLevel)

	newStorage, err := storage.New(newCfg)
	if err != nil {
		log.Fatal().Err(err).
			Str("config", newCfg.String()).
			Msg("unable to create new storage")
	}

	newService, err := service.New(newStorage, newCfg)
	if err != nil {
		log.Fatal().Err(err).
			Str("storage", fmt.Sprint(newStorage)).
			Str("config", newCfg.String()).
			Msg("unable to create new service")
	}

	myAPI := api.New(newService, newCfg)

	myAPI.Run()

}
