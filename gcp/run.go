package gcp

import (
	"context"
	"log"

	"cloud.google.com/go/bigquery"
	"github.com/emicklei/moneypenny/model"
	"google.golang.org/api/iterator"
)

func RunBigQuery(ctx context.Context, p model.Params, query string) (model.CostComputation, error) {
	cc := model.CostComputation{Lines: []map[string]bigquery.Value{}}
	client, err := bigquery.NewClient(ctx, p.BillingProjectID()) // assume execution project == billing record project
	cc.Query = query
	if p.Verbose {
		log.Println(query)
	}
	q := client.Query(query)
	q.Location = p.QueryExecutionRegion
	q.Labels = map[string]string{
		"opex":    p.QueryExecutionOpex,
		"service": "moneypenny",
	}
	job, err := q.Run(ctx)
	if err != nil {
		return cc, err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return cc, err
	}
	if err := status.Err(); err != nil {
		return cc, err
	}
	it, err := job.Read(ctx)
	for {
		var line map[string]bigquery.Value
		err := it.Next(&line)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("%#v, error:%v", line, err)
		} else {
			cc.Lines = append(cc.Lines, line)
		}
	}
	stats, _ := job.Status(ctx)
	cc.ByteProcessed = stats.Statistics.TotalBytesProcessed
	cc.ExecutionTime = stats.Statistics.EndTime.Sub(stats.Statistics.StartTime)
	if p.Verbose {
		log.Printf("%#v\n", stats.Statistics)
	}
	return cc, nil
}
