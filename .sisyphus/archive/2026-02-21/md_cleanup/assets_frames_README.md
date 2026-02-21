# Frames Directory

## Overview

This directory contains frame assets used for creating visual presentations, marketing materials, and demo environments. Frames provide consistent borders and backgrounds for visual content.

## Contents

### Device Frames

| File | Description | Dimensions |
|------|-------------|------------|
| phone-frame.png | Mobile phone frame | 1242x2688 |
| tablet-frame.png | Tablet frame | 2048x2732 |
| laptop-frame.png | Laptop frame | 2880x1800 |
| desktop-frame.png | Desktop frame | 3840x2160 |

### Browser Frames

| File | Description | Dimensions |
|------|-------------|------------|
| browser-chrome.png | Chrome browser | 1920x1080 |
| browser-safari.png | Safari browser | 1920x1080 |
| browser-firefox.png | Firefox browser | 1920x1080 |

### Presentation Frames

| File | Description | Style |
|------|-------------|-------|
| slide-16-9.svg | 16:9 slide frame | Vector |
| slide-4-3.svg | 4:3 slide frame | Vector |
| title-card.svg | Title card frame | Vector |

## Frame Types

### Device Mockup Frames
- Realistic device bezels
- Shadow and reflection effects
- Multiple angle options
- Color variants

### Browser Mockup Frames
- Consistent browser chrome
- Address bar customizable
- Multiple viewport sizes
- Dark/light themes

### Content Frames
- Clean borders
- Title areas
- Consistent padding
- Brand colors

## Usage

### Image Placement
```bash
# Place screenshot in device frame
./scripts/frame-add.sh screenshot.png phone-frame.png output.png
```

### Automated Generation
```python
from PIL import Image

def add_frame(content, frame):
    # Overlay content in frame
    frame.paste(content, (x, y))
    return frame
```

## Templates

### Quick Frame Template
```html
<div class="device-frame phone">
  <div class="screen">
    <img src="screenshot.png" alt="App screenshot">
  </div>
</div>
```

### CSS Styles
```css
.device-frame {
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
}
```

## Export Options

### Image Formats
- PNG (with transparency)
- JPEG (for presentations)
- WebP (for web)

### Resolution Presets
| Preset | Resolution | Use Case |
|--------|------------|----------|
| Thumbnail | 300px width | Cards |
| Standard | 1200px width | Blog |
| High | 2400px width | Print |

## Maintenance

- Update device frames annually
- Add new device types as needed
- Maintain consistent styling

## Tools

### Frame Creation
- **Figma**: Main design tool
- **Blender**: 3D device frames
- **Keyshot**: Photorealistic renders

### Image Compositing
```bash
# Composite command
convert screenshot.png phone-frame.png -gravity center -composite output.png
```

## Version Control

```bash
git lfs track "*.png"
git lfs track "*.psd"
```

## See Also

- [Images Directory](../images/)
- [Renders Directory](../renders/)
- [Dashboard Components](../dashboard/)
