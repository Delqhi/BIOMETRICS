# GITHUB.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Branch-, Review- und Schutzregeln folgen globalem Governance-Kern.
- PRs benötigen klare Nachweise für Tests, Doku und Mappings.
- Automationsregeln sind als Policy-as-Code zu bevorzugen.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universelle Repository-Governance für Branching, Reviews, Qualität und Nachvollziehbarkeit.

## Branching-Modell
- `main`: stabil
- `feature/*`: Funktionsarbeit
- `fix/*`: Fehlerbehebungen
- `chore/*`: Wartung

## Pull-Request-Regeln
1. klare Zielbeschreibung
2. Scope begrenzt
3. Tests/Checks angegeben
4. Risiken benannt
5. Doku-Änderungen enthalten

## Review-Gates
- mindestens ein qualifizierter Review
- alle CI-Gates grün
- keine offenen P0-Risiken

## Merge-Regeln
- squash/merge gemäß Teamregel
- kein Merge bei roten Checks
- kein direktes Pushen auf geschützte Branches

## Issue- und Task-Kopplung
- Änderungen referenzieren Task-ID
- Changelog und Meeting aktualisieren

## Sicherheitsaspekte
- Secret Scanning aktiv
- Dependabot/Abhängigkeitsupdates prüfen
- Branch Protection aktivieren

## Abnahme-Check GITHUB
1. Branch-/PR-/Merge-Regeln definiert
2. Review-Gates enthalten
3. Sicherheitsregeln enthalten
4. Task-Referenzierung geregelt

---
