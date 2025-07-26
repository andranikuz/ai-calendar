// API Response types
export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  offset: number;
  limit: number;
}

// Auth types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
}

export interface User {
  id: string;
  email: string;
  name: string;
  profile: Record<string, any>;
  settings: Record<string, any>;
  created_at: string;
  updated_at: string;
}

// Goal types
export interface Goal {
  id: string;
  user_id: string;
  title: string;
  description: string;
  category: 'health' | 'career' | 'education' | 'personal' | 'financial' | 'relationship';
  priority: 'low' | 'medium' | 'high' | 'critical';
  status: 'draft' | 'active' | 'paused' | 'completed' | 'cancelled';
  progress: number;
  deadline?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateGoalRequest {
  title: string;
  description: string;
  category: Goal['category'];
  priority: Goal['priority'];
  deadline?: string;
}

export interface Task {
  id: string;
  goal_id: string;
  title: string;
  description: string;
  priority: 'low' | 'medium' | 'high' | 'critical';
  status: 'pending' | 'in_progress' | 'active' | 'completed' | 'cancelled';
  estimated_duration?: number;
  deadline?: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface Milestone {
  id: string;
  goal_id: string;
  title: string;
  description: string;
  deadline?: string;
  target_value?: string;
  status: 'active' | 'completed' | 'cancelled';
  completed_at?: string;
  created_at: string;
  updated_at: string;
}

// Event types
export interface Event {
  id: string;
  user_id: string;
  goal_id?: string;
  title: string;
  description: string;
  start_time: string;
  end_time: string;
  timezone: string;
  recurrence?: any;
  location?: string;
  attendees: any[];
  status: 'tentative' | 'confirmed' | 'cancelled';
  external_id?: string;
  external_source?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateEventRequest {
  title: string;
  description: string;
  start_time: string;
  end_time: string;
  timezone?: string;
  location?: string;
  goal_id?: string;
}

// Mood types
export interface Mood {
  id: string;
  user_id: string;
  date: string;
  level: number; // 1-5
  notes?: string;
  tags: string[];
  recorded_at: string;
}

export interface CreateMoodRequest {
  date: string;
  level: number;
  notes?: string;
  tags?: string[];
}

// Google Integration types
export interface GoogleIntegration {
  id: string;
  email: string;
  name: string;
  enabled: boolean;
  calendar_id: string;
  created_at: string;
}

export interface GoogleCalendarSync {
  id: string;
  calendar_id: string;
  calendar_name: string;
  sync_direction: 'bidirectional' | 'from_google' | 'to_google';
  sync_status: 'active' | 'paused' | 'error' | 'disabled';
  last_sync_at?: string;
  last_sync_error?: string;
  settings: {
    sync_interval: number;
    auto_sync: boolean;
    sync_past_events: boolean;
    sync_future_events: boolean;
    conflict_resolution: string;
  };
  created_at: string;
  updated_at: string;
}

export interface GoogleCalendar {
  id: string;
  summary: string;
  description?: string;
  primary: boolean;
  access_role: string;
}