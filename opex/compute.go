package opex

import (
	"context"
	"time"

	"github.com/emicklei/moneypenny/gcp"
	"github.com/emicklei/moneypenny/model"
)

// computeCostPerOpex runs a job for the queryCostPerOpex BigQuery and returns the result
func computeCostPerOpex(p model.Params) (model.CostComputation, error) {
	return gcp.RunBigQuery(context.Background(), p, queryCostPerOpex(p.BillingTableFQN, p.Year, p.MonthIndex, p.Opex))
}

// computeCostPerComponent runs a job for the queryCostPerComponent BigQuery and returns the result
func computeCostPerComponent(p model.Params) (model.CostComputation, error) {
	return gcp.RunBigQuery(context.Background(), p, queryCostPerComponent(p.BillingTableFQN, p.Year, p.MonthIndex))
}

// computeCostPerOpexLastDay runs a job for the queryCostPerOpexLastDay BigQuery and returns the result
func computeCostPerOpexLastDay(p model.Params) (model.CostComputation, error) {
	return gcp.RunBigQuery(context.Background(), p, queryCostPerOpexLastDay(p.BillingTableFQN, time.Now().Add(-24*time.Hour), p.Opex))
}
