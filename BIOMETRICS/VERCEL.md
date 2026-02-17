# VERCEL.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Plattformbetrieb folgt globalen Deployment- und Security-Leitlinien.
- Release- und Rollback-Prozesse sind testbar und dokumentiert.
- Konfigurationsdrift wird regelmäßig geprüft und korrigiert.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Leitfaden für Vercel-Projektbetrieb, Environments und Deploy-Governance.

## Hinweis
Wenn Vercel nicht genutzt wird, Status auf `NOT_APPLICABLE` setzen und Begründung dokumentieren.

## Projekt-Metadaten (Template)
- Projektname: {PROJECT_NAME}
- Projekt-ID: {VERCEL_PROJECT_ID}
- Team: {VERCEL_TEAM}
- Owner: {OWNER_ROLE}

## Environment-Strategie
- development
- preview
- production

## Deploy-Regeln
1. Deploy nur über kontrollierten CI/CD-Pfad
2. Produktionsdeploy nur bei grünen Gates
3. Rollback-Option vor Deploy verifizieren

## Domain-/Routing-Template
- Primärdomain: {PRIMARY_DOMAIN}
- Zusätzliche Domains: {ADDITIONAL_DOMAINS}
- Redirect-Policy: {REDIRECT_POLICY}

## Sicherheitsregeln
- Secrets nur über Environment-Management
- Keine Secrets in Repo-Dateien
- Zugriff auf Production minimal halten

## Verifikation
- Preview Deploy erfolgreich
- Production Health-Check erfolgreich
- Rollback-Test dokumentiert

## Betriebscheckliste
### Pre-Deploy
1. Alle CI-Gates grün
2. Ziel-Environment bestätigt
3. Rollback-Referenz dokumentiert

### Deploy
1. Deploy-Trigger dokumentiert
2. Versionsreferenz notiert
3. Status-Checks beobachtet

### Post-Deploy
1. Health-Endpunkte geprüft
2. Kernjourney-Sanity-Check durchgeführt
3. Fehler-/Latenzbaseline geprüft

### Rollback
1. Triggerkriterium erreicht?
2. Rollback ausgelöst
3. Nach-Rollback Health-Check erfolgreich

## Abnahme-Check VERCEL
1. Projekt-/Environmentdaten als Platzhalter vorhanden
2. Deploy- und Rollback-Regeln dokumentiert
3. Sicherheitsregeln enthalten

---
