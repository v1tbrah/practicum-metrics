package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/service"
)

func setupLog() {
	zerolog.TimeFieldFormat = time.RFC822
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	setupLog()
	log.Debug().Str("application", "agent").Msg("main started")
	defer log.Debug().Str("application", "agent").Msg("main ended")

	cfgOptions := []string{config.WithFlag, config.WithEnv}
	newCfg, err := config.New(cfgOptions...)
	if err != nil {
		log.Fatal().
			Err(err).
			Strs("config options", cfgOptions).
			Msg("unable to create new config")
	}
	if newCfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	newService, err := service.New(newCfg)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("config", fmt.Sprint(newCfg)).
			Strs("config options", cfgOptions).
			Msg("unable to create new service")
	}

	newService.Run()
}
