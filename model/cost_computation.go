package model

import (
	"time"

	"cloud.google.com/go/bigquery"
)

// CostComputation is the result of running a query
type CostComputation struct {
	Lines         []map[string]bigquery.Value
	ByteProcessed int64
	ExecutionTime time.Duration
	Query         string
}
