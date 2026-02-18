# DESKTOP-APP.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Architekturentscheidungen folgen den globalen Regeln aus `AGENTS-GLOBAL.md`.
- Jede Strukturänderung benötigt Mapping- und Integrationsabgleich.
- Security-by-Default, NLM-First und Nachweisbarkeit sind Pflicht.

Status: ACTIVE  
Version: 1.0  
Stand: Februar 2026

---

## 1) Overview

This document describes the Electron desktop application architecture for BIOMETRICS. The desktop app provides advanced biometric management, system tray integration, and native OS features for enterprise environments.

### Core Features
- Biometric device management and monitoring
- System tray with quick actions
- Native file system integration
- Desktop notifications
- Auto-update functionality
- Offline mode with local database
- Multi-window management

---

## 2) Technology Stack

### Framework & Tools
```json
{
  "electron": "34.0.0",
  "electron-builder": "25.0.0",
  "typescript": "5.3.3",
  "react": "18.3.1",
  "vite": "6.0.0",
  "electron-vite": "3.0.0"
}
```

### Key Dependencies
```json
{
  "dependencies": {
    "electron-store": "^10.0.0",
    "electron-updater": "^6.3.0",
    "electron-log": "^5.2.0",
    "better-sqlite3": "^11.0.0",
    "node-notifier": "^10.0.1",
    "usb-detection": "^4.16.0",
    "fido2-lib": "^4.0.0",
    "openpgp": "^6.0.0",
    "ws": "^8.18.0"
  }
}
```

---

## 3) Project Structure

```
biometrics-desktop/
├── electron/
│   ├── main/                   # Main process
│   │   ├── index.ts            # Entry point
│   │   ├── window.ts           # Window management
│   │   ├── tray.ts            # System tray
│   │   ├── ipc.ts             # IPC handlers
│   │   ├── updater.ts         # Auto-update
│   │   ├── usb.ts             # USB device detection
│   │   └── database.ts        # SQLite operations
│   ├── preload/
│   │   └── index.ts           # Preload scripts
│   └── protocol.ts            # Custom protocol
├── src/
│   ├── renderer/               # React frontend
│   │   ├── App.tsx
│   │   ├── pages/
│   │   ├── components/
│   │   ├── hooks/
│   │   └── stores/
│   ├── shared/                # Shared types
│   └── main.ts                # Renderer entry
├── resources/                 # App resources
│   ├── icons/
│   └── entitlements/
├── build/                     # Build configs
├── dist/                      # Output
├── electron.vite.config.ts
├── tsconfig.json
├── package.json
└── electron-builder.yml
```

---

## 4) Main Process Implementation

### Main Entry Point

