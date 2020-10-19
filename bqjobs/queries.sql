with
    total_costs
    as
    (
        SELECT
            ROUND( SUM(5 * total_bytes_processed/POWER(2,40)),2) AS total_cost
        FROM
            `dataset.moneypenny_bigquery_job_history`
)
, user_project_costs as
(
SELECT
    email,
    project,
    ROUND( SUM(5 * total_bytes_processed/POWER(2,40)),2) AS cost,
    ROUND( SUM(total_bytes_processed/POWER(2,30)),2) AS GB,
    COUNT(*) AS job_count
FROM
    `dataset.moneypenny_bigquery_job_history`
 GROUP BY
  project,
  email
ORDER BY
  cost DESC
)
select
    email
, project
, cost
, round(cost / tc.total_cost, 3) as relative_cost
, GB
, job_count
from user_project_costs
    join total_costs tc on true
order by cost desc