-- Migration 005: Create moods table
-- This creates the mood tracking table

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
CREATE INDEX idx_moods_user_id ON moods(user_id);
CREATE INDEX idx_moods_date ON moods(date);