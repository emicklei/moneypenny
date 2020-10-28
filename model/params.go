package model

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

const TimestampDayLayout = "2006-01-02"

// Params captures all program flags
type Params struct {
	// JSON tag names must match flag option names
	Opex       string `json:"opex,omitempty" `
	DayInMonth int    `json:"day" `
	MonthIndex int    `json:"month" `
	Year       int    `json:"year" `

	// Billing
	BillingTableFQN      string `json:"billing-table,omitempty" `
	QueryExecutionRegion string `json:"query-execution-region,omitempty" `
	QueryExecutionOpex   string `json:"query-execution-opex,omitempty" `
	// Target
	TargetTableFQN         string `json:"target-table,omitempty" `
	TargetDatasetRegion    string `json:"target-region,omitempty" `
	TargetMetricsProjectID string `json:"metrics-project-id,omitempty" `

	DryRun  bool `json:"dryrun" `
	Verbose bool `json:"v" `
}

func ParamsFromContext(c *cli.Context) Params {
	p := Params{
		DryRun:  c.Bool("dryrun"),
		Verbose: c.Bool("v"),
		// When
		DayInMonth: c.Int("day"),
		MonthIndex: c.Int("month"),
		Year:       c.Int("year"),
		// Billing
		BillingTableFQN:      c.String("billing-table"),
		QueryExecutionRegion: c.String("query-execution-region"),
		// others
		TargetTableFQN: c.String("target-table"),
		Opex:           c.String("opex"),
	}
	if date := c.String("date"); len(date) > 0 {
		d, err := time.Parse(TimestampDayLayout, date)
		if err != nil {
			log.Fatalf("unable to parse date [%s], expect format [%s]\n", date, TimestampDayLayout)
		}
		p.DayInMonth = d.Day()
		p.MonthIndex = int(d.Month())
		p.Year = d.Year()
	}
	return p
}

func (p Params) JSON() string {
	data, _ := json.MarshalIndent(p, "", "\t")
	return string(data)
}

func (p Params) TargetProjectID() string {
	s := strings.Split(p.TargetTableFQN, ".")
	if len(s) != 3 {
		log.Fatalln("flag -target-table format must be PROJECT.DATASET.TABLE , got:", p.TargetTableFQN)
	}
	return s[0]
}

func (p Params) TargetDatasetID() string {
	s := strings.Split(p.TargetTableFQN, ".")
	if len(s) != 3 {
		log.Fatalln("flag -target-table format must be PROJECT.DATASET.TABLE , got:", p.TargetTableFQN)
	}
	return s[1]
}

func (p Params) TargetTableID() string {
	s := strings.Split(p.TargetTableFQN, ".")
	if len(s) != 3 {
		log.Fatalln("flag -target-table format must be PROJECT.DATASET.TABLE , got:", p.TargetTableFQN)
	}
	return s[2]
}

func (p Params) BillingProjectID() string {
	s := strings.Split(p.BillingTableFQN, ".")
	if len(s) != 3 {
		log.Fatalln("flag -billing-table format must be PROJECT.DATASET.TABLE , got:", p.BillingTableFQN)
	}
	return s[0]
}

func (p Params) Date() time.Time {
	return time.Date(p.Year, time.Month(p.MonthIndex), p.DayInMonth, 0, 0, 0, 0, time.Local)
}
