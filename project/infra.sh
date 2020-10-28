# example scripts to prepare GCP infrastructure

bq \
    --project_id=$PROJECT \
    --location=eu \
    mk --dataset \
    --label env:prd \
    --label opex:team-gcp-cost \
    --label service:moneypenny \
    --description "Used for moneypenny to compute costs and alert about anomalies" \
    $PROJECT:moneypenny_dataset

bq \
    --project_id=$PROJECT \
    --location=eu \
    mk --table \
    --time_partitioning_field event_creation_time \
    --label env:prd \
    --label opex:team-gcp-cost \
    --label service:moneypenny \
    $PROJECT:moneypenny_dataset.moneypenny_cost_anomaly_events event_id:STRING,event_creation_time:TIMESTAMP,project_id:STRING,project_name:STRING,charges:FLOAT64,charges_percentage:FLOAT64,credits:FLOAT64,mean:FLOAT64,stddev:FLOAT64,detection_day:TIMESTAMP,detector:STRING
                