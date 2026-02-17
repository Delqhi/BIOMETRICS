# WORKFLOW.md — Unified Skill Architecture & Meta-Builder Protocol

**Status:** ACTIVE  
**Version:** 1.0  
**Stand:** Februar 2026  
**Purpose:** Zentrale Dokumentation der Self-Building AI Agent Architektur

---

## 🧠 1. Core Concept: Brain vs Muscle

**BRAIN (AI/OpenClaw):**
- Entscheidungsfindung
- Interface zum User
- Skill-Orchestrierung
- Meta-Cognition (baut sich selbst neue Tools)

**MUSCLE (Supabase/n8n/SDKs):**
- Ausführung komplexer Logik
- Datenbank-Operationen
- Workflow-Automation
- Rechenintensive Tasks

---

## 🏗️ 2. Architecture Patterns

### Pattern A: Webhook Wrapper (n8n Integration)

**Use Case:** Multi-step Prozesse, externe APIs, komplexe Workflows

```
User Request → OpenClaw Skill → Webhook → n8n Workflow → Clean JSON Response
```

### Pattern B: Serverless Proxy (Supabase Edge Functions)

**Use Case:** Datenbank-Operationen, Auth, rechenintensive Tasks

```
User Request → OpenClaw Skill → Supabase Edge Function → Database → Typed Response
```

### Pattern C: SDK Native (Direct Library Usage)

**Use Case:** Lokale Operationen, einfache Tasks, maximale Geschwindigkeit

```
User Request → OpenClaw Skill → SDK/Library → Immediate Result
```

---

## 🤖 3. Meta-Builder Protocol (Advanced)

### Das ultimative Ziel

**Der Agent soll nicht nur Tools BENUTZEN, sondern sich selbst neue Tools BAUEN.**

### Der Loop

1. **DETECT** - Agent erkennt repetitive manuelle Task
2. **ARCHITECT** - Agent designed Lösung (Edge Function/n8n/Skill)
3. **BUILD & DEPLOY** - Agent implementiert und deployed via API
4. **INTEGRATE** - Agent registriert neuen Skill für sich selbst
5. **REPEAT** - Agent wird kontinuierlich mächtiger

### Master-Skills (Gott-Modus)

1. **deploy_n8n_workflow** - Erstellt autonome n8n Workflows
2. **deploy_supabase_function** - Deployt TypeScript Edge Functions
3. **register_openclaw_skill** - Fügt neue Skills zum eigenen Skill-Set hinzu

---

## 📋 4. Best Practices

### Strict Typing (Zod)
Jede Skill-Input/Output muss validiert werden.

### "Return for AI" (Clean Outputs)
Jede Execution muss AI-freundlich zurückgeben (clean JSON, keine Errors ohne Kontext).

### Idempotency
Skills müssen wiederholbar sein ohne Seiteneffekte.

---

## 🎯 5. Real-World Beispiel

**User:** "Überwach meine Konkurrenten auf Preisänderungen"

**Agent denkt:**
- Wiederkehrende Aufgabe (Polling alle 6h)
- Braucht autonomen Scraper

**Agent handelt:**
1. Baut Supabase Edge Function (Scraper)
2. Erstellt n8n Workflow (alle 6h Trigger)
3. Registriert Skill `check_competitor_prices`

**Agent Antwort:**
> "Erledigt. Scraper läuft alle 6 Stunden autonom. Frag mich jederzeit nach aktuellen Preisen."

---

## 📚 6. Referenzen

- Architecture: `ARCHITECTURE.md`
- Supabase: `SUPABASE.md`
- n8n: `N8N.md`
- OpenClaw: `OPENCLAW.md`
- Agents: `AGENTS-GLOBAL.md`

---

**Version:** 1.0  
**Stand:** 2026-02-17  
**Status:** PRODUCTION READY ✅
