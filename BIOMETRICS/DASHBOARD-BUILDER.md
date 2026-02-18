# DASHBOARD-BUILDER.md - Custom Dashboard Builder

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - BI/Analytics  
**Author:** BIOMETRICS Frontend Team

---

## 1. Overview

This document describes the custom dashboard builder for BIOMETRICS, enabling users to create personalized dashboards with drag-and-drop functionality.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| Dashboard Editor | React + DnD | UI builder |
| Widget Library | React | Chart components |
| State Management | Redux | Dashboard state |
| API | Node.js | Save/load |

### 2.2 User Flow

```
User → Create Dashboard → Add Widgets → Configure → Save → Share
```

## 3. Widget Types

### 3.1 Available Widgets

| Widget | Description | Data Source |
|--------|-------------|-------------|
| Metric Card | Single KPI display | Query |
| Line Chart | Time series | Query |
| Bar Chart | Categorical data | Query |
| Pie Chart | Distribution | Query |
| Table | Data grid | Query |
| Funnel | Conversion | Query |
| Heatmap | Correlation | Query |
| Map | Geographic | Query |
| Text | Rich text | Static |
| Image | Image display | URL |

### 3.2 Widget Configuration

```typescript
interface WidgetConfig {
  id: string;
  type: WidgetType;
  title: string;
  position: { x: number; y: number; w: number; h: number };
  dataSource: {
    type: 'query' | 'static' | 'api';
    query?: string;
    apiUrl?: string;
    staticData?: any;
  };
  visualization: {
    colors?: string[];
    showLegend?: boolean;
    showTooltip?: boolean;
    xAxis?: string;
    yAxis?: string;
    aggregation?: 'sum' | 'avg' | 'count' | 'min' | 'max';
  };
  refreshInterval?: number; // seconds
  filters?: Filter[];
}
```

## 4. Drag and Drop

### 4.1 Grid Layout

```tsx
import { Responsive, WidthProvider } from 'react-grid-layout';
import { useDrag, useDrop } from 'react-dnd';

const ResponsiveGridLayout = WidthProvider(Responsive);

const DashboardEditor = ({ widgets, onLayoutChange }) => {
  return (
    <ResponsiveGridLayout
      className="layout"
      layouts={{ lg: widgets.map(w => ({
        i: w.id,
        x: w.position.x,
        y: w.position.y,
        w: w.position.w,
        h: w.position.h,
      }))}}
      breakpoints={{ lg: 1200, md: 996, sm: 768 }}
      cols={{ lg: 12, md: 10, sm: 6 }}
      rowHeight={30}
      onLayoutChange={(layout) => onLayoutChange(layout)}
      draggableHandle=".drag-handle"
    >
      {widgets.map(widget => (
        <WidgetCard key={widget.id} widget={widget}>
          <WidgetRenderer widget={widget} />
        </WidgetCard>
      ))}
    </ResponsiveGridLayout>
  );
};
```

### 4.2 Widget Toolbar

```tsx
const WidgetToolbar = () => {
  const widgetTypes = [
    { type: 'metric', icon: '📊', label: 'Metric' },
    { type: 'line', icon: '📈', label: 'Line Chart' },
    { type: 'bar', icon: '📉', label: 'Bar Chart' },
    { type: 'table', icon: '📋', label: 'Table' },
    { type: 'text', icon: '📝', label: 'Text' },
    { type: 'image', icon: '🖼️', label: 'Image' },
  ];

  return (
    <div className="widget-toolbar">
      {widgetTypes.map(wt => (
        <DraggableWidget
          key={wt.type}
          type={wt.type}
          icon={wt.icon}
          label={wt.label}
        />
      ))}
    </div>
  );
};
```

## 5. Query Builder

### 5.1 Visual Query Builder

