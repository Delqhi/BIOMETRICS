# Enterprise AGENTS.md - AI Agent Rules and Policies

**Version:** 1.0.0  
**Status:** ACTIVE  
**Scope:** All AI Agents in Enterprise Project

---

## 1. CORE PRINCIPLES

### 1.1 Agent Philosophy
- **Autonomy First** - Agents work independently without waiting
- **Parallel Execution** - Multiple agents work simultaneously
- **Quality Over Speed** - Production-ready code always
- **Documentation Mandatory** - Every change must be documented

### 1.2 Decision Making
- Agents must verify before executing
- No assumptions - always confirm with tests
- Break complex tasks into parallel subtasks
- Escalate blockers immediately

---

## 2. AGENT TYPES

### 2.1 Architect Agent
**Purpose:** System design and architecture decisions

**Responsibilities:**
- Design scalable architectures
- Define data models
- Plan API structures
- Evaluate trade-offs

**Model:** Qwen 3.5 397B (NVIDIA NIM)

### 2.2 Developer Agent
**Purpose:** Code implementation

**Responsibilities:**
- Write production-ready code
- Follow coding standards
- Implement tests
- Document implementations

**Model:** Qwen 3.5 397B (NVIDIA NIM)

### 2.3 Researcher Agent
**Purpose:** Documentation and research

**Responsibilities:**
- Create documentation
- Research best practices
- Analyze requirements
- Gather technical information

**Model:** MiniMax M2.5 (parallel execution)

### 2.4 Tester Agent
**Purpose:** Quality assurance

**Responsibilities:**
- Write unit tests
- Perform integration testing
- Execute E2E tests
- Validate security

**Model:** Kimi K2.5 (for complex analysis)

### 2.5 Reviewer Agent
**Purpose:** Code review and quality control

**Responsibilities:**
- Review pull requests
- Check code quality
- Validate security
- Ensure standards compliance

**Model:** Qwen 3.5 397B

---

## 3. TEAM COLLABORATION

### 3.1 Swarm Protocol
All complex tasks must use swarm delegation:

```
1. Planner Agent - Breaks down task
2. Researcher Agent - Gathers context
3. Developer Agent - Implements solution
4. Tester Agent - Validates implementation
5. Reviewer Agent - Ensures quality
```

### 3.2 Communication Patterns
- **Direct** - Agent to Agent for clarifications
- **Through Coordinator** - Complex workflows via orchestrator
- **Documentation** - All decisions in project docs

### 3.3 Conflict Resolution
1. Agents disagree on approach
2. Escalate to Reviewer Agent
3. If, still unresolved escalate to Architect
4. Final decision documented

---

## 4. WORKFLOW STANDARDS

### 4.1 Task Execution
```
1. Analyze Task → Break into subtasks
2. Delegate → Parallel agent execution
3. Implement → Each agent works independently
4. Verify → Tests and validation
5. Review → Quality gate
6. Deploy → Production deployment
7. Document → Update all docs
```

### 4.2 Parallel Execution Rules
- Minimum 3 agents for complex tasks
- No agent waits for another
- Results merged after completion
- Failures handled individually

### 4.3 Error Handling
- Agents handle errors independently
- Log all errors with context
- Retry with exponential backoff
- Escalate after 3 failures

---

## 5. SECURITY POLICIES

### 5.1 Secrets Management
- NEVER commit secrets to git
- Use environment variables
- Store in Vault for production
- Rotate credentials regularly

### 5.2 Code Security
- No hardcoded credentials
- Input validation on all endpoints
- Output encoding for XSS prevention
- SQL parameterization required

### 5.3 Access Control
- Minimum privilege principle
- Role-based access control
- Audit all access attempts
- Session timeout enforcement

---

## 6. QUALITY STANDARDS

### 6.1 Code Quality
- TypeScript strict mode
- 80% minimum test coverage
- No critical lint warnings
- Performance benchmarks pass

### 6.2 Documentation Quality
- All exported functions documented
- API docs always up-to-date
- README for every module
- Changelog for releases

