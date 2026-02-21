# Video Transcoding

## Overview
Video transcoding pipeline for BIOMETRICS video content.

## Features
- Multi-format encoding (MP4, WebM)
- Adaptive bitrate streaming (HLS)
- Thumbnail generation
- Video preview generation

## Transcoding Pipeline

### Input Formats
- MP4, MOV, AVI, MKV, WebM

### Output Formats
| Format | Resolution | Bitrate | Use Case |
|--------|-----------|---------|----------|
| HLS | 1080p | 5Mbps | Adaptive streaming |
| MP4 | 1080p | 3Mbps | Download |
| MP4 | 720p | 1.5Mbps | Mobile |
| WebM | 480p | 0.8Mbps | Web preview |

### API Endpoints
```
POST /api/videos/transcode
Body: { videoUrl, outputs[] }

GET /api/videos/:id/status
GET /api/videos/:id/
```

## Processing
- FFmpegstream for encoding
- AWS Elemental MediaConvert alternative
- Queue via Supabase (background jobs)

## Storage
- Input: Supabase Storage `/videos/input/`
- Output: Supabase Storage `/videos/output/`
- Thumbnails: `/videos/thumbnails/`
