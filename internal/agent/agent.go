package agent

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/v1tbrah/metricsAndAlerting/internal/agent/metric"
)

type agent struct {
	client  *resty.Client
	options *options
	metrics *metric.Metrics
}

type options struct {
	pollInterval      time.Duration
	reportInterval    time.Duration
	srvAddr           string
	updateTemplateURL string
	getTemplateURL    string
	contentType       string
}

func newDefaultOptions() *options {
	srvAddr := "127.0.0.1:8080"
	return &options{
		pollInterval:      2 * time.Second,
		reportInterval:    10 * time.Second,
		srvAddr:           srvAddr,
		updateTemplateURL: "http://" + srvAddr + "/update/{typeM}/{nameM}/{valM}",
		getTemplateURL:    "http://" + srvAddr + "/value/{typeM}/{nameM}",
		contentType:       "text/plain",
	}
}

// Creates the agent instance with default settings.
func NewAgent() *agent {
	return &agent{
		client:  resty.New(),
		options: newDefaultOptions(),
		metrics: metric.New()}
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
			a.metrics.Update()
			log.Println("Metrics updated successfully.")
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
			log.Println("Metrics sent successfully.")
			mutex.Unlock()
		}
	}()

	<-shutdown
	os.Exit(0)

}

func (a *agent) sendAllMetrics() error {

	if err := a.sendMetric("Alloc"); err != nil {
		return err
	}

	if err := a.sendMetric("BuckHashSys"); err != nil {
		return err
	}

	if err := a.sendMetric("Frees"); err != nil {
		return err
	}

	if err := a.sendMetric("Frees"); err != nil {
		return err
	}

	if err := a.sendMetric("GCCPUFraction"); err != nil {
		return err
	}

	if err := a.sendMetric("GCSys"); err != nil {
		return err
	}

	if err := a.sendMetric("HeapAlloc"); err != nil {
		return err
	}

	if err := a.sendMetric("HeapIdle"); err != nil {
		return err
	}

	if err := a.sendMetric("HeapInuse"); err != nil {
		return err
	}

	if err := a.sendMetric("HeapObjects"); err != nil {
		return err
	}

	if err := a.sendMetric("HeapReleased"); err != nil {
		return err
	}

	if err := a.sendMetric("HeapSys"); err != nil {
		return err
	}

	if err := a.sendMetric("LastGC"); err != nil {
		return err
	}

	if err := a.sendMetric("Lookups"); err != nil {
		return err
	}

	if err := a.sendMetric("MCacheInuse"); err != nil {
		return err
	}

	if err := a.sendMetric("MCacheSys"); err != nil {
		return err
	}

	if err := a.sendMetric("MSpanInuse"); err != nil {
		return err
	}

	if err := a.sendMetric("MSpanSys"); err != nil {
		return err
	}

	if err := a.sendMetric("NextGC"); err != nil {
		return err
	}

	if err := a.sendMetric("NumForcedGC"); err != nil {
		return err
	}

	if err := a.sendMetric("NumGC"); err != nil {
		return err
	}

	if err := a.sendMetric("OtherSys"); err != nil {
		return err
	}

	if err := a.sendMetric("PauseTotalNs"); err != nil {
		return err
	}

	if err := a.sendMetric("StackInuse"); err != nil {
		return err
	}

	if err := a.sendMetric("StackSys"); err != nil {
		return err
	}

	if err := a.sendMetric("Sys"); err != nil {
		return err
	}

	if err := a.sendMetric("TotalAlloc"); err != nil {
		return err
	}

	if err := a.sendMetric("PollCount"); err != nil {
		return err
	}

	if err := a.sendMetric("RandomValue"); err != nil {
		return err
	}

	return nil
}

func (a *agent) sendMetric(nameM string) error {

	infoM, err := a.metrics.Info(nameM)
	if err != nil {
		return err
	}

	_, err = a.client.NewRequest().
		SetHeader("Content-Type", a.options.contentType).
		SetPathParams(map[string]string{
			"typeM": infoM.TypeM(),
			"nameM": infoM.NameM(),
			"valM":  infoM.ValM()}).
		Post(a.options.updateTemplateURL)

	if err != nil {
		return err
	}

	return nil
}

func (a *agent) getMetric(nameM string) (string, error) {

	infoM, err := a.metrics.Info(nameM)
	if err != nil {
		return "", err
	}

	resp, err := a.client.NewRequest().
		SetHeader("Content-Type", a.options.contentType).
		SetPathParams(map[string]string{
			"typeM": infoM.TypeM(),
			"nameM": infoM.NameM()}).
		Get(a.options.getTemplateURL)

	if err != nil {
		return "", err
	}

	return string(resp.Body()), nil
}
