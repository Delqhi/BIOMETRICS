# Video Outputs Directory

## Overview

This directory contains processed video outputs, including rendered videos, encoded files, and video processing results.

## Contents

### Rendered Videos

| File | Description | Resolution | Duration |
|------|-------------|------------|----------|
| demo-login.mp4 | Login flow demo | 1080p | 30s |
| tutorial-enroll.mp4 | Enrollment tutorial | 720p | 5min |
| feature-highlight.mp4 | Feature showcase | 4K | 60s |

### Processed Videos

| File | Description | Format |
|------|-------------|--------|
| compressed-demo.mp4 | Web-optimized | MP4 |
| preview.gif | Animated preview | GIF |
| thumbnail.jpg | Video thumbnail | JPEG |

### Encoding Outputs

| Directory | Description |
|-----------|-------------|
| web/ | Web-optimized versions |
| mobile/ | Mobile versions |
| archive/ | High-quality archives |

## Video Processing

### Encoding Script
```bash
# Process video for web
ffmpeg -i input.mp4 \
  -c:v libx264 -preset medium \
  -crf 23 -c:a aac \
  -movflags +faststart \
  output-web.mp4
```

### Thumbnail Generation
```bash
# Generate thumbnail
ffmpeg -i video.mp4 -ss 00:00:05 \
  -vframes 1 thumbnail.jpg
```

### GIF Creation
```bash
# Create GIF
ffmpeg -i video.mp4 -ss 10 -t 5 \
  -vf "fps=10,scale=320:-1:flags=lanczos" \
  output.gif
```

## Quality Presets

### Web (High)
- Resolution: 1080p
- Bitrate: 5 Mbps
- Format: MP4 (H.264)

### Web (Standard)
- Resolution: 720p
- Bitrate: 2.5 Mbps
- Format: MP4 (H.264)

### Mobile
- Resolution: 480p
- Bitrate: 1 Mbps
- Format: MP4 (H.264)

## Storage

### Organization
```
videos/
├── 2026-01/
│   ├── demo-login.mp4
│   └── demo-login-web.mp4
└── 2026-02/
    └── tutorial-enroll.mp4
```

### CDN Distribution
```yaml
cdn:
  videos:
    - path: /videos/outputs/*
      cache: 30d
      formats: [mp4, webm]
```

## Cleanup

### Regular Cleanup
```bash
# Clean temporary files
rm -f *.tmp

# Remove old processed files
find . -mtime +30 -delete
```

## Best Practices

1. **Version**: Use semantic naming
2. **Archive**: Keep high-quality masters
3. **Optimize**: Create web versions
4. **Document**: Track processing steps

## Maintenance

- Review storage monthly
- Compress old videos
- Remove duplicates

## See Also

- [Assets Videos](../assets/videos/)
- [Outputs Overview](../outputs/)
- [Assets Output](./assets/)
