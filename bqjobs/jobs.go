package bqjobs

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
)

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
		if job.LastStatus().Statistics.TotalBytesProcessed > 0 {
			bj := BigQueryJob{
				JobID:               job.ID(),
				Project:             project,
				Location:            job.Location(),
				Email:               job.Email(),
				TotalBytesProcessed: job.LastStatus().Statistics.TotalBytesProcessed,
				CreationTime:        job.LastStatus().Statistics.CreationTime,
				InsertionTime:       time.Now(),
			}
			jobsToInsert = append(jobsToInsert, bj)
			insertCount++
		}
	}
	timingOut, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if dryrun {
		log.Println("skip inserting jobs because dryrun, count:", jobCount)
		return
	}
	log.Println("inserting jobs", insertCount, "out of", jobCount)
	if err := inserter.Put(timingOut, jobsToInsert); err != nil {
		log.Println("ERROR: inserting rows", err)
		return
	}
}
