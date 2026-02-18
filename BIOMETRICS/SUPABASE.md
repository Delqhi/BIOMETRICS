# SUPABASE.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Supabase-Betrieb folgt globalen RLS-, Backup- und Audit-Anforderungen.
- Schema-/Policy-Änderungen erfordern verifizierte Impact- und Mapping-Updates.
- Datenhaltung und Retention bleiben compliance-fähig dokumentiert.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Leitfaden für Supabase als primäres Backend (DB/Auth/Storage/Edge).

## 1) Grundsatz
Supabase wird als zentrale Daten- und Auth-Schicht verwendet. Alle Strukturen sind projektagnostisch als Vorlage formuliert.

## 2) Domänenmodell (Template)

| Tabelle | Zweck | Primärschlüssel | Beziehungen | Sensitivität |
|---|---|---|---|---|
| users | Nutzerstammdaten | id | profiles, sessions | hoch |
| projects | Projektmetadaten | id | assets, tasks | mittel |
| assets | Content-Artefakte | id | projects | mittel |
| tasks | Aufgaben und Status | id | projects | niedrig |
| audit_logs | Nachvollziehbarkeit | id | users | hoch |

## 3) Spaltenstandard
Für jede Tabelle dokumentieren:
1. Spaltenname
2. Datentyp
3. Nullbarkeit
4. Default
5. Validierungsregel
6. Beispielwert

## 4) RLS-Standard
Für jede Tabelle definieren:
1. SELECT Policy
2. INSERT Policy
3. UPDATE Policy
4. DELETE Policy
5. Rollenbezug

## 5) Auth-Standard
- Rollenmodell: user, dev, admin, agent
- Session-Handling und Ablaufzeiten dokumentieren
- Recovery-Flows dokumentieren

## 6) Storage-Standard
- Bucket-Zwecke dokumentieren
- Zugriff je Rolle definieren
- Upload-/Download-Richtlinien festlegen

## 7) Edge Functions (Template)

| Function | Zweck | Trigger | Eingabe | Ausgabe | Sicherheitsgrenze |
|---|---|---|---|---|---|
| fn_task_sync | Task-Sync | API call | payload | status | role-check |
| fn_asset_validate | Asset-Prüfung | upload | asset_ref | score | policy-check |

## 8) Migrationsstrategie
1. Änderungen versioniert
2. Rollback-Pfad definiert
3. Datenintegritätscheck nach Migration

## 9) NLM-Datenbezug
- NLM-Assets in `assets` strukturieren
- Qualitätsscore pro Asset erfassen
- Verwendungsstatus pflegen (draft/review/approved/retired)

## 10) Backup/Restore (Template)
- Backup-Intervall: {BACKUP_INTERVAL}
- RPO/RTO Ziele: {RPO_RTO}
- Restore-Testintervall: {RESTORE_TEST_INTERVAL}

## 11) Verifikation
- Schema-Review
- RLS-Review
- Auth-Flow Review
- Migrationstest

## Abnahme-Check SUPABASE
1. Tabellenkatalog vollständig
2. RLS pro Tabelle definiert
3. Auth- und Storage-Strategie dokumentiert
4. Migrations-/Backupprozess vorhanden
5. NLM-Asset-Tracking berücksichtigt

---

## 13) Qwen 3.5 Edge Functions Integration

Qwen 3.5 (NVIDIA NIM) wird für spezialisierte KI-Aufgaben als Edge Functions integriert. Die Skills sind in `AGENTS.md` definiert und werden via Supabase Edge Functions ausgeführt.

### 13.1) Verfügbare Qwen 3.5 Skills

| Skill | Zweck | Trigger | Input | Output |
|---|---|---|---|---|
| qwen_vision_analysis | Bildanalyse | API call | Bilder (PNG, JPG, WebP) | Strukturierte Analyse mit Tags und Metriken |
| qwen_code_generation | Code-Generierung | API call | Natürliche Sprache / Spezifikation | Fertiger Code (Next.js, Go, Supabase) |
| qwen_document_ocr | Texterkennung | API call | PDF, Bilder mit Text | Extrahierter Text, Metadaten, Struktur |
| qwen_video_understanding | Video-Analyse | API call | Videos (MP4, MOV, WebM) | Szenenbeschreibung, Key-Frames, Metadaten |
| qwen_conversation | Konversations-KI | API call | Benutzer-Nachrichten, Kontext | Kontextbezogene Antworten |

### 13.2) Datenbank-Schema für AI-Responses

