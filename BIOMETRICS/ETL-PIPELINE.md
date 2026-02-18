# ETL-PIPELINE.md - Airflow Data Pipeline

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Data Engineering  
**Author:** BIOMETRICS Engineering Team

---

## 1. Overview

This document describes the ETL (Extract, Transform, Load) pipeline architecture using Apache Airflow for the BIOMETRICS platform. The pipeline handles data from multiple sources including biometric sensors, user interactions, and third-party integrations.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| Orchestrator | Apache Airflow 2.8+ | Workflow management |
| Metadata DB | PostgreSQL | DAG metadata storage |
| Executor | LocalExecutor | Task execution |
| UI | Airflow Webserver | Dashboard & monitoring |
| Scheduler | Airflow Scheduler | DAG scheduling |

### 2.2 Pipeline Flow

```
Sources → Extract → Transform → Load → Data Warehouse → Analytics
```

## 3. DAG Structure

### 3.1 Core DAGs

| DAG ID | Schedule | Description |
|--------|----------|-------------|
| biometric_daily_etl | 0 2 * * * | Daily biometric data processing |
| user_event_etl | */15 * * * * | Real-time user events |
| analytics_etl | 0 3 * * * | Analytics aggregation |
| reporting_etl | 0 4 * * * | Report generation |

### 3.2 DAG Example

```python
from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime, timedelta

default_args = {
    'owner': 'biometrics',
    'depends_on_past': False,
    'start_date': datetime(2026, 1, 1),
    'email_on_failure': True,
    'email_on_retry': False,
    'retries': 3,
    'retry_delay': timedelta(minutes=5),
}

with DAG(
    'biometric_daily_etl',
    default_args=default_args,
    schedule_interval='0 2 * * *',
    catchup=False,
    tags=['biometrics', 'etl'],
) as dag:

    extract_biometric_data = PythonOperator(
        task_id='extract_biometric_data',
        python_callable=extract_biometric_data,
    )

    transform_biometric_data = PythonOperator(
        task_id='transform_biometric_data',
        python_callable=transform_biometric_data,
    )

    load_to_warehouse = PythonOperator(
        task_id='load_to_warehouse',
        python_callable=load_to_warehouse,
    )

    validate_data_quality = PythonOperator(
        task_id='validate_data_quality',
        python_callable=validate_data_quality,
    )

    extract_biometric_data >> transform_biometric_data >> load_to_warehouse >> validate_data_quality
```

## 4. Task Operators

### 4.1 Python Operator

Used for custom data processing logic:

```python
def extract_biometric_data(**context):
    from src.extractors import BiometricExtractor
    extractor = BiometricExtractor()
    data = extractor.extract(
        start_date=context['ds'],
        end_date=context['tomorrow_ds']
    )
    context['task_instance'].xcom_push(key='biometric_data', value=data)
    return data
```

### 4.2 SQL Operator

Used for database transformations:

```python
transform_metrics = PostgresOperator(
    task_id='transform_metrics',
    postgres_conn_id='biometrics_warehouse',
    sql="""
        INSERT INTO daily_metrics (date, user_id, metric_type, value)
        SELECT 
            DATE(created_at) as date,
            user_id,
            metric_type,
            AVG(value) as value
        FROM raw_biometric_data
        WHERE DATE(created_at) = '{{ ds }}'
        GROUP BY DATE(created_at), user_id, metric_type
    """,
)
```

### 4.3 Bash Operator

Used for system commands and scripts:

```python
backup_data = BashOperator(
    task_id='backup_data',
    bash_command='python /scripts/backup_warehouse.py {{ ds }}',
)
```

## 5. Connections

### 5.1 Default Connections

| Conn ID | Type | Description |
|---------|------|-------------|
| biometrics_db | PostgreSQL | Main application database |
| biometrics_warehouse | PostgreSQL | Data warehouse |
| s3_etl_bucket | AWS S3 | ETL staging area |
| slack_notification | HTTP | Slack alerts |

### 5.2 Connection Configuration

```bash
airflow connections add 'biometrics_warehouse' \
    --conn-type 'postgres' \
    --conn-host 'room-03-postgres-master' \
    --conn-port '5432' \
    --conn-login 'biometrics_etl' \
    --conn-password 'secure_password' \
    --conn-schema 'warehouse'
```

## 6. Variables

### 6.1 Required Variables

| Variable | Value | Description |
|---------|-------|-------------|
| etl_staging_path | s3://biometrics-etl/staging | Staging directory |
| data_retention_days | 90 | Data retention period |
| batch_size | 10000 | Processing batch size |
| notification_channel | #etl-alerts | Slack channel |

## 7. Monitoring

### 7.1 Alerts

- Task failure → Slack notification
- DAG timeout → Email alert
- Data quality issues → PagerDuty escalation

### 7.2 Metrics

| Metric | Description |
|--------|-------------|
| task_duration | Time to complete each task |
| dag_success_rate | Success rate per DAG |
| data_volume | Records processed per run |
| error_rate | Failed records percentage |

## 8. Error Handling

### 8.1 Retry Strategy

```python
default_args = {
    'retries': 3,
    'retry_delay': timedelta(minutes=5),
    'retry_exponential_backoff': True,
    'max_retry_delay': timedelta(minutes=60),
}
```

### 8.2 Dead Letter Queue

Failed records are moved to a dead letter queue for manual review:

```python
def handle_failure(context):
    task_instance = context['task_instance']
    data = task_instance.xcom_pull(key='failed_data')
    # Push to DLQ
    dlq_client.push('etl_failures', data)
```

## 9. Deployment

### 9.1 Docker Compose

```yaml
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: airflow
      POSTGRES_PASSWORD: airflow
      POSTGRES_DB: airflow

  airflow-webserver:
    image: apache/airflow:2.8.0
    command: webserver
    environment:
      AIRFLOW__CORE__EXECUTOR: LocalExecutor
      AIRFLOW__DATABASE__SQL_ALCHEMY_CONN: postgresql+psycopg2://airflow:airflow@postgres/airflow
    volumes:
      - ./dags:/opt/airflow/dags
      - ./scripts:/opt/airflow/scripts
```

## 10. Security

### 10.1 RBAC

- Admin: Full access
- Editor: Edit DAGs, run tasks
- Viewer: Read-only access
- User: View own task runs

### 10.2 Secrets

All sensitive data stored in Vault or environment variables:

```python
import os
DB_PASSWORD = os.environ.get('AIRFLOW_DB_PASSWORD')
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
