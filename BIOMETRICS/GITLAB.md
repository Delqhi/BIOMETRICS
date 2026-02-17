# 📦 GITLAB MEDIA STORAGE - BEST PRACTICES 2026

**Status:** ✅ ACTIVE - MANDATORY FOR ALL MEDIA FILES  
**Effective:** 2026-02-18  
**Scope:** ALL AI Agents, ALL Projects, ALL Media Files

---

## 🚨 CRITICAL MANDATE: NO DIRECT MEDIA IN GITHUB REPOS

**PROBLEM:** GitHub has strict file size limits and cannot display large media files (videos, large images, PDFs) directly in repositories.

**SOLUTION:** ALL media files MUST be uploaded to GitLab and referenced via public URLs.

### ❌ VERBOTEN (FORBIDDEN)
- Large video files (>1MB) in GitHub repos
- Large PDF presentations (>5MB) in GitHub repos
- High-resolution images (>2MB) in GitHub repos
- Audio files in GitHub repos
- Binary assets that bloat repo size

### ✅ PFLICHT (MANDATORY)
- Upload ALL media to GitLab first
- Use public GitLab URLs in README files
- Store URLs in Supabase database tables
- Reference GitLab URLs in documentation
- Keep GitHub repos lean (code + small assets only)

---

## 📋 GITLAB SETUP FOR AGENTS

### Step 1: Create GitLab Project (Automated)

Agents MUST create a dedicated GitLab project for media storage:

```bash
curl --request POST \
  --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  --header "Content-Type: application/json" \
  --data '{"name":"biometrics-videos","description":"BIOMETRICS project video storage","visibility":"public"}' \
  "https://gitlab.com/api/v4/projects"
```

**Response contains:**
- `id`: Project ID (e.g., `79575238`)
- `path_with_namespace`: Full path (e.g., `zukunftsorientierte.energie/biometrics-videos`)
- `web_url`: Public URL (e.g., `https://gitlab.com/zukunftsorientierte.energie/biometrics-videos`)

### Step 2: Upload Media Files

#### For Videos:
```bash
curl --request POST \
  --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  --form "file=@/path/to/video.mp4" \
  "https://gitlab.com/api/v4/projects/$PROJECT_ID/uploads"
```

**Response:**
```json
{
  "id": 903612626,
  "alt": "video",
  "url": "/uploads/d23d181f4278365b97454a3c0602d132/video.mp4",
  "full_path": "/-/project/79575238/uploads/d23d181f4278365b97454a3c0602d132/video.mp4",
  "markdown": "![video](/uploads/d23d181f4278365b97454a3c0602d132/video.mp4)"
}
```

#### For Images:
```bash
curl --request POST \
  --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  --form "file=@/path/to/image.png" \
  "https://gitlab.com/api/v4/projects/$PROJECT_ID/uploads"
```

#### For PDFs:
```bash
curl --request POST \
  --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  --form "file=@/path/to/presentation.pdf" \
  "https://gitlab.com/api/v4/projects/$PROJECT_ID/uploads"
```

#### For Audio Files:
```bash
curl --request POST \
  --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  --form "file=@/path/to/audio.mp3" \
  "https://gitlab.com/api/v4/projects/$PROJECT_ID/uploads"
```

### Step 3: Construct Public URL

**Formula:**
```
https://gitlab.com/{path_with_namespace}/-/project/{project_id}/uploads/{upload_path}
```

**Example:**
```
https://gitlab.com/zukunftsorientierte.energie/biometrics-videos/-/project/79575238/uploads/d23d181f4278365b97454a3c0602d132/video.mp4
```

### Step 4: Use in README

```markdown
### 📹 Video

[![Video Thumbnail](BIOMETRICS/video-thumbnail.png)](https://gitlab.com/zukunftsorientierte.energie/biometrics-videos/-/project/79575238/uploads/d23d181f4278365b97454a3c0602d132/video.mp4)

**[Video abspielen ▶](https://gitlab.com/zukunftsorientierte.energie/biometrics-videos/-/project/79575238/uploads/d23d181f4278365b97454a3c0602d132/video.mp4)**
```

### Step 5: Store in Supabase

