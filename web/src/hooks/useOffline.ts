import { useState, useEffect } from 'react';
import { indexedDBManager } from '../utils/indexedDB';

export interface OfflineAction {
  id: number;
  type: 'create' | 'update' | 'delete';
  store: string;
  data: Record<string, unknown>;
  url: string;
  method: string;
  created_at: string;
  retry_count: number;
}

export interface OfflineHookReturn {
  isOnline: boolean;
  isOfflineMode: boolean;
  pendingActions: OfflineAction[];
  syncStatus: 'idle' | 'syncing' | 'error';
  offlineData: {
    goals: Array<Record<string, unknown>>;
    events: Array<Record<string, unknown>>;
    moods: Array<Record<string, unknown>>;
  };
  addOfflineAction: (action: Omit<OfflineAction, 'id' | 'created_at' | 'retry_count'>) => Promise<void>;
  syncPendingActions: () => Promise<void>;
  clearOfflineData: () => Promise<void>;
  refreshOfflineData: () => Promise<void>;
}

export const useOffline = (): OfflineHookReturn => {
  const [isOnline, setIsOnline] = useState(navigator.onLine);
  const [pendingActions, setPendingActions] = useState<OfflineAction[]>([]);
  const [syncStatus, setSyncStatus] = useState<'idle' | 'syncing' | 'error'>('idle');
  const [offlineData, setOfflineData] = useState({
    goals: [] as Array<Record<string, unknown>>,
    events: [] as Array<Record<string, unknown>>,
    moods: [] as Array<Record<string, unknown>>
  });

  // Track online/offline status
  useEffect(() => {
    const handleOnline = () => {
      setIsOnline(true);
      // Auto-sync when coming back online
      syncPendingActions();
    };

    const handleOffline = () => {
      setIsOnline(false);
    };

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  // Load pending actions and offline data on mount
  useEffect(() => {
    loadPendingActions();
    refreshOfflineData();
  }, []);

  // Listen for service worker messages
  useEffect(() => {
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.addEventListener('message', (event) => {
        if (event.data.type === 'OFFLINE_ACTION_ADDED') {
          loadPendingActions();
        }
        if (event.data.type === 'SYNC_COMPLETED') {
          loadPendingActions();
          refreshOfflineData();
        }
      });
    }
  }, []);

  const loadPendingActions = async () => {
    try {
      const actions = await indexedDBManager.getPendingActions();
      setPendingActions(actions as unknown as OfflineAction[]);
    } catch (error) {
      console.error('Failed to load pending actions:', error);
    }
  };

  const refreshOfflineData = async () => {
    try {
      const goals = await indexedDBManager.getAll('goals');
      const events = await indexedDBManager.getAll('events');
      const moods = await indexedDBManager.getAll('moods');

      setOfflineData({
        goals: (goals as unknown as Array<Record<string, unknown>>) || [],
        events: (events as unknown as Array<Record<string, unknown>>) || [],
        moods: (moods as unknown as Array<Record<string, unknown>>) || []
      });
    } catch (error) {
      console.error('Failed to refresh offline data:', error);
    }
  };

  const addOfflineAction = async (action: Omit<OfflineAction, 'id' | 'created_at' | 'retry_count'>) => {
    try {
      await indexedDBManager.addPendingAction({
        type: action.type,
        store: action.store,
        data: action.data,
        url: action.url,
        method: action.method
      });
      
      // Reload pending actions
      await loadPendingActions();
      
      // Try to register background sync
      if ('serviceWorker' in navigator) {
        try {
          const registration = await navigator.serviceWorker.ready;
          if ('sync' in registration) {
            await (registration as any).sync.register('sync-data');
          }
        } catch (error) {
          console.log('Background sync not supported:', error);
        }
      }
    } catch (error) {
      console.error('Failed to add offline action:', error);
      throw error;
    }
  };

  const syncPendingActions = async () => {
    if (!isOnline || syncStatus === 'syncing') return;

    setSyncStatus('syncing');
    
    try {
      const actions = await indexedDBManager.getPendingActions();
      
      for (const action of actions) {
        try {
          const response = await fetch(action.url as string, {
            method: action.method as string,
            headers: {
              'Content-Type': 'application/json',
              ...(action as any).headers
            },
            body: (action as any).body
          });

          if (response.ok) {
            // Remove successful action
            await indexedDBManager.removePendingAction(action.id as number);
            
            // Update offline data if it was a successful operation
            if (action.type === 'create' || action.type === 'update') {
              const responseData = await response.json();
              await updateOfflineData(action.store as string, responseData);
            } else if (action.type === 'delete') {
              await removeFromOfflineData(action.store as string, (action.data as any).id);
            }
          } else {
            // Increment retry count for failed actions
            (action as any).retry_count += 1;
            if ((action as any).retry_count < 3) {
              await indexedDBManager.put('pendingActions', action);
            } else {
              // Remove after 3 failed attempts
              await indexedDBManager.removePendingAction(action.id as number);
            }
          }
        } catch (error) {
          console.error('Failed to sync action:', action, error);
          // Increment retry count
          (action as any).retry_count += 1;
          if ((action as any).retry_count < 3) {
            await indexedDBManager.put('pendingActions', action);
          } else {
            await indexedDBManager.removePendingAction(action.id as number);
          }
        }
      }

      setSyncStatus('idle');
      await loadPendingActions();
      await refreshOfflineData();
    } catch (error) {
      console.error('Sync failed:', error);
      setSyncStatus('error');
      
      // Reset to idle after 5 seconds
      setTimeout(() => setSyncStatus('idle'), 5000);
    }
  };

  const updateOfflineData = async (store: string, data: Record<string, unknown>) => {
    try {
      await indexedDBManager.put(store, {
        ...data,
        sync_status: 'synced',
        synced_at: new Date().toISOString()
      });
    } catch (error) {
      console.error('Failed to update offline data:', error);
    }
  };

  const removeFromOfflineData = async (store: string, id: string) => {
    try {
      await indexedDBManager.delete(store, id);
    } catch (error) {
      console.error('Failed to remove from offline data:', error);
    }
  };

  const clearOfflineData = async () => {
    try {
      await indexedDBManager.clear('goals');
      await indexedDBManager.clear('events');
      await indexedDBManager.clear('moods');
      await indexedDBManager.clear('pendingActions');
      
      setOfflineData({
        goals: [],
        events: [],
        moods: []
      });
      setPendingActions([]);
    } catch (error) {
      console.error('Failed to clear offline data:', error);
      throw error;
    }
  };

  return {
    isOnline,
    isOfflineMode: !isOnline,
    pendingActions,
    syncStatus,
    offlineData,
    addOfflineAction,
    syncPendingActions,
    clearOfflineData,
    refreshOfflineData
  };
};