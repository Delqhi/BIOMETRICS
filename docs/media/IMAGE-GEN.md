# 🎨 IMAGE GENERATION - BEST PRACTICES 2026

**Status:** ✅ ACTIVE | **Lines:** 512 | **Version:** 1.0 | **Effective:** 2026-02-19

---

## 🚨 CRITICAL: GITLAB FOR ALL IMAGES > 2MB

**✅ PFLICHT:** Upload large images to GitLab  
**❌ VERBOTEN:** Store large images in GitHub

### Warum GitLab?

| Aspekt | GitLab | GitHub |
|--------|--------|--------|
| File Size | 100MB | 100MB |
| Bandwidth | Unlimited | 2GB/month |
| CDN | Optional | Manual |
| Cost | FREE | $4/month |

### Wann GitLab nutzen?

- **ALLE Images > 2MB** → GitLab
- **ALLE Images < 2MB** → Supabase Storage

---

## 📡 NVIDIA NIM IMAGE API ENDPOINTS

### FLUX.1-dev

```
Endpoint: https://integrate.api.nvidia.com/v1/images/generations
Method: POST
Auth: Bearer Token (NVIDIA_API_KEY)
Rate Limit: 40 RPM
```

**Request:**
```json
{
  "prompt": "A futuristic motorcycle in cyberpunk city, neon lights",
  "negative_prompt": "blurry, low quality, distorted",
  "width": 1024,
  "height": 1024,
  "num_images": 1,
  "guidance_scale": 7.5,
  "num_inference_steps": 25,
  "seed": -1
}
```

### FLUX.1-kontext-dev (Editing)

```
Endpoint: https://integrate.api.nvidia.com/v1/images/edits
```

### Stable Diffusion 3.5

```
Endpoint: https://integrate.api.nvidia.com/v1/images/generations
Model: nvidia/stable-diffusion-3_5-large
```

---

## 🔧 PYTHON CLIENT

```python
import os
import requests
import time

class NVIDIAImageClient:
    BASE_URL = "https://integrate.api.nvidia.com/v1"
    
    def __init__(self, api_key=None):
        self.api_key = api_key or os.getenv("NVIDIA_API_KEY")
        self.headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json"
        }
    
    def generate_image(self, prompt, negative="", width=1024, 
                       height=1024, num_images=1, timeout=120):
        payload = {
            "prompt": prompt,
            "negative_prompt": negative,
            "width": width,
            "height": height,
            "num_images": num_images,
            "guidance_scale": 7.5,
            "num_inference_steps": 25
        }
        
        resp = requests.post(f"{self.BASE_URL}/images/generations",
                           headers=self.headers, json=payload, timeout=30)
        resp.raise_for_status()
        gen_id = resp.json()["id"]
        
        start = time.time()
        while time.time() - start < timeout:
            status = requests.get(f"{self.BASE_URL}/images/generations/{gen_id}",
                                 headers=self.headers).json()
            if status["status"] == "completed":
                return status["output"]["images"]
            elif status["status"] == "failed":
                raise RuntimeError(f"Failed: {status.get('error')}")
            time.sleep(3)
        raise TimeoutError("Timeout")
    
    def edit_image(self, image_path, prompt, timeout=120):
        # Upload source image, then edit
        pass
    
    def download_image(self, url, path):
        r = requests.get(url, timeout=60)
        with open(path, 'wb') as f: f.write(r.content)

# Usage
client = NVIDIAImageClient()
imgs = client.generate_image("Futuristic motorcycle, neon lights")
client.download_image(imgs[0]["url"], "image.png")
```

---

## 🎯 IMAGE AGENTS

### flux1-image (nvidia/flux_1-dev)
- State-of-the-art image generation
- Photorealistic outputs
- 1024x1024 default

### flux1-image-edit (nvidia/flux_1-kontext-dev)
- In-context editing
- Brand consistency

