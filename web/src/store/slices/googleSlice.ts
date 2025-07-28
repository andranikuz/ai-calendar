import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { GoogleIntegration, GoogleCalendarSync, GoogleCalendar } from '../../types/api';
import { googleService } from '../../services/googleService';

interface GoogleState {
  integration: GoogleIntegration | null;
  calendars: GoogleCalendar[];
  calendarSyncs: GoogleCalendarSync[];
  authUrl: string | null;
  isConnected: boolean;
  isLoading: boolean;
  error: string | null;
}

const initialState: GoogleState = {
  integration: null,
  calendars: [],
  calendarSyncs: [],
  authUrl: null,
  isConnected: false,
  isLoading: false,
  error: null,
};

// Async thunks
export const getAuthUrl = createAsyncThunk(
  'google/getAuthUrl',
  async (_, { rejectWithValue }) => {
    try {
      const response = await googleService.getAuthUrl();
      return response.auth_url;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to get auth URL';
      return rejectWithValue(errorMessage);
    }
  }
);

export const handleCallback = createAsyncThunk(
  'google/handleCallback',
  async ({ code, state }: { code: string; state: string }, { rejectWithValue }) => {
    try {
      const integration = await googleService.handleCallback(code, state);
      return integration;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to connect Google account';
      return rejectWithValue(errorMessage);
    }
  }
);

export const getIntegration = createAsyncThunk(
  'google/getIntegration',
  async (_, { rejectWithValue }) => {
    try {
      const integration = await googleService.getIntegration();
      return integration;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to get integration';
      return rejectWithValue(errorMessage);
    }
  }
);

export const disconnect = createAsyncThunk(
  'google/disconnect',
  async (_, { rejectWithValue }) => {
    try {
      await googleService.disconnect();
      return null;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to disconnect';
      return rejectWithValue(errorMessage);
    }
  }
);

export const getCalendars = createAsyncThunk(
  'google/getCalendars',
  async (_, { rejectWithValue }) => {
    try {
      const response = await googleService.getCalendars();
      return response.calendars;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to get calendars';
      return rejectWithValue(errorMessage);
    }
  }
);

export const getCalendarSyncs = createAsyncThunk(
  'google/getCalendarSyncs',
  async (_, { rejectWithValue }) => {
    try {
      const response = await googleService.getCalendarSyncs();
      return response.syncs;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to get calendar syncs';
      return rejectWithValue(errorMessage);
    }
  }
);

export const createCalendarSync = createAsyncThunk(
  'google/createCalendarSync',
  async (syncData: Partial<GoogleCalendarSync>, { rejectWithValue }) => {
    try {
      const sync = await googleService.createCalendarSync(syncData);
      return sync;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create calendar sync';
      return rejectWithValue(errorMessage);
    }
  }
);

export const updateCalendarSync = createAsyncThunk(
  'google/updateCalendarSync',
  async ({ id, data }: { id: string; data: Partial<GoogleCalendarSync> }, { rejectWithValue }) => {
    try {
      const sync = await googleService.updateCalendarSync(id, data);
      return sync;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update calendar sync';
      return rejectWithValue(errorMessage);
    }
  }
);

export const deleteCalendarSync = createAsyncThunk(
  'google/deleteCalendarSync',
  async (syncId: string, { rejectWithValue }) => {
    try {
      await googleService.deleteCalendarSync(syncId);
      return syncId;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to delete calendar sync';
      return rejectWithValue(errorMessage);
    }
  }
);

export const triggerSync = createAsyncThunk(
  'google/triggerSync',
  async (syncId: string, { rejectWithValue }) => {
    try {
      const result = await googleService.triggerSync(syncId);
      return { syncId, result };
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to trigger sync';
      return rejectWithValue(errorMessage);
    }
  }
);

export const triggerSyncWithConflictDetection = createAsyncThunk(
  'google/triggerSyncWithConflictDetection',
  async (syncId: string, { rejectWithValue }) => {
    try {
      const result = await googleService.triggerSyncWithConflictDetection(syncId);
      return { syncId, result };
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to trigger sync with conflict detection';
      return rejectWithValue(errorMessage);
    }
  }
);

// Aliases for better naming in components
export const fetchGoogleIntegration = getIntegration;
export const disconnectGoogle = disconnect;
export const fetchGoogleCalendars = getCalendars;
export const fetchCalendarSyncs = getCalendarSyncs;
export const manualSync = triggerSync;
export const syncWithConflicts = triggerSyncWithConflictDetection;

const googleSlice = createSlice({
  name: 'google',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    clearAuthUrl: (state) => {
      state.authUrl = null;
    },
    setConnected: (state, action: PayloadAction<boolean>) => {
      state.isConnected = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      // Get auth URL
      .addCase(getAuthUrl.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(getAuthUrl.fulfilled, (state, action) => {
        state.isLoading = false;
        state.authUrl = action.payload;
      })
      .addCase(getAuthUrl.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Handle callback
      .addCase(handleCallback.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(handleCallback.fulfilled, (state, action) => {
        state.isLoading = false;
        state.integration = action.payload;
        state.isConnected = true;
        state.authUrl = null;
      })
      .addCase(handleCallback.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Get integration
      .addCase(getIntegration.fulfilled, (state, action) => {
        state.integration = action.payload;
        state.isConnected = !!action.payload;
      })
      .addCase(getIntegration.rejected, (state) => {
        state.integration = null;
        state.isConnected = false;
      })
      // Disconnect
      .addCase(disconnect.fulfilled, (state) => {
        state.integration = null;
        state.calendars = [];
        state.calendarSyncs = [];
        state.isConnected = false;
      })
      // Get calendars
      .addCase(getCalendars.fulfilled, (state, action) => {
        state.calendars = action.payload;
      })
      // Get calendar syncs
      .addCase(getCalendarSyncs.fulfilled, (state, action) => {
        state.calendarSyncs = action.payload;
      })
      // Create calendar sync
      .addCase(createCalendarSync.fulfilled, (state, action) => {
        state.calendarSyncs.push(action.payload);
      })
      // Update calendar sync
      .addCase(updateCalendarSync.fulfilled, (state, action) => {
        const index = state.calendarSyncs.findIndex(s => s.id === action.payload.id);
        if (index !== -1) {
          state.calendarSyncs[index] = action.payload;
        }
      })
      // Delete calendar sync
      .addCase(deleteCalendarSync.fulfilled, (state, action) => {
        state.calendarSyncs = state.calendarSyncs.filter(s => s.id !== action.payload);
      })
      // Trigger sync
      .addCase(triggerSync.fulfilled, (state, action) => {
        const sync = state.calendarSyncs.find(s => s.id === action.payload.syncId);
        if (sync) {
          sync.last_sync_at = new Date().toISOString();
          sync.sync_status = 'active';
          sync.last_sync_error = '';
        }
      });
  },
});

export const { clearError, clearAuthUrl, setConnected } = googleSlice.actions;
export default googleSlice.reducer;