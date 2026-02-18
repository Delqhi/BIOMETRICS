# LIVE-CHAT.md - Live Chat Support

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Support  
**Author:** BIOMETRICS Support Team

---

## 1. Overview

This document describes the live chat support system for BIOMETRICS, enabling real-time customer support through web and mobile interfaces.

## 2. Architecture

### 2.1 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| Chat Server | Socket.io | Real-time messaging |
| Widget | React | Chat widget |
| Agent Dashboard | React | Agent interface |
| Message Store | Redis | Session storage |
| Persistence | PostgreSQL | Message history |

### 2.2 Setup

```yaml
# docker-compose.yml
services:
  chat-server:
    build: ./chat-server
    environment:
      REDIS_URL: redis://redis:6379
      DB_URL: postgresql://user:pass@db:5432/chat
    ports:
      - "53054:3000"  # Port Sovereignty: 3000 → 53054

  redis:
    image: redis:7-alpine
    
  chat-db:
    image: postgres:15
```

## 3. Chat Widget

### 3.1 Widget Embed

```html
<!-- Chat widget embed code -->
<script>
  window.biometricsChat = {
    appId: 'biometrics-web',
    position: 'bottom-right',
    greeting: 'Hi! How can we help you?',
    color: '#1a73e8',
   上班族: {
      enabled: true,
      hours: 'Mon-Fri 9:00-18:00 CET',
    },
  };
</script>
<script src="https://chat.biometrics.com/widget.js" async></script>
```

### 3.2 Widget Configuration

```typescript
interface ChatConfig {
  appId: string;
  position: 'bottom-right' | 'bottom-left';
  greeting?: string;
  color?: string;
  language?: string;
 上班族?: {
    enabled: boolean;
    hours: string;
  };
  prechatForm?: {
    enabled: boolean;
    fields: FormField[];
  };
  sound?: {
    enabled: boolean;
    url?: string;
  };
}
```

## 4. Messaging

### 4.1 Message Types

| Type | Description | Payload |
|------|-------------|---------|
| text | Text message | { text: string } |
| image | Image attachment | { url: string, name: string } |
| file | File attachment | { url: string, name: string, size: number } |
| quick_reply | Quick response | { options: QuickReply[] } |
| carousel | Carousel cards | { cards: CarouselCard[] } |
| system | System message | { message: string } |

### 4.2 Message Handler

```typescript
interface Message {
  id: string;
  conversationId: string;
  senderId: string;
  senderType: 'user' | 'agent' | 'bot';
  type: MessageType;
  content: any;
  timestamp: Date;
  read: boolean;
}

class MessageHandler {
  async sendMessage(conversationId: string, message: CreateMessageRequest): Promise<Message> {
    const msg: Message = {
      id: generateUUID(),
      conversationId,
      senderId: message.senderId,
      senderType: message.senderType,
      type: message.type,
      content: message.content,
      timestamp: new Date(),
      read: false,
    };
    
    // Save to database
    await db.messages.insert(msg);
    
    // Send via socket
    io.to(conversationId).emit('message', msg);
    
    // Handle bot response
    if (message.senderType === 'user') {
      await handleBotResponse(conversationId, msg);
    }
    
    return msg;
  }
  
  async handleBotResponse(conversationId: string, userMessage: Message): Promise<void> {
    const response = await bot.processMessage(userMessage.content.text);
    
    if (response) {
      await this.sendMessage(conversationId, {
        senderId: 'bot',
        senderType: 'bot',
        type: response.type,
        content: response.content,
      });
    }
  }
}
```

## 5. Agent Interface

### 5.1 Dashboard

```tsx
const AgentDashboard = () => {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [activeConversation, setActiveConversation] = useState<Conversation | null>(null);
  
  return (
    <div className="agent-dashboard">
      <Sidebar>
        <ConversationList
          conversations={conversations}
          onSelect={setActiveConversation}
        />
      </Sidebar>
      <Main>
        {activeConversation && (
          <ChatWindow conversation={activeConversation} />
        )}
      </Main>
    </div>
  );
};
```

### 5.2 Agent Tools

| Tool | Description |
|------|-------------|
| Canned Responses | Quick replies |
| Internal Notes | Notes not visible to user |
| Transfer | Transfer to another agent |
| Prioritize | Mark as urgent |
| Close | Close conversation |

## 6. Routing

### 6.1 Conversation Routing

```typescript
const routeConversation = async (conversation: Conversation): Promise<void> => {
  // Get online agents
  const agents = await getOnlineAgents(conversation.groupId);
  
  if (agents.length === 0) {
    // Add to queue
    await addToQueue(conversation);
    await sendSystemMessage(conversation.id, 'You are in queue. An agent will be with you shortly.');
    return;
  }
  
  // Select agent
  const agent = selectAgent(agents, conversation);
  
  // Assign
  await assignToAgent(conversation.id, agent.id);
  
  // Notify agent
  await notifyAgent(agent.id, conversation);
};

const selectAgent = (agents: Agent[], conversation: Conversation): Agent => {
  // Load balance by current chats
  return agents
    .sort((a, b) => a.activeChats - b.activeChats)
    .filter(a => a.maxChats > a.activeChats)[0];
};
```

### 6.2 Queue Management

```typescript
class QueueManager {
  async addToQueue(conversation: Conversation): Promise<void> {
    const queueItem = {
      conversationId: conversation.id,
      addedAt: new Date(),
      priority: conversation.priority,
      estimatedWait: await this.estimateWaitTime(),
    };
    
    await this.queue.add('chat_queue', queueItem, {
      priority: queueItem.priority,
    });
  }
  
  async processQueue(): Promise<void> {
    const item = await this.queue.process(async (job) => {
      const conversation = await getConversation(job.data.conversationId);
      await this.routeToAgent(conversation);
    });
  }
}
```

## 7. Canned Responses

### 7.1 Response Library

```typescript
interface CannedResponse {
  id: string;
  title: string;
  content: string;
  category: string;
  tags: string[];
  shortcut: string;  // e.g., "/greeting"
}

const CANNED_RESPONSES: CannedResponse[] = [
  {
    id: '1',
    title: 'Greeting',
    shortcut: '/hi',
    content: 'Hi {{name}}! Thanks for reaching out to BIOMETRICS support. How can I help you today?',
    category: 'General',
    tags: ['greeting', 'start'],
  },
  {
    id: '2',
    title: 'Billing Question',
    shortcut: '/billing',
    content: 'I understand you have a billing question. Let me look into that for you. Can you please provide your account email?',
    category: 'Billing',
    tags: ['billing', 'payment'],
  },
];
```

## 8. Analytics

### 8.1 Chat Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| Wait Time | Time in queue | < 2 min |
| Response Time | First reply | < 1 min |
| Resolution Rate | Solved in chat | > 70% |
| CSAT | Satisfaction | > 90% |
| Chat Duration | Average length | < 10 min |

### 8.2 Real-time Dashboard

```typescript
const ChatDashboard = () => {
  const [metrics, setMetrics] = useState<ChatMetrics>({});
  
  useEffect(() => {
    const ws = connectWebSocket('/chat/metrics');
    
    ws.onmessage = (event) => {
      setMetrics(JSON.parse(event.data));
    };
    
    return () => ws.close();
  }, []);
  
  return (
    <div className="metrics-dashboard">
      <MetricCard label="Active Chats" value={metrics.activeChats} />
      <MetricCard label="Agents Online" value={metrics.agentsOnline} />
      <MetricCard label="Avg Wait Time" value={metrics.avgWaitTime} />
      <MetricCard label="Queue Size" value={metrics.queueSize} />
    </div>
  );
};
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
