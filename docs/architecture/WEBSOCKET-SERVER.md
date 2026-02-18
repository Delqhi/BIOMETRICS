# WEBSOCKET-SERVER.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

---

## 1) Overview

This document describes the WebSocket server architecture for BIOMETRICS. The WebSocket server provides real-time bidirectional communication for biometric events, session management, notifications, and collaborative features.

### Core Features
- Real-time biometric event streaming
- Live session monitoring
- Push notifications
- Presence indicators
- Message broadcasting
- Auto-reconnection
- Heartbeat/ping-pong
- Connection pooling

---

## 2) Technology Stack

### Framework & Libraries
```json
{
  "ws": "^8.18.0",
  "socket.io": "^4.7.0",
  "uWebSockets.js": "^20.0.0",
  "ws": "^8.18.0"
}
```

### Supporting Libraries
```json
{
  "ioredis": "^5.4.0",
  "jsonwebtoken": "^9.0.0",
  "uuid": "^9.0.0",
  "ws-heartbeat": "^2.0.0"
}
```

---

## 3) Server Architecture

### Main Server

```typescript
// src/websocket/server.ts
import { createServer } from 'http';
import { Server } from 'socket.io';
import { verifyToken } from './auth';
import { setupEventHandlers } from './handlers';
import { logger } from './logger';
import { RedisAdapter } from './adapters/redis';

const httpServer = createServer();

const io = new Server(httpServer, {
  cors: {
    origin: process.env.ALLOWED_ORIGINS?.split(',') || ['http://localhost:3000'],
    methods: ['GET', 'POST'],
    credentials: true,
  },
  transports: ['websocket', 'polling'],
  pingTimeout: 60000,
  pingInterval: 25000,
  maxHttpBufferSize: 1e8, // 100MB
  allowEIO3: true,
  path: '/ws',
});

// Redis adapter for scaling
io.adapter(RedisAdapter({
  host: process.env.REDIS_HOST || 'localhost',
  port: parseInt(process.env.REDIS_PORT || '6379'),
  password: process.env.REDIS_PASSWORD,
}));

// Authentication middleware
io.use(async (socket, next) => {
  const token = socket.handshake.auth.token || socket.handshake.headers.authorization?.replace('Bearer ', '');
  
  if (!token) {
    return next(new Error('Authentication required'));
  }

  try {
    const user = await verifyToken(token);
    socket.data.user = user;
    next();
  } catch (error) {
    next(new Error('Invalid token'));
  }
});

// Connection handling
io.on('connection', (socket) => {
  const userId = socket.data.user.id;
  
  logger.info(`Client connected: ${socket.id} (User: ${userId})`);

  // Join user's personal room
  socket.join(`user:${userId}`);

  // Setup event handlers
  setupEventHandlers(io, socket);

  // Handle disconnection
  socket.on('disconnect', (reason) => {
    logger.info(`Client disconnected: ${socket.id} (Reason: ${reason})`);
    socket.leave(`user:${userId}`);
  });

  // Handle errors
  socket.on('error', (error) => {
    logger.error(`Socket error: ${socket.id}`, error);
  });
});

const PORT = process.env.WS_PORT || 3001;

httpServer.listen(PORT, () => {
  logger.info(`WebSocket server running on port ${PORT}`);
});

export { io, httpServer };
```

---

## 4) Event Handlers

### Biometric Events

