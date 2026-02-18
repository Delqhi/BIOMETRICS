# IONOS.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Hosting-Betrieb folgt globalen Infra-, Security- und Recovery-Vorgaben.
- DNS, Secrets und Deploy-Konfigurationen bleiben revisionssicher dokumentiert.
- Provider-spezifische Risiken werden proaktiv nachverfolgt.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Leitfaden für Domain- und DNS-Verwaltung bei IONOS.

## Hinweis
Wenn IONOS nicht genutzt wird, Status auf `NOT_APPLICABLE` setzen und Begründung dokumentieren.

## Domain-Inventar

| Domain | Zweck | Umgebung | Status | Owner |
|---|---|---|---|---|
| biometrics.de | Hauptdomain | production | active | Simone Schulze |
| api.biometrics.de | Backend API | production | active | DevOps |
| cdn.biometrics.de | Static Assets (Supabase) | production | active | DevOps |
| staging.biometrics.de | Staging-Umgebung | staging | active | DevOps |

## DNS-Konfiguration

| Record | Typ | Wert | TTL |
|---|---|---|---|
| biometrics.de | A | 212.227.XXX.XXX (IONOS Server) | 3600 |
| api.biometrics.de | CNAME | biometrics.de | 3600 |
| cdn.biometrics.de | CNAME | biometrics.supabase.co | 3600 |
| staging.biometrics.de | A | 172.20.0.XX (Docker Internal) | 300 |
| www.biometrics.de | CNAME | biometrics.de | 3600 |

## Zertifikats- und TLS-Rahmen
- Zertifikatstyp: Let's Encrypt (automatisch über IONOS)
- Renewal-Prozess: Automatisch 30 Tage vor Ablauf
- Verantwortliche Rolle: DevOps

## Sicherheitsregeln
1. Zugang mit MFA absichern
2. Rollen- und Rechteprüfung monatlich
3. Änderungen protokolliert in CHANGELOG.md

## Deploy-Konfiguration

### Build & Deploy Pipeline
```yaml
deploy:
  provider: IONOS
  server: 212.227.XXX.XXX
  user: root
  key_path: ~/.ssh/ionos_deploy
  webroot: /var/www/biometrics
  ssl:
    lets_encrypt: true
    auto_renew: true
```

### Automatischer Deploy
1. Git Push auf main
2. CI Pipeline läuft (Lint → Test → Build)
3. Bei Erfolg: Deploy auf staging.biometrics.de
4. Manuell: Prüfung und Deploy auf biometrics.de

## Sicherheitsregeln
1. Zugang mit MFA absichern
2. Rollen- und Rechteprüfung regelmäßig
3. Änderungen protokollieren

## Verifikation
- DNS-Auflösung pro Domain geprüft
- Zertifikat gültig
- Delegation korrekt

## Betriebscheckliste
### DNS
1. notwendige Records vollständig
2. TTL-Werte bewusst gesetzt
3. Konfliktfreie Einträge bestätigt

### TLS
1. Zertifikat aktiv und gültig
2. Erneuerungsprozess terminiert
3. Verantwortlichkeit dokumentiert

### Zugriff
1. MFA aktiv
2. Rollen und Rechte geprüft
3. Änderungen protokolliert

## Abnahme-Check IONOS
1. Domain-Inventar vorhanden
2. DNS- und TLS-Rahmen dokumentiert
3. Sicherheitsregeln enthalten

---