```typescript
// electron/main/index.ts
import { app, BrowserWindow, ipcMain, Menu, Tray, nativeImage } from 'electron';
import { join } from 'path';
import { electronApp, optimizer, is } from '@electron-toolkit/utils';
import { initDatabase } from './database';
import { setupIpcHandlers } from './ipc';
import { createTray, updateTrayMenu } from './tray';
import { setupAutoUpdater } from './updater';
import { initUsbDetection } from './usb';
import log from 'electron-log';

let mainWindow: BrowserWindow | null = null;
let tray: Tray | null = null;

async function createWindow(): Promise<void> {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    minWidth: 900,
    minHeight: 600,
    show: false,
    title: 'BIOMETRICS Desktop',
    backgroundColor: '#000000',
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false,
      contextIsolation: true,
      nodeIntegration: false,
    },
  });

  mainWindow.on('ready-to-show', () => {
    mainWindow?.show();
    log.info('Main window shown');
  });

  mainWindow.on('close', (event) => {
    if (!app.isQuitting) {
      event.preventDefault();
      mainWindow?.hide();
    }
  });

  mainWindow.webContents.setWindowOpenHandler((details) => {
    require('electron').shell.openExternal(details.url);
    return { action: 'deny' };
  });

  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    await mainWindow.loadURL(process.env['ELECTRON_RENDERER_URL']);
  } else {
    await mainWindow.loadFile(join(__dirname, '../renderer/index.html'));
  }
}

function createMenu(): void {
  const template: Electron.MenuItemConstructorOptions[] = [
    {
      label: 'File',
      submenu: [
        {
          label: 'New Session',
          accelerator: 'CmdOrCtrl+N',
          click: () => mainWindow?.webContents.send('menu:new-session'),
        },
        { type: 'separator' },
        {
          label: 'Settings',
          accelerator: 'CmdOrCtrl+,',
          click: () => mainWindow?.webContents.send('menu:settings'),
        },
        { type: 'separator' },
        { role: 'quit' },
      ],
    },
    {
      label: 'Edit',
      submenu: [
        { role: 'undo' },
        { role: 'redo' },
        { type: 'separator' },
        { role: 'cut' },
        { role: 'copy' },
        { role: 'paste' },
        { role: 'selectAll' },
      ],
    },
    {
      label: 'View',
      submenu: [
        { role: 'reload' },
        { role: 'forceReload' },
        { role: 'toggleDevTools' },
        { type: 'separator' },
        { role: 'resetZoom' },
        { role: 'zoomIn' },
        { role: 'zoomOut' },
        { type: 'separator' },
        { role: 'togglefullscreen' },
      ],
    },
    {
      label: 'Window',
      submenu: [
        { role: 'minimize' },
        { role: 'zoom' },
        { type: 'separator' },
        { role: 'close' },
      ],
    },
    {
      label: 'Help',
      submenu: [
        {
          label: 'Documentation',
          click: () => require('electron').shell.openExternal('https://docs.biometrics.app'),
        },
        {
          label: 'About',
          click: () => mainWindow?.webContents.send('menu:about'),
        },
      ],
    },
  ];

  const menu = Menu.buildFromTemplate(template);
  Menu.setApplicationMenu(menu);
}

app.whenReady().then(async () => {
  log.initialize();
  log.info('Application starting...');

  electronApp.setAppUserModelId('com.biometrics.app');

  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window);
  });

  await initDatabase();
  log.info('Database initialized');

  setupIpcHandlers();
  log.info('IPC handlers registered');

  await createWindow();
  createMenu();

  tray = createTray(mainWindow!);
  log.info('System tray created');

  setupAutoUpdater(mainWindow!);
  initUsbDetection(mainWindow!);
  log.info('USB detection initialized');

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    } else {
      mainWindow?.show();
    }
  });
});

app.on('before-quit', () => {
  (app as any).isQuitting = true;
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});
```

---

## 5) System Tray Implementation

### Tray Manager

```typescript
// electron/main/tray.ts
import { Tray, Menu, nativeImage, app, BrowserWindow } from 'electron';
import { join } from 'path';
import log from 'electron-log';

let tray: Tray | null = null;

export function createTray(mainWindow: BrowserWindow): Tray {
  const iconPath = join(__dirname, '../../resources/icon.png');
  const icon = nativeImage.createFromPath(iconPath).resize({ width: 16, height: 16 });

  tray = new Tray(icon);
  tray.setToolTip('BIOMETRICS Desktop');

  updateTrayMenu(mainWindow);

  tray.on('click', () => {
    if (mainWindow.isVisible()) {
      mainWindow.hide();
    } else {
      mainWindow.show();
    }
  });

  tray.on('double-click', () => {
    mainWindow.show();
  });

  return tray;
}

export function updateTrayMenu(mainWindow: BrowserWindow): void {
  const contextMenu = Menu.buildFromTemplate([
    {
      label: 'Open BIOMETRICS',
      click: () => mainWindow.show(),
    },
    { type: 'separator' },
    {
      label: 'Quick Actions',
      submenu: [
        {
          label: 'New Biometric Scan',
          click: () => {
            mainWindow.show();
            mainWindow.webContents.send('action:new-scan');
          },
        },
        {
          label: 'View Logs',
          click: () => {
            mainWindow.show();
            mainWindow.webContents.send('action:view-logs');
          },
        },
      ],
    },
    { type: 'separator' },
    {
      label: 'Devices',
      submenu: [
        {
          label: 'Connected Devices',
          click: () => {
            mainWindow.show();
            mainWindow.webContents.send('action:devices');
          },
        },
        {
          label: 'Refresh Devices',
          click: () => {
            mainWindow.webContents.send('action:refresh-devices');
          },
        },
      ],
    },
    { type: 'separator' },
    {
      label: 'Settings',
      click: () => {
        mainWindow.show();
        mainWindow.webContents.send('menu:settings');
      },
    },
    { type: 'separator' },
    {
      label: 'Quit',
      click: () => {
        (app as any).isQuitting = true;
        app.quit();
      },
    },
  ]);

  tray?.setContextMenu(contextMenu);
}

export function showTrayNotification(title: string, body: string): void {
  if (tray) {
    tray.displayBalloon({
      title,
      content: body,
      iconType: 'info',
    });
  }
}
```

---

## 6) IPC Communication

