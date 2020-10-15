package project

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/emicklei/moneypenny/model"
)

func TestReportJSON(t *testing.T) {
	d, _ := time.Parse(model.TimestampDayLayout, "2020-10-15")
	r := ProjectStatsReport{
		LastDay: DailyCost{
			Day:         d,
			Charges:     12.34,
			Credits:     0.01,
			ProjectID:   "project-id-test",
			ProjectName: "project-name-test",
		},
		Detector:          BestSundaySky.String(),
		Mean:              8.32,
		StandardDeviation: 1.3230,
	}
	data, _ := json.MarshalIndent(r, "", "\t")
	if got, want := string(data), `{
	"last-day": {
		"consumption-day": "2020-10-15T00:00:00Z",
		"project-name": "project-name-test",
		"project-id": "project-id-test",
		"charges": 12.34,
		"credits": 0.01
	},
	"mean": 8.32,
	"stddev": 1.323,
	"detector": "sundaysky{relativeThreshold=1.25,stddevThreshold=2.00,absoluteThreshold=10.00}"
}`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}
