package project

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type DailyCost struct {
	Day         time.Time `biguery:"consumption_day" json:"consumption-day" `
	ProjectName string    `bigquery:"name" json:"project-name" `
	ProjectID   string    `bigquery:"id" json:"project-id" `
	Charges     float64   `bigquery:"charges" json:"charges"`
	Credits     float64   `bigquery:"credits" json:"credits"`
}

func DailyCostFrom(m map[string]bigquery.Value) DailyCost {
	return DailyCost{
		Day:         m["consumption_day"].(time.Time),
		ProjectName: m["name"].(string),
		ProjectID:   m["id"].(string),
		Charges:     m["charges"].(float64),
		Credits:     m["credits"].(float64),
	}
}

type ProjectStats struct {
	Daily             []DailyCost `json:"daily" `
	Mean              float64     `json:"mean" `
	StandardDeviation float64     `json:"stddev" `
}
