import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { Mood, CreateMoodRequest } from '../../types/api';
import { moodsService } from '../../services/moodsService';

interface MoodsState {
  moods: Mood[];
  currentMood: Mood | null;
  todayMood: Mood | null;
  stats: {
    average: number;
    total: number;
    trend: 'up' | 'down' | 'stable';
    weekAverage?: number;
    weeklyData?: Array<{ level: number; date: string } | null>;
    distribution?: { [key: number]: number };
  } | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: MoodsState = {
  moods: [],
  currentMood: null,
  todayMood: null,
  stats: null,
  isLoading: false,
  error: null,
};

// Async thunks
export const fetchMoods = createAsyncThunk(
  'moods/fetchMoods',
  async (params: { start?: string; end?: string } = {}, { rejectWithValue }) => {
    try {
      const response = await moodsService.getMoods(params);
      return response.moods;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch moods';
      return rejectWithValue(errorMessage);
    }
  }
);

export const getTodayMood = createAsyncThunk(
  'moods/getTodayMood',
  async (_, { rejectWithValue }) => {
    try {
      const mood = await moodsService.getTodayMood();
      return mood;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to get today mood';
      return rejectWithValue(errorMessage);
    }
  }
);

export const createMood = createAsyncThunk(
  'moods/createMood',
  async (moodData: CreateMoodRequest, { rejectWithValue }) => {
    try {
      const mood = await moodsService.createMood(moodData);
      return mood;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create mood';
      return rejectWithValue(errorMessage);
    }
  }
);

export const updateMood = createAsyncThunk(
  'moods/updateMood',
  async ({ id, data }: { id: string; data: Partial<Mood> }, { rejectWithValue }) => {
    try {
      const mood = await moodsService.updateMood(id, data);
      return mood;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update mood';
      return rejectWithValue(errorMessage);
    }
  }
);

export const getMoodStats = createAsyncThunk(
  'moods/getMoodStats',
  async (params: { days?: number } = {}, { rejectWithValue }) => {
    try {
      const stats = await moodsService.getMoodStats(params);
      return stats;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to get mood stats';
      return rejectWithValue(errorMessage);
    }
  }
);

export const getMoodsByDateRange = createAsyncThunk(
  'moods/getMoodsByDateRange',
  async (params: { startDate: string; endDate: string }, { rejectWithValue }) => {
    try {
      const response = await moodsService.getMoodsByDateRange(params.startDate, params.endDate);
      return response.moods;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to get moods by date range';
      return rejectWithValue(errorMessage);
    }
  }
);

const moodsSlice = createSlice({
  name: 'moods',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    setCurrentMood: (state, action: PayloadAction<Mood | null>) => {
      state.currentMood = action.payload;
    },
    addMoodEntry: (state, action: PayloadAction<Mood>) => {
      const existingIndex = state.moods.findIndex(m => m.date === action.payload.date);
      if (existingIndex !== -1) {
        state.moods[existingIndex] = action.payload;
      } else {
        state.moods.push(action.payload);
      }
      
      // Update today's mood if it's today's entry
      const today = new Date().toISOString().split('T')[0];
      if (action.payload.date === today) {
        state.todayMood = action.payload;
      }
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch moods
      .addCase(fetchMoods.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchMoods.fulfilled, (state, action) => {
        state.isLoading = false;
        state.moods = action.payload;
      })
      .addCase(fetchMoods.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Get today mood
      .addCase(getTodayMood.fulfilled, (state, action) => {
        state.todayMood = action.payload;
      })
      // Create mood
      .addCase(createMood.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(createMood.fulfilled, (state, action) => {
        state.isLoading = false;
        const existingIndex = state.moods.findIndex(m => m.date === action.payload.date);
        if (existingIndex !== -1) {
          state.moods[existingIndex] = action.payload;
        } else {
          state.moods.push(action.payload);
        }
        
        // Update today's mood if it's today's entry
        const today = new Date().toISOString().split('T')[0];
        if (action.payload.date === today) {
          state.todayMood = action.payload;
        }
      })
      .addCase(createMood.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Update mood
      .addCase(updateMood.fulfilled, (state, action) => {
        const index = state.moods.findIndex(m => m.id === action.payload.id);
        if (index !== -1) {
          state.moods[index] = action.payload;
        }
        if (state.currentMood?.id === action.payload.id) {
          state.currentMood = action.payload;
        }
        
        // Update today's mood if it's today's entry
        const today = new Date().toISOString().split('T')[0];
        if (action.payload.date === today) {
          state.todayMood = action.payload;
        }
      })
      // Get mood stats
      .addCase(getMoodStats.fulfilled, (state, action) => {
        state.stats = action.payload;
      });
  },
});

export const { clearError, setCurrentMood, addMoodEntry } = moodsSlice.actions;
export default moodsSlice.reducer;