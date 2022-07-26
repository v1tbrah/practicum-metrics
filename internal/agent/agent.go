package agent

import (
	"encoding/json"

	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/memory"
	"github.com/v1tbrah/metricsAndAlerting/internal/agent/metric"
)

type agent struct {
	client  *resty.Client
	options *options
	memory  *memory.MemStorage
}

type options struct {
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	SrvAddr        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
}

// NewAgent with default settings.
func NewAgent() *agent {
	agent := agent{
		client:  resty.New(),
		options: newDefaultOptions(),
		memory:  memory.NewMemStorage()}
	err := env.Parse(agent.options)
	if err != nil {
		log.Println(err)
	}
	return &agent
}

//Run agent updating the metrics once per pollInterval and sends them to the server once per reportInterval.
//For more about pollInterval and reportInterval, see agent.options.
//
//On signals SIGTERM, SIGINT, SIGQUIT, it exits with code 0.
func (a *agent) Run() {

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	var mutex sync.Mutex

	go func() {
		updateTime := time.NewTicker(a.options.PollInterval)
		for {
			<-updateTime.C
			mutex.Lock()
			a.memory.Metrics.Update()
			log.Println("AllMetrics updated successfully.")
			mutex.Unlock()
		}
	}()

	go func() {
		reportTime := time.NewTicker(a.options.ReportInterval)
		for {
			<-reportTime.C
			mutex.Lock()
			if err := a.sendAllMetricsJSON(); err != nil {
				log.Println(err)
			} else {
				log.Println("AllMetrics sent successfully.")
			}
			mutex.Unlock()
		}
	}()

	<-shutdown
	os.Exit(0)

}

func newDefaultOptions() *options {
	opt := &options{
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
		SrvAddr:        "127.0.0.1:8080",
	}
	return opt
}

func (a *agent) updateTemplateJSONURL() string {
	return "http://" + a.options.SrvAddr + "/update/"
}

func (a *agent) getTemplateJSONURL() string {
	return "http://" + a.options.SrvAddr + "/value/"
}

func (a *agent) sendAllMetricsJSON() error {

	for _, currMetric := range *a.memory.Metrics {
		if _, err := a.sendMetricJSON(currMetric); err != nil {
			return err
		}
	}

	return nil
}

func (a *agent) sendMetricJSON(metric metric.Metrics) (*resty.Response, error) {

	body, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(a.updateTemplateJSONURL())

	return resp, err
}

func (a *agent) getAllMetricJSON() error {

	for _, currMetric := range *a.memory.Metrics {
		if _, err := a.getMetricJSON(currMetric); err != nil {
			return err
		}
	}

	return nil
}

func (a *agent) getMetricJSON(metric metric.Metrics) (*resty.Response, error) {

	metric.Delta = nil
	metric.Value = nil

	body, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(a.getTemplateJSONURL())

	return resp, err
}
