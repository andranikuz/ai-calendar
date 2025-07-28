-- Create sync_conflicts table for handling Google Calendar synchronization conflicts

CREATE TABLE sync_conflicts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    calendar_sync_id VARCHAR(255) NOT NULL,
    conflict_type VARCHAR(50) NOT NULL CHECK (conflict_type IN ('time_overlap', 'content_diff', 'duplicate_event', 'deleted_event')),
    local_event JSONB,
    google_event JSONB,
    description TEXT NOT NULL,
    resolution TEXT,
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolved_by VARCHAR(255), -- "auto" or user ID
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'resolved', 'ignored')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_sync_conflicts_user_id ON sync_conflicts(user_id);
CREATE INDEX idx_sync_conflicts_calendar_sync_id ON sync_conflicts(calendar_sync_id);
CREATE INDEX idx_sync_conflicts_status ON sync_conflicts(status);
CREATE INDEX idx_sync_conflicts_conflict_type ON sync_conflicts(conflict_type);
CREATE INDEX idx_sync_conflicts_created_at ON sync_conflicts(created_at);
CREATE INDEX idx_sync_conflicts_user_status ON sync_conflicts(user_id, status);

-- Create partial index for pending conflicts (most common query)
CREATE INDEX idx_sync_conflicts_pending ON sync_conflicts(user_id, created_at) WHERE status = 'pending';

-- Add updated_at trigger
CREATE TRIGGER update_sync_conflicts_updated_at
    BEFORE UPDATE ON sync_conflicts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE sync_conflicts IS 'Stores conflicts that occur during Google Calendar synchronization';
COMMENT ON COLUMN sync_conflicts.conflict_type IS 'Type of conflict: time_overlap, content_diff, duplicate_event, deleted_event';
COMMENT ON COLUMN sync_conflicts.local_event IS 'Local event data as JSON';
COMMENT ON COLUMN sync_conflicts.google_event IS 'Google Calendar event data as JSON';
COMMENT ON COLUMN sync_conflicts.resolution IS 'How the conflict was resolved';
COMMENT ON COLUMN sync_conflicts.resolved_by IS 'Who resolved the conflict: auto or user ID';
COMMENT ON COLUMN sync_conflicts.status IS 'Current status: pending, resolved, ignored';