# BIOMETRICS SETUP CHECKLISTE

## ✅ Status: READY

**Projekt:** BIOMETRICS - AI Media Generation Pipeline  
**Version:** 1.0.0  
**Datum:** 2026-02-19

---

## 📋 Datei-Status

| Datei | Status | Zeilen |
|-------|--------|--------|
| `.env.example` | ✅ Vollständig | 31 |
| `oh-my-opencode.json` | ✅ Vollständig (9 Agents) | 136 |
| `requirements.txt` | ✅ Vollständig | 34 |
| `package.json` | ❌ Nicht vorhanden | - |

---

## 🔑 Konfigurationen

### .env.example (VOLLSTÄNDIG)

```
NVIDIA_API_KEY=nvapi-YOUR_NVIDIA_API_KEY_HERE
COSMOS_MODEL=nvidia/cosmos-transfer1-7b
GITLAB_TOKEN=glpat-YOUR_GITLAB_TOKEN_HERE
GITLAB_MEDIA_PROJECT_ID=79575238
GITLAB_MEDIA_NAMESPACE=your-namespace
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-anon-key
SUPABASE_SERVICE_KEY=your-supabase-service-key
CLOUDFLARE_TUNNEL_ID=your-tunnel-id
CLOUDFLARE_DOMAIN=your-domain.com
ENVIRONMENT=development
DEBUG=true
LOG_LEVEL=INFO
```

### oh-my-opencode.json (9 AGENTS KONFIGURIERT)

| Agent | Modell | Kategorie |
|-------|--------|-----------|
| flux1-image | nvidia/flux_1-dev | artistry |
| flux1-image-edit | nvidia/flux_1-kontext-dev | artistry |
| stable-diffusion-35 | nvidia/stable-diffusion-3_5-large | artistry |
| cosmos-video-gen | nvidia/cosmos-transfer1-7b | visual-engineering |
| cosmos-video-edit | nvidia/cosmos-predict1-5b | visual-engineering |
| trellis-3d | microsoft/trellis | visual-engineering |
| magpie-voice | nvidia/magpie-tts-multilingual | artistry |
| studio-voice | nvidia/studiovoice | artistry |
| qwen-vlm | qwen/qwen3.5-397b-a17b | ultrabrain |

---

## ⏳ User muss tun

### 1. API Keys holen

| Service | URL | Key-Typ |
|---------|-----|----------|
| NVIDIA NIM | https://build.nvidia.com/ | API Key |
| GitLab | https://gitlab.com/-/profile/personal_access_tokens | Token |
| Supabase | https://supabase.com/dashboard/ | URL + Keys |

### 2. Environment konfigurieren

```bash
cd /Users/jeremy/dev/BIOMETRICS
cp .env.example .env
nano .env
```

### 3. Python Dependencies installieren

```bash
pip install -r requirements.txt
```

### 4. NVIDIA NIM Authentifizierung

```bash
opencode auth add nvidia-nim
opencode models | grep nvidia
```

### 5. Verify

```bash
# Check NVIDIA models
opencode models | grep nvidia

# Check agents
cat oh-my-opencode.json | jq '.agents | keys'
```

---

## ✅ Ready when:

- [ ] `.env` mit echten Keys konfiguriert
- [ ] `pip install -r requirements.txt` ohne Fehler
- [ ] `opencode auth add nvidia-nim` erfolgreich
- [ ] `opencode models | grep nvidia` zeigt NVIDIA Modelle

---

## 🚀 Nächste Schritte

1. **NVIDIA API Key holen:** https://build.nvidia.com/
2. **.env konfigurieren:** Keys eintragen
3. **Dependencies installieren:** pip + npm (falls package.json benötigt)
4. **OpenCode authentifizieren:** opencode auth add nvidia-nim
5. **Testen:** opencode models

---

## 📞 Support

Bei Fragen: Siehe `/Users/jeremy/dev/BIOMETRICS/README.md`
