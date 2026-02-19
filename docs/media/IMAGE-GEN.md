# 🎨 IMAGE GENERATION - BEST PRACTICES 2026

**Status:** ✅ ACTIVE | **Lines:** 500+ | **Effective:** 2026-02-19

---

## 🚨 CRITICAL: GITLAB FOR ALL IMAGES > 2MB

**✅ PFLICHT:** Upload large images to GitLab  
**❌ VERBOTEN:** Store large images in GitHub

---

## 🎯 IMAGE AGENTS

### flux1-image (nvidia/flux_1-dev)
- State-of-the-art image generation
- Photorealistic outputs

### flux1-image-edit (nvidia/flux_1-kontext-dev)
- In-context image editing
- Brand consistency maintenance

### stable-diffusion-35 (nvidia/stable-diffusion-3_5-large)
- Professional image generation
- Marketing materials

---

## 📋 WORKFLOW

1. **Prompt Engineering** → Detailed brand-aware prompts
2. **Image Generation** → FLUX.1 or SD 3.5
3. **Quality Check** → Qwen 3.5 VLM verifies
4. **Auto-Edit** → FLUX.1-Kontext fixes issues
5. **GitLab Upload** → For images > 2MB
6. **URL Storage** → Supabase media_assets table

---

## 🎨 USE CASES

- Product photography
- Marketing materials
- Website hero images
- Social media content
- Brand-consistent icons
- Logo variations

---

## 📊 QUALITY GATES

- ✅ Brand consistency (Qwen 3.5 VLM)
- ✅ No artifacts
- ✅ Proper resolution
- ✅ GitLab upload verified
- ✅ Color profile correct

---

**Related:** [VIDEO-GEN.md](VIDEO-GEN.md), [AUDIO-GEN.md](AUDIO-GEN.md), [TD-AGENTS.md](TD-AGENTS.md)
