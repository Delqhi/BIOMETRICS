# {{PROJECT_NAME}}

> Quick project description goes here

## 📖 Introduction

**What is this project?**
Brief description of what the project does and its main purpose.

**Who is this for?**
Target audience and use cases.

**Why use this project?**
Key benefits and unique selling points.

---

## 🚀 Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/{{ORG}}/{{PROJECT_NAME}}.git
cd {{PROJECT_NAME}}

# 2. Install dependencies
npm install

# 3. Configure environment
cp .env.example .env
# Edit .env with your values

# 4. Start development
npm run dev

# 5. Run tests
npm run test
```

**Expected result:** Application starts on `http://localhost:50001`

---

## 📚 API Reference

Full API documentation available at `/docs/dev/`

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/v1/{{RESOURCE}}` | List {{RESOURCE}} |
| POST | `/api/v1/{{RESOURCE}}` | Create {{RESOURCE}} |
| GET | `/api/v1/{{RESOURCE}}/:id` | Get {{RESOURCE}} by ID |
| PUT | `/api/v1/{{RESOURCE}}/:id` | Update {{RESOURCE}} |
| DELETE | `/api/v1/{{RESOURCE}}/:id` | Delete {{RESOURCE}} |

---

## 📖 Tutorials

- [Getting Started](/docs/non-dev/getting-started.md)
- [Configuration Guide](/docs/dev/configuration.md)
- [Deployment Guide](/docs/dev/deployment.md)

---

## 🔧 Troubleshooting

### Common Issues

**Issue:** Database connection failed
**Solution:** Check `DATABASE_URL` in `.env`

**Issue:** Port already in use
**Solution:** Change `PORT` in `.env`

See [Troubleshooting Guide](/docs/dev/troubleshooting.md) for more.

---

## 📝 Changelog

All notable changes to this project will be documented in [CHANGELOG.md](/CHANGELOG.md).

---

## 📄 License

MIT License - see [LICENSE](/LICENSE) for details.
