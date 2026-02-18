# DATA-LINEAGE.md - Data Lineage Tracking

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Data Engineering  
**Author:** BIOMETRICS Data Architecture Team

---

## 1. Overview

This document describes the data lineage tracking system for BIOMETRICS, enabling end-to-end visibility of data flow from source to consumption.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| Lineage Collector | Python | Collect metadata from sources |
| Lineage Store | PostgreSQL | Store lineage metadata |
| Lineage API | FastAPI | Query lineage data |
| Lineage UI | React | Visual lineage explorer |
| Schema Registry | Apache Atlas | Schema management |

### 2.2 Data Flow

```
Sources → Collectors → Lineage Store → API → UI
                ↓
          Schema Registry
```

## 3. Data Model

### 3.1 Core Entities

| Entity | Description | Attributes |
|--------|-------------|------------|
| Dataset | Logical data collection | name, schema, owner |
| Column | Dataset field | name, type, description |
| Job | Processing job | name, schedule, owner |
| Task | Job component | name, query, dependencies |
| Source | Data source | type, connection |
| User | Data consumer | email, role |

### 3.2 Tables

```sql
-- Datasets
CREATE TABLE datasets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    schema JSONB,
    owner_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(name)
);

-- Columns
CREATE TABLE columns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dataset_id UUID REFERENCES datasets(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(50),
    description TEXT,
    is_pii BOOLEAN DEFAULT FALSE,
    
    UNIQUE(dataset_id, name)
);

-- Jobs
CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    schedule VARCHAR(100),
    owner_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(name)
);

-- Task Lineage
CREATE TABLE task_lineage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID REFERENCES jobs(id) ON DELETE CASCADE,
    task_name VARCHAR(255),
    query TEXT,
    input_datasets JSONB,
    output_datasets JSONB,
    executed_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50),
    duration_ms INTEGER
);
```

## 4. Lineage Collection

### 4.1 SQL Lineage Parser

```python
import sqlparse
from sqlparse.sql import Statement, Token

class SQLLineageExtractor:
    """Extract lineage from SQL queries"""
    
    def parse(self, sql: str) -> LineageResult:
        """Parse SQL and extract input/output tables"""
        
        parsed = sqlparse.parse(sql)
        
        inputs = set()
        outputs = set()
        
        for statement in parsed:
            # Find SELECT statements (inputs)
            if statement.get_type() == 'SELECT':
                inputs.update(self._extract_tables(statement, 'FROM'))
                inputs.update(self._extract_tables(statement, 'JOIN'))
            
            # Find INSERT/UPDATE/DELETE (outputs)
            elif statement.get_type() in ['INSERT', 'UPDATE', 'DELETE']:
                outputs.add(self._extract_target_table(statement))
        
        return LineageResult(
            inputs=list(inputs),
            outputs=list(outputs),
        )
    
    def _extract_tables(self, statement: Statement, clause: str) -> Set[str]:
        """Extract table names from FROM/JOIN clause"""
        
        tables = set()
        
        # Use sqlparse to extract identifiers
        for token in statement.tokens:
            if token.ttype is None and token.value.upper() == clause:
                # Get next significant token(s)
                tables.add(self._get_table_name(token))
        
        return tables
```

### 4.2 Auto Lineage

```python
class AutoLineageCollector:
    """Automatically collect lineage from data pipeline"""
    
    def collect_from_airflow(self, dag_id: str) -> List[TaskLineage]:
        """Collect lineage from Airflow DAG"""
        
        lineage_records = []
        
        dag = self.airflow_client.get_dag(dag_id)
        
        for task in dag.tasks:
            # Get task metadata
            task_info = self.airflow_client.get_task_info(dag_id, task.task_id)
            
            # Extract lineage from task
            lineage = TaskLineage(
                job_id=self.get_or_create_job(dag_id, task.task_id),
                task_name=task.task_id,
                query=self._extract_query(task),
                input_datasets=self._get_input_datasets(task),
                output_datasets=self._get_output_datasets(task),
                executed_at=datetime.now(),
            )
            
            lineage_records.append(lineage)
        
        return lineage_records
```

## 5. Column-Level Lineage

### 5.1 Column Lineage Model

```python
@dataclass
class ColumnLineage:
    """Column-level lineage"""
    source_column: str
    target_column: str
    transformation: str
    task_id: str
    confidence: float

class ColumnLineageExtractor:
    """Extract column-level lineage from transformations"""
    
    def extract_from_query(self, query: str) -> List[ColumnLineage]:
        """Extract column lineage from SQL"""
        
        lineage = []
        
        # Parse query
        parsed = sqlparse.parse(query)[0]
        
        # Find SELECT clause
        select_clause = self._find_select_clause(parsed)
        
        for token in select_clause.tokens:
            if token.ttype is None:  # Identifier
                # Check for transformations
                if ' AS ' in token.value.upper():
                    source, target = token.value.split(' AS ')
                    lineage.append(ColumnLineage(
                        source_column=source.strip(),
                        target_column=target.strip(),
                        transformation='direct',
                        task_id=get_current_task_id(),
                        confidence=1.0,
                    ))
        
        return lineage
```

