import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { apiService } from '../../services/api';

interface CalendarEvent {
  id: string;
  title: string;
  start_time: string;
  end_time: string;
  description?: string;
  location?: string;
}

export interface SyncConflict {
  id: string;
  user_id: string;
  calendar_sync_id: string;
  conflict_type: 'time_overlap' | 'content_diff' | 'duplicate_event' | 'deleted_event';
  local_event?: CalendarEvent;
  google_event?: CalendarEvent;
  description: string;
  resolution?: string;
  resolved_at?: string;
  resolved_by?: string;
  status: 'pending' | 'resolved' | 'ignored';
  created_at: string;
  updated_at: string;
}

export interface ConflictStats {
  stats: Record<string, number>;
  total: number;
  period: number;
}

export interface ConflictResolutionAction {
  action: 'use_local' | 'use_google' | 'merge' | 'ignore';
  event_data?: Record<string, unknown>;
  resolution?: string;
}

interface SyncConflictsState {
  conflicts: SyncConflict[];
  stats: ConflictStats | null;
  isLoading: boolean;
  error: string | null;
  resolving: Record<string, boolean>; // Track which conflicts are being resolved
}

const initialState: SyncConflictsState = {
  conflicts: [],
  stats: null,
  isLoading: false,
  error: null,
  resolving: {},
};

// Async thunks
export const fetchPendingConflicts = createAsyncThunk(
  'syncConflicts/fetchPending',
  async (_, { rejectWithValue }) => {
    try {
      const response = await apiService.get<{
        conflicts: SyncConflict[];
        count: number;
      }>('/sync-conflicts');
      return response.conflicts;
    } catch (error) {
      return rejectWithValue('Failed to fetch pending conflicts');
    }
  }
);

export const fetchConflictStats = createAsyncThunk(
  'syncConflicts/fetchStats',
  async (days: number = 30, { rejectWithValue }) => {
    try {
      const response = await apiService.get<ConflictStats>(`/sync-conflicts/stats?days=${days}`);
      return response;
    } catch (error) {
      return rejectWithValue('Failed to fetch conflict statistics');
    }
  }
);

export const resolveConflict = createAsyncThunk(
  'syncConflicts/resolve',
  async (
    { conflictId, action }: { conflictId: string; action: ConflictResolutionAction },
    { rejectWithValue }
  ) => {
    try {
      await apiService.post(`/sync-conflicts/${conflictId}/resolve`, action);
      return conflictId;
    } catch (error) {
      return rejectWithValue('Failed to resolve conflict');
    }
  }
);

export const bulkResolveConflicts = createAsyncThunk(
  'syncConflicts/bulkResolve',
  async (
    {
      conflictIds,
      action,
      resolution,
    }: {
      conflictIds: string[];
      action: 'use_local' | 'use_google' | 'ignore';
      resolution?: string;
    },
    { rejectWithValue }
  ) => {
    try {
      const response = await apiService.post<{
        successful: number;
        failed: number;
        results: Array<{ conflict_id: string; success: boolean; error?: string }>;
      }>('/sync-conflicts/bulk-resolve', {
        conflict_ids: conflictIds,
        action,
        resolution,
      });
      return { conflictIds, response };
    } catch (error) {
      return rejectWithValue('Failed to bulk resolve conflicts');
    }
  }
);

const syncConflictsSlice = createSlice({
  name: 'syncConflicts',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    setResolvingStatus: (state, action: PayloadAction<{ conflictId: string; isResolving: boolean }>) => {
      const { conflictId, isResolving } = action.payload;
      if (isResolving) {
        state.resolving[conflictId] = true;
      } else {
        delete state.resolving[conflictId];
      }
    },
    removeConflict: (state, action: PayloadAction<string>) => {
      state.conflicts = state.conflicts.filter(conflict => conflict.id !== action.payload);
    }
  },
  extraReducers: (builder) => {
    // Fetch pending conflicts
    builder
      .addCase(fetchPendingConflicts.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchPendingConflicts.fulfilled, (state, action) => {
        state.isLoading = false;
        state.conflicts = action.payload;
      })
      .addCase(fetchPendingConflicts.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });

    // Fetch conflict stats
    builder
      .addCase(fetchConflictStats.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchConflictStats.fulfilled, (state, action) => {
        state.isLoading = false;
        state.stats = action.payload;
      })
      .addCase(fetchConflictStats.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });

    // Resolve conflict
    builder
      .addCase(resolveConflict.pending, (state, action) => {
        const conflictId = action.meta.arg.conflictId;
        state.resolving[conflictId] = true;
        state.error = null;
      })
      .addCase(resolveConflict.fulfilled, (state, action) => {
        const conflictId = action.payload;
        delete state.resolving[conflictId];
        // Remove resolved conflict from list
        state.conflicts = state.conflicts.filter(conflict => conflict.id !== conflictId);
      })
      .addCase(resolveConflict.rejected, (state, action) => {
        const conflictId = action.meta.arg.conflictId;
        delete state.resolving[conflictId];
        state.error = action.payload as string;
      });

    // Bulk resolve conflicts
    builder
      .addCase(bulkResolveConflicts.pending, (state, action) => {
        const conflictIds = action.meta.arg.conflictIds;
        conflictIds.forEach(id => {
          state.resolving[id] = true;
        });
        state.error = null;
      })
      .addCase(bulkResolveConflicts.fulfilled, (state, action) => {
        const { conflictIds, response } = action.payload;
        
        // Clear resolving status for all conflicts
        conflictIds.forEach(id => {
          delete state.resolving[id];
        });

        // Remove successfully resolved conflicts
        const successfulIds = response.results
          .filter(result => result.success)
          .map(result => result.conflict_id);
        
        state.conflicts = state.conflicts.filter(
          conflict => !successfulIds.includes(conflict.id)
        );
      })
      .addCase(bulkResolveConflicts.rejected, (state, action) => {
        const conflictIds = action.meta.arg.conflictIds;
        conflictIds.forEach(id => {
          delete state.resolving[id];
        });
        state.error = action.payload as string;
      });
  },
});

export const { clearError, setResolvingStatus, removeConflict } = syncConflictsSlice.actions;
export default syncConflictsSlice.reducer;