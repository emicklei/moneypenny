package opex

import (
	"github.com/emicklei/moneypenny/gcp"
	"github.com/emicklei/moneypenny/model"
	"github.com/emicklei/moneypenny/util"
	"github.com/emicklei/tre"
	"github.com/urfave/cli/v2"
)

func ReportCostPerOpex(c *cli.Context, p model.Params) error {
	cc, err := computeCostPerOpex(p)
	if err != nil {
		return tre.New(err, "computeCostPerOpex")
	}
	return writeDetailReport(p, cc)
}

func ReportCostPerComponent(c *cli.Context, p model.Params) error {
	cc, err := computeCostPerComponent(p)
	if err != nil {
		return tre.New(err, "computeCostPerComponent")
	}
	return writeDetailReport(p, cc)
}

func MeasureCostPerOpexLastDay(c *cli.Context, p model.Params) error {
	metricsProjectID := c.String("metrics-project")
	util.CheckNonEmpty("metrics-project", metricsProjectID)

	cc, err := computeCostPerOpexLastDay(p)
	if err != nil {
		return tre.New(err, "computeCostPerOpexLastDay")
	}

	sd, err := gcp.NewStackDriver(metricsProjectID)
	if err != nil {
		return err
	}
	defer sd.Close()
	lines := []model.LabeledCost{}
	for _, each := range cc.Lines {
		lines = append(lines, model.LabeledCostFrom(each))
	}
	if err := sd.SendMetrics(p.Opex, lines, p.DryRun); err != nil {
		return err
	}
	return nil
}
