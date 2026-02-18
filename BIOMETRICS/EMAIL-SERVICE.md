# EMAIL-SERVICE.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- E-Mail-Versand erfolgt DSGVO-konform mit Consent-Management.
- Alle E-Mails werden protokolliert und sind auditierbar.
- Bounce-Handling und Unsubscribe sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

## Zweck
Dokumentation der E-Mail-Service-Architektur und -Konfiguration für BIOMETRICS.

## 1) E-Mail-Anbieter

| Anbieter | Status | Use Case | Kosten |
|----------|--------|----------|--------|
| Resend | ✅ Primary | Transaktional | €0-29/Monat |
| Gmail API | ✅ Fallback | Backup | €0 |
| SendGrid | ❌ | Nicht genutzt | - |

## 2) E-Mail-Typen

| Template | Trigger | Kategorie | Priority |
|----------|---------|-----------|----------|
| welcome | User Registration | Transaktional | HIGH |
| verify-email | Email Verification | Transaktional | HIGH |
| password-reset | Password Reset | Transaktional | HIGH |
| order-confirmation | Order Created | Transaktional | HIGH |
| order-shipped | Shipping Update | Transaktional | MEDIUM |
| invoice | Payment Received | Transaktional | HIGH |
| newsletter | Scheduled | Marketing | LOW |
| welcome-series | User Milestone | Marketing | MEDIUM |

## 3) Implementierung

### 3.1) Resend Client

```typescript
import { Resend } from 'resend'

const resend = new Resend(process.env.RESEND_API_KEY)

interface EmailOptions {
  to: string | string[]
  subject: string
  template: string
  data: object
  replyTo?: string
}

async function sendEmail(options: EmailOptions) {
  const { to, subject, template, data, replyTo } = options
  
  // Get template content
  const templateHtml = await getTemplate(template, data)
  
  const email = await resend.emails.send({
    from: 'BIOMETRICS <noreply@biometrics.de>',
    to: Array.isArray(to) ? to : [to],
    subject,
    html: templateHtml,
    reply_to: replyTo,
    headers: {
      'X-Entity-Ref-ID': crypto.randomUUID(),
      'X-Priority': '1'
    }
  })
  
  // Log for audit
  await logEmail({
    messageId: email.data?.id,
    to,
    subject,
    template,
    status: email.error ? 'failed' : 'sent'
  })
  
  return email
}
```

### 3.2) Template-System

```typescript
// templates/welcome.html
const getTemplate = (name: string, data: object): string => {
  const templates = {
    'welcome': `
      <h1>Welcome ${data.name}!</h1>
      <p>Thanks for joining BIOMETRICS.</p>
      <a href="${data.verifyUrl}">Verify your email</a>
    `,
    'order-confirmation': `
      <h1>Order Confirmed!</h1>
      <p>Order #${data.orderId}</p>
      <p>Total: €${data.total}</p>
    `
  }
  
  return templates[name] || templates['default']
}
```

## 4) DSGVO-Compliance

### 4.1) Consent Management

```typescript
interface ConsentRecord {
  userId: string
  email: string
  consents: {
    marketing: boolean
    newsletter: boolean
    thirdParty: boolean
  }
  grantedAt: Date
  source: 'registration' | 'preferences' | 'thirdParty'
}

// Check consent before sending
async function canSendEmail(userId: string, category: string) {
  const consent = await db.query(
    'SELECT * FROM email_consents WHERE user_id = $1',
    [userId]
  )
  
  if (category === 'marketing' && !consent.marketing) {
    return false
  }
  
  return true
}
```

### 4.2) Unsubscribe

```typescript
// One-click unsubscribe
app.get('/unsubscribe', async (req, res) => {
  const { token } = req.query
  
  // Verify token
  const user = await verifyUnsubscribeToken(token)
  
  // Update consent
  await db.query(
    'UPDATE email_consents SET marketing = false, newsletter = false WHERE user_id = $1',
    [user.id]
  )
  
  // List-Unsubscribe header support
  res.set('List-Unsubscribe', `<mailto:unsubscribe@biometrics.de?subject=unsubscribe>`)
  
  res.send('You have been unsubscribed.')
})
```

## 5) Bounce-Handling

### 5.1) Webhook

```typescript
// Resend Webhook Handler
app.post('/api/webhooks/resend', async (req, res) => {
  const { type, data } = req.body
  
  if (type === 'email.bounced') {
    await handleBounce({
      email: data.email,
      reason: data.bounce?.diagnosticCode,
      type: data.bounce?.bounceType // hard or soft
    })
  }
  
  if (type === 'email.complained') {
    await handleComplaint(data.email)
  }
  
  res.json({ received: true })
})

async function handleBounce(bounce: BounceData) {
  // Soft bounce: retry later
  if (bounce.type === 'soft') {
    await redis.zadd('email:retry', Date.now() + 3600000, bounce.email)
    return
  }
  
  // Hard bounce: disable email
  await db.query(
    'UPDATE users SET email_verified = false, email_bounced = true WHERE email = $1',
    [bounce.email]
  )
}
```

## 6) Monitoring

| Metric | Alert | Dashboard |
|--------|-------|-----------|
| emails.sent.total | - | Grafana |
| emails.sent.failed | > 5% | PagerDuty |
| emails.opened.rate | < 20% | Slack |
| emails.clicked.rate | < 5% | Slack |
| bounce.rate | > 10% | PagerDuty |

## 7) SPF/DKIM/DMARC

| Record | Value | Status |
|--------|-------|--------|
| SPF | v=spf1 include:resend.dev ~all | ✅ Configured |
| DKIM | resend._domainkey | ✅ Configured |
| DMARC | v=DMARC1; p=quarantine; rua=mailto:dmarc@biometrics.de | ✅ Configured |

---

## Abnahme-Check EMAIL-SERVICE
1. E-Mail-Anbieter konfiguriert
2. Templates definiert
3. DSGVO-Consent implementiert
4. Bounce-Handling vorhanden
5. Monitoring konfiguriert

---
