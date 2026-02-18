# MARKETING-AUTOMATION.md - Marketing Automation Platform

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Marketing  
**Author:** BIOMETRICS Marketing Team

---

## 1. Overview

This document describes the marketing automation platform for BIOMETRICS, enabling sophisticated multi-channel campaigns with lead scoring and behavioral targeting.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| Automation Engine | n8n | Workflow automation |
| CRM | HubSpot | Customer management |
| CDP | Segment | Customer data |
| Analytics | Mixpanel | Event tracking |
| Email | SendGrid | Email delivery |

### 2.2 Flow

```
User Action → Event → Segment → Automation → Action → Track → Optimize
```

## 3. Automation Workflows

### 3.1 Workflow Types

| Workflow | Trigger | Channels |
|----------|---------|----------|
| Lead Nurture | Form submit | Email + SMS |
| Re-engagement | 30 days inactive | Email |
| Welcome Series | New signup | Email + In-app |
| Abandonment Recovery | Cart abandoned | Email |
| Upsell Campaign | Usage threshold | Email + Push |
| Win-back | Churned user | Email + SMS |

### 3.2 n8n Workflows

```yaml
# n8n-workflow.yml
name: Lead Nurture Workflow
nodes:
  - name: Trigger
    type: webhook
    webhook: lead_nurture_trigger

  - name: Score Lead
    type: code
    code: |
      const score = calculateLeadScore($input.item.json);
      return [{ json: { ...$input.item.json, lead_score: score } }];

  - name: Route by Score
    type: switch
    conditions:
      - value1: "{{ $json.lead_score }}"
        operation: greater_than
        value2: 50
    output: high_score

  - name: High Score - Sales
    type: slack
    channel: #sales-leads
    text: "New high-score lead: {{ $json.email }}"

  - name: Medium Score - Nurture
    type: email
    to: "{{ $json.email }}"
    subject: "Thanks for your interest in BIOMETRICS"
    template: lead_nurture

  - name: Low Score - Educational
    type: email
    to: "{{ $json.email }}"
    subject: "Here's more about BIOMETRICS"
    template: educational
```

## 4. Lead Scoring

### 4.1 Scoring Model

```typescript
interface LeadScore {
  total: number;
  demographics: number;
  engagement: number;
  behavioral: number;
  firmographic: number;
}

const SCORING_RULES = {
  demographics: {
    'age_25_34': 10,
    'age_35_44': 15,
    'age_45_54': 20,
    'income_50k_75k': 10,
    'income_75k_100k': 15,
    'income_100k_plus': 20,
  },
  engagement: {
    'email_opened': 2,
    'email_clicked': 5,
    'page_viewed_pricing': 10,
    'page_viewed_demo': 15,
    'form_submitted': 20,
    'chat_started': 15,
  },
  behavioral: {
    'visited_3_plus_pages': 10,
    'return_visit': 5,
    'spent_5_minutes': 10,
    'returned_within_7_days': 15,
  },
};

const calculateLeadScore = (lead: Lead): LeadScore => {
  let demographics = 0;
  let engagement = 0;
  let behavioral = 0;

  // Demographics
  if (lead.age) demographics += SCORING_RULES.demographics[`age_${lead.age}`] || 0;
  if (lead.income) demographics += SCORING_RULES.demographics[`income_${lead.income}`] || 0;

  // Engagement
  if (lead.emailOpened) engagement += 2;
  if (lead.emailClicked) engagement += 5;
  if (lead.viewedPricing) engagement += 10;

  // Behavioral
  if (lead.pagesVisited >= 3) behavioral += 10;
  if (lead.returnVisit) behavioral += 5;

  return {
    demographics,
    engagement,
    behavioral,
    firmographic: 0,
    total: demographics + engagement + behavioral,
  };
};
```

### 4.2 Score Thresholds

| Score | Stage | Action |
|-------|-------|--------|
| 0-25 | Cold | Educational content |
| 26-50 | Warm | Nurture sequence |
| 51-75 | Hot | Sales outreach |
| 76-100 | Ready | Demo booking |

