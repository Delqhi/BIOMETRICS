#!/usr/bin/env python3
"""
Video Quality Check - Qwen 3.5 VLM
Verifies video quality: physical correctness, no glitches, brand consistency
"""

import requests
import os
import json
from dotenv import load_dotenv

load_dotenv()

def verify_video(video_path):
    """Qwen 3.5 VLM quality verification"""
    
    prompt = """
    Prüfe dieses Video auf:
    1. Physikalische Korrektheit (Schwerkraft, Licht, Schatten)
    2. Keine Artefakte/Glitches
    3. Konsistente Beleuchtung
    4. Marken-Identität gewahrt?
    
    Wenn FEHLER gefunden:
    → Liste ALLE Fehler auf
    → Empfehle Korrektur mit cosmos-video-edit
    
    Wenn PERFEKT:
    → Bestätige "APPROVED FOR PRODUCTION"
    """
    
    response = requests.post(
        "https://integrate.api.nvidia.com/v1/chat/completions",
        headers={
            "Authorization": f"Bearer {os.getenv('NVIDIA_API_KEY')}",
            "Content-Type": "application/json"
        },
        json={
            "model": "qwen/qwen3.5-397b-a17b",
            "messages": [{
                "role": "user",
                "content": [
                    {"type": "text", "text": prompt},
                    {"type": "file", "url": video_path}
                ]
            }]
        }
    )
    
    return response.json()

if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("Usage: video_quality_check.py <video_path>")
        sys.exit(1)
    
    result = verify_video(sys.argv[1])
    print(json.dumps(result, indent=2))
