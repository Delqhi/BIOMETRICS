# BIOMETRICS CLI Installation

## WICHTIG: Go Binary (NICHT Node.js!)

Seit Version 2.0 ist die biometrics-cli eine **Go-Anwendung** (nicht mehr Node.js).

## Installation

### Option 1: Automatisch (Empfohlen)

```bash
# Repo klonen
git clone https://github.com/Delqhi/BIOMETRICS.git
cd BIOMETRICS/biometrics-cli

# Dependencies installieren
go mod download

# Binary bauen
go build -o biometrics-cli

# Global installieren
sudo cp biometrics-cli /usr/local/bin/biometrics
sudo cp biometrics-cli /usr/local/bin/biometrics-onboard
sudo chmod +x /usr/local/bin/biometrics /usr/local/bin/biometrics-onboard
```

### Option 2: Manuelles Update (wenn bereits installiert)

Falls du eine alte Node.js-Version installiert hast:

```bash
# Alten Symlink entfernen
sudo rm -f /usr/local/bin/biometrics
sudo rm -f /usr/local/bin/biometrics-onboard

# Neue Go Binary installieren
cd /Users/jeremy/dev/BIOMETRICS/biometrics-cli
sudo cp biometrics-cli /usr/local/bin/biometrics
sudo cp biometrics-cli /usr/local/bin/biometrics-onboard
sudo chmod +x /usr/local/bin/biometrics /usr/local/bin/biometrics-onboard
```

## Testen

```bash
# Binary testen
biometrics
# ODER
biometrics-onboard
```

## Fehlerbehebung

### "Cannot find module '/Users/jeremy/dev/BIOMETRICS/biometrics-cli/src/index.js'"

**Ursache:** Alte Node.js-Installation noch aktiv

**Lösung:**
```bash
sudo rm -f /usr/local/bin/biometrics
sudo rm -f /usr/local/bin/biometrics-onboard
cd /Users/jeremy/dev/BIOMETRICS/biometrics-cli
sudo cp biometrics-cli /usr/local/bin/biometrics
sudo cp biometrics-cli /usr/local/bin/biometrics-onboard
```

### "Command not found: biometrics"

**Ursache:** `/usr/local/bin` nicht im PATH

**Lösung:**
```bash
export PATH="/usr/local/bin:$PATH"
# Oder zu ~/.zshrc hinzufügen:
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## Deinstallation

```bash
sudo rm -f /usr/local/bin/biometrics
sudo rm -f /usr/local/bin/biometrics-onboard
```

## Version

- **Aktuell:** 2.0.0 (Go + Bubbletea TUI)
- **Vorher:** 1.0.0 (Node.js + Inquirer)

## Tech Stack

- **Sprache:** Go 1.23
- **TUI:** Bubbletea (Charm.sh)
- **UI:** Bubbles + Lipgloss
- **Binary:** 4.3MB (static linked)
