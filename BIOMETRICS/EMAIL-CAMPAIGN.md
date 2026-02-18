# EMAIL-CAMPAIGN.md - Email Campaign Management

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Marketing  
**Author:** BIOMETRICS Marketing Team

---

## 1. Overview

This document describes the email campaign management system for BIOMETRICS, enabling targeted email marketing with automation and analytics.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| Email Service | SendGrid | Email delivery |
| Campaign Manager | Node.js | Campaign management |
| Template Engine | Handlebars | Email templates |
| Analytics | Mixpanel | Analytics |
| List Manager | PostgreSQL | Subscriber management |

### 2.2 Setup

```yaml
# docker-compose.yml
services:
  email-service:
    build: ./email-service
    environment:
      SENDGRID_API_KEY: ${SENDGRID_API_KEY}
      FROM_EMAIL: noreply@biometrics.com
      FROM_NAME: BIOMETRICS
    ports:
      - "53053:3000"  # Port Sovereignty: 3000 → 53053
```

## 3. Campaign Types

### 3.1 Campaign Categories

| Type | Description | Trigger |
|------|-------------|---------|
| Welcome | New user onboarding | Account created |
| Re-engagement | Inactive users | 30 days inactive |
| Promotional | Offers & deals | Manual trigger |
| Transactional | Orders, confirmations | Event trigger |
| Newsletter | Regular updates | Scheduled |
| Drip | Nurture sequence | Time-based |

### 3.2 Campaign Flow

```
Campaign Created → Audience Selected → Template → A/B Test → Schedule → Send → Track
```

## 4. Audience Management

### 4.1 Segmentation

```typescript
interface Segment {
  id: string;
  name: string;
  conditions: Condition[];
  count: number;
}

interface Condition {
  field: string;
  operator: 'equals' | 'contains' | 'greater_than' | 'less_than' | 'in';
  value: any;
  logic: 'AND' | 'OR';
}

class SegmentBuilder {
  buildSegment(conditions: Condition[]): Segment {
    let query = 'SELECT * FROM users WHERE ';
    
    conditions.forEach((condition, index) => {
      if (index > 0) query += ` ${condition.logic} `;
      
      switch (condition.operator) {
        case 'equals':
          query += `${condition.field} = '${condition.value}'`;
          break;
        case 'contains':
          query += `${condition.field} LIKE '%${condition.value}%'`;
          break;
        case 'greater_than':
          query += `${condition.field} > ${condition.value}`;
          break;
      }
    });
    
    return this.executeQuery(query);
  }
}
```

### 4.2 Audience Segments

| Segment | Conditions | Size |
|---------|------------|------|
| Active Users | last_active > 7 days | 50,000 |
| Inactive 30d | last_active < 30 days | 10,000 |
| Premium Users | subscription = premium | 15,000 |
| New This Month | created_at > 30 days | 5,000 |
| High Value | LTV > $500 | 8,000 |

## 5. Template System

### 5.1 Email Templates

```html
<!-- templates/welcome.html -->
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 0 auto; }
    .header { background: #1a73e8; padding: 20px; text-align: center; }
    .content { padding: 30px; }
    .button { 
      background: #1a73e8; 
      color: white; 
      padding: 12px 24px; 
      text-decoration: none; 
      border-radius: 4px;
    }
    .footer { background: #f5f5f5; padding: 20px; text-align: center; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <img src="{{logoUrl}}" alt="BIOMETRICS" width="150">
    </div>
    <div class="content">
      <h1>Welcome, {{firstName}}!</h1>
      <p>Thank you for joining BIOMETRICS. We're excited to help you monitor your health.</p>
      <p>Get started with these features:</p>
      <ul>
        <li>Track your heart rate in real-time</li>
        <li>Monitor sleep quality</li>
        <li>Set health goals</li>
      </ul>
      <p><a href="{{dashboardUrl}}" class="button">Go to Dashboard</a></p>
    </div>
    <div class="footer">
      <p>© 2026 BIOMETRICS. All rights reserved.</p>
      <p>
        <a href="{{unsubscribeUrl}}">Unsubscribe</a> | 
        <a href="{{preferencesUrl}}">Preferences</a>
      </p>
    </div>
  </div>
</body>
</html>
```

