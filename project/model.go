package project

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type DailyCost struct {
	Day         time.Time `biguery:"consumption_day" json:"consumption_day" `
	ProjectName string    `bigquery:"project_name" json:"project_name" `
	ProjectID   string    `bigquery:"project_id" json:"project_id" `
	Charges     float64   `bigquery:"charges" json:"charges"`
	Credits     float64   `bigquery:"credits" json:"credits"`
}

func DailyCostFrom(m map[string]bigquery.Value) DailyCost {
	return DailyCost{
		Day:         m["consumption_day"].(time.Time),
		ProjectName: m["project_name"].(string),
		ProjectID:   m["project_id"].(string),
		Charges:     m["charges"].(float64),
		Credits:     m["credits"].(float64),
	}
}

type ProjectStats struct {
	Daily             []DailyCost `json:"daily" `
	Mean              float64     `json:"mean" `
	StandardDeviation float64     `json:"stddev" `
}

type ProjectStatsReport struct {
	LastDay           DailyCost `json:"last_day" `
	Mean              float64   `json:"mean" `
	StandardDeviation float64   `json:"stddev" `
	Detector          string    `json:"detector" `
}
