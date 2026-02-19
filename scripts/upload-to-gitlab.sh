#!/bin/bash

# =============================================================================
# upload-to-gitlab.sh - Upload Media Files to GitLab and Store URL in Supabase
# =============================================================================
# 
# Description:
#   This script uploads media files larger than 1MB to GitLab and stores
#   the public URL in a Supabase media_assets table.
#
# Usage:
#   ./upload-to-gitlab.sh <file_path> <media_type>
#
# Arguments:
#   file_path   - Path to the file to upload
#   media_type  - Type of media: video, image, audio, or model3d
#
# Environment Variables (must be set):
#   GITLAB_TOKEN            - GitLab personal access token
#   GIT_PROJECT_ID          - GitLab project ID for media storage
#   GITLAB_API_URL          - GitLab API URL (default: https://gitlab.com/api/v4)
#   SUPABASE_URL            - Supabase project URL
#   SUPABASE_KEY            - Supabase anon/public key
#   SUPABASE_TABLE          - Table name for media assets (default: media_assets)
#
# Features:
#   - Only uploads files > 1MB (smaller files skipped)
#   - Supports video, image, audio, and 3D model files
#   - Generates public URL via GitLab raw file URL
#   - Stores metadata in Supabase
#   - Retry logic for failed uploads (3 attempts)
#   - Comprehensive logging
# =============================================================================

set -o pipefail

# =============================================================================
# Configuration
# =============================================================================

# Minimum file size to upload (1MB in bytes)
MIN_FILE_SIZE=1048576

# Maximum retry attempts
MAX_RETRIES=3

# Retry delay in seconds
RETRY_DELAY=2

# Log file location
LOG_DIR="${LOG_DIR:-/Users/jeremy/dev/BIOMETRICS/logs}"
LOG_FILE="${LOG_DIR}/upload-to-gitlab.log"

# GitLab API defaults
GITLAB_API_URL="${GITLAB_API_URL:-https://gitlab.com/api/v4}"

# Supabase defaults
SUPABASE_TABLE="${SUPABASE_TABLE:-media_assets}"

# =============================================================================
# Color Codes for Output
# =============================================================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# =============================================================================
# Logging Functions
# =============================================================================

log() {
    local level="$1"
    local message="$2"
    local timestamp
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    # Ensure log directory exists
    mkdir -p "$LOG_DIR"
    
    # Log to file with timestamp
    echo "[$timestamp] [$level] $message" >> "$LOG_FILE"
    
    # Output to stdout with color
    case "$level" in
        ERROR)
            echo -e "${RED}[$timestamp] ERROR: $message${NC}"
            ;;
        SUCCESS)
            echo -e "${GREEN}[$timestamp] SUCCESS: $message${NC}"
            ;;
        WARNING)
            echo -e "${YELLOW}[$timestamp] WARNING: $message${NC}"
            ;;
        INFO)
            echo -e "${BLUE}[$timestamp] INFO: $message${NC}"
            ;;
        *)
            echo "[$timestamp] $message"
            ;;
    esac
}

log_upload() {
    local file_path="$1"
    local url="$2"
    local status="$3"
    
    log "$status" "Upload: $file_path -> $url"
}

# =============================================================================
# Validation Functions
# =============================================================================

