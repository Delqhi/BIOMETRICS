# PWA-CONFIG.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

---

## 1) Overview

This document describes the Progressive Web App (PWA) configuration for BIOMETRICS. The PWA provides installable, offline-capable web experience with native-like features including push notifications, background sync, and home screen installation.

### Core Features
- Installable on desktop and mobile
- Offline functionality with Service Workers
- Push notifications via Web Push
- Background synchronization
- App shell architecture
- Cache strategies
- Web App Manifest

---

## 2) Technology Stack

### Framework & Tools
```json
{
  "next": "15.1.0",
  "typescript": "5.3.3",
  "next-pwa": "^5.6.0",
  "workbox-webpack-plugin": "^7.3.0",
  "web-push": "^3.6.0"
}
```

### Key Dependencies
```json
{
  "dependencies": {
    "@ducanh2912/next-pwa": "^10.2.0",
    "workbox-core": "^7.3.0",
    "workbox-precaching": "^7.3.0",
    "workbox-routing": "^7.3.0",
    "workbox-strategies": "^7.3.0",
    "workbox-background-sync": "^7.3.0",
    "workbox-cacheable-response": "^7.3.0",
    "workbox-expiration": "^7.3.0",
    "web-push": "^3.6.0"
  }
}
```

---

## 3) Project Structure

```
biometrics-pwa/
├── public/
│   ├── manifest.json          # Web App Manifest
│   ├── icons/
│   │   ├── icon-192x192.png
│   │   ├── icon-512x512.png
│   │   ├── icon-maskable.svg
│   │   └── apple-touch-icon.png
│   ├── sw.js                 # Service Worker
│   ├── workbox-*.js          # Workbox generated
│   └── offline.html          # Offline fallback
├── src/
│   ├── app/
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   └── ...
│   ├── components/
│   ├── hooks/
│   │   ├── usePushNotifications.ts
│   │   ├── useOnlineStatus.ts
│   │   └── useBackgroundSync.ts
│   ├── lib/
│   │   ├── push.ts           # Push notification utilities
│   │   ├── sync.ts           # Background sync
│   │   └── caching.ts        # Cache strategies
│   └── styles/
├── next.config.ts
├── tsconfig.json
└── package.json
```

---

## 4) Web App Manifest

### manifest.json

```json
{
  "name": "BIOMETRICS",
  "short_name": "BIOMETRICS",
  "description": "Advanced biometric authentication and identity management platform",
  "start_url": "/",
  "display": "standalone",
  "orientation": "portrait-primary",
  "background_color": "#000000",
  "theme_color": "#00FF00",
  "scope": "/",
  "prefer_related_applications": false,
  "categories": ["business", "productivity", "security"],
  "icons": [
    {
      "src": "/icons/icon-72x72.png",
      "sizes": "72x72",
      "type": "image/png",
      "purpose": "any maskable"
    },
    {
      "src": "/icons/icon-96x96.png",
      "sizes": "96x96",
      "type": "image/png",
      "purpose": "any maskable"
    },
    {
      "src": "/icons/icon-128x128.png",
      "sizes": "128x128",
      "type": "image/png",
      "purpose": "any maskable"
    },
    {
      "src": "/icons/icon-144x144.png",
      "sizes": "144x144",
      "type": "image/png",
      "purpose": "any maskable"
    },
    {
      "src": "/icons/icon-152x152.png",
      "sizes": "152x152",
      "type": "image/png",
      "purpose": "any maskable"
    },
    {
      "src": "/icons/icon-192x192.png",
      "sizes": "192x192",
      "type": "image/png",
      "purpose": "any maskable"
    },
    {
      "src": "/icons/icon-384x384.png",
      "sizes": "384x384",
      "type": "image/png",
      "purpose": "any maskable"
    },
    {
      "src": "/icons/icon-512x512.png",
      "sizes": "512x512",
      "type": "image/png",
      "purpose": "any maskable"
    }
  ],
  "screenshots": [
    {
      "src": "/screenshots/home.png",
      "sizes": "1280x720",
      "type": "image/png",
      "form_factor": "wide"
    },
    {
      "src": "/screenshots/dashboard-mobile.png",
      "sizes": "720x1280",
      "type": "image/png",
      "form_factor": "narrow"
    }
  ],
  "shortcuts": [
    {
      "name": "New Scan",
      "short_name": "Scan",
      "description": "Start a new biometric scan",
      "url": "/scan?action=new",
      "icons": [{ "src": "/icons/scan.png", "sizes": "96x96" }]
    },
    {
      "name": "View History",
      "short_name": "History",
      "description": "View your scan history",
      "url": "/history"
    },
    {
      "name": "Settings",
      "short_name": "Settings",
      "description": "Configure app settings",
      "url": "/settings"
    }
  ],
  "related_applications": [],
  "protocol_handlers": [
    {
      "protocol": "biometrics",
      "url": "/auth?token=%s"
    }
  ]
}
```

