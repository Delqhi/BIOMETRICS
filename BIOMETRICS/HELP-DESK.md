# HELP-DESK.md - Help Desk System

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Support  
**Author:** BIOMETRICS Support Team

---

## 1. Overview

This document describes the help desk system for BIOMETRICS, enabling efficient customer support ticket management and resolution.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| Ticketing | Zammad | Ticket management |
| Knowledge Base | Custom | Self-service |
| Email | IMAP/SMTP | Email support |
| Chat | LiveChat | Real-time support |
| Phone | Twilio | Voice support |

### 2.2 Setup

```yaml
# docker-compose.yml
services:
  zammad:
    image: zammad/zammad:latest
    environment:
      ZAMMAD_AUTO_EOF: "true"
      ZAMMAD_RAILSFATN: web
      ZAMMAD_RUN_JOBS: "true"
    ports:
      - "51303:80"
    volumes:
      - ./zammad-data:/opt/zammad
```

## 3. Ticket Workflow

### 3.1 Ticket Lifecycle

```
New → Open → Pending → Waiting → Solved → Closed
  ↓
Escalated → Priority Handling → Resolution → Follow-up
```

### 3.2 Ticket Creation

```typescript
interface Ticket {
  id: string;
  subject: string;
  description: string;
  category: TicketCategory;
  priority: Priority;
  status: TicketStatus;
  requesterId: string;
  assigneeId?: string;
  groupId: string;
  tags: string[];
  createdAt: Date;
  updatedAt: Date;
  resolvedAt?: Date;
}

enum TicketCategory {
  TECHNICAL = 'technical',
  BILLING = 'billing',
  ACCOUNT = 'account',
  FEATURE_REQUEST = 'feature_request',
  BUG_REPORT = 'bug_report',
  GENERAL = 'general',
}

enum Priority {
  URGENT = 1,
  HIGH = 2,
  NORMAL = 3,
  LOW = 4,
}

const createTicket = async (data: CreateTicketRequest): Promise<Ticket> => {
  const ticket: Ticket = {
    id: generateUUID(),
    subject: data.subject,
    description: data.description,
    category: data.category,
    priority: determinePriority(data),
    status: TicketStatus.NEW,
    requesterId: data.requesterId,
    groupId: getGroupForCategory(data.category),
    tags: extractTags(data),
    createdAt: new Date(),
    updatedAt: new Date(),
  };
  
  // Save to database
  await db.tickets.insert(ticket);
  
  // Notify assignee
  await notifyAssignee(ticket);
  
  // Create internal note
  await addInternalNote(ticket.id, `Ticket created via ${data.channel}`);
  
  return ticket;
};
```

## 4. Routing Rules

### 4.1 Auto-Assignment

```typescript
const autoAssignTicket = async (ticket: Ticket): Promise<void> => {
  // Get available agents
  const agents = await getAvailableAgents(ticket.groupId);
  
  if (agents.length === 0) {
    // Put in queue
    await addToQueue(ticket);
    return;
  }
  
  // Round-robin or skills-based
  const assignee = await selectBestAgent(agents, ticket);
  
  await assignTicket(ticket.id, assignee.id);
};

const selectBestAgent = async (agents: Agent[], ticket: Ticket): Promise<Agent> => {
  // Score each agent
  const scores = await Promise.all(
    agents.map(async (agent) => {
      let score = 0;
      
      // Skills match
      if (agent.skills.includes(ticket.category)) score += 50;
      
      // Workload
      const workload = await getAgentWorkload(agent.id);
      score -= workload * 5;
      
      // Language preference
      if (agent.languages.includes(ticket.language)) score += 20;
      
      return { agent, score };
    })
  );
  
  // Return highest score
  return scores.sort((a, b) => b.score - a.score)[0].agent;
};
```

### 4.2 Routing Rules

| Condition | Action |
|-----------|--------|
| Category = Billing | Route to Billing group |
| Priority = Urgent | Route to Senior agents |
| Language = DE | Route to German team |
| User VIP | Route to VIP team |
| Technical complex | Route to Technical team |

## 5. SLA Management

### 5.1 SLA Policies

