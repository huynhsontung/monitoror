package datadog

import (
	"github.com/monitoror/monitoror/monitorable/datadog/models"
)

type Repository interface {
	GetMetric(query string, timespan uint) (*models.DatadogMetric, error)
}
