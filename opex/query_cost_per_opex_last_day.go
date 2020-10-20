package opex

import (
	"fmt"
	"time"

	"github.com/emicklei/moneypenny/model"
	"github.com/emicklei/moneypenny/util"
)

func queryCostPerOpexLastDay(fqDatasetTableID string, date time.Time, opex string) string {
	util.CheckBigQueryTable(fqDatasetTableID)
	util.CheckOpex(opex)
	return fmt.Sprintf(`
# Moneypenny - queryCostPerOpexLastDay
#
# Author: EMicklei
# Params: fqDatasetTableID,dayString,opex
# Output: charges,project_name,project_id,gcp_service,credits
#
SELECT
  ROUND(SUM(cost), 2) AS charges,
  IFNULL(ROUND(SUM((SELECT SUM(amount) FROM UNNEST(credits))),2), 0) as credits,
  project.name AS project_name,
  project.id AS project_id,
  service.description AS gcp_service
FROM `+"`%s`,"+`
  UNNEST(labels) AS label
WHERE
  _PARTITIONTIME = TIMESTAMP('%s')      
  AND label.key = "opex"
  AND label.value = "%s"
GROUP BY
  project_id,
  project_name,
  gcp_service
ORDER BY
  charges DESC
`, fqDatasetTableID, date.Format(model.TimestampDayLayout), opex)
}
