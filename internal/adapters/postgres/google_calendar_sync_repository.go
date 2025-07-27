package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
)

type googleCalendarSyncRepository struct {
	db *pgxpool.Pool
}

func NewGoogleCalendarSyncRepository(db *pgxpool.Pool) repositories.GoogleCalendarSyncRepository {
	return &googleCalendarSyncRepository{db: db}
}

func (r *googleCalendarSyncRepository) Create(ctx context.Context, sync *entities.GoogleCalendarSync) error {
	settingsJSON, err := json.Marshal(sync.Settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	query := `
		INSERT INTO google_calendar_syncs (
			id, user_id, google_integration_id, calendar_id, calendar_name,
			sync_direction, sync_status, last_sync_at, last_sync_error,
			sync_token, settings, webhook_channel_id, webhook_url, 
			webhook_resource_id, webhook_expires_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	_, err = r.db.Exec(ctx, query,
		sync.ID,
		sync.UserID,
		sync.GoogleIntegrationID,
		sync.CalendarID,
		sync.CalendarName,
		sync.SyncDirection,
		sync.SyncStatus,
		sync.LastSyncAt,
		sync.LastSyncError,
		sync.SyncToken,
		settingsJSON,
		sync.WebhookChannelID,
		sync.WebhookURL,
		sync.WebhookResourceID,
		sync.WebhookExpiresAt,
		sync.CreatedAt,
		sync.UpdatedAt,
	)

	return err
}

func (r *googleCalendarSyncRepository) GetByID(ctx context.Context, id string) (*entities.GoogleCalendarSync, error) {
	query := `
		SELECT id, user_id, google_integration_id, calendar_id, calendar_name,
			   sync_direction, sync_status, last_sync_at, last_sync_error,
			   sync_token, settings, webhook_channel_id, webhook_url,
			   webhook_resource_id, webhook_expires_at, created_at, updated_at
		FROM google_calendar_syncs
		WHERE id = $1`

	return r.scanCalendarSync(r.db.QueryRow(ctx, query, id))
}

func (r *googleCalendarSyncRepository) GetByUserID(ctx context.Context, userID entities.UserID) ([]*entities.GoogleCalendarSync, error) {
	query := `
		SELECT id, user_id, google_integration_id, calendar_id, calendar_name,
			   sync_direction, sync_status, last_sync_at, last_sync_error,
			   sync_token, settings, webhook_channel_id, webhook_url,
			   webhook_resource_id, webhook_expires_at, created_at, updated_at
		FROM google_calendar_syncs
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCalendarSyncs(rows)
}

func (r *googleCalendarSyncRepository) GetByIntegrationID(ctx context.Context, integrationID entities.GoogleIntegrationID) ([]*entities.GoogleCalendarSync, error) {
	query := `
		SELECT id, user_id, google_integration_id, calendar_id, calendar_name,
			   sync_direction, sync_status, last_sync_at, last_sync_error,
			   sync_token, settings, webhook_channel_id, webhook_url,
			   webhook_resource_id, webhook_expires_at, created_at, updated_at
		FROM google_calendar_syncs
		WHERE google_integration_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, integrationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCalendarSyncs(rows)
}

func (r *googleCalendarSyncRepository) GetByCalendarID(ctx context.Context, calendarID string) (*entities.GoogleCalendarSync, error) {
	query := `
		SELECT id, user_id, google_integration_id, calendar_id, calendar_name,
			   sync_direction, sync_status, last_sync_at, last_sync_error,
			   sync_token, settings, webhook_channel_id, webhook_url,
			   webhook_resource_id, webhook_expires_at, created_at, updated_at
		FROM google_calendar_syncs
		WHERE calendar_id = $1`

	return r.scanCalendarSync(r.db.QueryRow(ctx, query, calendarID))
}

func (r *googleCalendarSyncRepository) GetByChannelID(ctx context.Context, channelID string) (*entities.GoogleCalendarSync, error) {
	query := `
		SELECT id, user_id, google_integration_id, calendar_id, calendar_name,
			   sync_direction, sync_status, last_sync_at, last_sync_error,
			   sync_token, settings, webhook_channel_id, webhook_url,
			   webhook_resource_id, webhook_expires_at, created_at, updated_at
		FROM google_calendar_syncs
		WHERE webhook_channel_id = $1`

	return r.scanCalendarSync(r.db.QueryRow(ctx, query, channelID))
}

func (r *googleCalendarSyncRepository) Update(ctx context.Context, sync *entities.GoogleCalendarSync) error {
	settingsJSON, err := json.Marshal(sync.Settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	query := `
		UPDATE google_calendar_syncs
		SET calendar_name = $2, sync_direction = $3, sync_status = $4,
			last_sync_at = $5, last_sync_error = $6, sync_token = $7,
			settings = $8, webhook_channel_id = $9, webhook_url = $10,
			webhook_resource_id = $11, webhook_expires_at = $12, updated_at = $13
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query,
		sync.ID,
		sync.CalendarName,
		sync.SyncDirection,
		sync.SyncStatus,
		sync.LastSyncAt,
		sync.LastSyncError,
		sync.SyncToken,
		settingsJSON,
		sync.WebhookChannelID,
		sync.WebhookURL,
		sync.WebhookResourceID,
		sync.WebhookExpiresAt,
		time.Now(),
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *googleCalendarSyncRepository) UpdateSyncStatus(ctx context.Context, id string, status entities.CalendarSyncStatus, lastSyncAt *time.Time, syncError string) error {
	query := `
		UPDATE google_calendar_syncs
		SET sync_status = $2, last_sync_at = $3, last_sync_error = $4, updated_at = $5
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id, status, lastSyncAt, syncError, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *googleCalendarSyncRepository) UpdateSyncToken(ctx context.Context, id string, syncToken string) error {
	query := `
		UPDATE google_calendar_syncs
		SET sync_token = $2, updated_at = $3
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id, syncToken, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *googleCalendarSyncRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM google_calendar_syncs WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *googleCalendarSyncRepository) GetNeedingSync(ctx context.Context) ([]*entities.GoogleCalendarSync, error) {
	query := `
		SELECT id, user_id, google_integration_id, calendar_id, calendar_name,
			   sync_direction, sync_status, last_sync_at, last_sync_error,
			   sync_token, settings, webhook_channel_id, webhook_url,
			   webhook_resource_id, webhook_expires_at, created_at, updated_at
		FROM google_calendar_syncs
		WHERE sync_status = 'active'
		  AND (last_sync_at IS NULL OR 
			   last_sync_at <= NOW() - INTERVAL '15 minutes')
		ORDER BY last_sync_at ASC NULLS FIRST`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCalendarSyncs(rows)
}

func (r *googleCalendarSyncRepository) GetActive(ctx context.Context) ([]*entities.GoogleCalendarSync, error) {
	query := `
		SELECT id, user_id, google_integration_id, calendar_id, calendar_name,
			   sync_direction, sync_status, last_sync_at, last_sync_error,
			   sync_token, settings, webhook_channel_id, webhook_url,
			   webhook_resource_id, webhook_expires_at, created_at, updated_at
		FROM google_calendar_syncs
		WHERE sync_status = 'active'
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCalendarSyncs(rows)
}

func (r *googleCalendarSyncRepository) scanCalendarSync(row interface{}) (*entities.GoogleCalendarSync, error) {
	type rowScanner interface {
		Scan(dest ...interface{}) error
	}

	scanner, ok := row.(rowScanner)
	if !ok {
		return nil, fmt.Errorf("invalid row type")
	}

	var sync entities.GoogleCalendarSync
	var settingsJSON []byte

	err := scanner.Scan(
		&sync.ID,
		&sync.UserID,
		&sync.GoogleIntegrationID,
		&sync.CalendarID,
		&sync.CalendarName,
		&sync.SyncDirection,
		&sync.SyncStatus,
		&sync.LastSyncAt,
		&sync.LastSyncError,
		&sync.SyncToken,
		&settingsJSON,
		&sync.WebhookChannelID,
		&sync.WebhookURL,
		&sync.WebhookResourceID,
		&sync.WebhookExpiresAt,
		&sync.CreatedAt,
		&sync.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal settings
	if err := json.Unmarshal(settingsJSON, &sync.Settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	return &sync, nil
}

func (r *googleCalendarSyncRepository) scanCalendarSyncs(rows interface{}) ([]*entities.GoogleCalendarSync, error) {
	type rowScanner interface {
		Next() bool
		Scan(dest ...interface{}) error
	}

	scanner, ok := rows.(rowScanner)
	if !ok {
		return nil, fmt.Errorf("invalid rows type")
	}

	var syncs []*entities.GoogleCalendarSync
	for scanner.Next() {
		var sync entities.GoogleCalendarSync
		var settingsJSON []byte

		err := scanner.Scan(
			&sync.ID,
			&sync.UserID,
			&sync.GoogleIntegrationID,
			&sync.CalendarID,
			&sync.CalendarName,
			&sync.SyncDirection,
			&sync.SyncStatus,
			&sync.LastSyncAt,
			&sync.LastSyncError,
			&sync.SyncToken,
			&settingsJSON,
			&sync.WebhookChannelID,
			&sync.WebhookURL,
			&sync.WebhookResourceID,
			&sync.WebhookExpiresAt,
			&sync.CreatedAt,
			&sync.UpdatedAt,
		)

		if err != nil {
			continue
		}

		// Unmarshal settings
		if err := json.Unmarshal(settingsJSON, &sync.Settings); err != nil {
			continue
		}

		syncs = append(syncs, &sync)
	}

	return syncs, nil
}