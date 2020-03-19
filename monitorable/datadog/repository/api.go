package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/datadog"
	"github.com/monitoror/monitoror/monitorable/datadog/models"
	"github.com/sourcegraph/httpcache"
)

type datadogRepository struct {
	httpClient *http.Client
	config     *config.Datadog
}

func NewDatadogRepository(config *config.Datadog) datadog.Repository {
	httpClient := &http.Client{
		Transport: httpcache.NewMemoryCacheTransport(),
		Timeout:   time.Duration(config.Timeout) * time.Millisecond,
	}

	return &datadogRepository{
		httpClient: httpClient,
		config:     config,
	}
}

// GetMetric with query and timespan in minutes
func (r *datadogRepository) GetMetric(query string, timespan uint) (*models.DatadogMetric, error) {
	if timespan == 0 {
		timespan = 15
	}
	start := time.Now().Add(-time.Duration(timespan) * time.Minute).Unix()
	end := time.Now().Unix()
	urlParams := url.Values{
		"query": {query},
		"from":  {strconv.FormatInt(start, 10)},
		"to":    {strconv.FormatInt(end, 10)},
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.datadoghq.com/api/v1/query?%s", urlParams.Encode()), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("DD-API-KEY", r.config.APIKey)
	req.Header.Add("DD-APPLICATION-KEY", r.config.ApplicationKey)
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var metricResp models.DatadogMetricRespond
	err = json.Unmarshal(bytes, &metricResp)
	if err != nil {
		return nil, err
	}
	if len(metricResp.Series) < 1 {
		if metricResp.Status == "ok" {
			return &models.DatadogMetric{}, nil
		} else {
			return nil, errors.New("No data received")
		}
	}

	return metricResp.Series[0], nil
}
