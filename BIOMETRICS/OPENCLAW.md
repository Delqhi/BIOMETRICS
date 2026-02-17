# OPENCLAW.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- OpenClaw-Konfiguration folgt globalen Provider- und Timeout-Regeln.
- Modell-ID-Disziplin, Fallback-Strategie und Betriebsmetriken sind Pflicht.
- Änderungen werden sicher, nachvollziehbar und testsynchron umgesetzt.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Integrationsleitfaden für OpenClaw als Connector- und Auth-Layer.

## Rolle im System
OpenClaw dient als Integrationsschicht für externe Plattformzugriffe und orchestrierte Connector-Flows.

## Grundprinzipien
1. Keine direkten unsicheren Provider-Zugriffe
2. Token- und Rollenhygiene
3. Klare Fehlerpfade
4. Wiederholbare Connector-Prozesse

## Betriebsmodi
- local: Entwicklungsmodus
- staging: Integrationsvalidierung
- production: stabiler Betrieb

## Konfigurationsmatrix (Template)

| Variable | Zweck | Pflicht | Beispiel |
|---|---|---|---|
| OPENCLAW_BASE_URL | Basis-URL | ja | {URL} |
| OPENCLAW_AUTH_MODE | Auth-Modus | ja | token |
| OPENCLAW_TIMEOUT_MS | Timeout | ja | 10000 |
| OPENCLAW_RETRY_COUNT | Retry Anzahl | ja | 3 |

## Auth- und Tokenfluss
1. Zugriffsanfrage wird validiert
2. Rolle und Scope geprüft
3. Token wird aus sicherer Quelle gelesen
4. Request wird signiert und gesendet
5. Response wird protokolliert und bewertet

## Fehler- und Retry-Strategie
- temporäre Fehler: Retry mit Backoff
- dauerhafte Fehler: sofortige Eskalation
- Auth-Fehler: Tokenfluss prüfen, ggf. Rotation

## Sicherheitsregeln
1. Keine Tokens im Repo
2. Kein Logging sensibler Daten
3. Least Privilege für Connectoren

## Verifikation
- Health Check für Connectorpfade
- Auth-Flow Tests
- Fehlerfalltests (Timeout, 401, 403, 5xx)

## Abnahme-Check OPENCLAW
1. Konfigurationsmatrix vorhanden
2. Authfluss dokumentiert
3. Retry-/Fehlerstrategie definiert
4. Sicherheitsregeln enthalten
5. Verifikationsplan vorhanden

---

---

## 12) OpenClaw Skills als Interfaces

OpenClaw Skills sind die "Interfaces" oder "Wrappers" für zugrundeliegende Logik (Supabase/n8n/SDKs).

**Skill Types:**
1. **Webhook Wrapper** - Trigger n8n Workflows
2. **Serverless Proxy** - Call Supabase Edge Functions
3. **SDK Native** - Direct library usage

**Ultimate Goal:**
OpenClaw soll sich selbst replizieren durch autonomes Erstellen neuer Skills via Meta-Builder Protocol.

**Master-Skills für Self-Replication:**
- `deploy_n8n_workflow`
- `deploy_supabase_function`
- `register_openclaw_skill`

**Siehe auch:** `WORKFLOW.md` für Skill creation architecture.
