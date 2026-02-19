# BIOMETRICS API Code Examples

Practical code examples for common BIOMETRICS API operations.

## Table of Contents

- [Health Check](#health-check)
- [Agent Management](#agent-management)
- [Session Management](#session-management)
- [Task Management](#task-management)
- [Model Management](#model-management)
- [Error Handling](#error-handling)
- [Advanced Patterns](#advanced-patterns)

---

## Health Check

### Basic Health Check

**cURL:**
```bash
curl https://api.biometrics.dev/v1/health
```

**JavaScript:**
```javascript
const response = await fetch('https://api.biometrics.dev/v1/health');
const data = await response.json();
console.log(data);
```

**Python:**
```python
import requests

response = requests.get('https://api.biometrics.dev/v1/health')
print(response.json())
```

**Go:**
```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type HealthResponse struct {
    Status    string `json:"status"`
    Version   string `json:"version"`
    Uptime    int64  `json:"uptime"`
    Timestamp string `json:"timestamp"`
}

func main() {
    resp, err := http.Get("https://api.biometrics.dev/v1/health")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    var health HealthResponse
    json.NewDecoder(resp.Body).Decode(&health)
    fmt.Printf("Status: %s, Version: %s\n", health.Status, health.Version)
}
```

---

## Agent Management

### List All Agents

**cURL:**
```bash
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://api.biometrics.dev/v1/agents?status=active&limit=10
```

**JavaScript:**
```javascript
const response = await fetch('https://api.biometrics.dev/v1/agents?status=active&limit=10', {
  headers: {
    'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`
  }
});

const data = await response.json();
console.log(`Found ${data.total} active agents`);
data.agents.forEach(agent => {
  console.log(`- ${agent.name} (${agent.role}): ${agent.status}`);
});
```

**Python:**
```python
import requests
import os

headers = {'Authorization': f"Bearer {os.environ['NVIDIA_API_KEY']}"}
params = {'status': 'active', 'limit': 10}

response = requests.get(
    'https://api.biometrics.dev/v1/agents',
    headers=headers,
    params=params
)

data = response.json()
print(f"Found {data['total']} active agents")
for agent in data['agents']:
    print(f"- {agent['name']} ({agent['role']}): {agent['status']}")
```

### Get Agent Details

**cURL:**
```bash
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://api.biometrics.dev/v1/agents/agent-sisyphus-001
```

**JavaScript:**
```javascript
const agentId = 'agent-sisyphus-001';
const response = await fetch(
  `https://api.biometrics.dev/v1/agents/${agentId}`,
  {
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`
    }
  }
);

const agent = await response.json();
console.log(`${agent.name} has completed ${agent.tasksCompleted} tasks`);
```

### Create New Agent

**cURL:**
```bash
curl -X POST \
  -H "Authorization: Bearer $NVIDIA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Custom-Coder",
    "role": "coder",
    "model": "qwen/qwen3.5-397b-a17b"
  }' \
  https://api.biometrics.dev/v1/agents
```

**JavaScript:**
```javascript
const response = await fetch('https://api.biometrics.dev/v1/agents', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'Custom-Coder',
    role: 'coder',
    model: 'qwen/qwen3.5-397b-a17b'
  })
});

const agent = await response.json();
console.log(`Created agent: ${agent.id}`);
```

### Get Agent Tasks

**cURL:**
```bash
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  "https://api.biometrics.dev/v1/agents/agent-sisyphus-001/tasks?status=in_progress"
```

**Python:**
```python
import requests
import os

agent_id = 'agent-sisyphus-001'
headers = {'Authorization': f"Bearer {os.environ['NVIDIA_API_KEY']}"}
params = {'status': 'in_progress'}

response = requests.get(
    f'https://api.biometrics.dev/v1/agents/{agent_id}/tasks',
    headers=headers,
    params=params
)

tasks = response.json()['tasks']
print(f"Agent has {len(tasks)} tasks in progress")
for task in tasks:
    print(f"- {task['title']} (Priority: {task['priority']})")
```

---

## Session Management

### Create Session

**cURL:**
```bash
curl -X POST \
  -H "Authorization: Bearer $NVIDIA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"title": "Build REST API", "agentId": "agent-sisyphus-001"}' \
  https://api.biometrics.dev/v1/sessions
```

**JavaScript:**
```javascript
const response = await fetch('https://api.biometrics.dev/v1/sessions', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    title: 'Build REST API',
    agentId: 'agent-sisyphus-001'
  })
});

const session = await response.json();
console.log(`Created session: ${session.id}`);
```

### Send Prompt

**cURL:**
```bash
curl -X POST \
  -H "Authorization: Bearer $NVIDIA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Build a REST API with Express and TypeScript"}' \
  https://api.biometrics.dev/v1/sessions/ses_abc123/prompt
```

**JavaScript (Streaming):**
```javascript
const response = await fetch(
  'https://api.biometrics.dev/v1/sessions/ses_abc123/prompt',
  {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      prompt: 'Build a REST API with Express and TypeScript',
      stream: true
    })
  }
);

const reader = response.body.getReader();
const decoder = new TextDecoder();

while (true) {
  const { done, value } = await reader.read();
  if (done) break;
  
  const chunk = decoder.decode(value);
  process.stdout.write(chunk);
}
```

**Python:**
```python
import requests
import os

session_id = 'ses_abc123'
headers = {
    'Authorization': f"Bearer {os.environ['NVIDIA_API_KEY']}",
    'Content-Type': 'application/json'
}

data = {
    'prompt': 'Build a REST API with Express and TypeScript'
}

response = requests.post(
    f'https://api.biometrics.dev/v1/sessions/{session_id}/prompt',
    headers=headers,
    json=data
)

result = response.json()
print(result['content'])
```

### Get Session Details

**cURL:**
```bash
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://api.biometrics.dev/v1/sessions/ses_abc123
```

**Go:**
```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Session struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    Status    string `json:"status"`
    CreatedAt string `json:"createdAt"`
}

func main() {
    req, _ := http.NewRequest("GET", 
        "https://api.biometrics.dev/v1/sessions/ses_abc123", nil)
    req.Header.Set("Authorization", 
        "Bearer "+os.Getenv("NVIDIA_API_KEY"))
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    var session Session
    json.NewDecoder(resp.Body).Decode(&session)
    fmt.Printf("Session: %s - %s\n", session.ID, session.Title)
}
```

---

## Task Management

### Create Task

**cURL:**
```bash
curl -X POST \
  -H "Authorization: Bearer $NVIDIA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Implement login endpoint",
    "description": "Create JWT-based authentication endpoint",
    "priority": "high",
    "agentId": "agent-sisyphus-001"
  }' \
  https://api.biometrics.dev/v1/tasks
```

**JavaScript:**
```javascript
const response = await fetch('https://api.biometrics.dev/v1/tasks', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    title: 'Implement login endpoint',
    description: 'Create JWT-based authentication endpoint',
    priority: 'high',
    agentId: 'agent-sisyphus-001'
  })
});

const task = await response.json();
console.log(`Created task: ${task.id} with priority ${task.priority}`);
```

### Get Task Status

**cURL:**
```bash
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://api.biometrics.dev/v1/tasks/task-001
```

**Python:**
```python
import requests
import os
import time

task_id = 'task-001'
headers = {'Authorization': f"Bearer {os.environ['NVIDIA_API_KEY']}"}

# Poll until task completes
while True:
    response = requests.get(
        f'https://api.biometrics.dev/v1/tasks/{task_id}',
        headers=headers
    )
    
    task = response.json()
    print(f"Status: {task['status']}")
    
    if task['status'] in ['completed', 'failed', 'cancelled']:
        break
    
    time.sleep(5)  # Wait 5 seconds before checking again

if task['status'] == 'completed':
    print(f"Result: {task['result']}")
else:
    print(f"Error: {task['error']}")
```

### Cancel Task

**cURL:**
```bash
curl -X POST \
  -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://api.biometrics.dev/v1/tasks/task-001/cancel
```

**JavaScript:**
```javascript
const response = await fetch(
  'https://api.biometrics.dev/v1/tasks/task-001/cancel',
  {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`
    }
  }
);

const task = await response.json();
console.log(`Task cancelled: ${task.id}`);
```

---

## Model Management

### List Available Models

**cURL:**
```bash
curl -H "Authorization: Bearer $NVIDIA_API_KEY" \
  https://api.biometrics.dev/v1/models
```

**JavaScript:**
```javascript
const response = await fetch('https://api.biometrics.dev/v1/models', {
  headers: {
    'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`
  }
});

const data = await response.json();
console.log('Available Models:');
data.models.forEach(model => {
  console.log(`- ${model.name} (${model.id})`);
  console.log(`  Context: ${model.contextLimit}, Output: ${model.outputLimit}`);
});
```

**Python:**
```python
import requests
import os

headers = {'Authorization': f"Bearer {os.environ['NVIDIA_API_KEY']}"}

response = requests.get(
    'https://api.biometrics.dev/v1/models',
    headers=headers
)

models = response.json()['models']
for model in models:
    print(f"{model['name']}")
    print(f"  ID: {model['id']}")
    print(f"  Context: {model['contextLimit']:,} tokens")
    print(f"  Output: {model['outputLimit']:,} tokens")
    print()
```

---

## Error Handling

### Basic Error Handling

**JavaScript:**
```javascript
try {
  const response = await fetch('https://api.biometrics.dev/v1/agents', {
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`
    }
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(`${error.error.code}: ${error.error.message}`);
  }
  
  const data = await response.json();
  console.log(data);
} catch (error) {
  console.error('API Error:', error.message);
  
  if (error.message.includes('RATE_LIMIT')) {
    console.log('Waiting 60 seconds before retry...');
    await sleep(60000);
    // Retry logic here
  }
}
```

**Python:**
```python
import requests
from requests.exceptions import HTTPError

try:
    response = requests.get(
        'https://api.biometrics.dev/v1/agents',
        headers={'Authorization': f"Bearer {os.environ['NVIDIA_API_KEY']}"}
    )
    
    response.raise_for_status()
    data = response.json()
    print(data)
    
except HTTPError as e:
    error_data = e.response.json()
    print(f"API Error: {error_data['error']['code']}")
    print(f"Message: {error_data['error']['message']}")
    
    if error_data['error']['code'] == 'RATE_LIMIT_EXCEEDED':
        print("Waiting 60 seconds...")
        time.sleep(60)
        # Retry logic here
        
except Exception as e:
    print(f"Unexpected error: {e}")
```

### Retry with Exponential Backoff

**JavaScript:**
```javascript
async function makeRequestWithRetry(url, options, maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      const response = await fetch(url, options);
      
      if (!response.ok) {
        const error = await response.json();
        
        if (error.error.code === 'RATE_LIMIT_EXCEEDED') {
          const waitTime = Math.pow(2, i) * 1000; // 1s, 2s, 4s
          console.log(`Rate limited. Waiting ${waitTime}ms...`);
          await sleep(waitTime);
          continue;
        }
        
        throw new Error(`${error.error.code}: ${error.error.message}`);
      }
      
      return await response.json();
      
    } catch (error) {
      if (i === maxRetries - 1) throw error;
      console.log(`Retry ${i + 1}/${maxRetries}`);
      await sleep(Math.pow(2, i) * 1000);
    }
  }
}

// Usage
const agents = await makeRequestWithRetry(
  'https://api.biometrics.dev/v1/agents',
  {
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`
    }
  }
);
```

**Python:**
```python
import time
import random

def make_request_with_retry(url, headers, max_retries=3):
    for i in range(max_retries):
        try:
            response = requests.get(url, headers=headers)
            response.raise_for_status()
            return response.json()
            
        except requests.exceptions.HTTPError as e:
            error_data = e.response.json()
            
            if error_data['error']['code'] == 'RATE_LIMIT_EXCEEDED':
                wait_time = (2 ** i) + random.random()
                print(f"Rate limited. Waiting {wait_time:.1f}s...")
                time.sleep(wait_time)
                continue
                
            raise Exception(f"{error_data['error']['code']}: {error_data['error']['message']}")
            
        except Exception as e:
            if i == max_retries - 1:
                raise
            print(f"Retry {i + 1}/{max_retries}")
            time.sleep((2 ** i) + random.random())
    
    raise Exception("Max retries exceeded")

# Usage
agents = make_request_with_retry(
    'https://api.biometrics.dev/v1/agents',
    {'Authorization': f"Bearer {os.environ['NVIDIA_API_KEY']}"}
)
```

---

## Advanced Patterns

### Parallel Agent Tasks

**JavaScript:**
```javascript
// Create multiple tasks in parallel
const tasks = await Promise.all([
  fetch('https://api.biometrics.dev/v1/tasks', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      title: 'Implement feature A',
      agentId: 'agent-sisyphus-001'
    })
  }),
  fetch('https://api.biometrics.dev/v1/tasks', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      title: 'Implement feature B',
      agentId: 'agent-prometheus-001'
    })
  }),
  fetch('https://api.biometrics.dev/v1/tasks', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      title: 'Write documentation',
      agentId: 'agent-librarian-001'
    })
  })
]);

