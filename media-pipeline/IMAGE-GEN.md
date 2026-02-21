# 🖼️ ENTERPRISE IMAGE PIPELINE (IMAGE-GEN.md)

<document_purpose>
  directive: "EXECUTION MANUAL FOR STATIC IMAGE GENERATION AND IN-CONTEXT EDITING."
  target_agent: "Technical Director (TD) / Qwen 3.5-397B-A17B"
  mandate: "You MUST strictly follow these pipelines to generate, verify, and edit high-fidelity images. Brand consistency is the highest priority."
</document_purpose>

<model_registry_and_routing>
  photorealistic_master:
    agent_name: "FLUX1-IMG"
    model: "black-forest-labs/flux_1-dev"
    endpoint: "https://build.nvidia.com/black-forest-labs/flux_1-dev/modelcard"
    use_case: "State-of-the-Art generation for photorealistic, highly detailed master images. Default choice for new visual assets."
  in_context_editing:
    agent_name: "FLUX1-IMG-Edit"
    model: "black-forest-labs/flux_1-kontext-dev"
    endpoint: "https://build.nvidia.com/black-forest-labs/flux_1-kontext-dev/modelcard"
    use_case: "Modify existing images, integrate products (`/inputs/brand_assets`) into new environments, or fix errors while strictly maintaining the original visual style and context."
  professional_visuals:
    agent_name: "Diffusion-Image"
    model: "stabilityai/stable-diffusion-3_5-large"
    endpoint: "https://build.nvidia.com/stabilityai/stable-diffusion-3_5-large"
    use_case: "Alternative generation for specific stylized visual assets or when FLUX constraints apply."
</model_registry_and_routing>

<brand_consistency_protocol>
  concept: "Images must never look like random AI generations. They must fit the corporate identity."
  execution: "When asked to place a product in a new scene, DO NOT use standard text-to-image. You MUST use `FLUX1-IMG-Edit` (In-Context) and provide the original brand asset as the visual anchor."
  prompting_rules: "Always specify exact lighting, camera angle, and material properties in the prompt to match the established brand guide."
</brand_consistency_protocol>

<quality_assurance_and_refinement_loop>
  concept: "Zero-Defect Delivery. Do not output the first generated image blindly."
  workflow:
    step_1_generate: "Generate the initial image via Python API call to `FLUX1-IMG` or `Diffusion-Image`. Save to `/assets/renders/temp_v1.jpg`."
    step_2_vision_inspection: "Load the generated image back into your Vision capabilities (Qwen 3.5 VLM). Analyze it for:"
      - "Artifacts (e.g., deformed hands, melting objects)"
      - "Text errors (Gibberish text on signs/products)"
      - "Brand compliance (Wrong colors, incorrect product details)"
    step_3_correction: "If ANY defect is found, generate an exact bounding-box or mask description of the error. Call `FLUX1-IMG-Edit` to surgically fix ONLY the broken area."
    step_4_finalization: "Repeat inspection (Max 3 iterations). Once perfect, move to `/outputs/final_image.jpg`."
</quality_assurance_and_refinement_loop>

<technical_execution_protocol>
  api_communication: "All calls MUST be executed via Python wrappers (`requests`)."
  payload_handling: "Images sent to the API for editing MUST be base64 encoded accurately. Ensure aspect ratios are preserved."
  logging: "Log the exact prompts, the flaws found during Vision Inspection, and the refinement steps in `/logs/thinking/image_qa_log.json`."
</technical_execution_protocol>