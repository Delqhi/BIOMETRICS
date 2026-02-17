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

## Domain-Inventar (Template)

| Domain | Zweck | Umgebung | Status | Owner |
|---|---|---|---|---|
| {DOMAIN_1} | {PURPOSE_1} | production | planned | {OWNER_ROLE} |

## DNS-Standard
- A/AAAA Records nur wenn erforderlich
- CNAME bevorzugen, wo sinnvoll
- TTL bewusst definieren

## Zertifikats- und TLS-Rahmen
- Zertifikatstyp: {CERT_TYPE}
- Renewal-Prozess: {RENEWAL_PROCESS}
- Verantwortliche Rolle: {OWNER_ROLE}

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
