# CHATBOT.md - Rasa Chatbot

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Support  
**Author:** BIOMETRICS AI Team

---

## 1. Overview

This document describes the Rasa chatbot implementation for BIOMETRICS, enabling AI-powered automated support and conversational interfaces.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| NLU | Rasa NLU | Intent recognition |
| Core | Rasa Core | Dialogue management |
| Actions | Python | Business logic |
| Tracker | Redis | Conversation state |
| Embedding | BERT | Sentence embeddings |

### 2.2 Setup

```yaml
# docker-compose.yml
services:
  rasa:
    image: rasa/rasa:3.6.0
    volumes:
      - ./rasa/models:/models
      - ./rasa/data:/data
      - ./rasa/actions:/actions
    ports:
      - "5005:5005"
    environment:
      RASA_MODEL: /models/current.tar.gz
      RASA_PORT: 5005

  action-server:
    build: ./rasa-actions
    ports:
      - "5055:5055"
```

## 3. NLU Training

### 3.1 Intent Definition

```yaml
# data/nlu.yml
version: "3.1"

nlu:
  - intent: greet
    examples: |
      - hey
      - hello
      - hi there
      - good morning
      - good evening
      - hey there

  - intent: goodbye
    examples: |
      - bye
      - goodbye
      - see you later
      - talk to you later
      - i have to go

  - intent: get_health_metrics
    examples: |
      - show my health metrics
      - what's my heart rate
      - how did i sleep last night
      - tell me about my health
      - show my daily stats

  - intent: report_issue
    examples: |
      - i have a problem
      - something is not working
      - i need help
      - there's an error
      - the app crashed

  - intent: subscription_help
    examples: |
      - how do i upgrade
      - what plans are available
      - change my subscription
      - cancel my plan
      - billing question
```

### 3.2 Entity Extraction

```yaml
  - intent: track_metric
    examples: |
      - track my [heart rate](metric_type)
      - log my [steps](metric_type) for [yesterday](date)
      - record my [sleep](metric_type) from [last night](date)
      - add [weight](metric_type) [today](date)

  - synonym: heart_rate
    examples:
      - heartrate
      - pulse
      - bpm

  - synonym: steps
    examples:
      - walking
      - footsteps
      - step count
```

## 4. Dialogue Management

### 4.1 Stories

```yaml
# data/stories.yml
version: "3.1"

stories:
  - story: greet and ask for help
    steps:
      - intent: greet
      - action: utter_greet
      - intent: get_help
      - action: action_offer_help

  - story: track health metric
    steps:
      - intent: track_metric
      - action: ask_metric_details
      - intent: provide_metric_value
        entities:
          - metric_type
          - value
      - action: action_save_metric
      - action: utter_confirm_saved

  - story: report issue
    steps:
      - intent: report_issue
      - action: ask_issue_details
      - intent: describe_issue
        entities:
          - issue_type
      - action: action_create_ticket
      - action: utter_ticket_created
```

### 4.2 Rules

```yaml
# data/rules.yml
version: "3.1"

rules:
  - rule: Always say goodbye
    steps:
      - intent: goodbye
      - action: utter_goodbye

  - rule: Handle unknown
    steps:
      - intent: nlu_fallback
      - action: utter_fallback
```

## 5. Custom Actions

### 5.1 Action Server