```typescript
// src/websocket/handlers/biometric.ts
import { Server, Socket } from 'socket.io';
import { logger } from '../logger';
import { publishBiometricEvent } from '../pubsub';

export function setupBiometricHandlers(io: Server, socket: Socket) {
  const userId = socket.data.user.id;

  // Subscribe to biometric events
  socket.on('biometric:subscribe', async (data: { recordId?: string }) => {
    if (data.recordId) {
      socket.join(`biometric:${data.recordId}`);
      logger.info(`User ${userId} subscribed to biometric record ${data.recordId}`);
    }
  });

  // Unsubscribe from biometric events
  socket.on('biometric:unsubscribe', async (data: { recordId?: string }) => {
    if (data.recordId) {
      socket.leave(`biometric:${data.recordId}`);
      logger.info(`User ${userId} unsubscribed from biometric record ${data.recordId}`);
    }
  });

  // Request biometric verification
  socket.on('biometric:verify', async (data: { recordId: string; challenge: string }) => {
    try {
      // Emit to specific room for processing
      io.to(`biometric:${data.recordId}`).emit('biometric:verification-started', {
        recordId: data.recordId,
        timestamp: Date.now(),
      });

      // Process verification (async)
      const result = await processBiometricVerification(userId, data);
      
      // Send result back
      socket.emit('biometric:verification-result', {
        recordId: data.recordId,
        success: result.success,
        confidence: result.confidence,
        timestamp: Date.now(),
      });

      // Publish event for other listeners
      publishBiometricEvent(userId, {
        type: 'verification',
        recordId: data.recordId,
        success: result.success,
      });
    } catch (error) {
      socket.emit('biometric:error', {
        message: error.message,
        recordId: data.recordId,
      });
    }
  });

  // Real-time streaming
  socket.on('biometric:stream-start', async (data: { deviceId: string }) => {
    socket.join(`stream:${data.deviceId}`);
    logger.info(`User ${userId} started streaming from device ${data.deviceId}`);
  });

  socket.on('biometric:stream-stop', async (data: { deviceId: string }) => {
    socket.leave(`stream:${data.deviceId}`);
    logger.info(`User ${userId} stopped streaming from device ${data.deviceId}`);
  });

  socket.on('biometric:stream-data', (data: { deviceId: string; frame: Buffer }) => {
    // Broadcast to all subscribers of this stream
    io.to(`stream:${data.deviceId}`).emit('biometric:stream-frame', {
      deviceId: data.deviceId,
      frame: data.frame.toString('base64'),
      timestamp: Date.now(),
      userId,
    });
  });
}

async function processBiometricVerification(userId: string, data: any) {
  // Implementation of biometric verification
  return { success: true, confidence: 0.95 };
}
```

---

## 5) Session Management

### Session Handlers

```typescript
// src/websocket/handlers/session.ts
import { Server, Socket } from 'socket.io';
import { logger } from '../logger';
import { getActiveSessions, createSession, endSession } from '../services/session';

export function setupSessionHandlers(io: Server, socket: Socket) {
  const userId = socket.data.user.id;
  const socketId = socket.id;

  // Get active sessions
  socket.on('session:list', async () => {
    try {
      const sessions = await getActiveSessions(userId);
      socket.emit('session:list-result', { sessions });
    } catch (error) {
      socket.emit('session:error', { message: error.message });
    }
  });

  // Create new session
  socket.on('session:create', async (data: { deviceInfo: string; ipAddress: string }) => {
    try {
      const session = await createSession({
        userId,
        deviceInfo: data.deviceInfo,
        ipAddress: data.ipAddress,
        socketId,
      });

      // Notify user
      socket.emit('session:created', { session });

      // Broadcast to user's other sessions
      io.to(`user:${userId}`).emit('session:new', {
        session,
        message: 'New session started',
      });

      logger.info(`New session created for user ${userId}: ${session.id}`);
    } catch (error) {
      socket.emit('session:error', { message: error.message });
    }
  });

  // End session
  socket.on('session:end', async (data: { sessionId: string }) => {
    try {
      await endSession(data.sessionId, userId);
      socket.emit('session:ended', { sessionId: data.sessionId });
      
      // Notify other sessions
      io.to(`user:${userId}`).emit('session:ended', {
        sessionId: data.sessionId,
      });

      logger.info(`Session ended for user ${userId}: ${data.sessionId}`);
    } catch (error) {
      socket.emit('session:error', { message: error.message });
    }
  });

  // Heartbeat
  socket.on('session:heartbeat', async () => {
    socket.data.lastHeartbeat = Date.now();
    await updateSessionHeartbeat(socketId);
  });

  // Request session sync
  socket.on('session:sync', async () => {
    const sessions = await getActiveSessions(userId);
    socket.emit('session:sync-result', { sessions });
  });
}
```

