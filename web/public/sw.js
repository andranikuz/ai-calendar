const CACHE_NAME = 'smart-goal-calendar-v1';
const API_CACHE_NAME = 'sgc-api-cache-v1';

// Files to cache for offline functionality
const STATIC_CACHE_FILES = [
  '/',
  '/index.html',
  '/manifest.json',
  '/offline.html',
  // Add main JS and CSS files (will be added dynamically during build)
];

// API endpoints to cache
const API_CACHE_PATTERNS = [
  '/api/v1/users/me',
  '/api/v1/goals',
  '/api/v1/events',
  '/api/v1/moods',
];

// Install event - cache static files
self.addEventListener('install', (event) => {
  console.log('Service Worker installing...');
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => {
        console.log('Caching static files');
        return cache.addAll(STATIC_CACHE_FILES);
      })
      .then(() => {
        console.log('Service Worker installed successfully');
        return self.skipWaiting();
      })
  );
});

// Activate event - clean up old caches
self.addEventListener('activate', (event) => {
  console.log('Service Worker activating...');
  event.waitUntil(
    caches.keys()
      .then((cacheNames) => {
        return Promise.all(
          cacheNames.map((cacheName) => {
            if (cacheName !== CACHE_NAME && cacheName !== API_CACHE_NAME) {
              console.log('Deleting old cache:', cacheName);
              return caches.delete(cacheName);
            }
          })
        );
      })
      .then(() => {
        console.log('Service Worker activated');
        return self.clients.claim();
      })
  );
});

// Fetch event - implement caching strategies
self.addEventListener('fetch', (event) => {
  const { request } = event;
  const url = new URL(request.url);

  // Handle API requests
  if (url.pathname.startsWith('/api/')) {
    event.respondWith(handleApiRequest(request));
    return;
  }

  // Handle static files
  if (request.method === 'GET') {
    event.respondWith(handleStaticRequest(request));
    return;
  }
});

// Handle API requests with network-first strategy and offline fallback
async function handleApiRequest(request) {
  const url = new URL(request.url);
  const method = request.method;
  
  // Check if this API endpoint should be cached
  const shouldCache = API_CACHE_PATTERNS.some(pattern => 
    url.pathname.includes(pattern)
  );

  // Handle write operations (POST, PUT, DELETE) when offline
  if (['POST', 'PUT', 'DELETE'].includes(method)) {
    try {
      const networkResponse = await fetch(request.clone());
      return networkResponse;
    } catch (error) {
      console.log('Network request failed for write operation, queuing for later:', error);
      
      // Queue the action for background sync
      const requestBody = await request.clone().text();
      const action = {
        url: request.url,
        method: method,
        headers: Object.fromEntries(request.headers.entries()),
        body: requestBody,
        type: method.toLowerCase()
      };
      
      const actionId = await addPendingAction(action);
      
      // Trigger background sync
      if ('serviceWorker' in self && 'sync' in self.registration) {
        try {
          await self.registration.sync.register('sync-data');
        } catch (syncError) {
          console.log('Background sync not available:', syncError);
        }
      }
      
      // Return optimistic response for user feedback
      return new Response(
        JSON.stringify({ 
          success: true,
          offline: true,
          message: 'Action queued for sync when online',
          actionId: actionId
        }),
        { 
          status: 202, // Accepted
          headers: { 'Content-Type': 'application/json' }
        }
      );
    }
  }

  // Handle read operations (GET)
  if (!shouldCache) {
    // For non-cacheable API requests, just try network
    try {
      return await fetch(request);
    } catch (error) {
      console.log('Network request failed:', error);
      return new Response(
        JSON.stringify({ error: 'Network unavailable' }),
        { 
          status: 503,
          headers: { 'Content-Type': 'application/json' }
        }
      );
    }
  }

  try {
    // Try network first
    const networkResponse = await fetch(request.clone());
    
    // If successful, cache the response and store in IndexedDB
    if (networkResponse.ok) {
      const cache = await caches.open(API_CACHE_NAME);
      cache.put(request, networkResponse.clone());
      
      // Store data in IndexedDB for offline access
      try {
        const responseData = await networkResponse.clone().json();
        await storeDataInIndexedDB(url.pathname, responseData);
      } catch (dbError) {
        console.log('Failed to store in IndexedDB:', dbError);
      }
    }
    
    return networkResponse;
  } catch (error) {
    console.log('Network request failed, trying cache and IndexedDB:', error);
    
    // Try IndexedDB first for fresh data
    try {
      const offlineData = await getDataFromIndexedDB(url.pathname);
      if (offlineData) {
        return new Response(
          JSON.stringify(offlineData),
          { 
            status: 200,
            headers: { 
              'Content-Type': 'application/json',
              'X-Offline-Data': 'true'
            }
          }
        );
      }
    } catch (dbError) {
      console.log('IndexedDB lookup failed:', dbError);
    }
    
    // Fallback to HTTP cache
    const cachedResponse = await caches.match(request);
    if (cachedResponse) {
      return cachedResponse;
    }
    
    // If no cache, return error response
    return new Response(
      JSON.stringify({ error: 'Data unavailable offline' }),
      { 
        status: 503,
        headers: { 'Content-Type': 'application/json' }
      }
    );
  }
}

