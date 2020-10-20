package opex

import (
	"fmt"

	"github.com/emicklei/moneypenny/util"
)

func queryCostPerComponent(fqDatasetTableID string, year, monthindex int) string {
	util.CheckBigQueryTable(fqDatasetTableID)
	util.CheckMonth(monthindex)
	return fmt.Sprintf(`
# Moneypenny - queryCostPerComponent
#
# Author: E.Micklei
# Params: fqDatasetTableID,year,monthindex
# Output: charges,credits,project_id,project_name,gcp_service,component,service,opex
#
SELECT
  ROUND(SUM(cost), 2) AS charges,
  IFNULL(ROUND(SUM((SELECT SUM(amount) FROM UNNEST(credits))),2), 0) as credits,
  project.name AS project_name,
  project.id AS project_id,
  bill.service.description AS gcp_service,
  (SELECT value FROM UNNEST(bill.labels) WHERE key = "component") AS component,
  (SELECT value FROM UNNEST(bill.labels) WHERE key = "service") AS service,
  (SELECT value FROM UNNEST(bill.labels) WHERE key = "opex") AS opex
FROM `+"`%s` bill"+`
WHERE
  EXTRACT(YEAR
  FROM
    _PARTITIONTIME) = %d
  AND EXTRACT(MONTH
      FROM
        _PARTITIONTIME) = %d      
  AND ARRAY_LENGTH(bill.labels) != 0
  AND (SELECT value FROM UNNEST(bill.labels) WHERE key = "opex") != ""
GROUP BY
  project_id,
  project_name,
  gcp_service,
  component,
  service,
  opex
ORDER BY
  charges DESC
`, fqDatasetTableID, year, monthindex)
}
