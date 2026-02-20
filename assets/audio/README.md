# Audio Directory

## Overview

This directory contains all audio assets used in the biometrics project. Audio files are used for voice biometrics, notifications, sound effects, and multimedia content.

## Contents

### Voice Samples

| File | Description | Format | Duration |
|------|-------------|--------|----------|
| sample-enroll.mp3 | Enrollment voice sample | MP3 | 30s |
| sample-verify.mp3 | Verification sample | MP3 | 10s |
| sample-commands.mp3 | Command samples | MP3 | 60s |

### Sound Effects

| File | Description | Format |
|------|-------------|--------|
| success.mp3 | Success notification | MP3 |
| error.mp3 | Error notification | MP3 |
| click.mp3 | Button click | WAV |
| notification.mp3 | Alert sound | MP3 |

### Background Audio

| File | Description | Format |
|------|-------------|--------|
| ambient.mp3 | Ambient background | MP3 |
| intro-music.mp3 | Intro jingle | MP3 |

## Audio Standards

### Format Guidelines
| Use Case | Format | Bitrate | Sample Rate |
|----------|--------|---------|-------------|
| Voice Samples | WAV | 256kbps | 44.1kHz |
| Effects | WAV | 128kbps | 44.1Hz |
| Web Audio | MP3 | 128kbps | 44.1kHz |
| Podcast | MP3 | 96kbps | 48kHz |

### Quality Requirements
- **Voice Samples**: High quality, noise-free
- **Effects**: Clean, no artifacts
- **Background**: Looping, seamless

## Voice Biometrics

### Recording Guidelines
- **Environment**: Quiet, no echo
- **Microphone**: High quality, close-talking
- **Format**: WAV, 16-bit, 44.1kHz
- **Volume**: Normalize to -3dB

### Sample Requirements
- Minimum 30 seconds enrollment
- Multiple phrases for variety
- Different acoustic environments

## Usage

### HTML Audio
```html
<audio controls>
  <source src="success.mp3" type="audio/mpeg">
</audio>
```

### JavaScript Audio
```javascript
const audio = new Audio('success.mp3');
audio.play();
```

## Audio Processing

### Normalization
```bash
ffmpeg -i input.wav -af "loudnorm=I=-16:TP=-1.5:LRA=11" output.wav
```

### Format Conversion
```bash
# Convert to MP3
ffmpeg -i input.wav -codec:a libmp3lame -qscale:a 2 output.mp3

# Convert to WAV
ffmpeg -i input.mp3 -codec:a pcm_s16le output.wav
```

### Compression
```bash
# Compress for web
ffmpeg -i input.wav -b:a 128k output.mp3
```

## Storage

### Git LFS
```bash
git lfs track "*.mp3"
git lfs track "*.wav"
git lfs track "*.ogg"
```

### CDN Distribution
```yaml
cdn:
  audio:
    - path: /audio/*
      cache: 30d
      formats: [mp3, ogg]
```

## Accessibility

### Captions
Provide transcripts for audio content:
```html
<audio controls>
  <source src="sample.mp3" type="audio/mpeg">
  <track kind="captions" src="transcript.vtt">
</audio>
```

## Maintenance

- Review audio quality regularly
- Remove unused files
- Update for new branding

## Tools

### Editing
- **Audacity**: Audio editor
- **Adobe Audition**: Pro audio editing
- **ffmpeg**: Command-line processing

### Analysis
- **Waveform**: Visualization
- **Sonic Visualizer**: Analysis

## See Also

- [Videos Directory](../videos/)
- [Icons Directory](../icons/)
- [Images Directory](../images/)
