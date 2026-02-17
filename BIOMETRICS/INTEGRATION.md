# INTEGRATION.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Integrationen folgen globalen Security-, Secrets- und Observability-Regeln.
- Jede Schnittstelle benötigt Vertrag, Ownership und Incident-Pfad.
- Änderungen sind inklusive Mapping und Betriebsnachweis zu dokumentieren.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Integrationsrahmen für externe Systeme und interne Automationspfade.

## Integrationsprinzipien
1. API-first
2. Explizite Verträge
3. Fehlerpfade dokumentieren
4. Idempotente Verarbeitung wo möglich
5. Auditierbare Übergaben

## Integrationsmatrix (Template)

| Integration | Zweck | Richtung | Auth | Kritikalität | Owner |
|---|---|---|---|---|---|
| OpenClaw | Connector/Auth Layer | bidirektional | token/role | hoch | {OWNER_ROLE} |
| n8n | Workflow Automation | bidirektional | service auth | mittel | {OWNER_ROLE} |
| NLM-CLI | Content-Erzeugung | outbound | local policy | hoch | {OWNER_ROLE} |
| Cloudflare | Netz-/Edge Layer | ingress | platform auth | hoch | {OWNER_ROLE} |

## 1) OpenClaw (Template)
- Zweck: {OPENCLAW_PURPOSE}
- Auth-Flow: {OPENCLAW_AUTH_FLOW}
- Retry-Strategie: {OPENCLAW_RETRY}
- Fehlerbehandlung: {OPENCLAW_ERROR_POLICY}

## 2) n8n (Template)
- Workflow-Kategorien: {N8N_WORKFLOW_TYPES}
- Trigger: {N8N_TRIGGERS}
- Recovery: {N8N_RECOVERY}
- Observability: {N8N_OBSERVABILITY}

## 3) NLM-CLI (Pflicht)
- Asset-Typen: Video, Infografik, Präsentation, Datentabelle
- Vorlagenquelle: `../∞Best∞Practices∞Loop.md`
- Qualitätsprüfung: NLM-Matrix 13/16, Korrektheit 2/2
- Delegationsprotokoll: `MEETING.md`

## 4) Fehler- und Eskalationsmodell
- P0: sofortige Eskalation
- P1: innerhalb der Session lösen/eskalieren
- P2: in nächsten Zyklus einplanen

## 5) Vertragsmodell pro Integration

| Feld | Beschreibung |
|---|---|
| Input Contract | Eingabestruktur und Validierung |
| Output Contract | erwartete Ausgabe |
| Error Contract | Fehlerklassen und Codes |
| Timeout/Retry | Robustheitsparameter |
| Security Boundary | Zugriffsgrenzen |

## 6) Benennungsstandard für Assets
Pattern:
`{project}_{asset-type}_{topic}_{locale}_{version}`

## 7) Verifikation
- Contract-Check je Integration
- Fehlerpfad-Test
- Security-Review
- Doku-Synchronität mit `COMMANDS.md` und `ENDPOINTS.md`

## Abnahme-Check INTEGRATION
1. Integrationsmatrix vollständig
2. NLM-CLI Prozess enthalten
3. Error/Retry Regeln dokumentiert
4. Vertragsmodell vorhanden
5. Asset-Naming standardisiert

---
