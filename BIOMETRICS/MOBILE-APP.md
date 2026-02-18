# MOBILE-APP.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

---

## 1) Overview

This document describes the React Native mobile application architecture for BIOMETRICS. The mobile app provides biometric authentication, user management, and real-time data synchronization with the backend services.

### Core Features
- Biometric enrollment and verification (Face ID, Touch ID, Fingerprint)
- Secure credential storage using device Keychain
- Real-time sync with Supabase backend
- Offline-first architecture with conflict resolution
- Push notifications for authentication events
- Multi-language support (i18n)

---

## 2) Technology Stack

### Framework & Language
```json
{
  "react-native": "0.76.9",
  "typescript": "5.3.3",
  "expo": "52.0.0"
}
```

### Key Dependencies
```json
{
  "dependencies": {
    "@react-navigation/native": "^7.0.0",
    "@react-navigation/native-stack": "^7.0.0",
    "@react-navigation/bottom-tabs": "^7.0.0",
    "react-native-keychain": "^9.2.0",
    "expo-local-authentication": "~15.0.0",
    "@supabase/supabase-js": "^2.45.0",
    "react-native-mmkv": "^4.0.0",
    "react-native-nitro-modules": "^0.22.0",
    "zustand": "^5.0.0",
    "react-native-worklets": "^1.0.0",
    "expo-notifications": "~0.29.0",
    "expo-splash-screen": "~0.29.0"
  }
}
```

---

## 3) Project Structure

```
biometrics-mobile/
├── src/
│   ├── app/                    # Expo Router file-based routing
│   │   ├── (auth)/            # Authentication flow
│   │   │   ├── login.tsx
│   │   │   ├── register.tsx
│   │   │   └── biometric-enroll.tsx
│   │   ├── (main)/            # Main app flow
│   │   │   ├── _layout.tsx
│   │   │   ├── home.tsx
│   │   │   ├── profile.tsx
│   │   │   └── settings.tsx
│   │   └── +layout.tsx        # Root layout
│   ├── components/
│   │   ├── ui/                # Reusable UI components
│   │   │   ├── Button.tsx
│   │   │   ├── Input.tsx
│   │   │   ├── Card.tsx
│   │   │   └── Loading.tsx
│   │   ├── biometric/
│   │   │   ├── BiometricPrompt.tsx
│   │   │   ├── FaceIdScanner.tsx
│   │   │   └── FingerprintSensor.tsx
│   │   └── auth/
│   │       ├── AuthGuard.tsx
│   │       └── SessionProvider.tsx
│   ├── lib/
│   │   ├── supabase.ts        # Supabase client
│   │   ├── keychain.ts        # Secure storage
│   │   ├── biometrics.ts      # Biometric API wrapper
│   │   └── sync.ts           # Offline sync logic
│   ├── stores/
│   │   ├── authStore.ts       # Zustand auth state
│   │   ├── settingsStore.ts   # App settings
│   │   └── syncStore.ts      # Sync state
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   ├── useBiometrics.ts
│   │   └── useOfflineSync.ts
│   ├── services/
│   │   ├── api.ts            # REST API client
│   │   ├── websocket.ts       # Real-time updates
│   │   └── notifications.ts   # Push notifications
│   ├── utils/
│   │   ├── crypto.ts         # Encryption utilities
│   │   ├── logger.ts         # Logging
│   │   └── validation.ts     # Input validation
│   ├── types/
│   │   └── index.ts          # TypeScript types
│   └── i18n/
│       ├── en.json
│       └── de.json
├── android/
├── ios/
├── app.json
├── babel.config.js
├── metro.config.js
├── tsconfig.json
└── package.json
```

---

## 4) Biometric Authentication Implementation

### BiometricService Class