### IPC Handlers

```typescript
// electron/main/ipc.ts
import { ipcMain, dialog, shell, app } from 'electron';
import { initDatabase, db } from './database';
import log from 'electron-log';
import * as fs from 'fs';
import * as path from 'path';

export function setupIpcHandlers(): void {
  // Database operations
  ipcMain.handle('db:query', async (_, sql: string, params?: any[]) => {
    try {
      return db.prepare(sql).all(...(params || []));
    } catch (error) {
      log.error('Database query error:', error);
      throw error;
    }
  });

  ipcMain.handle('db:run', async (_, sql: string, params?: any[]) => {
    try {
      return db.prepare(sql).run(...(params || []));
    } catch (error) {
      log.error('Database run error:', error);
      throw error;
    }
  });

  // File dialogs
  ipcMain.handle('dialog:openFile', async (_, options: Electron.OpenDialogOptions) => {
    const result = await dialog.showOpenDialog(options);
    return result;
  });

  ipcMain.handle('dialog:saveFile', async (_, options: Electron.SaveDialogOptions) => {
    const result = await dialog.showSaveDialog(options);
    return result;
  });

  // Shell operations
  ipcMain.handle('shell:openExternal', async (_, url: string) => {
    await shell.openExternal(url);
  });

  ipcMain.handle('shell:openPath', async (_, filePath: string) => {
    await shell.openPath(filePath);
  });

  // App info
  ipcMain.handle('app:getVersion', () => app.getVersion());
  ipcMain.handle('app:getPath', (_, name: string) => app.getPath(name as any));

  // File system
  ipcMain.handle('fs:readFile', async (_, filePath: string) => {
    return fs.readFileSync(filePath);
  });

  ipcMain.handle('fs:writeFile', async (_, filePath: string, data: Buffer) => {
    fs.writeFileSync(filePath, data);
  });

  ipcMain.handle('fs:exists', async (_, filePath: string) => {
    return fs.existsSync(filePath);
  });

  // Window controls
  ipcMain.on('window:minimize', (event) => {
    const win = BrowserWindow.fromWebContents(event.sender);
    win?.minimize();
  });

  ipcMain.on('window:maximize', (event) => {
    const win = BrowserWindow.fromWebContents(event.sender);
    if (win?.isMaximized()) {
      win.unmaximize();
    } else {
      win?.maximize();
    }
  });

  ipcMain.on('window:close', (event) => {
    const win = BrowserWindow.fromWebContents(event.sender);
    win?.close();
  });

  log.info('IPC handlers setup complete');
}
```

---

## 7) USB Device Detection

### Biometric Device Manager

```typescript
// electron/main/usb.ts
import { BrowserWindow } from 'electron';
import USB from 'usb-detection';
import log from 'electron-log';

interface BiometricDevice {
  id: string;
  vendorId: number;
  productId: number;
  productName: string;
  manufacturer: string;
  connected: boolean;
}

let detectedDevices: Map<string, BiometricDevice> = new Map();

export function initUsbDetection(mainWindow: BrowserWindow): void {
  USB.startMonitoring();

  USB.on('add', (device) => {
    log.info('USB device added:', device);
    
    const biometricDevice: BiometricDevice = {
      id: `${device.vendorId}:${device.productId}`,
      vendorId: device.vendorId,
      productId: device.productId,
      productName: device.productName || 'Unknown',
      manufacturer: device.manufacturer || 'Unknown',
      connected: true,
    };

    detectedDevices.set(biometricDevice.id, biometricDevice);
    
    mainWindow.webContents.send('usb:device-added', biometricDevice);
  });

  USB.on('remove', (device) => {
    log.info('USB device removed:', device);
    
    const deviceId = `${device.vendorId}:${device.productId}`;
    const existingDevice = detectedDevices.get(deviceId);
    
    if (existingDevice) {
      existingDevice.connected = false;
      detectedDevices.delete(deviceId);
      mainWindow.webContents.send('usb:device-removed', existingDevice);
    }
  });

  USB.on('change', (device) => {
    log.info('USB device changed:', device);
    mainWindow.webContents.send('usb:device-changed', device);
  });

  // Initial scan
  USB.find()
    .then((devices) => {
      log.info('Initial USB devices:', devices);
      devices.forEach((device) => {
        const biometricDevice: BiometricDevice = {
          id: `${device.vendorId}:${device.productId}`,
          vendorId: device.vendorId,
          productId: device.productId,
          productName: device.productName || 'Unknown',
          manufacturer: device.manufacturer || 'Unknown',
          connected: true,
        };
        detectedDevices.set(biometricDevice.id, biometricDevice);
      });
      mainWindow.webContents.send('usb:devices-list', Array.from(detectedDevices.values()));
    })
    .catch((err) => log.error('USB detection error:', err));

  log.info('USB detection initialized');
}

export function getDetectedDevices(): BiometricDevice[] {
  return Array.from(detectedDevices.values());
}

export function stopUsbDetection(): void {
  USB.stopMonitoring();
}
```

