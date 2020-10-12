package gcp

import (
	"context"
	"log"
	"strings"

	"google.golang.org/api/serviceusage/v1"
)

// IsBigQueryEnabled returns wheter a GCP project has bigquery.googleapis.com enabled.
func IsBigQueryEnabled(ctx context.Context, project string) bool {
	serviceusageService, err := serviceusage.NewService(ctx)
	call := serviceusageService.Services.List("projects/" + project).Filter("state:ENABLED")
	resp, err := call.Do()
	if err != nil {
		log.Println("ABORT: error getting enabled services for", project, err)
		return false
	}
	for _, each := range resp.Services {
		if strings.HasSuffix(each.Name, "bigquery.googleapis.com") {
			return true
		}
	}
	return false
}
