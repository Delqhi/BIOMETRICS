# 🎬 VIDEO SCRIPT: Complete Setup Guide

**Duration:** 15 minutes  
**Level:** Beginner  
**Prerequisites:** None

---

## SCENE 1: Introduction (0:00-0:45)

### Visual: Title Card
- **Text:** "BIOMETRICS Complete Setup Guide"
- **Subtitle:** "From Zero to Hero in 15 Minutes"
- **Background:** BIOMETRICS logo

### Voiceover:
"Willkommen zum kompletten Setup-Guide für BIOMETRICS! In diesem Video zeige ich dir Schritt für Schritt, wie du deine Entwicklungsumgebung so konfigurierst, dass alle Agenten perfekt funktionieren - genau wie auf unserem Haupt-Mac."

### Visual: Bullet Points
- ✅ NVIDIA API Key besorgen
- ✅ OpenCode installieren
- ✅ Provider konfigurieren
- ✅ Agenten testen
- ✅ Troubleshooting

### Voiceover:
"Wir werden einen NVIDIA API Key besorgen, OpenCode installieren, alle Provider konfigurieren und am Ende testen ob alles funktioniert. Los geht's!"

---

## SCENE 2: NVIDIA API Key Besorgen (0:45-3:00)

### Visual: Browser Recording
- **URL:** https://build.nvidia.com/
- **Action:** Navigate to website

### Voiceover:
"Öffne deinen Browser und gehe zu build.nvidia.com. Das ist die offizielle NVIDIA Developer Platform."

### Visual: Screen Recording
- **Action:** Click "Sign Up" or "Login"
- **Show:** Login form

### Voiceover:
"Wenn du noch keinen Account hast, klicke auf 'Sign Up' und erstelle einen kostenlosen Account. Falls du schon einen hast, einfach einloggen."

### Visual: Dashboard Navigation
- **Action:** Click on "API Keys" in left sidebar
- **Show:** API Keys dashboard

### Voiceover:
"Im Dashboard klickst du links auf 'API Keys'. Hier kannst du einen neuen Key erstellen."

### Visual: Create API Key
- **Action:** Click "Create New API Key"
- **Show:** Key generation (blur actual key)
- **Highlight:** Key starts with "nvapi-"

### Voiceover:
"Klicke auf 'Create New API Key'. Dein Key wird generiert und beginnt mit 'nvapi-'. **WICHTIG:** Kopiere diesen Key sofort und speichere ihn an einem sicheren Ort! Du kannst ihn später nicht mehr einsehen."

### Visual: Terminal
- **Command:** `nano ~/.zshrc`
- **Action:** Open .zshrc file

### Voiceover:
"Öffne jetzt dein Terminal und gib ein: nano ~/.zshrc Das öffnet deine Shell-Konfigurationsdatei."

### Visual: Editor
- **Action:** Scroll to end of file
- **Type:** `export NVIDIA_API_KEY="nvapi-YOUR_KEY_HERE"`

### Voiceover:
"Scrolle ans Ende der Datei und füge diese Zeile hinzu: export NVIDIA_API_KEY等于 dein-key-hier. Ersetze 'dein-key-hier' mit deinem kopierten Key."

### Visual: Save and Exit
- **Action:** Ctrl+O, Enter, Ctrl+X
- **Command:** `source ~/.zshrc`

### Voiceover:
"Speichere mit Strg+O, bestätige mit Enter, und schließe mit Strg+X. Lade dann die Konfiguration neu mit: source ~/.zshrc"

---

## SCENE 3: OpenCode Installation (3:00-5:30)

### Visual: Terminal
- **Command:** `node --version`
- **Show:** Node.js version output

### Voiceover:
"Bevor wir OpenCode installieren können, prüfen wir ob Node.js installiert ist. Gib ein: node --version. Wenn eine Versionsnummer erscheint, ist Node.js installiert."

### Visual: Installation Command
- **Command:** `brew install node` (if not installed)

### Voiceover:
"Falls 'command not found' erscheint, installiere Node.js mit: brew install node"

### Visual: OpenCode Install
- **Command:** `npm install -g opencode`
- **Show:** Installation progress

### Voiceover:
"Jetzt installieren wir OpenCode global mit: npm install -g opencode. Das kann ein paar Sekunden dauern."

### Visual: Verification
- **Command:** `opencode --version`
- **Show:** Version number

### Voiceover:
"Überprüfe die Installation mit: opencode --version. Du solltest eine Versionsnummer sehen."

### Visual: OpenClaw Install
- **Command:** `npm install -g @anthropic-ai/openclaw`

### Voiceover:
"Optional aber empfohlen: Installiere auch OpenClaw mit: npm install -g @anthropic-ai/openclaw"

---

## SCENE 4: Provider Authentifizierung (5:30-8:00)

### Visual: Terminal
- **Command:** `opencode auth add nvidia-nim`
- **Show:** Browser opens for OAuth

### Voiceover:
"Jetzt authentifizieren wir die Provider. Starte mit: opencode auth add nvidia-nim. Ein Browser-Fenster öffnet sich für die OAuth-Authentifizierung."

### Visual: Browser OAuth
- **Show:** OAuth consent screen
- **Action:** Click "Allow"

### Voiceover:
"Melde dich an und erlaube den Zugriff. Das verbindet deinen NVIDIA Account mit OpenCode."

### Visual: More Auth Commands
```bash
opencode auth add moonshot-ai
opencode auth add kimi-for-coding
opencode auth add opencode-zen
opencode auth add huggingface
opencode auth add groq
```