---

## 8) Auto-Updater

### Update Manager

```typescript
// electron/main/updater.ts
import { autoUpdater } from 'electron-updater';
import { BrowserWindow, dialog } from 'electron';
import log from 'electron-log';

let mainWindow: BrowserWindow | null = null;

export function setupAutoUpdater(window: BrowserWindow): void {
  mainWindow = window;

  autoUpdater.logger = log;
  autoUpdater.autoDownload = true;
  autoUpdater.autoInstallOnAppQuit = true;

  autoUpdater.on('checking-for-update', () => {
    log.info('Checking for update...');
    mainWindow?.webContents.send('updater:checking');
  });

  autoUpdater.on('update-available', (info) => {
    log.info('Update available:', info.version);
    mainWindow?.webContents.send('updater:available', info);
  });

  autoUpdater.on('update-not-available', (info) => {
    log.info('Update not available');
    mainWindow?.webContents.send('updater:not-available', info);
  });

  autoUpdater.on('download-progress', (progress) => {
    log.info(`Download progress: ${progress.percent}%`);
    mainWindow?.webContents.send('updater:progress', progress);
  });

  autoUpdater.on('update-downloaded', (info) => {
    log.info('Update downloaded:', info.version);
    mainWindow?.webContents.send('updater:downloaded', info);
    
    dialog.showMessageBox({
      type: 'info',
      title: 'Update Ready',
      message: `Version ${info.version} has been downloaded. The application will restart to install the update.`,
      buttons: ['Restart Now', 'Later'],
    }).then((result) => {
      if (result.response === 0) {
        autoUpdater.quitAndInstall();
      }
    });
  });

  autoUpdater.on('error', (error) => {
    log.error('Updater error:', error);
    mainWindow?.webContents.send('updater:error', error.message);
  });

  // Check for updates on startup (with delay)
  setTimeout(() => {
    autoUpdater.checkForUpdates().catch((err) => {
      log.error('Failed to check for updates:', err);
    });
  }, 5000);
}
```

---

## 9) Database Implementation

### SQLite Setup

```typescript
// electron/main/database.ts
import Database from 'better-sqlite3';
import { app } from 'electron';
import { join } from 'path';
import { existsSync, mkdirSync } from 'fs';
import log from 'electron-log';

export let db: Database.Database;

export async function initDatabase(): Promise<void> {
  const userDataPath = app.getPath('userData');
  const dbPath = join(userDataPath, 'biometrics.db');

  if (!existsSync(userDataPath)) {
    mkdirSync(userDataPath, { recursive: true });
  }

  db = new Database(dbPath);
  db.pragma('journal_mode = WAL');

  // Create tables
  db.exec(`
    CREATE TABLE IF NOT EXISTS sessions (
      id TEXT PRIMARY KEY,
      user_id TEXT NOT NULL,
      started_at INTEGER NOT NULL,
      ended_at INTEGER,
      device_info TEXT,
      ip_address TEXT,
      created_at INTEGER DEFAULT (strftime('%s', 'now'))
    );

    CREATE TABLE IF NOT EXISTS biometric_records (
      id TEXT PRIMARY KEY,
      user_id TEXT NOT NULL,
      device_id TEXT,
      biometric_type TEXT NOT NULL,
      template_data BLOB,
      quality_score REAL,
      enrolled_at INTEGER NOT NULL,
      last_used INTEGER,
      is_active INTEGER DEFAULT 1
    );

    CREATE TABLE IF NOT EXISTS audit_logs (
      id TEXT PRIMARY KEY,
      user_id TEXT,
      action TEXT NOT NULL,
      details TEXT,
      ip_address TEXT,
      user_agent TEXT,
      timestamp INTEGER NOT NULL
    );

    CREATE TABLE IF NOT EXISTS settings (
      key TEXT PRIMARY KEY,
      value TEXT,
      updated_at INTEGER
    );

    CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
    CREATE INDEX IF NOT EXISTS idx_biometric_records_user_id ON biometric_records(user_id);
    CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp);
  `);

  log.info('Database initialized at:', dbPath);
}

export function closeDatabase(): void {
  if (db) {
    db.close();
    log.info('Database closed');
  }
}
```

