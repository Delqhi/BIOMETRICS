# Media Assets Documentation

**Purpose:** Documentation and guides for media asset management

## Overview

This directory contains documentation for media assets used in the BIOMETRICS project, including guidelines for using, creating, and managing various media types.

## Contents

### Asset Guidelines

| Document | Description |
|----------|-------------|
| [-guidelines.md](asset-guidassetelines.md) | General guidelines for all media |
| [image-standards.md](image-standards.md) | Image specifications |
| [video-production.md](video-production.md) | Video production workflow |

## Supported Formats

### Images
- PNG (preferred for diagrams)
- JPEG (for photographs)
- SVG (for icons and logos)
- WebP (for web content)

### Videos
- MP4 (H.264)
- WebM (VP9)
- MOV (for editing)

### Audio
- MP3
- WAV (for high quality)
- OGG (for web)

## Asset Categories

### Dashboard Screenshots
- Format: PNG
- Resolution: 1920x1080 minimum
- Style: Dark mode preferred

### Architecture Diagrams
- Format: SVG (preferred) or PNG
- Style: Consistent color scheme
- Colors: See brand guidelines

### Video Content
- Resolution: 1080p minimum
- Frame rate: 30fps standard
- Codec: H.264

## Processing Pipeline

### Image Optimization
```bash
# Compress images
biometrics media optimize-images --input ./assets --quality 85

# Generate thumbnails
biometrics media thumbnail --input ./assets/images --size 300x200
```

### Video Processing
```bash
# Transcode video
biometrics media transcode --input video.mov --output video.mp4 --codec h264

# Generate preview
biometrics media preview --input video.mp4 --frames 10
```

## Maintenance

- Review assets quarterly
- Update deprecated formats
- Archive unused assets
- Back up originals

## Related Documentation

- [Brand Guidelines](../brand-guidelines.md)
- [Design System](../design-system.md)
- [Assets Location](../../assets/)