### stable-diffusion-35 (nvidia/stable-diffusion-3_5-large)
- Professional quality
- Marketing materials

### qwen-vlm-analysis
- Quality verification
- Brand compliance

---

## 📋 WORKFLOW

### 1. Prompt Engineering
```python
def build_prompt(product, style, mood):
    return f"{product} in {style} style, {mood}, 4K, professional"
```

### 2. Image Generation
```python
images = client.generate_image(
    prompt=build_prompt("watch", "luxury", "elegant"),
    negative="blurry, low quality",
    width=1024, height=1024
)
```

### 3. Quality Check
```python
def verify_quality(image_path):
    # Analyze with Qwen 3.5 VLM
    # Check: resolution, artifacts, composition
    return {"passed": True, "score": 0.95}
```

### 4. Auto-Edit
```python
# Use FLUX.1-kontext for fixes
fixed = client.edit_image(path, "improve lighting")
```

### 5. GitLab Upload
```python
def upload_gitlab(path):
    token = os.getenv("GITLAB_TOKEN")
    url = "https://gitlab.com/api/v4/projects/biometrics%2Fmedia-assets/uploads"
    r = requests.post(url, files={'file': open(path,'rb')},
                     headers={"PRIVATE-TOKEN": token})
    return r.json()["markdown_url"]
```

### 6. URL Storage
```python
def store_supabase(url, metadata):
    supabase = create_client(URL, KEY)
    supabase.table("media_assets").insert({
        "type": "image",
        "public_url": url,
        "metadata": metadata
    }).execute()
```

---

## 🖼️ IMAGE PROCESSING

### ImageMagick Commands

```bash
# Resize
convert input.png -resize 1920x1080 output.png

# Compress JPEG
convert input.png -quality 85 output.jpg

# Add watermark
convert input.png logo.png -gravity southeast -composite output.png

# Format conversion
convert input.png output.webp

# Extract colors
convert input.png -colors 5 -unique-colors palette.png

# Thumbnail
convert input.png -thumbnail 300x300 thumb.png

# Batch resize
mogrify -resize 50% *.png

# Add border
convert input.png -border 10 -bordercolor white output.png

# Optimize PNG
pngquant --quality=65-80 input.png --output output.png
```

### Python Image Processing

```python
from PIL import Image
import os

def resize_keep_aspect(path, max_w=1920, max_h=1080):
    img = Image.open(path)
    img.thumbnail((max_w, max_h), Image.LANCZOS)
    img.save(path)

def add_watermark(path, logo_path, position="se", opacity=0.5):
    base = Image.open(path).convert("RGBA")
    logo = Image.open(logo_path).convert("RGBA")
    logo.alpha = int(255 * opacity)
    
    positions = {"se": (base.width-logo.width-10, base.height-logo.height-10)}
    base.paste(logo, positions[position], logo)
    return base

def create_thumbnails(directory, sizes=[(150,150), (300,300), (600,600)]):
    for f in os.listdir(directory):
        if f.endswith(('.png','.jpg')):
            img = Image.open(f"{directory}/{f}")
            for w,h in sizes:
                img.copy().thumbnail((w,h)).save(f"{directory}/thumb_{w}_{h}_{f}")
```

---

## 📊 QUALITY GATES

| Gate | Criteria | Check |
|------|----------|-------|
| Resolution | >= 1024x1024 | PIL check |
| Artifacts | No compression issues | Edge detection |
| Brand | Colors, logo placement | Template match |
| Format | PNG/JPG/WebP correct | Magic bytes |
| Upload | GitLab accessible | HEAD request |

---

## 🎨 USE CASES

### Product Photography
```python
prompt = "Product on transparent background, 4K, studio lighting, white background"
```

### Marketing Banners
```python
prompt = "Hero banner, luxury watch, dark background, dramatic lighting, 1920x600"
```

### Social Media
```python
# Instagram (1:1)
# Facebook (16:9)
# Twitter (16:9)
# LinkedIn (1.91:1)
```

