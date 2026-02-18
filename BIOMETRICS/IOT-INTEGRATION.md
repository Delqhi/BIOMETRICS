# IOT-INTEGRATION.md - IoT Device Integration

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Advanced Tech  
**Author:** BIOMETRICS IoT Team

---

## 1. Overview

This document describes the IoT device integration for BIOMETRICS, enabling connection with wearable devices, medical sensors, and smart home health equipment.

## 2. Supported Devices

### 2.1 Device Categories

| Category | Devices | Protocol |
|----------|---------|----------|
| Wearables | Fitbit, Garmin, Apple Watch | Health API |
| Smart Scales | Withings, Eufy | BLE + Cloud |
| Blood Pressure | Omron, Withings | BLE |
| Glucose Monitors | Dexcom, Libre | BLE |
| Sleep | Sleep Number, Oura | Cloud API |
| Smart Home | Apple HomeKit, Google Home | Matter |

### 2.2 Device Protocol

```typescript
interface Device {
  id: string;
  type: DeviceType;
  name: string;
  manufacturer: string;
  model: string;
  firmware: string;
  status: DeviceStatus;
  lastSync: Date;
}

interface BiometricReading {
  deviceId: string;
  userId: string;
  type: MetricType;
  value: number;
  unit: string;
  timestamp: Date;
  metadata: Record<string, any>;
}
```

## 3. Bluetooth LE

### 3.1 Device Scanning

```javascript
// Scan for BLE devices
class BLEDeviceScanner {
  async scan(filters: DeviceFilter[]): Promise<Device[]> {
    const deviceFilter = filters.map(f => ({
      services: [f.serviceUUID],
      name: f.namePrefix,
    }));
    
    const device = await navigator.bluetooth.requestDevice({
      filters: deviceFilter,
      optionalServices: ['battery_service'],
    });
    
    return {
      id: device.id,
      name: device.name,
      gatt: await device.gatt.connect(),
    };
  }
  
  async discoverServices(device: BluetoothDevice) {
    const services = await device.gatt.getPrimaryServices();
    
    for (const service of services) {
      const characteristics = await service.getCharacteristics();
      
      for (const char of characteristics) {
        console.log(`Characteristic: ${char.uuid}`);
        
        if (char.properties.notify) {
          char.startNotifications();
          char.addEventListener('characteristicvaluechanged', 
            (event) => this.handleData(event, char)
          );
        }
      }
    }
  }
}
```

### 3.2 Characteristic Parsing

```javascript
// Parse device-specific data
class DeviceDataParser {
  parseHeartRate(data: DataView): number {
    const flags = data.getUint8(0);
    
    if (flags & 0x01) {
      // 16-bit heart rate
      return data.getUint16(1, true);
    } else {
      // 8-bit heart rate
      return data.getUint8(1);
    }
  }
  
  parseBatteryLevel(data: DataView): number {
    return data.getUint8(0);
  }
  
  parseSteps(data: DataView): number {
    return data.getUint32(0, true);
  }
  
  parseBloodPressure(data: DataView): { systolic: number; diastolic: number } {
    return {
      systolic: data.getUint16(0, true),
      diastolic: data.getUint16(2, true),
    };
  }
}
```

## 4. Cloud Integration

### 4.1 Fitbit Integration

```typescript
class FitbitClient {
  private accessToken: string;
  private baseUrl = 'https://api.fitbit.com';
  
  async getHeartRate(date: string, period: string = '1d') {
    const response = await fetch(
      `${this.baseUrl}/1/user/-/activities/heart/date/${date}/${period}.json`,
      { headers: { Authorization: `Bearer ${this.accessToken}` } }
    );
    
    return response.json();
  }
  
  async getSteps(date: string, period: string = '1d') {
    const response = await fetch(
      `${this.baseUrl}/1/user/-/activities/steps/date/${date}/${period}.json`,
      { headers: { Authorization: `Bearer ${this.accessToken}` } }
    );
    
    return response.json();
  }
  
  async getSleep(date: string) {
    const response = await fetch(
      `${this.baseUrl}/1.2/user/-/sleep/date/${date}.json`,
      { headers: { Authorization: `Bearer ${this.accessToken}` } }
    );
    
    return response.json();
  }
  
  async getActivitySummary(date: string) {
    const response = await fetch(
      `${this.baseUrl}/1/user/-/activities/date/${date}.json`,
      { headers: { Authorization: `Bearer ${this.accessToken}` } }
    );
    
    return response.json();
  }
}
```

