package project

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/emicklei/moneypenny/model"
	"github.com/google/uuid"
)

func appendEventsForAnomalies(anomalies []ProjectStatsReport, detector AnomalyDetector, p model.Params) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, p.TargetProjectID())
	if err != nil {
		return err
	}
	defer client.Close()
	inserter := client.Dataset(p.TargetDatasetID()).Table(p.TargetTableID()).Inserter()
	// build events
	events := []AnomalyEvent{}
	for _, each := range anomalies {
		event := AnomalyEvent{
			EventID:           uuid.New().String(),
			EventCreationTime: time.Now(),
			ProjectID:         each.LastDay.ProjectID,
			ProjectName:       each.LastDay.ProjectName,
			Charges:           fs(each.LastDay.Charges),
			ChargesPercentage: fs(each.ChargesPercentage),
			Credits:           fs(each.LastDay.Credits),
			Mean:              fs(each.Mean),
			StandardDeviation: fs(each.StandardDeviation),
			DetectionDay:      each.LastDay.Day,
			Detector:          detector.String(),
		}
		events = append(events, event)
	}
	log.Printf("appending %d events to %s\n", len(events), p.TargetTableFQN)
	return inserter.Put(ctx, events)
}

// this exists because the pkg cannot handle float64 fields directly. TODO
func fs(f float64) string { return fmt.Sprintf("%f", f) }