---

## 6) Presence System

### Presence Handler

```typescript
// src/websocket/handlers/presence.ts
import { Server, Socket } from 'socket.io';
import { Redis } from '../services/redis';
import { logger } from '../logger';

const PRESENCE_TTL = 30; // seconds

export function setupPresenceHandlers(io: Server, socket: Socket) {
  const userId = socket.data.user.id;

  // Update presence
  const updatePresence = async (status: 'online' | 'away' | 'offline') => {
    const key = `presence:${userId}`;
    await Redis.setex(key, PRESENCE_TTL, JSON.stringify({
      status,
      socketId: socket.id,
      timestamp: Date.now(),
    }));

    // Broadcast presence change
    io.to(`user:${userId}`).emit('presence:updated', {
      userId,
      status,
      timestamp: Date.now(),
    });
  };

  // Initial presence
  socket.on('presence:update', async (data: { status: 'online' | 'away' }) => {
    await updatePresence(data.status);
  });

  // Join room for team presence
  socket.on('presence:join-room', async (data: { roomId: string }) => {
    const roomKey = `room:${data.roomId}:members`;
    await Redis.sadd(roomKey, userId);
    socket.join(`room:${data.roomId}`);
    
    // Notify room members
    io.to(`room:${data.roomId}`).emit('presence:member-joined', {
      userId,
      roomId: data.roomId,
    });

    // Send current room members
    const members = await Redis.smembers(roomKey);
    socket.emit('presence:room-members', { roomId: data.roomId, members });
  });

  // Leave room
  socket.on('presence:leave-room', async (data: { roomId: string }) => {
    const roomKey = `room:${data.roomId}:members`;
    await Redis.srem(roomKey, userId);
    socket.leave(`room:${data.roomId}`);
    
    io.to(`room:${data.roomId}`).emit('presence:member-left', {
      userId,
      roomId: data.roomId,
    });
  });

  // Cleanup on disconnect
  socket.on('disconnect', async () => {
    await updatePresence('offline');
    
    // Clean up room memberships
    const roomKeys = await Redis.keys(`room:*:members`);
    for (const key of roomKeys) {
      await Redis.srem(key, userId);
    }
  });
}
```

---

## 7) Message Broadcasting

### Broadcast Handler

```typescript
// src/websocket/handlers/broadcast.ts
import { Server, Socket } from 'socket.io';
import { logger } from '../logger';

interface BroadcastMessage {
  type: string;
  payload: any;
  targetUsers?: string[];
  targetRoom?: string;
}

export function setupBroadcastHandlers(io: Server, socket: Socket) {
  const userId = socket.data.user.id;

  // Send to specific users
  socket.on('broadcast:to-users', async (message: BroadcastMessage) => {
    const { targetUsers, type, payload } = message;
    
    if (!targetUsers || targetUsers.length === 0) {
      socket.emit('broadcast:error', { message: 'No target users specified' });
      return;
    }

    for (const targetUserId of targetUsers) {
      io.to(`user:${targetUserId}`).emit('broadcast:message', {
        type,
        payload,
        from: userId,
        timestamp: Date.now(),
      });
    }

    logger.debug(`Broadcast from ${userId} to ${targetUsers.length} users`);
  });

  // Send to room
  socket.on('broadcast:to-room', async (message: BroadcastMessage) => {
    const { targetRoom, type, payload } = message;

    if (!targetRoom) {
      socket.emit('broadcast:error', { message: 'No target room specified' });
      return;
    }

    io.to(targetRoom).emit('broadcast:message', {
      type,
      payload,
      from: userId,
      timestamp: Date.now(),
    });

    logger.debug(`Broadcast from ${userId} to room ${targetRoom}`);
  });

  // Admin broadcast (with special permission)
  socket.on('broadcast:admin', async (message: BroadcastMessage) => {
    if (socket.data.user.role !== 'ADMIN') {
      socket.emit('broadcast:error', { message: 'Admin permission required' });
      return;
    }

    io.emit('broadcast:message', {
      type: message.type,
      payload: message.payload,
      from: 'admin',
      timestamp: Date.now(),
    });
  });
}
```

