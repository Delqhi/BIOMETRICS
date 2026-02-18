# VERCEL.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Plattformbetrieb folgt globalen Deployment- und Security-Leitlinien.
- Release- und Rollback-Prozesse sind testbar und dokumentiert.
- Konfigurationsdrift wird regelmäßig geprüft und korrigiert.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Leitfaden für Vercel-Projektbetrieb, Environments und Deploy-Governance.

## Hinweis
Wenn Vercel nicht genutzt wird, Status auf `NOT_APPLICABLE` setzen und Begründung dokumentieren.

## Projekt-Metadaten (Template)
- Projektname: {PROJECT_NAME}
- Projekt-ID: {VERCEL_PROJECT_ID}
- Team: {VERCEL_TEAM}
- Owner: {OWNER_ROLE}

## Environment-Strategie
- development
- preview
- production

## Deploy-Regeln
1. Deploy nur über kontrollierten CI/CD-Pfad
2. Produktionsdeploy nur bei grünen Gates
3. Rollback-Option vor Deploy verifizieren

## Domain-/Routing-Template
- Primärdomain: {PRIMARY_DOMAIN}
- Zusätzliche Domains: {ADDITIONAL_DOMAINS}
- Redirect-Policy: {REDIRECT_POLICY}

## Sicherheitsregeln
- Secrets nur über Environment-Management
- Keine Secrets in Repo-Dateien
- Zugriff auf Production minimal halten

## Verifikation
- Preview Deploy erfolgreich
- Production Health-Check erfolgreich
- Rollback-Test dokumentiert

## Betriebscheckliste
### Pre-Deploy
1. Alle CI-Gates grün
2. Ziel-Environment bestätigt
3. Rollback-Referenz dokumentiert

### Deploy
1. Deploy-Trigger dokumentiert
2. Versionsreferenz notiert
3. Status-Checks beobachtet

### Post-Deploy
1. Health-Endpunkte geprüft
2. Kernjourney-Sanity-Check durchgeführt
3. Fehler-/Latenzbaseline geprüft

### Rollback
1. Triggerkriterium erreicht?
2. Rollback ausgelöst
3. Nach-Rollback Health-Check erfolgreich

## Abnahme-Check VERCEL
1. Projekt-/Environmentdaten als Platzhalter vorhanden
2. Deploy- und Rollback-Regeln dokumentiert
3. Sicherheitsregeln enthalten

---

## Qwen 3.5 Edge Functions Deployment

### Projekt-Metadaten
- Projektname: BIOMETRICS
- Projekt-ID: biomet-rics-01
- Team: SIN-Enterprise
- Owner: DevOps

### Serverless Deployment
Qwen 3.5 Edge Functions werden als Vercel Serverless Functions deployed.

**Architektur:**
```
┌─────────────┐     ┌──────────────────┐     ┌─────────────┐
│   Client   │────▶│  Vercel Edge     │────▶│   NVIDIA    │
│  (Browser) │     │  Function        │     │   NIM API   │
└─────────────┘     │  (Qwen 3.5)     │     └─────────────┘
                    └──────────────────┘
```

**Edge Function Struktur:**
```
api/
├── qwen/
│   ├── chat/
│   │   └── route.ts        # Chat Completion
│   ├── vision/
│   │   └── route.ts        # Bildanalyse
│   ├── ocr/
│   │   └── route.ts        # Texterkennung
│   └── video/
│       └── route.ts        # Video-Analyse
```

### Environment-Konfiguration

**Vercel Environment Variables:**
```bash
# NVIDIA NIM (Pflicht)
NVIDIA_API_KEY=nvapi-xxxxxxxxxxxxxxxxxxxx

# Model Configuration
QWEN_MODEL_ID=qwen/qwen3.5-397b-a17b
QWEN_BASE_URL=https://integrate.api.nvidia.com/v1
QWEN_MAX_TOKENS=32768
QWEN_TIMEOUT=120000

# Feature Flags
ENABLE_STREAMING=true
ENABLE_VISION=true
ENABLE_VIDEO=true
```

**Environment-Strategie (Edge Functions):**
| Variable | development | preview | production |
|----------|-------------|---------|-------------|
| NVIDIA_API_KEY | ✅ | ✅ | ✅ |
| QWEN_MAX_TOKENS | 8192 | 16384 | 32768 |
| ENABLE_STREAMING | true | true | true |
| LOG_LEVEL | debug | info | warn |