```sql
INSERT INTO media_assets (
  title,
  type,
  gitlab_url,
  gitlab_project_id,
  upload_path,
  created_at
) VALUES (
  'BIOMETRICS Explainer Video',
  'video',
  'https://gitlab.com/zukunftsorientierte.energie/biometrics-videos/-/project/79575238/uploads/d23d181f4278365b97454a3c0602d132/video.mp4',
  '79575238',
  '/uploads/d23d181f4278365b97454a3c0602d132/video.mp4',
  NOW()
);
```

---

## 🔧 AGENT WORKFLOW (MANDATORY)

### When Creating Media with NLM CLI:

```bash
# 1. Generate video with NLM
nlm generate video "BIOMETRICS explainer" --output video.mp4

# 2. Upload to GitLab IMMEDIATELY
curl --request POST \
  --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  --form "file=@video.mp4" \
  "https://gitlab.com/api/v4/projects/$PROJECT_ID/uploads"

# 3. Extract public URL from response
# 4. Update README with GitLab URL (NOT local path)
# 5. Commit ONLY the thumbnail to GitHub (not the video!)
# 6. Store URL in Supabase media_assets table
```

### File Size Guidelines:

| File Type | Max GitHub Size | GitLab Required |
|-----------|----------------|-----------------|
| Videos | >1MB | ✅ ALWAYS |
| PDFs | >5MB | ✅ ALWAYS |
| Images (PNG/JPG) | >2MB | ✅ RECOMMENDED |
| Images (thumbnails) | <500KB | ❌ Can stay in GitHub |
| Audio | >1MB | ✅ ALWAYS |
| Binary assets | >1MB | ✅ ALWAYS |

---

## 📁 GITLAB PROJECT STRUCTURE

Recommended organization for media projects:

```
biometrics-videos/
├── videos/
│   ├── explainer.mp4
│   ├── tutorial-01.mp4
│   └── demo.mp4
├── presentations/
│   ├── pitch-deck.pdf
│   └── technical-overview.pdf
├── images/
│   ├── infografik.png
│   ├── screenshots/
│   └── thumbnails/
├── audio/
│   ├── podcast-01.mp3
│   └── voiceover.mp3
└── README.md (list of all assets with URLs)
```

---

## 🔐 SECURITY & ACCESS

### Visibility Settings:

- **Public Projects:** Anyone can view/download (default for media)
- **Internal Projects:** Only organization members
- **Private Projects:** Only project members

### API Token Permissions:

Required scopes for media upload:
- `api` - Full API access
- `write_repository` - Upload files
- `read_api` - Read project info

### Token Storage:

```bash
# Store in environment (NEVER commit to git)
export GITLAB_TOKEN="glpat-xxxxxxxxxxxxxxxxxxxx"

# Or use .env file (gitignored)
echo "GITLAB_TOKEN=glpat-xxxxxxxxxxxxxxxxxxxx" >> .env
```

---

## 🚀 AUTOMATED UPLOAD SCRIPT

Save as `upload-to-gitlab.sh`:

```bash
#!/bin/bash

# Configuration
GITLAB_TOKEN="${GITLAB_TOKEN:-}"
PROJECT_ID="${GITLAB_PROJECT_ID:-}"
GITLAB_URL="https://gitlab.com"

if [ -z "$GITLAB_TOKEN" ] || [ -z "$PROJECT_ID" ]; then
  echo "Error: GITLAB_TOKEN and GITLAB_PROJECT_ID must be set"
  exit 1
fi

# Upload file
upload_file() {
  local file="$1"
  local description="$2"
  
  if [ ! -f "$file" ]; then
    echo "Error: File not found: $file"
    return 1
  fi
  
  echo "Uploading: $file..."
  
  response=$(curl --request POST \
    --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
    --form "file=@$file" \
    --form "description=$description" \
    "$GITLAB_URL/api/v4/projects/$PROJECT_ID/uploads")
  
  # Extract URL from response
  url=$(echo "$response" | python3 -c "import sys, json; print(json.load(sys.stdin).get('full_path', ''))")
  
  if [ -n "$url" ]; then
    full_url="$GITLAB_URL$url"
    echo "✅ Success! Public URL: $full_url"
    
    # Append to assets list
    echo "- $description: $full_url" >> ASSETS.md
    
    return 0
  else
    echo "❌ Upload failed"
    echo "$response"
    return 1
  fi
}

# Main
if [ $# -lt 1 ]; then
  echo "Usage: $0 <file> [description]"
  exit 1
fi

upload_file "$1" "${2:-$(basename $1)}"
```