const results = await Promise.all(tasks.map(t => t.json()));
console.log(`Created ${results.length} parallel tasks`);
```

### Agent Swarm Workflow

**Python:**
```python
import requests
import os

API_KEY = os.environ['NVIDIA_API_KEY']
BASE_URL = 'https://api.biometrics.dev/v1'
headers = {'Authorization': f"Bearer {API_KEY}"}

# 1. Create session
session = requests.post(
    f'{BASE_URL}/sessions',
    headers=headers,
    json={'title': 'Build Complete Application'}
).json()

print(f"Created session: {session['id']}")

# 2. Send planning prompt
plan = requests.post(
    f'{BASE_URL}/sessions/{session["id"]}/prompt',
    headers=headers,
    json={'prompt': 'Plan a complete web application architecture'}
).json()

print(f"Plan generated: {len(plan['content'])} characters")

# 3. Create tasks based on plan
tasks = [
    {'title': 'Setup database schema', 'priority': 'high'},
    {'title': 'Implement API endpoints', 'priority': 'high'},
    {'title': 'Create frontend components', 'priority': 'medium'},
    {'title': 'Write tests', 'priority': 'medium'},
    {'title': 'Deploy to production', 'priority': 'low'}
]

for task_data in tasks:
    task = requests.post(
        f'{BASE_URL}/tasks',
        headers=headers,
        json={**task_data, 'sessionId': session['id']}
    ).json()
    print(f"Created task: {task['id']} - {task['title']}")

