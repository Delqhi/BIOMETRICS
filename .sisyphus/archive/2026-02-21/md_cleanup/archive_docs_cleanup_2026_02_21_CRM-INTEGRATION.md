# CRM-INTEGRATION.md - CRM System Integration

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Marketing  
**Author:** BIOMETRICS Integration Team

---

## 1. Overview

This document describes the CRM integration architecture for BIOMETRICS, enabling seamless synchronization with customer relationship management systems.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| CRM | HubSpot | Customer management |
| Integration Layer | Node.js | Sync logic |
| Event Bus | Kafka | Real-time events |
| Cache | Redis | Performance |

### 2.2 Integration Flow

```
BIOMETRICS App → Event Bus → Integration Service → CRM API
                                         ↓
                                   Webhook → BIOMETRICS
```

## 3. Sync Strategy

### 3.1 Data Types

| Entity | Direction | Frequency | Priority |
|--------|-----------|-----------|----------|
| Contacts | Bidirectional | Real-time | High |
| Companies | BIOMETRICS → CRM | Hourly | Medium |
| Deals | Bidirectional | Real-time | High |
| Tickets | CRM → BIOMETRICS | Real-time | High |
| Notes | Bidirectional | Real-time | Medium |
| Activities | BIOMETRICS → CRM | Daily | Low |

### 3.2 Sync Modes

| Mode | Use Case | Latency |
|------|----------|---------|
| Real-time | Contact updates | < 1s |
| Batch | Bulk operations | Hourly |
| Webhook | External triggers | Immediate |

## 4. Contact Sync

### 4.1 Contact Mapping

```typescript
interface ContactMapping {
  biometricsField: string;
  crmField: string;
  transform?: (value: any) => any;
}

const CONTACT_MAPPING: ContactMapping[] = [
  { biometricsField: 'email', crmField: 'email', transform: (v) => v.toLowerCase() },
  { biometricsField: 'firstName', crmField: 'firstname' },
  { biometricsField: 'lastName', crmField: 'lastname' },
  { biometricsField: 'phone', crmField: 'phone' },
  { biometricsField: 'company', crmField: 'company' },
  { biometricsField: 'jobTitle', crmField: 'jobtitle' },
  { biometricsField: 'createdAt', crmField: 'createdate' },
  { biometricsField: 'subscription.plan', crmField: 'subscription_plan' },
  { biometricsField: 'subscription.status', crmField: 'subscription_status' },
  { biometricsField: 'lifetimeValue', crmField: 'lifetime_value' },
  { biometricsField: 'healthScore', crmField: 'health_score' },
];
```

### 4.2 Sync Implementation

```typescript
class ContactSync {
  async syncContact(biometricsContact: Contact): Promise<SyncResult> {
    const crmContact = this.mapToCRM(biometricsContact);
    
    // Check if contact exists
    const existing = await this.crm.getContactByEmail(crmContact.email);
    
    if (existing) {
      // Update
      const result = await this.crm.updateContact(existing.id, crmContact);
      await this.logSync(existing.id, 'update', result);
      return { success: true, action: 'update', crmId: existing.id };
    } else {
      // Create
      const result = await this.crm.createContact(crmContact);
      await this.logSync(result.id, 'create', result);
      return { success: true, action: 'create', crmId: result.id };
    }
  }

  private mapToCRM(contact: Contact): CRMContact {
    const crmContact: any = {};
    
    for (const mapping of CONTACT_MAPPING) {
      const value = get(contact, mapping.biometricsField);
      crmContact[mapping.crmField] = mapping.transform 
        ? mapping.transform(value) 
        : value;
    }
    
    return crmContact;
  }
}
```

## 5. Deal Management

### 5.1 Deal Pipeline

```typescript
interface Deal {
  id: string;
  name: string;
  stage: DealStage;
  amount: number;
  probability: number;
  closeDate: Date;
  contactId: string;
  ownerId: string;
}

enum DealStage {
  APPOINTMENT = 'appointment',
  QUALIFICATION = 'qualification',
  PROPOSAL = 'proposal',
  NEGOTIATION = 'negotiation',
  CLOSED_WON = 'closedwon',
  CLOSED_LOST = 'closedlost',
}

const createDealFromSubscription = async (subscription: Subscription): Promise<Deal> => {
  return {
    name: `${subscription.user.name} - ${subscription.plan.name}`,
    stage: DealStage.QUALIFICATION,
    amount: subscription.plan.price,
    probability: getProbability(subscription.plan.name),
    closeDate: subscription.startDate,
    contactId: subscription.user.crmContactId,
    ownerId: getSalesRep(subscription.user.region),
  };
};
```