### Brand Consistency
```python
BRAND_COLORS = ["#FF0000", "#00FF00", "#0000FF"]
BRAND_LOGO_POSITION = "bottom-right"
```

---

## 🔗 INTEGRATIONS

### n8n
```json
{
  "nodes": [
    {"name": "Trigger", "type": "webhook"},
    {"name": "NVIDIA", "type": "httpRequest"},
    {"name": "Supabase", "type": "supabase", "operation": "insert"}
  ]
}
```

### Supabase
```sql
CREATE TABLE images (
    id UUID PRIMARY KEY,
    url TEXT NOT NULL,
    width INT,
    height INT,
    format VARCHAR(10),
    metadata JSONB,
    created_at TIMESTAMPTZ
);
```

---

## 🛠️ TROUBLESHOOTING

| Issue | Cause | Solution |
|-------|-------|----------|
| 400 Bad Request | Prompt too long | <500 chars |
| Low quality | Low steps | Increase to 50 |
| Artifacts | Low guidance | Increase to 10 |
| Wrong colors | Negative prompt | Add "wrong colors" |
| Rate limit | >40/min | Wait 60s |

---

## 📈 BENCHMARKS

| Model | Resolution | Time | Quality |
|-------|------------|------|---------|
| FLUX.1 | 1024x1024 | 10-20s | Excellent |
| SD 3.5 | 1024x1024 | 15-30s | Good |
| FLUX Kontext | 1024x1024 | 20-40s | Excellent |

---

## 🎯 PROMPT ENGINEERING GUIDE

### Structure

```
[Subject] + [Style] + [Lighting] + [Composition] + [Quality] + [Negatives]
```

Beispiel:
```
A luxury watch (Subject)
in minimalist style (Style)
studio lighting, soft shadows (Lighting)
centered composition (Composition)
8K, photorealistic (Quality)
blurry, low quality, distorted (Negatives)
```

### Style Keywords

| Style | Keywords |
|-------|-----------|
| Product | e-commerce, white background, clean, professional |
| Portrait | bokeh, shallow depth, professional lighting |
| Landscape | epic, panoramic, golden hour |
| Cyberpunk | neon, dark, futuristic, RGB |
| Vintage | film grain, warm tones, nostalgic |

### Lighting Setups

```python
LIGHTING_PRESETS = {
    "product": {
        "key": "softbox left",
        "fill": "reflector right",
        "back": "strip light",
        "ambient": "low"
    },
    "portrait": {
        "key": "softbox 45°",
        "fill": "ambient",
        "rim": "hair light",
        "ambient": "medium"
    },
    "moody": {
        "key": "single hard light",
        "fill": "none",
        "back": "none",
        "ambient": "dark"
    }
}
```

### Composition Rules

| Rule | Description | Example |
|------|-------------|----------|
| Rule of Thirds | Grid 3x3 | Subject at intersection |
| Golden Ratio | 1.618 spiral | Natural flow |
| Leading Lines | Draw eye | Roads, fences |
| Symmetry | Mirror image | Architecture |
| Frame | Natural frame | Doors, windows |

---

## 🎨 COLOR THEORY

### Brand Colors

```python
BRAND_PALETTE = {
    "primary": "#FF6B6B",
    "secondary": "#4ECDC4",
    "accent": "#FFE66D",
    "neutral": "#2C3E50",
    "background": "#FFFFFF"
}

def validate_colors(image_path):
    """Check if image matches brand colors"""
    import colorsys
    
    img = Image.open(image_path)
    colors = img.getcolors(maxcolors=1000)
    
    # Extract dominant colors
    dominant = [c[1] for c in sorted(colors, reverse=True)[:5]]
    
    # Check brand match
    brand_colors = [hex_to_rgb(c) for c in BRAND_PALETTE.values()]
    matches = sum(1 for c in dominant if min_color_distance(c, brand_colors) < 50)
    
    return matches >= 2
```