```typescript
// src/lib/biometrics.ts
import * as LocalAuthentication from 'expo-local-authentication';
import * as Keychain from 'react-native-keychain';

export interface BiometricResult {
  success: boolean;
  error?: string;
  biometricType?: 'facial' | 'fingerprint' | 'iris';
}

export interface EnrolledBiometric {
  id: string;
  type: 'facial' | 'fingerprint';
  enrolledAt: Date;
  lastUsed: Date;
}

class BiometricService {
  private readonly SERVICE_NAME = 'com.biometrics.app';

  async checkAvailability(): Promise<{
    available: boolean;
    biometricTypes: string[];
  }> {
    const compatible = await LocalAuthentication.hasHardwareAsync();
    const enrolled = await LocalAuthentication.isEnrolledAsync();
    const types = await LocalAuthentication.supportedAuthenticationTypesAsync();

    return {
      available: compatible && enrolled,
      biometricTypes: types.map((type) => {
        switch (type) {
          case LocalAuthentication.AuthenticationType.FACIAL_RECOGNITION:
            return 'facial';
          case LocalAuthentication.AuthenticationType.FINGERPRINT:
            return 'fingerprint';
          case LocalAuthentication.AuthenticationType.IRIS:
            return 'iris';
          default:
            return 'unknown';
        }
      }),
    };
  }

  async authenticate(reason: string): Promise<BiometricResult> {
    try {
      const result = await LocalAuthentication.authenticateAsync({
        promptMessage: reason,
        fallbackLabel: 'Use Passcode',
        cancelLabel: 'Cancel',
        disableDeviceFallback: false,
        fallback: 'passcode',
      });

      return {
        success: result.success,
        error: result.error,
        biometricType: result.authenticationType 
          ? this.mapAuthType(result.authenticationType)
          : undefined,
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      };
    }
  }

  async enrollBiometric(userId: string): Promise<BiometricResult> {
    const result = await this.authenticate('Enroll your biometric credential');
    
    if (result.success) {
      await this.storeBiometricCredential(userId, result.biometricType!);
    }
    
    return result;
  }

  private async storeBiometricCredential(
    userId: string,
    biometricType: string
  ): Promise<void> {
    const credential = {
      userId,
      biometricType,
      enrolledAt: new Date().toISOString(),
    };

    await Keychain.setGenericPassword(
      `biometric_${userId}`,
      JSON.stringify(credential),
      {
        service: this.SERVICE_NAME,
        accessControl: Keychain.ACCESS_CONTROL.BIOMETRY_ANY,
        accessible: Keychain.ACCESSIBLE.WHEN_PASSCODE_SET_THIS_DEVICE_ONLY,
      }
    );
  }

  private mapAuthType(type: number): 'facial' | 'fingerprint' | 'iris' {
    switch (type) {
      case LocalAuthentication.AuthenticationType.FACIAL_RECOGNITION:
        return 'facial';
      case LocalAuthentication.AuthenticationType.FINGERPRINT:
        return 'fingerprint';
      case LocalAuthentication.AuthenticationType.IRIS:
        return 'iris';
      default:
        return 'fingerprint';
    }
  }
}

export const biometricService = new BiometricService();
```

---

## 5) Secure Storage Implementation

### KeychainWrapper

```typescript
// src/lib/keychain.ts
import * as Keychain from 'react-native-keychain';

export interface SecureCredential {
  username: string;
  password: string;
  accessToken?: string;
  refreshToken?: string;
}

class KeychainWrapper {
  private readonly SERVICE_NAME = 'com.biometrics.secure';

  async storeCredentials(credentials: SecureCredential): Promise<boolean> {
    try {
      await Keychain.setGenericPassword(
        credentials.username,
        JSON.stringify({
          password: credentials.password,
          accessToken: credentials.accessToken,
          refreshToken: credentials.refreshToken,
        }),
        {
          service: this.SERVICE_NAME,
          accessControl: Keychain.ACCESS_CONTROL.BIOMETRY_ANY_OR_DEVICE_PASSCODE,
          accessible: Keychain.ACCESSIBLE.WHEN_UNLOCKED_THIS_DEVICE_ONLY,
        }
      );
      return true;
    } catch (error) {
      console.error('Failed to store credentials:', error);
      return false;
    }
  }

  async getCredentials(): Promise<SecureCredential | null> {
    try {
      const credentials = await Keychain.getGenericPassword({
        service: this.SERVICE_NAME,
        authenticationPrompt: {
          title: 'Authenticate',
          subtitle: 'Access your secure data',
          description: 'Use biometric to access credentials',
        },
      });

      if (credentials) {
        return {
          username: credentials.username,
          ...JSON.parse(credentials.password),
        };
      }
      return null;
    } catch (error) {
      console.error('Failed to get credentials:', error);
      return null;
    }
  }

  async deleteCredentials(): Promise<boolean> {
    try {
      await Keychain.resetGenericPassword({
        service: this.SERVICE_NAME,
      });
      return true;
    } catch (error) {
      console.error('Failed to delete credentials:', error);
      return false;
    }
  }

  async isBiometricAvailable(): Promise<boolean> {
    const biometryType = await Keychain.getSupportedBiometryType();
    return biometryType !== null;
  }
}

export const keychainWrapper = new KeychainWrapper();
```

---

## 6) Offline-First Data Sync

### SyncManager

