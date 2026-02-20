# AGENTS-PLAN.md - 20-TASK INFINITY LOOP

**Created:** 2026-02-20  
**Status:** ✅ ACTIVE  
**Version:** 1.0 "DEQLHI-LOOP EDITION"  
**Reference:** `∞Best∞Practices∞Loop.md` (Operating System for Agents)

---

## 🎯 CORE PRINCIPLE: INFINITE WORK LOOP

> **"Ein Task endet, fünf neue beginnen"**  
> **"Kein Warten, nur Arbeiten"**  
> **"Kein Fertig, nur Weiter"**

After EVERY completed task → IMMEDIATELY add 5 new tasks  
NEVER stop → PERPETUAL LOOP until all 100 done → Then 100 new tasks

---

## 📊 INFINITY LOOP STRUCTURE

### Loop Mechanics

```
START (20 Tasks)
  ↓
Task 1 Complete → +5 New Tasks = 24 Tasks
  ↓
Task 2 Complete → +5 New Tasks = 28 Tasks
  ↓
Task 10 Complete → +5 New Tasks = 60 Tasks
  ↓
Task 20 Complete → +5 New Tasks = 100 Tasks ← Goal!
  ↓
Task 100 Complete → +5 New Tasks = 105 Tasks ← Loop continues!
```

### Task Categories

1. **Chaos Cleanup** (Phase 2) - ✅ COMPLETE
2. **24/7 Agent Loop** (Phase 3) - ⏳ IN PROGRESS
3. **Onboarding** (Phase 4) - ⏳ PENDING
4. **Templates** (Phase 5) - ⏳ PENDING
5. **Quality Assurance** (Phase 6) - ⏳ PENDING

---

## 🔄 CURRENT LOOP STATUS

### Phase 2: CHAOS CLEANUP ✅ COMPLETE

**Completed Tasks:**
- ✅ Phase 2.1: Archive Sprint 5 packages (14 packages + 3 empty dirs)
- ✅ Phase 2.2: Consolidate MD files (20+ → 6 essential files)
- ✅ Phase 2.3: Update AGENTS-PLAN.md (this file)

**Impact:**
- 70% reduction in root MD files
- 61% reduction in pkg/ packages
- 100% empty directories eliminated
- 246 files changed, 24,490 insertions, 3,999 deletions

**Commit:** `e769ef5` - "chore: chaos cleanup phase 2.1-2.2 complete"

---

### Phase 3: 24/7 AGENT LOOP ⏳ IN PROGRESS

**Goal:** Build autonomous orchestrator that works 24/7 without user intervention

**Tasks:**
- ⏳ Phase 3.1: Create orchestrator.go with session monitoring
- ⏳ Phase 3.2: Implement 'Sicher?' check logic
- ⏳ Phase 3.3: Implement massive prompt generator for sub-agents
- ⏳ Phase 3.4: Test orchestrator with 3 parallel agents (different models)
- ⏳ Phase 3.5: Deploy orchestrator as Docker container

**Model Assignment Rules:**
- Qwen 3.5 397B: MAX 1 agent
- Kimi K2.5: MAX 1 agent
- MiniMax M2.5: MAX 1 agent
- **TOTAL: MAX 3 agents parallel (all different models)**

---

### Phase 4: ONBOARDING ⏳ PENDING

**Goal:** One-click setup for new users

**Tasks:**
- ⏳ Phase 4.1: Create ONBOARDING.md (complete setup guide)
- ⏳ Phase 4.2: Create bootstrap.sh (automated setup script)
- ⏳ Phase 4.3: Create API key setup wizard
- ⏳ Phase 4.4: Test onboarding with fresh clone
- ⏳ Phase 4.5: Create video tutorial (5 min)

---

### Phase 5: TEMPLATES ⏳ PENDING

**Goal:** Copy-paste templates for new projects

**Tasks:**
- ⏳ Phase 5.1: OpenCode Standard Template (DONE - verify)
- ⏳ Phase 5.2: OpenCode Minimal Template (DONE - verify)
- ⏳ Phase 5.3: OpenClaw Standard Template
- ⏳ Phase 5.4: OpenClaw Enterprise Template
- ⏳ Phase 5.5: Test all templates with fresh projects

---

### Phase 6: QUALITY ASSURANCE ⏳ PENDING

**Goal:** Tesla/Apple level quality, no "school project" code

**Tasks:**
- ⏳ Phase 6.1: Review all created files against ∞Best∞Practices∞Loop.md
- ⏳ Phase 6.2: Crashtest orchestrator (edge cases, failures)
- ⏳ Phase 6.3: Performance benchmarks (<1s response time)
- ⏳ Phase 6.4: Security audit (OWASP Top 10)
- ⏳ Phase 6.5: User acceptance testing

---

## 🎯 TASK GENERATION RULES

### After EVERY Task Completion:

1. **Mark task complete** in this file
2. **Add 5 new tasks** immediately
3. **Commit changes** with conventional commit
4. **Continue to next task** (no pause!)

