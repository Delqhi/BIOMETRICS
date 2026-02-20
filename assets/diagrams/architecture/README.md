# Architecture Diagrams

## Overview

This directory contains detailed architecture diagrams that illustrate the system design, component relationships, and data flows for the biometrics platform.

## Diagram Categories

### High-Level Architecture

| File | Description |
|------|-------------|
| high-level.svg | Overall system architecture |
| microservices.svg | Microservices breakdown |
| integration.svg | External integrations |

### Component Architecture

| File | Description |
|------|-------------|
| auth-service.svg | Authentication service design |
| biometric-engine.svg | Biometric processing engine |
| storage-service.svg | Data storage architecture |
| api-gateway.svg | API gateway configuration |

### Data Architecture

| File | Description |
|------|-------------|
| data-model.svg | Data model relationships |
| schema.svg | Database schema |
| cache-strategy.svg | Caching architecture |

## Diagram Standards

### Visual Style
- **Style**: Clean, modern, minimal
- **Colors**: Corporate palette (blue/slate)
- **Lines**: Consistent stroke width (2pt)
- **Spacing**: 20pt grid alignment

### Components

| Component Type | Symbol |
|-----------------|--------|
| Service | Rounded rectangle |
| Database | Cylinder |
| Queue | Parallel lines |
| API | Circle |
| User | Person icon |
| External | Dashed border |

### Connectors

| Relationship | Line Style |
|--------------|------------|
| Data flow | Solid arrow |
| Async message | Dashed arrow |
| Bidirectional | Double arrow |
| External | Globe icon |

## Detailed Diagrams

### System Overview Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        Clients                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐       │
│  │   Web   │  │ Mobile  │  │   API   │  │  CLI    │       │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘       │
└───────┼────────────┼────────────┼────────────┼─────────────┘
        │            │            │            │
        └────────────┴────────────┴────────────┘
                          │
                    ┌─────▼─────┐
                    │   CDN     │
                    └─────┬─────┘
                          │
                  ┌───────▼────────┐
                  │  Load Balancer  │
                  └───────┬────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
  ┌─────▼─────┐    ┌─────▼─────┐    ┌─────▼─────┐
  │ API GW 1  │    │ API GW 2  │    │ API GW 3  │
  └─────┬─────┘    └─────┬─────┘    └─────┬─────┘
        │                 │                 │
        └─────────────────┼─────────────────┘
                          │
                    ┌─────▼─────┐
                    │  Services │
                    └───────────┘
```

### Service Dependencies

```
┌─────────────┐     ┌─────────────┐
│    Auth     │────►│   Users     │
└─────────────┘     └─────────────┘
       │                   │
       ▼                   ▼
┌─────────────┐     ┌─────────────┐
│ Biometrics  │────►│  Analytics  │
└─────────────┘     └─────────────┘
       │
       ▼
┌─────────────┐     ┌─────────────┐
│  Storage    │◄────│   Queue     │
└─────────────┘     └─────────────┘
```

## Tooling

### Primary Tool
- **Draw.io**: Main diagramming tool
- **Files**: `.drawio` format (editable)

### Export Options
```bash
# Export all diagrams
for f in *.drawio; do
  drawio export --format svg --output "${f%.drawio}.svg" "$f"
done
```

## Maintenance

### Update Triggers
- New feature releases
- Architecture changes
- Quarterly reviews

### Version Control
```bash
git lfs track "*.drawio"
git lfs track "*.svg"
```

## Integration

### Documentation Usage
```markdown
## System Architecture

![High Level Architecture](architecture/high-level.svg)

### Component Details

![Auth Service](architecture/auth-service.svg)
```

## Best Practices

1. **Keep Updated**: Update with every architecture change
2. **Version Control**: Track source files
3. **Consistency**: Use standardized symbols
4. **Clarity**: Keep diagrams readable

## See Also

- [Diagrams Overview](../diagrams/)
- [3D Assets](../3d/)
- [Renders](../renders/)
