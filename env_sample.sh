# github
export GITHUB_USER="your-name"
export GITHUB_TOKEN="your-token"

# DB
export DB_MASTER_HOST="db"
export DB_MASTER_USER="hoge"
export DB_MASTER_PASSWORD="moge"
export DB_MASTER_DATABASE="mimosa"
export DB_SLAVE_HOST="db"
export DB_SLAVE_USER="hoge"
export DB_SLAVE_PASSWORD="moge"
export DB_LOG_MODE="true"
export DB_PORT="3306"
export DB_SCHEMA="mimosa"

# AWS
export AWS_REGION="ap-northeast-1"
export AWS_ACCESS_KEY_ID="hogehoge"
export AWS_SECRET_ACCESS_KEY="hugahuga"
export AWS_SESSION_TOKEN="piyopiyo"
export PRIVATE_EXPOSE_QUEUE_NAME="osint-privateexpose"
export PRIVATE_EXPOSE_QUEUE_URL="http://sqs:9324/queue/osint-privateexpose"
export SUBDOMAIN_QUEUE_NAME="osint-subdomain"
export SUBDOMAIN_QUEUE_URL="http://sqs:9324/queue/osint-subdomain"
export SQS_ENDPOINT="http://sqs:9324"

# mimosa
export FINDING_SVC_ADDR="finding:8001"
export ALERT_SVC_ADDR="alert:8004"
export OSINT_SVC_ADDR="osint:18081"

# theHarvester
export RESULT_PATH="/tmp"
export HARVESTER_PATH="/opt/git/mimosa-osint/src/private-expose/theHarvester"