### 6.3 Review Quality
- All PRs reviewed by 2+ agents
- Security review for changes
- Performance review for optimizations
- Documentation review for accuracy

---

## 7. DEPLOYMENT POLICIES

### 7.1 Environment Strategy
- Development → Local testing
- Staging → Integration tests
- Production → Full validation

### 7.2 Deployment Rules
- All tests must pass
- Security scan clean
- Performance benchmarks met
- Documentation updated

### 7.3 Rollback Procedure
- Automatic on critical failure
- Manual approval for minor issues
- Post-mortem required
- Fix and redeploy

---

## 8. MONITORING & OBSERVABILITY

### 8.1 Logging Standards
- Structured JSON logging
- Correlation IDs for tracing
- Log levels: ERROR, WARN, INFO, DEBUG
- Retention: 90 days production

### 8.2 Metrics
- Request latency (P50, P95, P99)
- Error rates by type
- System health metrics
- Business KPIs

### 8.3 Alerts
- Critical: PagerDuty
- Warning: Slack
- Info: Email digest

---

## 9. COMPLIANCE REQUIREMENTS

### 9.1 Data Protection
- GDPR compliance for EU data
- Encryption at rest and in transit
- Data retention policies
- Right to deletion support

### 9.2 Audit Trail
- All data access logged
- Configuration changes tracked
- User actions auditable
- 7-year retention

### 9.3 Regular Audits
- Quarterly security review
- Annual compliance audit
- Penetration testing
- Vulnerability scanning

---

## 10. AGENT MODEL ASSIGNMENT

### 10.1 Task-Based Model Selection

| Task Type | Model | Provider | Use Case |
|-----------|-------|----------|----------|
| Code Implementation | Qwen 3.5 397B | NVIDIA NIM | Best code quality |
| Research/Docs | MiniMax M2.5 | MiniMax | Fast, parallel |
| Complex Analysis | Kimi K2.5 | Kimi | Deep reasoning |
| Vision/Images | Kimi K2.5 | Kimi | Multimodal |

### 10.2 Rate Limiting
- Qwen 3.5: 1 concurrent (NVIDIA limit)
- MiniMax: 10 concurrent (parallel work)
- Kimi K2.5: 3 concurrent (balanced)

---

## 11. PERFORMANCE TARGETS

### 11.1 Response Times
- API P95: < 200ms
- Page Load: < 3 seconds
- Database Query: < 100ms

### 11.2 Resource Usage
- CPU: < 70% average
- Memory: < 80% average
- Disk: < 90% maximum

### 11.3 Availability
- SLA: 99.9% uptime
- RTO: 1 hour
- RPO: 15 minutes

---

## 12. ESCALATION PROCEDURES

### 12.1 Blocker Escalation
1. Agent identifies blocker
2. Documents in project issues
3. Escalates to Reviewer
4. Reviewer resolves or escalates

### 12.2 Incident Response
1. Detection → Alert sent
2. Assessment → Severity determined
3. Containment → Impact limited
4. Resolution → Fix implemented
5. Review → Post-mortem

---

## 13. CONTINUOUS IMPROVEMENT

### 13.1 Retrospectives
- Weekly team review
- Identify improvements
- Update processes
- Document lessons learned

### 13.2 Knowledge Transfer
- Document all decisions
- Share learnings
- Update AGENTS.md
- Train new agents

---

## 14. PROHIBITED ACTIONS

### 14.1 Absolute Prohibitions
- Delete production data
- Commit secrets to git
- Bypass security controls
- Ignore test failures

### 14.2 Strong Recommendations
- No single points of failure
- Avoid blocking operations
- No hardcoded configuration
- Never skip documentation

---

## 15. SUCCESS CRITERIA

### 15.1 Project Success
- All features delivered on time
- Quality gates passed
- Documentation complete
- Client satisfied

### 15.2 Process Success
- Zero security incidents
- < 1% error rate
- > 99.9% uptime
- Team productivity maintained

---

**Document Version:** 1.0.0  
**Last Updated:** 2026-02-20  
**Review Cycle:** Monthly  
**Owner:** Enterprise Architecture Team
