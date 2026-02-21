# GITHUB.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Branch-, Review- und Schutzregeln folgen globalem Governance-Kern.
- PRs benötigen klare Nachweise für Tests, Doku und Mappings.
- Automationsregeln sind als Policy-as-Code zu bevorzugen.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universelle Repository-Governance für Branching, Reviews, Qualität und Nachvollziehbarkeit.

## Branching-Modell
- `main`: stabil
- `feature/*`: Funktionsarbeit
- `fix/*`: Fehlerbehebungen
- `chore/*`: Wartung

## Pull-Request-Regeln
1. klare Zielbeschreibung
2. Scope begrenzt
3. Tests/Checks angegeben
4. Risiken benannt
5. Doku-Änderungen enthalten

## Review-Gates
- mindestens ein qualifizierter Review
- alle CI-Gates grün
- keine offenen P0-Risiken

## Merge-Regeln
- squash/merge gemäß Teamregel
- kein Merge bei roten Checks
- kein direktes Pushen auf geschützte Branches

## Issue- und Task-Kopplung
- Änderungen referenzieren Task-ID
- Changelog und Meeting aktualisieren

## Sicherheitsaspekte
- Secret Scanning aktiv
- Dependabot/Abhängigkeitsupdates prüfen
- Branch Protection aktivieren

## Abnahme-Check GITHUB
1. Branch-/PR-/Merge-Regeln definiert
2. Review-Gates enthalten
3. Sicherheitsregeln enthalten
4. Task-Referenzierung geregelt

---

## Mirror-Sync-Setup

### Git Remote hinzufügen

```bash
# Ins Projektverzeichnis wechseln
cd /Users/jeremy/dev/BIOMETRICS/BIOMETRICS

# GitHub Remote hinzufügen (Mirror)
git remote add github git@github.com:Delqhi/BIOMETRICS.git

# Oder HTTPS (wenn kein SSH-Key):
# git remote add github https://github.com/Delqhi/BIOMETRICS.git

# Remotes anzeigen
git remote -v
```

### Manueller Sync

```bash
# Push zu GitHub (Mirror-Mode)
git push github main --mirror

# Oder nur main branch:
git push github main
```

### Automatisierter Sync-Workflow

**.github/workflows/sync-mirror.yml:**

```yaml
name: Mirror to GitHub

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Mirror Push
        env:
          GH_TOKEN: ${{ secrets.GH_TOKEN }}
        run: |
          git push github main --mirror
```

### Automation mit GitLab CI/CD

**.gitlab-ci.yml:**

```yaml
sync_to_github:
  image: alpine:latest
  before_script:
    - apk add git
  script:
    - git remote add github git@github.com:Delqhi/BIOMETRICS.git
    - git push github main --mirror
  only:
    - main
```

### Cron-Job Sync (täglich)

```bash
# Täglicher Sync-Cronjob einrichten
crontab -e

# Folgende Zeile hinzufügen (täglich um 6 Uhr):
0 6 * * * cd /Users/jeremy/dev/BIOMETRICS/BIOMETRICS && git push github main --mirror
```

### GitHub Token Setup

1. **GitHub → Settings → Developer settings → Personal access tokens**
2. **Token generieren (classic)** mit `repo` Scope
3. **GitLab CI/CD Variable setzen:** `GH_TOKEN = <token>`

### Verifikation

```bash
# Prüfen ob Sync funktioniert hat
git fetch github
git log github/main --oneline -5

# Unterschiede anzeigen
git diff main github/main
```

---

**Letzte Aktualisierung:** Februar 2026