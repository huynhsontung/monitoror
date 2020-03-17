package usecase

import (
	"fmt"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/datadog"
	datadogModels "github.com/monitoror/monitoror/monitorable/datadog/models"
)

type datadogUsecase struct {
	repository datadog.Repository
}

func NewDatadogUsecase(repository datadog.Repository) datadog.Usecase {
	return &datadogUsecase{repository}
}

func (du *datadogUsecase) Metric(params *datadogModels.MetricParams) (*models.Tile, error) {
	tile := models.NewTile(datadog.DatadogMetricTileType).WithValue(models.NumberUnit)
	tile.Label = "Datadog Metric"

	metric, err := du.repository.GetMetric(params.Query, params.Timespan)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile}
	}
	if len(metric.Pointlist) < 1 {
		return nil, &models.MonitororError{Tile: tile, Message: "No data received"}
	}

	if metric.Metric != "" {
		tile.Label = metric.Metric
	}
	latestData := metric.Pointlist[len(metric.Pointlist)-1].Value
	if params.Threshold != 0 && latestData > float64(params.Threshold) {
		tile.Status = models.WarningStatus
	} else {
		tile.Status = models.SuccessStatus
	}
	tile.Value.Values = append(tile.Value.Values, fmt.Sprintf("%.2f", latestData))
	return tile, nil
}
