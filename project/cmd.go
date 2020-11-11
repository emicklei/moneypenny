package project

import (
	"context"
	"log"
	"time"

	"github.com/emicklei/moneypenny/gcp"
	"github.com/emicklei/moneypenny/model"
	"github.com/emicklei/moneypenny/util"
	"github.com/urfave/cli/v2"
	"gonum.org/v1/gonum/stat"
)

// DetectProjectCostAnomalies collects 30 days of costs per project to detect a cost
// Write a report to [DetectProjectCostAnomalies.json] if at lease one anomaly was found.
func DetectProjectCostAnomalies(c *cli.Context, p model.Params) error {
	targetTable := p.TargetTableFQN
	// optionally, anomaly events are written to a table
	if len(targetTable) > 0 {
		util.CheckBigQueryTable(targetTable)
	}
	util.CheckBigQueryTable(p.BillingTableFQN)

	detector := BestSundaySky
	if stddev := c.Float64("sundaysky.stddev"); stddev > 0 {
		detector.stddevThreshold = stddev
	}
	// date is a YYYYMMDD with zero time
	// dayTo must be yesterday
	dayTo := p.Date().Add(-1 * time.Second)
	// 30 days back
	dayFrom := dayTo.Add(-time.Duration(detector.windowDays) * time.Hour * 24)

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

	log.Println("detecting cost anomalies on", dayTo, "...")

	anomalies := []ProjectStatsReport{}
	for id, each := range statsMap {
		if detector.IsAnomaly(each) {
			report := ProjectStatsReport{
				LastDay:           each.Daily[0], // bounds are checked
				Detector:          detector.String(),
				Mean:              each.Mean,
				StandardDeviation: each.StandardDeviation,
			}
			if each.Mean > 0 {
				// bounds are checked
				report.ChargesPercentage = (each.Daily[0].Charges - each.Mean) * 100.0 / each.Mean
			}
			anomalies = append(anomalies, report)
			log.Println("id:", id, "cost:", each.Daily[0].Charges, "avg:", each.Mean, "stddev:", each.StandardDeviation, "day:", each.Daily[0].Day.String())
		}
	}
	// only export report if at least one anomaly found"
	if len(anomalies) > 0 {
		root := map[string]interface{}{}
		root["anomalies"] = anomalies
		root["project_count"] = len(statsMap)
		root["statistics_days"] = detector.windowDays
		if err := util.ExportJSON(root, "DetectProjectCostAnomalies.json"); err != nil {
			return err
		}

		// see if we need to store the non-empty list of events
		if len(targetTable) > 0 {
			return appendEventsForAnomalies(anomalies, detector, p)
		}
	}
	return nil
}
