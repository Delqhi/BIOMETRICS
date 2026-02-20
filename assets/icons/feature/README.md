# Feature Icons

## Overview

This directory contains feature-specific icons that represent the core biometric features and capabilities of the platform.

## Contents

### Biometric Types

| File | Feature | Description |
|------|---------|-------------|
| face-recognition.svg | Face Recognition | Facial biometric |
| fingerprint.svg | Fingerprint | Fingerprint scan |
| voice-recognition.svg | Voice ID | Voice biometric |
| iris-recognition.svg | Iris Scan | Iris/eye scan |
| palm-vein.svg | Palm Vein | Palm vein pattern |

### Security Features

| File | Feature | Description |
|------|---------|-------------|
| liveness-detection.svg | Liveness | Spoofing prevention |
| multi-factor.svg | Multi-Factor | MFA support |
| encryption.svg | Encryption | Data encryption |
| secure-enclave.svg | Secure Enclave | Hardware security |

### Platform Features

| File | Feature | Description |
|------|---------|-------------|
| mobile-sdk.svg | Mobile SDK | iOS/Android support |
| web-auth.svg | WebAuthn | Web authentication |
| api-integration.svg | API | REST API |
| sdk.svg | SDK | Software development kit |

## Icon Design

### Style
- **Type**: Filled icons
- **Detail**: High detail for recognition
- **Size**: 48px for feature displays
- **Color**: Brand blue (#2563EB)

### Representation
- Face: Simplified face outline
- Fingerprint: Finger pad pattern
- Voice: Sound wave
- Iris: Eye with concentric circles

## Usage

### Feature Pages
```html
<div class="feature">
  <svg class="feature-icon">...</svg>
  <h3>Face Recognition</h3>
  <p>Advanced facial biometric authentication</p>
</div>
```

### Comparison Tables
```html
<table>
  <thead>
    <tr>
      <th>Feature</th>
      <th>Support</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><svg class="icon">...</svg> Face ID</td>
      <td>✓</td>
    </tr>
  </tbody>
</table>
```

### Marketing Materials
```html
<div class="feature-grid">
  <div class="feature-item">
    <img src="face-recognition.svg" alt="Face Recognition">
  </div>
  <div class="feature-item">
    <img src="fingerprint.svg" alt="Fingerprint">
  </div>
</div>
```

## Color Usage

### Default State
```css
.feature-icon {
  fill: #2563EB;
}
```

### Active State
```css
.feature-icon.active {
  fill: #10B981;
}
```

### Inactive State
```css
.feature-icon.inactive {
  fill: #CBD5E1;
}
```

## Accessibility

### Alt Text
```html
<img src="face-recognition.svg" 
     alt="Face Recognition: Advanced facial biometric authentication">
```

## Maintenance

- Review feature coverage
- Update for new features
- Remove deprecated features

## Best Practices

1. Use consistent sizing across features
2. Maintain visual hierarchy
3. Provide clear descriptions
4. Support all screen readers

## See Also

- [Icons Overview](../icons/)
- [Social Icons](./social/)
- [Action Icons](./action/)
