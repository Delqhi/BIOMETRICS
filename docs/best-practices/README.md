# 📖 BIOMETRICS Documentation - Best Practices

**Mandates, workflows, compliance, and troubleshooting.**

---

## 📁 Best Practices Documents

| File | Description | Priority |
|------|-------------|----------|
| [∞Best∞Practices∞Loop.md](∞Best∞Practices∞Loop.md) | Infinite work loop | 🔴 CRITICAL |
| [BLUEPRINT.md](BLUEPRINT.md) | Blueprint template | 🔴 CRITICAL |
| [COMPLIANCE.md](COMPLIANCE.md) | Compliance requirements | 🔴 CRITICAL |
| [SECURITY.md](SECURITY.md) | Security protocols | 🔴 CRITICAL |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Contribution guide | 🟠 HIGH |
| [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) | Code of conduct | 🟠 HIGH |
| [TESTING.md](TESTING.md) | Testing standards | 🟠 HIGH |
| [TESTING-SUITE.md](TESTING-SUITE.md) | Test suite guide | 🟠 HIGH |
| [TROUBLESHOOTING.md](TROUBLESHOOTING.md) | Common issues | 🟡 MEDIUM |
| [CHANGELOG.md](CHANGELOG.md) | Project changelog | 🟡 MEDIUM |
| [MEETING.md](MEETING.md) | Meeting notes | 🟡 MEDIUM |
| [GREENBOOK.md](GREENBOOK.md) | Green book guide | 🟡 MEDIUM |
| [SETUP-COMPLETE.md](SETUP-COMPLETE.md) | Setup completion | 🟡 MEDIUM |
| [DOCUMENTATION-TEMPLATE.md](DOCUMENTATION-TEMPLATE.md) | Doc template | 🟡 MEDIUM |
| [GDPR-GOVERNANCE.md](GDPR-GOVERNANCE.md) | GDPR compliance | 🟠 HIGH |
| [SECURITY-AUDIT.md](SECURITY-AUDIT.md) | Security audits | 🟠 HIGH |
| [PENETRATION-TESTING.md](PENETRATION-TESTING.md) | Pen testing | 🟠 HIGH |

---

## 🔥 DEQLHI-LOOP (Infinite Work Mode)

**Core Principle:** After EACH completed task → Add 5 new tasks immediately

### Work Rules (ABSOLUT BINDEND):
1. NIEMALS warten auf Agenten → Immer parallel weiterarbeiten
2. NIEMALS delegate_task mit run_in_background=false → Immer background
3. HAUPTSÄCHLICH selbst coden → Nur kritisches delegieren
4. IMMER 5 neue Tasks nach jeder Completion → Todo-Liste nie leer
5. IMMER dokumentieren → Jede Änderung in lastchanges.md + AGENTS.md
6. IMMER visuell prüfen → Screenshots, Browser-Checks, CDP Logs
7. IMMER Crashtests → Keine Annahmen, nur harte Fakten
8. IMMER Best Practices 2026 → CEO-Elite Niveau, nichts Halbfertiges

---

## 📋 33 Core Mandates

### Critical Mandates (Top 5):
1. **IMMUTABILITY OF KNOWLEDGE** - Never delete without backup
2. **MODULAR SWARM SYSTEM** - 5+ agents minimum for complex tasks
3. **REALITY OVER PROTOTYPE** - No mocks, real code only
4. **OMNISCIENCE BLUEPRINT** - 500+ line documentation
5. **DOCKER SOVEREIGNTY** - Local persistence, save images

### Operational Mandates:
- **FREE-FIRST PHILOSOPHY** - Self-host everything possible
- **NO STANDARD PORTS** - Use unique ports (50000-59999)
- **GIT COMMIT DISCIPLINE** - After every change
- **TODO CONTINUATION** - Track all tasks
- **PARALLEL EXECUTION** - Never block, always parallel

---

## 🛡️ Security Best Practices

1. **NEVER commit secrets** - Use environment variables
2. **File permissions** - `chmod 600` for sensitive configs
3. **Zero-trust architecture** - Verify everything
4. **Regular security audits** - Penetration testing
5. **GDPR compliance** - Data protection by design

---

## 🧪 Testing Standards

### Test Coverage Requirements:
- **Unit Tests:** 80%+ coverage
- **Integration Tests:** All API endpoints
- **E2E Tests:** Critical user flows
- **Performance Tests:** Load testing 10k+ users

### Test Types:
- Unit tests (Jest, Vitest)
- Integration tests (Supertest)
- E2E tests (Playwright)
- Performance tests (k6)
- Security tests (OWASP ZAP)

---

## 🔧 Troubleshooting

### Common Issues:

#### 1. Models Not Found
```bash
# Solution: Restart terminal
exec zsh
```

#### 2. Timeout in Config
```bash
# Check (MUST BE EMPTY!)
grep -r "timeout" ~/.config/opencode/opencode.json
```

#### 3. HTTP 429 Rate Limit
```bash
# Solution: Wait 60 seconds + use fallbacks
```

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for complete guide.

---

**Last Updated:** 2026-02-18  
**Status:** ✅ Production-Ready
