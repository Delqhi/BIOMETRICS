# Interactive Quizzes

**Purpose:** Knowledge assessment quizzes for learning and certification

## Overview

This directory contains interactive quizzes designed to test knowledge and understanding of BIOMETRICS concepts, best practices, and technical details.

## Quiz Categories

### Technical Quizzes

| Quiz | Topic | Difficulty |
|------|-------|------------|
| cli-basics | CLI fundamentals | Beginner |
| config-management | Configuration handling | Intermediate |
| security | Security best practices | Advanced |
| deployment | Deployment procedures | Intermediate |
| troubleshooting | Problem diagnosis | Advanced |

### Mandate Quizzes

| Quiz | Mandate | Focus |
|------|---------|-------|
| mandate-essentials | MANDATE 0.x | Core requirements |
| swarm-operations | MANDATE 0.1 | Agent swarm |
| quality-standards | MANDATE 0.2 | Code quality |

## Usage

### Running Quizzes

```bash
# List available quizzes
biometrics quiz list

# Start specific quiz
biometrics quiz start cli-basics

# Review results
biometrics quiz results --latest
```

### Quiz Format

Each quiz contains:
- Multiple choice questions
- Code analysis challenges
- Scenario-based problems
- Timed sections

## Scoring

| Score | Grade |
|-------|-------|
| 90-100% | Expert |
| 75-89% | Proficient |
| 60-74% | Competent |
| <60% | Needs Improvement |

## Creating Custom Quizzes

### Quiz Structure
```yaml
# quiz-name.yaml
title: "Quiz Title"
description: "Description"
questions:
  - id: 1
    type: multiple_choice
    question: "Question text?"
    options:
      - A. Option A
      - B. Option B
    correct: A
    explanation: "Why A is correct"
```

## Integration with Learning

### Recommended Flow
1. Read documentation
2. Take relevant quiz
3. Review incorrect answers
4. Re-read weak areas
5. Retake quiz

## Related Documentation

- [Best Practices](../best-practices/)
- [Agent Guide](../agents/)
- [Architecture](../architecture/)
