package api

import (
	"bytes"
	"fmt"
	"sort"
	"text/template"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
)

const pageTpl = `
<!DOCTYPE page>
<page>
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<h1> Metrics </h1>
		<ul>
		{{range .Data}}<li>{{ . }}</li>{{else}}<div><strong>Have no information</strong></div>{{end}}
		</ul>
	</body>
</page>`

type dataForPage struct {
	Title string
	Data  []string
}

func newDataForPage() *dataForPage {
	log.Debug().Msg("service.NewDataForPage started")
	defer log.Debug().Msg("service.NewDataForPage ended")

	return &dataForPage{Title: "Metrics"}
}

func (m *dataForPage) page() (string, error) {
	log.Debug().Msg("service.Page started")
	defer log.Debug().Msg("service.Page ended")

	buf := &bytes.Buffer{}
	t, err := template.New("webpage").Parse(pageTpl)
	if err != nil {
		return "", err
	}
	err = t.Execute(buf, m)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}

func fillMetricsForPage(dataForPage *[]string, listMetrics model.Data) {
	log.Debug().Str("metrics", fmt.Sprint(listMetrics)).Msg("api.fillMetricsForPage started")
	defer log.Debug().Msg("api.fillMetricsForPage ended")

	*dataForPage = append(*dataForPage, sortedMetricsForPage(listMetrics)...)
}

func sortedMetricsForPage(listMetrics model.Data) []string {
	log.Debug().Str("metrics", fmt.Sprint(listMetrics)).Msg("api.sortedMetricsForPage started")
	defer log.Debug().Msg("api.sortedMetricsForPage ended")

	sortedMetrics := make([]string, len(listMetrics))
	var i int
	for _, currMetric := range listMetrics {
		if currMetric.MType == "gauge" {
			sortedMetrics[i] = currMetric.ID + ": " + fmt.Sprintf("%f", *currMetric.Value)
		} else if currMetric.MType == "counter" {
			sortedMetrics[i] = currMetric.ID + ": " + fmt.Sprintf("%v", *currMetric.Delta)
		}
		i++
	}

	sort.Slice(sortedMetrics, func(i, j int) bool { return sortedMetrics[i] < sortedMetrics[j] })
	return sortedMetrics
}
