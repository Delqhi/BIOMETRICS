# 🎙️ ENTERPRISE AUDIO & SYNC PIPELINE (AUDIO-GEN.md)

<document_purpose>
  directive: "EXECUTION MANUAL FOR AUDIO SYNTHESIS, ENHANCEMENT, AND VIDEO SYNCHRONIZATION."
  target_agent: "Technical Director (TD) / Qwen 3.5-397B-A17B"
  mandate: "You MUST strictly follow these pipelines to generate, process, and merge audio. Never output raw, unpolished audio."
</document_purpose>

<model_registry_and_routing>
  voice_generation:
    agent_name: "BrandVoice-Gen"
    models: ["nvidia/magpie-tts-zeroshot", "nvidia/magpie-tts-multilingual"]
    use_case: "Generate voiceovers using a specific brand voice reference file from `/inputs/references/`."
  audio_enhancement:
    agent_name: "Studio-Optimizer"
    models: ["nvidia/studiovoice", "nvidia/background-noise-removal"]
    use_case: "Upscale generated or raw audio to professional studio quality and remove artifacts."
  facial_animation_sync:
    agent_name: "Audio2Face-Sync"
    model: "nvidia/audio2face-3d"
    use_case: "Generate facial animation data from an audio file to animate 3D avatars."
</model_registry_and_routing>

<brand_voice_protocol>
  concept: "The brand must sound identical across all media. No random synthetic voices."
  execution: "When generating speech, ALWAYS use `magpie-tts-zeroshot` and provide the designated reference audio file (e.g., `/inputs/references/brand_voice_sample.wav`) as the zero-shot anchor."
  tonality_check: "Use `<think>` mode to analyze the script before generation. Ensure the pacing, pauses (using SSML if applicable), and tone match the brand identity."
</brand_voice_protocol>

<studio_quality_pipeline>
  mandate: "Never use the raw output from a TTS model directly in a final video."
  workflow:
    step_1_tts: "Generate the raw voiceover via NIM API. Save to `/assets/audio/raw_voice.wav`."
    step_2_enhancement: "Pass `raw_voice.wav` through the `nvidia/studiovoice` API to normalize EQ and add professional studio resonance. Save to `/assets/audio/studio_voice.wav`."
    step_3_music_mixing: "If a background track is provided in `/inputs/brand_assets/`, use a Python script to duck the music volume by -15dB during voiceover segments. Save combined track to `/assets/audio/final_mix.wav`."
</studio_quality_pipeline>

<ffmpeg_synchronization_protocol>
  concept: "Audio and Video must be merged programmatically without human intervention."
  execution: "Write and execute a strict FFmpeg Python script to combine the visual output and the audio mix."
  command_structure: `ffmpeg -i /outputs/generated_video.mp4 -i /assets/audio/final_mix.wav -c:v copy -c:a aac -strict experimental -map 0:v:0 -map 1:a:0 -shortest /outputs/final_master_video.mp4`
  verification: "Verify that the final `.mp4` file exists, has an audio stream, and the duration matches the shortest input."
</ffmpeg_synchronization_protocol>

<technical_execution_protocol>
  api_communication: "All NIM API calls MUST be executed via Python wrappers (`requests`)."
  file_formats: "Always request uncompressed `.wav` from the audio APIs for processing. Only compress to `aac` during the final FFmpeg merge."
  error_handling: "If the audio API fails or times out, do NOT merge a silent video. Stop the pipeline, log the error in `/logs/thinking/audio_error.log`, and retry automatically."
</technical_execution_protocol>
