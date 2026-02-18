# MCP.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- MCP-Betrieb folgt globalen Sicherheits- und Wrapper-Regeln.
- Konfigurationen bleiben nachvollziehbar, testbar und ohne Blindänderungen.
- Troubleshooting ist ticket- und evidenzbasiert zu führen.

Status: ACTIVE  
Version: 1.2 (Qwen 3.5 Integration)  
Stand: Februar 2026

## Zweck
Globales MCP-Betriebshandbuch im Projektkontext: welche MCP-Server existieren, wie sie installiert, aktiviert und sicher genutzt werden.

## Pfadstandard
`BIOMETRICS/MCP.md` ist der kanonische Speicherort für das MCP-Book dieses Projekts.

## Grundregeln
1. MCP-Nutzung ist dokumentationspflichtig.
2. Vor Einsatz immer Zweck, Risiko und erwartetes Ergebnis klären.
3. Serverzugriffe mit Least-Privilege konfigurieren.
4. Änderungen und Nutzung in `BIOMETRICS/MEETING.md` protokollieren.

## Server-Register

| Server | Zweck | Installation | Aktivierung | Primäre Use Cases | Sicherheitsniveau | Status |
|---|---|---|---|---|---|---|
| Serena MCP | Codekontext, Projektaktivierung | {SERENA_INSTALL} | {SERENA_ACTIVATE} | Refactoring, Strukturarbeit | hoch | active |
| CDP MCP | Browser/Live-UI Analyse | {CDP_INSTALL} | {CDP_ACTIVATE} | UI-Checks, Live-Testflows | hoch | active |
| Tavily MCP | Recherche/Knowledge | {TAVILY_INSTALL} | {TAVILY_ACTIVATE} | Wissensrecherche, Vergleich | mittel | active |
| **Qwen 3.5 MCP** | **Bildanalyse, Code-Gen, OCR, Video** | **NVIDIA NIM Integration** | **opencode.json config** | **Vision, Code, Document, Video** | **hoch** | **active** |

## Qwen 3.5 MCP Server (NVIDIA NIM)

### Server-Übersicht
Qwen 3.5 (397B) via NVIDIA NIM für spezialisierte KI-Aufgaben.

### Konfiguration

**opencode.json Integration:**
```json
{
  "mcp": {
    "qwen-3.5": {
      "type": "local",
      "command": ["node", "/path/to/qwen-mcp-wrapper.js"],
      "enabled": true,
      "environment": {
        "NVIDIA_API_KEY": "${NVIDIA_API_KEY}",
        "MODEL_ID": "qwen/qwen3.5-397b-a17b"
      }
    }
  }
}
```

### Verfügbare Tools

| Tool | Beschreibung | Input | Output |
|---|---|---|---|
| `qwen_vision_analysis` | Bildanalyse und visuelle Erkennung | Bilder (PNG, JPG, WebP) | Strukturierte Analyse mit Tags |
| `qwen_code_generation` | Full-Stack Code-Generierung | Natürliche Sprache | Fertiger Code |
| `qwen_document_ocr` | Texterkennung aus Dokumenten | PDF, Bilder | Extrahierter Text |
| `qwen_video_understanding` | Video-Inhaltsanalyse | Videos (MP4, MOV) | Szenenbeschreibung |
| `qwen_conversation` | Natürliche Konversation | Benutzer-Nachrichten | Kontextbezogene Antworten |

### Installation

1. NVIDIA API Key besorgen: https://build.nvidia.com/
2. Environment Variable setzen:
   ```bash
   export NVIDIA_API_KEY="nvapi-xxx"
   ```
3. Wrapper installieren (falls vorhanden):
   ```bash
   npm install -g qwen-mcp-wrapper
   ```
4. opencode.json mit Qwen MCP config aktualisieren

### Aktivierungsbefehl

```bash
# Health Check
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
     https://integrate.api.nvidia.com/v1/models

# Server aktivieren (via opencode.json)
# Nach Änderung: opencode restart
```

### Deaktivierungsbefehl

```bash
# MCP in opencode.json deaktivieren
# enabled: false setzen
# opencode restart
```

## Installations-Standard
Für jeden Server dokumentieren:
1. Voraussetzungen
2. Installationsschritte
3. Konfigurationsparameter
4. Verbindungsprüfung
5. Rollback bei Fehlinstallation

## Aktivierungs-Standard
Für jeden Server dokumentieren:
1. Aktivierungsbefehl
2. Auth- und Rechtekontext
3. Health-Check
4. Deaktivierungsbefehl

## Nutzungsmatrix

| Aufgabe | Primärer MCP | Optionaler MCP | Warum |
|---|---|---|---|
| Code-Änderungen im Repo | Serena MCP | **Qwen 3.5 MCP** | Struktur + Code-Generierung |
| Browser-/UI-Probleme | CDP MCP | Serena MCP | reproduzierbarer UI-Test |
| Externe Recherche | Tavily MCP | Serena MCP | Quellenorientierung |
| Bildanalyse/OCR | **Qwen 3.5 MCP** | Tavily MCP | Vision-Analyse |
| Video-Verarbeitung | **Qwen 3.5 MCP** | - | Video-Understanding |
| NLM Asset Governance | Serena MCP | CDP MCP | Doku + Integrationscheck |

## Sicherheitsregeln
1. Keine Secrets in Prompt- oder Doku-Texten.
2. Keine unkontrollierten Schreibzugriffe.
3. Jede MCP-gestützte Änderung braucht Evidenz.
4. Bei P0-Risiko sofort eskalieren.
5. NVIDIA API Key NIEMALS in Code committen.

## Troubleshooting

| Problem | Ursache (vermutet) | Diagnose | Maßnahme | Eskalation |
|---|---|---|---|---|
| Server nicht erreichbar | Fehlkonfiguration | Health-Check | Konfiguration prüfen | P1 |
| Auth schlägt fehl | Token/Rechte | Auth-Log prüfen | Rechte korrigieren | P0 |
| Unklare Ergebnisse | falscher Scope | Scope-Review | Auftrag schärfen | P1 |
| **Qwen Rate Limit** | **40 RPM Limit** | **Response 429** | **60s warten + Fallback** | **P1** |
| **Qwen Latenz hoch** | **397B Modell** | **Response Time** | **Timeout auf 120s setzen** | **P1** |

### Qwen 3.5 Specific Troubleshooting

**Rate Limiting (HTTP 429):**
- Maximale Requests: 40 pro Minute
- Lösung: Fallback-Chain nutzen oder 60 Sekunden warten

**Hohe Latenz:**
- Qwen 3.5 397B: 70-90 Sekunden
- Lösung: Timeout auf 120000ms setzen

**Modell-Verfügbarkeit:**
- Korrekte ID: `qwen/qwen3.5-397b-a17b`
- NICHT: `qwen2.5` (falsches Modell)

## Änderungslog

| Datum | Änderung | Betroffene Server | Owner |
|---|---|---|---|
| {DATE} | Initialer Universal-Stand | Serena/CDP/Tavily | {OWNER_ROLE} |
| 2026-02-18 | Qwen 3.5 MCP Integration hinzugefügt | Qwen 3.5 MCP | AI Orchestrator |

## Abnahme-Check MCP
1. Server-Register vorhanden
2. Installation/Aktivierung dokumentiert
3. Nutzungsmatrix vorhanden
4. Sicherheitsregeln enthalten
5. Troubleshooting vorhanden
6. Qwen 3.5 Integration dokumentiert

---
