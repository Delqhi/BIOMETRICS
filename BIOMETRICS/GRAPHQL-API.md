# GRAPHQL-API.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

---

## 1) Overview

This document describes the GraphQL API architecture for BIOMETRICS. GraphQL provides a flexible, type-safe API layer with efficient data fetching, real-time subscriptions, and strong typing for client-server communication.

### Core Features
- Type-safe schema with code generation
- Real-time subscriptions via WebSocket
- N+1 query prevention with DataLoader
- Rate limiting and caching
- Query complexity analysis
- Automatic persisted queries

---

## 2) Technology Stack

### Framework & Tools
```json
{
  "@apollo/server": "^4.11.0",
  "graphql": "^16.9.0",
  "@graphql-tools/schema": "^10.0.0",
  "dataloader": "^2.2.2",
  "graphql-shield": "^8.0.0",
  "pothos": "^3.0.0"
}
```

### Build Tools
```json
{
  "@graphql-codegen/cli": "^6.0.0",
  "@graphql-codegen/typescript": "^6.0.0",
  "@graphql-codegen/typescript-urql": "^4.0.0",
  "@graphql-codegen/typescript-operations": "^6.0.0"
}
```

---

## 3) Schema Definition

### Main Schema

```graphql
# schema.graphql
scalar DateTime
scalar JSON
scalar Upload

type Query {
  # User queries
  me: User
  user(id: ID!): User
  users(filter: UserFilter, limit: Int, offset: Int): UserConnection!
  
  # Biometric queries
  biometricRecord(id: ID!): BiometricRecord
  biometricRecords(userId: ID!, limit: Int, offset: Int): BiometricRecordConnection!
  biometricDevices(userId: ID!): [BiometricDevice!]!
  
  # Session queries
  session(id: ID!): Session
  sessions(userId: ID!, limit: Int, offset: Int): SessionConnection!
  activeSessions(userId: ID!): [Session!]!
  
  # Analytics queries
  analytics(interval: AnalyticsInterval!, startDate: DateTime!, endDate: DateTime!): Analytics!
  userActivity(userId: ID!, days: Int!): [ActivityPoint!]!
}

type Mutation {
  # Auth mutations
  login(input: LoginInput!): AuthPayload!
  register(input: RegisterInput!): AuthPayload!
  logout: Boolean!
  refreshToken(token: String!): AuthPayload!
  forgotPassword(email: String!): Boolean!
  resetPassword(token: String!, password: String!): Boolean!
  
  # User mutations
  updateUser(input: UpdateUserInput!): User!
  deleteUser(id: ID!): Boolean!
  
  # Biometric mutations
  enrollBiometric(input: EnrollBiometricInput!): BiometricRecord!
  updateBiometric(id: ID!, input: UpdateBiometricInput!): BiometricRecord!
  deleteBiometric(id: ID!): Boolean!
  verifyBiometric(input: VerifyBiometricInput!): BiometricVerificationResult!
  
  # Session mutations
  createSession(input: CreateSessionInput!): Session!
  endSession(id: ID!): Session!
  
  # Admin mutations
  createApiKey(input: CreateApiKeyInput!): ApiKey!
  revokeApiKey(id: ID!): Boolean!
}

type Subscription {
  # Real-time events
  biometricEvent(userId: ID!): BiometricEvent!
  sessionEvent(userId: ID!): SessionEvent!
  notification(userId: ID!): Notification!
}

# Types
type User {
  id: ID!
  email: String!
  name: String
  avatar: String
  role: UserRole!
  createdAt: DateTime!
  updatedAt: DateTime!
  biometricRecords: [BiometricRecord!]!
  sessions: [Session!]!
  settings: UserSettings!
}

type BiometricRecord {
  id: ID!
  user: User!
  device: BiometricDevice
  type: BiometricType!
  template: String!
  qualityScore: Float
  enrolledAt: DateTime!
  lastUsed: DateTime
  isActive: Boolean!
  metadata: JSON
}

type BiometricDevice {
  id: ID!
  name: String!
  type: DeviceType!
  manufacturer: String
  model: String
  isTrusted: Boolean!
  lastSeen: DateTime!
}

type Session {
  id: ID!
  user: User!
  deviceInfo: String!
  ipAddress: String!
  startedAt: DateTime!
  endedAt: DateTime
  isActive: Boolean!
  location: String
}

type AuthPayload {
  accessToken: String!
  refreshToken: String!
  user: User!
  expiresIn: Int!
}

type Analytics {
  totalUsers: Int!
  totalBiometricEnrollments: Int!
  totalSessions: Int!
  successRate: Float!
  averageSessionDuration: Float!
  dailyBreakdown: [DailyMetric!]!
  topDevices: [DeviceMetric!]!
}

# Inputs
input LoginInput {
  email: String!
  password: String!
  deviceInfo: String
}

input RegisterInput {
  email: String!
  password: String!
  name: String
}

input EnrollBiometricInput {
  type: BiometricType!
  deviceId: ID
  template: String!
  metadata: JSON
}

input VerifyBiometricInput {
  biometricRecordId: ID!
  verificationData: String!
}

input UpdateUserInput {
  name: String
  avatar: String
  settings: UserSettingsInput
}

# Enums
enum UserRole {
  USER
  ADMIN
  SERVICE
}

enum BiometricType {
  FINGERPRINT
  FACE
  IRIS
  VOICE
  MULTI
}

enum DeviceType {
  MOBILE
  DESKTOP
  TABLET
  IOT
}

enum AnalyticsInterval {
  DAILY
  WEEKLY
  MONTHLY
  YEARLY
}

# Connections
interface Connection {
  edges: [Edge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

interface Edge {
  node: Node!
  cursor: String!
}

interface Node {
  id: ID!
}

type UserConnection implements Connection
type BiometricRecordConnection implements Connection
type SessionConnection implements Connection

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}
```