```typescript
// src/lib/sync.ts
import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { supabase } from './supabase';

export interface SyncQueueItem {
  id: string;
  table: string;
  operation: 'create' | 'update' | 'delete';
  data: Record<string, unknown>;
  timestamp: number;
  retryCount: number;
}

interface SyncState {
  isOnline: boolean;
  lastSyncTime: number | null;
  pendingChanges: SyncQueueItem[];
  conflictResolution: 'server-wins' | 'client-wins' | 'manual';
}

interface SyncActions {
  setOnline: (online: boolean) => void;
  addToQueue: (item: Omit<SyncQueueItem, 'id' | 'timestamp' | 'retryCount'>) => void;
  processQueue: () => Promise<void>;
  resolveConflict: (itemId: string, resolution: 'server' | 'client') => Promise<void>;
}

const useSyncStore = create<SyncState & SyncActions>()(
  persist(
    (set, get) => ({
      isOnline: true,
      lastSyncTime: null,
      pendingChanges: [],
      conflictResolution: 'server-wins',

      setOnline: (online) => set({ isOnline: online }),

      addToQueue: (item) =>
        set((state) => ({
          pendingChanges: [
            ...state.pendingChanges,
            {
              ...item,
              id: crypto.randomUUID(),
              timestamp: Date.now(),
              retryCount: 0,
            },
          ],
        })),

      processQueue: async () => {
        const { pendingChanges, isOnline } = get();
        
        if (!isOnline || pendingChanges.length === 0) return;

        const failedItems: SyncQueueItem[] = [];

        for (const item of pendingChanges) {
          try {
            let result;
            
            switch (item.operation) {
              case 'create':
                result = await supabase.from(item.table).insert(item.data);
                break;
              case 'update':
                result = await supabase
                  .from(item.table)
                  .update(item.data)
                  .eq('id', item.data.id);
                break;
              case 'delete':
                result = await supabase
                  .from(item.table)
                  .delete()
                  .eq('id', item.data.id);
                break;
            }

            if (result.error) {
              if (result.error.code === '409') {
                failedItems.push(item);
              } else if (item.retryCount < 3) {
                item.retryCount++;
                failedItems.push(item);
              }
            }
          } catch (error) {
            console.error('Sync error:', error);
            failedItems.push(item);
          }
        }

        set({
          pendingChanges: failedItems,
          lastSyncTime: Date.now(),
        });
      },

      resolveConflict: async (itemId, resolution) => {
        const item = get().pendingChanges.find((i) => i.id === itemId);
        if (!item) return;

        if (resolution === 'server-wins') {
          set((state) => ({
            pendingChanges: state.pendingChanges.filter((i) => i.id !== itemId),
          }));
        } else {
          await get().processQueue();
        }
      },
    }),
    {
      name: 'biometrics-sync',
      storage: createJSONStorage(() => AsyncStorage),
    }
  )
);

export const useSyncStore = useSyncStore;
```

---

## 7) Navigation Structure

### Root Layout with Auth Guard

```typescript
// src/app/_layout.tsx
import { useEffect } from 'react';
import { Stack } from 'expo-router';
import { useAuthStore } from '@/stores/authStore';
import { useSyncStore } from '@/lib/sync';
import { Loading } from '@/components/ui/Loading';
import { NetworkProvider } from '@/components/network/NetworkProvider';

export default function RootLayout() {
  const { isLoading, isAuthenticated } = useAuthStore();
  const { setOnline } = useSyncStore();

  useEffect(() => {
    const unsubscribe = NetInfo.addEventListener((state) => {
      setOnline(state.isConnected ?? false);
    });

    return () => unsubscribe();
  }, []);

  if (isLoading) {
    return <Loading fullScreen message="Initializing..." />;
  }

  return (
    <NetworkProvider>
      <Stack
        screenOptions={{
          headerShown: false,
          animation: 'fade',
        }}
      >
        {isAuthenticated ? (
          <Stack.Screen name="(main)" />
        ) : (
          <Stack.Screen name="(auth)" />
        )}
      </Stack>
    </NetworkProvider>
  );
}
```

---

## 8) Push Notifications

### NotificationService