### 4.2 Withings Integration

```typescript
class WithingsClient {
  private clientId: string;
  private clientSecret: string;
  private accessToken: string;
  
  async getMeasurements(userId: string) {
    const response = await fetch(
      `https://wbsapi.withings.net/v3/measure`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
          action: 'getmeasure',
          access_token: this.accessToken,
          userid: userId,
          lastupdate: Math.floor(Date.now() / 1000).toString(),
        }),
      }
    );
    
    return response.json();
  }
  
  parseMeasurements(data: any): BiometricReading[] {
    const readings: BiometricReading[] = [];
    
    for (const measure of data.body.measures) {
      switch (measure.type) {
        case 1: // Weight (kg)
          readings.push({
            type: 'weight',
            value: measure.value * Math.pow(10, measure.unit),
            unit: 'kg',
          });
          break;
        case 4: // Height (m)
          readings.push({
            type: 'height',
            value: measure.value * Math.pow(10, measure.unit),
            unit: 'm',
          });
          break;
        case 6: // Heart rate (bpm)
          readings.push({
            type: 'heart_rate',
            value: measure.value,
            unit: 'bpm',
          });
          break;
      }
    }
    
    return readings;
  }
}
```

## 5. Device Management

### 5.1 Device Registry

```sql
CREATE TABLE devices (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    device_type VARCHAR(50) NOT NULL,
    manufacturer VARCHAR(100),
    model VARCHAR(100),
    serial_number VARCHAR(255),
    firmware_version VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    last_sync TIMESTAMP,
    registered_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(serial_number)
);

CREATE TABLE device_readings (
    id UUID PRIMARY KEY,
    device_id UUID REFERENCES devices(id),
    user_id UUID REFERENCES users(id),
    reading_type VARCHAR(50) NOT NULL,
    value DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20),
    timestamp TIMESTAMP NOT NULL,
    metadata JSONB,
    
    INDEX idx_device_timestamp (device_id, timestamp)
);
```

### 5.2 Sync Service

```typescript
class DeviceSyncService {
  private devices: Map<string, DeviceConnection> = new Map();
  
  async registerDevice(userId: string, device: DeviceConfig) {
    const connection = await this.connectDevice(device);
    
    this.devices.set(device.id, {
      connection,
      userId,
      lastSync: new Date(),
    });
    
    // Start polling
    this.startPolling(device.id);
  }
  
  private async syncDevice(deviceId: string) {
    const device = this.devices.get(deviceId);
    
    try {
      // Fetch latest data
      const readings = await device.connection.fetchReadings();
      
      // Save to database
      await this.saveReadings(device.userId, readings);
      
      // Update last sync
      device.lastSync = new Date();
    } catch (error) {
      console.error(`Sync failed for device ${deviceId}:`, error);
      await this.handleSyncError(deviceId, error);
    }
  }
  
  private startPolling(deviceId: string) {
    const device = this.devices.get(deviceId);
    
    setInterval(() => {
      this.syncDevice(deviceId);
    }, device.pollInterval || 300000); // Default 5 minutes
  }
}
```

## 6. Data Normalization

### 6.1 Unified Data Model

```typescript
interface NormalizedReading {
  userId: string;
  metric: MetricType;
  value: number;
  unit: string;
  source: DataSource;
  timestamp: Date;
  quality: DataQuality;
}

class DataNormalizer {
  normalize(data: any, source: DataSource): NormalizedReading {
    switch (source) {
      case DataSource.FITBIT:
        return this.normalizeFitbit(data);
      case DataSource.WITHINGS:
        return this.normalizeWithings(data);
      case DataSource.APPLE_WATCH:
        return this.normalizeAppleHealth(data);
      case DataSource.OMRON:
        return this.normalizeOmron(data);
      default:
        throw new Error(`Unknown source: ${source}`);
    }
  }
  
  private normalizeFitbit(data: any): NormalizedReading[] {
    return data['activities-heart'].map(entry => ({
      metric: MetricType.HEART_RATE,
      value: parseInt(entry.value),
      unit: 'bpm',
      source: DataSource.FITBIT,
      timestamp: new Date(entry.dateTime),
      quality: DataQuality.HIGH,
    }));
  }
}
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