```sql
-- AI-Request-Log
CREATE TABLE ai_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  skill_name TEXT NOT NULL,
  input_type TEXT NOT NULL,
  input_data JSONB,
  output_data JSONB,
  model_used TEXT NOT NULL,
  tokens_used INTEGER,
  processing_time_ms INTEGER,
  status TEXT NOT NULL DEFAULT 'pending',
  error_message TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  completed_at TIMESTAMPTZ,
  user_id UUID REFERENCES auth.users(id),
  session_id TEXT,
  metadata JSONB DEFAULT '{}'::jsonb
);

-- AI-Response-Cache
CREATE TABLE ai_response_cache (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cache_key TEXT NOT NULL UNIQUE,
  skill_name TEXT NOT NULL,
  input_hash TEXT NOT NULL,
  output_data JSONB NOT NULL,
  hit_count INTEGER DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  last_hit_at TIMESTAMPTZ,
  expires_at TIMESTAMPTZ DEFAULT (NOW() + INTERVAL '24 hours')
);

-- AI-Analysis-Results
CREATE TABLE ai_analysis_results (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  request_id UUID REFERENCES ai_requests(id),
  analysis_type TEXT NOT NULL,
  confidence_score REAL,
  tags TEXT[],
  metadata JSONB DEFAULT '{}'::jsonb,
  raw_response JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- RLS Policies für AI-Tabellen
ALTER TABLE ai_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_response_cache ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_analysis_results ENABLE ROW LEVEL SECURITY;

-- User können ihre eigenen Requests sehen
CREATE POLICY "Users can view own ai_requests"
  ON ai_requests FOR SELECT
  USING (auth.uid() = user_id);

-- AI-Service kann alle Requests lesen
CREATE POLICY "AI service can manage ai_requests"
  ON ai_requests FOR ALL
  USING (auth.jwt()->>'role' = 'ai_service');
```

### 13.3) Edge Function Implementation (Beispiel: qwen_vision_analysis)

```typescript
// supabase/functions/qwen-vision-analysis/index.ts
import { serve } from 'https://deno.land/std@0.168.0/http/server.ts'
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2'

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type',
}

serve(async (req) => {
  if (req.method === 'OPTIONS') {
    return new Response('ok', { headers: corsHeaders })
  }

  try {
    const { image_url, analysis_type, options } = await req.json()
    
    // Supabase Client
    const supabaseAdmin = createClient(
      Deno.env.get('SUPABASE_URL')!,
      Deno.env.get('SUPABASE_SERVICE_ROLE_KEY')!
    )

    // Request loggen
    const { data: requestLog, error: logError } = await supabaseAdmin
      .from('ai_requests')
      .insert({
        skill_name: 'qwen_vision_analysis',
        input_type: 'image',
        input_data: { image_url, analysis_type, options },
        model_used: 'qwen/qwen3.5-397b-a17b',
        status: 'processing',
      })
      .select()
      .single()

    if (logError) throw logError

    // Qwen 3.5 API Aufruf (NVIDIA NIM)
    const nvidiaResponse = await fetch(
      'https://integrate.api.nvidia.com/v1/chat/completions',
      {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${Deno.env.get('NVIDIA_API_KEY')}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          model: 'qwen/qwen3.5-397b-a17b',
          messages: [
            {
              role: 'user',
              content: [
                { type: 'text', text: `Analyze this image for ${analysis_type}. Provide structured JSON with tags, confidence scores, and metadata.` },
                { type: 'image_url', image_url: { url: image_url } }
              ]
            }
          ],
          temperature: 0.2,
          max_tokens: 2048,
        }),
      }
    )

    if (!nvidiaResponse.ok) {
      throw new Error(`NVIDIA API error: ${nvidiaResponse.status}`)
    }

    const nvidiaData = await nvidiaResponse.json()
    const analysisResult = JSON.parse(nvidiaData.choices[0].message.content)

    // Ergebnisse speichern
    const { data: result, error: resultError } = await supabaseAdmin
      .from('ai_analysis_results')
      .insert({
        request_id: requestLog.id,
        analysis_type,
        confidence_score: analysisResult.confidence,
        tags: analysisResult.tags,
        metadata: analysisResult.metadata || {},
        raw_response: nvidiaData,
      })
      .select()
      .single()

    if (resultError) throw resultError

    // Request als abgeschlossen markieren
    await supabaseAdmin
      .from('ai_requests')
      .update({
        status: 'completed',
        output_data: analysisResult,
        tokens_used: nvidiaData.usage?.total_tokens,
        processing_time_ms: Date.now() - new Date(requestLog.created_at).getTime(),
        completed_at: new Date().toISOString(),
      })
      .eq('id', requestLog.id)

    return new Response(
      JSON.stringify({ success: true, result, analysis: analysisResult }),
      { headers: { ...corsHeaders, 'Content-Type': 'application/json' } }
    )

  } catch (error) {
    return new Response(
      JSON.stringify({ success: false, error: error.message }),
      { status: 500, headers: { ...corsHeaders, 'Content-Type': 'application/json' } }
    )
  }
})
```

### 13.4) Deployment

```bash
# Edge Function deployen
supabase functions deploy qwen-vision-analysis

# Mit Secrets
supabase secrets set NVIDIA_API_KEY=nvapi-xxx
```

### 13.5) Aufruf aus Frontend

```typescript
const response = await fetch('https://[PROJECT].supabase.co/functions/v1/qwen-vision-analysis', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${session.access_token}`,
  },
  body: JSON.stringify({
    image_url: 'https://example.com/product.jpg',
    analysis_type: 'product_quality',
    options: { strict: true }
  })
})

const { success, result, analysis } = await response.json()
```

Edge Functions sind die primäre Ausführungsschicht für OpenClaw Skills im Serverless Proxy Pattern.

**Integration Pattern:**
- OpenClaw Skill ruft Edge Function via HTTP auf
- Edge Function führt Datenbank-Operationen aus
- Rückgabe als typed JSON für AI-Consumption

**Meta-Builder Protocol:**
Der Agent kann autonom neue Edge Functions schreiben und deployen via `deploy_supabase_function` Master-Skill.

**Siehe auch:** `WORKFLOW.md` für vollständige Architektur-Dokumentation.
