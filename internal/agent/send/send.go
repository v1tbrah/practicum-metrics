package send

import (
	"github.com/v1tbrah/metricsAndAlerting/internal/metric"
	"net/http"
)

const (
	addr        = "127.0.0.1:8080"
	contentType = "text/plain"
)

func AllMetrics(m *metric.Metrics) error {
	if err := Metric(m, "Alloc"); err != nil {
		return err
	}
	if err := Metric(m, "BuckHashSys"); err != nil {
		return err
	}
	if err := Metric(m, "Frees"); err != nil {
		return err
	}
	if err := Metric(m, "GCCPUFraction"); err != nil {
		return err
	}
	if err := Metric(m, "GCSys"); err != nil {
		return err
	}
	if err := Metric(m, "HeapAlloc"); err != nil {
		return err
	}
	if err := Metric(m, "HeapIdle"); err != nil {
		return err
	}
	if err := Metric(m, "HeapInuse"); err != nil {
		return err
	}
	if err := Metric(m, "HeapObjects"); err != nil {
		return err
	}
	if err := Metric(m, "HeapReleased"); err != nil {
		return err
	}
	if err := Metric(m, "HeapSys"); err != nil {
		return err
	}
	if err := Metric(m, "LastGC"); err != nil {
		return err
	}
	if err := Metric(m, "Lookups"); err != nil {
		return err
	}
	if err := Metric(m, "MCacheInuse"); err != nil {
		return err
	}
	if err := Metric(m, "MCacheSys"); err != nil {
		return err
	}
	if err := Metric(m, "MSpanInuse"); err != nil {
		return err
	}
	if err := Metric(m, "MSpanSys"); err != nil {
		return err
	}
	if err := Metric(m, "NextGC"); err != nil {
		return err
	}
	if err := Metric(m, "NumForcedGC"); err != nil {
		return err
	}
	if err := Metric(m, "NumGC"); err != nil {
		return err
	}
	if err := Metric(m, "OtherSys"); err != nil {
		return err
	}
	if err := Metric(m, "PauseTotalNs"); err != nil {
		return err
	}
	if err := Metric(m, "StackInuse"); err != nil {
		return err
	}
	if err := Metric(m, "StackSys"); err != nil {
		return err
	}
	if err := Metric(m, "Sys"); err != nil {
		return err
	}
	if err := Metric(m, "TotalAlloc"); err != nil {
		return err
	}
	if err := Metric(m, "PollCount"); err != nil {
		return err
	}
	if err := Metric(m, "RandomValue"); err != nil {
		return err
	}
	return nil
}

func Metric(m *metric.Metrics, nameM string) error {
	request, err := request(m, nameM)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func request(m *metric.Metrics, nameM string) (*http.Request, error) {
	infoMetric, err := m.Info(nameM)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodPost, url(infoMetric.TypeM(), infoMetric.NameM(), infoMetric.ValM()), nil)
}

func url(typeM, nameM, valM string) string {
	return "http://" + addr + "/" + "update/" + typeM + "/" + nameM + "/" + valM
}