validate_environment() {
    local missing_vars=()
    
    # Check required environment variables
    if [[ -z "$GITLAB_TOKEN" ]]; then
        missing_vars+=("GITLAB_TOKEN")
    fi
    
    if [[ -z "$GIT_PROJECT_ID" ]]; then
        missing_vars+=("GIT_PROJECT_ID")
    fi
    
    if [[ -z "$SUPABASE_URL" ]]; then
        missing_vars+=("SUPABASE_URL")
    fi
    
    if [[ -z "$SUPABASE_KEY" ]]; then
        missing_vars+=("SUPABASE_KEY")
    fi
    
    if [[ ${#missing_vars[@]} -gt 0 ]]; then
        log ERROR "Missing required environment variables: ${missing_vars[*]}"
        return 1
    fi
    
    log INFO "Environment validation passed"
    return 0
}

validate_file() {
    local file_path="$1"
    
    # Check if file exists
    if [[ ! -f "$file_path" ]]; then
        log ERROR "File does not exist: $file_path"
        return 1
    fi
    
    # Check if file is readable
    if [[ ! -r "$file_path" ]]; then
        log ERROR "File is not readable: $file_path"
        return 1
    fi
    
    # Get file size
    local file_size
    file_size=$(stat -f%z "$file_path" 2>/dev/null || stat -c%s "$file_path" 2>/dev/null)
    
    if [[ $file_size -lt $MIN_FILE_SIZE ]]; then
        log INFO "File is smaller than 1MB ($file_size bytes), skipping upload"
        return 2
    fi
    
    log INFO "File size: $file_size bytes (eligible for upload)"
    return 0
}

validate_media_type() {
    local media_type="$1"
    
    case "$media_type" in
        video|image|audio|model3d)
            return 0
            ;;
        *)
            log ERROR "Invalid media type: $media_type (must be: video, image, audio, model3d)"
            return 1
            ;;
    esac
}

# =============================================================================
# GitLab Upload Functions
# =============================================================================

upload_file_to_gitlab() {
    local file_path="$1"
    local attempt=1
    local response
    local http_code
    
    # Get filename from path
    local filename
    filename=$(basename "$file_path")
    
    # Sanitize filename for GitLab
    local sanitized_filename
    sanitized_filename=$(echo "$filename" | sed 's/[^a-zA-Z0-9._-]/_/g')
    
    while [[ $attempt -le $MAX_RETRIES ]]; do
        log INFO "Upload attempt $attempt of $MAX_RETRIES: $file_path"
        
        # Upload file to GitLab using multipart form upload
        response=$(curl -s -w "\n%{http_code}" \
            --request POST \
            --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
            --form "file=@${file_path};type=application/octet-stream" \
            --form "branch=main" \
            "${GITLAB_API_URL}/projects/${GIT_PROJECT_ID}/repository/files/${sanitized_filename}" \
            2>&1)
        
        # Extract HTTP status code (last line)
        http_code=$(echo "$response" | tail -n1)
        
        # Get response body (all but last line)
        response=$(echo "$response" | sed '$d')
        
        if [[ "$http_code" == "200" || "$http_code" == "201" ]]; then
            log SUCCESS "File uploaded successfully to GitLab"
            echo "$response"
            return 0
        elif [[ "$http_code" == "400" && "$response" == *"already exists"* ]]; then
            # File already exists, try to update it
            log INFO "File already exists, attempting to update..."
            response=$(curl -s -w "\n%{http_code}" \
                --request PUT \
                --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
                --form "file=@${file_path};type=application/octet-stream" \
                --form "branch=main" \
                --form "commit_message=Update ${sanitized_filename}" \
                "${GITLAB_API_URL}/projects/${GIT_PROJECT_ID}/repository/files/${sanitized_filename}" \
                2>&1)
            
            http_code=$(echo "$response" | tail -n1)
            response=$(echo "$response" | sed '$d')
            
            if [[ "$http_code" == "200" || "$http_code" == "201" ]]; then
                log SUCCESS "File updated successfully in GitLab"
                echo "$response"
                return 0
            fi
        fi
        
        log WARNING "Upload attempt $attempt failed (HTTP $http_code): $response"
        
        if [[ $attempt -lt $MAX_RETRIES ]]; then
            log INFO "Waiting ${RETRY_DELAY}s before retry..."
            sleep $RETRY_DELAY
        fi
        
        ((attempt++))
    done
    
    log ERROR "All $MAX_RETRIES upload attempts failed"
    return 1
}

get_public_url() {
    local file_path="$1"
    
    # Get filename from path
    local filename
    filename=$(basename "$file_path")
    
    # Sanitize filename for URL
    local sanitized_filename
    sanitized_filename=$(echo "$filename" | sed 's/[^a-zA-Z0-9._-]/_/g')
    
    # URL encode the filename
    local encoded_filename
    encoded_filename=$(python3 -c "import urllib.parse; print(urllib.parse.quote('$sanitized_filename'))" 2>/dev/null || echo "$sanitized_filename")
    
    # Construct raw file URL
    # Note: This assumes public access is enabled on the repository
    local project_path_encoded
    project_path_encoded=$(python3 -c "import urllib.parse; print(urllib.parse.quote('$GIT_PROJECT_ID'))" 2>/dev/null || echo "$GIT_PROJECT_ID")
    
    local public_url="https://gitlab.com/api/v4/projects/${project_path_encoded}/repository/files/${encoded_filename}/raw?ref=main"
    
    log INFO "Generated public URL: $public_url"
    echo "$public_url"
    return 0
}

# =============================================================================
# Supabase Functions
# =============================================================================

save_to_supabase() {
    local url="$1"
    local media_type="$2"
    local file_path="$3"
    
    # Get file metadata
    local filename
    filename=$(basename "$file_path")
    
    local file_size
    file_size=$(stat -f%z "$file_path" 2>/dev/null || stat -c%s "$file_path" 2>/dev/null)
    
    local file_extension
    file_extension="${filename##*.}"
    
    local mime_type
    case "$file_extension" in
        mp4|webm|avi|mov|mkv)
            mime_type="video/${file_extension}"
            ;;
        jpg|jpeg|png|gif|webp|svg|bmp)
            mime_type="image/${file_extension}"
            ;;
        mp3|wav|ogg|flac|aac|m4a)
            mime_type="audio/${file_extension}"
            ;;
        obj|fbx|stl|gltf|glb|3ds|blend)
            mime_type="model/${file_extension}"
            ;;
        *)
            mime_type="application/octet-stream"
            ;;
    esac
    
    # Get file checksum
    local checksum
    checksum=$(shasum -a 256 "$file_path" 2>/dev/null | awk '{print $1}')
    
    # Prepare JSON payload
    local json_payload
    json_payload=$(cat <<EOF
{
    "url": "$url",
    "media_type": "$media_type",
    "filename": "$filename",
    "file_size": $file_size,
    "mime_type": "$mime_type",
    "checksum": "$checksum",
    "uploaded_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
    "gitlab_project_id": "$GIT_PROJECT_ID"
}
EOF
)
    
    log INFO "Saving to Supabase: $SUPABASE_TABLE"
    
    # Insert into Supabase
    local response
    local http_code
    
    response=$(curl -s -w "\n%{http_code}" \
        --request POST \
        --header "apikey: $SUPABASE_KEY" \
        --header "Authorization: Bearer $SUPABASE_KEY" \
        --header "Content-Type: application/json" \
        --header "Prefer: return=representation" \
        --data "$json_payload" \
        "${SUPABASE_URL}/rest/v1/${SUPABASE_TABLE}" \
        2>&1)
    
    http_code=$(echo "$response" | tail -n1)
    response=$(echo "$response" | sed '$d')
    
    if [[ "$http_code" == "200" || "$http_code" == "201" ]]; then
        log SUCCESS "Saved to Supabase successfully"
        return 0
    else
        log ERROR "Failed to save to Supabase (HTTP $http_code): $response"
        return 1
    fi
}