```typescript
// src/services/notifications.ts
import * as Notifications from 'expo-notifications';
import Constants from 'expo-constants';
import { supabase } from '@/lib/supabase';
import { useAuthStore } from '@/stores/authStore';

Notifications.setNotificationHandler({
  handleNotification: async () => ({
    shouldShowAlert: true,
    shouldPlaySound: true,
    shouldSetBadge: true,
    shouldShowBanner: true,
    shouldShowList: true,
  }),
});

class NotificationService {
  private expoToken: string | null = null;

  async initialize(): Promise<void> {
    const { status: existingStatus } =
      await Notifications.getPermissionsAsync();
    
    let finalStatus = existingStatus;
    
    if (existingStatus !== 'granted') {
      const { status } = await Notifications.requestPermissionsAsync();
      finalStatus = status;
    }

    if (finalStatus !== 'granted') {
      console.log('Notification permissions not granted');
      return;
    }

    this.expoToken = (
      await Notifications.getExpoPushTokenAsync({
        projectId: Constants.expoConfig?.extra?.eas?.projectId,
      })
    ).data;

    await this.registerToken();
  }

  private async registerToken(): Promise<void> {
    const { userId } = useAuthStore.getState();
    
    if (!userId || !this.expoToken) return;

    await supabase.from('push_tokens').upsert({
      user_id: userId,
      token: this.expoToken,
      platform: 'mobile',
      created_at: new Date().toISOString(),
    });
  }

  async scheduleBiometricReminder(): Promise<void> {
    await Notifications.scheduleNotificationAsync({
      content: {
        title: 'Biometric Verification Required',
        body: 'Please verify your identity to continue using the app',
        data: { type: 'biometric_reminder' },
      },
      trigger: {
        type: Notifications.SchedulableTriggerInputTypes.TIME_INTERVAL,
        seconds: 3600,
        repeats: true,
      },
    });
  }

  async cancelAllNotifications(): Promise<void> {
    await Notifications.cancelAllScheduledNotificationsAsync();
  }
}

export const notificationService = new NotificationService();
```

---

## 9) State Management with Zustand

### Auth Store

```typescript
// src/stores/authStore.ts
import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { supabase } from '@/lib/supabase';
import { keychainWrapper } from '@/lib/keychain';

interface User {
  id: string;
  email: string;
  createdAt: string;
}

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  biometricEnabled: boolean;
  sessionExpiry: number | null;
}

interface AuthActions {
  signIn: (email: string, password: string) => Promise<boolean>;
  signUp: (email: string, password: string) => Promise<boolean>;
  signOut: () => Promise<void>;
  enableBiometric: () => Promise<boolean>;
  refreshSession: () => Promise<boolean>;
}

export const useAuthStore = create<AuthState & AuthActions>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isLoading: true,
      biometricEnabled: false,
      sessionExpiry: null,

      signIn: async (email, password) => {
        set({ isLoading: true });
        
        try {
          const { data, error } = await supabase.auth.signInWithPassword({
            email,
            password,
          });

          if (error) {
            console.error('Sign in error:', error);
            set({ isLoading: false });
            return false;
          }

          await keychainWrapper.storeCredentials({
            username: email,
            password,
            accessToken: data.session?.access_token,
            refreshToken: data.session?.refresh_token,
          });

          set({
            user: data.user,
            isAuthenticated: true,
            isLoading: false,
            sessionExpiry: Date.now() + 3600000,
          });

          return true;
        } catch (error) {
          console.error('Sign in error:', error);
          set({ isLoading: false });
          return false;
        }
      },

      signUp: async (email, password) => {
        set({ isLoading: true });

        try {
          const { data, error } = await supabase.auth.signUp({
            email,
            password,
          });

          if (error) {
            console.error('Sign up error:', error);
            set({ isLoading: false });
            return false;
          }

          set({ isLoading: false });
          return true;
        } catch (error) {
          console.error('Sign up error:', error);
          set({ isLoading: false });
          return false;
        }
      },

      signOut: async () => {
        await supabase.auth.signOut();
        await keychainWrapper.deleteCredentials();
        
        set({
          user: null,
          isAuthenticated: false,
          biometricEnabled: false,
          sessionExpiry: null,
        });
      },

      enableBiometric: async () => {
        const { user } = get();
        if (!user) return false;

        const isAvailable = await keychainWrapper.isBiometricAvailable();
        if (!isAvailable) return false;

        const result = await biometricService.enrollBiometric(user.id);
        
        if (result.success) {
          set({ biometricEnabled: true });
        }

        return result.success;
      },

      refreshSession: async () => {
        const { data, error } = await supabase.auth.refreshSession();
        
        if (error) {
          console.error('Session refresh error:', error);
          return false;
        }

        set({
          sessionExpiry: Date.now() + 3600000,
        });

        return true;
      },
    }),
    {
      name: 'biometrics-auth',
      storage: createJSONStorage(() => AsyncStorage),
      partialize: (state) => ({
        biometricEnabled: state.biometricEnabled,
      }),
    }
  )
);
```

---

## 10) Testing Strategy

