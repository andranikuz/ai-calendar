-- Migration 006: Create Google integrations tables
-- This creates tables for Google Calendar integration

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

-- Update triggers
CREATE TRIGGER update_google_integrations_updated_at 
    BEFORE UPDATE ON google_integrations 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_google_calendar_syncs_updated_at 
    BEFORE UPDATE ON google_calendar_syncs 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();