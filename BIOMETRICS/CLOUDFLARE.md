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

## Tunnel Services

| Service | Tunnel-Name | Subdomain | Ziel-Port | Status |
|---------|-------------|------------|-----------|--------|
| OpenCode | opencode-tunnel | opencode.delqhi.com | 18789 | active |
| OpenClaw | openclaw-tunnel | openclaw.delqhi.com | 18789 | active |
| n8n | n8n-tunnel | n8n.delqhi.com | 5678 | active |
| Supabase | supabase-tunnel | supabase.delqhi.com | 54321 | active |

### OpenCode Tunnel
```yaml
tunnel: opencode-tunnel
ingress:
  - hostname: opencode.delqhi.com
    service: http://localhost:18789
  - service: http_status:404
```

### OpenClaw Tunnel
```yaml
tunnel: openclaw-tunnel
ingress:
  - hostname: openclaw.delqhi.com
    service: http://localhost:18789
  - service: http_status:404
```

### n8n Tunnel
```yaml
tunnel: n8n-tunnel
ingress:
  - hostname: n8n.delqhi.com
    service: http://localhost:5678
  - service: http_status:404
```

### Supabase Tunnel
```yaml
tunnel: supabase-tunnel
ingress:
  - hostname: supabase.delqhi.com
    service: http://localhost:54321
  - service: http_status:404
```

## Abnahme-Check CLOUDFLARE
1. Konfigurationsmatrix vollständig
2. Sicherheitsregeln dokumentiert
3. Verifikationsschritte vorhanden

---