# =============================================================================
# Main Functions
# =============================================================================

main() {
    local file_path="$1"
    local media_type="$2"
    
    # Validate arguments
    if [[ -z "$file_path" || -z "$media_type" ]]; then
        echo "Usage: $0 <file_path> <media_type>"
        echo ""
        echo "Arguments:"
        echo "  file_path   - Path to the file to upload"
        echo "  media_type  - Type of media: video, image, audio, or model3d"
        echo ""
        echo "Environment Variables (must be set):"
        echo "  GITLAB_TOKEN           - GitLab personal access token"
        echo "  GIT_PROJECT_ID        - GitLab project ID for media storage"
        echo "  GITLAB_API_URL        - GitLab API URL (optional, default: https://gitlab.com/api/v4)"
        echo "  SUPABASE_URL          - Supabase project URL"
        echo "  SUPABASE_KEY          - Supabase anon/public key"
        echo "  SUPABASE_TABLE        - Table name (optional, default: media_assets)"
        echo ""
        echo "Example:"
        echo "  GITLAB_TOKEN=xxx GIT_PROJECT_ID=123 SUPABASE_URL=https://xxx.supabase.co SUPABASE_KEY=xxx \\"
        echo "    $0 /path/to/video.mp4 video"
        exit 1
    fi
    
    log INFO "=========================================="
    log INFO "Starting upload: $file_path (type: $media_type)"
    log INFO "=========================================="
    
    # Validate environment
    if ! validate_environment; then
        log_upload "$file_path" "" "FAILED"
        exit 1
    fi
    
    # Validate media type
    if ! validate_media_type "$media_type"; then
        log_upload "$file_path" "" "FAILED"
        exit 1
    fi
    
    # Validate file
    if ! validate_file "$file_path"; then
        # Return code 2 means file is too small, not an error
        if [[ $? -eq 2 ]]; then
            exit 0
        fi
        log_upload "$file_path" "" "FAILED"
        exit 1
    fi
    
    # Upload to GitLab
    if ! upload_file_to_gitlab "$file_path"; then
        log_upload "$file_path" "" "FAILED"
        exit 1
    fi
    
    # Get public URL
    local public_url
    if ! public_url=$(get_public_url "$file_path"); then
        log_upload "$file_path" "" "FAILED"
        exit 1
    fi
    
    # Save to Supabase
    if ! save_to_supabase "$public_url" "$media_type" "$file_path"; then
        log_upload "$file_path" "$public_url" "PARTIAL"
        log WARNING "File uploaded to GitLab but failed to save to Supabase"
        exit 1
    fi
    
    log_upload "$file_path" "$public_url" "SUCCESS"
    log SUCCESS "=========================================="
    log SUCCESS "Upload completed successfully!"
    log SUCCESS "URL: $public_url"
    log SUCCESS "=========================================="
    
    exit 0
}

# =============================================================================
# Script Entry Point
# =============================================================================

main "$@"
