# CONTRIBUTING.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Contributions müssen Rule-/Control-konform und auditierbar sein.
- Dokumentations- und Mapping-Updates sind verpflichtender Teil jedes Beitrags.
- Sicherheits- und Qualitätsverletzungen blockieren Merge.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Beitragspfad für konsistente, hochwertige Änderungen.

## Grundregeln
1. Erst lesen, dann ändern
2. kleine, prüfbare Änderungen
3. keine Done-Behauptung ohne Evidenz
4. Doku immer mitführen

## Beitragsablauf
1. Task auswählen
2. Scope klären
3. relevante Dateien lesen
4. Änderung umsetzen
5. Tests/Checks ausführen
6. Doku aktualisieren
7. PR mit Nachweisen erstellen

## Pflicht-Checks
- `pnpm lint`
- `pnpm typecheck`
- `pnpm test`
- `pnpm build`
- `go test ./...` (falls zutreffend)

## Dokumentationspflicht
Bei relevanten Änderungen aktualisieren:
- `CHANGELOG.md`
- `MEETING.md`
- betroffene Fachdoku

## NLM-Content-Beiträge
- ausschließlich via NLM-CLI
- passende Vorlage nutzen
- Qualitätsmatrix anwenden
- Delegation protokollieren

## PR-Template (Kurz)
1. Was geändert wurde
2. Warum
3. Prüfungen
4. Risiken
5. Nächste Schritte

## Abnahme-Check CONTRIBUTING
1. Ablauf klar und reproduzierbar
2. Pflicht-Checks genannt
3. Doku- und NLM-Regeln enthalten

---
