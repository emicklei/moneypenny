package opex

import (
	"fmt"

	"github.com/emicklei/moneypenny/util"
)

func queryCostPerOpex(fqDatasetTableID string, year, monthindex int, opex string) string {
	util.CheckBigQueryTable(fqDatasetTableID)
	util.CheckMonth(monthindex)
	util.CheckOpex(opex)
	return fmt.Sprintf(`
# Moneypenny - queryCostPerOpex
#
# Author: EMicklei
# Params: fqDatasetTableID,year,monthindex,opex
# Output: charges,project_id,project_name,gcp_service,credits
#
SELECT
  ROUND(SUM(cost), 2) AS charges,
  IFNULL(ROUND(SUM((SELECT SUM(amount) FROM UNNEST(credits))),2), 0) as credits,
  project.id AS project_id,
  project.name AS project_name,  
  service.description AS gcp_service
FROM `+"`%s`,"+`
  UNNEST(labels) AS label
WHERE
  EXTRACT(YEAR
  FROM
    _PARTITIONTIME) = %d
  AND EXTRACT(MONTH
      FROM
        _PARTITIONTIME) = %d       
  AND label.key = "opex"
  AND label.value = "%s"
GROUP BY
  project_id,
  project_name,
  gcp_service
ORDER BY
  charges DESC
`, fqDatasetTableID, year, monthindex, opex)
}
