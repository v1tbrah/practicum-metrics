package agent

import (
	"github.com/go-resty/resty/v2"
	"github.com/v1tbrah/metricsAndAlerting/internal/metric"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	pollInterval      = 2 * time.Second
	reportInterval    = 10 * time.Second
	srvAddr           = "127.0.0.1:8080"
	updateTemplateURL = "http://" + srvAddr + "/update/{typeM}/{nameM}/{valM}"
	getTemplateURL    = "http://" + srvAddr + "/value/{typeM}/{nameM}"
	contentType       = "text/plain"
)

type agent struct {
	client  *resty.Client
	metrics *metric.Metrics
}

func NewAgent() *agent {
	return &agent{
		client:  resty.New(),
		metrics: metric.New()}
}

func (a *agent) Run() {

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		updateTime := time.NewTicker(pollInterval)
		for {
			<-updateTime.C
			a.metrics.Update()
		}
	}()

	go func() {
		reportTime := time.NewTicker(reportInterval)
		for {
			<-reportTime.C
			if err := a.SendAllMetrics(); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	<-shutdown
	os.Exit(0)

}

func (a *agent) SendAllMetrics() error {

	if err := a.SendMetric("Alloc"); err != nil {
		return err
	}

	if err := a.SendMetric("BuckHashSys"); err != nil {
		return err
	}

	if err := a.SendMetric("Frees"); err != nil {
		return err
	}

	if err := a.SendMetric("Frees"); err != nil {
		return err
	}

	if err := a.SendMetric("GCCPUFraction"); err != nil {
		return err
	}

	if err := a.SendMetric("GCSys"); err != nil {
		return err
	}

	if err := a.SendMetric("HeapAlloc"); err != nil {
		return err
	}

	if err := a.SendMetric("HeapIdle"); err != nil {
		return err
	}

	if err := a.SendMetric("HeapInuse"); err != nil {
		return err
	}

	if err := a.SendMetric("HeapObjects"); err != nil {
		return err
	}

	if err := a.SendMetric("HeapReleased"); err != nil {
		return err
	}

	if err := a.SendMetric("HeapSys"); err != nil {
		return err
	}

	if err := a.SendMetric("LastGC"); err != nil {
		return err
	}

	if err := a.SendMetric("Lookups"); err != nil {
		return err
	}

	if err := a.SendMetric("MCacheInuse"); err != nil {
		return err
	}

	if err := a.SendMetric("MCacheSys"); err != nil {
		return err
	}

	if err := a.SendMetric("MSpanInuse"); err != nil {
		return err
	}

	if err := a.SendMetric("MSpanSys"); err != nil {
		return err
	}

	if err := a.SendMetric("NextGC"); err != nil {
		return err
	}

	if err := a.SendMetric("NumForcedGC"); err != nil {
		return err
	}

	if err := a.SendMetric("NumGC"); err != nil {
		return err
	}

	if err := a.SendMetric("OtherSys"); err != nil {
		return err
	}

	if err := a.SendMetric("PauseTotalNs"); err != nil {
		return err
	}

	if err := a.SendMetric("StackInuse"); err != nil {
		return err
	}

	if err := a.SendMetric("StackSys"); err != nil {
		return err
	}

	if err := a.SendMetric("Sys"); err != nil {
		return err
	}

	if err := a.SendMetric("TotalAlloc"); err != nil {
		return err
	}

	if err := a.SendMetric("PollCount"); err != nil {
		return err
	}

	if err := a.SendMetric("RandomValue"); err != nil {
		return err
	}

	return nil
}

func (a *agent) SendMetric(nameM string) error {

	infoM, err := a.metrics.Info(nameM)
	if err != nil {
		return err
	}

	_, err = a.client.NewRequest().
		SetHeader("Content-Type", contentType).
		SetPathParams(map[string]string{
			"typeM": infoM.TypeM(),
			"nameM": infoM.NameM(),
			"valM":  infoM.ValM()}).
		Post(updateTemplateURL)

	if err != nil {
		return err
	}

	return nil
}

func (a *agent) GetMetric(nameM string) (string, error) {

	infoM, err1 := a.metrics.Info(nameM)
	if err1 != nil {
		return "", err1
	}

	resp, err2 := a.client.NewRequest().
		SetHeader("Content-Type", contentType).
		SetPathParams(map[string]string{
			"typeM": infoM.TypeM(),
			"nameM": infoM.NameM()}).
		Get(getTemplateURL)

	if err2 != nil {
		return "", err2
	}

	return string(resp.Body()), nil
}