---

## 4) Apollo Server Setup

### Server Configuration

```typescript
// src/graphql/server.ts
import { ApolloServer } from '@apollo/server';
import { expressMiddleware } from '@apollo/server/express4';
import { ApolloServerPluginDrainHttpServer } from '@apollo/server/plugin/drainHttpServer';
import { makeExecutableSchema } from '@graphql-tools/schema';
import { WebSocketServer } from 'ws';
import { useServer } from 'graphql-ws/lib/use/ws';
import express from 'express';
import { createServer } from 'http';
import { authenticate } from './auth';
import {授权} from './resolvers/authorization';
import { queries } from './resolvers/queries';
import { mutations } from './resolvers/mutations';
import { subscriptions } from './resolvers/subscriptions';
import { typeDefs } from './schema';
import { logger } from './logger';

const app = express();
const httpServer = createServer(app);

const schema = makeExecutableSchema({ typeDefs, resolvers });

const wsServer = new WebSocketServer({
  server: httpServer,
  path: '/graphql',
});

const serverCleanup = useServer(
  {
    schema,
    context: async (ctx) => {
      const token = ctx.connectionParams?.authorization;
      const user = await authenticate(token);
      return { user };
    },
    onConnect: async (ctx) => {
      const token = ctx.connectionParams?.authorization;
      return !!token;
    },
    onDisconnect: (ctx) => {
      logger.info('Client disconnected', { ctx });
    },
  },
  wsServer
);

const server = new ApolloServer({
  schema,
  plugins: [
    ApolloServerPluginDrainHttpServer({ httpServer }),
    {
      async requestDidStart() {
        return {
          async willSendResponse({ response }) {
            logger.info('GraphQL request completed', {
              status: response.http?.status,
            });
          },
        };
      },
    },
  ],
  formatError: (error) => {
    logger.error('GraphQL Error:', error);
    return {
      message: error.message,
      locations: error.locations,
      path: error.path,
    };
  },
});

await server.start();

app.use(
  '/graphql',
  express.json(),
  expressMiddleware(server, {
    context: async ({ req }) => {
      const token = req.headers.authorization;
      const user = await authenticate(token);
      return { user };
    },
  })
);

await new Promise<void>((resolve) => httpServer.listen({ port: 4000 }, resolve));

logger.info('GraphQL server ready at http://localhost:4000/graphql');
```

---

## 5) DataLoader Implementation

### N+1 Prevention

