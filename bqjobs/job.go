package bqjobs

import "time"

type BigQueryJob struct {
	JobID               string    `bigquery:"job_id"`
	Project             string    `bigquery:"project"`
	Location            string    `bigquery:"location"`
	Email               string    `bigquery:"email"`
	CreationTime        time.Time `bigquery:"creation_time"`
	InsertionTime       time.Time `bigquery:"insertion_time"`
	TotalBytesProcessed int64     `bigquery:"total_bytes_processed"`
}
