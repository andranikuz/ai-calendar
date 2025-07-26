-- Migration 004: Create events table
-- This creates the calendar events table

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

-- Indexes for performance
CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_goal_id ON events(goal_id);
CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_external_id ON events(external_id);

-- Update trigger
CREATE TRIGGER update_events_updated_at 
    BEFORE UPDATE ON events 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();