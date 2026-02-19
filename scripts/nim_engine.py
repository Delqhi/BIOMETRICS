#!/usr/bin/env python3
"""
NVIDIA NIM Engine - Central API Wrapper
Orchestrates all NVIDIA NIM model calls (Cosmos, FLUX, TRELLIS, etc.)
"""

import requests
import os
import json
from dotenv import load_dotenv

load_dotenv()

class NIMEngine:
    def __init__(self):
        self.api_key = os.getenv('NVIDIA_API_KEY')
        self.base_url = 'https://integrate.api.nvidia.com/v1'
        self.headers = {
            'Authorization': f'Bearer {self.api_key}',
            'Content-Type': 'application/json'
        }
    
    def generate_video(self, model, prompt, duration=5, fps=30):
        """Generate video with Cosmos models"""
        response = requests.post(
            f'{self.base_url}/video/generate',
            headers=self.headers,
            json={
                'model': model,
                'prompt': prompt,
                'duration': duration,
                'fps': fps,
                'resolution': '1920x1080'
            }
        )
        return response.json()
    
    def edit_video(self, model, input_video, prompt):
        """Edit/refine video with Cosmos-Predict"""
        response = requests.post(
            f'{self.base_url}/video/predict',
            headers=self.headers,
            json={
                'model': model,
                'input_video': input_video,
                'prompt': prompt
            }
        )
        return response.json()
    
    def generate_image(self, model, prompt):
        """Generate image with FLUX or Stable Diffusion"""
        response = requests.post(
            f'{self.base_url}/images/generate',
            headers=self.headers,
            json={
                'model': model,
                'prompt': prompt
            }
        )
        return response.json()

if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument('--model', required=True)
    parser.add_argument('--prompt', required=True)
    parser.add_argument('--output', default='output.mp4')
    args = parser.parse_args()
    
    engine = NIMEngine()
    
    if 'cosmos' in args.model:
        result = engine.generate_video(args.model, args.prompt)
        print(f"Video generated: {result}")
    elif 'flux' in args.model or 'stable' in args.model:
        result = engine.generate_image(args.model, args.prompt)
        print(f"Image generated: {result}")
