import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { Event, CreateEventRequest } from '../../types/api';
import { eventsService } from '../../services/eventsService';
import dayjs from 'dayjs';

interface EventsState {
  events: Event[];
  currentEvent: Event | null;
  calendarView: 'month' | 'week' | 'day';
  currentDate: string;
  isLoading: boolean;
  error: string | null;
}

const initialState: EventsState = {
  events: [],
  currentEvent: null,
  calendarView: 'month',
  currentDate: dayjs().format('YYYY-MM-DD'),
  isLoading: false,
  error: null,
};

// Async thunks
export const fetchEvents = createAsyncThunk(
  'events/fetchEvents',
  async (params: { start?: string; end?: string } = {}, { rejectWithValue }) => {
    try {
      const response = await eventsService.getEvents(params);
      return response.events;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch events';
      return rejectWithValue(errorMessage);
    }
  }
);

export const createEvent = createAsyncThunk(
  'events/createEvent',
  async (eventData: CreateEventRequest, { rejectWithValue }) => {
    try {
      const event = await eventsService.createEvent(eventData);
      return event;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create event';
      return rejectWithValue(errorMessage);
    }
  }
);

export const updateEvent = createAsyncThunk(
  'events/updateEvent',
  async ({ id, data }: { id: string; data: Partial<Event> }, { rejectWithValue }) => {
    try {
      const event = await eventsService.updateEvent(id, data);
      return event;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update event';
      return rejectWithValue(errorMessage);
    }
  }
);

export const deleteEvent = createAsyncThunk(
  'events/deleteEvent',
  async (id: string, { rejectWithValue }) => {
    try {
      await eventsService.deleteEvent(id);
      return id;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to delete event';
      return rejectWithValue(errorMessage);
    }
  }
);

export const moveEvent = createAsyncThunk(
  'events/moveEvent',
  async ({ id, start, end }: { id: string; start: string; end: string }, { rejectWithValue }) => {
    try {
      const event = await eventsService.moveEvent(id, start, end);
      return event;
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to move event';
      return rejectWithValue(errorMessage);
    }
  }
);

const eventsSlice = createSlice({
  name: 'events',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    setCurrentEvent: (state, action: PayloadAction<Event | null>) => {
      state.currentEvent = action.payload;
    },
    setCalendarView: (state, action: PayloadAction<'month' | 'week' | 'day'>) => {
      state.calendarView = action.payload;
    },
    setCurrentDate: (state, action: PayloadAction<string>) => {
      state.currentDate = action.payload;
    },
    navigateDate: (state, action: PayloadAction<'prev' | 'next' | 'today'>) => {
      const current = dayjs(state.currentDate);
      switch (action.payload) {
        case 'prev':
          if (state.calendarView === 'month') {
            state.currentDate = current.subtract(1, 'month').format('YYYY-MM-DD');
          } else if (state.calendarView === 'week') {
            state.currentDate = current.subtract(1, 'week').format('YYYY-MM-DD');
          } else {
            state.currentDate = current.subtract(1, 'day').format('YYYY-MM-DD');
          }
          break;
        case 'next':
          if (state.calendarView === 'month') {
            state.currentDate = current.add(1, 'month').format('YYYY-MM-DD');
          } else if (state.calendarView === 'week') {
            state.currentDate = current.add(1, 'week').format('YYYY-MM-DD');
          } else {
            state.currentDate = current.add(1, 'day').format('YYYY-MM-DD');
          }
          break;
        case 'today':
          state.currentDate = dayjs().format('YYYY-MM-DD');
          break;
      }
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch events
      .addCase(fetchEvents.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchEvents.fulfilled, (state, action) => {
        state.isLoading = false;
        state.events = action.payload;
      })
      .addCase(fetchEvents.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Create event
      .addCase(createEvent.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(createEvent.fulfilled, (state, action) => {
        state.isLoading = false;
        state.events.push(action.payload);
      })
      .addCase(createEvent.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Update event
      .addCase(updateEvent.fulfilled, (state, action) => {
        const index = state.events.findIndex(e => e.id === action.payload.id);
        if (index !== -1) {
          state.events[index] = action.payload;
        }
        if (state.currentEvent?.id === action.payload.id) {
          state.currentEvent = action.payload;
        }
      })
      // Delete event
      .addCase(deleteEvent.fulfilled, (state, action) => {
        state.events = state.events.filter(e => e.id !== action.payload);
        if (state.currentEvent?.id === action.payload) {
          state.currentEvent = null;
        }
      })
      // Move event
      .addCase(moveEvent.fulfilled, (state, action) => {
        const index = state.events.findIndex(e => e.id === action.payload.id);
        if (index !== -1) {
          state.events[index] = action.payload;
        }
      });
  },
});

export const { 
  clearError, 
  setCurrentEvent, 
  setCalendarView, 
  setCurrentDate, 
  navigateDate 
} = eventsSlice.actions;
export default eventsSlice.reducer;