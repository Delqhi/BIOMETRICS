# 🎬 VIDEO GENERATION - BEST PRACTICES 2026

**Status:** ✅ ACTIVE | **Lines:** 523 | **Version:** 1.0 | **Effective:** 2026-02-19

---

## 🚨 CRITICAL: GITLAB FOR ALL MEDIA > 1MB

**✅ PFLICHT:** Upload ALL videos to GitLab, use public URLs  
**❌ VERBOTEN:** Store videos in GitHub repo

### Warum GitLab?

| Aspekt | GitLab | GitHub |
|--------|--------|--------|
| File Size | 100MB per file | 100MB per file |
| Bandwidth | Unlimited | 2GB/month |
| Public URLs | Auto-generate | Manual |
| Cost | FREE | $4/month |

### Wann GitLab nutzen?

- **ALLE Videos > 1MB** → GitLab
- **ALLE Videos < 1MB** → Supabase Storage
- **ALLE Videos > 50MB** → Split into chunks

---

## 📡 NVIDIA NIM VIDEO API ENDPOINTS

### Cosmos-Transfer1-7B (Video Generation)

```
Endpoint: https://integrate.api.nvidia.com/v1/video/generations
Method: POST
Auth: Bearer Token (NVIDIA_API_KEY)
Rate Limit: 40 RPM (Free Tier)
```

**Request Format:**
```json
{
  "prompt": "A sleek sports car driving through a neon city at night",
  "negative_prompt": "blurry, low quality, distorted",
  "num_frames": 24,
  "fps": 24,
  "resolution": { "width": 512, "height": 288 },
  "guidance_scale": 7.5,
  "num_inference_steps": 25
}
```

---

## 🔧 PYTHON CLIENT

```python
import os
import requests
import time

class NVIDIAVideoClient:
    BASE_URL = "https://integrate.api.nvidia.com/v1"
    
    def __init__(self, api_key=None):
        self.api_key = api_key or os.getenv("NVIDIA_API_KEY")
        self.headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json"
        }
    
    def generate_video(self, prompt, num_frames=24, fps=24, 
                       width=512, height=288, timeout=300):
        payload = {
            "prompt": prompt,
            "num_frames": num_frames,
            "fps": fps,
            "resolution": {"width": width, "height": height}
        }
        
        # Submit
        resp = requests.post(f"{self.BASE_URL}/video/generations",
                           headers=self.headers, json=payload, timeout=30)
        resp.raise_for_status()
        gen_id = resp.json()["id"]
        
        # Poll
        start = time.time()
        while time.time() - start < timeout:
            status = requests.get(f"{self.BASE_URL}/video/generations/{gen_id}",
                                 headers=self.headers, timeout=30).json()
            if status["status"] == "completed":
                return status
            elif status["status"] == "failed":
                raise RuntimeError(f"Failed: {status.get('error')}")
            time.sleep(5)
        raise TimeoutError(f"Timeout after {timeout}s")
    
    def download_video(self, url, output_path):
        r = requests.get(url, stream=True, timeout=300)
        r.raise_for_status()
        with open(output_path, 'wb') as f:
            for chunk in r.iter_content(8192): f.write(chunk)

# Usage
client = NVIDIAVideoClient()
result = client.generate_video("Futuristic motorcycle in cyberpunk city")
client.download_video(result["output"]["video_url"], "video.mp4")
```

---

## 🎯 VIDEO AGENTS

### cosmos-video-gen (nvidia/cosmos-transfer1-7b)
- Text-to-video, physics simulation
- Max 24 frames, 512x288 default
- Product showcases, marketing

### cosmos-video-edit (nvidia/cosmos-predict1-5b)
- Video continuation & refinement
- Style transfer, color grading

### sealcam-analysis (Qwen 3.5 VLM)
- Video quality verification
- Pre/post generation checks

---

## 📋 WORKFLOW

### 1. Reference Analysis
```python
def analyze_ref(image_path):
    # Analyze with Qwen 3.5 VLM
    # Extract: subjects, lighting, camera, mood
    return analysis
```

