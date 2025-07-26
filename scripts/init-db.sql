-- Initialize Smart Goal Calendar Database
-- This script creates the initial database structure

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types
CREATE TYPE goal_category AS ENUM ('health', 'career', 'education', 'personal', 'financial', 'relationship');
CREATE TYPE goal_status AS ENUM ('draft', 'active', 'paused', 'completed', 'cancelled');
CREATE TYPE task_status AS ENUM ('pending', 'in_progress', 'completed', 'cancelled');
CREATE TYPE event_status AS ENUM ('tentative', 'confirmed', 'cancelled');
CREATE TYPE priority_level AS ENUM ('low', 'medium', 'high', 'critical');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    profile JSONB DEFAULT '{}',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Goals table
CREATE TABLE goals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category goal_category NOT NULL,
    priority priority_level NOT NULL DEFAULT 'medium',
    status goal_status NOT NULL DEFAULT 'draft',
    progress INTEGER DEFAULT 0 CHECK (progress >= 0 AND progress <= 100),
    deadline TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Milestones table
CREATE TABLE milestones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    goal_id UUID NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    target_date TIMESTAMP WITH TIME ZONE,
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tasks table
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    goal_id UUID NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    priority priority_level NOT NULL DEFAULT 'medium',
    status task_status NOT NULL DEFAULT 'pending',
    estimated_duration INTEGER, -- in minutes
    due_date TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Events table
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    goal_id UUID REFERENCES goals(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    timezone VARCHAR(50) DEFAULT 'UTC',
    recurrence JSONB, -- RRULE as JSON
    location TEXT,
    attendees JSONB DEFAULT '[]',
    status event_status NOT NULL DEFAULT 'confirmed',
    external_id VARCHAR(255), -- For Google Calendar sync
    external_source VARCHAR(50), -- 'google', 'outlook', etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Moods table
CREATE TABLE moods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    level INTEGER NOT NULL CHECK (level >= 1 AND level <= 5),
    notes TEXT,
    tags JSONB DEFAULT '[]',
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, date)
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_goals_user_id ON goals(user_id);
CREATE INDEX idx_goals_status ON goals(status);
CREATE INDEX idx_goals_deadline ON goals(deadline);
CREATE INDEX idx_milestones_goal_id ON milestones(goal_id);
CREATE INDEX idx_tasks_goal_id ON tasks(goal_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_goal_id ON events(goal_id);
CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_external_id ON events(external_id);
CREATE INDEX idx_moods_user_id ON moods(user_id);
CREATE INDEX idx_moods_date ON moods(date);

-- Google integrations table
CREATE TABLE google_integrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    google_user_id VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    token_type VARCHAR(50) DEFAULT 'Bearer',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    scopes JSONB DEFAULT '[]',
    calendar_id VARCHAR(255),
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id)
);

-- Google calendar sync configurations
CREATE TABLE google_calendar_syncs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    google_integration_id UUID NOT NULL REFERENCES google_integrations(id) ON DELETE CASCADE,
    calendar_id VARCHAR(255) NOT NULL,
    calendar_name VARCHAR(255) NOT NULL,
    sync_direction VARCHAR(20) NOT NULL CHECK (sync_direction IN ('bidirectional', 'from_google', 'to_google')),
    sync_status VARCHAR(20) NOT NULL CHECK (sync_status IN ('active', 'paused', 'error', 'disabled')),
    last_sync_at TIMESTAMP WITH TIME ZONE,
    last_sync_error TEXT,
    sync_token VARCHAR(1024),
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, calendar_id)
);

-- Indexes for Google integrations
CREATE INDEX idx_google_integrations_user_id ON google_integrations(user_id);
CREATE INDEX idx_google_integrations_google_user_id ON google_integrations(google_user_id);
CREATE INDEX idx_google_integrations_expires_at ON google_integrations(expires_at);
CREATE INDEX idx_google_integrations_enabled ON google_integrations(enabled);

-- Indexes for Google calendar syncs
CREATE INDEX idx_google_calendar_syncs_user_id ON google_calendar_syncs(user_id);
CREATE INDEX idx_google_calendar_syncs_integration_id ON google_calendar_syncs(google_integration_id);
CREATE INDEX idx_google_calendar_syncs_calendar_id ON google_calendar_syncs(calendar_id);
CREATE INDEX idx_google_calendar_syncs_status ON google_calendar_syncs(sync_status);
CREATE INDEX idx_google_calendar_syncs_last_sync_at ON google_calendar_syncs(last_sync_at);

-- Update triggers for updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_goals_updated_at BEFORE UPDATE ON goals FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON tasks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_google_integrations_updated_at BEFORE UPDATE ON google_integrations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_google_calendar_syncs_updated_at BEFORE UPDATE ON google_calendar_syncs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();