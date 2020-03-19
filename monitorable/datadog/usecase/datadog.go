package usecase

import (
	"encoding/json"

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
	tile := models.NewTile(datadog.DatadogMetricTileType).WithValue(models.TrendUnit)
	tile.Label = "Datadog Metric"

	metric, err := du.repository.GetMetric(params.Query, params.Timespan)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile}
	}
	dataLen := len(metric.Pointlist)
	if dataLen < 1 {
		tile.Value.Values = append(tile.Value.Values, "[]")
		return tile, nil
	}

	if metric.Metric != "" {
		tile.Label = metric.Metric
	}
	data := make([]float64, dataLen)
	for i, v := range metric.Pointlist {
		data[i] = v.Value
	}
	latestData := metric.Pointlist[dataLen-1].Value
	if params.Threshold != 0 && latestData > float64(params.Threshold) {
		tile.Status = models.WarningStatus
	} else {
		tile.Status = models.SuccessStatus
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	tile.Value.Values = append(tile.Value.Values, string(bytes))
	return tile, nil
}
