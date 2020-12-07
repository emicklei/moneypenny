package project

import (
	"bytes"
	"encoding/json"
	"html/template"
	"testing"
	"time"

	"github.com/emicklei/moneypenny/model"
)

func TestReportJSON(t *testing.T) {
	data, _ := json.MarshalIndent(testReport(), "", "\t")
	if got, want := string(data), `{
	"last_day": {
		"consumption_day": "2020-10-15T00:00:00Z",
		"project_name": "project-name-test",
		"project_id": "project-id-test",
		"charges": 12.34,
		"credits": 0.01
	},
	"charges_percentage": 0,
	"mean": 8.32,
	"stddev": 1.323,
	"detector": "sundaysky{relativeThreshold=1.25,stddevThreshold=2.00,absoluteThreshold=10.00,windowDays=30}"
}`; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func testReport() ProjectStatsReport {
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
	return r
}

func TestReportTemplate(t *testing.T) {
	root := map[string]interface{}{
		"anomalies": []ProjectStatsReport{testReport(), testReport()},
	}
	data, _ := json.MarshalIndent(root, "", "\t")
	mapdata := map[string]interface{}{}
	err := json.Unmarshal(data, &mapdata)
	if err != nil {
		t.Fatal(err)
	}
	tmp, err := template.New("test").Parse(`
	{{range .anomalies}} 
		{{.detector}} 
		{{ .last_day.consumption_day }}
		{{ .last_day.project_name }}
		{{ .last_day.project_id }}
		{{ .last_day.charges }}
		{{ .last_day.credits }}
		{{.mean}} 
		{{.stddev}} 
	{{ end }}`)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	err = tmp.Execute(buf, mapdata)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buf.String())
}
