# Videos Directory

## Overview

This directory contains all video assets used in the biometrics project. Videos are used for marketing, tutorials, demonstrations, and documentation.

## Contents

### Marketing Videos

| File | Description | Duration | Resolution |
|------|-------------|----------|------------|
| product-intro.mp4 | Product introduction | 60s | 4K |
| feature-showcase.mp4 | Feature highlights | 30s | 1080p |
| testimonial-1.mp4 | Customer testimonial | 45s | 1080p |

### Tutorial Videos

| File | Description | Duration |
|------|-------------|----------|
| getting-started.mp4 | Quick start guide | 5min |
| advanced-setup.mp4 | Advanced configuration | 15min |
| api-tutorial.mp4 | API integration | 10min |

### Demo Videos

| File | Description | Format |
|------|-------------|--------|
| demo-login.mp4 | Login flow demo | MP4 |
| demo-enrollment.mp4 | Enrollment demo | MP4 |
| demo-admin.mp4 | Admin panel demo | MP4 |

## Video Standards

### Format Guidelines
| Use Case | Format | Codec | Bitrate |
|----------|--------|-------|---------|
| Web | MP4 | H.264 | 5-10 Mbps |
| Archive | MOV | ProRes | 100+ Mbps |
| Mobile | MP4 | H.264 | 2-4 Mbps |

### Resolution Standards
| Name | Resolution | Aspect |
|------|------------|--------|
| 4K | 3840x2160 | 16:9 |
| 1080p | 1920x1080 | 16:9 |
| 720p | 1280x720 | 16:9 |
| Vertical | 1080x1920 | 9:16 |
| Square | 1080x1080 | 1:1 |

### Audio Standards
- Format: AAC
- Sample Rate: 48kHz
- Channels: Stereo
- Bitrate: 192kbps

## Usage

### HTML5 Video
```html
<video controls poster="thumbnail.jpg" width="1920">
  <source src="video.mp4" type="video/mp4">
  <track kind="captions" src="captions.vtt" srclang="en">
</video>
```

### Responsive Video
```css
.video-container {
  position: relative;
  padding-bottom: 56.25%; /* 16:9 */
  height: 0;
}
.video-container video {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}
```

## Thumbnails

### Auto-generation
```bash
ffmpeg -i input.mp4 -ss 00:00:05 -vframes 1 thumbnail.jpg
```

### Sizes
| Type | Size |
|------|------|
| Preview | 320x180 |
| Standard | 640x360 |
| Large | 1280x720 |

## Captioning

### VTT Format
```vtt
WEBVTT

00:00:00.000 --> 00:00:05.000
Welcome to Biometrics

00:00:05.500 --> 00:00:10.000
This video will show you how to get started
```

## Encoding

### Recommended Settings
```bash
# H.264 encoding
ffmpeg -i input.mov -c:v libx264 -preset slow \
  -crf 23 -c:a aac -b:a 192k \
  -movflags +faststart output.mp4

# WebM (optional)
ffmpeg -i input.mov -c:v libvpx-vp9 \
  -crf 30 -b:v 0 -c:a libopus output.webm
```

## Storage

### Git LFS
```bash
git lfs track "*.mp4"
git lfs track "*.mov"
```

### CDN Distribution
```yaml
# CDN configuration
cdn:
  videos:
    - path: /videos/*
      cache: 7d
      formats: [mp4, webm]
```

## Maintenance

- Review videos quarterly
- Update for new features
- Compress for web optimization

## See Also

- [Audio Directory](../audio/)
- [Images Directory](../images/)
- [Renders Directory](../renders/)
