# VOICE-ASSISTANT.md - Voice Assistant Integration

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Support  
**Author:** BIOMETRICS AI Team

---

## 1. Overview

This document describes the voice assistant integration for BIOMETRICS, enabling voice-based interactions through Alexa, Google Assistant, and custom voice interfaces.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| Voice Platform | Custom + Alexa + GAssistant | Platform integration |
| STT | Whisper | Speech to text |
| TTS | Edge TTS | Text to speech |
| NLU | Rasa | Intent recognition |
| Integration | API | BIOMETRICS backend |

### 2.2 Voice Flow

```
User Speech → STT → NLU → Intent → Action → Response → TTS → Audio
```

## 3. Alexa Integration

### 3.1 Skill Definition

```json
{
  "manifest": {
    "apis": {
      "custom": {
        "endpoint": {
          "uri": "arn:aws:lambda:us-east-1:account:function:biometrics-alexa"
        },
        "interfaces": [
          {
            "type": "AUDIO_PLAYER"
          }
        ]
      }
    },
    "publishingInformation": {
      "name": "BIOMETRICS Health",
      "description": "Track your health metrics with voice",
      "category": "HEALTH_AND_FITNESS"
    }
  }
}
```

### 3.2 Intent Handlers

```python
# alexa/handlers.py
from ask_sdk_core.dispatch_components import AbstractRequestHandler
from ask_sdk_core.utils import get_slot_value

class GetHealthMetricsHandler(AbstractRequestHandler):
    def can_handle(self, handler_input):
        return handler_input.request_envelope.request.type == "LaunchRequest" or \
               (handler_input.request_envelope.request.type == "IntentRequest" and
                handler_input.request_envelope.request.intent.name == "GetHealthMetricsIntent")

    def handle(self, handler_input):
        user_id = handler_input.request_envelope.context.system.user.user_id
        metrics = get_health_metrics(user_id)
        
        speech = f"Your latest health metrics are: " \
                  f"Heart rate {metrics.heart_rate} beats per minute, " \
                  f"you took {metrics.steps} steps today, " \
                  f"and slept {metrics.sleep_hours} hours last night."
        
        return handler_input.response_builder.speak(speech).response

class TrackMetricHandler(AbstractRequestHandler):
    def can_handle(self, handler_input):
        return handler_input.request_envelope.request.intent.name == "TrackMetricIntent"

    def handle(self, handler_input):
        metric_type = get_slot_value(handler_input, "metricType")
        value = get_slot_value(handler_input, "value")
        
        user_id = handler_input.request_envelope.context.system.user.user_id
        save_metric(user_id, metric_type, value)
        
        speech = f"I've recorded your {metric_type} as {value}."
        
        return handler_input.response_builder.speak(speech).response
```

### 3.3 Interaction Model

```json
{
  "intents": [
    {
      "name": "GetHealthMetricsIntent",
      "slots": [],
      "samples": [
        "what are my health metrics",
        "how am i doing",
        "tell me my stats",
        "my health summary"
      ]
    },
    {
      "name": "TrackMetricIntent",
      "slots": [
        {
          "name": "metricType",
          "type": "METRIC_TYPE"
        },
        {
          "name": "value",
          "type": "AMAZON.NUMBER"
        }
      ],
      "samples": [
        "log my {metricType} as {value}",
        "record {value} for {metricType}",
        "track {metricType} {value}"
      ]
    },
    {
      "name": "GetReminderIntent",
      "slots": [
        {
          "name": "reminderType",
          "type": "REMINDER_TYPE"
        }
      ],
      "samples": [
        "remind me to {reminderType}",
        "set a {reminderType} reminder"
      ]
    }
  ]
}
```

## 4. Google Assistant

### 4.1 Dialogflow Agent

```yaml
# dialogflow/agent.yaml
version: 2

intent:
  - name: "health.metrics"
    auto: true
    contexts: []
    responses:
      - defaultResponse:
          messages:
            - speech: "Your health metrics are: Heart rate $session:heartRate bpm, Steps $session:steps steps, Sleep $session:sleep hours."

intent:
  - name: "track.metric"
    auto: true
    parameters:
      - name: "metricType"
        entity: "metricType"
        required: true
      - name: "value"
        entity: "number"
        required: true
    responses:
      - defaultResponse:
          messages:
            - speech: "Recorded. Your $metricType is $value."

fulfillment:
  enabled: true
  webhook: "https://api.biometrics.com/dialogflow/fulfillment"
```

### 4.2 Fulfillment

