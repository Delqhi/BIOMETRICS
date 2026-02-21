# COST-OPTIMIZATION.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Kosten-Nutzen-Analyse ist Pflicht vor jeder Ressourcen-Allocierung.
- Ungenutzte Ressourcen werden innerhalb von 7 Tagen identifiziert und abgebaut.
- Free-First-Philosophie: Bevorzugt kostenlose oder selbstgehostete Lösungen.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

## Zweck
Strategien und Maßnahmen zur Kostenoptimierung für BIOMETRICS-Infrastruktur.

## 1) Kostenübersicht (Monatlich)

| Service | Anbieter | Geschätzte Kosten | Free Tier | Optimierungspotenzial |
|---------|----------|-------------------|-----------|----------------------|
| Supabase | Supabase | ~€50 | 500MB / 2GB | Mittel |
| Vercel | Vercel | ~€20 | 100GB Bandwidth | Mittel |
| n8n | Self-Hosted | €0 | ∞ | Keines |
| GitLab | Self-Hosted | €0 | ∞ | Keines |
| Cloudflare | Cloudflare | ~€5 | 100% | Niedrig |
| IONOS | IONOS | ~€15 | - | Niedrig |
| NVIDIA NIM | NVIDIA | ~€50 | 40 RPM | Hoch |
| **Gesamt** | | **~€140** | | |

## 2) Kostenoptimierungsmaßnahmen

### 2.1) Supabase
- **Database:**
  - Nutzung von Connection Pooling aktivieren
  - Ungenutzte Tabellen identifizieren und archivieren
  - Indexe auf häufige Queries optimieren
- **Storage:**
  - Objekte über 90 Tage in Glacier verschieben
  - Automatische Lifecycle-Policies aktivieren
- **API:**
  - Cache-Strategien für statische Daten
  - Query-Optimierung reduziert Compute-Nutzung

### 2.2) Vercel
- **Bandwidth:**
  - Statische Assets komprimieren (Brotli)
  - Images über next/image optimieren
  - Caching-Headers korrekt setzen
- **Serverless:**
  - Funktionen unter 10MB keep cold
  - Keine unnötigen Dependencies
  - Lambda-Laufzeit unter 30 Sekunden

### 2.3) NVIDIA NIM (Hoch priorisiert)
- **Caching:**
  - Modell-Caching für wiederkehrende Requests
  - Batch-Verarbeitung nutzen
- **Fallbacks:**
  - Kostenlose Modelle als Fallback (Mistral, Groq)
  - Nur Premium bei Bedarf
- **Monitoring:**
  - Token-Nutzung pro Request tracken
  - Budget-Alerts bei 80%

### 2.4) IONOS
- **Server:**
  - Nur notwendige Ressourcen allokieren
  - Auto-Scaling bei Bedarf aktivieren
- **Backup:**
  - Inkrementelle statt vollständige Backups
  - Unnötige Snapshots löschen

## 3) Kostenüberwachung

### 3.1) Dashboard
| Metrik | Quelle | Update-Frequenz |
|--------|--------|-----------------|
| Supabase Compute | Supabase Dashboard | Täglich |
| Vercel Bandwidth | Vercel Analytics | Echtzeit |
| NVIDIA API | NVIDIA Console | Täglich |
| IONOS | IONOS Panel | Täglich |

### 3.2) Alerts
| Alert | Threshold | Eskalation |
|-------|-----------|------------|
| Budget 80% erreicht | €112/Monat | Email |
| Ungewöhnlicher Anstieg | >50% vs Vormonat | Slack |
| Supabase Storage > 80% | 1.6GB | Slack |

## 4) Automatisierung

### 4.1) Cost-Saving Rules
```
- Storage Lifecycle: Objekte >90 Tage → Glacier
- Compute Idle: Keine Nutzung >1 Stunde → Pause
- Snapshot Cleanup: Älter als 30 Tage → Löschen
```

### 4.2) Reporting
- Wöchentlicher Cost-Report via n8n
- Monatliche Review im Team
- Quartalsweise Strategie-Anpassung

## 5) Einsparpotenzial

| Maßnahme | Einsparung/Jahr | Aufwand |
|----------|-----------------|---------|
| NVIDIA Fallback auf free Tier | ~€400 | Niedrig |
| Supabase Storage Lifecycle | ~€50 | Niedrig |
| Vercel Caching optimieren | ~€30 | Mittel |
| Ungenutzte Ressourcen abbauen | ~€100 | Niedrig |
| **Gesamtpotenzial** | **~€580/Jahr** | |

## 6) Verantwortlichkeiten

| Aufgabe | Verantwortlich | Frequenz |
|---------|----------------|----------|
| Cost-Monitoring | DevOps | Täglich |
| Report-Erstellung | Product Owner | Wöchentlich |
| Strategie-Anpassung | CTO | Quartalsweise |

---

## Abnahme-Check COST-OPTIMIZATION
1. Kostenübersicht dokumentiert und aktuell
2. Optimierungsmaßnahmen identifiziert
3. Monitoring und Alerts konfiguriert
4. Reporting-Schedule definiert
5. Verantwortlichkeiten geklärt

---
