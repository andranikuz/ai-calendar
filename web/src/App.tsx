import React, { useEffect, Suspense, lazy } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Provider } from 'react-redux';
import { ConfigProvider, App as AntApp } from 'antd';
import { store } from './store';
import { useAppDispatch, useAppSelector } from './hooks/redux';
import { getCurrentUser } from './store/slices/authSlice';
import { authService } from './services/authService';

// Critical components (loaded immediately)
import Layout from './components/Layout/Layout';
import LoadingSpinner from './components/Common/LoadingSpinner';

// Lazy loaded pages
const LoginPage = lazy(() => import('./pages/LoginPage'));
const RegisterPage = lazy(() => import('./pages/RegisterPage'));
const DashboardPage = lazy(() => import('./pages/DashboardPage'));
const CalendarPage = lazy(() => import('./pages/CalendarPage'));
const GoalsPage = lazy(() => import('./pages/GoalsPage'));
const MoodsPage = lazy(() => import('./pages/MoodsPage'));
const SettingsPage = lazy(() => import('./pages/SettingsPage'));

// Styles
import './App.css';

const AppContent: React.FC = () => {
  const dispatch = useAppDispatch();
  const { isAuthenticated, isLoading, user } = useAppSelector(state => state.auth);

  useEffect(() => {
    // Check if user has stored tokens and try to get current user
    const tokens = authService.getStoredTokens();
    if (tokens && !user) {
      dispatch(getCurrentUser());
    }
  }, [dispatch, user]);

  if (isLoading) {
    return <LoadingSpinner />;
  }

  return (
    <Router>
      <Suspense fallback={<LoadingSpinner />}>
        <Routes>
          {!isAuthenticated ? (
            // Unauthenticated routes
            <>
              <Route path="/login" element={<LoginPage />} />
              <Route path="/register" element={<RegisterPage />} />
              <Route path="*" element={<Navigate to="/login" replace />} />
            </>
          ) : (
            // Authenticated routes
            <Route path="/" element={<Layout />}>
              <Route index element={<Navigate to="/dashboard" replace />} />
              <Route path="/dashboard" element={<DashboardPage />} />
              <Route path="/calendar" element={<CalendarPage />} />
              <Route path="/goals" element={<GoalsPage />} />
              <Route path="/moods" element={<MoodsPage />} />
              <Route path="/settings" element={<SettingsPage />} />
              <Route path="*" element={<Navigate to="/dashboard" replace />} />
            </Route>
          )}
        </Routes>
      </Suspense>
    </Router>
  );
};

const App: React.FC = () => {
  return (
    <Provider store={store}>
      <ConfigProvider
        theme={{
          token: {
            colorPrimary: '#1677ff',
            borderRadius: 6,
          },
        }}
      >
        <AntApp>
          <AppContent />
        </AntApp>
      </ConfigProvider>
    </Provider>
  );
};

export default App;
