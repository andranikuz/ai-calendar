import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { Goal, Task, Milestone, CreateGoalRequest } from '../../types/api';
import { goalsService } from '../../services/goalsService';

interface GoalsState {
  goals: Goal[];
  currentGoal: Goal | null;
  tasks: Task[];
  milestones: Milestone[];
  isLoading: boolean;
  error: string | null;
}

const initialState: GoalsState = {
  goals: [],
  currentGoal: null,
  tasks: [],
  milestones: [],
  isLoading: false,
  error: null,
};

// Async thunks
export const fetchGoals = createAsyncThunk(
  'goals/fetchGoals',
  async (_, { rejectWithValue }) => {
    try {
      const response = await goalsService.getGoals();
      return response.goals;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch goals';
      return rejectWithValue(errorMessage);
    }
  }
);

export const createGoal = createAsyncThunk(
  'goals/createGoal',
  async (goalData: CreateGoalRequest, { rejectWithValue }) => {
    try {
      const goal = await goalsService.createGoal(goalData);
      return goal;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create goal';
      return rejectWithValue(errorMessage);
    }
  }
);

export const updateGoal = createAsyncThunk(
  'goals/updateGoal',
  async ({ id, data }: { id: string; data: Partial<Goal> }, { rejectWithValue }) => {
    try {
      const goal = await goalsService.updateGoal(id, data);
      return goal;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update goal';
      return rejectWithValue(errorMessage);
    }
  }
);

export const deleteGoal = createAsyncThunk(
  'goals/deleteGoal',
  async (id: string, { rejectWithValue }) => {
    try {
      await goalsService.deleteGoal(id);
      return id;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to delete goal';
      return rejectWithValue(errorMessage);
    }
  }
);

export const fetchGoalTasks = createAsyncThunk(
  'goals/fetchGoalTasks',
  async (goalId: string, { rejectWithValue }) => {
    try {
      const response = await goalsService.getGoalTasks(goalId);
      return response.tasks;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch tasks';
      return rejectWithValue(errorMessage);
    }
  }
);

const goalsSlice = createSlice({
  name: 'goals',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    setCurrentGoal: (state, action: PayloadAction<Goal | null>) => {
      state.currentGoal = action.payload;
    },
    updateGoalProgress: (state, action: PayloadAction<{ id: string; progress: number }>) => {
      const goal = state.goals.find(g => g.id === action.payload.id);
      if (goal) {
        goal.progress = action.payload.progress;
      }
      if (state.currentGoal?.id === action.payload.id) {
        state.currentGoal.progress = action.payload.progress;
      }
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch goals
      .addCase(fetchGoals.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchGoals.fulfilled, (state, action) => {
        state.isLoading = false;
        state.goals = action.payload;
      })
      .addCase(fetchGoals.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Create goal
      .addCase(createGoal.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(createGoal.fulfilled, (state, action) => {
        state.isLoading = false;
        state.goals.push(action.payload);
      })
      .addCase(createGoal.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Update goal
      .addCase(updateGoal.fulfilled, (state, action) => {
        const index = state.goals.findIndex(g => g.id === action.payload.id);
        if (index !== -1) {
          state.goals[index] = action.payload;
        }
        if (state.currentGoal?.id === action.payload.id) {
          state.currentGoal = action.payload;
        }
      })
      // Delete goal
      .addCase(deleteGoal.fulfilled, (state, action) => {
        state.goals = state.goals.filter(g => g.id !== action.payload);
        if (state.currentGoal?.id === action.payload) {
          state.currentGoal = null;
        }
      })
      // Fetch goal tasks
      .addCase(fetchGoalTasks.fulfilled, (state, action) => {
        state.tasks = action.payload;
      });
  },
});

export const { clearError, setCurrentGoal, updateGoalProgress } = goalsSlice.actions;
export default goalsSlice.reducer;