### Task Format:

```markdown
- [ ] Phase X.Y: Task name
  - **Goal:** Clear objective
  - **Success Criteria:** Measurable outcome
  - **Estimated Time:** X hours
  - **Priority:** high/medium/low
```

---

## 🚨 CRITICAL RULES (NEVER BREAK)

### 1. Model Assignment

❌ **VERBOTEN:**
- 2 agents with same model parallel
- Qwen 3.5 running on >1 agent
- Ignoring model limits

✅ **PFLICHT:**
- MAX 3 agents parallel
- Each agent = different model
- Verify before spawning

### 2. Evidence-Based Completion

❌ **VERBOTEN:**
- Claiming "fertig" without files
- Lying about tests
- No verification

✅ **PFLICHT:**
- Show file contents
- Run tests
- "Sicher?" check
- Git commit

### 3. Read-First Policy

❌ **VERBOTEN:**
- Creating files without checking existence
- Overwriting without reading
- Duplikate erstellen

✅ **PFLICHT:**
- glob() before create
- Read completely (bis letzte Zeile!)
- Reuse existing structure

### 4. Active Orchestration

❌ **VERBOTEN:**
- Waiting passively for agents
- Not monitoring sessions
- Ignoring stuck agents

✅ **PFLICHT:**
- Monitor all sessions
- Intervene if stuck
- Verify work actively
- Ask "Sicher?" if unsure

---

## 📊 PROGRESS TRACKING

### Metrics

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| Tasks Completed | ∞ | 3 | 🔄 Starting |
| Tasks Added | ∞ | 15 | 🔄 Active |
| Chaos Reduction | 90% | 70% | ✅ Good |
| Agent Loop | 24/7 | 0% | ⏳ Pending |
| Onboarding Time | <5 min | N/A | ⏳ Pending |

### Velocity

- **Tasks/Hour:** Target 5+ (currently ~3)
- **Commits/Hour:** Target 3+ (currently ~2)
- **Code Lines/Hour:** Target 100+ (currently ~50)

---

## 🔄 LOOP MAINTENANCE

### Weekly Review (Every Sunday)

1. **Archive completed tasks** (older than 7 days)
2. **Consolidate duplicates**
3. **Update priorities**
4. **Add 20 new tasks** (if queue < 50)

### Monthly Cleanup

1. **Review all active tasks**
2. **Remove obsolete tasks**
3. **Update AGENTS-PLAN.md structure**
4. **Commit cleanup**

---

## 📝 SESSION LOG

### Session 2026-02-20 (Current)

**Started:** 21:12 CET  
**Agent:** Orchestrator (Main)  
**Focus:** Phase 2 Chaos Cleanup

**Accomplished:**
- ✅ Phase 2.1: Archived 14 Sprint 5 packages
- ✅ Phase 2.2: Consolidated root MD files (20+ → 6)
- ✅ Phase 2.3: Created this AGENTS-PLAN.md
- ✅ Git commit: e769ef5 (246 files changed)

**Next:** Phase 3.1 - Create orchestrator.go

---

## 🎯 NEXT 5 TASKS (Ready to Start)

1. **Phase 3.1:** Create orchestrator.go with session monitoring
   - **Goal:** Main orchestrator that spawns/monitors agents
   - **Success:** Running binary that can spawn 3 agents
   - **Time:** 2-3 hours
   - **Priority:** HIGH

2. **Phase 3.2:** Implement 'Sicher?' check logic
   - **Goal:** Verify agent work before marking complete
   - **Success:** Automated verification system
   - **Time:** 1-2 hours
   - **Priority:** HIGH

3. **Phase 3.3:** Massive prompt generator
   - **Goal:** Generate context-rich prompts for sub-agents
   - **Success:** Prompts include ALL context files
   - **Time:** 2 hours
   - **Priority:** HIGH

4. **Phase 3.4:** Test with 3 parallel agents
   - **Goal:** Verify model assignment works
   - **Success:** 3 agents, different models, no conflicts
   - **Time:** 1 hour
   - **Priority:** HIGH

5. **Phase 3.5:** Docker deployment
   - **Goal:** Run orchestrator 24/7 in container
   - **Success:** Container running, spawning agents
   - **Time:** 2 hours
   - **Priority:** MEDIUM

---

## 🔗 REFERENCES

- **∞Best∞Practices∞Loop.md:** Operating System for Agents (archived)
- **docs/ORCHESTRATOR-MANDATE.md:** Orchestrator workflow
- **docs/agents/AGENT-MODEL-MAPPING.md:** Model assignment rules
- **archive/CHAOS-CLEANUP-SUMMARY.md:** Phase 2 cleanup report
- **.sisyphus/boulder.json:** Active plan tracking

---

**Last Updated:** 2026-02-20 21:30 CET  
**Next Update:** After Phase 3.1 complete  
**Status:** ✅ ACTIVE INFINITY LOOP