```typescript
interface SLAPolicy {
  id: string;
  name: string;
  conditions: Condition[];
  firstResponseTime: number;  // minutes
  nextResponseTime: number;
  resolutionTime: number;
  businessHours: BusinessHours;
}

const SLA_POLICIES: SLAPolicy[] = [
  {
    id: 'premium_support',
    conditions: [{ field: 'user.plan', equals: 'premium' }],
    firstResponseTime: 15,
    nextResponseTime: 30,
    resolutionTime: 240,
    businessHours: '24x7',
  },
  {
    id: 'standard_support',
    conditions: [{ field: 'user.plan', equals: 'standard' }],
    firstResponseTime: 60,
    nextResponseTime: 120,
    resolutionTime: 1440,
    businessHours: '9x5',
  },
  {
    id: 'bug_report',
    conditions: [{ field: 'category', equals: 'bug_report' }],
    firstResponseTime: 240,
    nextResponseTime: 480,
    resolutionTime: 2880,
    businessHours: '9x5',
  },
];
```

### 5.2 SLA Tracking

```typescript
const checkSLA = async (ticketId: string): Promise<SLAStatus> => {
  const ticket = await getTicket(ticketId);
  const policy = await getPolicyForTicket(ticket);
  const now = new Date();
  
  const firstResponseDeadline = addMinutes(ticket.createdAt, policy.firstResponseTime);
  const resolutionDeadline = addMinutes(ticket.createdAt, policy.resolutionTime);
  
  const status: SLAStatus = {
    firstResponse: {
      deadline: firstResponseDeadline,
      remaining: now < firstResponseDeadline ? firstResponseDeadline.getTime() - now.getTime() : 0,
      breached: now > firstResponseDeadline && !ticket.firstResponseAt,
    },
    resolution: {
      deadline: resolutionDeadline,
      remaining: now < resolutionDeadline ? resolutionDeadline.getTime() - now.getTime() : 0,
      breached: now > resolutionDeadline && !ticket.resolvedAt,
    },
  };
  
  if (status.firstResponse.breached || status.resolution.breached) {
    await notifySLABreach(ticket, status);
  }
  
  return status;
};
```

## 6. Knowledge Base

### 6.1 Article Structure

```typescript
interface KBArticle {
  id: string;
  title: string;
  content: string;
  category: string;
  tags: string[];
  authorId: string;
  status: 'draft' | 'published' | 'archived';
  views: number;
  helpful: number;
  notHelpful: number;
  createdAt: Date;
  updatedAt: Date;
}

const createArticle = async (data: CreateArticleRequest): Promise<KBArticle> => {
  const article: KBArticle = {
    id: generateUUID(),
    title: data.title,
    content: await renderMarkdown(data.content),
    category: data.category,
    tags: data.tags,
    authorId: data.authorId,
    status: 'draft',
    views: 0,
    helpful: 0,
    notHelpful: 0,
    createdAt: new Date(),
    updatedAt: new Date(),
  };
  
  await db.kb_articles.insert(article);
  return article;
};
```

### 6.2 KB Categories

| Category | Articles | Description |
|----------|----------|-------------|
| Getting Started | 15 | Onboarding |
| Account | 20 | Account management |
| Billing | 18 | Subscription & payments |
| Technical | 45 | Troubleshooting |
| API | 25 | Developer docs |
| Security | 12 | Security settings |

## 7. Reporting

### 7.1 Support Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| First Response Time | Time to first reply | < 60 min |
| Resolution Time | Time to solve | < 24 hours |
| CSAT | Customer satisfaction | > 90% |
| Ticket Volume | Daily tickets | Monitor |
| Backlog | Open tickets | < 100 |

### 7.2 Reports

```sql
-- Weekly support report
SELECT 
    DATE_TRUNC('week', created_at) as week,
    COUNT(*) as total_tickets,
    AVG(EXTRACT(EPOCH FROM (first_response_at - created_at))/60) as avg_first_response_min,
    AVG(EXTRACT(EPOCH FROM (resolved_at - created_at))/3600) as avg_resolution_hours,
    COUNT(*) FILTER (WHERE status = 'closed') as resolved,
    COUNT(*) FILTER (WHERE priority = 1) as urgent
FROM tickets
WHERE created_at >= NOW() - INTERVAL '4 weeks'
GROUP BY DATE_TRUNC('week', created_at)
ORDER BY week;
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
