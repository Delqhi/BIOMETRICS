# COMPLETE-SETUP.md - Vollständige Einrichtung für BIOMETRICS

**Stand:** 2026-02-18  
**Status:** ✅ PRODUCTION READY  
**Wichtig:** VOR dem Klonen des Repos ausführen!

---

## 🚨 WARUM DIESES DOKUMENT?

**Problem:** Agenten funktionieren auf dem Haupt-Mac, aber NICHT auf geklonten Macs!

**Grund:** Kritische Konfigurationen sind NICHT im Repository:
- ❌ NVIDIA_API_KEY in ~/.zshrc
- ❌ opencode.json in ~/.config/opencode/
- ❌ oh-my-opencode.json in ~/.config/opencode/
- ❌ openclaw.json in ~/.openclaw/
- ❌ Keine Setup-Dokumentation im Repo

**Lösung:** Diese Anleitung führt dich durch ALLE notwendigen Schritte!

---

## 📋 INHALTSVERZEICHNIS

1. [Phase 1: Base Tools installieren](#phase-1-base-tools-installieren)
2. [Phase 2: Environment Variables setzen](#phase-2-environment-variables-setzen)
3. [Phase 3: OpenCode konfigurieren](#phase-3-opencode-konfigurieren)
4. [Phase 4: Provider authentifizieren](#phase-4-provider-authentifizieren)
5. [Phase 5: oh-my-opencode konfigurieren](#phase-5-oh-my-opencode-konfigurieren)
6. [Phase 6: Plugins installieren](#phase-6-plugins-installieren)
7. [Phase 7: Setup verifizieren](#phase-7-setup-verifizieren)
8. [Phase 8: BIOMETRICS Repo klonen](#phase-8-biometrics-repo-klonen)

---

## PHASE 1: BASE TOOLS INSTALLIEREN (15 Min)

### 1.1 Homebrew installieren

```bash
# Prüfen ob bereits installiert
which brew

# Falls NICHT installiert:
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Verifizieren
brew --version
```

**✅ Checkliste:**
- [ ] Homebrew installiert
- [ ] Version angezeigt

### 1.2 Node.js 20+ installieren

```bash
# Installieren
brew install node@20

# Verifizieren
node --version  # Muss v20.x.x sein
npm --version   # Muss 10.x.x sein
```

**✅ Checkliste:**
- [ ] Node.js v20.x.x installiert
- [ ] npm 10.x.x installiert

### 1.3 pnpm installieren

```bash
# Installieren
brew install pnpm

# Verifizieren
pnpm --version  # Muss 8.x.x oder höher sein
```

**✅ Checkliste:**
- [ ] pnpm 8.x.x installiert

### 1.4 Python 3.11+ installieren

```bash
# Installieren
brew install python@3.11

# Verifizieren
python3 --version  # Muss 3.11.x sein
pip3 --version     # Muss vorhanden sein
```

**✅ Checkliste:**
- [ ] Python 3.11.x installiert
- [ ] pip3 vorhanden

### 1.5 Git installieren/aktualisieren

```bash
# Installieren (falls nicht vorhanden)
brew install git

# Verifizieren
git --version  # Muss 2.x.x sein
```

**✅ Checkliste:**
- [ ] Git 2.x.x installiert

---

## PHASE 2: ENVIRONMENT VARIABLES SETZEN (10 Min)

### 2.1 NVIDIA API Key besorgen

**SCHRITT 1: Registrieren**

```bash
# Browser öffnen
open https://build.nvidia.com/
```

1. Auf "Sign Up" klicken
2. Mit E-Mail registrieren
3. E-Mail verifizieren
4. Einloggen

**SCHRITT 2: API Key generieren**

1. Im Dashboard auf "API Keys" klicken
2. "Generate New Key" klicken
3. Key kopieren (beginnt mit `nvapi-`)
4. **Sofort speichern!** (kann nicht wieder angezeigt werden)

**✅ Checkliste:**
- [ ] NVIDIA Build Account erstellt
- [ ] API Key generiert
- [ ] Key sicher gespeichert

### 2.2 Environment Variables zu ~/.zshrc hinzufügen

```bash
# ~/.zshrc öffnen
nano ~/.zshrc
```

**Folgende Zeilen am ENDE hinzufügen:**

```bash
# NVIDIA API Keys (BIOMETRICS)
export NVIDIA_API_KEY="nvapi-DEIN_KEY_HIER"
export NVIDIA_NIM_API_KEY="nvapi-DEIN_KEY_HIER"

# PATH Erweiterungen (falls benötigt)
export PATH="$HOME/bin:$HOME/.local/bin:$PATH"

# Opencode Aliases
alias oc='opencode'
alias oc-models='opencode models'
alias oc-auth='opencode auth'
```

**WICHTIG:** Ersetze `DEIN_KEY_HIER` mit deinem echten Key!

**✅ Checkliste:**
- [ ] NVIDIA_API_KEY hinzugefügt
- [ ] NVIDIA_NIM_API_KEY hinzugefügt
- [ ] Keys korrekt eingefügt (keine Tippfehler!)

### 2.3 Shell NEU LADEN (KRITISCH!)

```bash
# KOMPLETTE Shell neu laden (MANDATORY!)
exec zsh

# Alternative (wenn exec nicht funktioniert):
# source ~/.zshrc

# Verifizieren dass Variablen geladen sind:
echo $NVIDIA_API_KEY
echo $NVIDIA_NIM_API_KEY
```

**⚠️ ACHTUNG:** Wenn du diesen Schritt überspringst, funktionieren die Keys NICHT!

**✅ Checkliste:**
- [ ] Shell mit `exec zsh` neu geladen
- [ ] Keys werden mit `echo` angezeigt

### 2.4 Weitere Environment Variables (Optional)

```bash
# Zu ~/.zshrc hinzufügen (falls benötigt):

# OpenCode Konfiguration
export OPENCODE_CONFIG="$HOME/.config/opencode"

# OpenClaw Konfiguration
export OPENCLAW_CONFIG="$HOME/.openclaw"

# Node.js Optimierungen
export NODE_OPTIONS="--max-old-space-size=4096"
```

---

## PHASE 3: OPENCODE KONFIGURIEREN (10 Min)

### 3.1 Opencode installieren

```bash
# Installieren
brew install opencode

# Verifizieren
opencode --version
```

**✅ Checkliste:**
- [ ] Opencode installiert
- [ ] Version wird angezeigt

### 3.2 opencode.json Existenz prüfen

```bash
# Prüfen ob Datei existiert
ls -la ~/.config/opencode/opencode.json

# Falls NICHT vorhanden → Opencode initialisieren:
opencode --version
```

**✅ Checkliste:**
- [ ] opencode.json existiert ODER wurde erstellt

### 3.3 opencode.json auf TIMEOUT prüfen (MANDATE 0.35!)

```bash
# KRITISCH: Auf timeout-Einträge prüfen
grep -r "timeout" ~/.config/opencode/opencode.json

# MUSS LEER SEIN! Falls Einträge gefunden:
nano ~/.config/opencode/opencode.json
# ALLE Zeilen mit "timeout" LÖSCHEN!

# JSON Syntax prüfen:
cat ~/.config/opencode/opencode.json | python3 -m json.tool
```

**⚠️ WICHTIG:** timeout in opencode.json ist VERBOTEN per MANDATE 0.35!

**✅ Checkliste:**
- [ ] KEINE timeout-Einträge in opencode.json
- [ ] JSON Syntax ist valide

### 3.4 opencode.json Inhalt prüfen

```bash
# Datei anzeigen:
cat ~/.config/opencode/opencode.json

# Sollte mindestens enthalten:
# - provider section
# - models section
# - KEIN timeout!
```

**Beispiel (korrekt):**

```json
{
  "provider": {
    "moonshot-ai": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Moonshot AI",
      "options": {
        "baseURL": "https://api.moonshot.ai/v1"
      }
    }
  },
  "models": {
    "moonshotai/kimi-k2.5": {
      "name": "Kimi K2.5",
      "limit": {
        "context": 1048576,
        "output": 65536
      }
    }
  }
}
```

---

## PHASE 4: PROVIDER AUTHENTIFIZIEREN (10 Min)

### 4.1 Opencode Authentifizierung

```bash
# OAuth Flow starten (öffnet Browser)
opencode auth login

# Im Browser:
# 1. Einloggen mit GitHub/Google
# 2. Zugriff erlauben
# 3. Zurück zum Terminal
```

**✅ Checkliste:**
- [ ] OAuth erfolgreich abgeschlossen

### 4.2 Provider hinzufügen

```bash
# moonshot-ai hinzufügen
opencode auth add moonshot-ai

# kimi-for-coding hinzufügen
opencode auth add kimi-for-coding

# opencode-zen hinzufügen (FREE!)
opencode auth add opencode-zen
```

**✅ Checkliste:**
- [ ] moonshot-ai hinzugefügt
- [ ] kimi-for-coding hinzugefügt
- [ ] opencode-zen hinzugefügt

### 4.3 Verifizierung

```bash
# Alle verfügbaren Modelle anzeigen
opencode models

# Erwartete Ausgabe (mindestens):
# - moonshotai/kimi-k2.5
# - qwen/qwen3.5-397b-a17b
# - opencode/minimax-m2.5-free
# - opencode/minimax-m2.5-free
# - opencode/minimax-m2.5-free
```

**✅ Checkliste:**
- [ ] Mindestens 5 Modelle sichtbar
- [ ] Alle 3 Provider vertreten

---

## PHASE 5: OH-MY-OPENCODE KONFIGURIEREN (5 Min)

### 5.1 oh-my-opencode.json prüfen

```bash
# Existenz prüfen
ls -la ~/.config/opencode/oh-my-opencode.json

# Auf timeout prüfen (MUSS LEER SEIN!)
grep -r "timeout" ~/.config/opencode/oh-my-opencode.json

# Inhalt anzeigen:
cat ~/.config/opencode/oh-my-opencode.json | python3 -m json.tool
```

**✅ Checkliste:**
- [ ] oh-my-opencode.json existiert
- [ ] KEINE timeout-Einträge
- [ ] JSON Syntax valide

### 5.2 Agent-Modelle konfigurieren

Die Datei sollte folgende Agent-Zuweisungen enthalten:

```json
{
  "agents": {
    "sisyphus": {
      "model": {
        "providerID": "moonshot-ai",
        "modelID": "kimi-k2.5"
      }
    },
    "librarian": {
      "model": {
        "providerID": "opencode-zen",
        "modelID": "zen/big-pickle"
      }
    },
    "explore": {
      "model": {
        "providerID": "opencode-zen",
        "modelID": "zen/big-pickle"
      }
    }
  }
}
```

**✅ Checkliste:**
- [ ] sisyphus → kimi-k2.5
- [ ] librarian → zen/big-pickle (FREE)
- [ ] explore → zen/big-pickle (FREE)

---

## PHASE 6: PLUGINS INSTALLIEREN (5 Min)

### 6.1 Google Antigravity Plugin

```bash
# Plugin installieren
opencode plugin add opencode-antigravity-auth

# OAuth durchführen
opencode auth login

# Status prüfen
opencode auth status
```

**⚠️ WICHTIG:** Private Gmail verwenden, NICHT Google Workspace!

**✅ Checkliste:**
- [ ] Antigravity Plugin installiert
- [ ] OAuth abgeschlossen
- [ ] Status zeigt "authenticated"

### 6.2 Oh-My-Opencode Plugin

```bash
# Plugin installieren (falls nicht automatisch)
opencode plugin add oh-my-opencode

# Verifizieren
opencode plugin list
```

**✅ Checkliste:**
- [ ] oh-my-opencode Plugin installiert
- [ ] In Plugin-Liste sichtbar

---

## PHASE 7: SETUP VERIFIZIEREN (5 Min)

### 7.1 Vollständige Verifizierung

```bash
# Verifikationsskript erstellen:
cat > /tmp/verify-biometrics.sh << 'VERIFY_EOF'
#!/bin/bash
echo "=== BIOMETRICS SETUP VERIFIKATION ==="
echo ""

# 1. NVIDIA API Key
if [ -z "$NVIDIA_API_KEY" ]; then
  echo "❌ NVIDIA_API_KEY NICHT gesetzt"
  echo "   Lösung: ~/.zshrc prüfen + 'exec zsh'"
  exit 1
else
  echo "✅ NVIDIA_API_KEY: ${#NVIDIA_API_KEY} Zeichen"
fi

# 2. Opencode
if ! command -v opencode &> /dev/null; then
  echo "❌ Opencode NICHT installiert"
  echo "   Lösung: 'brew install opencode'"
  exit 1
else
  echo "✅ Opencode: $(opencode --version)"
fi

# 3. Modelle
model_count=$(opencode models 2>/dev/null | wc -l | tr -d ' ')
if [ "$model_count" -lt 5 ]; then
  echo "❌ Zu wenige Modelle (< 5)"
  echo "   Lösung: 'opencode auth add <provider>'"
  exit 1
else
  echo "✅ Modelle: $model_count konfiguriert"
fi

# 4. Timeout Check (MANDATE 0.35)
timeout_count=$(grep -r "timeout" ~/.config/opencode/opencode.json 2>/dev/null | wc -l | tr -d ' ')
if [ "$timeout_count" -gt 0 ]; then
  echo "❌ MANDATE 0.35 VERLETZUNG: timeout in opencode.json"
  echo "   Lösung: timeout-Einträge SOFORT löschen!"
  exit 1
else
  echo "✅ Kein timeout in opencode.json"
fi

# 5. OpenClaw (optional)
if [ -f ~/.openclaw/openclaw.json ]; then
  echo "✅ OpenClaw konfiguriert"
  if command -v openclaw &> /dev/null; then
    openclaw doctor --fix 2>/dev/null && echo "✅ OpenClaw Health Check bestanden"
  fi
else
  echo "⚠️  OpenClaw nicht konfiguriert (optional)"
fi

# 6. Shell Config
if grep -q "NVIDIA_API_KEY" ~/.zshrc 2>/dev/null; then
  echo "✅ NVIDIA_API_KEY in ~/.zshrc"
else
  echo "⚠️  NVIDIA_API_KEY NICHT in ~/.zshrc"
fi

echo ""
echo "=== ALLE CHECKS BESTANDEN ==="
echo "Dein Mac ist BEREIT für BIOMETRICS!"
VERIFY_EOF

# Skript ausführen:
chmod +x /tmp/verify-biometrics.sh
/tmp/verify-biometrics.sh
```

**✅ Checkliste:**
- [ ] Alle Verifikations-Checks bestanden
- [ ] KEINE Fehlermeldungen

### 7.2 Manuelle Verifikation

```bash
# Diese Befehle sollten OHNE FEHLER laufen:

# 1. Opencode Version
opencode --version

# 2. Modelle anzeigen
opencode models | head -10

# 3. Auth Status
opencode auth status

# 4. Environment
echo "NVIDIA: ${NVIDIA_API_KEY:0:20}..."
```

---

## PHASE 8: BIOMETRICS REPO KLONEN (5 Min)

### 8.1 Repository klonen

```bash
# In gewünschtes Verzeichnis navigieren
cd ~/dev  # Oder dein bevorzugtes Verzeichnis

# Repo klonen
git clone https://github.com/Delqhi/BIOMETRICS.git

# In Repo wechseln
cd BIOMETRICS

# Verzeichnisstruktur prüfen
ls -la
```

**✅ Checkliste:**
- [ ] Repo erfolgreich geklont
- [ ] In BIOMETRICS Verzeichnis gewechselt

### 8.2 BIOMETRICS CLI installieren

```bash
# In CLI Verzeichnis wechseln
cd biometrics-cli

# Dependencies installieren
pnpm install

# Global verlinken
pnpm link --global

# CLI testen
biometrics --version
```

**✅ Checkliste:**
- [ ] CLI installiert
- [ ] Version wird angezeigt

### 8.3 Onboarding starten

```bash
# Onboarding Wizard starten
biometrics

# Oder vollständige Version:
biometrics-onboard
```

**✅ Checkliste:**
- [ ] Onboarding erfolgreich
- [ ] Alle System-Checks bestanden

---

## 🚨 HÄUFIGE FEHLER & LÖSUNGEN

### Fehler: "NVIDIA_API_KEY not found"

**Ursache:** Shell wurde nicht neu geladen

**Lösung:**
```bash
exec zsh
echo $NVIDIA_API_KEY  # Sollte Key anzeigen
```

### Fehler: "Provider not found"

**Ursache:** Provider nicht authentifiziert

**Lösung:**
```bash
opencode auth add moonshot-ai
opencode auth add kimi-for-coding
opencode auth add opencode-zen
```

### Fehler: "timeout in config"

**Ursache:** MANDATE 0.35 Verletzung

**Lösung:**
```bash
nano ~/.config/opencode/opencode.json
# ALLE Zeilen mit "timeout" löschen!
```

### Fehler: "Models not loading"

**Ursache:** API Key nicht geladen

**Lösung:**
```bash
# 1. ~/.zshrc prüfen
grep NVIDIA ~/.zshrc

# 2. Shell neu laden
exec zsh

# 3. Verifizieren
echo $NVIDIA_API_KEY
```

---

## ✅ ABSCHLUSS-CHECKLISTE

**Phase 1: Base Tools**
- [ ] Homebrew installiert
- [ ] Node.js 20+ installiert
- [ ] pnpm installiert
- [ ] Python 3.11+ installiert
- [ ] Git installiert

**Phase 2: Environment Variables**
- [ ] NVIDIA API Key besorgt
- [ ] ~/.zshrc aktualisiert
- [ ] Shell neu geladen (exec zsh)
- [ ] Keys verifiziert

**Phase 3: OpenCode**
- [ ] Opencode installiert
- [ ] opencode.json existiert
- [ ] KEIN timeout in opencode.json
- [ ] JSON Syntax valide

**Phase 4: Provider**
- [ ] OAuth durchgeführt
- [ ] moonshot-ai hinzugefügt
- [ ] kimi-for-coding hinzugefügt
- [ ] opencode-zen hinzugefügt
- [ ] Modelle verifiziert (5+)

**Phase 5: oh-my-opencode**
- [ ] oh-my-opencode.json existiert
- [ ] KEIN timeout
- [ ] Agent-Modelle konfiguriert

**Phase 6: Plugins**
- [ ] Antigravity Plugin installiert
- [ ] OAuth abgeschlossen
- [ ] oh-my-opencode Plugin installiert

**Phase 7: Verifikation**
- [ ] Verifikationsskript bestanden
- [ ] Manuelle Checks bestanden
- [ ] KEINE Fehlermeldungen

**Phase 8: BIOMETRICS**
- [ ] Repo geklont
- [ ] CLI installiert
- [ ] Onboarding abgeschlossen

---

## 📚 WEITERFÜHRENDE DOKUMENTATION

- **PROVIDER-SETUP.md:** Detaillierte Provider-Konfiguration
- **SHELL-SETUP.md:** Shell-Konfiguration und Environment
- **VERIFICATION.md:** Testbefehle und Troubleshooting
- **BIOMETRICS/README.md:** Projekt-Übersicht
- **BIOMETRICS/PROVIDER.md:** Provider-Referenz

---

## 🎯 NÄCHSTE SCHRITTE

1. **BIOMETRICS/README.md lesen** - Projekt verstehen
2. **BIOMETRICS/AGENTS.md lesen** - Agenten-Nutzung
3. **Erstes Projekt erstellen** - `biometrics create`
4. **Dokumentation pflegen** - Changes dokumentieren

**Willkommen bei BIOMETRICS!** 🚀

---

**Dokument Version:** 1.0  
**Letztes Update:** 2026-02-18  
**Status:** ✅ PRODUCTION READY
