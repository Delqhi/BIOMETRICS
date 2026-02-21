# AGENT-MODEL-MAPPING.md

**Status:** ACTIVE  
**Version:** 1.0 (2026-02-19)  
**Mandate:** CRITICAL - NIEMALS 2 Agents mit gleichem Modell parallel!

---

## 🚨 KRITISCHE REGEL

**NIEMALS mehr als 1 Agent mit demselben Modell parallel laufen lassen!**

❌ **FALSCH:**
```typescript
// ALLE nutzen Qwen 3.5 - BLOCKING!
task(category="visual-engineering", prompt="...") // Qwen 3.5
task(category="visual-engineering", prompt="...") // Qwen 3.5 BLOCKED!
task(category="deep", prompt="...") // Qwen 3.5 BLOCKED!
```

✅ **RICHTIG:**
```typescript
// Verschiedene Modelle - PARALLEL!
task(category="visual-engineering", prompt="...") // Qwen 3.5
task(category="deep", model="opencode/kimi-k2.5-free", prompt="...") // Kimi K2.5
task(category="quick", model="opencode/minimax-m2.5-free", prompt="...") // MiniMax M2.5
```

---

## 📊 VERFÜGBARE MODELLE

| Modell | Provider | Category | Use Case | Max Parallel |
|--------|----------|----------|----------|--------------|
| **qwen/qwen3.5-397b-a17b** | NVIDIA NIM | build, visual-engineering, ultrabrain, artistry, writing, general | Haupt-Code, Complex Tasks, Planung | **1** |
| **opencode/kimi-k2.5-free** | OpenCode ZEN | deep | Heavy Lifting, Setup, Deep Research | **1** |
| **opencode/minimax-m2.5-free** | OpenCode ZEN | quick, explore, librarian, writing | **SUCHE, LESEN, MD-DATEIEN ERSTELLEN** - Schnell! | **10** |

---

## 🎯 NEUE REGEL: MD-DATEIEN ERSTELLEN

**WICHTIG:** Alle MD-Dateien sollen mit **MiniMax M2.5** erstellt werden!

### Warum?
- MiniMax ist SCHNELL
- MiniMax kann 10x parallel laufen
- MiniMax ist KOSTENLOS
- Qwen 3.5 ist zu langsam für Dokumentation

### Workflow:
1. **MiniMax** → Suchen, Lesen, MD-Dateien erstellen
2. **Qwen 3.5** → Code-Umsetzung, Planung
3. **Kimi K2.5** → Deep Analysis (wenn nötig)

### Beispiel:
```typescript
// Phase 1: MiniMax erstellt MD
task(
  category="writing", 
  model="opencode/minimax-m2.5-free", 
  run_in_background=true,
  prompt="Create VIDEO-GEN.md with 500+ lines"
)

// Phase 2: Qwen 3.5 setzt um
task(
  category="build", 
  model="qwen/qwen3.5-397b-a17b", 
  prompt="Implement video generation based on VIDEO-GEN.md"
)
```

---

## 🎯 MODEL-ZUWEISUNG (ORCHESTRATOR PFLICHT)

### BEI 1 AGENT:
```typescript
task(category="visual-engineering", prompt="...") // Qwen 3.5
```

### BEI 2 AGENTS PARALLEL:
```typescript
task(category="visual-engineering", prompt="...") // Qwen 3.5
task(
  category="deep",
  model="opencode/kimi-k2.5-free",
  prompt="..."
) // Kimi K2.5
```

### BEI 3 AGENTS PARALLEL (MAXIMUM):
```typescript
task(category="visual-engineering", prompt="...") // Qwen 3.5
task(
  category="deep",
  model="opencode/kimi-k2.5-free",
  prompt="..."
) // Kimi K2.5
task(
  category="quick",
  model="opencode/minimax-m2.5-free",
  prompt="..."
) // MiniMax M2.5
```

---

## ⚠️ ORCHESTRATOR CHECKLISTE (VOR JEDER DELEGATION)

```markdown
## MODEL-MAPPING CHECK

**Laufende Agents:**
- [ ] Agent A: [Name] → [Modell]
- [ ] Agent B: [Name] → [Modell]
- [ ] Agent C: [Name] → [Modell]

**Neuer Agent:**
- [ ] Benötigtes Modell: [Qwen 3.5 / Kimi K2.5 / MiniMax M2.5]
- [ ] Modell bereits vergeben? JA/NEIN
- [ ] Wenn JA → Warte oder anderes Modell wählen!

**Entscheidung:**
- [ ] Agent starten (Modell frei)
- [ ] Warten bis Agent fertig (Modell belegt)
- [ ] Anderes Modell wählen
```

---

## 🔄 AGENT-RECYCLING

**Wenn alle 3 Modelle belegt sind:**

1. **Warte** auf Completion eines Agents
2. **Lese Session** des fertigen Agents
3. **Prüfe** "Sicher?"-Check
4. **Recycle** Modell für nächsten Agent

---

## 📋 BEISPIEL: BIOMETRICS MEDIA PIPELINE

### WAVE 1: Foundation (3 Agents parallel)

```typescript
// Agent 1: Qwen 3.5 - VIDEO-GEN.md
task(
  category="writing",
  prompt="Create VIDEO-GEN.md"
)

// Agent 2: Kimi K2.5 - Setup & Config
task(
  category="deep",
  model="opencode/kimi-k2.5-free",
  prompt="Complete BIOMETRICS setup"
)

// Agent 3: MiniMax M2.5 - Quick Configs
task(
  category="quick",
  model="opencode/minimax-m2.5-free",
  prompt="Create .env.example"
)
```

### WAVE 2: Video Agents (warte auf WAVE 1 Completion)

```typescript
// Nach WAVE 1 completion:
// Agent 4: Qwen 3.5 - cosmos_video_gen.py
task(
  category="visual-engineering",
  prompt="Create cosmos_video_gen.py"
)

// Agent 5: Kimi K2.5 - Integration Tests
task(
  category="deep",
  model="opencode/kimi-k2.5-free",
  prompt="Integration tests"
)
```

---

## 🚨 HÄUFIGE FEHLER

### ❌ FEHLER 1: Alle Agents Qwen 3.5
```typescript
// BLOCKING!
task(category="visual-engineering", prompt="...") // Qwen 3.5
task(category="visual-engineering", prompt="...") // Qwen 3.5 - WARTET!
task(category="deep", prompt="...") // Qwen 3.5 - WARTET!
```

### ❌ FEHLER 2: Falsche Category für Modell
```typescript
// deep category braucht Kimi K2.5, nicht Qwen 3.5!
task(category="deep", prompt="...") // FALSCH!
task(category="deep", model="opencode/kimi-k2.5-free", prompt="...") // RICHTIG!
```

### ❌ FEHLER 3: Model nicht angegeben
```typescript
// MiniMax M2.5 wird nicht genutzt!
task(category="quick", prompt="...") // Nutzt default (Qwen 3.5)!
task(category="quick", model="opencode/minimax-m2.5-free", prompt="...") // RICHTIG!
```

---

## ✅ BEST PRACTICES

1. **Immer Modell explizit angeben** bei deep, quick, explore, librarian
2. **Max 3 Agents parallel** (je 1 pro Modell)
3. **Sessions laufend lesen** während Agents arbeiten
4. **"Sicher?"-Check** nach jeder Completion
5. **Modell recyclen** nach Completion

---

**Last Updated:** 2026-02-19  
**Orchestrator:** MUST READ before EVERY delegation!
