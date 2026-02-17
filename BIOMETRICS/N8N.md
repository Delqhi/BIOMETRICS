# N8N.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Workflow-Automation folgt globaler Governance und Secrets-Disziplin.
- Trigger, Fehlerpfade und Retry-Verhalten müssen dokumentiert sein.
- Kritische Flows benötigen Incident- und Recovery-Referenzen.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Betriebsleitfaden für n8n-Workflows in produktionsnahen Umgebungen.

## Prinzipien
1. Workflows sind versioniert
2. Fehlerpfade sind explizit
3. Idempotenz wird berücksichtigt
4. Recovery ist dokumentiert

## Workflow-Katalog (Template)

| Workflow-ID | Zweck | Trigger | Kritikalität | Owner | Status |
|---|---|---|---|---|---|
| WF-001 | {PURPOSE} | {TRIGGER} | P1 | {OWNER_ROLE} | active |

## Trigger-Typen
- webhook
- schedule
- event-driven
- manual

## Input/Output-Contract
Für jeden Workflow dokumentieren:
1. Eingabefelder
2. Validierungsregeln
3. Ausgabefelder
4. Fehlercodes
5. Nebenwirkungen

## Fehlerbehandlung
- Retry bei transienten Fehlern
- Dead-letter/Quarantäne bei dauerhaften Fehlern
- Alarmierung bei kritischen Ausfällen

## Recovery-Plan
1. Fehlerfall erkennen
2. betroffenen Workflow pausieren
3. Ursache isolieren
4. Korrektur deployen
5. sicheren Neustart durchführen

## Observability
- Laufzeitmetriken pro Workflow
- Fehlerquote
- Durchsatz
- Queue-Länge

## NLM-Bezug
- NLM-generierte Asset-Aufgaben können n8n-triggerbar modelliert werden
- Freigabe nur nach Qualitätsmatrix
- Ergebnisprotokoll in `MEETING.md`

## Abnahme-Check N8N
1. Workflow-Katalog vorhanden
2. Input/Output-Contract je Workflow definiert
3. Fehler- und Recoverypfade dokumentiert
4. Observability-Basis enthalten

---

---

## 12) n8n als Heavy Lifting Muscle für AI Skills

n8n Workflows sind die "Automation Muscle" die von AI Skills getriggert werden im Webhook Wrapper Pattern.

**Integration Pattern:**
- OpenClaw Skill trigger n8n Webhook
- n8n führt multi-step Workflow aus
- Clean JSON Response für AI

**Design Principles:**
- Workflows müssen AI-triggbar sein (Webhook endpoint)
- Fehler müssen AI-freundlich zurückgegeben werden
- Idempotenz für wiederholbare Execution

**Meta-Builder Protocol:**
Der Agent kann autonom neue n8n Workflows erstellen und deployen via `deploy_n8n_workflow` Master-Skill.

**Siehe auch:** `WORKFLOW.md` für vollständige Architektur-Dokumentation.