**Usage:**
```bash
chmod +x upload-to-gitlab.sh
./upload-to-gitlab.sh video.mp4 "BIOMETRICS Explainer Video"
```

---

## 📊 TRACKING & METRICS

### Supabase Table Schema:

```sql
CREATE TABLE media_assets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('video', 'image', 'pdf', 'audio')),
  gitlab_url TEXT NOT NULL,
  gitlab_project_id TEXT,
  upload_path TEXT,
  file_size BIGINT,
  mime_type TEXT,
  thumbnail_url TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  metadata JSONB DEFAULT '{}'
);

-- Index for fast lookups
CREATE INDEX idx_media_assets_type ON media_assets(type);
CREATE INDEX idx_media_assets_created ON media_assets(created_at DESC);
```

### Example Queries:

```sql
-- Get all videos
SELECT title, gitlab_url, created_at 
FROM media_assets 
WHERE type = 'video' 
ORDER BY created_at DESC;

-- Get total storage used
SELECT type, SUM(file_size) as total_bytes 
FROM media_assets 
GROUP BY type;

-- Find assets by project
SELECT * FROM media_assets 
WHERE gitlab_project_id = '79575238';
```

---

## 🔄 MIGRATION FROM GITHUB

If you have media files already in GitHub:

```bash
# 1. Download from GitHub
git lfs pull  # If using LFS
# OR
wget https://github.com/user/repo/raw/main/video.mp4

# 2. Upload to GitLab
./upload-to-gitlab.sh video.mp4 "Migrated from GitHub"

# 3. Update all references in code/docs
# Replace: https://github.com/user/repo/raw/main/video.mp4
# With: https://gitlab.com/.../video.mp4

# 4. Remove from GitHub
git rm video.mp4
git commit -m "chore: migrate video to GitLab storage"
git push
```

---

## ✅ VERIFICATION CHECKLIST

Before considering a media file "done":

- [ ] File uploaded to GitLab successfully
- [ ] Public URL accessible in browser (test it!)
- [ ] URL added to README/documentation
- [ ] URL stored in Supabase media_assets table
- [ ] File NOT committed to GitHub (check .gitignore)
- [ ] Thumbnail created and committed to GitHub (if applicable)
- [ ] Metadata recorded (file size, type, upload date)

---

## 🚫 COMMON MISTAKES

| Mistake | Why Wrong | Correct Approach |
|---------|-----------|------------------|
| Commit video to GitHub | Bloats repo, GitHub blocks large files | Upload to GitLab, commit thumbnail only |
| Use raw GitHub URLs | Links break, files blocked | Use GitLab public URLs |
| Forget to test URL | URL might be private | Always open URL in incognito browser |
| No thumbnail for videos | Ugly README preview | Generate thumbnail with ffmpeg |
| Hardcode project ID | Can't reuse across projects | Use environment variables |

---

## 📚 REFERENCES

- **GitLab Upload API:** https://docs.gitlab.com/ee/api/project_level_variables.html
- **GitLab LFS:** https://docs.gitlab.com/ee/user/project/lfs/
- **GitLab Releases:** https://docs.gitlab.com/ee/user/project/releases/
- **Supabase Storage:** https://supabase.com/docs/guides/storage

---

## 🆘 TROUBLESHOOTING

### "403 Forbidden" on upload
- Check API token has `write_repository` scope
- Verify project ID is correct
- Ensure project visibility is public or you have access

### "File too large" error
- GitLab has 2GB per file limit (free tier)
- Compress video with ffmpeg: `ffmpeg -i input.mp4 -vcodec libx264 -crf 28 output.mp4`

### URL not accessible publicly
- Check project visibility: Settings → General → Visibility → Public
- Verify upload completed successfully

### Token authentication failed
- Regenerate token: https://gitlab.com/-/profile/personal_access_tokens
- Export token: `export GITLAB_TOKEN="glpat-..."`

---

**MANDATE VERSION:** 1.0  
**LAST UPDATED:** 2026-02-18  
**COMPLIANCE:** MANDATORY FOR ALL AGENTS  
**ENFORCEMENT:** Automatic via CI/CD checks

---

*"GitLab first, GitHub second - Keep repos lean, URLs clean."*