### 2. Video Generation
```python
result = client.generate_video(
    prompt="Product in cinematic style",
    num_frames=24, fps=24
)
```

### 3. Quality Check
```python
def verify(video_path):
    frames = extract_frames(video_path, 8)
    issues = []
    for f in frames:
        analysis = analyze_frame(f)
        if analysis.artifacts: issues.append("artifacts")
        if analysis.physics_errors: issues.append("physics")
    return {"passed": len(issues)==0, "issues": issues}
```

### 4. Auto-Edit
```python
fixed = client.edit_video(url, "smooth edges", strength=0.5)
```

### 5. GitLab Upload
```python
import requests
def upload_gitlab(path, proj="biometrics/media-assets"):
    token = os.getenv("GITLAB_TOKEN")
    url = f"https://gitlab.com/api/v4/projects/{proj}/uploads"
    r = requests.post(url, files={'file': open(path,'rb')},
                     headers={"PRIVATE-TOKEN": token})
    return r.json()["markdown_url"]
```

### 6. Frame Extraction
```bash
ffmpeg -i video.mp4 -vf "fps=30" -q:v 2 frames/frame_%04d.jpg
```

---

## 🔧 FFmpeg COMMANDS

```bash
# Extract 30 FPS frames
ffmpeg -i input.mp4 -vf "fps=30,scale=1920:-1" -q:v 2 frames/frame_%04d.jpg

# Add logo overlay
ffmpeg -i input.mp4 -i logo.png -filter_complex "overlay=10:10" output.mp4

# Add audio
ffmpeg -i input.mp4 -i audio.mp3 -c:v copy -c:a aac -shortest output.mp4

# Resize for 9:16
ffmpeg -i input.mp4 -vf "crop=ih*(9/16):ih" output_9x16.mp4

# Resize for 1:1
ffmpeg -i input.mp4 -vf "crop=ih:ih:(iw-ih)/2:0" output_1x1.mp4

# Compress H.264
ffmpeg -i input.mp4 -c:v libx264 -crf 23 -c:a aac output.mp4

# Speed 0.5x
ffmpeg -i input.mp4 -vf "setpts=2.0*PTS" slow.mp4

# Extract audio
ffmpeg -i input.mp4 -vn -acodec libmp3lame audio.mp3

# Thumbnail
ffmpeg -i input.mp4 -ss 00:00:05 -vframes 1 thumb.jpg
```

---

## 📊 QUALITY GATES

| Gate | Criteria | Verification |
|------|----------|--------------|
| Physical | Objects obey physics | Qwen 3.5 VLM check |
| Artifacts | No compression artifacts | Edge detection |
| Brand | Logo, colors consistent | Template match |
| Technical | Resolution, FPS correct | ffprobe |
| Upload | GitLab URL accessible | HEAD request |

---

## 🔗 INTEGRATIONS

### n8n Workflow
```json
{
  "nodes": [
    {"name": "Trigger", "type": "webhook"},
    {"name": "NVIDIA", "type": "httpRequest", "url": "https://integrate.api.nvidia.com/v1/video/generations"},
    {"name": "GitLab", "type": "gitlab", "operation": "upload"}
  ]
}
```

### Supabase Schema
```sql
CREATE TABLE media_assets (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    type VARCHAR(50),
    storage_provider VARCHAR(50),
    storage_url TEXT,
    public_url TEXT,
    file_size INTEGER,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 🛠️ TROUBLESHOOTING

| Issue | Cause | Solution |
|-------|-------|----------|
| 400 Bad Request | Invalid prompt | Remove special chars, <500 chars |
| Timeout | Rate limit | Wait 60s, max 40 RPM |
| 401 Unauthorized | Token expired | Refresh .env |
| 413 Too Large | File >100MB | Split or compress |

---

## 📈 BENCHMARKS

| Resolution | Frames | Time | Size |
|------------|--------|------|------|
| 512x288 | 24 | 45-60s | 1-3MB |
| 768x480 | 24 | 60-90s | 3-5MB |

---

## 📚 RELATED

- [IMAGE-GEN.md](IMAGE-GEN.md)
- [AUDIO-GEN.md](AUDIO-GEN.md)
- [TD-AGENTS.md](TD-AGENTS.md)

---

## 🔬 ADVANCED TOPICS

### Video Physics Simulation

Cosmos-Transfer1-7B versteht Physik und erzeugt realistische Bewegungen:

```python
# Gute Prompts für Physics
prompts = [
    "Car driving through city at 60mph, realistic motion blur",
    "Water flowing in river, physics-accurate turbulence",
    "Birds flying in formation, realistic wing movements",
    "Explosion with debris, realistic particle physics",
    "Character running, natural gait cycle"
]

