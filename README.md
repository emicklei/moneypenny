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

### Optional flags for detect-project-cost-anomalies

Add this option to override the threshold (2.0) of the sundaysky detection method.

```bash
    -sundaysky.stddev=3.0
```

Add this option to store all anomaly events in a BigQuery table. See `infra.sh` how to create it.

```bash
    -target-table PROJECT.moneypenny_dataset.moneypenny_cost_anomaly_events
```

&copy; 2020, MIT Licensed. [ernestmicklei.com](http://ernestmicklei.com)
