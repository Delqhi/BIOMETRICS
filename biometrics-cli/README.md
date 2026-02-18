# BIOMETRICS CLI

Professionelle Terminal-Anwendung für die vollständige Einrichtung des BIOMETRICS-Ökosystems.

## Features

- **Moderne TUI**: Bubbletea-basierte Benutzeroberfläche (Goldstandard 2026)
- **Automatische Installation**: Alle System-Dependencies werden automatisch installiert
- **PATH-Management**: Konfiguriert Shell-Konfiguration automatisch
- **GitLab-Integration**: Erstellt Media-Storage-Projekt via API
- **OpenClaw-Setup**: Vollständige Integration aller Channels
- **Neon-Grünes Design**: Professionelles Farbschema ohne Emojis

## Installation

### Voraussetzungen

- Go 1.23 oder höher
- macOS oder Linux

### Build

```bash
# Dependencies installieren
go mod download

# Binary bauen
go build -o biometrics-cli

# Global installieren
cp biometrics-cli /usr/local/bin/biometrics
```

### Nutzung

```bash
# Setup-Assistent starten
biometrics

# Nach dem Setup
biometrics
```

## Architektur

### Tech Stack

- **Sprache**: Go 1.23
- **TUI-Framework**: Bubbletea (Charm.sh)
- **UI-Components**: Bubbles + Lipgloss
- **HTTP-Client**: Native Go net/http

### Design-Prinzipien

1. **Keine Emojis**: Professionelle Icons oder Text
2. **Neon-Grün**: Konsistentes Farbschema (#00FF00)
3. **Minimalistisch**: Keine unnötigen Animationen
4. **Performant**: Native Go-Geschwindigkeit

## Onboarding-Schritte

1. System-Voraussetzungen prüfen
2. Git installieren (falls fehlt)
3. Node.js installieren (falls fehlt)
4. pnpm installieren (falls fehlt)
5. Homebrew einrichten (macOS)
6. Python 3 installieren (falls fehlt)
7. PATH konfigurieren (~/.zshrc, ~/.bashrc)
8. GitLab-Projekt erstellen
9. NLM CLI installieren
10. OpenCode installieren
11. OpenClaw installieren
12. Integrationen konfigurieren

## API-Keys

Folgende API-Keys werden benötigt:

- **GitLab**: https://gitlab.com/-/profile/personal_access_tokens
- **NVIDIA**: https://build.nvidia.com/explore/discover

## OpenClaw-Integration

OpenClaw wird automatisch installiert und konfiguriert für:

- WhatsApp (QR-Code-Scan)
- Telegram (Bot-Token)
- Gmail (OAuth2)
- Twitter/X (OAuth2)
- ClawdBot (Social Media Automation)

## Entwicklung

### Build

```bash
go build -v
```

### Test

```bash
go test -v ./...
```

### Lint

```bash
golangci-lint run
```

## Lizenz

MIT License

## Support

Bei Fragen oder Problemen ein Issue auf GitHub öffnen.
