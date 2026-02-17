# CI-CD-SETUP.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Pipeline-Gates müssen Policy-as-Code und Doku-Sync erzwingen.
- Sicherheits-, Mapping- und Qualitätschecks sind Default-Blocker bei Verstößen.
- Deployment-Freigaben erfolgen nur mit Nachweis kompletter Kontrollkette.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Pipeline-Standard für Qualitätssicherung und kontrollierte Auslieferung.

## Pipeline-Prinzipien
1. Jeder Merge ist überprüfbar
2. Qualität vor Geschwindigkeit
3. Rollback muss jederzeit möglich sein
4. Artefakte sind nachvollziehbar versioniert

## Standard-Stages
1. Install
2. Lint
3. Typecheck
4. Tests
5. Build
6. Security Checks
7. Packaging
8. Deploy
9. Post-Deploy Verification

## Gate-Definitionen
- Gate 1: `pnpm lint` grün
- Gate 2: `pnpm typecheck` grün
- Gate 3: `pnpm test` grün
- Gate 4: `pnpm build` grün
- Gate 5: `go test ./...` grün (falls Backend vorhanden)
- Gate 6: Security-Check ohne kritische Findings

## Branch-/Merge-Regeln
- PR Pflicht für geschützte Branches
- Keine direkte Änderung auf protected branch
- Merge nur bei grünen Gates

## Deploy-Strategie (Template)
- Staging: automatisch nach Merge in Integrationsbranch
- Production: kontrollierter Trigger mit Freigabe
- Rollback: definierter Fallback auf letzte stabile Version

## Artefakt-Management
- Build-Artefakte versionieren
- Prüfsummen und Metadaten speichern
- Release-Referenz im Changelog dokumentieren

## Post-Deploy Checks
1. Health Endpoints
2. Kernjourneys
3. Fehlerquote
4. Performance-Basiswerte

## Incident-Integration
- Bei P0/P1 Fehlern automatischer Eskalationspfad
- Postmortem-Pflicht mit Follow-up Tasks

## Abnahme-Check CI-CD
1. Stages und Gates vollständig
2. Merge-Regeln dokumentiert
3. Rollbackprozess vorhanden
4. Post-Deploy Checks definiert

---
