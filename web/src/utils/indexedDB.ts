// IndexedDB wrapper for Smart Goal Calendar offline storage

interface DBConfig {
  name: string;
  version: number;
  stores: {
    [key: string]: {
      keyPath: string;
      autoIncrement?: boolean;
      indexes?: Array<{
        name: string;
        keyPath: string | string[];
        unique?: boolean;
      }>;
    };
  };
}

const DB_CONFIG: DBConfig = {
  name: 'SmartGoalCalendarDB',
  version: 1,
  stores: {
    goals: {
      keyPath: 'id',
      indexes: [
        { name: 'user_id', keyPath: 'user_id' },
        { name: 'status', keyPath: 'status' },
        { name: 'sync_status', keyPath: 'sync_status' }
      ]
    },
    events: {
      keyPath: 'id',
      indexes: [
        { name: 'user_id', keyPath: 'user_id' },
        { name: 'start_time', keyPath: 'start_time' },
        { name: 'sync_status', keyPath: 'sync_status' }
      ]
    },
    moods: {
      keyPath: 'id',
      indexes: [
        { name: 'user_id', keyPath: 'user_id' },
        { name: 'date', keyPath: 'date' },
        { name: 'sync_status', keyPath: 'sync_status' }
      ]
    },
    pendingActions: {
      keyPath: 'id',
      autoIncrement: true,
      indexes: [
        { name: 'type', keyPath: 'type' },
        { name: 'created_at', keyPath: 'created_at' }
      ]
    },
    syncMetadata: {
      keyPath: 'store',
      indexes: [
        { name: 'last_sync', keyPath: 'last_sync' }
      ]
    }
  }
};

export class IndexedDBManager {
  private db: IDBDatabase | null = null;
  private dbPromise: Promise<IDBDatabase> | null = null;

  constructor() {
    this.dbPromise = this.initDB();
  }

