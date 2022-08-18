package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/api"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/config"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/repo"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/service"
)

func setupLog() {
	zerolog.TimeFieldFormat = time.RFC822
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	setupLog()
	log.Debug().Str("application", "server").Msg("main started")
	defer log.Debug().Str("application", "server").Msg("main ended")

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

	newStorage, err := repo.New(newCfg)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("config", fmt.Sprint(newCfg)).
			Msg("unable to create new storage")
	}

	newService, err := service.New(newStorage, newCfg)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("storage", fmt.Sprint(newCfg)).
			Str("config", fmt.Sprint(newCfg)).
			Msg("unable to create new service")
	}

	myAPI := api.New(newService)

	myAPI.Run()

}
