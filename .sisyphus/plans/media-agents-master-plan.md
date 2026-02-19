# 🎬 MEDIA AGENTS MASTER PLAN - VIDEO, IMAGE, AUDIO, 3D GENERATION

**Created:** 2026-02-19  
**Status:** READY FOR EXECUTION  
**Priority:** CRITICAL - Enterprise God-Mode Media Pipeline  
**Execution:** Run `/start-work media-agents-master-plan`

---

## 🎯 OBJECTIVE

Create complete media generation infrastructure for BIOMETRICS project with:
- ✅ Video generation agents (Cosmos-Transfer, Cosmos-Predict)
- ✅ Image generation agents (FLUX.1, Stable Diffusion 3.5)
- ✅ Audio generation agents (Magpie-TTS, StudioVoice)
- ✅ 3D generation agents (Microsoft TRELLIS)
- ✅ Qwen 3.5 VLM quality checks for ALL media
- ✅ GitLab integration (ALL media > 1MB to GitLab)
- ✅ TD-Agent orchestration scripts
- ✅ Documentation (VIDEO-GEN.md, IMAGE-GEN.md, AUDIO-GEN.md, TD-AGENTS.md)

---

## 📋 EXECUTION WAVES

### WAVE 1: FOUNDATION (4 tasks, parallel)
1. **Create oh-my-opencode.json** with ALL media agents
2. **Create project structure** (/inputs, /outputs, /assets, /scripts, /logs, /skills)
3. **Create requirements.txt** with all Python dependencies
4. **Create .env.example** with all required API keys

### WAVE 2: VIDEO AGENTS (6 tasks, parallel)
5. **Create VIDEO-GEN.md** documentation
6. **Create cosmos-video-gen agent** script (NVIDIA Cosmos-Transfer1-7B)
7. **Create cosmos-video-edit agent** script (NVIDIA Cosmos-Predict1-5B)
8. **Create sealcam_analysis.py** for video analysis
9. **Create video_quality_check.py** (Qwen 3.5 VLM verification)
10. **Create upload-to-gitlab.sh** for video uploads

### WAVE 3: IMAGE AGENTS (5 tasks, parallel)
11. **Create IMAGE-GEN.md** documentation
12. **Create flux1-image agent** script (FLUX.1-dev)
13. **Create flux1-image-edit agent** script (FLUX.1-Kontext-dev)
14. **Create stable-diffusion-35 agent** script (SD 3.5 Large)
15. **Create image_quality_check.py** (Qwen 3.5 VLM verification)

### WAVE 4: AUDIO AGENTS (4 tasks, parallel)
16. **Create AUDIO-GEN.md** documentation
17. **Create magpie-voice agent** script (Magpie-TTS Multilingual)
18. **Create studio-voice agent** script (Audio optimization)
19. **Create audio_quality_check.py** (Qwen 3.5 VLM verification)
20. **Create audio-sync.py** for video+audio merging

### WAVE 5: 3D AGENTS (3 tasks, parallel)
21. **Create trellis-3d agent** script (Microsoft TRELLIS)
22. **Create 3d_quality_check.py** (Qwen 3.5 VLM verification)
23. **Create render-frames.py** for 3D→Video pipeline

### WAVE 6: TD-AGENT ORCHESTRATION (4 tasks, parallel)
24. **Create TD-AGENTS.md** master documentation
25. **Create nim_engine.py** central API wrapper
26. **Create video_processor.py** FFmpeg automation
27. **Create scroll-animation.js** for web (Apple-effect)

### WAVE 7: INTEGRATION & TESTING (3 tasks, parallel)
28. **Create complete-pipeline.sh** end-to-end test
29. **Test GitLab upload** for all media types
30. **Create example website** with scroll animations

---

## 🔧 TECHNICAL REQUIREMENTS

### oh-my-opencode.json Configuration
```json
{
  "agents": {
    "flux1-image": {
      "model": "nvidia/flux_1-dev",
      "category": "artistry"
    },
    "flux1-image-edit": {
      "model": "nvidia/flux_1-kontext-dev",
      "category": "artistry"
    },
    "stable-diffusion-35": {
      "model": "nvidia/stable-diffusion-3_5-large",
      "category": "artistry"
    },
    "cosmos-video-gen": {
      "model": "nvidia/cosmos-transfer1-7b",
      "category": "visual-engineering"
    },
    "cosmos-video-edit": {
      "model": "nvidia/cosmos-predict1-5b",
      "category": "visual-engineering"
    },
    "trellis-3d": {
      "model": "microsoft/trellis",
      "category": "visual-engineering"
    },
    "magpie-voice": {
      "model": "nvidia/magpie-tts-multilingual",
      "category": "artistry"
    },
    "studio-voice": {
      "model": "nvidia/studiovoice",
      "category": "artistry"
    }
  }
}
```

