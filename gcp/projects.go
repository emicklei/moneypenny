package gcp

import (
	"context"
	"log"

	cloudresourcemanagerv1 "google.golang.org/api/cloudresourcemanager/v1"
)

func AllProjects(ctx context.Context) (list []string) {
	// https://cloud.google.com/resource-manager/reference/rest/v1/projects
	cloudresourcemanagerService, err := cloudresourcemanagerv1.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}
	pcall := cloudresourcemanagerService.Projects.List()
	presp, err := pcall.Do()
	if err != nil {
		log.Println("ABORT: error getting all projects", err)
		return
	}
	for _, each := range presp.Projects {
		list = append(list, each.ProjectId)
	}
	return
}
