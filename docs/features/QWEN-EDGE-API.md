# QWEN-EDGE-API.md

## Übersicht

Dokumentation für Qwen 3.5 Edge Functions auf Vercel.

**Base URL:** `https://biomet-rics-01.vercel.app`

---

## Endpoints

### 1. Chat Completion

**POST** `/api/qwen/chat`

Chat-Interaktion mit Qwen 3.5.

**Request:**
```json
{
  "messages": [
    { "role": "user", "content": "Hallo, wer bist du?" }
  ],
  "temperature": 0.7,
  "max_tokens": 2048
}
```

**Response:**
```json
{
  "id": "chatcmpl-xxx",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "Ich bin Qwen 3.5, ein KI-Modell von Alibaba Cloud."
      }
    }
  ]
}
```

---

### 2. Vision Analysis

**POST** `/api/qwen/vision`

Bildanalyse mit Qwen 3.5 Vision.

**Request:**
```json
{
  "image": "data:image/jpeg;base64,...",
  "prompt": "Beschreibe den Inhalt dieses Bildes"
}
```

**Response:**
```json
{
  "choices": [
    {
      "message": {
        "content": "Das Bild zeigt eine Landschaft mit Bergen..."
      }
    }
  ]
}
```

---

### 3. OCR (Texterkennung)

**POST** `/api/qwen/ocr`

Texterkennung aus Bildern.

**Request:**
```json
{
  "image": "data:image/jpeg;base64,...",
  "language": "de"
}
```

**Response:**
```json
{
  "choices": [
    {
      "message": {
        "content": "Erkannter Text aus dem Bild..."
      }
    }
  ]
}
```

---

### 4. Video Understanding

**POST** `/api/qwen/video`

Video-Inhaltsanalyse (via Frame-Extraction).

**Request:**
```json
{
  "video_frames": [
    "data:image/jpeg;base64,...",
    "data:image/jpeg;base64,..."
  ],
  "prompt": "Beschreibe was in diesem Video passiert"
}
```

**Response:**
```json
{
  "choices": [
    {
      "message": {
        "content": "Das Video zeigt eine Person, die..."
      }
    }
  ]
}
```

---

## Fehlercodes

| Code | Beschreibung |
|------|-------------|
| 400 | Ungültige Anfrage |
| 401 | Authentifizierung fehlgeschlagen |
| 429 | Rate Limit erreicht |
| 500 | Server-Fehler |
| 503 | Service nicht verfügbar |

---

## Rate Limits

- **RPM:** 40 Requests/Minute (NVIDIA NIM)
- **Vercel:** 1000 concurrent executions

---

## Beispiele

### cURL

```bash
# Chat
curl -X POST https://biomet-rics-01.vercel.app/api/qwen/chat \
  -H "Authorization: Bearer $NVIDIA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"messages":[{"role":"user","content":"Hallo"}]}'

# Vision
curl -X POST https://biomet-rics-01.vercel.app/api/qwen/vision \
  -H "Authorization: Bearer $NVIDIA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"image":"data:image/jpeg;base64,...","prompt":"Beschreibe"}'
```

### JavaScript

```javascript
const response = await fetch('https://biomet-rics-01.vercel.app/api/qwen/chat', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    messages: [{ role: 'user', content: 'Hallo' }],
  }),
});

const data = await response.json();
console.log(data.choices[0].message.content);
```