### Project Structure
```
/PROJECT_ROOT/
├── .env
├── requirements.txt
├── oh-my-opencode.json
│
├── /inputs/
│   ├── /references/       # Original videos (SealCam analysis)
│   └── /brand_assets/     # Product images for 3D gen
│
├── /outputs/
│   ├── /videos/           # Final 4K videos
│   └── /assets/           # All media assets
│
├── /assets/
│   ├── /3d/               # .glb/.usd from TRELLIS
│   ├── /renders/          # 3D renders (360° frames)
│   ├── /frames/           # JPG sequences (30 FPS)
│   ├── /images/           # Generated images
│   └── /audio/            # Generated audio
│
├── /scripts/
│   ├── nim_engine.py          # Central API wrapper
│   ├── video_processor.py     # FFmpeg automation
│   ├── sealcam_analysis.py    # Video analysis
│   ├── video_quality_check.py # Qwen 3.5 VLM
│   ├── image_quality_check.py # Qwen 3.5 VLM
│   ├── audio_quality_check.py # Qwen 3.5 VLM
│   ├── 3d_quality_check.py    # Qwen 3.5 VLM
│   ├── upload-to-gitlab.sh    # GitLab uploads
│   ├── audio-sync.py          # Video+Audio merge
│   ├── render-frames.py       # 3D→Video pipeline
│   └── complete-pipeline.sh   # End-to-end test
│
├── /logs/
│   └── /thinking/         # Qwen 3.5 <think> logs
│
└── /skills/
    └── production_skill.md # TD-Agent knowledge base
```

### requirements.txt
```txt
requests>=2.31.0
python-dotenv>=1.0.0
ffmpeg-python>=0.2.0
opencv-python>=4.8.0
Pillow>=10.0.0
numpy>=1.24.0
```

### .env.example
```bash
# NVIDIA NIM API
NVIDIA_API_KEY=nvapi-YOUR_KEY

# GitLab Media Storage
GITLAB_TOKEN=glpat-YOUR_TOKEN
GITLAB_MEDIA_PROJECT_ID=79575238

# Supabase Database
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-anon-key
```

---

## 📚 DOCUMENTATION REQUIREMENTS

### VIDEO-GEN.md Must Include:
- ✅ Cosmos-Transfer1-7B usage (text→video)
- ✅ Cosmos-Predict1-5B usage (video→video)
- ✅ SealCam Framework analysis
- ✅ Qwen 3.5 VLM quality checks
- ✅ GitLab upload automation
- ✅ FFmpeg scripts (30 FPS extraction, overlays, audio sync)
- ✅ Scroll animation setup (Apple-effect)
- ✅ Complete pipeline examples

### IMAGE-GEN.md Must Include:
- ✅ FLUX.1-dev usage (text→image)
- ✅ FLUX.1-Kontext-dev usage (in-context editing)
- ✅ Stable Diffusion 3.5 Large usage
- ✅ Qwen 3.5 VLM quality checks
- ✅ GitLab upload automation
- ✅ Brand consistency workflows
- ✅ Logo overlay automation

### AUDIO-GEN.md Must Include:
- ✅ Magpie-TTS Multilingual (brand voice)
- ✅ StudioVoice (audio optimization)
- ✅ Background Noise Removal
- ✅ Qwen 3.5 VLM quality checks
- ✅ GitLab upload automation
- ✅ Audio+Video sync workflows

### TD-AGENTS.md Must Include:
- ✅ Technical Director role definition
- ✅ Complete pipeline orchestration
- ✅ nim_engine.py documentation
- ✅ video_processor.py documentation
- ✅ Scroll animation JavaScript
- ✅ Error handling & diagnostics
- ✅ Skill documentation format

---

## ✅ ACCEPTANCE CRITERIA

### Code Quality:
- [ ] All scripts have error handling
- [ ] All API calls use .env variables (no hardcoded keys)
- [ ] All scripts have docstrings and comments
- [ ] LSP diagnostics: 0 errors, 0 warnings

### Testing:
- [ ] Video generation test (Cosmos-Transfer)
- [ ] Video editing test (Cosmos-Predict)
- [ ] Image generation test (FLUX.1)
- [ ] Image editing test (FLUX.1-Kontext)
- [ ] Audio generation test (Magpie-TTS)
- [ ] 3D generation test (TRELLIS)
- [ ] GitLab upload test (all media types)
- [ ] Qwen 3.5 VLM quality check test
- [ ] Scroll animation test (30 FPS)

### Documentation:
- [ ] VIDEO-GEN.md complete (500+ lines)
- [ ] IMAGE-GEN.md complete (500+ lines)
- [ ] AUDIO-GEN.md complete (500+ lines)
- [ ] TD-AGENTS.md complete (500+ lines)
- [ ] All scripts documented
- [ ] Example workflows included

### GitLab Integration:
- [ ] ALL videos > 1MB uploaded to GitLab
- [ ] ALL images > 2MB uploaded to GitLab
- [ ] ALL audio > 1MB uploaded to GitLab
- [ ] Public URLs working
- [ ] URLs stored in Supabase media_assets table
- [ ] NO media files in GitHub repo

### Quality Gates:
- [ ] Qwen 3.5 VLM verifies ALL videos
- [ ] Qwen 3.5 VLM verifies ALL images
- [ ] Qwen 3.5 VLM verifies ALL audio
- [ ] Qwen 3.5 VLM verifies ALL 3D assets
- [ ] Zero defects in final outputs
- [ ] Brand consistency verified

