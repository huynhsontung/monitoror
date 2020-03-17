package models

import (
	"encoding/json"
	"fmt"
)

type (
	DatadogMetricRespond struct {
		Status string           `json:"status"`
		Series []*DatadogMetric `json:"series"`
	}

	DatadogMetric struct {
		Start     uint64             `json:"start"`
		End       uint64             `json:"end"`
		Interval  uint               `json:"interval"`
		Metric    string             `json:"metric"`
		Scope     string             `json:"scope"`
		Pointlist []*MetricDatapoint `json:"pointlist"`
	}

	MetricDatapoint struct {
		Timestamp float64
		Value     float64
	}
)

func (d *MetricDatapoint) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&d.Timestamp, &d.Value}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in MetricDatapoint: %d != %d", g, e)
	}
	return nil
}