---

## 5) Next.js PWA Configuration

### next.config.ts

```typescript
import type { NextConfig } from 'next';
const withPWA = require('@ducanh2912/next-pwa').default({
  dest: 'public',
  register: true,
  skipWaiting: true,
  disable: process.env.NODE_ENV === 'development',
  runtimeCaching: [
    {
      urlPattern: /^https:\/\/fonts\.googleapis\.com\/.*/i,
      handler: 'CacheFirst',
      options: {
        cacheName: 'google-fonts-cache',
        expiration: {
          maxEntries: 10,
          maxAgeSeconds: 60 * 60 * 24 * 365,
        },
        cacheableResponse: {
          statuses: [0, 200],
        },
      },
    },
    {
      urlPattern: /^https:\/\/fonts\.gstatic\.com\/.*/i,
      handler: 'CacheFirst',
      options: {
        cacheName: 'gstatic-fonts-cache',
        expiration: {
          maxEntries: 10,
          maxAgeSeconds: 60 * 60 * 24 * 365,
        },
        cacheableResponse: {
          statuses: [0, 200],
        },
      },
    },
    {
      urlPattern: /\.(?:js|css|woff|woff2|ttf|eot|svg|ico|png|jpg|jpeg|webp|gif|avif)$/,
      handler: 'StaleWhileRevalidate',
      options: {
        cacheName: 'static-resources-cache',
        expiration: {
          maxEntries: 100,
          maxAgeSeconds: 60 * 60 * 24 * 30,
        },
        cacheableResponse: {
          statuses: [0, 200],
        },
      },
    },
    {
      urlPattern: /^https:\/\/api\.biometrics\.app\/.*/i,
      handler: 'NetworkFirst',
      options: {
        cacheName: 'api-cache',
        networkTimeoutSeconds: 10,
        expiration: {
          maxEntries: 100,
          maxAgeSeconds: 60 * 60 * 24,
        },
        cacheableResponse: {
          statuses: [0, 200],
        },
      },
    },
    {
      urlPattern: /^https:\/\/images\.biometrics\.app\/.*/i,
      handler: 'CacheFirst',
      options: {
        cacheName: 'images-cache',
        expiration: {
          maxEntries: 50,
          maxAgeSeconds: 60 * 60 * 24 * 7,
        },
        cacheableResponse: {
          statuses: [0, 200],
        },
      },
    },
  ],
});

const nextConfig: NextConfig = {
  reactStrictMode: true,
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'images.biometrics.app',
      },
      {
        protocol: 'https',
        hostname: 'avatars.githubusercontent.com',
      },
    ],
  },
  headers: async () => [
    {
      source: '/sw.js',
      headers: [
        {
          key: 'Cache-Control',
          value: 'no-cache, no-store, must-revalidate',
        },
        {
          key: 'Service-Worker-Allowed',
          value: '/',
        },
      ],
    },
    {
      source: '/manifest.json',
      headers: [
        {
          key: 'Access-Control-Allow-Origin',
          value: '*',
        },
      ],
    },
  ],
};

export default withPWA(nextConfig);
```

---

## 6) Service Worker Registration

### Registration Component