## 5. Behavioral Targeting

### 5.1 Event Tracking

```typescript
interface TrackingEvent {
  event: string;
  userId: string;
  timestamp: Date;
  properties: Record<string, any>;
  context: {
    page: string;
    source: string;
    device: string;
  };
}

// Track user actions
analytics.track('viewed_pricing', {
  userId: 'user123',
  plan: 'premium',
  features: ['biometric', 'analytics'],
});

analytics.track('clicked_cta', {
  userId: 'user123',
  cta: 'start_trial',
  location: 'pricing_page',
});
```

### 5.2 Audience Triggers

| Trigger | Conditions | Action |
|---------|------------|--------|
| High Intent | Viewed pricing 3+ times | Sales notification |
| Price Sensitive | Viewed pricing, no conversion | Discount offer |
| Power User | Daily active > 30 min, 7+ days | Upsell premium |
| At Risk | No login in 30 days | Re-engagement |
| New Lead | First 7 days | Onboarding sequence |

## 6. Multi-Channel Orchestration

### 6.1 Channel Matrix

| Channel | Use Case | Timing | Personalization |
|---------|----------|--------|-----------------|
| Email | Nurture, announcements | Day 0, 3, 7 | High |
| SMS | Urgent, OTPs | Real-time | Medium |
| Push | Re-engagement | Real-time | Medium |
| In-App | Guidance | Contextual | High |
| Chat | Support | On-request | High |

### 6.2 Cross-Channel Flow

```typescript
const orchestrateCampaign = async (userId: string, trigger: string) => {
  const user = await getUser(userId);
  const channel = determineBestChannel(user, trigger);

  switch (channel) {
    case 'email':
      await sendEmail(user, trigger);
      break;
    case 'sms':
      await sendSMS(user, trigger);
      break;
    case 'push':
      await sendPush(user, trigger);
      break;
    case 'in_app':
      await showInAppMessage(user, trigger);
      break;
  }

  // Log for optimization
  await logCampaignAction(userId, trigger, channel);
};
```

## 7. Analytics & Attribution

### 7.1 Attribution Model

```sql
-- Multi-touch attribution
WITH touchpoints AS (
  SELECT 
    user_id,
    campaign_id,
    channel,
    event_date,
    ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY event_date) as touch_num
  FROM campaign_touchpoints
  WHERE event_date >= NOW() - INTERVAL '30 days'
)
SELECT 
  channel,
  COUNT(DISTINCT user_id) as users,
  SUM(conversion_value) as revenue,
  ROUND(SUM(conversion_value) / NULLIF(COUNT(DISTINCT user_id), 0), 2) as aov
FROM touchpoints t
LEFT JOIN conversions c ON t.user_id = c.user_id
GROUP BY channel;
```

### 7.2 Campaign Analytics

| Metric | Description | Target |
|--------|-------------|--------|
| MQLs | Marketing qualified leads | 500/month |
| SQLs | Sales qualified leads | 100/month |
| CAC | Customer acquisition cost | < $50 |
| LTV | Lifetime value | > $500 |
| ROI | Return on investment | > 300% |

## 8. Integration

### 8.1 CRM Integration

```typescript
const syncToCRM = async (lead: Lead) => {
  const hubspot = new HubSpotClient();

  // Create or update contact
  await hubspot.contacts.createOrUpdate({
    email: lead.email,
    properties: {
      firstname: lead.firstName,
      lastname: lead.lastName,
      company: lead.company,
      lead_score: lead.score,
      lifecycle_stage: getLifecycleStage(lead.score),
    },
  });

  // Add to list
  await hubspot.lists.addContact(
    getListForScore(lead.score),
    lead.email
  );
};
```

### 8.2 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/campaigns | GET | List campaigns |
| /api/campaigns | POST | Create campaign |
| /api/automations | GET | List automations |
| /api/analytics | GET | Get analytics |
| /api/audiences | GET | List audiences |

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
