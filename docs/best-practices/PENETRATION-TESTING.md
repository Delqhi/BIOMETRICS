# PENETRATION-TESTING.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Pentests werden halbjährlich von qualifizierten Personen durchgeführt.
- Kritische Findings (P0) werden innerhalb von 24h adressiert.
- Alle Exploits werden dokumentiert und mit Fix verifiziert.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

## Zweck
Definition von Penetrationstest-Prozessen, -Scope und -Ergebnissen für BIOMETRICS.

## 1) Pentest-Plan

### 1.1) Frequenz
| Typ | Frequenz | Anbieter | Kosten |
|-----|----------|----------|--------|
| Intern | Monatlich | DevOps Team | €0 |
| Extern | Halbjährlich | {ANBIETER} | ~€5000 |

### 1.2) Scope
- Frontend: Vercel gehostete Anwendung
- Backend: Supabase Edge Functions, n8n
- API: Alle öffentlichen Endpunkte
- Auth: Login, Session-Management
- Infrastructure: Cloudflare, IONOS

### 1.3) Ausschlüsse
- DoS-Tests (außer sanfte Rate-Limit-Tests)
- Social Engineering
- Phishing-Kampagnen
- Physicher Zugang

## 2) OWASP Top 10 (2021)

| # | Kategorie | Test-Methode | Letzter Test | Status |
|---|-----------|--------------|--------------|--------|
| A01 | Broken Access Control | Manuelle Prüfung + autom. Scanner | {DATUM} | PASS |
| A02 | Cryptographic Failures | Code Review | {DATUM} | PASS |
| A03 | Injection | SQL/Command Injection Tests | {DATUM} | PASS |
| A04 | Insecure Design | Architektur-Review | {DATUM} | PASS |
| A05 | Security Misconfig | Config-Analyse | {DATUM} | PASS |
| A06 | Vulnerable Components | Dependency Scan | {DATUM} | PASS |
| A07 | Auth Failures | Brute-Force, Session Hijack | {DATUM} | PASS |
| A08 | Data Integrity Failures | API-Integrität | {DATUM} | PASS |
| A09 | Logging Failures | Log-Analyse | {DATUM} | PASS |
| A10 | SSRF | Internal Scan | {DATUM} | PASS |

## 3) Test-Methodik

### 3.1) Reconnaissance
```bash
# Subdomain Enumeration
subfinder -d biometrics.de -o subdomains.txt

# Port Scanning
nmap -sV -p- biometrics.de

# Technology Fingerprinting
wappalyzer https://biometrics.de
```

### 3.2) Vulnerability Assessment
```bash
# Automated Scan
nikto -h https://biometrics.de

# Nuclei Templates
nuclei -t cves/ -u https://biometrics.de

# SQL Injection
sqlmap -u "https://biometrics.de/api?id=1"
```

### 3.3) Manual Testing
- Business Logic Flaws
- Authorization Bypass
- Race Conditions
- JWT Token Analysis

## 4) Findings-Register

### Offene Findings

| ID | Severity | Kategorie | Description | Status | Due Date |
|----|----------|-----------|-------------|--------|----------|
| PT-001 | CRITICAL | SQLi | Parameterized Queries fehlen | IN_PROGRESS | {DATUM} |
| PT-002 | HIGH | XSS | CSP nicht gesetzt | OPEN | {DATUM} |
| PT-003 | MEDIUM | Info | Debug-Mode aktiv | OPEN | {DATUM} |

### Behobene Findings

| ID | Severity | Finding | Fixed | Re-Test |
|----|----------|---------|-------|---------|
| PT-000 | CRITICAL | Admin ohne MFA | {DATUM} | {DATUM} |

## 5) Exploit-Documentation (Template)

```markdown
## Exploit Report: PT-001

### Details
- **Finding:** SQL Injection in /api/users
- **Severity:** CRITICAL
- **CVSS:** 9.8

### Proof of Concept
```bash
curl "https://biometrics.de/api/users?id=1' OR '1'='1"
```

### Impact
- Full Database Access
- User Data Exfiltration
- Potential RCE

### Remediation
1. Parameterized Queries nutzen
2. Input Validation
3. WAF Rule

### Status
- [ ] Fix implementiert
- [ ] Re-Test bestanden
- [ ] Closed
```

## 6) Re-Test-Prozess

```
1. Fix wird implementiert
2. Re-Test Request erstellen
3. Gleiche Testszenarien durchführen
4. Bei Erfolg: Finding als "Verified Fixed" markieren
5. Bei Misserfolg: Finding bleibt offen mit neuer Due Date
```

## 7) Tools

| Tool | Lizenz | Use Case |
|------|--------|----------|
| Burp Suite | Pro | Full Penetration Test |
| SQLMap | Open Source | SQL Injection |
| Nuclei | Open Source | Vulnerability Scanning |
| subfinder | Open Source | Subdomain Enum |
| nikto | Open Source | Web Server Scan |

---

## Abnahme-Check PENETRATION-TESTING
1. Pentest-Schedule definiert
2. Scope dokumentiert
3. OWASP Top 10 Tests durchgeführt
4. Findings-Register gepflegt
5. Re-Test-Prozess definiert

---
