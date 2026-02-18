# SECURITY.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Security-by-Default und Secret-Hygiene sind nicht verhandelbar.
- Controls, Findings und Re-Tests folgen dem globalen Audit-Modell.
- Jede Ausnahme benötigt Compensating Control und Ablaufdatum.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Sicherheitsrahmen für Entwicklung, Betrieb und Content-Workflows.

## 1) Sicherheitsprinzipien
1. Least Privilege
2. Zero Trust Denkweise
3. Defense in Depth
4. Secure by Default
5. Nachweisbarkeit und Auditierbarkeit

## 2) Schutzobjekte
- Quellcode und Konfiguration
- Benutzer- und Betriebsdaten
- Integrationszugänge
- NLM-generierte Artefakte
- Betriebs- und Audit-Protokolle

## 3) Bedrohungsmodell (Template)

| Bedrohung | Angriffsvektor | Auswirkung | Gegenmaßnahme | Priorität |
|---|---|---|---|---|
| unautorisierter Zugriff | schwache Rechtevergabe | Datenabfluss | Rollenmodell + MFA | P0 |
| Secret Leak | Fehlkonfiguration | Account kompromittiert | Secret Hygiene + Rotation | P0 |
| Prompt Injection | ungeprüfte Inputs | fehlerhafte Artefakte | Input-Validierung + Review | P1 |
| Overclaim in Content | ungeprüfte NLM-Ausgabe | Compliance-Risiko | Qualitätsmatrix + Freigabe | P1 |

## 4) Zugriffskontrolle
- Rollen: user, dev, admin, agent
- Trennung von Leserechten und Schreibrechten
- Kritische Aktionen nur mit hoher Rolle

## 5) Secret-Management
1. Keine Secrets in Repo-Dateien
2. ENV-basiertes Secret Handling
3. Rotation bei Verdacht oder Incident
4. Zugriff nur für notwendige Rollen

## 5.1) NVIDIA API Key Management
Die NVIDIA NIM API ermöglicht Zugang zu High-Performance KI-Modellen wie Qwen 3.5 397B. Die folgenden Standards sind für alle NVIDIA API Keys verbindlich.

### 5.1.1) Key Rotation Schedule
- **Automatische Rotation:** Alle 90 Tage
- **Manuelle Rotation:** Bei Verdacht auf Kompromittierung sofort
- **Ablaufüberwachung:** Wöchentlicher Check 14 Tage vor Ablauf
- **Dokumentation:** Jede Rotation in Audit-Log mit Zeitstempel

### 5.1.2) Environment Variable Best Practices
```bash
# NVIDIA API Key - NIEMALS in Code oder Config-Files
export NVIDIA_API_KEY="nvapi-xxxxxxxxxxxx"

# Prefer via Vault (empfohlen)
eval $(vault env nvidia-api)
```
- **Regel 1:** Niemals hardcodieren
- **Regel 2:** Nur über Environment-Variablen nutzen
- **Regel 3:** Nicht in .env-Dateien speichern (Ausnahme: lokal mit .gitignore)
- **Regel 4:** In Docker-Containern nur via --env-file oder Docker Secrets

### 5.1.3) Vault Integration
- **Primärer Speicher:** HashiCorp Vault (room-02-tresor-vault)
- **Zugriffspfad:** `secret/data/nvidia-api`
- **Authentifizierung:** Kubernetes Service Account oder IAM Role
- **Rotation-Automation:** Vault Agent Injector für automatische Rotation

### 5.1.4) API Key Storage
| Speicherort | Typ | Zugriff |
|-------------|-----|---------|
| HashiCorp Vault | Primär | vault CLI, Kubernetes |
| Kubernetes Secrets | Sekundär | Nur für Pods |
| .env.local | Lokal DEV | Nur Developer-Machine |
| CI/CD Secrets | Pipeline | GitHub Actions, n8n |

**GitHub Secrets Konfiguration:**
```yaml
# .github/workflows/ci.yml
env:
  NVIDIA_API_KEY: ${{ secrets.NVIDIA_API_KEY }}
```

### 5.1.5) Monitoring und Alerts
- **Rate-Limit-Warnung:** Bei >80% Nutzung
- **Kosten-Alert:** Bei >80% des monatlichen Budgets
- **Fehler-Alert:** Bei 429 (Rate Limited) Status
- **Rotation-Erinnerung:** 7 Tage vor geplanter Rotation

## 6) API-Sicherheit
- Authentifizierung für produktive Endpunkte
- Autorisierung je Rolle
- Input-Validierung und Fehlerhärtung
- Rate Limits definieren

## 7) NLM-spezifische Sicherheit
1. NLM nur über NLM-CLI ausführen
2. Keine sensiblen Rohdaten ungefiltert an NLM
3. Output vor Übernahme auf Fakten und Compliance prüfen
4. Verworfenes dokumentieren

## 8) Compliance-Hinweise (Template)
- Rechtsraum: {COMPLIANCE_SCOPE}
- Datenklassen: {DATA_CLASSES}
- Aufbewahrung: {RETENTION_POLICY}
- Löschung: {DELETION_POLICY}

## 9) Incident-Response Quickflow
1. Severity klassifizieren (P0/P1/P2)
2. Schaden begrenzen
3. Secrets rotieren falls nötig
4. Ursache analysieren
5. Fix validieren
6. Postmortem dokumentieren

## 10) Security-Checks pro Zyklus
- P0-Risiken offen? nein
- Kritische Endpunkte geschützt? ja
- Secret-Policy eingehalten? ja
- NLM-Output geprüft? ja

## 11) Nachweise
- Security-Checklisten
- Audit-Logeinträge
- Freigabeprotokolle

## Abnahme-Check SECURITY
1. Threat Model vorhanden
2. Secret-Policy klar
3. Rollenmodell dokumentiert
4. NLM-Sicherheitsregeln enthalten
5. Incident-Prozess definiert

---
