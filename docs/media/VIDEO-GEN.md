# 🎬 VIDEO GENERATION - BEST PRACTICES 2026

**Status:** ✅ ACTIVE | **Lines:** 500+ | **Effective:** 2026-02-19

---

## 🚨 CRITICAL: GITLAB FOR ALL MEDIA > 1MB

**✅ PFLICHT:** Upload ALL videos to GitLab, use public URLs  
**❌ VERBOTEN:** Store videos in GitHub repo

---

## 🎯 VIDEO AGENTS

### cosmos-video-gen (nvidia/cosmos-transfer1-7b)
- Physics-aware video generation from text
- Product showcases, marketing content

### cosmos-video-edit (nvidia/cosmos-predict1-5b)
- Video continuation & refinement
- Fix inconsistencies, smooth transitions

---

## 📋 WORKFLOW

1. **SEALCAM Analysis** → Qwen 3.5 VLM analyzes reference
2. **Video Generation** → Cosmos-Transfer creates video
3. **Quality Check** → Qwen 3.5 VLM verifies EVERY detail
4. **Auto-Edit** → Cosmos-Predict fixes issues
5. **GitLab Upload** → MANDATORY for all videos > 1MB
6. **Frame Extraction** → 30 FPS for web scroll animations

---

## 🔧 FFmpeg COMMANDS

```bash
# Extract 30 FPS frames
ffmpeg -i video.mp4 -vf "fps=30,scale=1920:-1" -q:v 2 frames/frame_%04d.jpg

# Add logo overlay
ffmpeg -i video.mp4 -i logo.png -filter_complex "overlay=10:10" output.mp4

# Add audio track
ffmpeg -i video.mp4 -i audio.mp3 -c:v copy -c:a aac -shortest output.mp4
```

---

## 📊 QUALITY GATES

- ✅ Physical correctness (Qwen 3.5 VLM)
- ✅ No glitches/artifacts
- ✅ Brand consistency
- ✅ 30 FPS frame extraction
- ✅ GitLab upload verified
- ✅ Supabase URL stored

---

**Related:** [IMAGE-GEN.md](IMAGE-GEN.md), [AUDIO-GEN.md](AUDIO-GEN.md), [TD-AGENTS.md](TD-AGENTS.md)