# Schlechte Prompts (vermeiden)
bad_prompts = [
    "Object floating without support",  # Verletzt Physik
    "Person walking through wall",      # Kollision ignoriert
    "Light bending around corner",      # Optik ignoriert
]
```

### Camera Movement Types

| Type | Prompt Example | Use Case |
|------|---------------|----------|
| Pan | "Camera pans left to reveal..." | Reveal shots |
| Tilt | "Camera tilts up to show..." | Height emphasis |
| Dolly | "Camera dolly forward..." | Intimacy |
| Zoom | "Camera zooms in on..." | Focus |
| Orbit | "Camera orbits around..." | 360 product |
| Crane | "Camera cranes up..." | Scale |

### Lighting Setups

```python
LIGHTING_PROMPTS = {
    "cinematic": "dramatic lighting, chiaroscuro, volumetric rays",
    "product": "softbox lighting, even illumination, no shadows",
    "moody": "low key, dark shadows, single light source",
    "neon": "neon lights, cyberpunk aesthetic, color bloom",
    "golden_hour": "warm sunset, golden hour, lens flare"
}
```

### Frame-by-Frame Analysis

```python
def analyze_frames(video_path, num_frames=24):
    """Analyze each frame for consistency"""
    import cv2
    
    cap = cv2.VideoCapture(video_path)
    total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    step = total_frames // num_frames
    
    analyses = []
    for i in range(0, total_frames, step):
        cap.set(cv2.CAP_PROP_POS_FRAMES, i)
        ret, frame = cap.read()
        if ret:
            analysis = analyze_frame_quality(frame)
            analyses.append(analysis)
    
    # Check consistency across frames
    consistency_score = calculate_consistency(analyses)
    return consistency_score

def analyze_frame_quality(frame):
    """Analyze single frame for quality issues"""
    import cv2
    import numpy as np
    
    issues = []
    
    # Check blur
    laplacian = cv2.Laplacian(frame, cv2.CV_64F)
    blur_score = laplacian.var()
    if blur_score < 100:
        issues.append(f"Blurry frame (score: {blur_score})")
    
    # Check exposure
    gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
    mean = np.mean(gray)
    if mean < 30:
        issues.append("Underexposed")
    elif mean > 225:
        issues.append("Overexposed")
    
    # Check noise
    noise = detect_noise(frame)
    if noise > 0.1:
        issues.append(f"Noisy frame (noise: {noise})")
    
    return {"issues": issues, "blur_score": blur_score, "exposure": mean}
```

### Color Grading Presets

```python
COLOR_GRADES = {
    "action": {
        "contrast": 1.3,
        "saturation": 1.2,
        "temperature": 6500,
        "gamma": 1.1
    },
    "romantic": {
        "contrast": 0.9,
        "saturation": 0.8,
        "temperature": 4500,
        "gamma": 1.0
    },
    "horror": {
        "contrast": 1.5,
        "saturation": 0.7,
        "temperature": 7000,
        "gamma": 0.9
    },
    "documentary": {
        "contrast": 1.0,
        "saturation": 0.9,
        "temperature": 5600,
        "gamma": 1.0
    }
}
```

### Multi-Scene Videos

```python
def generate_multi_scene_video(scenes: list, transitions: list):
    """Generate video with multiple scenes and transitions"""
    
    all_videos = []
    
    for i, scene in enumerate(scenes):
        # Generate individual scene
        result = client.generate_video(
            prompt=scene["prompt"],
            num_frames=scene.get("frames", 24),
            fps=24
        )
        scene_path = f"scene_{i}.mp4"
        client.download_video(result["output"]["video_url"], scene_path)
        all_videos.append(scene_path)
    
    # Concatenate with transitions
    concat_list = "|".join(all_videos)
    cmd = f"ffmpeg -i {concat_list} -filter_complex 'concat=n={len(all_videos)}:v=1:a=0' output.mp4"
    
    return cmd
