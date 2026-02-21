# REPORT-GENERATOR.md - Automated Report Generation

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - BI/Analytics  
**Author:** BIOMETRICS Engineering Team

---

## 1. Overview

This document describes the automated report generation system for BIOMETRICS, enabling scheduled and on-demand report creation in multiple formats.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| Report Engine | Python | Generate reports |
| Template Engine | Jinja2 | HTML/PDF templates |
| Scheduler | Airflow | Schedule reports |
| Storage | S3 | Report storage |
| Delivery | SendGrid | Email delivery |

### 2.2 Report Flow

```
Request → Queue → Generate → Store → Notify → Deliver
```

## 3. Report Types

### 3.1 Available Reports

| Report | Frequency | Format | Recipients |
|--------|-----------|--------|------------|
| Daily Summary | Daily 6AM | PDF | All users |
| Weekly Analytics | Monday 9AM | PDF | Managers |
| Monthly Business | 1st of month | PDF + XLSX | Leadership |
| Custom | On-demand | PDF + XLSX + CSV | Requester |
| Real-time | On-demand | PDF | Requester |

### 3.2 Report Templates

```python
class ReportTemplate:
    """Base report template"""
    
    TEMPLATES = {
        'daily_summary': {
            'name': 'Daily Summary Report',
            'description': 'Daily activity summary',
            'sections': [
                'header',
                'kpi_summary',
                'user_activity',
                'health_metrics',
                'alerts',
                'footer'
            ],
            'filters': ['date'],
        },
        'weekly_analytics': {
            'name': 'Weekly Analytics Report',
            'description': 'Weekly analytics overview',
            'sections': [
                'header',
                'executive_summary',
                'user_growth',
                'revenue_analysis',
                'feature_usage',
                'trends',
                'recommendations',
                'footer'
            ],
            'filters': ['date_range', 'segment'],
        },
        'monthly_business': {
            'name': 'Monthly Business Review',
            'description': 'Monthly business metrics',
            'sections': [
                'header',
                'financial_summary',
                'kpis',
                'user_analysis',
                'product_metrics',
                'market_analysis',
                'forecast',
                'footer'
            ],
            'filters': ['month', 'region'],
        },
    }
```

## 4. Report Engine

### 4.1 Core Engine

```python
from jinja2 import Environment, FileSystemLoader
from weasyprint import HTML
import pandas as pd
from datetime import datetime
import boto3

class ReportGenerator:
    """Generate reports in multiple formats"""
    
    def __init__(self):
        self.jinja_env = Environment(
            loader=FileSystemLoader('templates/')
        )
        self.s3_client = boto3.client('s3')
    
    def generate(
        self,
        template_name: str,
        params: dict,
        format: str = 'pdf'
    ) -> bytes:
        """Generate report"""
        
        # Fetch data
        data = self.fetch_data(template_name, params)
        
        # Render template
        if format == 'pdf':
            return self.render_pdf(template_name, data, params)
        elif format == 'html':
            return self.render_html(template_name, data, params)
        elif format == 'xlsx':
            return self.render_excel(template_name, data, params)
        elif format == 'csv':
            return self.render_csv(template_name, data, params)
    
    def fetch_data(self, template_name: str, params: dict) -> dict:
        """Fetch data for report"""
        
        queries = {
            'daily_summary': self.get_daily_summary_query,
            'weekly_analytics': self.get_weekly_analytics_query,
            'monthly_business': self.get_monthly_business_query,
        }
        
        query_func = queries.get(template_name)
        return query_func(params) if query_func else {}
```

### 4.2 PDF Generation

```python
def render_pdf(self, template_name: str, data: dict, params: dict) -> bytes:
    """Render PDF using WeasyPrint"""
    
    # Get template
    template = self.jinja_env.get_template(f'{template_name}.html')
    
    # Render HTML
    html_content = template.render(
        data=data,
        params=params,
        generated_at=datetime.now(),
        company_name='BIOMETRICS',
    )
    
    # Convert to PDF
    pdf = HTML(string=html_content).write_pdf(
        stylesheets=[
            CSS(string='@page { size: A4; margin: 2cm; }'),
            CSS(filename='templates/pdf/styles.css'),
        ]
    )
    
    return pdf
```

## 5. Scheduling

### 5.1 Airflow DAG

