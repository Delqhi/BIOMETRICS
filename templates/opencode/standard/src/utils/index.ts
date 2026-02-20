/**
 * Utility functions for OpenCode project
 */

export function delay(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export function generateId(): string {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

export function parseEnv<T extends Record<string, string>>(
  required: (keyof T)[],
  optional: (keyof T)[] = []
): T {
  const result = {} as T;
  
  for (const key of required) {
    const value = process.env[key as string];
    if (!value) {
      throw new Error(`Missing required environment variable: ${key}`);
    }
    result[key] = value as T[keyof T];
  }
  
  for (const key of optional) {
    result[key] = (process.env[key as string] || '') as T[keyof T];
  }
  
  return result;
}

export function isValidUrl(url: string): boolean {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
}

export function truncate(str: string, maxLength: number): string {
  if (str.length <= maxLength) return str;
  return str.slice(0, maxLength - 3) + '...';
}