## 6. API Endpoints

### 6.1 Lineage API

```python
from fastapi import FastAPI, Query
from typing import List, Optional

app = FastAPI(title="Data Lineage API")

@app.get("/api/lineage/dataset/{dataset_name}")
def get_dataset_lineage(
    dataset_name: str,
    direction: str = Query("both", regex="^(upstream|downstream|both)$")
) -> dict:
    """Get full lineage for a dataset"""
    
    if direction in ['upstream', 'both']:
        upstream = query_upstream_lineage(dataset_name)
    else:
        upstream = []
    
    if direction in ['downstream', 'both']:
        downstream = query_downstream_lineage(dataset_name)
    else:
        downstream = []
    
    return {
        'dataset': dataset_name,
        'upstream': upstream,
        'downstream': downstream,
    }

@app.get("/api/lineage/column")
def get_column_lineage(
    dataset: str = Query(...),
    column: str = Query(...),
) -> List[ColumnLineage]:
    """Get column-level lineage"""
    
    return query_column_lineage(dataset, column)

@app.get("/api/lineage/job/{job_name}")
def get_job_lineage(job_name: str) -> dict:
    """Get lineage for a specific job"""
    
    return query_job_lineage(job_name)
```

## 7. Visualization

### 7.1 Lineage Graph

```javascript
// Lineage visualization component
const LineageGraph = ({ datasetName }) => {
  const [lineage, setLineage] = useState(null);
  
  useEffect(() => {
    fetchLineage(datasetName).then(setLineage);
  }, [datasetName]);
  
  const nodes = [
    ...lineage.upstream.map(d => ({ id: d, type: 'upstream' })),
    { id: datasetName, type: 'current' },
    ...lineage.downstream.map(d => ({ id: d, type: 'downstream' })),
  ];
  
  const edges = [
    ...lineage.upstream.map(d => ({ from: d, to: datasetName })),
    ...lineage.downstream.map(d => ({ from: datasetName, to: d })),
  ];
  
  return (
    <ReactFlow nodes={nodes} edges={edges}>
      <Controls />
    </ReactFlow>
  );
};
```

### 7.2 Data Flow Diagram

```
┌─────────────────┐
│  User Profile   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────────┐
│ Biometric Data │────►│  Daily Metrics  │
│   Collector     │     │    Aggregator   │
└─────────────────┘     └────────┬────────┘
                                  │
                                  ▼
                        ┌─────────────────┐
                        │   Analytics     │
                        │   Warehouse     │
                        └────────┬────────┘
                                 │
              ┌──────────────────┼──────────────────┐
              ▼                  ▼                  ▼
     ┌────────────────┐ ┌──────────────┐  ┌────────────────┐
     │  Metabase      │ │  ML Models   │  │   Reports     │
     │  Dashboards    │ │  Training    │  │   Generator   │
     └────────────────┘ └──────────────┘  └────────────────┘
```

## 8. Use Cases

### 8.1 Impact Analysis

```python
def analyze_impact(dataset_name: str) -> ImpactReport:
    """Analyze impact of changes to a dataset"""
    
    downstream = get_downstream_datasets(dataset_name)
    downstream_jobs = get_downstream_jobs(dataset_name)
    downstream_dashboards = get_downstream_dashboards(dataset_name)
    
    report = ImpactReport(
        affected_datasets=len(downstream),
        affected_jobs=len(downstream_jobs),
        affected_dashboards=len(downstream_dashboards),
        risk_score=calculate_risk_score(downstream),
        recommendations=get_migration_recommendations(dataset_name),
    )
    
    return report
```

### 8.2 Root Cause Analysis

```python
def find_root_cause(issue_dataset: str, error_time: datetime) -> RootCause:
    """Find root cause of data quality issue"""
    
    # Get all upstream dependencies
    upstream = get_upstream_lineage(issue_dataset)
    
    # Check each upstream for issues
    for dataset in upstream:
        issues = check_data_quality(dataset, error_time)
        
        if issues:
            return RootCause(
                source_dataset=dataset,
                issues=issues,
                propagation_path=get_propagation_path(dataset, issue_dataset),
            )
    
    return RootCause(source_dataset=None, issues=None)
```

## 9. Governance

### 9.1 Ownership

| Dataset | Owner | Steward |
|---------|-------|---------|
| Raw data | Engineering | Data Team |
| Processed data | Data Engineering | Analytics |
| Business metrics | Product | Business |

### 9.2 Glossary

| Term | Definition |
|------|------------|
| Upstream | Data sources that flow into a dataset |
| Downstream | Datasets that consume from a given dataset |
| Column lineage | Flow of data at column level |
| Task lineage | Flow through ETL/ELT tasks |

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