### 5.2 Dynamic Content

```javascript
// Handlebars helpers
const template = Handlebars.compile(emailContent);

// Personalization
const data = {
  firstName: user.firstName,
  lastName: user.lastName,
  logoUrl: 'https://biometrics.com/logo.png',
  dashboardUrl: `https://biometrics.com/dashboard?user=${user.id}`,
  unsubscribeUrl: `https://biometrics.com/unsubscribe?user=${user.id}`,
  preferencesUrl: `https://biometrics.com/preferences?user=${user.id}`,
  healthScore: user.healthScore,
  recentAchievements: user.achievements.slice(0, 3),
  nextMilestone: user.nextMilestone,
};

const html = template(data);
```

## 6. A/B Testing

### 6.1 Test Types

| Test | Description | Variants |
|------|-------------|----------|
| Subject | Email subject line | 3 variants |
| Content | Email body content | 2 variants |
| Send Time | Delivery time | 4 time slots |
| CTA | Call-to-action | 2 variants |

### 6.2 Test Configuration

```typescript
interface ABTest {
  id: string;
  campaign_id: string;
  test_type: 'subject' | 'content' | 'time' | 'cta';
  variants: Variant[];
  traffic_split: number[];  // e.g., [33, 33, 34]
  duration_hours: number;
  winner_metric: 'open_rate' | 'click_rate' | 'conversion';
}

const createABTest = async (test: ABTest) => {
  // Split audience
  const audience = await getAudience(test.campaign_id);
  const splits = splitAudience(audience, test.traffic_split);
  
  // Create variants
  for (const [index, variant] of test.variants.entries()) {
    await sendVariant(variant, splits[index]);
  }
  
  // Schedule winner selection
  scheduleWinnerSelection(test);
};
```

## 7. Analytics

### 7.1 Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Open Rate | % opened | > 20% |
| Click Rate | % clicked | > 3% |
| Unsubscribe Rate | % unsubscribed | < 0.5% |
| Bounce Rate | % bounced | < 2% |
| Conversion Rate | % converted | > 1% |

### 7.2 Tracking

```javascript
// Open tracking
<img src="https://track.biometrics.com/open?campaign={{campaignId}}&user={{userId}}" />

// Click tracking
<a href="https://track.biometrics.com/click?campaign={{campaignId}}&user={{userId}}&url={{encodedUrl}}">

// Analytics query
const getCampaignAnalytics = async (campaignId: string) => {
  return await db.query(`
    SELECT 
      campaign_id,
      SUM(sent) as sent,
      SUM(delivered) as delivered,
      SUM(opened) as opened,
      SUM(clicked) as clicked,
      SUM(bounced) as bounced,
      SUM(unsubscribed) as unsubscribed,
      SUM(opened) / SUM(delivered)::float as open_rate,
      SUM(clicked) / SUM(delivered)::float as click_rate
    FROM campaign_analytics
    WHERE campaign_id = $1
    GROUP BY campaign_id
  `, [campaignId]);
};
```

## 8. Automation

### 8.1 Drip Campaigns

```yaml
name: Welcome Series
trigger: user_created
steps:
  - delay: 0
    email: welcome_1
  - delay: 2 days
    email: welcome_2
  - delay: 5 days
    email: welcome_3
  - delay: 10 days
    email: welcome_4
```

### 8.2 Trigger Rules

```javascript
const triggers = {
  'user_created': {
    condition: 'user.created_at >= NOW() - INTERVAL 1 hour',
    delay: 0,
    template: 'welcome',
  },
  'inactive_30d': {
    condition: 'user.last_active < NOW() - INTERVAL 30 DAY',
    delay: 0,
    template: 're-engagement',
  },
  'purchase_completed': {
    condition: 'order.status = completed',
    delay: 0,
    template: 'order_confirmation',
  },
  'subscription_expiring': {
    condition: 'subscription.end_date < NOW() + 7 DAYS',
    delay: 0,
    template: 'subscription_expiring',
  },
};
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
