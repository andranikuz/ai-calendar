# Google Calendar Integration Guide

## Overview

Smart Goal Calendar supports two-way synchronization with Google Calendar, allowing you to:
- Import events from Google Calendar
- Export local events to Google Calendar
- Keep both calendars in sync automatically

## Setup Instructions

### 1. Configure Google OAuth2 Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google Calendar API:
   - Navigate to "APIs & Services" > "Library"
   - Search for "Google Calendar API"
   - Click "Enable"
4. Create OAuth2 credentials:
   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth client ID"
   - Choose "Web application"
   - Add authorized redirect URI: `http://localhost:8080/api/v1/google/callback`
   - Save the Client ID and Client Secret

### 2. Configure the Application

Update your `config/config.yaml` file with Google credentials:

```yaml
google:
  client_id: "YOUR_GOOGLE_CLIENT_ID"
  client_secret: "YOUR_GOOGLE_CLIENT_SECRET"
  redirect_url: "http://localhost:8080/api/v1/google/callback"
```

Or set environment variables:
```bash
export SMART_CALENDAR_GOOGLE_CLIENT_ID="YOUR_GOOGLE_CLIENT_ID"
export SMART_CALENDAR_GOOGLE_CLIENT_SECRET="YOUR_GOOGLE_CLIENT_SECRET"
export SMART_CALENDAR_GOOGLE_REDIRECT_URL="http://localhost:8080/api/v1/google/callback"
```

### 3. Database Setup

Ensure the Google integration tables are created by running the SQL migration:
```bash
psql -U postgres -d smart_calendar -f scripts/init-db.sql
```

## API Usage

### 1. Connect Google Account

#### Step 1: Get Authorization URL
```bash
curl -X GET http://localhost:8080/api/v1/google/auth-url \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:
```json
{
  "auth_url": "https://accounts.google.com/o/oauth2/v2/auth?...",
  "state": "unique-state-parameter"
}
```

#### Step 2: Redirect User
Direct the user to the `auth_url` to authorize access to their Google Calendar.

#### Step 3: Handle Callback
After authorization, Google redirects back with a code. Send it to complete the integration:

```bash
curl -X POST http://localhost:8080/api/v1/google/callback \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "AUTHORIZATION_CODE_FROM_GOOGLE",
    "state": "STATE_FROM_STEP_1"
  }'
```

### 2. List Available Calendars

```bash
curl -X GET http://localhost:8080/api/v1/google/calendars \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:
```json
{
  "calendars": [
    {
      "id": "primary",
      "summary": "My Calendar",
      "description": "Primary calendar",
      "primary": true,
      "access_role": "owner"
    }
  ]
}
```

### 3. Configure Calendar Sync

Create a sync configuration for a specific calendar:

```bash
curl -X POST http://localhost:8080/api/v1/google/calendar-syncs \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "calendar_id": "primary",
    "calendar_name": "My Primary Calendar",
    "sync_direction": "bidirectional",
    "settings": {
      "sync_interval": 900000000000,
      "auto_sync": true,
      "sync_past_events": false,
      "sync_future_events": true,
      "conflict_resolution": "google_wins"
    }
  }'
```

Sync directions:
- `bidirectional` - Two-way sync between Google and local calendars
- `from_google` - Only import events from Google Calendar
- `to_google` - Only export local events to Google Calendar

### 4. Trigger Manual Sync

```bash
curl -X POST http://localhost:8080/api/v1/google/calendar-syncs/{SYNC_ID}/sync \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:
```json
{
  "message": "Sync completed successfully",
  "synced_count": 15,
  "synced_at": "2025-07-26T12:00:00Z"
}
```

### 5. View Integration Status

```bash
curl -X GET http://localhost:8080/api/v1/google/integration \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 6. Disconnect Google Account

```bash
curl -X DELETE http://localhost:8080/api/v1/google/integration \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Sync Settings

The sync configuration supports the following settings:

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `sync_interval` | Duration | 15 minutes | How often to automatically sync |
| `auto_sync` | Boolean | true | Enable automatic background sync |
| `sync_past_events` | Boolean | false | Include past events in sync |
| `sync_future_events` | Boolean | true | Include future events in sync |
| `conflict_resolution` | String | "google_wins" | How to resolve conflicts ("google_wins", "local_wins", "manual") |

## Event Mapping

When syncing between Google Calendar and Smart Goal Calendar:

| Google Calendar | Smart Goal Calendar |
|-----------------|---------------------|
| Summary | Title |
| Description | Description |
| Location | Location |
| Start | StartTime |
| End | EndTime |
| ID | ExternalID |

## Troubleshooting

### Common Issues

1. **"Failed to refresh Google token"**
   - The refresh token may have expired
   - User needs to re-authorize the application

2. **"Sync completed with errors"**
   - Check the `error_detail` in the response
   - Some events may have failed to sync due to validation errors

3. **"Google integration not found"**
   - User needs to connect their Google account first

### Debug Tips

1. Check integration status:
   ```bash
   curl -X GET http://localhost:8080/api/v1/google/integration \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

2. View sync configurations:
   ```bash
   curl -X GET http://localhost:8080/api/v1/google/calendar-syncs \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

3. Check application logs for detailed error messages

## Security Considerations

1. **Token Storage**: Access and refresh tokens are encrypted in the database
2. **Scope Limitations**: The app only requests calendar read/write permissions
3. **User Consent**: Users must explicitly authorize access to their Google Calendar
4. **Token Refresh**: Tokens are automatically refreshed before expiration

## Future Enhancements

- Real-time sync using Google Calendar webhooks
- Support for recurring events with complex patterns
- Selective sync based on event categories or labels
- Conflict resolution UI for manual intervention
- Support for multiple Google accounts