# DATA-QUALITY.md - Data Quality Framework

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Data Engineering  
**Author:** BIOMETRICS Data Team

---

## 1. Overview

This document outlines the data quality framework for BIOMETRICS, ensuring all data meets defined standards for accuracy, completeness, consistency, and timeliness.

## 2. Quality Dimensions

### 2.1 Core Dimensions

| Dimension | Description | Threshold |
|-----------|-------------|----------|
| Accuracy | Data matches real-world values | ≥ 99% |
| Completeness | No missing values | ≥ 98% |
| Consistency | Data uniform across systems | ≥ 99% |
| Timeliness | Data available when needed | ≥ 95% |
| Validity | Data conforms to rules | ≥ 99% |
| Uniqueness | No duplicate records | ≥ 99.5% |

### 2.2 Quality Rules

```python
class DataQualityRules:
    """Define data quality validation rules"""
    
    RULES = {
        'biometric_heart_rate': {
            'min_value': 30,
            'max_value': 220,
            'nullable': False,
            'data_type': 'integer',
        },
        'biometric_blood_pressure_systolic': {
            'min_value': 60,
            'max_value': 250,
            'nullable': False,
            'data_type': 'integer',
        },
        'biometric_blood_pressure_diastolic': {
            'min_value': 40,
            'max_value': 150,
            'nullable': False,
            'data_type': 'integer',
        },
        'user_email': {
            'pattern': r'^[\w\.-]+@[\w\.-]+\.\w+$',
            'nullable': False,
            'data_type': 'string',
        },
    }
```

## 3. Validation Framework

### 3.1 Great Expectations Integration

```python
import great_expectations as ge

def validate_biometric_data(df: pd.DataFrame) -> dict:
    """Validate biometric data using Great Expectations"""
    
    # Create expectations
    df_ge = ge.from_pandas(df)
    
    # Heart rate validation
    df_ge.expect_column_values_to_be_between(
        column='heart_rate',
        min_value=30,
        max_value=220,
    )
    
    # Blood pressure validation
    df_ge.expect_column_values_to_be_between(
        column='blood_pressure_systolic',
        min_value=60,
        max_value=250,
    )
    
    # Email format validation
    df_ge.expect_column_values_to_match_regex(
        column='email',
        regex=r'^[\w\.-]+@[\w\.-]+\.\w+$',
    )
    
    # Not null validation
    df_ge.expect_column_values_to_not_be_null(column='user_id')
    df_ge.expect_column_values_to_not_be_null(column='timestamp')
    
    # Run validation
    results = df_ge.validate()
    return results
```

### 3.2 Custom Validators

```python
from dataclasses import dataclass
from typing import List, Optional

@dataclass
class ValidationResult:
    """Result of a validation check"""
    rule_name: str
    passed: bool
    failed_count: int
    total_count: int
    sample_failures: Optional[List] = None

class BiometricValidator:
    """Custom validators for biometric data"""
    
    def validate_heart_rate_consistency(self, df: pd.DataFrame) -> ValidationResult:
        """Validate heart rate doesn't change unrealistically fast"""
        df = df.sort_values('timestamp')
        df['rate_change'] = df['heart_rate'].diff().abs()
        
        # Heart rate shouldn't change more than 30 bpm per minute
        time_diff = df['timestamp'].diff().dt.total_seconds() / 60
        max_change = 30 * time_diff
        
        failures = df[df['rate_change'] > max_change]
        
        return ValidationResult(
            rule_name='heart_rate_consistency',
            passed=len(failures) == 0,
            failed_count=len(failures),
            total_count=len(df),
            sample_failures=failures.head(10).to_dict('records'),
        )
    
    def validate_bmi_range(self, df: pd.DataFrame) -> ValidationResult:
        """Validate BMI is in valid range"""
        bmi = df['weight_kg'] / ((df['height_cm'] / 100) ** 2)
        failures = df[(bmi < 10) | (bmi > 60)]
        
        return ValidationResult(
            rule_name='bmi_range',
            passed=len(failures) == 0,
            failed_count=len(failures),
            total_count=len(df),
        )
```

## 4. Data Profiling

### 4.1 Profiling Tasks

| Task | Frequency | Description |
|------|----------|-------------|
| Column profiling | Daily | Statistics per column |
| Distribution analysis | Weekly | Value distribution |
| Correlation analysis | Weekly | Cross-column relationships |
| Anomaly detection | Daily | Unusual patterns |