```typescript
// src/components/PWAProvider.tsx
'use client';

import { useEffect, useState } from 'react';

interface PWADeferredPrompt extends Event {
  userPrompted: boolean;
  prompt: () => Promise<void>;
}

export function PWAProvider({ children }: { children: React.ReactNode }) {
  const [deferredPrompt, setDeferredPrompt] = useState<PWADeferredPrompt | null>(null);
  const [isInstallable, setIsInstallable] = useState(false);
  const [isInstalled, setIsInstalled] = useState(false);

  useEffect(() => {
    // Check if already installed
    const checkInstalled = () => {
      if (window.matchMedia('(display-mode: standalone)').matches) {
        setIsInstalled(true);
      }
    };
    checkInstalled();

    // Listen for beforeinstallprompt
    const handleBeforeInstall = (e: Event) => {
      e.preventDefault();
      setDeferredPrompt(e as PWADeferredPrompt);
      setIsInstallable(true);
    };

    window.addEventListener('beforeinstallprompt', handleBeforeInstall);

    // Listen for successful install
    const handleAppInstalled = () => {
      setIsInstalled(true);
      setIsInstallable(false);
      setDeferredPrompt(null);
    };

    window.addEventListener('appinstalled', handleAppInstalled);

    return () => {
      window.removeEventListener('beforeinstallprompt', handleBeforeInstall);
      window.removeEventListener('appinstalled', handleAppInstalled);
    };
  }, []);

  const handleInstall = async () => {
    if (!deferredPrompt) return;

    deferredPrompt.prompt();
    const { outcome } = await deferredPrompt.userChoice;

    if (outcome === 'accepted') {
      setIsInstalled(true);
    }

    setDeferredPrompt(null);
    setIsInstallable(false);
  };

  return (
    <>
      {children}
      {isInstallable && !isInstalled && (
        <div className="pwa-install-banner">
          <button onClick={handleInstall}>Install BIOMETRICS App</button>
        </div>
      )}
    </>
  );
}

// Custom hook for PWA installation
export function usePWAInstall() {
  const [deferredPrompt, setDeferredPrompt] = useState<PWADeferredPrompt | null>(null);
  const [isInstallable, setIsInstallable] = useState(false);

  useEffect(() => {
    const handleBeforeInstall = (e: Event) => {
      e.preventDefault();
      setDeferredPrompt(e as PWADeferredPrompt);
      setIsInstallable(true);
    };

    window.addEventListener('beforeinstallprompt', handleBeforeInstall);

    return () => {
      window.removeEventListener('beforeinstallprompt', handleBeforeInstall);
    };
  }, []);

  const install = async () => {
    if (!deferredPrompt) return false;

    deferredPrompt.prompt();
    const { outcome } = await deferredPrompt.userChoice;

    if (outcome === 'accepted') {
      return true;
    }

    return false;
  };

  return { isInstallable, install };
}
```

---

## 7) Push Notifications

### Push Notification Hook

```typescript
// src/hooks/usePushNotifications.ts
'use client';

import { useEffect, useState, useCallback } from 'react';

interface PushSubscription {
  endpoint: string;
  keys: {
    p256dh: string;
    auth: string;
  };
}

interface PushNotificationPayload {
  title: string;
  options?: NotificationOptions & {
    data?: any;
    actions?: NotificationAction[];
  };
}

export function usePushNotifications() {
  const [subscription, setSubscription] = useState<PushSubscription | null>(null);
  const [isSupported, setIsSupported] = useState(false);
  const [permission, setPermission] = useState<NotificationPermission>('default');

  useEffect(() => {
    if ('Notification' in window && 'serviceWorker' in navigator) {
      setIsSupported(true);
      setPermission(Notification.permission);
    }
  }, []);

  const requestPermission = useCallback(async (): Promise<boolean> => {
    if (!isSupported) return false;

    const result = await Notification.requestPermission();
    setPermission(result);
    return result === 'granted';
  }, [isSupported]);

  const subscribe = useCallback(async (): Promise<PushSubscription | null> => {
    if (!isSupported || permission !== 'granted') return null;

    const registration = await navigator.serviceWorker.ready;
    const existingSubscription = await registration.pushManager.getSubscription();

    if (existingSubscription) {
      setSubscription(existingSubscription);
      return existingSubscription;
    }

    const vapidPublicKey = process.env.NEXT_PUBLIC_VAPID_PUBLIC_KEY;
    if (!vapidPublicKey) {
      console.error('VAPID public key not configured');
      return null;
    }

    const newSubscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(vapidPublicKey),
    });

    setSubscription(newSubscription);

    // Send subscription to server
    await fetch('/api/push/subscribe', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newSubscription),
    });

    return newSubscription;
  }, [isSupported, permission]);

  const unsubscribe = useCallback(async (): Promise<boolean> => {
    if (!subscription) return false;

    const result = await subscription.unsubscribe();
    
    if (result) {
      // Notify server
      await fetch('/api/push/unsubscribe', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ endpoint: subscription.endpoint }),
      });
      
      setSubscription(null);
    }

    return result;
  }, [subscription]);

  const showLocalNotification = useCallback(
    (title: string, options?: NotificationOptions) => {
      if (permission !== 'granted') return;

      navigator.serviceWorker.ready.then((registration) => {
        registration.showNotification(title, {
          icon: '/icons/icon-192x192.png',
          badge: '/icons/badge-72x72.png',
          ...options,
        });
      });
    },
    [permission]
  );

  return {
    isSupported,
    permission,
    subscription,
    requestPermission,
    subscribe,
    unsubscribe,
    showLocalNotification,
  };
}

// Utility function to convert VAPID key
function urlBase64ToUint8Array(base64String: string): Uint8Array {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');

  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);

  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i);
  }

  return outputArray;
}
```

