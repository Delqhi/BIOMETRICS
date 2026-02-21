# ONBOARDING.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Onboarding vermittelt den vollständigen globalen Regelkern verpflichtend.
- NLM-First, No-Blind-Delete und Evidence-Disziplin sind Basiskompetenzen.
- Rollen und Eskalationspfade werden explizit trainiert.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Schneller, rollenbasierter Einstieg für User, Dev und Admin in ein universelles Projektsetup.

## 10-Minuten Schnellstart

### Schritt 1: Projekt verstehen (2 Min)
1. `../README.md` lesen - Übersicht und Vision
2. `../∞Best∞Practices∞Loop.md` lesen - Arbeitsmethodik
3. `AGENTS-PLAN.md` lesen - Aktuelle Tasks und Prioritäten

### Schritt 2: Rolle definieren (1 Min)
- **User:** Möchtest du das System nutzen → User-Pfad
- **Developer:** Möchtest du entwickeln → Dev-Pfad  
- **Admin:** Möchtest du verwalten → Admin-Pfad

### Schritt 3: Setup durchführen (5 Min)
```bash
# CLI installieren
cd ../biometrics-cli
pnpm install
pnpm link --global

# Onboarding starten (automatisches Setup)
biometrics
```

### Schritt 4: Erste Aufgabe ausführen (2 Min)
- User: `USER-PLAN.md` öffnen → P0 Task erledigen
- Dev: Aus `AGENTS-PLAN.md` einen Task übernehmen
- Admin: Governance-Dateien prüfen

## Rollenpfade

### User-Pfad
1. `USER-PLAN.md` öffnen
2. offene P0 User-Tasks erledigen
3. Verifikation dokumentieren
4. Fortschritt in `MEETING.md` eintragen

### Dev-Pfad
1. `AGENTS.md` lesen - Arbeitsregeln verstehen
2. Task aus `AGENTS-PLAN.md` übernehmen
3. Read-before-write: Alle relevanten Dateien lesen
4. Änderung durchführen + Tests + Doku-Update
5. Evidenz in `CHANGELOG.md` dokumentieren

### Admin-Pfad
1. Governance-Dateien prüfen (`AGENTS-GLOBAL.md`, `AGENTS.md`)
2. Security-/Compliance-Status prüfen (`SECURITY.md`)
3. Task-20 Abschlussreport freigeben
4. `MEETING.md` für Team-Koordination nutzen

## Betriebsbefehle (universell)
- `pnpm install`
- `pnpm lint`
- `pnpm typecheck`
- `pnpm test`
- `pnpm build`
- `pnpm dev`
- `go test ./...`
- `go vet ./...`

## NLM-CLI Onboarding
1. NLM-Aufgabenart bestimmen (Video/Infografik/Präsentation/Tabelle)
2. passende Vorlage aus `../∞Best∞Practices∞Loop.md` nutzen
3. Ergebnis gegen NLM-Qualitätsmatrix prüfen
4. nur verifizierte Inhalte übernehmen
5. Delegation in `MEETING.md` protokollieren

## Done-Definition pro Onboarding-Schritt
Ein Schritt ist erst done wenn:
1. Ziel erreicht
2. Verifikation durchgeführt
3. Doku aktualisiert

## Häufige Fehler
1. Direkt editieren ohne vorher zu lesen
2. NLM-Output ungeprüft übernehmen
3. Done ohne Evidenz melden
4. Doku-Updates vergessen

## Troubleshooting-Verweise
- Prozessfragen: `AGENTS.md`
- Taskfragen: `AGENTS-PLAN.md`
- User-Aufgaben: `USER-PLAN.md`
- Historie: `CHANGELOG.md`

## Abnahme-Check ONBOARDING
1. Alle 3 Rollen abgedeckt
2. Schnellstart klar
3. NLM-CLI Pfad enthalten
4. Verifikationslogik enthalten

---
