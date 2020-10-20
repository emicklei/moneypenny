package model

import (
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/emicklei/moneypenny/util"
)

// LabeledCost represents a row of the cost query
type LabeledCost struct {
	// The tag names must match those in the sql query
	Charges float64 `bigquery:"charges" json:"charges"`
	Credits float64 `bigquery:"credits" json:"credits"`

	ProjectName string `bigquery:"name" json:"project-name"`
	ProjectID   string `bigquery:"id" json:"project-id"`

	GCPService bigquery.NullString `bigquery:"gcp_service" json:"gcp_service"`
	Component  bigquery.NullString `bigquery:"component" json:"component,omitempty"`
	Service    bigquery.NullString `bigquery:"service" json:"service,omitempty"`
	Opex       bigquery.NullString `bigquery:"opex" json:"opex,omitempty"`
}

// GCPServiceMonitorLabel returns a Stackdriver friendly display label for the GCPService
func (c LabeledCost) GCPServiceMonitorLabel() string {
	return strings.ToLower(strings.ReplaceAll(c.GCPService.StringVal, " ", "-"))
}

// CostComputation is the result of running a query
type CostComputation struct {
	Lines         []map[string]bigquery.Value
	ByteProcessed int64
	ExecutionTime time.Duration
	Query         string
}

func LabeledCostFrom(m map[string]bigquery.Value) LabeledCost {
	return LabeledCost{
		Charges:    util.Float64(m["charges"]),
		Credits:    util.Float64(m["credits"]),
		GCPService: util.BQNullString(m["gcp_service"]),
		Component:  util.BQNullString(m["component"]),
		Service:    util.BQNullString(m["service"]),
		Opex:       util.BQNullString(m["opex"]),
	}
}
