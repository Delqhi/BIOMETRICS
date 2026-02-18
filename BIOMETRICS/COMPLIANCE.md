# COMPLIANCE.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Compliance ist nicht verhandelbar - Bußgelder und Reputationsschäden drohen.
- Dokumentationspflichten werden strikt eingehalten.
- Externe Audits werden koordiniert und Ergebnisse archiviert.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

## Zweck
Übersicht über Compliance-Anforderungen, -Status und -Maßnahmen für BIOMETRICS.

## 1) Relevante Standards

| Standard | Anwendungsbereich | Status | Nächster Audit |
|----------|------------------|--------|-----------------|
| DSGVO | EU Datenschutz | PARTIAL | {DATUM} |
| ISO 27001 | Informationssicherheit | NOT_STARTED | - |
| SOC 2 | US Cloud Security | NOT_STARTED | - |

## 2) DSGVO-Compliance

### 2.1) Datenverarbeitung

| Prozess | Rechtsgrundlage | Datenkategorie | Aufbewahrung |
|---------|-----------------|----------------|--------------|
| User Registration | Einwilligung | Name, Email, IP | 2 Jahre Inaktivität |
| Analytics | Berechtigtes Interesse | Aggregiert | 14 Monate |
| Support Tickets | Vertragserfüllung | Kommunikation | 3 Jahre |
| Payment Data | Gesetzliche Pflicht | Finanzdaten | 10 Jahre |

### 2.2) Betroffenenrechte

| Recht | Implementierung | Status |
|-------|----------------|--------|
| Auskunft | /api/users/me Endpoint | ✅ |
| Löschung | /api/users/:id DELETE | ✅ |
| Berichtigung | /api/users/:id PATCH | ✅ |
| Datenportabilität | /api/users/:id/export | ✅ |
| Widerruf | /api/auth/revoke | ✅ |

### 2.3) Technische Maßnahmen

| Maßnahme | Implementierung | Geprüft |
|----------|-----------------|---------|
| Verschlüsselung at Rest | Supabase Encryption | ✅ |
| Verschlüsselung in Transit | TLS 1.3 | ✅ |
| Pseudonymisierung | UUID für User IDs | ✅ |
| Zwei-Faktor-Auth | TOTP (in Planung) | ❌ |
| Logging | Alle Zugriffe geloggt | ✅ |

### 2.4) Organisatorische Maßnahmen

| Maßnahme | Verantwortlich | Frequenz |
|----------|---------------|----------|
| Datenschutzbeauftragter | {NAME} | - |
| Privacy by Design Review | Product Team | Bei jedem Feature |
| DSFA | Compliance Officer | Jährlich |
| Mitarbeiter-Schulung | HR | Halbjährlich |

### 2.5) Dokumentation

| Dokument | Status | Letzte Aktualisierung |
|----------|--------|---------------------|
| Verarbeitungsverzeichnis | ✅ Aktuell | {DATUM} |
| DSFA | ✅ Aktuell | {DATUM} |
| Technische Dokumentation | ✅ Aktuell | {DATUM} |
| Incident Response Plan | ⚠️ Draft | {DATUM} |

## 3) ISO 27001 (Zukünftig)

### 3.1) Gap-Analyse

| Control | Status | Aufwand |
|---------|--------|---------|
| A.5 Information Security Policies | PARTIAL | Mittel |
| A.6 Organization of Information Security | PARTIAL | Niedrig |
| A.7 Human Resource Security | PARTIAL | Mittel |
| A.8 Asset Management | PARTIAL | Niedrig |
| A.9 Access Control | PARTIAL | Mittel |
| A.10 Cryptography | PARTIAL | Niedrig |

### 3.2) Roadmap

| Phase | Zeitraum | Scope |
|-------|----------|-------|
| Phase 1 | Q3 2026 | Gap-Analyse + Maßnahmenplan |
| Phase 2 | Q4 2026 | Implementierung |
| Phase 3 | Q1 2027 | Interne Audit |
| Phase 4 | Q2 2027 | Zertifizierung |

## 4) SOC 2 (Zukünftig)

### Trust Service Criteria

| Kriterium | Status | Maßnahmen |
|-----------|--------|-----------|
| Security | PARTIAL | MFA, Encryption, Logging |
| Availability | PARTIAL | Monitoring, Backups |
| Processing Integrity | PARTIAL | Validation, Testing |
| Confidentiality | PARTIAL | Encryption, Access Control |
| Privacy | PARTIAL | DSGVO + Policies |

## 5) Audit-Log-Anforderungen

### 5.1) Pflicht-Events

| Event | Retention | Format |
|-------|----------|--------|
| User Login | 2 Jahre | JSON |
| Data Access | 1 Jahr | JSON |
| Data Modification | 2 Jahre | JSON |
| Data Deletion | 2 Jahre | JSON |
| Admin Actions | 3 Jahre | JSON |
| Security Events | 3 Jahre | JSON |

### 5.2) Export

```bash
# Audit Log Export (Admin)
GET /api/admin/audit-log?from={DATE}&to={DATE}&format=json
```

## 6) Non-Compliance-Konsequenzen

| Verstoß | Bußgeld (max) | Reputationsschaden |
|---------|---------------|-------------------|
| DSGVO Art. 83 | €20M oder 4% Umsatz | Hoch |
| DSGVO Art. 83 | €10M oder 2% Umsatz | Mittel |
| Datenschutzverletzung | Meldepflicht 72h | Hoch |

---

## Abnahme-Check COMPLIANCE
1. Relevante Standards identifiziert
2. DSGVO-Status dokumentiert
3. Roadmap für ISO 27001/SOC 2 vorhanden
4. Audit-Log-Anforderungen definiert
5. Verantwortlichkeiten geklärt

---
