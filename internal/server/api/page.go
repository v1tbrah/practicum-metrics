package api

import (
	"fmt"
	"sort"

	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
)

func fillMetricsForPage(dataForPage *[]string, listMetrics model.Data) {
	log.Debug().Str("metrics", fmt.Sprint(listMetrics)).Msg("api.fillMetricsForPage started")
	log.Debug().Msg("api.fillMetricsForPage ended")

	*dataForPage = append(*dataForPage, sortedMetricsForPage(listMetrics)...)
}

func sortedMetricsForPage(listMetrics model.Data) []string {
	log.Debug().Str("metrics", fmt.Sprint(listMetrics)).Msg("api.sortedMetricsForPage started")
	log.Debug().Msg("api.sortedMetricsForPage ended")

	sortedMetrics := []string{}
	for _, currMetric := range listMetrics {
		if currMetric.MType == "gauge" {
			sortedMetrics = append(sortedMetrics, currMetric.ID+": "+fmt.Sprintf("%f", *currMetric.Value))
		} else if currMetric.MType == "counter" {
			sortedMetrics = append(sortedMetrics, currMetric.ID+": "+fmt.Sprintf("%v", *currMetric.Delta))
		}
	}

	sort.Slice(sortedMetrics, func(i, j int) bool { return sortedMetrics[i] < sortedMetrics[j] })
	return sortedMetrics
}
