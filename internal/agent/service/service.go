package service

//go:generate mockery --all

import (
	"context"
	"errors"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type service struct {
	client *resty.Client
	cfg    Config
	data   Data
}

// New returns service.
func New(data Data, cfg Config) (*service, error) {
	log.Debug().Str("cfg", cfg.String()).Msg("service.New started")
	defer log.Debug().Msg("service.New ended")

	if cfg == nil {
		return nil, errors.New("config is empty")
	}

	newService := service{
		client: resty.New(),
		cfg:    cfg,
		data:   data}

	return &newService, nil
}

//Run service updating the metrics once per pollInterval and sends them to the server once per reportInterval.
func (s *service) Run() {
	log.Printf("service.Run started")
	defer log.Printf("service.Run ended")

	log.Info().Msg("agent starting")
	rand.Seed(time.Now().UnixNano())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return s.updateBasicDataWithInterval(s.cfg.PollInterval(), shutdown)
	})

	g.Go(func() error {
		return s.updateAdditionalDataWithInterval(s.cfg.PollInterval(), shutdown)
	})

	g.Go(func() error {
		return s.reportDataWithInterval(s.cfg.ReportInterval(), shutdown)
	})

	if err := g.Wait(); err != nil {
		log.Error().Err(err).Msg("agent ended")
		return
	}

	log.Info().Msg("agent ended")
	os.Exit(0)

}

func (s *service) updateBasicDataWithInterval(interval time.Duration, shutdown chan os.Signal) error {
	log.Debug().Dur("interval", interval).Msg("service.updateBasicDataWithInterval started")
	defer log.Debug().Msg("service.updateBasicDataWithInterval ended")

	keyForUpdateHashMetric := s.cfg.HashKey()
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			s.data.UpdateBasic(keyForUpdateHashMetric)
		case _, ok := <-shutdown:
			if ok {
				close(shutdown)
			}
			return nil
		}
	}
}

func (s *service) updateAdditionalDataWithInterval(interval time.Duration, shutdown chan os.Signal) error {
	log.Debug().Dur("interval", interval).Msg("service.updateAdditionalDataWithInterval started")
	defer log.Debug().Msg("service.updateAdditionalDataWithInterval ended")

	keyForUpdateHashMetric := s.cfg.HashKey()
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			if err := s.data.UpdateAdditional(keyForUpdateHashMetric); err != nil {
				return err
			}
		case _, ok := <-shutdown:
			if ok {
				close(shutdown)
			}
			return nil
		}
	}
}

func (s *service) reportDataWithInterval(interval time.Duration, shutdown chan os.Signal) error {
	log.Debug().Dur("interval", interval).Msg("service.reportDataWithInterval started")
	defer log.Debug().Msg("service.reportDataWithInterval ended")

	var urlError *url.Error
	tried := 0
	ticker := time.NewTicker(interval)

loop:
	for {
		select {
		case <-ticker.C:
			for _, currMetric := range s.data.GetData() {
				if _, err := s.reportMetric(currMetric); err != nil {
					if errors.As(err, &urlError) {
						tried++
						log.Error().Err(err).Msg("unable reporting metric")
						if tried == 3 {
							return err
						}
						continue loop
					}
					log.Error().Err(err).Str("metric", currMetric.String()).Msg("unable reporting metric")
				} else {
					log.Info().Str("metric", currMetric.String()).Msg("metric reported")
				}
			}
		case _, ok := <-shutdown:
			if ok {
				close(shutdown)
			}
			return nil
		}
	}
}
