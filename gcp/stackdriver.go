package gcp

import (
	"context"
	"log"
	"time"

	stackmoni "cloud.google.com/go/monitoring/apiv3"
	"github.com/emicklei/moneypenny/model"
	"github.com/emicklei/moneypenny/util"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// StackDriver provides the api to send a moneypenny result
type StackDriver struct {
	metricsClient *stackmoni.MetricClient
	projectID     string
	ctx           context.Context
}

// NewStackDriver create a connected StackDriver for a given project for which metrics are created.
func NewStackDriver(projectID string) (*StackDriver, error) {
	util.CheckGCPCredentials()

	ctx := context.Background()
	metricsClient, err := stackmoni.NewMetricClient(ctx)
	if err != nil {
		return nil, err
	}
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return &StackDriver{metricsClient: metricsClient, projectID: projectID, ctx: ctx}, nil
}

// Close closes the metrics client
func (s *StackDriver) Close() error {
	return s.metricsClient.Close()
}

// SendReport will sends metrics to StackDriver using measurements of a samples.
func (s *StackDriver) SendMetrics(opex string, lines []model.LabeledCost, dryrun bool) error {
	if len(lines) == 0 {
		log.Println("nothing to report")
		return nil
	}
	metricType := s.metricType()
	resource := s.newResource(opex)

	reportTime := time.Now()
	timeSeries := []*monitoringpb.TimeSeries{}

	for _, each := range lines {
		log.Printf("%#v\n", each)
		if each.Charges > 1.0 {
			metric := &metricpb.Metric{
				Type: metricType,
				Labels: map[string]string{
					"project": each.ProjectID,
					"service": each.GCPServiceMonitorLabel(),
					"opex":    opex,
				},
			}
			dataPoint := newDatapoint(reportTime, each.Charges)
			timeSeries = append(timeSeries, newTimeSeries(dataPoint, metric, resource))
		}
	}
	if dryrun {
		log.Println("skip sending metrics to Stackdriver because dryrun, time series count:", len(timeSeries))
		return nil
	}
	return s.createTimeSeries(timeSeries)
}

func (s *StackDriver) metricType() string {
	return "custom.googleapis.com/moneypenny"
}

func (s *StackDriver) newResource(opex string) *monitoredrespb.MonitoredResource {
	resourceLabels := map[string]string{
		"project_id": s.projectID}
	resourceType := "global"
	return &monitoredrespb.MonitoredResource{
		Type:   resourceType,
		Labels: resourceLabels,
	}
}

func newTimeSeries(dataPoint *monitoringpb.Point,
	metric *metricpb.Metric,
	resource *monitoredrespb.MonitoredResource) *monitoringpb.TimeSeries {
	return &monitoringpb.TimeSeries{
		Metric:   metric,
		Resource: resource,
		Points:   []*monitoringpb.Point{dataPoint},
	}
}

func (s *StackDriver) createTimeSeries(timeSeries []*monitoringpb.TimeSeries) error {
	if len(timeSeries) == 0 {
		return nil
	}
	return s.metricsClient.CreateTimeSeries(s.ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name:       stackmoni.MetricProjectPath(s.projectID),
		TimeSeries: timeSeries,
	})
}

func newDatapoint(when time.Time, d float64) *monitoringpb.Point {
	return &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			// for gauge metric StartTime must be the same as EndTime or zero
			EndTime: &googlepb.Timestamp{Seconds: when.Unix()},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: d},
		},
	}
}