---

## 8) Notification System

### Notification Handler

```typescript
// src/websocket/handlers/notification.ts
import { Server, Socket } from 'socket.io';
import { createNotification, getUnreadNotifications } from '../services/notification';

export function setupNotificationHandlers(io: Server, socket: Socket) {
  const userId = socket.data.user.id;

  // Get unread notifications
  socket.on('notification:get-unread', async () => {
    try {
      const notifications = await getUnreadNotifications(userId);
      socket.emit('notification:unread-list', { notifications });
    } catch (error) {
      socket.emit('notification:error', { message: error.message });
    }
  });

  // Mark notification as read
  socket.on('notification:mark-read', async (data: { notificationId: string }) => {
    try {
      await markNotificationRead(data.notificationId, userId);
      socket.emit('notification:marked-read', { notificationId: data.notificationId });
    } catch (error) {
      socket.emit('notification:error', { message: error.message });
    }
  });

  // Mark all as read
  socket.on('notification:mark-all-read', async () => {
    try {
      await markAllNotificationsRead(userId);
      socket.emit('notification:all-marked-read');
    } catch (error) {
      socket.emit('notification:error', { message: error.message });
    }
  });

  // Subscribe to notification categories
  socket.on('notification:subscribe', async (data: { categories: string[] }) => {
    for (const category of data.categories) {
      socket.join(`notifications:${category}`);
    }
    socket.emit('notification:subscribed', { categories: data.categories });
  });

  // Real-time notification handler
  socket.on('notification:send', async (data: {
    targetUserId: string;
    type: string;
    title: string;
    body: string;
    data?: any;
  }) => {
    try {
      const notification = await createNotification({
        userId: data.targetUserId,
        type: data.type,
        title: data.title,
        body: data.body,
        data: data.data,
      });

      // Send real-time
      io.to(`user:${data.targetUserId}`).emit('notification:new', notification);
    } catch (error) {
      socket.emit('notification:error', { message: error.message });
    }
  });
}
```

---

## 9) Connection Pool

### Connection Manager

```typescript
// src/websocket/connection-manager.ts
import { Socket } from 'socket.io';
import { Redis } from '../services/redis';

interface ConnectionInfo {
  socketId: string;
  userId: string;
  connectedAt: number;
  lastActivity: number;
  rooms: Set<string>;
}

class ConnectionManager {
  private connections: Map<string, ConnectionInfo> = new Map();
  private userSockets: Map<string, Set<string>> = new Map();

  addConnection(socket: Socket, userId: string): void {
    const info: ConnectionInfo = {
      socketId: socket.id,
      userId,
      connectedAt: Date.now(),
      lastActivity: Date.now(),
      rooms: new Set(socket.rooms),
    };

    this.connections.set(socket.id, info);

    if (!this.userSockets.has(userId)) {
      this.userSockets.set(userId, new Set());
    }
    this.userSockets.get(userId)!.add(socket.id);

    // Store in Redis for scaling
    this.syncToRedis(socket.id, info);
  }

  removeConnection(socketId: string): void {
    const info = this.connections.get(socketId);
    if (info) {
      this.connections.delete(socketId);

      const userSockets = this.userSockets.get(info.userId);
      if (userSockets) {
        userSockets.delete(socketId);
        if (userSockets.size === 0) {
          this.userSockets.delete(info.userId);
        }
      }

      this.removeFromRedis(socketId);
    }
  }

  getUserSockets(userId: string): string[] {
    const sockets = this.userSockets.get(userId);
    return sockets ? Array.from(sockets) : [];
  }

  getConnectionCount(): number {
    return this.connections.size;
  }

  getUserCount(): number {
    return this.userSockets.size;
  }

  private async syncToRedis(socketId: string, info: ConnectionInfo): Promise<void> {
    await Redis.hset('ws:connections', socketId, JSON.stringify(info));
    await Redis.sadd('ws:users', info.userId);
  }

  private async removeFromRedis(socketId: string): Promise<void> {
    const info = this.connections.get(socketId);
    if (info) {
      await Redis.hdel('ws:connections', socketId);
    }
  }
}

export const connectionManager = new ConnectionManager();
```

