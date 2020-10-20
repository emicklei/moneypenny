# moneypenny

Tool for reporting and monitoring Google Cloud Platform costs.

## detect-project-cost-anomalies

For each project linked to a billing account, the last day cost is compared to the cost history of the past 30 days using the mean and standard deviation value. See SundaySky:IsAnomaly for the algorithm.

```bash
moneypenny \
    -billing-table PROJECT.DATASET.gcp_billing_export_v1_00000000 \
    detect-project-cost-anomalies
```

This command will produce a `DetectProjectCostAnomalies.json` file that looks like:

```json
    { "anomalies": [
        {
            "last_day": {
                "consumption_day": "2020-10-15T00:00:00Z",
                "project_name": "project-name-test",
                "project_id": "project-id-test",
                "charges": 12.34,
                "credits": 0.01
            },
            "mean": 8.32,
            "stddev": 1.323,
            "detector": "sundaysky{relativeThreshold=1.25,stddevThreshold=2.00,absoluteThreshold=10.00}"
        }]}
```

&copy; 2020, MIT Licensed. http://ernestmicklei.com