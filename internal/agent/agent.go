package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	pollInterval          time.Duration
	reportInterval        time.Duration
	srvAddr               string
	updateTemplateURL     string
	getTemplateURL        string
	contentTypeTextPlain  string
	updateTemplateJSONURL string
	getTemplateJSONURL    string
	contentTypeJSON       string
}

// Creates the agent instance with default settings.
func NewAgent() *agent {
	return &agent{
		client:  resty.New(),
		options: newDefaultOptions(),
		memory:  memory.NewMemStorage()}
}

//The agent starts updating the metrics once per pollInterval and sends them to the server once per reportInterval.
//For more about pollInterval and reportInterval, see agent.options.
//
//On signals SIGTERM, SIGINT, SIGQUIT, it exits with code 0.
func (a *agent) Run() {

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	var mutex sync.Mutex

	go func() {
		updateTime := time.NewTicker(a.options.pollInterval)
		for {
			<-updateTime.C
			mutex.Lock()
			a.memory.Metrics.Update()
			log.Println("AllMetrics updated successfully.")
			mutex.Unlock()
		}
	}()

	go func() {
		reportTime := time.NewTicker(a.options.reportInterval)
		for {
			<-reportTime.C
			mutex.Lock()
			if err := a.sendAllMetrics(); err != nil {
				log.Fatalln(err)
			}
			log.Println("AllMetrics sent successfully.")
			mutex.Unlock()
		}
	}()

	<-shutdown
	os.Exit(0)

}

func newDefaultOptions() *options {
	srvAddr := "127.0.0.1:8080"
	return &options{
		pollInterval:          2 * time.Second,
		reportInterval:        10 * time.Second,
		srvAddr:               srvAddr,
		updateTemplateURL:     "http://" + srvAddr + "/update/{typeM}/{nameM}/{valM}",
		getTemplateURL:        "http://" + srvAddr + "/value/{typeM}/{nameM}",
		contentTypeTextPlain:  "text/plain",
		updateTemplateJSONURL: "http://" + srvAddr + "/update/",
		getTemplateJSONURL:    "http://" + srvAddr + "/value/",
		contentTypeJSON:       "application/json",
	}
}

func (a *agent) sendAllMetrics() error {

	for _, metric := range *a.memory.Metrics {
		if _, err := a.sendMetric(metric); err != nil {
			return err
		}
	}

	return nil
}

func (a *agent) sendMetric(metric metric.Metrics) (*resty.Response, error) {

	var valM string
	if metric.MType == "gauge" {
		valM = fmt.Sprintf("%f", *metric.Value)
	} else {
		valM = fmt.Sprintf("%v", *metric.Delta)
	}

	resp, err := a.client.NewRequest().
		SetHeader("Content-Type", a.options.contentTypeTextPlain).
		SetPathParams(map[string]string{
			"typeM": metric.MType,
			"nameM": metric.ID,
			"valM":  valM}).
		Post(a.options.updateTemplateURL)

	return resp, err
}

func (a *agent) getMetric(metric metric.Metrics) (*resty.Response, error) {

	resp, err := a.client.NewRequest().
		SetHeader("Content-Type", a.options.contentTypeJSON).
		SetPathParams(map[string]string{
			"typeM": metric.MType,
			"nameM": metric.ID}).
		Get(a.options.getTemplateURL)

	return resp, err
}

func (a *agent) sendAllMetricsJSON() error {

	for _, metric := range *a.memory.Metrics {
		if _, err := a.sendMetricJSON(metric); err != nil {
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
		SetHeader("Content-Type", a.options.contentTypeJSON).
		SetBody(body).
		Post(a.options.updateTemplateJSONURL)

	return resp, err
}

func (a *agent) getMetricJSON(metric metric.Metrics) (*resty.Response, error) {

	metric.Delta = nil
	metric.Value = nil

	body, err := json.Marshal(metric)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.NewRequest().
		SetHeader("Content-Type", a.options.contentTypeJSON).
		SetBody(body).
		Post(a.options.getTemplateJSONURL)

	return resp, err
}