---

## 10) Heartbeat & Health Check

### Health Monitor

```typescript
// src/websocket/health.ts
import { Server } from 'socket.io';
import { logger } from './logger';
import { connectionManager } from './connection-manager';

const HEARTBEAT_INTERVAL = 30000;
const MAX_HEARTBEAT_MISS = 3;

export function setupHealthMonitor(io: Server) {
  const missedHeartbeats = new Map<string, number>();

  setInterval(() => {
    const now = Date.now();
    const sockets = io.sockets.sockets;

    for (const [socketId, socket] of sockets) {
      const lastHeartbeat = socket.data.lastHeartbeat || socket.handshake.time;
      const missed = Math.floor((now - lastHeartbeat) / HEARTBEAT_INTERVAL);

      if (missed >= MAX_HEARTBEAT_MISSED) {
        logger.warn(`Socket ${socketId} missed ${missed} heartbeats, disconnecting`);
        socket.disconnect(true);
        missedHeartbeats.delete(socketId);
      } else if (missed > 0) {
        missedHeartbeats.set(socketId, missed);
      }
    }

    // Emit health metrics
    io.emit('server:health', {
      connections: connectionManager.getConnectionCount(),
      users: connectionManager.getUserCount(),
      timestamp: now,
    });
  }, HEARTBEAT_INTERVAL);

  // Periodic cleanup
  setInterval(() => {
    const keys = missedHeartbeats.keys();
    for (const key of keys) {
      if (!io.sockets.sockets.has(key)) {
        missedHeartbeats.delete(key);
      }
    }
  }, 60000);
}
```

---

## 11) Redis Adapter

### Scaling Adapter

```typescript
// src/websocket/adapters/redis.ts
import { Adapter, SocketId } from 'socket.io-adapter';
import { createAdapter } from 'socket.io-redis';
import { Redis } from '../../services/redis';

const pubClient = Redis.createClient({
  url: `redis://${process.env.REDIS_HOST}:${process.env.REDIS_PORT}`,
  password: process.env.REDIS_PASSWORD,
});

const subClient = pubClient.duplicate();

export const RedisAdapter = createAdapter({
  pubClient,
  subClient,
  parser: require('redis-parser'),
  requestsTimeout: 5000,
});

export class CustomRedisAdapter extends Adapter {
  async publish(room: string, event: string, ...args: any[]): Promise<void> {
    await this.pubClient.publish(`socket:${room}`, JSON.stringify({ event, args }));
  }

  async subscribe(socket: SocketId, room: string): Promise<void> {
    await this.subClient.subscribe(`socket:${room}`, () => {});
    super.subscribe(socket, room);
  }

  async unsubscribe(socket: SocketId, room: string): Promise<void> {
    await super.unsubscribe(socket, room);
  }
}
```

---

## 12) Error Handling

### Error Handler

```typescript
// src/websocket/errors.ts
import { Socket } from 'socket.io';
import { logger } from './logger';

