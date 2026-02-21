# ENGINE.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Engine-Regeln folgen globalen Betriebs-, Sicherheits- und Qualitätsvorgaben.
- Ausführungslogik muss testbar, observierbar und dokumentiert bleiben.
- Kritische Änderungen brauchen Incident- und Recovery-Anschlussfähigkeit.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Rahmen für eine projektspezifische Core-Engine oder zentrale Verarbeitungslogik.

## Hinweis
Falls keine Engine-Komponente existiert, Status auf `NOT_APPLICABLE` setzen und begründen.

## Engine-Verantwortung (Template)
- Inputverarbeitung
- Regelanwendung
- Ergebnisorchestrierung
- Fehlerhärtung

## Architektur-Schnittstellen

| Schnittstelle | Eingabe | Ausgabe | Fehlerverhalten |
|---|---|---|---|
| API -> Engine | payload | response | validation errors |
| Engine -> Storage | normalized data | write status | retry policy |
| Engine -> Integrationslayer | job context | integration result | escalation |

## Laufzeitgrenzen
- erwartete Lastprofile: {LOAD_PROFILE}
- Timeout-Strategie: {TIMEOUT_POLICY}
- Retry-Strategie: {RETRY_POLICY}

## Qualitätsziele
1. deterministisches Verhalten in Kernpfaden
2. klar dokumentierte Fehlerzustände
3. reproduzierbare Ergebnisse

## Observability
- zentrale Metriken
- Fehlerklassifikation
- Durchsatz-/Latenzwerte

## Verifikation
- Core-Path Tests
- Fehlerpfadtests
- Lastprofil-Sanity-Checks

## Abnahme-Check ENGINE
1. Verantwortungen klar
2. Schnittstellen dokumentiert
3. Laufzeit- und Qualitätsziele vorhanden
4. Verifikation definiert

## Qwen 3.5 Code Generation

### qwen_code_generation
Automatische Code-Generierung für die Engine-Komponente.

| Feature | Beschreibung | Output |
|---------|--------------|--------|
| API-Routen | REST-Endpoints generieren | Go/Next.js Code |
| Datenbank-Schema | Supabase Tables/Functions | SQL + TypeScript |
| Edge-Functions | Serverless Logik | Deno/TypeScript |
| Tests | Unit- und Integrationstests | Jest/Testing Library |

### Code-Generation Pipeline
```
Natürliche Sprache Spezifikation
           ↓
    qwen_code_generation
           ↓
┌──────────┼──────────┐
↓          ↓          ↓
API Route  DB Schema  Edge Function
           ↓
    Automatische Tests
           ↓
    Deployment-Ready
```

### API-Integration
```typescript
// Code-Generierung
const generated = await fetch('/api/qwen/chat', {
  method: 'POST',
  body: JSON.stringify({
    messages: [{
      role: 'user',
      content: `Generiere eine Go-Handler-Funktion für:
      - Endpoint: POST /api/engine/process
      - Input: JSON mit user_id, action, payload
      - Output: JSON mit status, result, error
      - Auth: JWT-Token erforderlich`
    }],
    skill: 'qwen_code_generation',
    options: {
      language: 'go',
      framework: 'gin',
      generateTests: true
    }
  })
});
```

### Engine-spezifische Use Cases
1. **Batch-Verarbeitung**: Generiere Worker-Pool mit Retry-Logik
2. **Event-Handler**: Erstelle Event-Consumer für Queue
3. **Data-Transformer**: Konvertiere zwischen Formaten
4. **Validation-Layer**: Regelsätze für Input-Prüfung

### Qualitäts-Gates
- Keine Syntax-Fehler (TypeScript strict)
- Tests bestehen > 90% Coverage
- Security-Scan ohne Findings
- Performance: < 100ms Latenz

---
