# 🎥 ENTERPRISE VIDEO & 3D PIPELINE (VIDEO-GEN.md)

<document_purpose>
  directive: "EXECUTION MANUAL FOR VIDEO, 3D ASSETS, AND FRAME ANIMATION."
  target_agent: "Technical Director (TD) / Qwen 3.5-397B-A17B"
  mandate: "You MUST strictly follow these pipelines for any moving image or 3D generation. Do not deviate from the assigned NVIDIA NIM models."
</document_purpose>

<model_registry_and_routing>
  3d_generation:
    model: "microsoft/trellis"
    endpoint: "https://build.nvidia.com/microsoft/trellis/modelcard"
    use_case: "Generate high-quality 3D assets (.glb/.usd) from text or 2D brand images (`/inputs/brand_assets`)."
  video_from_text:
    model: "nvidia/cosmos-transfer1-7b"
    endpoint: "https://build.nvidia.com/nvidia/cosmos-transfer1-7b/modelcard"
    use_case: "Generate physics-aware video world states strictly from text prompts."
  video_from_image_or_video:
    model: "nvidia/cosmos-predict1-5b"
    endpoint: "https://build.nvidia.com/nvidia/cosmos-predict1-5b/modelcard"
    use_case: "Animate a static image (e.g., a 3D render) or extend an existing video with cinematic, physics-compliant movement."
</model_registry_and_routing>

<sealcam_analysis_framework>
  concept: "Never prompt blindly. If referencing an existing video (`/inputs/references`), you MUST first analyze it using the SealCam framework."
  execution: "Write a Python script to pass the reference video to your Vision capabilities, extracting the following strict JSON schema into `/logs/thinking/`:"
  schema:
    Subject: "Who/what is the main focus? (e.g., 'A matte-silver premium headphone')"
    Environment: "Where does it take place? (e.g., 'Dark, neon-lit cyberpunk street')"
    Action: "What exactly happens? (e.g., 'Headphone floats and disassembles')"
    Lighting: "Light setup? (e.g., 'High-contrast rim lighting, top-down spotlight')"
    Camera: "Lens/Shot type? (e.g., '50mm macro, shallow depth of field')"
    Movement: "Angle and movement? (e.g., 'Slow orbit tracking shot, low angle')"
  application: "Use this exact JSON to construct the new prompt for Cosmos, swapping the 'Subject' with our specific brand asset."
</sealcam_analysis_framework>

<pipeline_3d_to_web_animation>
  mandate: "When generating interactive 3D scroll animations (The 'Apple Effect'), execute this exact 4-step sequence:"
  step_1_trellis:
    action: "Call Microsoft TRELLIS via NIM API to convert `/inputs/brand_assets/product.png` into a 3D asset. Save to `/assets/3d/`."
  step_2_master_frame:
    action: "Write a script to render the 3D model into a high-res starting frame (PNG). Save to `/assets/renders/`."
  step_3_cosmos_animation:
    action: "Feed the master frame into `cosmos-predict1-5b`. Prompt for a smooth, physics-correct transformation (e.g., explosion view or 360 spin). Save output to `/outputs/`."
  step_4_ffmpeg_extraction:
    action: "Execute strict FFmpeg CLI command to extract exactly 30 FPS for web scrolling."
    command: `ffmpeg -i /outputs/generated_video.mp4 -vf "fps=30,scale=1920:-1" -q:v 2 /assets/frames/frame_%04d.jpg`
    verification: "Count the generated frames. Ensure naming is sequential (frame_0001.jpg, etc.)."
</pipeline_3d_to_web_animation>

<physics_and_quality_correction>
  trigger: "Visual artifacts, unnatural gravity, or lighting inconsistencies in Cosmos output."
  action: 
    - "Activate `<think>` mode."
    - "Diagnose the exact physics failure (e.g., 'Shadows move in opposite direction to light source')."
    - "Rewrite the Cosmos prompt to explicitly enforce physical constraints (e.g., 'Strictly maintain static top-down lighting, lock shadow vectors')."
    - "Re-run the NIM API call automatically (Max 3 retries)."
</physics_and_quality_correction>