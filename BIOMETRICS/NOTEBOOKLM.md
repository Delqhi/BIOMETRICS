# NOTEBOOKLM.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- NLM-First ist verbindlich und vor Implementierung zu berücksichtigen.
- Quellen-/Asset-Hygiene folgt den globalen Duplicate- und Sync-Regeln.
- Jeder Asset-Lebenszyklus muss reproduzierbar dokumentiert sein.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Betriebsleitfaden für NotebookLM-Nutzung über NLM-CLI in Projekten.

## Pflichtsatz
NotebookLM wird vollständig über NLM-CLI genutzt. Alle Outputs werden mit Qualitätsmatrix bewertet und nur verifiziert übernommen.

## Universalitäts-Regeln
1. Keine projektspezifischen Secrets in dieser Datei.
2. IDs, Quellen und Artefakte nur mit Platzhaltern dokumentieren.
3. Jeder NLM-Output braucht Zweck, Quelle, Score und Freigabe-Status.

## Platzhalter
- {NLM_WORKSPACE_ID}
- {NLM_NOTEBOOK_ID}
- {ASSET_ID}
- {SOURCE_REF}
- {OWNER_ROLE}

## NLM-CLI Prozess
1. Aufgabe klassifizieren: Video | Infografik | Präsentation | Datentabelle | Bericht | Mindmap | Podcast
2. Passende Vorlage aus `../∞Best∞Practices∞Loop.md` auswählen
3. Quellenrahmen definieren
4. NLM-Generierung durchführen
5. Qualitätsmatrix anwenden
6. Freigabe/Verwerfung dokumentieren
7. Freigegebenes Artefakt in `BIOMETRICS/NLM-ASSETS/...` ablegen und referenzieren

## Quellenregister (Template)

| Source-ID | Typ | Beschreibung | Gültig ab | Owner |
|---|---|---|---|---|
| SRC-001 | doc | {SOURCE_REF} | {DATE} | {OWNER_ROLE} |

## Notebook-Register (Template)

| Notebook-ID | Zweck | Scope | Status | Owner |
|---|---|---|---|---|
| {NLM_NOTEBOOK_ID} | {PURPOSE} | {SCOPE} | ACTIVE | {OWNER_ROLE} |

## Asset-Register (Template)

| Asset-ID | Typ | Thema | Qualitäts-Score | Status | Verwendungsort |
|---|---|---|---|---|---|
| {ASSET_ID} | video | {TOPIC} | {SCORE} | review | {TARGET_DOC} |

## Asset-Typen (Pflichtkatalog)
1. video
2. infographic
3. presentation
4. table
5. report
6. mindmap
7. podcast

## Ablagepfade (kanonisch)
- `BIOMETRICS/NLM-ASSETS/videos/`
- `BIOMETRICS/NLM-ASSETS/infographics/`
- `BIOMETRICS/NLM-ASSETS/presentations/`
- `BIOMETRICS/NLM-ASSETS/tables/`
- `BIOMETRICS/NLM-ASSETS/reports/`
- `BIOMETRICS/NLM-ASSETS/mindmaps/`
- `BIOMETRICS/NLM-ASSETS/podcasts/`

## README-Einbindungspflicht
Für jedes freigegebene NLM-Artefakt:
1. Kurzbeschreibung in `../README.md`
2. Verweis auf Asset-Pfad
3. Zuordnung zu Zielseite/Zielflow

## Qualitätsmatrix
- Korrektheit
- Konsistenz
- Verständlichkeit
- Zielgruppenfit
- Umsetzbarkeit
- Wiederverwendbarkeit
- Evidenzbezug

Freigabe:
- Mindestens 13/16
- Korrektheit 2/2

## Delegationsprotokoll (Template)
```text
Delegation-ID:
Datum:
Asset-Typ:
Verwendete Vorlage:
Quellen:
Output-Score:
Entscheidung: übernommen | verworfen
Begründung:
Nächste Iteration:
Owner:
```

## Verbotene Zustände
1. Übernahme ungeprüfter NLM-Ausgaben
2. Fehlende Quellenbezüge
3. Nicht dokumentierte Verwerfungen

## Abnahme-Check NOTEBOOKLM
1. NLM-CLI Prozess dokumentiert
2. Register für Quellen/Notebooks/Assets vorhanden
3. Qualitätsmatrix und Freigaberegel enthalten
4. Delegationsprotokoll enthalten
5. Mindmap- und Podcast-Regeln enthalten

---
