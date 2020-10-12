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
# Output: charges,project,gcp_service
#
SELECT
  ROUND(SUM(cost), 2) AS charges,
  project.name AS project,
  service.description AS gcp_service
FROM `+"`%s`,"+`
  UNNEST(labels) AS label
WHERE
  _PARTITIONTIME = TIMESTAMP('%s')      
  AND label.key = "opex"
  AND label.value = "%s"
GROUP BY
  project,
  gcp_service
ORDER BY
  charges DESC
`, fqDatasetTableID, date.Format(model.TimestampDayLayout), opex)
}