```typescript
// src/graphql/dataloaders.ts
import DataLoader from 'dataloader';
import { prisma } from '../db';

export const createDataLoaders = () => ({
  userLoader: new DataLoader(async (userIds: string[]) => {
    const users = await prisma.user.findMany({
      where: { id: { in: userIds } },
    });

    const userMap = new Map(users.map((user) => [user.id, user]));
    return userIds.map((id) => userMap.get(id) || null);
  }),

  biometricRecordLoader: new DataLoader(async (recordIds: string[]) => {
    const records = await prisma.biometricRecord.findMany({
      where: { id: { in: recordIds } },
      include: { user: true, device: true },
    });

    const recordMap = new Map(records.map((record) => [record.id, record]));
    return recordIds.map((id) => recordMap.get(id) || null);
  }),

  sessionLoader: new DataLoader(async (sessionIds: string[]) => {
    const sessions = await prisma.session.findMany({
      where: { id: { in: sessionIds } },
      include: { user: true },
    });

    const sessionMap = new Map(sessions.map((session) => [session.id, session]));
    return sessionIds.map((id) => sessionMap.get(id) || null);
  }),

  userBiometricRecordsLoader: new DataLoader(async (userIds: string[]) => {
    const records = await prisma.biometricRecord.findMany({
      where: { userId: { in: userIds } },
    });

    const recordsByUser = new Map<string, typeof records>();
    for (const userId of userIds) {
      recordsByUser.set(
        userId,
        records.filter((r) => r.userId === userId)
      );
    }

    return userIds.map((id) => recordsByUser.get(id) || []);
  }),

  userSessionsLoader: new DataLoader(async (userIds: string[]) => {
    const sessions = await prisma.session.findMany({
      where: { userId: { in: userIds } },
    });

    const sessionsByUser = new Map<string, typeof sessions>();
    for (const userId of userIds) {
      sessionsByUser.set(
        userId,
        sessions.filter((s) => s.userId === userId)
      );
    }

    return userIds.map((id) => sessionsByUser.get(id) || []);
  }),
});
```

---

## 6) Authentication Resolver

### Auth Implementation

```typescript
// src/graphql/resolvers/auth.ts
import { GraphQLError } from 'graphql';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import { prisma } from '../../db';
import { JWT_SECRET, JWT_EXPIRES_IN } from '../../config';

export const authResolvers = {
  Mutation: {
    login: async (_: any, { input }: { input: LoginInput }) => {
      const user = await prisma.user.findUnique({
        where: { email: input.email },
      });

      if (!user || !user.password) {
        throw new GraphQLError('Invalid credentials', {
          extensions: { code: 'UNAUTHENTICATED' },
        });
      }

      const isValid = await bcrypt.compare(input.password, user.password);
      if (!isValid) {
        throw new GraphQLError('Invalid credentials', {
          extensions: { code: 'UNAUTHENTICATED' },
        });
      }

      const accessToken = jwt.sign(
        { userId: user.id, role: user.role },
        JWT_SECRET,
        { expiresIn: '15m' }
      );

      const refreshToken = jwt.sign(
        { userId: user.id, type: 'refresh' },
        JWT_SECRET,
        { expiresIn: '7d' }
      );

      return {
        accessToken,
        refreshToken,
        user,
        expiresIn: 900,
      };
    },

    register: async (_: any, { input }: { input: RegisterInput }) => {
      const existingUser = await prisma.user.findUnique({
        where: { email: input.email },
      });

      if (existingUser) {
        throw new GraphQLError('Email already in use', {
          extensions: { code: 'BAD_USER_INPUT' },
        });
      }

      const hashedPassword = await bcrypt.hash(input.password, 12);

      const user = await prisma.user.create({
        data: {
          email: input.email,
          password: hashedPassword,
          name: input.name,
        },
      });

      const accessToken = jwt.sign(
        { userId: user.id, role: user.role },
        JWT_SECRET,
        { expiresIn: '15m' }
      );

      const refreshToken = jwt.sign(
        { userId: user.id, type: 'refresh' },
        JWT_SECRET,
        { expiresIn: '7d' }
      );

      return {
        accessToken,
        refreshToken,
        user,
        expiresIn: 900,
      };
    },

    refreshToken: async (_: any, { token }: { token: string }) => {
      try {
        const decoded = jwt.verify(token, JWT_SECRET) as any;

        if (decoded.type !== 'refresh') {
          throw new GraphQLError('Invalid token type');
        }

        const user = await prisma.user.findUnique({
          where: { id: decoded.userId },
        });

        if (!user) {
          throw new GraphQLError('User not found');
        }

        const accessToken = jwt.sign(
          { userId: user.id, role: user.role },
          JWT_SECRET,
          { expiresIn: '15m' }
        );

        const newRefreshToken = jwt.sign(
          { userId: user.id, type: 'refresh' },
          JWT_SECRET,
          { expiresIn: '7d' }
        );

        return {
          accessToken,
          refreshToken: newRefreshToken,
          user,
          expiresIn: 900,
        };
      } catch (error) {
        throw new GraphQLError('Invalid refresh token', {
          extensions: { code: 'UNAUTHENTICATED' },
        });
      }
    },

    logout: async (_: any, __: any, { user }: Context) => {
      return true;
    },
  },
};
```

