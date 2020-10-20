package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/emicklei/moneypenny/alert"
	"github.com/emicklei/moneypenny/bqjobs"
	"github.com/emicklei/moneypenny/model"
	"github.com/emicklei/moneypenny/opex"
	"github.com/emicklei/moneypenny/project"
	"github.com/urfave/cli/v2"
)

var version = time.Now().String()

func main() {
	if err := newApp().Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Version = version
	app.EnableBashCompletion = true
	app.Name = "moneypenny"
	app.Usage = `Google Cloud Platform cost reporting tool

	see https://github.com/emicklei/moneypenny for documentation.
`
	// override -v
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "v",
			Usage: "verbose logging",
		},
		&cli.BoolFlag{
			Name:  "dryrun",
			Usage: "only log and do not perform any queries",
		},
		&cli.StringFlag{
			Name:  "billing-table",
			Usage: "project.dataset.table, full qualified identifier",
		},
		&cli.IntFlag{
			Name:  "day",
			Value: int(time.Now().Day()),
			Usage: "DD",
		},
		&cli.IntFlag{
			Name:  "month",
			Value: int(time.Now().Month()),
			Usage: "MM, month index [1..12]",
		},
		&cli.IntFlag{
			Name:  "year",
			Value: int(time.Now().Year()),
			Usage: "YYYY, e.g year since billing records are available",
		},
		&cli.StringFlag{
			Name:  "date",
			Usage: "YYYY-MM-DD",
		},
		&cli.StringFlag{
			Name:  "opex",
			Usage: "team or person that has operation expenditure responsibility",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:  "detect-project-cost-anomalies",
			Usage: "detect-project-cost-anomalies",
			Action: func(c *cli.Context) error {
				p := model.ParamsFromContext(c)
				defer logBegin(c)()
				return logEnd(c, project.DetectProjectCostAnomalies(c, p))
			},
		},
		{
			Name:  "cost-per-opex",
			Usage: "cost-per-opex",
			Action: func(c *cli.Context) error {
				p := model.ParamsFromContext(c)
				defer logBegin(c)()
				return logEnd(c, opex.ReportCostPerOpex(c, p))
			},
		},
		{
			Name:  "update-bigquery-jobs",
			Usage: "update-bigquery-jobs",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "bigquery-jobs-table",
					Usage: "project.dataset.table, full qualified identifier",
				},
				&cli.StringFlag{
					Name:  "project-id",
					Usage: "if empty then update for all visible projects",
				},
			},
			Action: func(c *cli.Context) error {
				p := model.ParamsFromContext(c)
				defer logBegin(c)()
				return logEnd(c, bqjobs.CollectJobHistory(c, p))
			},
		},
		{
			Name:  "measure-cost-per-opex-yesterday",
			Usage: "Publish metrics to Google Monitoring about the cost per opex of yesterday",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "metrics-project",
					Usage: "GCP project id",
				},
			},
			Action: func(c *cli.Context) error {
				p := model.ParamsFromContext(c)
				defer logBegin(c)()
				return logEnd(c, opex.MeasureCostPerOpexLastDay(c, p))
			},
		},
		{
			Name:  "send-email",
			Usage: "Send a HTML email by processing a Go template with a JSON document",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:   "api-key",
					Usage:  "sendgrid.com API Key",
					Hidden: true,
				},
				&cli.StringFlag{
					Name:  "from",
					Usage: "from email address",
				},
				&cli.StringFlag{
					Name:  "to",
					Usage: "to email address",
				},
				&cli.StringFlag{
					Name:  "subject",
					Usage: "email subject",
				},
				&cli.StringFlag{
					Name:  "html-template-file",
					Usage: "file with Go template to produce HTML",
				},
				&cli.StringFlag{
					Name:  "json-file",
					Usage: "file with JSON document to process with template",
				},
			},
			Action: func(c *cli.Context) error {
				defer logBegin(c)()
				return logEnd(c, alert.SendEmail(c.String("subject"),
					c.String("from"),
					c.String("to"),
					c.String("json-file"),
					c.String("html-template-file"),
					c.String("api-key")))
			},
		},
	}
	return app
}

func logBegin(c *cli.Context) func() {
	buf := new(bytes.Buffer)
	fmt.Fprint(buf, "[moneypenny] ")

	appendFlag := func(each cli.Flag) {
		fv := reflect.ValueOf(each)
		hide := reflect.Indirect(fv).FieldByName("Hidden").Bool()
		name := each.Names()[0]
		value := c.Generic(name)
		if hide {
			value = "**hidden**"
		}
		fmt.Fprintf(buf, " %s=%v", name, value)
	}
	for _, each := range c.App.Flags {
		appendFlag(each)
	}
	fmt.Fprintf(buf, " %s", c.Command.Name)
	for _, each := range c.Command.Flags {
		appendFlag(each)
	}
	log.Println(buf.String())
	return func() {
		if err := recover(); err != nil {
			// no way to communicate error to cli so exit here.
			log.Fatalln(c.Command.Name, "failed because:", err)
		}
	}
}

func logEnd(c *cli.Context, err error) error {
	if err != nil {
		log.Printf("[moneypenny] %s failed:%v\n", c.Command.Name, err)
	} else {
		log.Printf("[moneypenny] %s done\n", c.Command.Name)
	}
	return err
}
