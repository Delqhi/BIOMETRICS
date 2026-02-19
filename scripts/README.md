# Automation Scripts

**Purpose:** Utility scripts for automation, deployment, and system management

## Overview

This directory contains various shell and Python scripts used for automation, deployment, configuration, and system management tasks throughout the BIOMETRICS project.

## Scripts

### Deployment Scripts

| Script | Purpose |
|--------|---------|
| `blue-green-deploy.sh` | Blue-green deployment automation |
| `setup.sh` | Initial project setup and configuration |

### Configuration Scripts

| Script | Purpose |
|--------|---------|
| `validate-config.sh` | Validate configuration files |
| `upload-to-gitlab.sh` | Upload artifacts to GitLab |

### Python Utilities

| Script | Purpose |
|--------|---------|
| `cosmos_video_gen.py` | Video generation utilities |
| `nim_engine.py` | NVIDIA NIM engine integration |
| `sealcam_analysis.py` | Seal detection analysis |
| `video_quality_check.py` | Video quality validation |

## Usage

### Shell Scripts

#### Setup Script
```bash
cd /Users/jeremy/dev/BIOMETRICS
./scripts/setup.sh
```

#### Deployment
```bash
./scripts/blue-green-deploy.sh --env production
```

#### Configuration Validation
```bash
./scripts/validate-config.sh --config config.yaml
```

### Python Scripts

#### Basic Usage
```bash
python3 scripts/cosmos_video_gen.py --input video.mp4
python3 scripts/nim_engine.py --model <model-name>
python3 scripts/sealcam_analysis.py --image seal.jpg
python3 scripts/video_quality_check.py --file video.mp4
```

## Requirements

### System Requirements
- Bash 4.0+
- Python 3.8+
- curl, wget

### Python Dependencies

Install via:
```bash
pip install -r requirements.txt
```

Or via pnpm:
```bash
pnpm install --script-prefix=""
```

## Common Tasks

### Running Full Setup
```bash
# Complete project setup
./scripts/setup.sh --full

# Quick setup (minimal)
./scripts/setup.sh --quick
```

### Deployment Workflow
```bash
# Validate configuration first
./scripts/validate-config.sh

# Deploy to staging
./scripts/blue-green-deploy.sh --env staging

# Deploy to production
./scripts/blue-green-deploy.sh --env production
```

### Video Processing
```bash
# Generate video with custom settings
python3 scripts/cosmos_video_gen.py --input input.mp4 --output output.mp4 --quality high

# Check quality
python3 scripts/video_quality_check.py --file output.mp4 --report
```

## Maintenance

### Adding New Scripts

1. Use appropriate shebang: `#!/usr/bin/env bash` or `#!/usr/bin/env python3`
2. Add executable permission: `chmod +x script.sh`
3. Document in this README
4. Add to .gitignore if contains secrets

### Script Standards

- Include usage help: `--help` or `-h`
- Use `set -e` for error handling in bash
- Log with timestamps
- Return appropriate exit codes

## Troubleshooting

### Permission Denied
```bash
chmod +x scripts/*.sh
```

### Missing Dependencies
```bash
pip install -r requirements.txt
```

### Path Issues
Run from project root:
```bash
cd /Users/jeremy/dev/BIOMETRICS
./scripts/script-name.sh
```

## Related Documentation

- [Deployment Guide](../docs/deployment.md)
- [Configuration Reference](../docs/configuration.md)
- [Development Setup](../docs/development-setup.md)
- [CI/CD Pipeline](../.github/workflows/)