---

## 7) Biometric Resolvers

### Biometric Operations

```typescript
// src/graphql/resolvers/biometric.ts
import { GraphQLError } from 'graphql';
import { prisma } from '../../db';
import { verifyBiometricTemplate } from '../../services/biometric';

export const biometricResolvers = {
  Query: {
    biometricRecord: async (_: any, { id }: { id: string }, { user }: Context) => {
      const record = await prisma.biometricRecord.findUnique({
        where: { id },
        include: { user: true, device: true },
      });

      if (!record) {
        throw new GraphQLError('Biometric record not found');
      }

      if (record.userId !== user.id && user.role !== 'ADMIN') {
        throw new GraphQLError('Not authorized');
      }

      return record;
    },

    biometricRecords: async (
      _: any,
      { userId, limit = 20, offset = 0 }: { userId: string; limit?: number; offset?: number },
      { user }: Context
    ) => {
      if (user.id !== userId && user.role !== 'ADMIN') {
        throw new GraphQLError('Not authorized');
      }

      const [records, totalCount] = await Promise.all([
        prisma.biometricRecord.findMany({
          where: { userId },
          take: limit,
          skip: offset,
          orderBy: { enrolledAt: 'desc' },
        }),
        prisma.biometricRecord.count({ where: { userId } }),
      ]);

      return {
        edges: records.map((record) => ({
          node: record,
          cursor: Buffer.from(record.id).toString('base64'),
        })),
        pageInfo: {
          hasNextPage: offset + records.length < totalCount,
          hasPreviousPage: offset > 0,
          startCursor: records[0]?.id,
          endCursor: records[records.length - 1]?.id,
        },
        totalCount,
      };
    },
  },

  Mutation: {
    enrollBiometric: async (
      _: any,
      { input }: { input: EnrollBiometricInput },
      { user }: Context
    ) => {
      const template = Buffer.from(input.template, 'base64');

      const record = await prisma.biometricRecord.create({
        data: {
          userId: user.id,
          deviceId: input.deviceId,
          type: input.type,
          template: template.toString('base64'),
          metadata: input.metadata,
          qualityScore: input.metadata?.qualityScore || 0.8,
        },
        include: { user: true, device: true },
      });

      return record;
    },

    verifyBiometric: async (
      _: any,
      { input }: { input: VerifyBiometricInput },
      { user }: Context
    ) => {
      const record = await prisma.biometricRecord.findUnique({
        where: { id: input.biometricRecordId },
      });

      if (!record) {
        throw new GraphQLError('Biometric record not found');
      }

      if (record.userId !== user.id) {
        throw new GraphQLError('Not authorized');
      }

      const isValid = await verifyBiometricTemplate(
        record.template,
        input.verificationData
      );

      await prisma.biometricRecord.update({
        where: { id: record.id },
        data: { lastUsed: new Date() },
      });

      return {
        success: isValid,
        record,
        confidence: isValid ? 0.95 : 0,
      };
    },

    deleteBiometric: async (_: any, { id }: { id: string }, { user }: Context) => {
      const record = await prisma.biometricRecord.findUnique({
        where: { id },
      });

      if (!record) {
        throw new GraphQLError('Biometric record not found');
      }

      if (record.userId !== user.id && user.role !== 'ADMIN') {
        throw new GraphQLError('Not authorized');
      }

      await prisma.biometricRecord.delete({
        where: { id },
      });

      return true;
    },
  },
};
```

---

## 8) Real-time Subscriptions

### Subscription Handlers