  private async initDB(): Promise<IDBDatabase> {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open(DB_CONFIG.name, DB_CONFIG.version);

      request.onerror = () => {
        console.error('Failed to open IndexedDB:', request.error);
        reject(request.error);
      };

      request.onsuccess = () => {
        this.db = request.result;
        console.log('IndexedDB opened successfully');
        resolve(request.result);
      };

      request.onupgradeneeded = (event) => {
        const db = (event.target as IDBOpenDBRequest).result;
        console.log('Upgrading IndexedDB schema...');

        // Create object stores
        Object.entries(DB_CONFIG.stores).forEach(([storeName, storeConfig]) => {
          if (!db.objectStoreNames.contains(storeName)) {
            const store = db.createObjectStore(storeName, {
              keyPath: storeConfig.keyPath,
              autoIncrement: storeConfig.autoIncrement
            });

            // Create indexes
            storeConfig.indexes?.forEach(index => {
              store.createIndex(index.name, index.keyPath, { unique: index.unique });
            });

            console.log(`Created object store: ${storeName}`);
          }
        });
      };
    });
  }

  private async getDB(): Promise<IDBDatabase> {
    if (this.db) return this.db;
    if (this.dbPromise) return this.dbPromise;
    this.dbPromise = this.initDB();
    return this.dbPromise;
  }

  // Generic CRUD operations
  async add<T>(storeName: string, data: T): Promise<void> {
    const db = await this.getDB();
    const transaction = db.transaction([storeName], 'readwrite');
    const store = transaction.objectStore(storeName);
    
    // Add sync status if not present
    const dataWithSync = {
      ...data,
      sync_status: 'pending',
      local_updated_at: new Date().toISOString()
    };

    return new Promise((resolve, reject) => {
      const request = store.add(dataWithSync);
      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  }

  async put<T>(storeName: string, data: T): Promise<void> {
    const db = await this.getDB();
    const transaction = db.transaction([storeName], 'readwrite');
    const store = transaction.objectStore(storeName);
    
    // Update sync status
    const dataWithSync = {
      ...data,
      sync_status: 'pending',
      local_updated_at: new Date().toISOString()
    };

    return new Promise((resolve, reject) => {
      const request = store.put(dataWithSync);
      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  }

  async get<T>(storeName: string, key: IDBValidKey): Promise<T | undefined> {
    const db = await this.getDB();
    const transaction = db.transaction([storeName], 'readonly');
    const store = transaction.objectStore(storeName);

    return new Promise((resolve, reject) => {
      const request = store.get(key);
      request.onsuccess = () => resolve(request.result);
      request.onerror = () => reject(request.error);
    });
  }

  async getAll<T>(storeName: string, indexName?: string, query?: IDBValidKey | IDBKeyRange): Promise<T[]> {
    const db = await this.getDB();
    const transaction = db.transaction([storeName], 'readonly');
    const store = transaction.objectStore(storeName);
    const source = indexName ? store.index(indexName) : store;

    return new Promise((resolve, reject) => {
      const request = query ? source.getAll(query) : source.getAll();
      request.onsuccess = () => resolve(request.result);
      request.onerror = () => reject(request.error);
    });
  }

  async delete(storeName: string, key: IDBValidKey): Promise<void> {
    const db = await this.getDB();
    const transaction = db.transaction([storeName], 'readwrite');
    const store = transaction.objectStore(storeName);

    return new Promise((resolve, reject) => {
      const request = store.delete(key);
      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  }

  async clear(storeName: string): Promise<void> {
    const db = await this.getDB();
    const transaction = db.transaction([storeName], 'readwrite');
    const store = transaction.objectStore(storeName);

    return new Promise((resolve, reject) => {
      const request = store.clear();
      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  }

  // Sync-specific methods
  async getPendingItems<T>(storeName: string): Promise<T[]> {
    return this.getAll<T>(storeName, 'sync_status', 'pending');
  }

  async markAsSynced(storeName: string, key: IDBValidKey): Promise<void> {
    const item = await this.get(storeName, key);
    if (item) {
      await this.put(storeName, {
        ...item,
        sync_status: 'synced',
        synced_at: new Date().toISOString()
      });
    }
  }

  async addPendingAction(action: {
    type: 'create' | 'update' | 'delete';
    store: string;
    data: Record<string, unknown>;
    url: string;
    method: string;
  }): Promise<void> {
    await this.add('pendingActions', {
      ...action,
      created_at: new Date().toISOString(),
      retry_count: 0
    });
  }

  async getPendingActions(): Promise<Array<Record<string, unknown>>> {
    return this.getAll('pendingActions');
  }

  async removePendingAction(id: number): Promise<void> {
    await this.delete('pendingActions', id);
  }

  // Sync metadata management
  async updateSyncMetadata(store: string, metadata: {
    last_sync: string;
    sync_token?: string;
  }): Promise<void> {
    await this.put('syncMetadata', {
      store,
      ...metadata,
      updated_at: new Date().toISOString()
    });
  }

  async getSyncMetadata(store: string): Promise<{ lastSync: string; version: number } | null> {
    const result = await this.get('syncMetadata', store) as { lastSync: string; version: number } | undefined;
    return result || null;
  }

  // Bulk operations for sync
  async bulkPut<T>(storeName: string, items: T[]): Promise<void> {
    const db = await this.getDB();
    const transaction = db.transaction([storeName], 'readwrite');
    const store = transaction.objectStore(storeName);

    return new Promise((resolve, reject) => {
      items.forEach(item => {
        store.put({
          ...item,
          sync_status: 'synced',
          synced_at: new Date().toISOString()
        });
      });

      transaction.oncomplete = () => resolve();
      transaction.onerror = () => reject(transaction.error);
    });
  }

  // Search functionality
  async search<T>(storeName: string, searchTerm: string, fields: string[]): Promise<T[]> {
    const allItems = await this.getAll<T>(storeName);
    const lowercaseSearch = searchTerm.toLowerCase();

    return allItems.filter(item => {
      return fields.some(field => {
        const value = (item as Record<string, unknown>)[field];
        if (typeof value === 'string') {
          return value.toLowerCase().includes(lowercaseSearch);
        }
        return false;
      });
    });
  }

  // Cleanup old data
  async cleanupOldData(storeName: string, daysToKeep: number = 30): Promise<void> {
    const cutoffDate = new Date();
    cutoffDate.setDate(cutoffDate.getDate() - daysToKeep);
    
    const allItems = await this.getAll(storeName);
    const db = await this.getDB();
    const transaction = db.transaction([storeName], 'readwrite');
    const store = transaction.objectStore(storeName);

    allItems.forEach(item => {
      const itemRecord = item as Record<string, unknown>;
      const updatedAt = new Date((itemRecord.updated_at || itemRecord.created_at) as string);
      if (updatedAt < cutoffDate && itemRecord.sync_status === 'synced') {
        store.delete(itemRecord.id as string);
      }
    });
  }

  // Database size estimation
  async estimateStorageUsage(): Promise<{ usage: number; quota: number }> {
    if ('storage' in navigator && 'estimate' in navigator.storage) {
      const estimate = await navigator.storage.estimate();
      return {
        usage: estimate.usage || 0,
        quota: estimate.quota || 0
      };
    }
    return { usage: 0, quota: 0 };
  }

  // Close database connection
  close(): void {
    if (this.db) {
      this.db.close();
      this.db = null;
      this.dbPromise = null;
    }
  }
}

// Singleton instance
export const indexedDBManager = new IndexedDBManager();