---

## 8) Service Worker Implementation

### Custom Service Worker

```typescript
// public/sw.js
import { precacheAndRoute, cleanupOutdatedCaches } from 'workbox-precaching';
import { registerRoute } from 'workbox-routing';
import { StaleWhileRevalidate, NetworkFirst, CacheFirst } from 'workbox-strategies';
import { ExpirationPlugin } from 'workbox-expiration';
import { CacheableResponsePlugin } from 'workbox-cacheable-response';
import { BackgroundSyncPlugin } from 'workbox-background-sync';

declare const self: ServiceWorkerGlobalScope;

// Precache static assets
precacheAndRoute(self.__WB_MANIFEST);

// Cleanup old caches
cleanupOutdatedCaches();

// API requests - Network first with fallback
registerRoute(
  ({ url }) => url.pathname.startsWith('/api/'),
  new NetworkFirst({
    cacheName: 'api-cache',
    plugins: [
      new CacheableResponsePlugin({
        statuses: [0, 200],
      }),
      new ExpirationPlugin({
        maxEntries: 100,
        maxAgeSeconds: 24 * 60 * 60,
      }),
    ],
  })
);

// Static assets - Stale while revalidate
registerRoute(
  ({ request }) =>
    request.destination === 'style' ||
    request.destination === 'script' ||
    request.destination === 'image',
  new StaleWhileRevalidate({
    cacheName: 'static-resources',
    plugins: [
      new CacheableResponsePlugin({
        statuses: [0, 200],
      }),
      new ExpirationPlugin({
        maxEntries: 200,
        maxAgeSeconds: 30 * 24 * 60 * 60,
      }),
    ],
  })
);

// Fonts - Cache first
registerRoute(
  ({ request }) => request.destination === 'font',
  new CacheFirst({
    cacheName: 'fonts-cache',
    plugins: [
      new CacheableResponsePlugin({
        statuses: [0, 200],
      }),
      new ExpirationPlugin({
        maxEntries: 10,
        maxAgeSeconds: 365 * 24 * 60 * 60,
      }),
    ],
  })
);

// Background sync for offline form submissions
registerRoute(
  ({ url }) => url.pathname.startsWith('/api/sync/'),
  new NetworkFirst({
    cacheName: 'sync-cache',
    plugins: [
      new BackgroundSyncPlugin('form-submissions', {
        maxRetentionTime: 24 * 60,
      }),
    ],
  })
);

// Push notification click handler
self.addEventListener('push', (event) => {
  if (!event.data) return;

  const data = event.data.json();

  const options: NotificationOptions = {
    body: data.body,
    icon: '/icons/icon-192x192.png',
    badge: '/icons/badge-72x72.png',
    vibrate: [100, 50, 100],
    data: {
      url: data.url || '/',
      dateOfArrival: Date.now(),
    },
    actions: data.actions || [
      { action: 'open', title: 'Open' },
      { action: 'close', title: 'Close' },
    ],
  };

  event.waitUntil(self.registration.showNotification(data.title, options));
});

// Notification click handler
self.addEventListener('notificationclick', (event) => {
  event.notification.close();

  if (event.action === 'open' || !event.action) {
    const urlToOpen = event.notification.data?.url || '/';

    event.waitUntil(
      self.clients.matchAll({ type: 'window', includeUncontrolled: true }).then((clientList) => {
        for (const client of clientList) {
          if ('focus' in client) {
            client.focus();
            return client.navigate(urlToOpen);
          }
        }
        return self.clients.openWindow(urlToOpen);
      })
    );
  }
});

// Message handler for skipping waiting
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
});

// Activate handler to claim clients
self.addEventListener('activate', (event) => {
  event.waitUntil(self.clients.claim());
});
```