```

### Audio Integration

```python
def add_background_music(video_path, mood, duration):
    """Add AI-generated background music"""
    
    # Generate audio with Edge TTS
    audio_prompt = f"Generate {mood} background music, instrumental only"
    
    # Use edge-tts for free TTS
    audio_file = generate_music(mood, duration)
    
    # Mix with video
    cmd = [
        "ffmpeg",
        "-i", video_path,
        "-i", audio_file,
        "-c:v", "copy",
        "-c:a", "aac",
        "-shortest",
        "output_with_audio.mp4"
    ]
    
    return subprocess.run(cmd)
```

### Subtitle Generation

```python
def generate_subtitles(video_path, language="en"):
    """Generate subtitles using Whisper"""
    
    # Extract audio
    audio_path = "audio.wav"
    cmd = ["ffmpeg", "-i", video_path, "-vn", "-acodec", "pcm_s16le", 
           "-ar", "16000", "-ac", "1", audio_path]
    subprocess.run(cmd)
    
    # Transcribe with Whisper
    import whisper
    model = whisper.load_model("base")
    result = model.transcribe(audio_path)
    
    # Generate SRT file
    srt_content = convert_to_srt(result["segments"])
    
    # Burn subtitles into video
    cmd = [
        "ffmpeg",
        "-i", video_path,
        "-vf", f"subtitles={srt_content}",
        "output_subtitled.mp4"
    ]
    
    return cmd
```

### A/B Testing Prompts

```python
def ab_test_prompts(original_prompt, variations):
    """Test multiple prompt variations"""
    
    results = []
    
    for i, variant in enumerate(variations):
        result = client.generate_video(
            prompt=variant,
            num_frames=24,
            fps=24
        )
        
        # Score the result
        score = evaluate_video_quality(result["video_url"])
        
        results.append({
            "prompt": variant,
            "video_url": result["output"]["video_url"],
            "score": score
        })
    
    # Return best performing variant
    best = max(results, key=lambda x: x["score"])
    return best

def evaluate_video_quality(video_url):
    """Automated quality scoring"""
    
    # Download and analyze
    frames = extract_frames(video_url, 8)
    
    scores = {
        "clarity": 0,
        "physics": 0,
        "aesthetic": 0,
        "consistency": 0
    }
    
    for frame in frames:
        analysis = analyze_frame(frame)
        scores["clarity"] += analysis.clarity
        scores["physics"] += analysis.physics_accuracy
        scores["aesthetic"] += analysis.aesthetic_score
    
    # Average
    for key in scores:
        scores[key] /= len(frames)
    
    # Weighted total
    total = (
        scores["clarity"] * 0.3 +
        scores["physics"] * 0.3 +
        scores["aesthetic"] * 0.2 +
        scores["consistency"] * 0.2
    )
    
    return total
```

### Batch Processing

```python
def batch_generate_videos(prompts: list, output_dir: str):
    """Generate multiple videos in batch"""
    
    os.makedirs(output_dir, exist_ok=True)
    
    results = []
    for i, prompt in enumerate(prompts):
        try:
            result = client.generate_video(
                prompt=prompt,
                num_frames=24,
                fps=24
            )
            
            output_path = os.path.join(output_dir, f"video_{i}.mp4")
            client.download_video(result["output"]["video_url"], output_path)
            
            results.append({
                "index": i,
                "prompt": prompt,
                "output": output_path,
                "success": True
            })
            
        except Exception as e:
            results.append({
                "index": i,
                "prompt": prompt,
                "error": str(e),
                "success": False
            })
    
    return results