### Unit Tests Example

```typescript
// __tests__/lib/biometrics.test.ts
import { biometricService } from '@/lib/biometrics';

jest.mock('expo-local-authentication', () => ({
  hasHardwareAsync: jest.fn().mockResolvedValue(true),
  isEnrolledAsync: jest.fn().mockResolvedValue(true),
  supportedAuthenticationTypesAsync: jest.fn().mockResolvedValue([1]),
  authenticateAsync: jest.fn().mockResolvedValue({
    success: true,
    authenticationType: 1,
  }),
}));

jest.mock('react-native-keychain', () => ({
  setGenericPassword: jest.fn().mockResolvedValue(true),
  getGenericPassword: jest.fn().mockResolvedValue(false),
  resetGenericPassword: jest.fn().mockResolvedValue(true),
  ACCESS_CONTROL: {
    BIOMETRY_ANY: 'biometryAny',
  },
  ACCESSIBLE: {
    WHEN_PASSCODE_SET_THIS_DEVICE_ONLY: 'whenPasscodeSetThisDeviceOnly',
  },
}));

describe('BiometricService', () => {
  it('should check biometric availability', async () => {
    const result = await biometricService.checkAvailability();
    
    expect(result.available).toBe(true);
    expect(result.biometricTypes).toContain('fingerprint');
  });

  it('should authenticate successfully', async () => {
    const result = await biometricService.authenticate('Test authentication');
    
    expect(result.success).toBe(true);
  });
});
```

---

## 11) Build Configuration

### app.json

```json
{
  "expo": {
    "name": "BIOMETRICS",
    "slug": "biometrics-app",
    "version": "1.0.0",
    "orientation": "portrait",
    "icon": "./assets/icon.png",
    "userInterfaceStyle": "automatic",
    "newArchEnabled": true,
    "splash": {
      "image": "./assets/splash.png",
      "resizeMode": "contain",
      "backgroundColor": "#00FF00"
    },
    "ios": {
      "supportsTablet": true,
      "bundleIdentifier": "com.biometrics.app",
      "infoPlist": {
        "NSFaceIDUsageDescription": "We use Face ID for secure authentication"
      }
    },
    "android": {
      "adaptiveIcon": {
        "foregroundImage": "./assets/adaptive-icon.png",
        "backgroundColor": "#00FF00"
      },
      "package": "com.biometrics.app",
      "permissions": [
        "USE_BIOMETRIC",
        "USE_FINGERPRINT"
      ]
    },
    "plugins": [
      [
        "expo-notifications",
        {
          "icon": "./assets/notification-icon.png",
          "color": "#00FF00"
        }
      ]
    ],
    "extra": {
      "eas": {
        "projectId": "biometrics-project-id"
      }
    }
  }
}
```

---

## 12) Deployment

### EAS Build Commands

```bash
# Install EAS CLI
npm install -g eas-cli

# Login to EAS
eas login

# Configure project
eas build:configure

# Build for iOS (Development)
eas build -p ios --profile development

# Build for iOS (Production)
eas build -p ios --profile production

# Build for Android
eas build -p android --profile production

# Submit to App Store
eas submit -p ios

# Submit to Play Store
eas submit -p android
```

---

## 13) Security Considerations

1. **Biometric Data**: Never store raw biometric data; use platform keychain
2. **Token Storage**: Access tokens stored in Keychain with biometric protection
3. **Certificate Pinning**: Implement for all API calls
4. **Obfuscation**: Enable ProGuard/R8 for Android, Code Obfuscation for iOS
5. **Jailbreak Detection**: Implement runtime integrity checks
6. **Secure Enclave**: Utilize iOS Secure Enclave for key operations

---

## 14) Performance Optimization

- **MMKV**: Ultra-fast key-value storage for local data
- **Hermes**: JavaScript engine for better performance
- **Lazy Loading**: Code splitting for routes
- **Image Optimization**: WebP format with progressive loading
- **Memoization**: React.memo and useMemo for expensive computations
- **List Virtualization**: FlashList for long lists

---

## 15) Monitoring & Analytics

- **Crashlytics**: Firebase Crashlytics for crash reporting
- **Sentry**: Error tracking and performance monitoring
- **Mixpanel**: User analytics and event tracking
- **Amplitude**: Product analytics

---

## 16) Compliance

- **GDPR**: Data minimization, right to erasure
- **CCPA**: California Consumer Privacy Act compliance
- **SOC 2**: Security controls documentation
- **ISO 27001**: Information security management

---

Status: APPROVED  
Version: 1.0  
Last Updated: Februar 2026
