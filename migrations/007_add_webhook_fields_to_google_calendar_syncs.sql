-- Add webhook fields to google_calendar_syncs table
-- Migration: 007_add_webhook_fields_to_google_calendar_syncs.sql

ALTER TABLE google_calendar_syncs 
ADD COLUMN webhook_channel_id VARCHAR(255),
ADD COLUMN webhook_url VARCHAR(512),
ADD COLUMN webhook_resource_id VARCHAR(255),
ADD COLUMN webhook_expires_at TIMESTAMP;

-- Create index for webhook channel lookup
CREATE INDEX idx_google_calendar_syncs_webhook_channel_id 
ON google_calendar_syncs(webhook_channel_id) 
WHERE webhook_channel_id IS NOT NULL;

-- Update last_sync_at column to be NOT NULL with default
ALTER TABLE google_calendar_syncs 
ALTER COLUMN last_sync_at SET DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN last_sync_at SET NOT NULL;