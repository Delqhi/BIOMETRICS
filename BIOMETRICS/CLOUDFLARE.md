# CLOUDFLARE.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Edge-/Security-Einstellungen folgen globalen Kontrollvorgaben.
- Änderungen an WAF, Routing und Caching werden evidenzbasiert dokumentiert.
- Drift und Fehlkonfigurationen sind als Governance-Risiko zu behandeln.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Leitfaden für Cloudflare-basierte Netzwerk- und Zugriffssicherung.

## Scope
- DNS
- Tunnel
- Access-Regeln
- Edge-Sicherheitsregeln

## Betriebsprinzipien
1. minimal notwendige Exposition
2. klare Trennung von öffentlich/privat
3. Logging und Nachvollziehbarkeit

## Konfigurationsmatrix (Template)

| Komponente | Zweck | Status | Owner |
|---|---|---|---|
| DNS-Zone | Namensauflösung | planned | {OWNER_ROLE} |
| Tunnel | sicherer Zugang | planned | {OWNER_ROLE} |
| Access Policy | Zugriffsschutz | planned | {OWNER_ROLE} |
| WAF Rules | Angriffsreduktion | planned | {OWNER_ROLE} |

## Sicherheitsregeln
- keine sensiblen Ursprungsdienste direkt exponieren
- Zugriff über definierte Policies absichern
- Änderungen versioniert dokumentieren

## Verifikation
- DNS-Auflösung geprüft
- Tunnel-Health geprüft
- Access-Regeln geprüft
- Baseline-WAF-Regeln aktiv

## Abnahme-Check CLOUDFLARE
1. Konfigurationsmatrix vollständig
2. Sicherheitsregeln dokumentiert
3. Verifikationsschritte vorhanden

---