// Store API response data in IndexedDB
async function storeDataInIndexedDB(endpoint, data) {
  try {
    const db = await openIndexedDB();
    
    // Determine which store to use based on endpoint
    let storeName = null;
    if (endpoint.includes('/goals')) storeName = 'goals';
    else if (endpoint.includes('/events')) storeName = 'events';
    else if (endpoint.includes('/moods')) storeName = 'moods';
    
    if (!storeName) return;
    
    const transaction = db.transaction([storeName], 'readwrite');
    const store = transaction.objectStore(storeName);
    
    // Handle both arrays and single objects
    const items = Array.isArray(data) ? data : [data];
    
    items.forEach(item => {
      store.put({
        ...item,
        sync_status: 'synced',
        synced_at: new Date().toISOString()
      });
    });
  } catch (error) {
    console.error('Failed to store data in IndexedDB:', error);
  }
}

// Retrieve data from IndexedDB
async function getDataFromIndexedDB(endpoint) {
  try {
    const db = await openIndexedDB();
    
    // Determine which store to use based on endpoint
    let storeName = null;
    if (endpoint.includes('/goals')) storeName = 'goals';
    else if (endpoint.includes('/events')) storeName = 'events';
    else if (endpoint.includes('/moods')) storeName = 'moods';
    
    if (!storeName) return null;
    
    const transaction = db.transaction([storeName], 'readonly');
    const store = transaction.objectStore(storeName);
    
    return new Promise((resolve) => {
      const request = store.getAll();
      request.onsuccess = () => {
        const data = request.result || [];
        resolve(data.length === 1 ? data[0] : data);
      };
      request.onerror = () => resolve(null);
    });
  } catch (error) {
    console.error('Failed to get data from IndexedDB:', error);
    return null;
  }
}

// Handle static requests with cache-first strategy
async function handleStaticRequest(request) {
  try {
    // Try cache first
    const cachedResponse = await caches.match(request);
    if (cachedResponse) {
      return cachedResponse;
    }
    
    // If not in cache, try network
    const networkResponse = await fetch(request);
    
    // Cache successful responses
    if (networkResponse.ok) {
      const cache = await caches.open(CACHE_NAME);
      cache.put(request, networkResponse.clone());
    }
    
    return networkResponse;
  } catch (error) {
    console.log('Request failed:', error);
    
    // For navigation requests, try index.html first, then offline page
    if (request.mode === 'navigate') {
      const cachedIndex = await caches.match('/index.html');
      if (cachedIndex) {
        return cachedIndex;
      }
      
      // Fallback to offline page
      const offlinePage = await caches.match('/offline.html');
      if (offlinePage) {
        return offlinePage;
      }
    }
    
    // Return offline page or error
    const offlinePage = await caches.match('/offline.html');
    if (offlinePage) {
      return offlinePage;
    }
    
    return new Response('Offline', { status: 503 });
  }
}

// Handle background sync
self.addEventListener('sync', (event) => {
  console.log('Background sync triggered:', event.tag);
  
  if (event.tag === 'sync-data') {
    event.waitUntil(syncPendingData());
  }
});

// Sync pending data when back online
async function syncPendingData() {
  try {
    // Get pending data from IndexedDB or localStorage
    const pendingActions = await getPendingActions();
    
    for (const action of pendingActions) {
      try {
        await fetch(action.url, {
          method: action.method,
          headers: action.headers,
          body: action.body
        });
        
        // Remove from pending actions if successful
        await removePendingAction(action.id);
      } catch (error) {
        console.log('Failed to sync action:', action, error);
      }
    }
  } catch (error) {
    console.log('Sync failed:', error);
  }
}

