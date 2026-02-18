# SMS-SERVICE.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- SMS-Versand erfolgt DSGVO-konform mit Consent-Management.
- Nur für kritische Benachrichtigungen (2FA, Lieferung, Notfall).
- Opt-Out-Mechanismus ist Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

## Zweck
Dokumentation der SMS-Service-Architektur und -Konfiguration für BIOMETRICS.

## 1) SMS-Anbieter

| Anbieter | Status | Use Case | Kosten |
|----------|--------|----------|--------|
| Twilio | ✅ Primary | Alle SMS | ~€0.08/SMS |
| MessageBird | ❌ Fallback | Nicht konfiguriert | - |

## 2) SMS-Typen

| Template | Trigger | Kategorie | Priorität |
|----------|---------|-----------|-----------|
| verification-code | 2FA Login | Sicherheit | CRITICAL |
| order-shipped | Versandbestätigung | Transaktional | HIGH |
| delivery-otp | Lieferung verifizieren | Sicherheit | HIGH |
| payment-alert | Zahlungsproblem | Transaktional | HIGH |
| account-alert | Sicherheitswarnung | Sicherheit | CRITICAL |
| marketing-promo | Werbe-SMS | Marketing | LOW |

## 3) Implementierung

### 3.1) Twilio Client

```typescript
import twilio from 'twilio'

const client = twilio(
  process.env.TWILIO_ACCOUNT_SID,
  process.env.TWILIO_AUTH_TOKEN
)

interface SMSOptions {
  to: string
  template: string
  variables?: Record<string, string>
  type: 'verification' | 'transactional' | 'marketing'
}

async function sendSMS(options: SMSOptions) {
  const { to, template, variables, type } = options
  
  // Validate phone number
  if (!validatePhoneNumber(to)) {
    throw new Error('Invalid phone number')
  }
  
  // Get consent
  const hasConsent = await checkSMSConsent(to, type)
  if (!hasConsent) {
    throw new Error('No consent for SMS')
  }
  
  // Get template content
  const message = await getSMSTemplate(template, variables)
  
  // Send via Twilio
  const result = await client.messages.create({
    body: message,
    from: process.env.TWILIO_PHONE_NUMBER,
    to
  })
  
  // Log for audit
  await logSMS({
    sid: result.sid,
    to,
    type,
    status: result.status,
    template
  })
  
  return result
}
```

### 3.2) 2FA (Verification Code)

```typescript
// Generate and send 2FA code
async function sendVerificationCode(userId: string, phone: string) {
  const code = Math.floor(100000 + Math.random() * 900000).toString()
  
  // Store code with expiry
  await redis.setex(
    `sms:verify:${userId}`,
    300, // 5 minutes
    code
  )
  
  // Send SMS
  await sendSMS({
    to: phone,
    template: 'verification-code',
    variables: { code },
    type: 'verification'
  })
  
  return { success: true, expiresIn: 300 }
}

// Verify code
async function verifyCode(userId: string, code: string): Promise<boolean> {
  const stored = await redis.get(`sms:verify:${userId}`)
  return stored === code
}
```

## 4) Consent Management

### 4.1) Database Schema

```sql
CREATE TABLE sms_consents (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id),
  phone VARCHAR(20) NOT NULL,
  transactional BOOLEAN DEFAULT false,
  marketing BOOLEAN DEFAULT false,
  granted_at TIMESTAMP DEFAULT NOW(),
  source VARCHAR(50) -- registration, preferences, verbal
);

CREATE INDEX idx_sms_consents_phone ON sms_consents(phone);
```

### 4.2) Consent Check

```typescript
async function checkSMSConsent(phone: string, type: string): Promise<boolean> {
  const consent = await db.query(
    `SELECT * FROM sms_consents WHERE phone = $1`,
    [phone]
  )
  
  if (type === 'verification' || type === 'transactional') {
    return consent?.transactional || false
  }
  
  if (type === 'marketing') {
    return consent?.marketing || false
  }
  
  return false
}
```

## 5) Opt-Out

### 5.1) STOP-Keyword

```typescript
// Handle incoming SMS (Twilio Webhook)
app.post('/api/webhooks/sms/incoming', async (req, res) => {
  const { From, Body } = req.body
  
  const message = Body.trim().toUpperCase()
  
  if (['STOP', 'UNSUBSCRIBE', 'NO'].includes(message)) {
    await updateSMSConsent(From, {
      transactional: false,
      marketing: false
    })
    
    res.send('You have been unsubscribed from SMS.')
    return
  }
  
  res.send('')
})
```

## 6) Monitoring

| Metric | Alert | Dashboard |
|--------|-------|-----------|
| sms.sent.total | - | Grafana |
| sms.sent.failed | > 2% | PagerDuty |
| sms.delivered.rate | < 95% | Slack |
| sms.2fa.success | < 90% | Slack |
| sms.bulk.credits | < 100 remaining | Email |

## 7) Kosten-Kontrolle

### 7.1) Budget

| Limit | Wert | Alert |
|-------|------|-------|
| Tageslimit | 100 SMS | 80% |
| Monatslimit | 1000 SMS | 80% |
| 2FA-Max/Tag | 10/User | Automatisch |

### 7.2) Implementierung

```typescript function checkSMSBudget
async(userId: string): Promise<boolean> {
  const today = new Date().toDateString()
  const key = `sms:count:${userId}:${today}`
  
  const count = await redis.incr(key)
  await redis.expire(key, 86400) // 24h
  
  if (count > 10) {
    throw new Error('SMS limit exceeded')
  }
  
  return true
}
```

---

## Abnahme-Check SMS-SERVICE
1. SMS-Anbieter konfiguriert
2. Templates definiert
3. Consent-Management implementiert
4. Opt-Out vorhanden
5. Kosten-Kontrolle aktiv

---