---

## 10) Preload Script

### Context Bridge

```typescript
// electron/preload/index.ts
import { contextBridge, ipcRenderer } from 'electron';

const api = {
  // Database
  db: {
    query: (sql: string, params?: any[]) => ipcRenderer.invoke('db:query', sql, params),
    run: (sql: string, params?: any[]) => ipcRenderer.invoke('db:run', sql, params),
  },

  // Dialogs
  dialog: {
    openFile: (options: any) => ipcRenderer.invoke('dialog:openFile', options),
    saveFile: (options: any) => ipcRenderer.invoke('dialog:saveFile', options),
  },

  // Shell
  shell: {
    openExternal: (url: string) => ipcRenderer.invoke('shell:openExternal', url),
    openPath: (path: string) => ipcRenderer.invoke('shell:openPath', path),
  },

  // App
  app: {
    getVersion: () => ipcRenderer.invoke('app:getVersion'),
    getPath: (name: string) => ipcRenderer.invoke('app:getPath', name),
  },

  // File System
  fs: {
    readFile: (path: string) => ipcRenderer.invoke('fs:readFile', path),
    writeFile: (path: string, data: Buffer) => ipcRenderer.invoke('fs:writeFile', path, data),
    exists: (path: string) => ipcRenderer.invoke('fs:exists', path),
  },

  // Window
  window: {
    minimize: () => ipcRenderer.send('window:minimize'),
    maximize: () => ipcRenderer.send('window:maximize'),
    close: () => ipcRenderer.send('window:close'),
  },

  // Events
  on: (channel: string, callback: (...args: any[]) => void) => {
    const subscription = (_event: any, ...args: any[]) => callback(...args);
    ipcRenderer.on(channel, subscription);
    return () => ipcRenderer.removeListener(channel, subscription);
  },

  // USB
  usb: {
    onDeviceAdded: (callback: (device: any) => void) => 
      ipcRenderer.on('usb:device-added', (_event, device) => callback(device)),
    onDeviceRemoved: (callback: (device: any) => void) =>
      ipcRenderer.on('usb:device-removed', (_event, device) => callback(device)),
    onDevicesList: (callback: (devices: any[]) => void) =>
      ipcRenderer.on('usb:devices-list', (_event, devices) => callback(devices)),
  },

  // Updater
  updater: {
    onChecking: (callback: () => void) => ipcRenderer.on('updater:checking', callback),
    onAvailable: (callback: (info: any) => void) => ipcRenderer.on('updater:available', (_event, info) => callback(info)),
    onProgress: (callback: (progress: any) => void) => ipcRenderer.on('updater:progress', (_event, progress) => callback(progress)),
    onDownloaded: (callback: (info: any) => void) => ipcRenderer.on('updater:downloaded', (_event, info) => callback(info)),
  },
};

contextBridge.exposeInMainWorld('api', api);

export type ApiType = typeof api;
```

---

## 11) Build Configuration

### electron-builder.yml

```yaml
appId: com.biometrics.app
productName: BIOMETRICS Desktop
copyright: Copyright © 2026 BIOMETRICS

directories:
  output: dist
  buildResources: resources

files:
  - "!**/.vscode/*"
  - "!src/*"
  - "!electron.vite.config.*"
  - "!{.eslintignore,.eslintrc.cjs,.prettierignore,.prettierrc.yaml,dev-app-update.yml,CHANGELOG.md,README.md}"
  - "!{tsconfig.json,tsconfig.node.json,tsconfig.web.json}"

asarUnpack:
  - resources/**
  - "**/*.node"

win:
  executableName: biometrics-desktop
  target:
    - target: nsis
      arch:
        - x64

nsis:
  oneClick: false
  perMachine: true
  allowToChangeInstallationDirectory: true
  deleteAppDataOnUninstall: true
  installerIcon: resources/icon.ico
  uninstallerIcon: resources/icon.ico
  installerHeaderIcon: resources/icon.ico
  createDesktopShortcut: true
  createStartMenuShortcut: true

mac:
  entitlementsInherit: build/entitlements.mac.plist
  extendInfo:
    - NSCameraUsageDescription: Required for biometric scanning
    - NSFaceIDUsageDescription: Required for Face ID authentication
  target:
    - target: dmg
      arch:
        - x64
        - arm64

dmg:
  contents:
    - x: 130
      y: 220
    - x: 410
      y: 220
      type: link
      path: /Applications

linux:
  target:
    - target: AppImage
      arch:
        - x64
  maintainer: biometrics.app
  category: Utility

publish:
  provider: generic
  url: https://updates.biometrics.app
```