---

## 9) Offline Support

### Offline Page

```typescript
// public/offline.html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Offline - BIOMETRICS</title>
  <style>
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      background: #000;
      color: #fff;
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 20px;
    }
    .container {
      text-align: center;
      max-width: 400px;
    }
    .icon {
      width: 120px;
      height: 120px;
      margin-bottom: 24px;
    }
    h1 {
      font-size: 24px;
      margin-bottom: 12px;
      color: #00FF00;
    }
    p {
      font-size: 16px;
      color: #888;
      margin-bottom: 24px;
      line-height: 1.5;
    }
    button {
      background: #00FF00;
      color: #000;
      border: none;
      padding: 12px 24px;
      font-size: 16px;
      font-weight: 600;
      border-radius: 8px;
      cursor: pointer;
      transition: opacity 0.2s;
    }
    button:hover {
      opacity: 0.9;
    }
  </style>
</head>
<body>
  <div class="container">
    <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="#00FF00" stroke-width="1.5">
      <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z"/>
      <path d="M4 12a8 8 0 0 1 8-8"/>
      <path d="M4 12h8"/>
      <circle cx="12" cy="12" r="3"/>
    </svg>
    <h1>You're Offline</h1>
    <p>It looks like you've lost your internet connection. Some features may not be available.</p>
    <button onclick="window.location.reload()">Try Again</button>
  </div>
</body>
</html>
```

---

## 10) Online Status Hook

```typescript
// src/hooks/useOnlineStatus.ts
'use client';

import { useState, useEffect } from 'react';

export function useOnlineStatus() {
  const [isOnline, setIsOnline] = useState(true);

  useEffect(() => {
    setIsOnline(navigator.onLine);

    const handleOnline = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  return isOnline;
}

// Component wrapper for offline handling
export function OfflineWrapper({ children }: { children: React.ReactNode }) {
  const isOnline = useOnlineStatus();

  if (!isOnline) {
    return (
      <div className="offline-notice">
        <span>You are currently offline. Some features may be limited.</span>
      </div>
    );
  }

  return <>{children}</>;
}
```

---

## 11) Background Sync

### Background Sync Hook

```typescript
// src/hooks/useBackgroundSync.ts
'use client';

import { useEffect, useState, useCallback } from 'react';

interface SyncQueueItem {
  id: string;
  type: 'POST' | 'PUT' | 'DELETE';
  url: string;
  data?: any;
  timestamp: number;
}

export function useBackgroundSync() {
  const [isSupported, setIsSupported] = useState(false);
  const [isSyncing, setIsSyncing] = useState(false);
  const [pendingCount, setPendingCount] = useState(0);

  useEffect(() => {
    setIsSupported('sync' in window.ServiceWorkerRegistration.prototype);
  }, []);

  const queueSync = useCallback(
    async (type: SyncQueueItem['type'], url: string, data?: any): Promise<boolean> => {
      const item: SyncQueueItem = {
        id: crypto.randomUUID(),
        type,
        url,
        data,
        timestamp: Date.now(),
      };

      // Store locally first
      const queue = getQueue();
      queue.push(item);
      localStorage.setItem('sync-queue', JSON.stringify(queue));
      setPendingCount(queue.length);

      // Try background sync if supported
      if (isSupported) {
        try {
          const registration = await navigator.serviceWorker.ready;
          await registration.sync.register('sync-data');
          return true;
        } catch (error) {
          console.error('Background sync registration failed:', error);
          return false;
        }
      }

      return false;
    },
    [isSupported]
  );

  const manualSync = useCallback(async () => {
    const queue = getQueue();
    if (queue.length === 0) return;

    setIsSyncing(true);

    for (const item of queue) {
      try {
        const response = await fetch(item.url, {
          method: item.type,
          headers: { 'Content-Type': 'application/json' },
          body: item.data ? JSON.stringify(item.data) : undefined,
        });

        if (response.ok) {
          // Remove from queue on success
          removeFromQueue(item.id);
        }
      } catch (error) {
        console.error('Sync failed for item:', item.id, error);
      }
    }

    setPendingCount(getQueue().length);
    setIsSyncing(false);
  }, []);

  return {
    isSupported,
    isSyncing,
    pendingCount,
    queueSync,
    manualSync,
  };
}

function getQueue(): SyncQueueItem[] {
  try {
    return JSON.parse(localStorage.getItem('sync-queue') || '[]');
  } catch {
    return [];
  }
}

function removeFromQueue(id: string): void {
  const queue = getQueue().filter((item) => item.id !== id);
  localStorage.setItem('sync-queue', JSON.stringify(queue));
}
```

