#!/bin/bash

# 🚨 INTERACTIVE OPENCODE.JSON GENERATOR
# This script interactively creates ~/.config/opencode/opencode.json

echo "START: Interactive OpenCode.json Generator"
echo "======================================"
echo ""

# Step 1: Get NVIDIA API Key
echo "STEP: Schritt 1: NVIDIA API Key"
echo "----------------------------"
read -p "Hast du bereits einen NVIDIA API Key? (y/n): " has_key

if [ "$has_key" != "y" ]; then
    echo ""
    echo "OPEN: Öffne https://build.nvidia.com/ in deinem Browser"
    echo "   1. Einloggen oder Account erstellen"
    echo "   2. Auf 'API Keys' klicken"
    echo "   3. 'Create New API Key' klicken"
    echo "   4. Key kopieren (beginnt mit nvapi-...)"
    echo ""
    read -p "Drück Enter wenn du den Key kopiert hast: "
fi

echo ""
read -p "Füge deinen NVIDIA API Key ein: " nvidia_key

# Validate key format
if [[ ! "$nvidia_key" =~ ^nvapi- ]]; then
    echo "ERROR: Fehler: Key muss mit 'nvapi-' beginnen!"
    exit 1
fi

echo "SUCCESS: Key validiert!"
echo ""

# Step 2: Create directory
echo "DIR: Schritt 2: Verzeichnis erstellen"
echo "------------------------------------"
mkdir -p ~/.config/opencode
echo "SUCCESS: Verzeichnis erstellt: ~/.config/opencode"
echo ""

# Step 3: Generate opencode.json
echo "CONFIG:  Schritt 3: opencode.json generieren"
echo "----------------------------------------"

cat > ~/.config/opencode/opencode.json << EOF
{
  "\$schema": "https://opencode.ai/config.json",
  "model": "nvidia-nim/qwen-3.5-397b",
  "default_agent": "sisyphus",
  "provider": {
    "nvidia-nim": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "NVIDIA NIM (Qwen 3.5)",
      "options": {
        "baseURL": "https://integrate.api.nvidia.com/v1"
      },
      "models": {
        "qwen-3.5-397b": {
          "name": "Qwen 3.5 397B (NVIDIA NIM)",
          "id": "qwen/qwen3.5-397b-a17b",
          "limit": {
            "context": 262144,
            "output": 32768
          }
        }
      }
    }
  },
  "plugin": [
    "opencode-antigravity-auth@latest",
    "oh-my-opencode"
  ]
}
EOF

echo "SUCCESS: opencode.json erstellt!"
echo ""

# Step 4: Add to .zshrc
echo "ENV: Schritt 4: Environment Variable in ~/.zshrc"
echo "-----------------------------------------------"

if ! grep -q "export NVIDIA_API_KEY" ~/.zshrc 2>/dev/null; then
    echo "" >> ~/.zshrc
    echo "# NVIDIA NIM Configuration (added by BIOMETRICS setup)" >> ~/.zshrc
    echo "export NVIDIA_API_KEY=\"$nvidia_key\"" >> ~/.zshrc
    echo "SUCCESS: NVIDIA_API_KEY zu ~/.zshrc hinzugefügt"
else
    echo "WARNING:  NVIDIA_API_KEY ist bereits in ~/.zshrc"
fi

echo ""

# Step 5: Verification
echo "SUCCESS: VERIFIKATION"
echo "--------------"
echo ""
echo "Deine Konfiguration wurde erstellt!"
echo ""
echo "Nächste Schritte:"
echo "1. Shell neu laden: exec zsh"
echo "2. Verifizieren: opencode auth add nvidia-nim"
echo "3. Testen: opencode models | grep nvidia"
echo ""
echo "DOCS: Vollständige Anleitung: docs/setup/COMPLETE-SETUP.md"
echo ""

read -p "Shell jetzt neu laden? (y/n): " reload

if [ "$reload" = "y" ]; then
    exec zsh
fi
