# DISASTER-RECOVERY.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Disaster Recovery ist kritische Infrastruktur - kein Opt-In.
- RTO/RPO-Ziele sind verbindlich und werden getestet.
- Dokumentation muss aktuell bleiben und regelmäßig geübt werden.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

## Zweck
Definierung von Disaster-Recovery-Strategien, -Prozessen und -Verantwortlichkeiten für BIOMETRICS.

## 1) RTO/RPO-Ziele

| Service | RTO (Recovery Time Objective) | RPO (Recovery Point Objective) | Priorität |
|---------|------------------------------|-------------------------------|-----------|
| Supabase DB | 15 Minuten | 5 Minuten | P0 |
| Vercel App | 5 Minuten | 0 (statisch) | P0 |
| n8n Workflows | 30 Minuten | 1 Stunde | P1 |
| GitLab Storage | 1 Stunde | 24 Stunden | P2 |
| Cloudflare DNS | 0 Minuten | 0 (multi) | P0 |

## 2) Backup-Strategie

### 2.1) Supabase Database
- **Automatische Backups:** Täglich um 02:00 UTC
- **Point-in-Time Recovery:** 7 Tage Aufbewahrung
- **Manual Snapshots:** Vor jedem Deployment
- **Cross-Region:** Backup-Storage in EU-West

### 2.2) Vercel
- **Git-Integration:** Automatisches Deployment bei Push
- **Environments:** Preview, Production
- **Rollback:** Ein-Klick Rollback möglich

### 2.3) n8n Workflows
- **Export:** Wöchentlich alle Workflows als JSON
- **Versionierung:** GitLab Storage
- **Backup-Schedule:** Sonntags 03:00 UTC

### 2.4) Konfigurationsdateien
- **Sekunden:** Environment-Variablen in Vault
- **Secrets:** HashiCorp Vault mit Versionskontrolle
- **Dokumentation:** Git-Versioniert

## 3) Recovery-Prozesse

### 3.1) Database Recovery (Supabase)
```
1. Identify Failure Scope
   └─> Check Supabase Dashboard Status

2. Initiate Recovery
   └─> Select Point-in-Time (max 7 days)
   └─> Restore to new instance

3. Validate
   └─> Run migration check script
   └─> Verify data integrity

4. Switchover
   └─> Update connection string if needed
   └─> Notify stakeholders
```

### 3.2) Application Recovery (Vercel)
```
1. Check Deployment Status
   └─> Vercel Dashboard / GitHub Actions

2. Rollback if Needed
   └─> GitHub: Previous successful deployment
   └─> Vercel: Click "Rollback"

3. Verify
   └─> Health check endpoint
   └─> Smoke tests
```

### 3.3) Complete System Failure
```
1. Declare Incident (P0)
2. Activate War Room
3. Follow Recovery Priority:
   a) Cloudflare DNS (always up)
   b) Vercel (statisch, schnell)
   c) Supabase (kritisch, RTO 15min)
   d) n8n (P1, RTO 30min)
4. Post-Incident Review
```

## 4) Monitoring & Alerts

| Alert-Type | Trigger | Eskalation |
|------------|---------|------------|
| Database Down | 1 min connection failure | P0: SMS + Call |
| API Error Rate > 5% | 5 min sustained | P1: Slack + Email |
| Backup Failure | Any failure | P1: Email |
| High Latency > 2s | 10 min sustained | P2: Slack |

## 5) Runbook-Inhalte (Kurzversion)

### 5.1) DB-Wiederherstellung
- Supabase Dashboard → Settings → Database → Point-in-Time Recovery
- Neuen Host erstellen, Connection String aktualisieren
- Health Check: `GET /api/health`

### 5.2) Vercel-Rollback
- GitHub Actions → Workflow Runs → Neuester erfolgreicher
- Oder: Vercel Dashboard → Deployments → Rollback

### 5.3) n8n-Recovery
- GitLab: Workflow-Export abrufen
- n8n: Import → JSON einfügen
- Aktivieren und testen

## 6) Testing-Schedule

| Test | Frequenz | Verantwortlich |
|------|----------|----------------|
| Backup-Verifizierung | Monatlich | DevOps |
| Restore-Drill | Quartalsweise | Team Lead |
| Failover-Simulation | Halbjährlich | CTO |

## 7) Kontakte & Verantwortlichkeiten

| Rolle | Name | Kontakt | Verantwortung |
|-------|------|---------|---------------|
| DR-Lead | {NAME} | {EMAIL} | Gesamtkoordination |
| DB-Admin | {NAME} | {EMAIL} | Supabase Recovery |
| App-Lead | {NAME} | {EMAIL} | Vercel Recovery |

---

## Abnahme-Check DISASTER-RECOVERY
1. RTO/RPO dokumentiert und realistisch
2. Backup-Schedule definiert und automatisiert
3. Recovery-Prozesse dokumentiert (Runbooks)
4. Monitoring-Alerts konfiguriert
5. Testing-Schedule festgelegt

---