---

## 12) PWA Installation Detection

### Installation Detection

```typescript
// src/hooks/usePWAInstallPrompt.ts
'use client';

import { useState, useEffect } from 'react';

interface BeforeInstallPromptEvent extends Event {
  readonly platforms: string[];
  readonly userChoice: Promise<{
    outcome: 'accepted' | 'dismissed';
    platform: string;
  }>;
  prompt(): Promise<void>;
}

export function usePWAInstallPrompt() {
  const [deferredPrompt, setDeferredPrompt] = useState<BeforeInstallPromptEvent | null>(null);
  const [isInstallable, setIsInstallable] = useState(false);
  const [isInstalled, setIsInstalled] = useState(false);

  useEffect(() => {
    // Check if already in standalone mode
    const checkStandalone = () => {
      if (window.matchMedia('(display-mode: standalone)').matches) {
        setIsInstalled(true);
        return true;
      }
      return false;
    };

    if (checkStandalone()) return;

    const handleBeforeInstall = (e: Event) => {
      e.preventDefault();
      setDeferredPrompt(e as BeforeInstallPromptEvent);
      setIsInstallable(true);
    };

    const handleAppInstalled = () => {
      setIsInstalled(true);
      setIsInstallable(false);
      setDeferredPrompt(null);
    };

    window.addEventListener('beforeinstallprompt', handleBeforeInstall);
    window.addEventListener('appinstalled', handleAppInstalled);

    return () => {
      window.removeEventListener('beforeinstallprompt', handleBeforeInstall);
      window.removeEventListener('appinstalled', handleAppInstalled);
    };
  }, []);

  const install = async (): Promise<boolean> => {
    if (!deferredPrompt) return false;

    deferredPrompt.prompt();
    const { outcome } = await deferredPrompt.userChoice;

    if (outcome === 'accepted') {
      setIsInstalled(true);
      return true;
    }

    return false;
  };

  const dismiss = (): void => {
    setDeferredPrompt(null);
    setIsInstallable(false);
  };

  return {
    isInstallable,
    isInstalled,
    install,
    dismiss,
  };
}
```

---

## 13) App Shell Architecture

### Layout Component

```typescript
// src/app/layout.tsx
import type { Metadata, Viewport } from 'next';
import { PWAProvider } from '@/components/PWAProvider';
import './globals.css';

export const metadata: Metadata = {
  manifest: '/manifest.json',
  title: 'BIOMETRICS',
  description: 'Advanced biometric authentication platform',
  appleWebApp: {
    capable: true,
    statusBarStyle: 'black-translucent',
    title: 'BIOMETRICS',
  },
};

export const viewport: Viewport = {
  themeColor: '#00FF00',
  width: 'device-width',
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <link rel="apple-touch-icon" href="/icons/apple-touch-icon.png" />
        <link rel="apple-touch-startup-image" href="/icons/splash.png" />
      </head>
      <body>
        <PWAProvider>{children}</PWAProvider>
      </body>
    </html>
  );
}
```

---

## 14) Caching Strategies

### Cache Configuration