```python
# actions/actions.py
from typing import Any, Text, Dict, List
from rasa_sdk import Action, Tracker
from rasa_sdk.events import SlotSet, FollowupAction
from rasa_sdk.executor import CollectingDispatcher
import requests

class ActionTrackMetric(Action):
    def name(self) -> Text:
        return "action_save_metric"

    def run(
        self,
        dispatcher: CollectingDispatcher,
        tracker: Tracker,
        domain: Dict[Text, Any],
    ) -> List[Dict[Text, Any]]:
        
        metric_type = tracker.get_slot("metric_type")
        value = tracker.get_slot("value")
        date = tracker.get_slot("date")
        
        # Save to database
        save_metric(
            user_id=tracker.sender_id,
            metric_type=metric_type,
            value=float(value),
            date=date or datetime.now()
        )
        
        dispatcher.utter_message(
            text=f"Got it! I've recorded your {metric_type} as {value}."
        )
        
        return []

class ActionGetHealthSummary(Action):
    def name(self) -> Text:
        return "action_health_summary"

    def run(
        self,
        dispatcher: CollectingDispatcher,
        tracker: Tracker,
        domain: Dict[Text, Any],
    ) -> List[Dict[Text, Any]]:
        
        user_id = tracker.sender_id
        summary = get_health_summary(user_id)
        
        message = f"""
Here's your health summary:
❤️ Heart Rate: {summary.heart_rate_avg} bpm avg
😴 Sleep: {summary.sleep_hours} hours avg
👟 Steps: {summary.steps_avg} daily avg
        """
        
        dispatcher.utter_message(text=message)
        
        return []
```

### 5.2 Forms

```python
class ReportIssueForm(FormAction):
    def name(self) -> Text:
        return "report_issue_form"

    @staticmethod
    def required_slots(tracker: Tracker) -> List[Text]:
        return ["issue_type", "description", "email"]

    def slot_mappings(self) -> Dict[Text, Any]:
        return {
            "issue_type": [
                self.from_entity(entity="issue_type"),
                self.from_text(intent="report_issue"),
            ],
            "description": [self.from_text()],
            "email": [self.from_entity(entity="email")],
        }

    def submit(
        self,
        dispatcher: CollectingDispatcher,
        tracker: Tracker,
        domain: Dict[Text, Any],
    ) -> List[Dict]:
        
        # Create ticket
        ticket_id = create_support_ticket(
            issue_type=tracker.get_slot("issue_type"),
            description=tracker.get_slot("description"),
            email=tracker.get_slot("email")
        )
        
        dispatcher.utter_message(
            text=f"I've created a support ticket for you. Ticket #{ticket_id}"
        )
        
        return []
```

## 6. Response Templates

### 6.1 Utterances

```yaml
# responses/utterances.yml
responses:
  utter_greet:
    - text: "Hello! Welcome to BIOMETRICS. How can I help you today?"
    - text: "Hi there! I'm your BIOMETRICS assistant. What can I do for you?"

  utter_offer_help:
    - text: "I can help you with:"
      buttons:
        - title: "Track health metrics"
          payload: "/track_metric"
        - title: "View my summary"
          payload: "/get_health_metrics"
        - title: "Report an issue"
          payload: "/report_issue"
        - title: "Subscription help"
          payload: "/subscription_help"

  utter_fallback:
    - text: "I'm not sure I understand. Could you please rephrase that?"
    - text: "I'm still learning! Could you try asking in a different way?"
```

## 7. Deployment

### 7.1 Training

```bash
# Train Rasa model
rasa train --nlu data/nlu.yml \
           --domain domain.yml \
           --stories data/stories.yml \
           --config config.yml \
           --out models

# Start server
rasa run --model models/current.tar.gz \
         --endpoints endpoints.yml \
         --credentials credentials.yml
```

### 7.2 Docker Deployment

```dockerfile
# Dockerfile
FROM rasa/rasa:3.6.0

COPY . /app
WORKDIR /app

RUN rasa train --force

CMD ["rasa", "run", "--model", "/app/models", "--endpoints", "/app/endpoints.yml"]
```

## 8. Analytics

### 8.1 Conversation Analytics

| Metric | Description | Target |
|--------|-------------|--------|
| Success Rate | Task completed | > 80% |
| Fallback Rate | Unrecognized | < 10% |
| Avg Turns | Turns per conversation | < 5 |
| Resolution Rate | Self-service resolved | > 60% |

### 8.2 Evaluation

```bash
# Test model
rasa test nlu --nlu data/nlu.yml --model models/current.tar.gz

# Test conversation
rasa test core --stories data/test_stories.yml --model models/current.tar.gz
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
