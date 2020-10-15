package bqjobs

import "time"

// BigQueryJob represents a row in $PROJECT:moneypenny_dataset.moneypenny_bigquery_job_history
type BigQueryJob struct {
	JobID               string    `bigquery:"job_id"`
	Project             string    `bigquery:"project"`
	Location            string    `bigquery:"location"`
	Email               string    `bigquery:"email"`
	CreationTime        time.Time `bigquery:"creation_time"`
	InsertionTime       time.Time `bigquery:"insertion_time"`
	TotalBytesProcessed int64     `bigquery:"total_bytes_processed"`
	Query               string    `bigquery:"query"`
}
