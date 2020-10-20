package project

import (
	"fmt"
	"time"

	"github.com/emicklei/moneypenny/model"
)

// QueryPastDays returns the query that collects daily cost records between start and end (both including).
func QueryPastDays(fqDatasetTableID string, start, end time.Time) string {
	return fmt.Sprintf(`#standardSQL
   # Moneypenny - dailyproject.QueryPastDays
   #
   # Author: EMicklei
   # Params: fqDatasetTableID,start,end
   # Output: consumption_day,project_name,project_id,charges,credits
   #   
SELECT
   _PARTITIONTIME as consumption_day,
   project.name as project_name,
   project.id as project_id,
   ROUND(SUM(cost), 2) as charges,
   IFNULL(ROUND(SUM((SELECT SUM(amount) FROM UNNEST(credits))),2), 0) as credits
FROM `+"`%s`"+`
WHERE 
  project.id IS NOT NULL
  AND _PARTITIONTIME >= TIMESTAMP("%s")
  AND _PARTITIONTIME <= TIMESTAMP("%s")
GROUP BY consumption_day, project_name, project_id
ORDER BY consumption_day DESC
`, fqDatasetTableID, start.Format(model.TimestampDayLayout), end.Format(model.TimestampDayLayout))
}