---

## 🚀 DELEGATION STRATEGY

### Wave 1 (Foundation):
```typescript
// Task 1-4: Quick setup tasks
task(category="quick", load_skills=["git-master"], run_in_background=true, prompt="Create oh-my-opencode.json with ALL media agents...")
task(category="quick", load_skills=[], run_in_background=true, prompt="Create project directory structure...")
task(category="quick", load_skills=[], run_in_background=true, prompt="Create requirements.txt with dependencies...")
task(category="quick", load_skills=[], run_in_background=true, prompt="Create .env.example with API keys...")
```

### Wave 2-5 (Media Agents):
```typescript
// Each wave: 1 doc + 3-4 agent scripts + quality check
task(category="writing", load_skills=[], run_in_background=true, prompt="Create VIDEO-GEN.md documentation...")
task(category="visual-engineering", load_skills=[], run_in_background=true, prompt="Create cosmos-video-gen agent script...")
task(category="visual-engineering", load_skills=[], run_in_background=true, prompt="Create cosmos-video-edit agent script...")
task(category="ultrabrain", load_skills=[], run_in_background=true, prompt="Create sealcam_analysis.py...")
// ... etc for each wave
```

### Wave 6 (TD-Orchestration):
```typescript
// TD-Agent master scripts
task(category="ultrabrain", load_skills=[], run_in_background=true, prompt="Create TD-AGENTS.md master doc...")
task(category="deep", load_skills=[], run_in_background=true, prompt="Create nim_engine.py central wrapper...")
task(category="deep", load_skills=[], run_in_background=true, prompt="Create video_processor.py FFmpeg automation...")
task(category="visual-engineering", load_skills=[], run_in_background=true, prompt="Create scroll-animation.js...")
```

### Wave 7 (Integration):
```typescript
// End-to-end testing
task(category="deep", load_skills=[], run_in_background=true, prompt="Create complete-pipeline.sh test...")
task(category="quick", load_skills=[], run_in_background=true, prompt="Test GitLab upload for all media...")
task(category="visual-engineering", load_skills=["playwright"], run_in_background=true, prompt="Create example website with scroll animations...")
```

---

## 🔍 QUALITY VERIFICATION (QWEN 3.5 VLM)

### For EACH generated media:
```typescript
task(
  category="ultrabrain",
  model="qwen/qwen3.5-397b-a17b",
  prompt=`
## 🎯 QUALITY CHECK: [MEDIA TYPE]

PRÜFE JEDES DETAIL:
1. [Type-specific checks]
2. Physikalische Korrektheit?
3. Keine Artefakte/Glitches?
4. Brand Identity gewahrt?
5. Technische Spezifikationen erfüllt?

WENN FEHLER:
→ Liste ALLE Fehler auf
→ Empfehle Korrektur mit [EDIT-AGENT]
→ Auto-fix einleiten

NICHT ABNEHMEN bevor PERFEKT!
`
)
```

---

## 📊 SUCCESS METRICS

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Video Quality** | 100% Qwen approved | Quality check logs |
| **Image Quality** | 100% Qwen approved | Quality check logs |
| **Audio Quality** | 100% Qwen approved | Quality check logs |
| **3D Quality** | 100% Qwen approved | Quality check logs |
| **GitLab Upload** | 100% success rate | Upload logs |
| **Scroll Smoothness** | 30 FPS constant | Browser DevTools |
| **Brand Consistency** | 100% verified | Qwen 3.5 VLM |
| **Documentation** | 2000+ lines total | Line count |
| **Test Coverage** | 100% passing | Test results |

---

## 🚫 FORBIDDEN ACTIONS

- ❌ Store media > 1MB in GitHub (ALWAYS GitLab!)
- ❌ Use local file paths (ALWAYS GitLab URLs!)
- ❌ Skip Qwen 3.5 VLM quality checks (MANDATORY!)
- ❌ Hardcode API keys (ALWAYS .env!)
- ❌ Manual FFmpeg commands (ALWAYS automated scripts!)
- ❌ Create duplicate files (ALWAYS check existence first!)
- ❌ Skip SealCam analysis (ALWAYS analyze first!)

---

## 🎯 FINAL DELIVERABLES

1. ✅ oh-my-opencode.json with 8+ media agents
2. ✅ Complete project structure (6 directories)
3. ✅ VIDEO-GEN.md (500+ lines)
4. ✅ IMAGE-GEN.md (500+ lines)
5. ✅ AUDIO-GEN.md (500+ lines)
6. ✅ TD-AGENTS.md (500+ lines)
7. ✅ 12+ Python/Bash scripts
8. ✅ 1 JavaScript scroll animation
9. ✅ All tests passing
10. ✅ GitLab integration working
11. ✅ Example website with animations
12. ✅ Zero media files in GitHub repo

---

**Plan Status:** ✅ READY FOR EXECUTION  
**Next Step:** Run `/start-work media-agents-master-plan`  
**Estimated Duration:** 45-60 minutes (all waves parallel)