```python
from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime, timedelta

dag = DAG(
    'report_generation',
    schedule_interval='0 6 * * *',
    start_date=datetime(2026, 1, 1),
    catchup=False,
)

def generate_daily_reports(**context):
    """Generate all daily reports"""
    
    generator = ReportGenerator()
    ds = context['ds']
    
    reports = [
        {'template': 'daily_summary', 'format': 'pdf'},
    ]
    
    for report in reports:
        # Generate
        pdf_data = generator.generate(
            template_name=report['template'],
            params={'date': ds},
            format=report['format']
        )
        
        # Store
        s3_key = f"reports/{ds}/{report['template']}.{report['format']}"
        upload_to_s3(pdf_data, s3_key)
        
        # Send email
        send_report_email(
            recipient='team@biometrics.com',
            subject=f'Daily Report - {ds}',
            attachment_data=pdf_data,
            attachment_name=f"{report['template']}_{ds}.{report['format']}"
        )

generate_reports = PythonOperator(
    task_id='generate_daily_reports',
    python_callable=generate_daily_reports,
    dag=dag,
)
```

### 5.2 Scheduled Reports Table

```sql
CREATE TABLE scheduled_reports (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    template VARCHAR(100) NOT NULL,
    schedule_cron VARCHAR(100) NOT NULL,
    format VARCHAR(20) DEFAULT 'pdf',
    recipients JSONB NOT NULL,
    params JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    last_run TIMESTAMP,
    next_run TIMESTAMP,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## 6. Email Delivery

### 6.1 Email Service

```python
class ReportDelivery:
    """Deliver reports via email"""
    
    def send_report(
        self,
        recipients: List[str],
        subject: str,
        body: str,
        attachment_data: bytes,
        attachment_name: str
    ) -> bool:
        """Send report email"""
        
        msg = MIMEMultipart()
        msg['From'] = 'reports@biometrics.com'
        msg['To'] = ', '.join(recipients)
        msg['Subject'] = subject
        
        # Add body
        msg.attach(MIMEText(body, 'html'))
        
        # Add attachment
        part = MIMEBase('application', 'octet-stream')
        part.set_payload(attachment_data)
        encoders.encode_base64(part)
        part.add_header(
            'Content-Disposition',
            f'attachment; filename= {attachment_name}'
        )
        msg.attach(part)
        
        # Send
        self.smtp.send_message(msg)
        
        return True
```

## 7. On-Demand Reports

### 7.1 API Endpoint

```python
from fastapi import APIRouter, BackgroundTasks
from pydantic import BaseModel

router = APIRouter()

class ReportRequest(BaseModel):
    template: str
    params: dict
    format: str = 'pdf'
    email_to: Optional[List[str]] = None

@router.post('/api/reports/generate')
async def generate_report(
    request: ReportRequest,
    background_tasks: BackgroundTasks
):
    """Generate report on-demand"""
    
    # Create request
    report_request = ReportRequest(
        template=request.template,
        params=request.params,
        format=request.format,
        requested_by=get_current_user(),
        status='queued'
    )
    
    # Queue generation
    background_tasks.add_task(
        process_report_request,
        report_request.id
    )
    
    return {
        'request_id': report_request.id,
        'status': 'queued',
        'estimated_time': '5 minutes'
    }

@router.get('/api/reports/{request_id}/download')
async def download_report(request_id: str):
    """Download generated report"""
    
    report = get_report_request(request_id)
    
    if report.status != 'completed':
        raise HTTPException(400, 'Report not ready')
    
    return FileResponse(
        path=report.s3_key,
        media_type='application/pdf',
        filename=f'{report.template}.{report.format}'
    )
```

## 8. Templates

### 8.1 HTML Template

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{{ params.title }}</title>
    <style>
        body { font-family: Arial, sans-serif; }
        .header { background: #1a73e8; color: white; padding: 20px; }
        .section { margin: 20px 0; }
        .kpi-card { 
            border: 1px solid #ddd; 
            padding: 15px; 
            margin: 10px;
            display: inline-block;
        }
        table { width: 100%; border-collapse: collapse; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        .footer { margin-top: 30px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{ params.title }}</h1>
        <p>Generated: {{ generated_at }}</p>
    </div>
    
    {% for section in data.sections %}
    <div class="section">
        <h2>{{ section.title }}</h2>
        {{ section.content | safe }}
    </div>
    {% endfor %}
    
    <div class="footer">
        <p>BIOMETRICS - Confidential</p>
    </div>
</body>
</html>
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
