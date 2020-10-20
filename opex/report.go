package opex

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/emicklei/moneypenny/model"
)

func writeDetailReport(input model.Params, cc model.CostComputation) error {
	byTeam := map[string][]model.LabeledCost{}
	for _, each := range cc.Lines {
		eachCost := model.LabeledCostFrom(each)
		if eachCost.Charges < 0.01 {
			continue
		}
		list, ok := byTeam[eachCost.Opex.StringVal]
		if ok {
			byTeam[eachCost.Opex.StringVal] = append(list, eachCost)
		} else {
			byTeam[eachCost.Opex.StringVal] = []model.LabeledCost{eachCost}
		}
	}

	for opex, lines := range byTeam {
		file := fmt.Sprintf("%s-%d-%d-component-breakdown.json", opex, input.Year, input.MonthIndex)
		out, err := os.Create(file)
		if err != nil {
			return err
		}
		defer out.Close()
		if abs, err := filepath.Abs(file); err == nil {
			log.Println("writing report", abs)
		}

		// TODO create struct for this?
		doc := map[string]interface{}{}
		doc["bigquery"] = map[string]interface{}{
			"bytes_processed":   cc.ByteProcessed,
			"execution_time_ms": cc.ExecutionTime.Milliseconds(),
			"query":             cc.Query,
		}
		doc["input"] = input
		doc["output"] = lines
		doc["report_date"] = time.Now()
		enc := json.NewEncoder(out)
		enc.SetIndent("", "\t")
		if err := enc.Encode(doc); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
