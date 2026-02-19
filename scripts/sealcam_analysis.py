#!/usr/bin/env python3
"""
SEALCAM Video Analysis - Qwen 3.5 VLM
Analyzes videos using SealCam Framework (Subject, Environment, Action, Lighting, Camera, Movement)
"""

import requests
import os
import json
from dotenv import load_dotenv

load_dotenv()

def analyze_video(video_path):
    """Analyze video using Qwen 3.5 VLM + SealCam Framework"""
    
    prompt = """
    Analyze this video using SealCam Framework:
    1. Subject: Who/what is the focus?
    2. Environment: Where does it take place?
    3. Action: What happens?
    4. Lighting: Light setup?
    5. Camera: Lens type/setting?
    6. Movement: Angle and camera motion?
    
    Output as structured JSON with scene count and duration.
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
        print("Usage: sealcam_analysis.py <video_path>")
        sys.exit(1)
    
    result = analyze_video(sys.argv[1])
    print(json.dumps(result, indent=2))