```typescript
// src/graphql/resolvers/subscriptions.ts
import { PubSub } from 'graphql-subscriptions';

export const pubsub = new PubSub();

const BIOMETRIC_EVENT = 'BIOMETRIC_EVENT';
const SESSION_EVENT = 'SESSION_EVENT';
const NOTIFICATION = 'NOTIFICATION';

export const subscriptionResolvers = {
  Subscription: {
    biometricEvent: {
      subscribe: async function* (_: any, { userId }: { userId: string }) {
        const channel = `${BIOMETRIC_EVENT}:${userId}`;
        for await (const event of pubsub.asyncIterator([channel])) {
          yield { biometricEvent: event };
        }
      },
    },

    sessionEvent: {
      subscribe: async function* (_: any, { userId }: { userId: string }) {
        const channel = `${SESSION_EVENT}:${userId}`;
        for await (const event of pubsub.asyncIterator([channel])) {
          yield { sessionEvent: event };
        }
      },
    },

    notification: {
      subscribe: async function* (_: any, { userId }: { userId: string }) {
        const channel = `${NOTIFICATION}:${userId}`;
        for await (const event of pubsub.asyncIterator([channel])) {
          yield { notification: event };
        }
      },
    },
  },
};

// Helper functions to publish events
export const publishBiometricEvent = (userId: string, event: any) => {
  pubsub.publish(`${BIOMETRIC_EVENT}:${userId}`, event);
};

export const publishSessionEvent = (userId: string, event: any) => {
  pubsub.publish(`${SESSION_EVENT}:${userId}`, event);
};

export const publishNotification = (userId: string, notification: any) => {
  pubsub.publish(`${NOTIFICATION}:${userId}`, notification);
};
```

---

## 9) Rate Limiting

### Rate Limiter

```typescript
// src/graphql/rateLimit.ts
import { rateLimit } from 'graphql-rate-limit';
import { GraphQLError } from 'graphql';

const fieldExtensions = {
  complexity: (args: any) => {
    let complexity = 1;

    if (args.limit) {
      complexity = Math.min(args.limit, 100);
    }

    return complexity;
  },
};

const rateLimitDirective = rateLimit(
  {
    identifyContext: (context) => context.user?.id || context.ip,
    resultsLoaderFactory: (loader) => loader,
    maximum: 100,
    window: '1minute',
    onLimit: (current, limit, window) => {
      throw new GraphQLError(
        `Rate limit exceeded. Try again in ${window}`,
        {
          extensions: {
            code: 'RATE_LIMITED',
            retryAfter: window,
          },
        }
      );
    },
  },
  {
    fieldExtensions,
  }
);

export { rateLimitDirective };
```

---

## 10) Code Generation

### Codegen Config

```typescript
// codegen.ts
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  overwrite: true,
  schema: './src/graphql/schema.ts',
  documents: './src/**/*.graphql',
  generates: {
    './src/graphql/generated/graphql.ts': {
      plugins: [
        'typescript',
        'typescript-operations',
        'typescript-urql',
      ],
      config: {
        withHooks: true,
        withComponent: false,
        withHOC: false,
        scalars: {
          DateTime: 'string',
          JSON: 'object',
          Upload: 'File',
        },
      },
    },
    './src/graphql/generated/schema.graphql': {
      plugins: ['schema-ast'],
    },
  },
  hooks: {
    afterAllFileWrite: ['prettier --write'],
  },
};

export default config;
```

---

## 11) Query Complexity Analysis

### Complexity Calculator

```typescript
// src/graphql/complexity.ts
import { getComplexity, simpleEstimator, fieldExtensionsEstimator } from 'graphql-query-complexity';

export const queryComplexityPlugin = {
  onParseComplete({ schema, document }) {
    const complexity = getComplexity({
      schema,
      document,
      estimators: [
        fieldExtensionsEstimator(),
        simpleEstimator({ defaultComplexity: 1 }),
      ],
      variables: {},
      operationName: '',
    });

    if (complexity > 100) {
      throw new Error(
        `Query is too complex: ${complexity}. Maximum allowed complexity is 100.`
      );
    }

    return complexity;
  },
};
```

---

## 12) Caching Strategy

### Response Caching

