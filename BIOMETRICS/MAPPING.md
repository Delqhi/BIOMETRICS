# MAPPING.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Mapping-Pflege ist mandatory und gilt als Teil der Definition of Done.
- Jede strukturelle Änderung muss hier referenziert und geprüft werden.
- Drift zwischen Domänen wird sofort als Risiko behandelt.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Zentrale Mapping-Governance für Konsistenz zwischen Produkt, Frontend, Backend, Datenbank, Commands, Endpoints und NLM-Assets.

## Pflichtregel
Jede relevante Änderung wird gegen die Mapping-Dateien geprüft. Ohne Mapping-Konsistenz ist eine Änderung nicht done.

## Mapping-Module
1. `MAPPING-COMMANDS-ENDPOINTS.md`
2. `MAPPING-FRONTEND-BACKEND.md`
3. `MAPPING-DB-API.md`
4. `MAPPING-NLM-ASSETS.md`

## Prüfzyklus
1. Änderung identifizieren
2. betroffenes Mapping-Modul aktualisieren
3. Konsistenzcheck durchführen
4. Ergebnisse in `MEETING.md` und `CHANGELOG.md` dokumentieren

## Abnahme-Check MAPPING
1. Alle Mapping-Module vorhanden
2. Zuordnungstabellen ausgefüllt
3. Offene Deltas explizit markiert
4. Prüfstatus dokumentiert

---
