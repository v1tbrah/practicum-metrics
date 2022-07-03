package send

import (
	"fmt"
	"github.com/v1tbrah/metricsAndAlerting/internal/metric"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	receiverAddr = "127.0.0.1"
	receiverPort = "8080"
	contentType  = "text/play"
)

func URL(typeM, valM string) string {
	return "http://" + receiverAddr + ":" + receiverPort + "/" + "update/" + typeM + "/" + valM
}

func AllMetrics(m metric.Metrics) {
	MetricAlloc(m)
	MetricBuckHashSys(m)
	MetricFrees(m)
	MetricGCCPUFraction(m)
	MetricGCSys(m)
	MetricHeapAlloc(m)
	MetricHeapIdle(m)
	MetricHeapInuse(m)
	MetricHeapObjects(m)
	MetricHeapReleased(m)
	MetricHeapSys(m)
	MetricLastGC(m)
	MetricLookups(m)
	MetricMCacheInuse(m)
	MetricMCacheSys(m)
	MetricMSpanInuse(m)
	MetricMSpanSys(m)
	MetricNextGC(m)
	MetricNumForcedGC(m)
	MetricNumGC(m)
	MetricOtherSys(m)
	MetricPauseTotalNs(m)
	MetricStackInuse(m)
	MetricStackSys(m)
	MetricSys(m)
	MetricTotalAlloc(m)
	MetricPollCount(m)
	MetricRandomValue(m)
}

func MetricAlloc(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("Alloc", fmt.Sprintf("%f", m.Alloc)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	fmt.Println(string(payload))

}

func MetricBuckHashSys(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("BuckHashSys", fmt.Sprintf("%f", m.BuckHashSys)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal(err)
	}
}

func MetricFrees(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("Frees", fmt.Sprintf("%f", m.Frees)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal(err)
	}
}

func MetricGCCPUFraction(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("GCCPUFraction", fmt.Sprintf("%f", m.GCCPUFraction)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal(err)
	}
}

func MetricGCSys(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("GCSys", fmt.Sprintf("%f", m.GCSys)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricHeapAlloc(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("HeapAlloc", fmt.Sprintf("%f", m.HeapAlloc)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricHeapIdle(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("HeapIdle", fmt.Sprintf("%f", m.HeapIdle)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricHeapInuse(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("HeapInuse", fmt.Sprintf("%f", m.HeapInuse)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricHeapObjects(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("HeapObjects", fmt.Sprintf("%f", m.HeapObjects)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricHeapReleased(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("HeapReleased", fmt.Sprintf("%f", m.HeapReleased)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricHeapSys(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("HeapSys", fmt.Sprintf("%f", m.HeapSys)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricLastGC(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("LastGC", fmt.Sprintf("%f", m.LastGC)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricLookups(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("Lookups", fmt.Sprintf("%f", m.Lookups)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricMCacheInuse(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("MCacheInuse", fmt.Sprintf("%f", m.MCacheInuse)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricMCacheSys(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("MCacheSys", fmt.Sprintf("%f", m.MCacheSys)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricMSpanInuse(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("MSpanInuse", fmt.Sprintf("%f", m.MSpanInuse)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricMSpanSys(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("MSpanSys", fmt.Sprintf("%f", m.MSpanSys)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricNextGC(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("NextGC", fmt.Sprintf("%f", m.NextGC)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricNumForcedGC(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("NumForcedGC", fmt.Sprintf("%f", m.NumForcedGC)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricNumGC(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("NumGC", fmt.Sprintf("%f", m.NumGC)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricOtherSys(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("OtherSys", fmt.Sprintf("%f", m.OtherSys)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricPauseTotalNs(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("PauseTotalNs", fmt.Sprintf("%f", m.PauseTotalNs)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricStackInuse(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("StackInuse", fmt.Sprintf("%f", m.StackInuse)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricStackSys(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("StackSys", fmt.Sprintf("%f", m.StackSys)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricSys(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("Sys", fmt.Sprintf("%f", m.Sys)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricTotalAlloc(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("TotalAlloc", fmt.Sprintf("%f", m.TotalAlloc)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricPollCount(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("PollCount", strconv.FormatInt(m.PollCount, 10)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}

func MetricRandomValue(m metric.Metrics) {
	request, err := http.NewRequest(http.MethodPost, URL("RandomValue", fmt.Sprintf("%f", m.RandomValue)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", contentType)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		//TODO
	}
}
