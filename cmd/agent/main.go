package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/service"
)

func main() {
	log.Debug().Str("application", "agent").Msg("main started")
	defer log.Debug().Str("application", "agent").Msg("main ended")

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

	newData := memory.New()

	newService, err := service.New(newData, newCfg)
	if err != nil {
		log.Fatal().Err(err).Str("config", newCfg.String()).
			Strs("config options", cfgOptions).
			Msg("unable to create new service")
	}

	newService.Run()
}