```python
# dialogflow/fulfillment.py
from flask import Flask, request, jsonify
import dialogflow_v2 as dialogflow

app = Flask(__name__)

@app.route('/dialogflow/fulfillment', methods=['POST'])
def fulfillment():
    req = request.get_json()
    
    intent_name = req.get('queryResult', {}).get('intent', {}).get('displayName')
    
    if intent_name == 'health.metrics':
        user_id = get_user_id_from_session(req)
        metrics = get_health_metrics(user_id)
        
        return jsonify({
            'fulfillmentText': f"Your health metrics: Heart rate {metrics.heart_rate} bpm, {metrics.steps} steps, {metrics.sleep_hours} hours sleep."
        })
    
    elif intent_name == 'track.metric':
        metric_type = req['queryResult']['parameters']['metricType']
        value = req['queryResult']['parameters']['value']
        user_id = get_user_id_from_session(req)
        
        save_metric(user_id, metric_type, value)
        
        return jsonify({
            'fulfillmentText': f"Recorded your {metric_type} as {value}."
        })
```

## 5. Custom Voice

### 5.1 STT Implementation

```python
# voice/stt.py
import whisper

class SpeechToText:
    def __init__(self, model_size="base"):
        self.model = whisper.load_model(model_size)
    
    def transcribe(self, audio_data: bytes) -> str:
        import io
        import numpy as np
        import soundfile as sf
        
        # Load audio from bytes
        audio = self.load_audio(audio_data)
        
        # Transcribe
        result = self.model.transcribe(audio)
        
        return result['text']
    
    def load_audio(self, audio_bytes: bytes) -> np.ndarray:
        import io
        import soundfile as sf
        
        buffer = io.BytesIO(audio_bytes)
        audio, samplerate = sf.read(buffer)
        
        # Convert to mono if stereo
        if len(audio.shape) > 1:
            audio = audio.mean(axis=1)
        
        # Resample if needed
        if samplerate != 16000:
            # Resample to 16kHz
            import resampy
            audio = resampy.resample(audio, samplerate, 16000)
        
        return audio
```

### 5.2 TTS Implementation

```python
# voice/tts.py
import edge_tts
import asyncio
import io

class TextToSpeech:
    def __init__(self, voice="en-US-JennyNeural"):
        self.voice = voice
    
    async def synthesize(self, text: str) -> bytes:
        communicate = edge_tts.Communicate(text, self.voice)
        
        audio_buffer = io.BytesIO()
        
        async for chunk in communicate.stream():
            if chunk["type"] == "audio":
                audio_buffer.write(chunk["data"])
        
        return audio_buffer.getvalue()
    
    def synthesize_sync(self, text: str) -> bytes:
        return asyncio.run(self.synthesize(text))
```

## 6. Voice Commands

### 6.1 Command Schema

| Command | Intent | Parameters | Example |
|---------|--------|------------|---------|
| Get metrics | health_metrics | - | "How am I doing?" |
| Track metric | track_metric | type, value | "Log my heart rate as 72" |
| Set reminder | set_reminder | type, time | "Remind me to exercise at 7am" |
| Get summary | daily_summary | date | "Give me yesterday's summary" |
| Get alerts | get_alerts | - | "Any health alerts?" |

### 6.2 Command Handler

```python
class VoiceCommandHandler:
    def __init__(self):
        self.stt = SpeechToText()
        self.tts = TextToSpeech()
        self.nlu = RasaClient()
    
    async def process_voice_input(self, audio: bytes) -> bytes:
        # Convert speech to text
        text = self.stt.transcribe(audio)
        
        # Process with NLU
        intent = self.nlu.parse(text)
        
        # Execute action
        response = await self.execute_intent(intent)
        
        # Convert response to speech
        response_audio = self.tts.synthesize(response)
        
        return response_audio
    
    async def execute_intent(self, intent: dict) -> str:
        handler = self.intent_handlers.get(intent['name'])
        
        if handler:
            return await handler(intent)
        else:
            return "I'm sorry, I didn't understand that."
    
    async def handle_health_metrics(self, intent: dict) -> str:
        metrics = get_health_metrics(intent['user_id'])
        
        return (
            f"Your current health status: "
            f"Heart rate is {metrics.heart_rate} beats per minute. "
            f"You've taken {metrics.steps} steps today. "
            f"And you slept {metrics.sleep_hours} hours last night."
        )
```

## 7. Voice App Features

### 7.1 Capabilities

| Feature | Alexa | Google | Custom |
|---------|-------|--------|--------|
| Get metrics | ✅ | ✅ | ✅ |
| Track metrics | ✅ | ✅ | ✅ |
| Health alerts | ✅ | ✅ | ✅ |
| Reminders | ✅ | ✅ | ✅ |
| Coaching | ❌ | ❌ | ✅ |
| Emergency | ❌ | ✅ | ✅ |

### 7.2 User Verification

```python
class VoiceVerification:
    def verify_user(self, audio: bytes) -> Optional[str]:
        # Use voice biometrics
        voice_model = self.extract_voice_features(audio)
        
        # Match against registered voices
        for user_id, voice_template in self.voice_templates.items():
            if self.compare_voices(voice_model, voice_template):
                return user_id
        
        return None
    
    def compare_voices(self, voice1: np.ndarray, voice2: np.ndarray) -> bool:
        similarity = cosine_similarity(voice1, voice2)
        return similarity > 0.85
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
