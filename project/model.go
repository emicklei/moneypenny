package project

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type DailyCost struct {
	Day         time.Time `biguery:"consumption_day"`
	ProjectName string    `bigquery:"name" `
	ProjectID   string    `bigquery:"id"  `
	Charges     float64   `bigquery:"charges" json:"charges"`
	Credits     float64   `bigquery:"credits" json:"cost"`
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
	Daily             []DailyCost
	Mean              float64
	StandardDeviation float64
}