---

## 12) React Frontend Integration

### Main App Component

```typescript
// src/renderer/App.tsx
import { useEffect, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { MainLayout } from './layouts/MainLayout';
import { AuthPage } from './pages/AuthPage';
import { DashboardPage } from './pages/DashboardPage';
import { useAuthStore } from './stores/authStore';
import './styles/global.css';

function App() {
  const { isAuthenticated, checkAuth } = useAuthStore();
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      await checkAuth();
      setIsLoading(false);
    };
    init();
  }, []);

  if (isLoading) {
    return <div className="loading-screen">Loading...</div>;
  }

  return isAuthenticated ? <MainLayout /> : <AuthPage />;
}

const container = document.getElementById('root');
if (container) {
  const root = createRoot(container);
  root.render(<App />);
}
```

---

## 13) Window Management

### Window State Manager

```typescript
// electron/main/window.ts
import { BrowserWindow, screen } from 'electron';
import Store from 'electron-store';

interface WindowState {
  x: number;
  y: number;
  width: number;
  height: number;
  isMaximized: boolean;
  isFullScreen: boolean;
}

const store = new Store<{ windowState: WindowState }>({
  defaults: {
    windowState: {
      width: 1200,
      height: 800,
      x: undefined,
      y: undefined,
      isMaximized: false,
      isFullScreen: false,
    },
  },
});

export function getWindowState(): WindowState {
  return store.get('windowState');
}

export function saveWindowState(win: BrowserWindow): void {
  const bounds = win.getBounds();
  const isMaximized = win.isMaximized();
  const isFullScreen = win.isFullScreen();

  store.set('windowState', {
    x: bounds.x,
    y: bounds.y,
    width: bounds.width,
    height: bounds.height,
    isMaximized,
    isFullScreen,
  });
}

export function ensureVisibleOnScreen(win: BrowserWindow): void {
  const state = getWindowState();
  const displays = screen.getAllDisplays();

  let isVisible = false;
  for (const display of displays) {
    const { x, y, width, height } = display.bounds;
    if (
      state.x !== undefined &&
      state.y !== undefined &&
      state.x >= x &&
      state.x < x + width &&
      state.y >= y &&
      state.y < y + height
    ) {
      isVisible = true;
      break;
    }
  }

  if (!isVisible) {
    const primaryDisplay = screen.getPrimaryDisplay();
    const { width, height } = primaryDisplay.workAreaSize;
    win.setBounds({
      x: Math.round((width - state.width) / 2),
      y: Math.round((height - state.height) / 2),
      width: state.width,
      height: state.height,
    });
  }
}
```

---

## 14) Logging & Monitoring

### Log Configuration

```typescript
// electron/main/logging.ts
import log from 'electron-log';
import { app } from 'electron';
import { join } from 'path';

export function initLogging(): void {
  log.transports.file.level = 'info';
  log.transports.console.level = 'debug';
  log.transports.file.maxSize = 10 * 1024 * 1024; // 10MB
  log.transports.file.format = '[{y}-{m}-{d} {h}:{i}:{s}.{ms}] [{level}] {text}';
  log.transports.file.resolvePathFn = () => join(app.getPath('logs'), 'biometrics.log');

  log.catchErrors({
    showDialog: false,
    onError: (error) => {
      log.error('Uncaught exception:', error);
    },
  });

  process.on('unhandledRejection', (reason) => {
    log.error('Unhandled promise rejection:', reason);
  });

  log.info('Logging initialized');
}
```

---

## 15) Security Features

1. **Context Isolation**: Enabled by default
2. **Node Integration**: Disabled in renderer
3. **Sandbox**: Enabled for renderer process
4. **CSP**: Content Security Policy configured
5. **Code Signing**: Required for production builds
6. **Hardened Runtime**: Enabled for macOS
7. **Secure Enclave**: Integration for key operations

---

## 16) Performance Optimization

- **Lazy Loading**: Dynamic imports for routes
- **Code Splitting**: By route and component
- **Caching**: HTTP and file system caching
- **Web Workers**: Heavy computations in workers
- **Virtual Scrolling**: For large lists
- **Memoization**: React performance optimizations

---

Status: APPROVED  
Version: 1.0  
Last Updated: Februar 2026