print(f"\nSwarm workflow started with {len(tasks)} tasks")
```

### Multi-Modal Prompt (Text + Image)

**JavaScript:**
```javascript
const fs = require('fs');

// Read image and convert to base64
const imageBuffer = fs.readFileSync('captcha.png');
const base64Image = imageBuffer.toString('base64');

const response = await fetch(
  'https://api.biometrics.dev/v1/sessions/ses_abc123/prompt',
  {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${process.env.NVIDIA_API_KEY}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      prompt: 'Solve this CAPTCHA',
      parts: [
        {
          type: 'text',
          text: 'What text is shown in this image?'
        },
        {
          type: 'file',
          mime: 'image/png',
          filename: 'captcha.png',
          url: `data:image/png;base64,${base64Image}`
        }
      ]
    })
  }
);

const result = await response.json();
console.log(`CAPTCHA solution: ${result.content}`);
```

---

## Support

- **Documentation**: https://github.com/Delqhi/BIOMETRICS/docs
- **Examples**: https://github.com/Delqhi/BIOMETRICS/tree/main/docs/api/examples
- **Issues**: https://github.com/Delqhi/BIOMETRICS/issues
- **Discord**: https://discord.gg/biometrics

---

**Last Updated:** February 2026 | **Version:** 1.0.0