// IndexedDB functions for pending actions
async function getPendingActions() {
  try {
    const db = await openIndexedDB();
    const transaction = db.transaction(['pendingActions'], 'readonly');
    const store = transaction.objectStore('pendingActions');
    return new Promise((resolve) => {
      const request = store.getAll();
      request.onsuccess = () => resolve(request.result || []);
      request.onerror = () => resolve([]);
    });
  } catch (error) {
    console.error('Failed to get pending actions:', error);
    return [];
  }
}

async function removePendingAction(id) {
  try {
    const db = await openIndexedDB();
    const transaction = db.transaction(['pendingActions'], 'readwrite');
    const store = transaction.objectStore('pendingActions');
    return new Promise((resolve) => {
      const request = store.delete(id);
      request.onsuccess = () => resolve();
      request.onerror = () => resolve();
    });
  } catch (error) {
    console.error('Failed to remove pending action:', error);
  }
}

async function addPendingAction(action) {
  try {
    const db = await openIndexedDB();
    const transaction = db.transaction(['pendingActions'], 'readwrite');
    const store = transaction.objectStore('pendingActions');
    const actionWithId = {
      ...action,
      id: Date.now() + Math.random(),
      created_at: new Date().toISOString(),
      retry_count: 0
    };
    return new Promise((resolve) => {
      const request = store.add(actionWithId);
      request.onsuccess = () => resolve(actionWithId.id);
      request.onerror = () => resolve(null);
    });
  } catch (error) {
    console.error('Failed to add pending action:', error);
    return null;
  }
}

// IndexedDB initialization
async function openIndexedDB() {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open('SmartGoalCalendarDB', 1);
    
    request.onerror = () => reject(request.error);
    request.onsuccess = () => resolve(request.result);
    
    request.onupgradeneeded = (event) => {
      const db = event.target.result;
      
      // Create pendingActions store
      if (!db.objectStoreNames.contains('pendingActions')) {
        const store = db.createObjectStore('pendingActions', { 
          keyPath: 'id',
          autoIncrement: true 
        });
        store.createIndex('type', 'type');
        store.createIndex('created_at', 'created_at');
      }
      
      // Create offline data stores
      if (!db.objectStoreNames.contains('goals')) {
        const goalsStore = db.createObjectStore('goals', { keyPath: 'id' });
        goalsStore.createIndex('user_id', 'user_id');
        goalsStore.createIndex('sync_status', 'sync_status');
      }
      
      if (!db.objectStoreNames.contains('events')) {
        const eventsStore = db.createObjectStore('events', { keyPath: 'id' });
        eventsStore.createIndex('user_id', 'user_id');
        eventsStore.createIndex('sync_status', 'sync_status');
      }
      
      if (!db.objectStoreNames.contains('moods')) {
        const moodsStore = db.createObjectStore('moods', { keyPath: 'id' });
        moodsStore.createIndex('user_id', 'user_id');
        moodsStore.createIndex('sync_status', 'sync_status');
      }
    };
  });
}

// Handle push notifications
self.addEventListener('push', (event) => {
  console.log('Push notification received:', event);
  
  const options = {
    body: 'You have new updates in Smart Goal Calendar',
    icon: '/icons/icon-192x192.png',
    badge: '/icons/badge-72x72.png',
    data: {
      url: '/'
    },
    actions: [
      {
        action: 'open',
        title: 'Open App'
      },
      {
        action: 'dismiss',
        title: 'Dismiss'
      }
    ]
  };

  if (event.data) {
    try {
      const data = event.data.json();
      options.body = data.message || options.body;
      options.data = data;
    } catch (error) {
      console.log('Failed to parse push data:', error);
    }
  }

  event.waitUntil(
    self.registration.showNotification('Smart Goal Calendar', options)
  );
});

// Handle notification clicks
self.addEventListener('notificationclick', (event) => {
  console.log('Notification clicked:', event);
  
  event.notification.close();
  
  if (event.action === 'dismiss') {
    return;
  }
  
  const url = event.notification.data?.url || '/';
  
  event.waitUntil(
    clients.openWindow(url)
  );
});

// Handle message from main thread
self.addEventListener('message', (event) => {
  console.log('Message received in SW:', event.data);
  
  if (event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }
  
  if (event.data.type === 'GET_VERSION') {
    event.ports[0].postMessage({ version: CACHE_NAME });
  }
});

console.log('Service Worker script loaded');