### API Routes

#### 1. Chat Completion
**Endpoint:** `POST /api/qwen/chat`

```typescript
// api/qwen/chat/route.ts
import { NextRequest, NextResponse } from 'next/server';

export async function POST(req: NextRequest) {
  const { messages, model, temperature, max_tokens } = await req.json();
  
  const response = await fetch(process.env.QWEN_BASE_URL + '/chat/completions', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      model: process.env.QWEN_MODEL_ID,
      messages,
      temperature: temperature ?? 0.7,
      max_tokens: max_tokens ?? process.env.QWEN_MAX_TOKENS,
      stream: process.env.ENABLE_STREAMING === 'true',
    }),
  });

  if (process.env.ENABLE_STREAMING === 'true') {
    return new Response(response.body, {
      headers: { 'Content-Type': 'text/event-stream' },
    });
  }

  return NextResponse.json(await response.json());
}
```

#### 2. Vision Analysis
**Endpoint:** `POST /api/qwen/vision`

```typescript
// api/qwen/vision/route.ts
export async function POST(req: NextRequest) {
  const { image, prompt } = await req.json();
  
  const response = await fetch(process.env.QWEN_BASE_URL + '/chat/completions', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      model: process.env.QWEN_MODEL_ID,
      messages: [
        {
          role: 'user',
          content: [
            { type: 'text', text: prompt },
            { type: 'image_url', image_url: { url: image } },
          ],
        },
      ],
      max_tokens: 4096,
    }),
  });

  return NextResponse.json(await response.json());
}
```

#### 3. OCR (Texterkennung)
**Endpoint:** `POST /api/qwen/ocr`

```typescript
// api/qwen/ocr/route.ts
export async function POST(req: NextRequest) {
  const { image, language = 'auto' } = await req.json();
  
  const response = await fetch(process.env.QWEN_BASE_URL + '/chat/completions', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      model: process.env.QWEN_MODEL_ID,
      messages: [
        {
          role: 'user',
          content: [
            { 
              type: 'text', 
              text: 'Extract all text from this image. Return only the text content.' 
            },
            { type: 'image_url', image_url: { url: image } },
          ],
        },
      ],
      max_tokens: 8192,
    }),
  });

  return NextResponse.json(await response.json());
}
```

#### 4. Video Understanding
**Endpoint:** `POST /api/qwen/video`

```typescript
// api/qwen/video/route.ts
export async function POST(req: NextRequest) {
  const { video_frames, prompt } = await req.json();
  
  const content = video_frames.map((frame: string) => ({
    type: 'image_url',
    image_url: { url: frame },
  }));
  content.unshift({ type: 'text', text: prompt });

  const response = await fetch(process.env.QWEN_BASE_URL + '/chat/completions', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      model: process.env.QWEN_MODEL_ID,
      messages: [{ role: 'user', content }],
      max_tokens: 16384,
    }),
  });

  return NextResponse.json(await response.json());
}
```

### Deploy-Prozess

**1. Vercel CLI Installation:**
```bash
npm i -g vercel
vercel login
```

**2. Projekt verknüpfen:**
```bash
cd /Users/jeremy/dev/BIOMETRICS/BIOMETRICS
vercel link
```

**3. Environment Variables setzen:**
```bash
vercel env add NVIDIA_API_KEY
vercel env add QWEN_MODEL_ID
vercel env add QWEN_BASE_URL
```

**4. Deployment:**
```bash
# Development
vercel dev

# Preview (automatisch bei PR)
vercel --prod

# Production
vercel --prod --yes
```

### Monitoring & Logs

**Vercel Dashboard:**
- Functions: `/dashboard/biomet-rics-01/functions`
- Logs: `/dashboard/biomet-rics-01/functions/.../logs`

**Edge Function Metrics:**
- Latency (P50, P95, P99)
- Invocations
- Error Rate
- Cold Starts

### Rate Limits

**NVIDIA NIM Limits (Free Tier):**
- RPM: 40 requests/minute
- Fallback: Queue mit Retry-Logic

**Vercel Edge Functions:**
- Concurrent Executions: 1000
- Duration: 30s (max)

---

## Abnahme-Check Qwen Edge
1. ✅ Environment Variables dokumentiert
2. ✅ API Routes implementiert
3. ✅ NVIDIA NIM Integration konfiguriert
4. ✅ Deploy-Prozess dokumentiert
5. ✅ Monitoring-Sektion vorhanden
