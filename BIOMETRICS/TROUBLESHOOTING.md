# TROUBLESHOOTING.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Troubleshooting wird ticketbasiert, reproduzierbar und evidenzgestützt geführt.
- Root Cause, Corrective und Preventive Actions sind Pflichtbestandteile.
- Learnings fließen in Regeln, Playbooks und Training zurück.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universelle Fehlerdiagnose- und Behebungsleitfäden für Betrieb und Entwicklung.

---

## 🧠 NLM CLI COMMANDS

```bash
# Create notebook
nlm notebook create "Title"

# List sources
nlm source list <notebook-id>

# Delete old source (BEFORE adding new!)
nlm source delete <source-id> -y

# Add new source
nlm source add <notebook-id> --file "file.md" --wait
```

**⚠️ DUPLICATE PREVENTION:** ALWAYS run `nlm source list` before `nlm source add`!

---

## 🔄 DEQLHI-LOOP (INFINITE WORK MODE)

- After each completed task → Add 5 new tasks immediately
- Never "done" - only "next task"
- Always document → Every change in files
- Git commit + push after EVERY change
- Parallel execution ALWAYS (run_in_background=true)

### Loop Mechanism:
1. Task N Complete
2. Git Commit + Push
3. Update Docs
4. Add 5 New Tasks
5. Next Task N+1
6. Repeat infinitely

## Diagnoseprinzip
1. Reproduzieren
2. Isolieren
3. Root-Cause finden
4. Fix minimal-invasiv umsetzen
5. Regression verhindern
6. Doku aktualisieren

## Fehlerkatalog (Template)

### Fall 01: Build schlägt fehl
Symptom:
- Build-Prozess endet mit Fehler.

Wahrscheinliche Ursachen:
- inkonsistente Abhängigkeiten
- Typfehler
- ungültige Konfiguration

Diagnose:
- `pnpm build`
- `pnpm typecheck`

Lösung:
1. Fehlerursache priorisieren
2. fixen
3. erneut builden

Verifikation:
- Build erfolgreich

### Fall 02: Tests schlagen fehl
Symptom:
- Test-Suite nicht grün.

Wahrscheinliche Ursachen:
- Regression im Code
- geänderte Schnittstellen

Diagnose:
- `pnpm test`
- `go test ./...`

Lösung:
1. failing tests isolieren
2. root cause beheben
3. relevante Tests erneut ausführen

Verifikation:
- betroffene Tests grün

### Fall 03: API-Fehler 401/403
Symptom:
- Zugriff verweigert.

Wahrscheinliche Ursachen:
- falsche Rolle
- fehlende Auth-Konfiguration

Diagnose:
- Auth-Kontext prüfen
- Rollenmapping prüfen

Lösung:
1. Rollen-/Policy-Konfiguration korrigieren
2. Auth-Flow verifizieren

Verifikation:
- Endpoint mit korrekter Rolle erreichbar

### Fall 04: NLM-Output unbrauchbar
Symptom:
- inkonsistent, unklar oder faktisch falsch.

Wahrscheinliche Ursachen:
- unklarer Prompt
- fehlende Quellen
- fehlende Qualitätsprüfung

Diagnose:
- Prompt gegen Vorlagen prüfen
- Quellabdeckung prüfen

Lösung:
1. Prompt präzisieren
2. zweite Iteration erzeugen
3. Qualitätsmatrix anwenden

Verifikation:
- Score >= 13/16 und Korrektheit 2/2

### Fall 05: Dokumente widersprechen sich
Symptom:
- `COMMANDS.md` und `ENDPOINTS.md` inkonsistent.

Wahrscheinliche Ursachen:
- asynchrone Updates
- fehlende Cross-Doc Prüfung

Diagnose:
- Mapping-Check durchführen

Lösung:
1. Primärquelle festlegen
2. Gegenstück synchronisieren
3. in Changelog dokumentieren

Verifikation:
- 1:1 Mapping konsistent