### Color Grading

```python
def apply_color_grade(image_path, preset):
    """Apply color grading preset"""
    
    grades = {
        "warm": {"temperature": 6500, "tint": 10},
        "cool": {"temperature": 5600, "tint": -10},
        "vintage": {"temperature": 4500, "tint": 15, "sepia": 30},
        "dramatic": {"contrast": 1.3, "saturation": 1.2}
    }
    
    img = Image.open(image_path)
    # Apply grade
    return img
```

---

## 🖼️ IMAGE OPTIMIZATION

### Web Optimization

```python
def optimize_for_web(image_path, target_size_kb=100):
    """Optimize image for web"""
    
    img = Image.open(image_path)
    
    # Progressive reduction
    quality = 95
    while True:
        img.save("output.webp", quality=quality, optimize=True)
        size_kb = os.path.getsize("output.webp") / 1024
        
        if size_kb <= target_size_kb or quality < 50:
            break
        quality -= 5
    
    return quality

def generate_srcset(image_path, sizes=[320, 640, 1024, 1920]):
    """Generate srcset for responsive images"""
    
    img = Image.open(image_path)
    base_name = os.path.splitext(image_path)[0]
    
    srcset = []
    for size in sizes:
        # Resize maintaining aspect ratio
        aspect = img.height / img.width
        new_height = int(size * aspect)
        
        resized = img.resize((size, new_height), Image.LANCZOS)
        output_path = f"{base_name}-{size}w.webp"
        resized.save(output_path, quality=85, optimize=True)
        
        srcset.append(f"{output_path} {size}w")
    
    return ", ".join(srcset)
```

### Format Selection

| Format | Use Case | Compression |
|--------|----------|------------|
| JPEG | Photos, complex images | Lossy |
| PNG | Transparency, graphics | Lossless |
| WebP | All web images | Both |
| AVIF | Modern browsers | Best |
| SVG | Logos, icons | Vector |

---

## 🔬 ADVANCED TECHNIQUES

### Inpainting

```python
def inpaint_image(image_path, mask_path, prompt):
    """Inpaint masked areas"""
    
    # Using FLUX.1-kontext
    result = client.edit_image(
        input_image=image_path,
        mask_image=mask_path,
        prompt=prompt,
        strength=0.8
    )
    
    return result["output"]["image_url"]
```

### Outpainting

```python
def extend_image(image_path, direction, prompt):
    """Extend image beyond original boundaries"""
    
    extensions = {
        "left": {"x_offset": -512},
        "right": {"x_offset": 512},
        "top": {"y_offset": -512},
        "bottom": {"y_offset": 512}
    }
    
    result = client.generate_image(
        prompt=f"Extended view of: {prompt}",
        init_image=image_path,
        strength=0.6
    )
    
    return result["output"]["images"][0]["url"]
```

### Style Transfer

```python
def transfer_style(content_image, style_reference):
    """Transfer style from reference to content"""
    
    result = client.edit_image(
        input_image=content_image,
        prompt=f"Style of {style_reference}",
        strength=0.75,
        style="preserve_content"
    )
    
    return result["output"]["image_url"]
```

### Multi-Subject Generation

```python
def generate_multi_subject(prompt, subjects: list):
    """Generate image with multiple specific subjects"""
    
    # Enhanced prompt with subject details
    enhanced = f"{prompt}, featuring: " + ", ".join([
        f"{s['name']}: {s['description']}"
        for s in subjects
    ])
    
    result = client.generate_image(
        prompt=enhanced,
        num_images=4,  # Generate variations
        width=1024,
        height=1024
    )
    
    return result["output"]["images"]
```

### Background Removal

```python
def remove_background(image_path):
    """Remove background using ML"""
    
    # Using rembg or similar
    from rembg import remove
    
    img = Image.open(image_path)
    output = remove(img)
    
    # Save as PNG with alpha
    output.save("output.png")
    
    return output
```

