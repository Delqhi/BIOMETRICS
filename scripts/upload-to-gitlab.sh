#!/bin/bash
# Upload video/image/audio to GitLab (MANDATORY for media > 1MB)

set -e

if [ $# -lt 1 ]; then
    echo "Usage: upload-to-gitlab.sh <file> [project_id]"
    exit 1
fi

FILE="$1"
PROJECT_ID="${2:-$GITLAB_MEDIA_PROJECT_ID}"
GITLAB_TOKEN="$GITLAB_TOKEN"

echo "📤 Uploading $FILE to GitLab..."

response=$(curl --request POST \
  --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  --form "file=@$FILE" \
  "https://gitlab.com/api/v4/projects/$PROJECT_ID/uploads")

PUBLIC_URL=$(echo $response | jq -r '.full_path')
FULL_URL="https://gitlab.com$PUBLIC_URL"

echo "✅ Uploaded to: $FULL_URL"

# Store in Supabase
if command -v psql &> /dev/null; then
    FILENAME=$(basename "$FILE")
    FILETYPE=$(file -b --mime-type "$FILE")
    
    psql -c "INSERT INTO media_assets (title, type, gitlab_url, created_at) 
             VALUES ('$FILENAME', '$FILETYPE', '$FULL_URL', NOW());"
    echo "✅ Stored in Supabase database"
fi

echo "$FULL_URL"