### 5.2 Stage Automation

```typescript
const updateDealStage = async (dealId: string, newStage: DealStage) => {
  const deal = await crm.getDeal(dealId);
  
  // Update stage
  await crm.updateDeal(dealId, { stage: newStage });
  
  // Trigger actions based on stage
  switch (newStage) {
    case DealStage.CLOSED_WON:
      await triggerOnboarding(deal);
      await createSubscription(deal);
      break;
    case DealStage.CLOSED_LOST:
      await sendFeedbackSurvey(deal);
      break;
    case DealStage.NEGOTIATION:
      await notifySalesManager(deal);
      break;
  }
};
```

## 6. Webhook Handling

### 6.1 CRM Webhooks

```typescript
// Handle CRM webhooks
const handleCRMWebhook = async (event: WebhookEvent) => {
  const { type, objectId, properties } = event;
  
  switch (type) {
    case 'contact.create':
    case 'contact.update':
      await handleContactUpdate(objectId, properties);
      break;
    case 'deal.create':
    case 'deal.update':
      await handleDealUpdate(objectId, properties);
      break;
    case 'ticket.create':
      await handleTicketCreate(objectId, properties);
      break;
    case 'note.create':
      await handleNoteCreate(objectId, properties);
      break;
  }
};

const handleContactUpdate = async (contactId: string, properties: any) => {
  const crmContact = await crm.getContact(contactId);
  const biometricsUser = await db.getUserByEmail(crmContact.email);
  
  if (biometricsUser) {
    await db.updateUser(biometricsUser.id, {
      firstName: crmContact.firstname,
      lastName: crmContact.lastname,
      phone: crmContact.phone,
      company: crmContact.company,
    });
  }
};
```

### 6.2 Webhook Events

| Event | Source | Action |
|-------|--------|--------|
| Contact created | CRM | Create user |
| Contact updated | CRM | Update user |
| Deal stage changed | CRM | Update subscription |
| Deal closed won | CRM | Activate subscription |
| Deal closed lost | CRM | Cancel subscription |

## 7. Activity Tracking

### 7.1 Activity Logging

```typescript
const logActivityToCRM = async (userId: string, activity: UserActivity) => {
  const crmContactId = await getCRMContactId(userId);
  
  if (!crmContactId) return;
  
  const activityNote = {
    body: formatActivityAsNote(activity),
    timestamp: activity.timestamp,
  };
  
  await crm.createNote(crmContactId, activityNote);
};

const formatActivityAsNote = (activity: UserActivity): string => {
  return `
Activity: ${activity.type}
Date: ${activity.timestamp}
Details: ${JSON.stringify(activity.details)}
  `.trim();
};
```

### 7.2 Activity Types

| Activity | CRM Note | Priority |
|----------|----------|----------|
| Subscription upgrade | Yes | High |
| Feature usage | Yes | Medium |
| Support ticket | Yes | High |
| Login | No | Low |
| Page view | No | Low |

## 8. Performance

### 8.1 Caching Strategy

```typescript
class CRMCache {
  private cache: Redis;
  
  async getContact(crmId: string): Promise<Contact | null> {
    const cached = await this.cache.get(`crm:contact:${crmId}`);
    if (cached) return JSON.parse(cached);
    
    const contact = await crm.getContact(crmId);
    await this.cache.setex(`crm:contact:${crmId}`, 3600, JSON.stringify(contact));
    return contact;
  }
  
  async invalidateContact(crmId: string): Promise<void> {
    await this.cache.del(`crm:contact:${crmId}`);
  }
}
```

### 8.2 Rate Limiting

| Operation | Limit | Window |
|-----------|-------|--------|
| API calls | 100 | 10 seconds |
| Bulk sync | 10,000 | 1 hour |
| Webhooks | 50 | 10 seconds |

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