### Fall 06: Qwen 3.5 API Errors
Symptom:
- HTTP 429 Too Many Requests
- HTTP 401 Unauthorized
- HTTP 500 Internal Server Error

Wahrscheinliche Ursachen:
- Rate Limit überschritten (40 RPM Free Tier)
- Ungültiger API Key
- NVIDIA NIM Service outage

Diagnose:
- `curl -H "Authorization: Bearer $NVIDIA_API_KEY" https://integrate.api.nvidia.com/v1/models`
- NVIDIA API Status prüfen

Lösung:
1. Bei HTTP 429: 60 Sekunden warten + Fallback-Modelle nutzen
2. API Key verifizieren: `echo $NVIDIA_API_KEY`
3. Fallback-Chain: qwen/qwen3.5-397b-a17b → qwen2.5-coder-32b
4. Timeout in Config auf 120000ms erhöhen

Verifikation:
- API Call erfolgreich
- Modell antwortet

### Fall 07: Qwen 3.5 Timeout Issues
Symptom:
- Request hängt > 60 Sekunden
- Modell antwortet nicht

Wahrscheinliche Ursachen:
- Qwen 3.5 397B hat extreme Latenz (70-90s)
- Netzwerkprobleme
- Modell overloaded

Diagnose:
- Timeout-Messung durchführen
- Ping zu NVIDIA API

Lösung:
1. Timeout auf 120000ms (120s) setzen (PFLICHT für Qwen 3.5!)
2. Fallback zu schnelleren Modellen: qwen2.5-coder-7b
3. Request-Queue reduzieren
4. Streaming deaktivieren (nicht supported bei NVIDIA NIM)

Verifikation:
- Request innerhalb von 120s abgeschlossen
- Keine Timeout-Fehler mehr

### Fall 08: Qwen 3.5 Model Loading Problems
Symptom:
- Modell kann nicht geladen werden
- "Model not found" Fehler

Wahrscheinliche Ursachen:
- Falsche Modell-ID verwendet
- Modell nicht verfügbar im NVIDIA NIM

Diagnose:
- Verfügbare Modelle listen: `curl https://integrate.api.nvidia.com/v1/models`

Lösung:
1. Korrekte Modell-ID verwenden: `qwen/qwen3.5-397b-a17b`
2. NICHT `qwen2.5` verwenden (falsches Modell!)
3. Provider-Konfiguration prüfen
4. Model-Limit in opencode.json verifizieren

Verifikation:
- Modell geladen und antwortfähig

### Fall 09: Qwen 3.5 Context Window Exceeded
Symptom:
- "context_length_exceeded" Fehler
- Modell schneidet Antworten ab

Wahrscheinliche Ursachen:
- Input überschreitet 262K Token Limit
- Zu viele Nachrichten im Chat-Verlauf

Diagnose:
- Token-Count der Anfrage prüfen
- Chat-History analysieren

Lösung:
1. Context-Mode: summarize oder restart verwenden
2. Nachrichten-History kürzen (nur letzte N behalten)
3. Input kürzen: Nicht alle Dokumente gleichzeitig
4. Chunking: Große Dokumente in Teile aufteilen
5. Besseres Modell: qwen3.5-397b hat 262K, qwen2.5-coder-32b nur 128K

Verifikation:
- Keine Context-Fehler
- Vollständige Antworten

## Eskalations-Playbook
- P0: sofortige Eskalation und Schadensbegrenzung
- P1: innerhalb Session beheben/eskalieren
- P2: nächsten Zyklus planen

## Postmortem-Vorlage
```text
Incident-ID:
Severity:
Timeline:
Root Cause:
Impact:
Fix:
Prävention:
Owner:
Follow-up Tasks:
```

## Abnahme-Check TROUBLESHOOTING
1. Root-Cause Ansatz beschrieben
2. NLM-Fehlerfall enthalten
3. API/Auth-Fälle enthalten
4. Verifikationsschritte je Fall vorhanden
5. Eskalationspfad vorhanden

---
