package project

import (
	"context"
	"log"
	"time"

	"github.com/emicklei/moneypenny/gcp"
	"github.com/emicklei/moneypenny/model"
	"github.com/urfave/cli/v2"
	"gonum.org/v1/gonum/stat"
)

func DetectProjectCostAnomalies(c *cli.Context, p model.Params) error {
	log.Println("DetectProjectCostAnomalies", p)
	dayTo := p.Date().Add(-24 * time.Hour)
	// 30 days back
	dayFrom := dayTo.Add(-30 * time.Hour * 24)

	q := QueryPastDays(p.BillingTableFQN, dayFrom, dayTo)

	cost, err := gcp.RunBigQuery(context.Background(), p, q)
	if err != nil {
		return err
	}
	log.Println("daily cost entries:", len(cost.Lines))
	statsMap := map[string]*ProjectStats{}
	for _, each := range cost.Lines {
		dc := DailyCostFrom(each)
		dstats, ok := statsMap[dc.ProjectID]
		if !ok {
			dstats = new(ProjectStats)
			statsMap[dc.ProjectID] = dstats
		}
		dstats.Daily = append(dstats.Daily, dc)
	}
	log.Println("projects:", len(statsMap))

	for _, each := range statsMap {
		charges := []float64{}
		for i, other := range each.Daily {
			// all but end date; we will compare against this
			if i > 0 {
				charges = append(charges, other.Charges)
			}
		}
		avg, stddev := stat.MeanStdDev(charges, nil)
		each.Mean = avg
		each.StandardDeviation = stddev
	}

	log.Println("detecting anomalies...")
	detector := BestSundaySky
	anomalies := []*ProjectStats{}
	for id, each := range statsMap {
		if detector.IsAnomaly(each) {
			anomalies = append(anomalies, each)
			log.Println("id:", id, "cost:", each.Daily[0].Charges, "avg:", each.Mean, "stddev:", each.StandardDeviation, "day:", each.Daily[0].Day.String())
		}
	}
	util.ExportJSON(anomalies, "DetectProjectCostAnomalies.json")
	log.Println("done")
	return nil
}
