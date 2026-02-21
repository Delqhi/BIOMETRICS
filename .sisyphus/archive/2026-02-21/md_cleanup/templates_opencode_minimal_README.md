# OpenCode Minimal Template

Quick prototype template for small projects.

## Quick Start

```bash
# 1. Install dependencies
npm install

# 2. Copy .env and add your NVIDIA_API_KEY
cp .env.example .env

# 3. Run OpenCode
npx opencode
```

## Minimal Setup

- **Provider:** NVIDIA NIM (Qwen 3.5)
- **API Key:** Only `NVIDIA_API_KEY` required
- **Language:** TypeScript

## Project Structure

```
├── src/
│   └── index.ts       # Entry point
├── .env.example       # Environment template
├── opencode.json      # OpenCode config
├── package.json       # Dependencies
└── README.md          # This file
```

## Usage

Edit `src/index.ts` to add your tasks, then run:

```bash
npx opencode
```
