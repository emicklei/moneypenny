package bqjobs

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/emicklei/moneypenny/gcp"
	"github.com/emicklei/moneypenny/model"
	"github.com/emicklei/moneypenny/util"
	"github.com/urfave/cli/v2"
)

func CollectJobHistory(c *cli.Context, p model.Params) error {
	ctx := context.Background()
	tableID := c.String("bigquery-jobs-table")
	util.CheckBigQueryTable(tableID)
	parts := strings.Split(c.String("bigquery-jobs-table"), ".")
	tableProject := parts[0]
	tableDataset := parts[1]
	tableName := parts[2]
	log.Println("client project:", tableProject, "dataset:", tableDataset, "table:", tableName)

	client, err := bigquery.NewClient(ctx, tableProject)
	if err != nil {
		return err
	}
	defer client.Close()
	inserter := client.Dataset(tableDataset).Table(tableName).Inserter()
	for _, each := range gcp.AllProjects(ctx) {
		if gcp.IsBigQueryEnabled(ctx, each) {
			queryAndAppend(ctx, client, each, inserter, p.DryRun)
		}
	}
	return nil
}
