package datadog

import (
	"github.com/monitoror/monitoror/models"
	datadogModels "github.com/monitoror/monitoror/monitorable/datadog/models"
)

const (
	DatadogMetricTileType models.TileType = "DATADOG-METRIC"
)

type (
	Usecase interface {
		Metric(params *datadogModels.MetricParams) (*models.Tile, error)
	}
)