### Voiceover:
"Wiederhole das für alle required Provider: moonshot-ai, kimi-for-coding, opencode-zen, huggingface, und groq."

### Visual: Verify Auth
- **Command:** `opencode auth list`
- **Show:** List of authenticated providers

### Voiceover:
"Überprüfe alle Authentifizierungen mit: opencode auth list. Du solltest alle Provider in der Liste sehen."

---

## SCENE 5: opencode.json Konfiguration (8:00-11:00)

### Visual: Terminal
- **Command:** `mkdir -p ~/.config/opencode`

### Voiceover:
"Erstelle das Konfigurationsverzeichnis: mkdir -p ~/.config/opencode"

### Visual: Editor
- **Command:** `cat > ~/.config/opencode/opencode.json << 'EOF'`
- **Show:** Complete JSON config (from COMPLETE-SETUP.md)

### Voiceover:
"Jetzt erstellen wir die opencode.json. Kopiere diesen kompletten JSON-Block aus der Dokumentation. Er enthält alle Provider-Konfigurationen, insbesondere den NVIDIA NIM Provider mit Qwen 3.5 397B."

### Visual: Highlight Key Sections
- **Highlight:** nvidia-nim provider section
- **Highlight:** model ID: qwen/qwen3.5-397b-a17b
- **Highlight:** context limit: 262144

### Voiceover:
"Achte besonders auf den nvidia-nim Provider. Hier ist die model ID 'qwen/qwen3.5-397b-a17b' mit einem gigantischen Context-Limit von 262K Tokens!"

### Visual: Save and Verify
- **Command:** `cat ~/.config/opencode/opencode.json`

### Voiceover:
"Speichere die Datei und verifiziere mit cat, dass alles korrekt ist."

---

## SCENE 6: oh-my-opencode.json (11:00-13:00)

### Visual: Terminal
- **Command:** `cat > ~/.config/opencode/oh-my-opencode.json << 'EOF'`
- **Show:** Complete agent configuration

### Voiceover:
"Jetzt konfigurieren wir die Agent-Zuweisungen. Diese Datei bestimmt welcher Agent welches Modell verwendet."

### Visual: Highlight Agent Mappings
- **Highlight:** sisyphus → nvidia-nim/qwen-3.5-397b
- **Highlight:** prometheus → nvidia-nim/qwen-3.5-397b
- **Highlight:** Categories with temperatures

### Voiceover:
"Alle wichtigen Agenten - sisyphus, prometheus, metis, oracle - verwenden alle das NVIDIA Qwen 3.5 Modell. Die Categories definieren Temperature und Thinking-Budget."

---

## SCENE 7: Shell Reload & Verification (13:00-14:30)

### Visual: Terminal
- **Command:** `exec zsh`
- **Show:** Terminal restarts

### Voiceover:
"**KRITISCHER SCHRITT:** Lade deine Shell komplett neu mit: exec zsh. Ohne diesen Schritt werden die neuen Konfigurationen nicht wirksam!"

### Visual: Environment Check
```bash
echo $NVIDIA_API_KEY
opencode models | grep nvidia
```

### Voiceover:
"Verifiziere dass alles geladen wurde: echo $NVIDIA_API_KEY sollte deinen Key anzeigen. opencode models | grep nvidia sollte die NVIDIA Modelle auflisten."

### Visual: Test Command
- **Command:** `opencode --model nvidia-nim/qwen-3.5-397b "Hello, this is a test"`
- **Show:** Model response

### Voiceover:
"Teste das Modell mit einem einfachen Prompt. Du solltest eine Antwort von Qwen 3.5 erhalten."

---

## SCENE 8: Conclusion & Next Steps (14:30-15:00)

### Visual: Summary Slide
- ✅ NVIDIA API Key gesetzt
- ✅ OpenCode installiert
- ✅ Provider authentifiziert
- ✅ Konfigurationen erstellt
- ✅ Shell neu geladen
- ✅ Tests erfolgreich

### Voiceover:
"Perfekt! Deine Umgebung ist jetzt komplett eingerichtet. Alle Agenten sollten jetzt funktionieren."

### Visual: Next Steps
- Clone BIOMETRICS repo
- Run `npm install`
- Start building!

### Voiceover:
"Jetzt kannst du das BIOMETRICS Repo klonen und loslegen. Die Links zu allen Konfigurationsdateien findest du in der Beschreibung."

### Visual: End Card
- **Text:** "Danke fürs Zuschauen!"
- **Links:** docs/setup/COMPLETE-SETUP.md
- **Links:** docs/UNIVERSAL-BLUEPRINT.md

### Voiceover:
"Viel Erfolg beim Builden! Wenn du Fragen hast, schau in unsere Dokumentation. Bis zum nächsten Video!"

---

## PRODUCTION NOTES

### Recording Software:
- **Screen Recording:** OBS Studio or QuickTime
- **Audio:** Blue Yeti or similar
- **Resolution:** 1920x1080
- **FPS:** 30

### Post-Production:
- **Editing:** DaVinci Resolve or Final Cut Pro
- **Captions:** Auto-generate + manual correction
- **Thumbnail:** Canva template
- **Upload:** YouTube + Vimeo

### Distribution:
- **YouTube:** BIOMETRICS channel
- **Vimeo:** Private link for docs
- **GitHub:** Embed in docs
- **Duration:** 15:00 minutes exactly

---

**SCRIPT COMPLETE** ✅
