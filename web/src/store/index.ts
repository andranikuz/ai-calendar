import { configureStore } from '@reduxjs/toolkit';
import authSlice from './slices/authSlice';
import goalsSlice from './slices/goalsSlice';
import eventsSlice from './slices/eventsSlice';
import moodsSlice from './slices/moodsSlice';
import googleSlice from './slices/googleSlice';

export const store = configureStore({
  reducer: {
    auth: authSlice,
    goals: goalsSlice,
    events: eventsSlice,
    moods: moodsSlice,
    google: googleSlice,
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