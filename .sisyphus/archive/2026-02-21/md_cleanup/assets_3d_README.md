# 3D Assets Directory

## Overview

This directory contains all 3D model assets used in the biometrics project. These assets are used for visualizations, presentations, and product demonstrations.

## Contents

### 3D Models

| File | Format | Description | Usage |
|------|--------|-------------|-------|
| biometric-scanner.obj | OBJ | 3D scanner model | Demos, marketing |
| face-model.obj | OBJ | Face biometric model | Visualizations |
| fingerprint-scanner.glb | GLB | Fingerprint scanner | Product demos |
| iris-scanner.glb | GLB | Iris recognition device | Presentations |
| device-mockup.glb | GLB | Hardware device mockup | Marketing materials |

### Model Specifications

#### Face Model
- **Vertices**: 50,000
- **Textures**: 4K resolution
- **Format**: OBJ with MTL
- **License**: Proprietary

#### Fingerprint Scanner
- **Format**: GLB (glTF binary)
- **Animations**: Idle animation included
- **PBR Materials**: Yes
- **Poly Count**: 10,000 triangles

### Usage Guidelines

#### Rendering
```bash
# Convert to render-ready format
blender --background --python render.py input.obj
```

#### Web Usage
```javascript
// Using Three.js
const loader = new GLTFLoader();
loader.load('face-model.glb', (gltf) => {
  scene.add(gltf.scene);
});
```

## Quality Standards

### Texture Requirements
- **Resolution**: Minimum 2048x2048
- **Format**: PNG with alpha or JPEG
- **Color Space**: sRGB

### Model Requirements
- **Cleanup**: No duplicate vertices
- **Normals**: Recalculated and smoothed
- **Origin**: Centered at origin (0,0,0)
- **Scale**: Real-world units (meters)

## Tools

### Recommended Tools
- Blender 3.x - Modeling and export
- Maya - Complex animations
- Substance Painter - Texturing

### Export Formats
- **OBJ**: Universal format
- **FBX**: Unity/Unreal compatible
- **GLB**: Web/Three.js optimized

## Version Control

Track 3D assets in Git LFS:
```bash
git lfs track "*.obj"
git lfs track "*.glb"
git lfs track "*.fbx"
```

## Maintenance

### Regular Updates
- Review models quarterly
- Update textures for new branding
- Optimize for performance

### Archival
- Keep original source files
- Version tagged releases
- Document changes

## See Also

- [Renders Directory](../renders/)
- [Images Directory](../images/)
- [Logos Directory](../logos/)
