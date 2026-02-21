# NLM-ASSETS README

This directory contains all NotebookLM-generated content assets for the BIOMETRICS project.

## Purpose

The NLM-ASSETS directory serves as the central repository for AI-generated multimedia content that supplements the project's documentation. All video, infographic, presentation, report, table, mindmap, and podcast assets are stored here following a strict organizational structure.

## Directory Structure

```
NLM-ASSETS/
├── README.md                 # This file
├── videos/                   # Video content and recordings
├── infographics/            # Visual information graphics
├── presentations/           # Slide decks and presentations
├── reports/                # Generated reports and analyses
├── tables/                 # Data tables and spreadsheets
├── mindmaps/               # Concept visualizations
└── podcasts/               # Audio content and podcasts
```

## Asset Types

### videos/
Educational videos, product demonstrations, and tutorial content. All videos should be in MP4 format with maximum 1080p resolution for web optimization.

**Naming Convention:** `{project}-{type}-{date}.mp4`
**Example:** `biometrics-demo-onboarding-2026-02-18.mp4`

### infographics/
Visual representations of data and concepts. Preferred formats are PNG for static and SVG for scalable graphics.

**Naming Convention:** `{project}-{topic}-{version}.{ext}`
**Example:** `biometrics-architecture-overview-v2.png`

### presentations/
Executive decks, pitch presentations, and meeting materials. Stored in PDF format for universal compatibility.

**Naming Convention:** `{project}-{audience}-{date}.pdf`
**Example:** `biometrics-investor-deck-2026-02.pdf`

### reports/
AI-generated analysis reports, status updates, and documentation. Markdown and PDF formats supported.

**Naming Convention:** `{project}-{report-type}-{date}.md`
**Example:** `biometrics-q1-analysis-2026-02.md`

### tables/
Structured data exports, comparison matrices, and KPI tables. CSV and Markdown formats preferred.

**Naming Convention:** `{project}-{table-type}-{date}.csv`
**Example:** `biometrics-kpi-dashboard-2026-02.csv`

### mindmaps/
Conceptual visualizations showing relationships between ideas. PNG and SVG formats supported.

**Naming Convention:** `{project}-{concept}-{version}.png`
**Example:** `biometrics-system-components-v1.png`

### podcasts/
Audio content including interviews, discussions, and updates. MP3 format with 128kbps minimum quality.

**Naming Convention:** `{project}-{episode}-{date}.mp3`
**Example:** `biometrics-podcast-episode-01-2026-02-18.mp3`

## Upload Process (NLM-CLI)

All assets must be uploaded to NotebookLM using the official NLM-CLI tool. Follow these steps:

### Prerequisites

```bash
# Install NLM-CLI if not already installed
npm install -g @notebooklm/cli

# Authenticate with your account
nlm auth login
```

### Upload Workflow

```bash
# Navigate to the project directory
cd /Users/jeremy/dev/BIOMETRICS/BIOMETRICS

# List existing sources to avoid duplicates
nlm source list

# If duplicate exists, delete old version first
nlm source delete <source-id> -y

# Upload new asset
nlm source add <notebook-id> --file "NLM-ASSETS/videos/biometrics-demo.mp4" --wait

# Verify upload
nlm source list <notebook-id>
```

### Best Practices

1. **Check for Duplicates:** Always run `nlm source list` before uploading to prevent duplicate entries
2. **Delete Old Versions:** Use `nlm source delete` before adding new versions of the same content
3. **Use Wait Flag:** Always use `--wait` flag to ensure upload completes before proceeding
4. **Verify After Upload:** Check source list to confirm successful upload

## Dateiformatkonvention

All dates in filenames must follow ISO 8601 format: `YYYY-MM-DD`

This ensures chronological sorting and prevents ambiguity.

## Versionierung

Version suffixes use lowercase `v` followed by a number: `v1`, `v2`, `v3`

Major changes increment the number, minor updates may use decimal notation: `v1.1`, `v1.2`

## Metadaten-Standards

Each asset should be accompanied by metadata including:

- Creation date
- Source prompt or input
- NLM quality score
- Purpose and target audience
- Associated project phase

## Qualitätskriterien

All NLM-generated content must meet the following minimum standards:

- **Completeness:** All sections or frames present
- **Accuracy:** Factual correctness verified
- **Clarity:** Clear and understandable
- **Target Fit:** Appropriate for intended audience
- **Technical Quality:** No rendering or playback issues

Minimum quality score: 13/16 on NLM quality matrix with 2/2 in correctness.

## Wartung

- Quarterly review of all assets for relevance
- Annual archive of outdated content
- Immediate removal of deprecated material
- Regular backup to GitLab media storage

## Siehe auch

- `NOTEBOOKLM.md` - NLM integration and usage guidelines
- `../∞Best∞Practices∞Loop.md` - Prompt templates and best practices
- `AGENTS-PLAN.md` - Task definitions for asset creation

---

**Letzte Aktualisierung:** 2026-02-18  
**Verantwortlich:** BIOMETRICS Orchestrator  
**Status:** ACTIVE
