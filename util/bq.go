package util

import "cloud.google.com/go/bigquery"

func BQNullString(sornil bigquery.Value) bigquery.NullString {
	if s, ok := sornil.(string); ok {
		return bigquery.NullString{StringVal: s, Valid: true}
	}
	return bigquery.NullString{}
}

func Float64(fornil interface{}) float64 {
	if f, ok := fornil.(float64); ok {
		return f
	}
	return 0.0
}