### 4.2 Profiling Script

```python
from ydata_profiling import ProfileReport

def generate_data_profile(df: pd.DataFrame, output_path: str):
    """Generate comprehensive data profile"""
    
    profile = ProfileReport(
        df,
        title="BIOMETRICS Data Profile",
        minimal=False,
        explorative=True,
        correlations={
            "pearson": {"calculate": True},
            "spearman": {"calculate": True},
            "kendall": {"calculate": False},
        },
    )
    
    profile.to_file(output_path)
    return profile
```

## 5. Monitoring & Alerts

### 5.1 Alert Rules

| Alert Type | Condition | Notification |
|------------|-----------|-------------|
| Critical | Quality score < 90% | Slack + PagerDuty |
| Warning | Quality score < 95% | Slack |
| Info | Quality score < 98% | Email daily digest |

### 5.2 Dashboard

```sql
SELECT 
    dataset_name,
    quality_dimension,
    score,
    last_checked,
    CASE 
        WHEN score < 90 THEN 'critical'
        WHEN score < 95 THEN 'warning'
        WHEN score < 98 THEN 'info'
        ELSE 'pass'
    END as status
FROM data_quality_metrics
WHERE last_checked > NOW() - INTERVAL '24 hours'
ORDER BY score ASC;
```

## 6. Data Quality Metrics

### 6.1 Metrics Table

```sql
CREATE TABLE data_quality_metrics (
    id SERIAL PRIMARY KEY,
    dataset_name VARCHAR(255) NOT NULL,
    quality_dimension VARCHAR(50) NOT NULL,
    score DECIMAL(5,2) NOT NULL,
    total_records INTEGER NOT NULL,
    valid_records INTEGER NOT NULL,
    failed_records INTEGER NOT NULL,
    check_timestamp TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT unique_dataset_check UNIQUE (dataset_name, quality_dimension, DATE(check_timestamp))
);
```

### 6.2 Daily Quality Report

```python
def generate_quality_report():
    """Generate daily data quality report"""
    
    report = {
        'date': datetime.now().date(),
        'overall_score': 0,
        'datasets': [],
    }
    
    for dataset in ['biometric_readings', 'user_profiles', 'health_metrics']:
        score = calculate_quality_score(dataset)
        report['datasets'].append({
            'name': dataset,
            'score': score,
            'issues': get_issues(dataset),
        })
    
    # Calculate weighted overall score
    weights = {'biometric_readings': 0.4, 'user_profiles': 0.3, 'health_metrics': 0.3}
    report['overall_score'] = sum(
        d['score'] * weights[d['name']] 
        for d in report['datasets']
    )
    
    send_quality_report(report)
    return report
```

## 7. Data Cleaning

### 7.1 Cleaning Rules

| Issue | Strategy | Automation |
|-------|----------|------------|
| Missing values | Impute or flag | Auto |
| Duplicates | Remove | Auto |
| Outliers | Cap or investigate | Manual review |
| Invalid format | Transform | Auto |
| Stale data | Archive or delete | Auto |

### 7.2 Cleaning Pipeline

```python
class DataCleaner:
    """Data cleaning operations"""
    
    def clean_biometric_data(self, df: pd.DataFrame) -> pd.DataFrame:
        """Clean biometric data"""
        
        # Remove duplicates
        df = df.drop_duplicates(subset=['user_id', 'timestamp'])
        
        # Handle missing values
        df['heart_rate'] = df['heart_rate'].fillna(df['heart_rate'].median())
        
        # Cap outliers
        df['heart_rate'] = df['heart_rate'].clip(30, 220)
        
        # Standardize formats
        df['timestamp'] = pd.to_datetime(df['timestamp']).dt.tz_localize(None)
        
        return df
```

## 8. Governance

### 8.1 Ownership

| Dataset | Owner | Steward |
|---------|-------|---------|
| biometric_readings | Data Science | John Doe |
| user_profiles | Product | Jane Smith |
| health_metrics | Medical | Dr. Johnson |

### 8.2 SLAs

| Metric | SLA | Enforcement |
|--------|-----|------------|
| Data availability | 99.9% | Auto-monitoring |
| Quality score | 95% | Daily review |
| Issue resolution | 48 hours | Escalation |

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