```typescript
// src/lib/caching.ts
import { CacheStrategy } from './types';

export const cacheStrategies: Record<string, CacheStrategy> = {
  // API calls - Network first
  api: {
    strategy: 'network-first',
    maxAge: 60 * 60 * 24, // 24 hours
    maxEntries: 100,
  },

  // Static assets - Stale while revalidate
  static: {
    strategy: 'stale-while-revalidate',
    maxAge: 60 * 60 * 24 * 30, // 30 days
    maxEntries: 200,
  },

  // Images - Cache first
  images: {
    strategy: 'cache-first',
    maxAge: 60 * 60 * 24 * 7, // 7 days
    maxEntries: 50,
  },

  // Fonts - Cache first with long expiry
  fonts: {
    strategy: 'cache-first',
    maxAge: 60 * 60 * 24 * 365, // 1 year
    maxEntries: 10,
  },

  // User data - Network first with short cache
  userData: {
    strategy: 'network-first',
    maxAge: 60 * 5, // 5 minutes
    maxEntries: 20,
  },
};

export function getCacheStrategy(url: string): CacheStrategy {
  if (url.startsWith('/api/')) return cacheStrategies.api;
  if (url.match(/\.(js|css)$/)) return cacheStrategies.static;
  if (url.match(/\.(png|jpg|jpeg|svg|webp|gif)$/)) return cacheStrategies.images;
  if (url.match(/\.(woff|woff2|ttf|eot)$/)) return cacheStrategies.fonts;
  return cacheStrategies.api;
}
```

---

## 15) IndexedDB Integration

### Offline Storage

```typescript
// src/lib/indexeddb.ts
import { openDB, DBSchema, IDBPDatabase } from 'idb';

interface BiometricsDB extends DBSchema {
  sessions: {
    key: string;
    value: {
      id: string;
      userId: string;
      startedAt: Date;
      endedAt?: Date;
      synced: boolean;
    };
    indexes: { 'by-user': string };
  };
  biometricRecords: {
    key: string;
    value: {
      id: string;
      userId: string;
      type: string;
      data: ArrayBuffer;
      createdAt: Date;
      synced: boolean;
    };
    indexes: { 'by-user': string };
  };
  syncQueue: {
    key: string;
    value: {
      id: string;
      operation: 'add' | 'update' | 'delete';
      store: string;
      data: any;
      timestamp: number;
    };
  };
}

let db: IDBPDatabase<BiometricsDB> | null = null;

export async function initDB(): Promise<IDBPDatabase<BiometricsDB>> {
  if (db) return db;

  db = await openDB<BiometricsDB>('biometrics-offline', 1, {
    upgrade(database) {
      // Sessions store
      const sessionStore = database.createObjectStore('sessions', { keyPath: 'id' });
      sessionStore.createIndex('by-user', 'userId');

      // Biometric records store
      const biometricStore = database.createObjectStore('biometricRecords', {
        keyPath: 'id',
      });
      biometricStore.createIndex('by-user', 'userId');

      // Sync queue store
      database.createObjectStore('syncQueue', { keyPath: 'id' });
    },
  });

  return db;
}

export async function addSession(session: BiometricsDB['sessions']['value']): Promise<void> {
  const database = await initDB();
  await database.put('sessions', session);
}

export async function getUnsyncedSessions(): Promise<BiometricsDB['sessions']['value'][]> {
  const database = await initDB();
  return database.getAllFromIndex('sessions', 'synced', false as any);
}

export async function addToSyncQueue(
  operation: BiometricsDB['syncQueue']['value']['operation'],
  store: string,
  data: any
): Promise<void> {
  const database = await initDB();
  await database.put('syncQueue', {
    id: crypto.randomUUID(),
    operation,
    store,
    data,
    timestamp: Date.now(),
  });
}
```

---

## 16) Security Considerations

1. **HTTPS Only**: Service workers require HTTPS
2. **CSP Headers**: Content Security Policy configured
3. **No Sensitive Data**: Don't cache sensitive info in SW
4. **Token Refresh**: Implement token refresh in sync
5. **Secure Storage**: Use IndexedDB for sensitive data
6. **Update Strategy**: Careful SW update handling

---

Status: APPROVED  
Version: 1.0  
Last Updated: Februar 2026
