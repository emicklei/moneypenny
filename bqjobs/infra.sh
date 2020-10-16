# example scripts to prepara GCP infrastructure

bq \
    --project_id=$PROJECT \
    --location=eu \
    mk --dataset \
    --label env:prd \
    --label opex:guild-venom \
    --label service:moneypenny \
    --description "Used for moneypenny to compute BigQuery Job costs per email" \
    $PROJECT:moneypenny_dataset

bq \
    --project_id=$PROJECT \
    --location=eu \
    mk --table \
    --time_partitioning_field creation_time \
    --label env:prd \
    --label opex:guild-venom \
    --label service:moneypenny \
    $PROJECT:moneypenny_dataset.moneypenny_bigquery_job_history job_id:STRING,project:STRING,email:STRING,creation_time:TIMESTAMP,insertion_time:TIMESTAMP,total_bytes_processed:NUMERIC,location:STRING,query:STRING,query_hash:STRING