export function setupErrorHandlers(io: Server) {
  io.on('connect_error', (error) => {
    logger.error('Connection error:', error.message);
  });

  io.engine.on('connection_error', (error) => {
    logger.error('Engine connection error:', {
      code: error.code,
      message: error.message,
      context: error.context,
    });
  });

  process.on('uncaughtException', (error) => {
    logger.error('Uncaught exception in WebSocket server:', error);
  });

  process.on('unhandledRejection', (reason) => {
    logger.error('Unhandled rejection in WebSocket server:', reason);
  });
}

export function handleSocketError(socket: Socket, error: Error): void {
  logger.error(`Socket ${socket.id} error:`, error);
  
  socket.emit('error', {
    message: error.message,
    code: error.name,
  });
}
```

---

## 13) Client Library

### Client Integration

```typescript
// src/lib/websocket-client.ts
import { io, Socket } from 'socket.io-client';

class BiometricWebSocket {
  private socket: Socket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;

  connect(token: string): Promise<void> {
    return new Promise((resolve, reject) => {
      this.socket = io(process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:3001', {
        auth: { token },
        transports: ['websocket', 'polling'],
        reconnection: true,
        reconnectionAttempts: this.maxReconnectAttempts,
        reconnectionDelay: 1000,
        reconnectionDelayMax: 5000,
        timeout: 20000,
      });

      this.socket.on('connect', () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
        resolve();
      });

      this.socket.on('connect_error', (error) => {
        console.error('Connection error:', error.message);
        this.reconnectAttempts++;
        
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
          reject(new Error('Max reconnection attempts reached'));
        }
      });

      this.socket.on('disconnect', (reason) => {
        console.log('WebSocket disconnected:', reason);
      });

      this.socket.on('error', (error) => {
        console.error('Socket error:', error);
      });
    });
  }

  // Biometric events
  subscribeToBiometric(recordId: string): void {
    this.socket?.emit('biometric:subscribe', { recordId });
  }

  verifyBiometric(recordId: string, challenge: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.socket?.emit('biometric:verify', { recordId, challenge });
      
      this.socket?.once('biometric:verification-result', (result) => {
        resolve(result);
      });

      this.socket?.once('biometric:error', (error) => {
        reject(new Error(error.message));
      });
    });
  }

  // Session events
  onSessionCreated(callback: (session: any) => void): void {
    this.socket?.on('session:created', callback);
  }

  // Notifications
  onNotification(callback: (notification: any) => void): void {
    this.socket?.on('notification:new', callback);
  }

  // Presence
  updatePresence(status: 'online' | 'away'): void {
    this.socket?.emit('presence:update', { status });
  }

  disconnect(): void {
    this.socket?.disconnect();
    this.socket = null;
  }
}

export const wsClient = new BiometricWebSocket();
```

---

## 14) Security Best Practices

1. **Token Authentication**: All connections authenticated
2. **Rate Limiting**: Connection and message rate limits
3. **Input Validation**: All socket events validated
4. **TLS/SSL**: WSS for production
5. **Origin Validation**: CORS restricted
6. **Message Size Limits**: 100MB max
7. **Heartbeat**: 30-second intervals

---

## 15) Monitoring

### Metrics Collection

```typescript
// src/websocket/metrics.ts
import { Server } from 'socket.io';

interface Metrics {
  connections: number;
  messagesPerSecond: number;
  errorsPerMinute: number;
  averageLatency: number;
}

export function setupMetrics(io: Server) {
  let messageCount = 0;
  let errorCount = 0;
  let latencies: number[] = [];

  io.on('connection', () => {
    messageCount++;
  });

  setInterval(() => {
    const metrics: Metrics = {
      connections: io.engine.clientsCount,
      messagesPerSecond: messageCount / 60,
      errorsPerMinute: errorCount,
      averageLatency: latencies.length > 0 
        ? latencies.reduce((a, b) => a + b) / latencies.length 
        : 0,
    };

    console.log('WebSocket Metrics:', metrics);
    
    // Reset counters
    messageCount = 0;
    errorCount = 0;
    latencies = [];
  }, 60000);
}
```

---

Status: APPROVED  
Version: 1.0  
Last Updated: Februar 2026