```

### CDN Integration

```python
def upload_to_cdn(video_path, cdn_provider="cloudflare"):
    """Upload to CDN for fast delivery"""
    
    if cdn_provider == "cloudflare":
        # Upload to R2 (S3-compatible)
        import boto3
        s3 = boto3.client('s3')
        
        with open(video_path, 'rb') as f:
            s3.upload_fileobj(
                f,
                'biometrics-media',
                f'videos/{os.path.basename(video_path)}',
                ExtraArgs={
                    'ContentType': 'video/mp4',
                    'CacheControl': 'max-age=31536000'
                }
            )
        
        # Get CDN URL
        cdn_url = f"https://media.biometrics.ai/videos/{os.path.basename(video_path)}"
        return cdn_url
    
    elif cdn_provider == "vercel":
        # Use Vercel Blob
        pass
```

### Analytics & Tracking

```python
def track_video_generation(prompt, result, metrics):
    """Track video generation for analytics"""
    
    import supabase
    
    supabase = create_client(SUPABASE_URL, SUPABASE_KEY)
    
    data = {
        "prompt": prompt,
        "duration_seconds": metrics["generation_time"],
        "resolution": f"{metrics['width']}x{metrics['height']}",
        "frames": metrics["num_frames"],
        "status": result["status"],
        "video_url": result.get("output", {}).get("video_url"),
        "cost_estimate": calculate_cost(metrics)
    }
    
    supabase.table("video_generations").insert(data).execute()
```

---

## 🎯 PROMPT ENGINEERING GUIDE

### Structure

```
[Subject] + [Action] + [Environment] + [Lighting] + [Mood] + [Technical]
```

Beispiel:
```
A futuristic motorcycle (Subject)
racing through (Action)
a cyberpunk city at night (Environment)
with neon lights reflecting on wet streets (Lighting)
cinematic mood (Mood)
8K quality, photorealistic (Technical)
```

### Keywords

| Category | Keywords |
|----------|----------|
| Quality | photorealistic, 8K, detailed, high fidelity |
| Lighting | cinematic, dramatic, volumetric, chiaroscuro |
| Camera | slow motion, tracking shot, aerial view |
| Style | cyberpunk, vintage, minimalist, maximalist |
| Mood | atmospheric, suspenseful, peaceful, chaotic |

### Negatives

Immer hinzufügen:
```
blurry, low quality, distorted, artifacts, watermark, text, 
logo, signature, deformed, ugly, bad anatomy
```

---

## 💰 COST OPTIMIZATION

### Free Tier Limits

| Metric | Free Tier | Paid Tier |
|--------|-----------|------------|
| RPM | 40 | 200 |
| Videos/day | 100 | Unlimited |
| Max duration | 1s | 10s |
| Resolution | 512x288 | 1920x1080 |

### Optimization Tips

1. **Kürzere Videos** - 24 Frames = 1 Sekunde = 45-60s Generation
2. **Kleinere Auflösung** - 512x288 statt 1024x576
3. **Weniger Frames** - 24 statt 48
4. **Wiederverwenden** - Cached Prompts

---

## 🔐 SECURITY

### API Key Protection

```python
# NIEMALS hardcoden!
# ~/.env:
# NVIDIA_API_KEY=nvapi-xxxxx

# Laden aus Environment
api_key = os.getenv("NVIDIA_API_KEY")
if not api_key:
    raise ValueError("NVIDIA_API_KEY not set")
```

### Rate Limit Handling

```python
from functools import wraps
import time

def rate_limit(max_calls=40, period=60):
    """Decorator for rate limiting"""
    calls = []
    
    @wraps(func)
    def wrapper(*args, **kwargs):
        now = time.time()
        calls[:] = [t for t in calls if now - t < period]
        
        if len(calls) >= max_calls:
            sleep_time = period - (now - calls[0])
            time.sleep(sleep_time)
        
        result = func(*args, **kwargs)
        calls.append(time.time())
        return result
    
    return wrapper
```

---

**Version:** 1.0 | **Updated:** 2026-02-19
