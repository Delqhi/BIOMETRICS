# INFRASTRUCTURE.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Infrastruktur folgt Port Sovereignty, Least Privilege und Drift-Kontrolle.
- Änderungen erfordern Runbook-, Backup- und Recovery-Nachweise.
- Betriebsrisiken sind aktiv im Risk Register zu führen.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universelle Betriebsvorlage für Infrastruktur, Laufzeit, Netzwerk und Observability.

---

## 🧠 NLM CLI COMMANDS

```bash
# Create notebook
nlm notebook create "Title"

# List sources
nlm source list <notebook-id>

# Delete old source (BEFORE adding new!)
nlm source delete <source-id> -y

# Add new source
nlm source add <notebook-id> --file "file.md" --wait
```

**⚠️ DUPLICATE PREVENTION:** ALWAYS run `nlm source list` before `nlm source add`!

---

## 🔄 DEQLHI-LOOP (INFINITE WORK MODE)

- After each completed task → Add 5 new tasks immediately
- Never "done" - only "next task"
- Always document → Every change in files
- Git commit + push after EVERY change
- Parallel execution ALWAYS (run_in_background=true)

### Loop Mechanism:
1. Task N Complete
2. Git Commit + Push
3. Update Docs
4. Add 5 New Tasks
5. Next Task N+1
6. Repeat infinitely

## Scope
- Laufzeitumgebungen: local, staging, production
- Compute, Netzwerk, Storage
- Monitoring, Backup/Restore

## Prinzipien
1. Reproduzierbar
2. Sicher
3. Beobachtbar
4. Skalierbar
5. Rückbaubar

## Umgebungsmodell (Template)

| Umgebung | Zweck | Verfügbarkeit | Zugang | Owner |
|---|---|---|---|---|
| local | Entwicklung | n/a | dev-only | {OWNER_ROLE} |
| staging | Vorabvalidierung | mittel | restricted | {OWNER_ROLE} |
| production | Livebetrieb | hoch | strict | {OWNER_ROLE} |

## Compute-Modell
- Frontend Runtime: {FRONTEND_RUNTIME}
- Backend Runtime: {BACKEND_RUNTIME}
- Worker Runtime: {WORKER_RUNTIME}

## Netzwerkmodell
- Eingangswege: {INGRESS_PATHS}
- Ausgehende Integrationen: {EGRESS_TARGETS}
- Segmentierung: {NETWORK_SEGMENTATION}

## Storage-Modell
- Datenhaltung: Supabase
- Asset Storage: {ASSET_STORAGE}
- Audit Storage: {AUDIT_STORAGE}

## Observability
- Logs: {LOG_SINKS}
- Metriken: {METRIC_SET}
- Alerts: {ALERT_RULES}
- Dashboards: {DASHBOARD_SET}

## Backup/Restore
- Backup-Intervall: {BACKUP_INTERVAL}
- RPO: {RPO}
- RTO: {RTO}
- Restore-Testintervall: {RESTORE_TEST_INTERVAL}

## Skalierungsstrategie
- horizontale Skalierung für API/Worker
- Lastprofile definieren
- Engpassanalyse pro Zyklus

## Sicherheitsanker
- Secret-Handling via ENV/Secret-Store
- Least Privilege Zugriffsmodell
- Netzwerkhärtung

## Verifikation
- Infrastruktur-Checkliste
- Restore-Test
- Alert-Simulation

## Abnahme-Check INFRASTRUCTURE
1. Umgebungsmodell vollständig
2. Observability definiert
3. Backup/Restore dokumentiert
4. Sicherheitsanker enthalten
5. Verifikation definiert

---
