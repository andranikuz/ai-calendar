import { configureStore } from '@reduxjs/toolkit';
import authSlice from './slices/authSlice';
import goalsSlice from './slices/goalsSlice';
import eventsSlice from './slices/eventsSlice';
import moodsSlice from './slices/moodsSlice';
import googleSlice from './slices/googleSlice';
import syncConflictsSlice from './slices/syncConflictsSlice';

export const store = configureStore({
  reducer: {
    auth: authSlice,
    goals: goalsSlice,
    events: eventsSlice,
    moods: moodsSlice,
    google: googleSlice,
    syncConflicts: syncConflictsSlice,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ['persist/PERSIST'],
      },
    }),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;