```tsx
const QueryBuilder = ({ onQueryChange }) => {
  const [tables, setTables] = useState(['users', 'biometric_events', 'subscriptions']);
  const [selectedTable, setSelectedTable] = useState(null);
  const [columns, setColumns] = useState([]);
  const [filters, setFilters] = useState([]);
  const [aggregations, setAggregations] = useState([]);

  const generateSQL = () => {
    let sql = 'SELECT ';
    
    if (aggregations.length > 0) {
      sql += aggregations.map(a => 
        `${a.aggregation}(${a.column}) as ${a.alias}`
      ).join(', ');
    } else {
      sql += columns.join(', ');
    }
    
    sql += ` FROM ${selectedTable}`;
    
    if (filters.length > 0) {
      sql += ' WHERE ' + filters.map(f => 
        `${f.column} ${f.operator} '${f.value}'`
      ).join(' AND ');
    }
    
    return sql;
  };

  return (
    <div className="query-builder">
      <TableSelector 
        tables={tables}
        onSelect={setSelectedTable}
      />
      <ColumnSelector 
        columns={columns}
        onChange={setColumns}
      />
      <FilterBuilder 
        filters={filters}
        onChange={setFilters}
      />
      <AggregationBuilder 
        aggregations={aggregations}
        onChange={setAggregations}
      />
      <SQLPreview sql={generateSQL()} />
    </div>
  );
};
```

### 5.2 Chart Configuration

```typescript
interface ChartConfig {
  type: 'line' | 'bar' | 'pie' | 'area' | 'scatter';
  xAxis: {
    field: string;
    label?: string;
    type?: 'category' | 'time' | 'value';
  };
  yAxis: {
    field: string;
    label?: string;
    aggregation?: 'sum' | 'avg' | 'count';
  };
  colors?: string[];
  showLegend?: boolean;
  showGrid?: boolean;
  tooltip?: {
    show: boolean;
    fields: string[];
  };
}
```

## 6. Dashboard State

### 6.1 Save/Load

```typescript
interface Dashboard {
  id: string;
  name: string;
  description?: string;
  layout: WidgetLayout[];
  widgets: Widget[];
  filters: DashboardFilter[];
  variables: Variable[];
  createdAt: Date;
  updatedAt: Date;
  createdBy: string;
  isPublic: boolean;
}

const saveDashboard = async (dashboard: Dashboard) => {
  const response = await api.post('/dashboards', dashboard);
  return response.data;
};

const loadDashboard = async (id: string) => {
  const response = await api.get(`/dashboards/${id}`);
  return response.data;
};
```

### 6.2 Version History

```sql
CREATE TABLE dashboard_versions (
    id UUID PRIMARY KEY,
    dashboard_id UUID REFERENCES dashboards(id),
    version INT NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    change_description TEXT,
    
    UNIQUE(dashboard_id, version)
);
```

## 7. Real-time Updates

### 7.1 WebSocket Connection

```typescript
class DashboardSocket {
  private ws: WebSocket;
  private subscriptions: Set<string> = new Set();

  connect(dashboardId: string) {
    this.ws = new WebSocket(`wss://api.biometrics.com/ws/dashboards/${dashboardId}`);
    
    this.ws.onmessage = (event) => {
      const update = JSON.parse(event.data);
      this.handleUpdate(update);
    };
  }

  subscribe(widgetId: string) {
    this.subscriptions.add(widgetId);
    this.ws.send(JSON.stringify({
      action: 'subscribe',
      widgetId,
    }));
  }

  private handleUpdate(update: WidgetUpdate) {
    store.dispatch(updateWidgetData(update));
  }
}
```

## 8. Sharing & Embedding

### 8.1 Sharing Options

| Option | Description | Access Control |
|--------|-------------|---------------|
| Public | Anyone with link | Read-only |
| Internal | Logged in users | By role |
| Specific | Named users | Direct grant |
| Time-limited | Expiring access | Expiry date |

### 8.2 Embed Code

```html
<iframe 
  src="https://biometrics.com/embed/dashboard/123"
  width="100%" 
  height="800"
  frameborder="0"
></iframe>

<script>
  window.biometricsEmbed.init({
    container: '#dashboard-container',
    dashboardId: '123',
    theme: 'light',
    filters: {
      user_id: 'current'
    }
  });
</script>
```

## 9. Templates

### 9.1 Template Library

| Template | Description | Widgets |
|----------|-------------|---------|
| Executive Overview | KPIs + trends | 6 |
| Health Monitor | Biometric real-time | 8 |
| Sales Dashboard | Revenue + pipeline | 5 |
| Operations | System metrics | 4 |

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