```typescript
// src/graphql/cache.ts
import { InMemoryCache, NormalizedCacheObject } from '@apollo/client';
import { CachePersistor } from 'apollo-cache-persist';

const cache: InMemoryCache = new InMemoryCache({
  typePolicies: {
    Query: {
      fields: {
        users: {
          keyArgs: ['filter'],
          merge(existing = { edges: [] }, incoming) {
            return {
              ...incoming,
              edges: [...existing.edges, ...incoming.edges],
            };
          },
        },
        biometricRecords: {
          keyArgs: ['userId'],
          merge(existing = { edges: [] }, incoming) {
            return {
              ...incoming,
              edges: [...existing.edges, ...incoming.edges],
            };
          },
        },
      },
    },
    User: {
      keyFields: ['id'],
    },
    BiometricRecord: {
      keyFields: ['id'],
    },
    Session: {
      keyFields: ['id'],
    },
  },
});

// For offline support
export const setupCachePersistor = async (storage: Storage) => {
  const persistor = new CachePersistor({
    cache,
    storage: storage as any,
    maxSize: 5242880, // 5MB
    debug: process.env.NODE_ENV === 'development',
  });

  await persistor.restore();

  return { cache, persistor };
};

export { cache };
```

---

## 13) Error Handling

### Error Formatter

```typescript
// src/graphql/errors.ts
import { GraphQLError, GraphQLFormattedError } from 'graphql';

interface ErrorExtension {
  code: string;
  statusCode?: number;
  details?: any;
}

export const formatError = (
  error: GraphQLFormattedError<ErrorExtension>
): GraphQLFormattedError<ErrorExtension> => {
  const originalError = error.originalError;

  if (originalError) {
    return {
      message: error.message,
      locations: error.locations,
      path: error.path,
      extensions: {
        code: originalError.extensions?.code || 'INTERNAL_SERVER_ERROR',
        statusCode: originalError.extensions?.statusCode || 500,
        details: originalError.extensions?.details,
      },
    };
  }

  return {
    message: error.message,
    locations: error.locations,
    path: error.path,
    extensions: {
      code: 'INTERNAL_SERVER_ERROR',
      statusCode: 500,
    },
  };
};

export class AppError extends GraphQLError {
  constructor(message: string, code: string, statusCode: number = 500) {
    super(message, undefined, undefined, undefined, undefined, undefined, {
      code,
      statusCode,
    });
  }
}

export class NotFoundError extends AppError {
  constructor(message: string = 'Resource not found') {
    super(message, 'NOT_FOUND', 404);
  }
}

export class UnauthorizedError extends AppError {
  constructor(message: string = 'Unauthorized') {
    super(message, 'UNAUTHENTICATED', 401);
  }
}

export class ForbiddenError extends AppError {
  constructor(message: string = 'Forbidden') {
    super(message, 'FORBIDDEN', 403);
  }
}
```

---

## 14) Client Integration

### Apollo Client Setup

```typescript
// src/lib/apollo-client.ts
import { ApolloClient, InMemoryCache, createHttpLink, split } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import { getMainDefinition } from '@apollo/client/utilities';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { createClient } from 'graphql-ws';
import { cache } from './cache';

const httpLink = createHttpLink({
  uri: process.env.NEXT_PUBLIC_GRAPHQL_URL || 'http://localhost:4000/graphql',
});

const wsLink = typeof window !== 'undefined' 
  ? new GraphQLWsLink(
      createClient({
        url: process.env.NEXT_PUBLIC_GRAPHQL_WS_URL || 'ws://localhost:4000/graphql',
        connectionParams: () => {
          const token = localStorage.getItem('accessToken');
          return { authorization: token ? `Bearer ${token}` : '' };
        },
      })
    )
  : null;

const splitLink = wsLink
  ? split(
      ({ query }) => {
        const definition = getMainDefinition(query);
        return (
          definition.kind === 'OperationDefinition' &&
          definition.operation === 'subscription'
        );
      },
      wsLink,
      httpLink
    )
  : httpLink;

const authLink = setContext((_, { headers }) => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('accessToken') : '';
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : '',
    },
  };
});

export const apolloClient = new ApolloClient({
  link: authLink.concat(splitLink),
  cache,
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'cache-and-network',
    },
    query: {
      fetchPolicy: 'cache-first',
    },
  },
});
```

---

## 15) Security Best Practices

1. **Query Depth Limiting**: Maximum query depth of 10
2. **Cost Analysis**: Block expensive queries
3. **Introspection**: Disable in production
4. **CSRF Protection**: Proper token validation
5. **Input Validation**: All inputs sanitized
6. **Rate Limiting**: Per-user rate limits

---

Status: APPROVED  
Version: 1.0  
Last Updated: Februar 2026
