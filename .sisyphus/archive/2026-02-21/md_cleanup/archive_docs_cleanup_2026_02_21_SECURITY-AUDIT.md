# SECURITY-AUDIT.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Security-Audits sind Pflicht und werden halbjährlich durchgeführt.
- Findings werden priorisiert und mit Tracking-IDs versehen.
- Re-Audits sind verbindlich nach Fix.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

## Zweck
Dokumentation von Security-Audits, Findings und deren Behebung für BIOMETRICS.

## 1) Audit-Frequenz und Scope

| Audit-Typ | Frequenz | Scope | Verantwortlich |
|-----------|----------|-------|---------------|
| Intern | Monatlich | Code + Config | DevOps |
| Externe Abhängigkeiten | Wöchentlich | npm/dependencies | Snyk/Dependabot |
| Penetrationstest | Halbjährlich | Full Stack | Externer Anbieter |
| Compliance | Jährlich | GxP, DSGVO | Compliance Officer |

## 2) Letzter Audit

**Datum:** {DATUM}  
**Typ:** Intern  
**Auditor:** {NAME}  
**Ergebnis:** {PASS/FAIL}

## 3) Findings-Register

### Offene Findings (P0-P2)

| ID | Severity | Kategorie | Finding | Status | Due Date |
|----|----------|-----------|---------|--------|----------|
| SA-001 | P0 | Auth | MFA nicht erzwungen | OPEN | {DATUM} |
| SA-002 | P1 | Secrets | API-Keys in Logs | IN_PROGRESS | {DATUM} |
| SA-003 | P2 | Config | CORS zu permissiv | OPEN | {DATUM} |

### Geschlossene Findings

| ID | Severity | Kategorie | Finding | Closed | Verified |
|----|----------|-----------|---------|--------|----------|
| SA-000 | P0 | Secrets | Env in Git | {DATUM} | {DATUM} |

## 4) Audit-Checkliste (Intern)

### 4.1) Authentication & Authorization
- [ ] MFA für alle Admin-Konten aktiviert
- [ ] Session-Timeout ≤ 30 Minuten
- [ ] Passwort-Policy durchgesetzt (12+ Zeichen, komplex)
- [ ] Role-Based Access Control (RBAC) implementiert

### 4.2) Data Protection
- [ ] Encryption at Rest (Supabase)
- [ ] Encryption in Transit (TLS 1.3)
- [ ] Sensitive Data in ENV nicht in Logs
- [ ] PII-Felder verschlüsselt

### 4.3) API Security
- [ ] Rate Limiting aktiviert
- [ ] Input Validation auf allen Endpunkten
- [ ] SQL-Injection Schutz (Parameterized Queries)
- [ ] XXS-Schutz (Content Security Policy)

### 4.4) Infrastructure
- [ ] Firewall-Regeln korrekt
- [ ] Unused Ports geschlossen
- [ ] Docker-Images aktuell
- [ ] Secrets in Vault, nicht in Code

### 4.5) Monitoring
- [ ] Logging aktiviert für sicherheitsrelevante Events
- [ ] Alerts für Anomalien konfiguriert
- [ ] Log-Retention ≥ 90 Tage

## 5) Vulnerability Scanning

### 5.1) Dependencies (Snyk/Dependabot)
```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: "npm"
    schedule: "weekly"
    open-pull-requests-limit: 10
```

### 5.2) Container Scanning
```bash
# Trivy Scan (wöchentlich)
trivy image --severity HIGH,CRITICAL biometrics-app:latest
```

### 5.3) Secret Scanning
```yaml
# .github/workflows/secret-scan.yml
name: Secret Scanning
on: [push, pull_request]
jobs:
  secret-scanning:
    runs-on: ubuntu-latest
    steps:
      - uses: trufflesecurity/trufflehog@main
        with:
          base: ${{ github.event.repository.default_branch }}
          head: ${{ github.head_ref }}
```

## 6) Compliance-Checks

| Standard | Status | Letzter Check | Nächster Check |
|----------|--------|---------------|----------------|
| DSGVO | PARTIAL | {DATUM} | {DATUM} |
| OWASP Top 10 | COMPLIANT | {DATUM} | {DATUM} |

## 7) Reporting-Template

```markdown
# Security Audit Report

## Executive Summary
- Audit-Datum: {DATUM}
- Auditor: {NAME}
- Gesamt-Risiko: {LOW/MEDIUM/HIGH/CRITICAL}

## Findings
[P0-P2 mit Details]

## Recommendations
[Priorisierte Maßnahmen]

## Timeline
[Fix-Deadlines]
```

---

## Abnahme-Check SECURITY-AUDIT
1. Audit-Schedule definiert
2. Findings-Register gepflegt
3. Checklisten vollständig
4. Vulnerability Scanning aktiviert
5. Compliance-Status dokumentiert

---
