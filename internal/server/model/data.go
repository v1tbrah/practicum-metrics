package model

import (
	"encoding/json"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

type Data struct {
	Metrics map[string]metric.Metrics
	sync.Mutex
}

// NewData returns new data.
func NewData() *Data {
	log.Debug().Msg("model.NewData started")
	defer log.Debug().Msg("model.NewData ended")

	return &Data{Metrics: map[string]metric.Metrics{}}
}

func (d *Data) MarshalJSON() ([]byte, error) {
	log.Debug().Msg("model.MarshalJSON started")
	defer log.Debug().Msg("model.MarshalJSON ended")

	jsonMetrics, err := json.Marshal(d.Metrics)
	if err != nil {
		return nil, err
	}
	return jsonMetrics, nil
}

func (d *Data) UnmarshalJSON(data []byte) error {
	log.Debug().Str("data", string(data)).Msg("model.UnmarshalJSON started")
	defer log.Debug().Msg("model.UnmarshalJSON ended")

	if err := json.Unmarshal(data, &d.Metrics); err != nil {
		return err
	}
	return nil
}
