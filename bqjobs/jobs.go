package bqjobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
)

const bytesProcessedThreshold = 1e+9 // 1 Gb TODO make flag

// Job information is available for a six month period after creation
// Requires the Can View project role, or the Is Owner project role if you set the allUsers property.
func queryAndAppend(ctx context.Context, client *bigquery.Client, project string, inserter *bigquery.Inserter, dryrun bool) {
	// query all jobs
	it := client.Jobs(ctx)
	it.ProjectID = project
	it.State = bigquery.Done
	today := time.Now()
	yesterday := today.Add(-24 * time.Hour)
	begin := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	end := begin.Add(24*time.Hour - 1*time.Nanosecond)
	it.MinCreationTime = begin
	it.MaxCreationTime = end
	it.AllUsers = true

	// iterate through results
	jobCount := 0
	insertCount := 0
	log.Println("quering job history of users in project", project)
	jobsToInsert := []BigQueryJob{}
	for {
		job, err := it.Next()
		if err != nil {
			break
		}
		jobCount++
		if job.LastStatus().Statistics.TotalBytesProcessed > bytesProcessedThreshold {
			bj := BigQueryJob{
				JobID:               job.ID(),
				Project:             project,
				Location:            job.Location(),
				Email:               job.Email(),
				TotalBytesProcessed: job.LastStatus().Statistics.TotalBytesProcessed,
				CreationTime:        job.LastStatus().Statistics.CreationTime,
				InsertionTime:       time.Now(),
				Query:               queryForJob(job),
			}
			jobsToInsert = append(jobsToInsert, bj)
			insertCount++
		}
	}
	timingOut, cancel := context.WithTimeout(ctx, 60*time.Second) // TODO add flag for this
	defer cancel()
	if dryrun {
		log.Println("skip inserting jobs because dryrun, count:", jobCount)
		return
	}
	log.Println("inserting jobs", insertCount, "with bytes processed larger than ", bytesProcessedThreshold/1e+9, "GB out of total jobs", jobCount)

	if err := inserter.Put(timingOut, jobsToInsert); err != nil {
		log.Printf("WARNING: inserting %d rows, error:%v\n", insertCount, err)
		batch := 100 // TODO add flag for this
		log.Println("try batching the rows but give up on the first error, batch size:", batch)
		for i := 0; i < len(jobsToInsert); i += batch {
			end := i + batch
			if end > len(jobsToInsert) {
				end = len(jobsToInsert)
			}
			subset := jobsToInsert[i:end]
			subTimeout, cancel := context.WithTimeout(ctx, 10*time.Second) // TODO add flag for this
			defer cancel()
			if err := inserter.Put(subTimeout, subset); err != nil {
				log.Printf("ERROR: inserting batch of %d rows, batch error:%v\n", len(subset), err)
				return
			}
			log.Println("completed insert of batch of rows:", len(subset))
		}
		return
	}
}

func queryForJob(j *bigquery.Job) string {
	cfg, err := j.Config()
	if err != nil {
		return fmt.Sprintf("// not available, error:%s", err.Error())
	}
	switch c := cfg.(type) {
	case *bigquery.CopyConfig:
		return fmt.Sprintf("// copy to %s.%s.%s", c.Dst.ProjectID, c.Dst.DatasetID, c.Dst.TableID)
	case *bigquery.ExtractConfig:
		return fmt.Sprintf("// extract from %s.%s.%s to %s", c.Src.ProjectID, c.Src.DatasetID, c.Src.TableID,
			c.Dst.URIs)
	case *bigquery.LoadConfig:
		return fmt.Sprintf("// load to %s.%s.%s from %v", c.Dst.ProjectID, c.Dst.DatasetID, c.Dst.TableID, c.Src)
	case *bigquery.QueryConfig:
		return c.Q
	default:
		return fmt.Sprintf("// not available, unknown type:%T", c)
	}
}
