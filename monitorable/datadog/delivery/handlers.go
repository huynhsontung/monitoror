package delivery

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/datadog"
	datadogModels "github.com/monitoror/monitoror/monitorable/datadog/models"
)

type DatadogDelivery struct {
	datadogUsecase datadog.Usecase
}

func NewDatadogDelivery(u datadog.Usecase) *DatadogDelivery {
	return &DatadogDelivery{u}
}

func (d *DatadogDelivery) GetMetric(c echo.Context) error {
	// Bind / check Params
	params := &datadogModels.MetricParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := d.datadogUsecase.Metric(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
