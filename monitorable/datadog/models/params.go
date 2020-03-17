package models

type (
	MetricParams struct {
		Query     string
		Timespan  uint
		Threshold int
	}
)

func (p *MetricParams) IsValid() bool {
	return p.Query != ""
}
