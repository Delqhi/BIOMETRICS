# Support Tickets

## Overview
Customer support ticket system for BIOMETRICS.

## Features
- Ticket creation and management
- Priority levels (low, medium, high, urgent)
- Status workflow (open, in_progress, resolved, closed)
- Assignment to support agents
- Internal notes
- File attachments

## Ticket Schema
```typescript
// types/support-ticket.ts
interface SupportTicket {
  id: string;
  userId: string;
  title: string;
  description: string;
  category: 'billing' | 'technical' | 'account' | 'other';
  priority: 'low' | 'medium' | 'high' | 'urgent';
  status: 'open' | 'in_progress' | 'resolved' | 'closed';
  assignedTo?: string;
  attachments: string[];
  createdAt: Date;
  updatedAt: Date;
}

interface TicketReply {
  id: string;
  ticketId: string;
  userId: string;
  isInternal: boolean;
  message: string;
  attachments: string[];
}
```

## API Endpoints

### Tickets
```
POST /api/support/tickets
GET /api/support/tickets
GET /api/support/tickets/:id
PUT /api/support/tickets/:id
DELETE /api/support/tickets/:id
```

### Replies
```
POST /api/support/tickets/:id/reply
GET /api/support/tickets/:id/replies
```

### Admin
```
GET /api/support/tickets?status=open
PUT /api/support/tickets/:id/assign
PUT /api/support/tickets/:id/priority
```

## Email Integration
- Incoming: Forward support@biometrics.app to tickets
- Outgoing: Email notification on new reply

## SLA
| Priority | Response Time | Resolution Time |
|----------|---------------|-----------------|
| Urgent | 1 hour | 4 hours |
| High | 4 hours | 24 hours |
| Medium | 24 hours | 72 hours |
| Low | 72 hours | 7 days |

## Storage
- Tickets: Supabase `support_tickets` table
- Attachments: Supabase Storage `/support/attachments/`