### Face Detection & Enhancement

```python
def enhance_faces(image_path):
    """Detect and enhance faces"""
    
    import cv2
    
    # Load cascade
    face_cascade = cv2.CascadeClassifier(
        cv2.data.haarcascades + 'haarcascade_frontalface_default.xml'
    )
    
    img = cv2.imread(image_path)
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    
    faces = face_cascade.detectMultiScale(gray, 1.1, 4)
    
    # Enhance each face
    for (x, y, w, h) in faces:
        face = img[y:y+h, x:x+w]
        # Apply enhancement
        enhanced = enhance_skin(face)
        img[y:y+h, x:x+w] = enhanced
    
    return Image.fromarray(img)
```

---

## 📐 RESPONSIVE IMAGES

### Picture Element

```html
<picture>
  <source srcset="image-320w.webp 320w,
              image-640w.webp 640w,
              image-1024w.webp 1024w"
          sizes="(max-width: 320px) 100vw,
                 (max-width: 640px) 50vw,
                 33vw"
          type="image/webp">
  <source srcset="image-320w.jpg 320w,
              image-640w.jpg 640w"
          type="image/jpeg">
  <img src="image-640w.jpg" alt="Description">
</picture>
```

### CSS Responsive

```css
.responsive-image {
  max-width: 100%;
  height: auto;
  srcset: "image-320w.webp 320w,
           image-640w.webp 640w,
           image-1024w.webp 1024w";
  sizes: "(max-width: 320px) 100vw,
          (max-width: 640px) 50vw,
          33vw";
}
```

---

## ♿ ACCESSIBILITY

### Alt Text Generation

```python
def generate_alt_text(image_path):
    """Generate accessibility alt text using AI"""
    
    # Analyze image
    analysis = analyze_image(image_path)
    
    # Generate descriptive alt text
    alt = f"{analysis.subject} in {analysis.setting}. "
    alt += f"{analysis.action if analysis.action else ''}. "
    alt += f"Colors: {analysis.primary_colors[:3]}."
    
    return alt
```

### ARIA Labels

```html
<img src="product.jpg"
     alt="Premium wireless headphones"
     role="img"
     aria-label="Product showcase: Premium wireless headphones in black">
```

---

## 📊 ANALYTICS

### Image Performance

```python
def track_image_performance(image_path, metrics):
    """Track image generation metrics"""
    
    data = {
        "filename": os.path.basename(image_path),
        "size_kb": os.path.getsize(image_path) / 1024,
        "dimensions": f"{img.width}x{img.height}",
        "format": img.format,
        "generation_time": metrics["time"],
        "prompt_length": len(metrics["prompt"]),
        "model": metrics["model"]
    }
    
    supabase.table("image_analytics").insert(data).execute()
```

---

## 💰 COST OPTIMIZATION

| Metric | Free | Paid |
|--------|------|------|
| RPM | 40 | 200 |
| Images/day | 100 | Unlimited |
| Max size | 1024x1024 | 2048x2048 |

### Tips

1. **Kleinere Dimensionen** - 512x512 für Thumbnails
2. **Weniger Steps** - 20 statt 50 (geringer Qualitätsverlust)
3. **Wiederverwenden** - Cached Ergebnisse

---

## 🔐 BEST PRACTICES

### Prompt Security

```python
# NIEMALS sensitive info in Prompts!
BAD = "User John Doe's face at 123 Main St"
GOOD = "Professional headshot, business casual"
```

### Rate Limiting

```python
from time import sleep

def rate_limited_call(func, max_per_minute=40):
    calls = []
    
    def wrapper(*args):
        now = time.time()
        calls[:] = [t for t in calls if now - t < 60]
        
        if len(calls) >= max_per_minute:
            sleep(60 - (now - calls[0]))
        
        result = func(*args)
        calls.append(time.time())
        return result
    
    return wrapper
```

---

**Version:** 1.0 | **Updated:** 2026-02-19
