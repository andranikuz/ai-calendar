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

type googleIntegrationRepository struct {
	db *pgxpool.Pool
}

func NewGoogleIntegrationRepository(db *pgxpool.Pool) repositories.GoogleIntegrationRepository {
	return &googleIntegrationRepository{db: db}
}

func (r *googleIntegrationRepository) Create(ctx context.Context, integration *entities.GoogleIntegration) error {
	scopesJSON, err := json.Marshal(integration.Scopes)
	if err != nil {
		return fmt.Errorf("failed to marshal scopes: %w", err)
	}

	query := `
		INSERT INTO google_integrations (
			id, user_id, google_user_id, email, name, access_token, 
			refresh_token, token_type, expires_at, scopes, calendar_id, 
			enabled, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err = r.db.Exec(ctx, query,
		integration.ID,
		integration.UserID,
		integration.GoogleUserID,
		integration.Email,
		integration.Name,
		integration.AccessToken,
		integration.RefreshToken,
		integration.TokenType,
		integration.ExpiresAt,
		scopesJSON,
		integration.CalendarID,
		integration.Enabled,
		integration.CreatedAt,
		integration.UpdatedAt,
	)

	return err
}

func (r *googleIntegrationRepository) GetByID(ctx context.Context, id entities.GoogleIntegrationID) (*entities.GoogleIntegration, error) {
	query := `
		SELECT id, user_id, google_user_id, email, name, access_token, 
			   refresh_token, token_type, expires_at, scopes, calendar_id, 
			   enabled, created_at, updated_at
		FROM google_integrations
		WHERE id = $1`

	return r.scanIntegration(r.db.QueryRow(ctx, query, id))
}

func (r *googleIntegrationRepository) GetByUserID(ctx context.Context, userID entities.UserID) (*entities.GoogleIntegration, error) {
	query := `
		SELECT id, user_id, google_user_id, email, name, access_token, 
			   refresh_token, token_type, expires_at, scopes, calendar_id, 
			   enabled, created_at, updated_at
		FROM google_integrations
		WHERE user_id = $1`

	return r.scanIntegration(r.db.QueryRow(ctx, query, userID))
}

func (r *googleIntegrationRepository) GetByGoogleUserID(ctx context.Context, googleUserID string) (*entities.GoogleIntegration, error) {
	query := `
		SELECT id, user_id, google_user_id, email, name, access_token, 
			   refresh_token, token_type, expires_at, scopes, calendar_id, 
			   enabled, created_at, updated_at
		FROM google_integrations
		WHERE google_user_id = $1`

	return r.scanIntegration(r.db.QueryRow(ctx, query, googleUserID))
}

func (r *googleIntegrationRepository) Update(ctx context.Context, integration *entities.GoogleIntegration) error {
	scopesJSON, err := json.Marshal(integration.Scopes)
	if err != nil {
		return fmt.Errorf("failed to marshal scopes: %w", err)
	}

	query := `
		UPDATE google_integrations
		SET google_user_id = $2, email = $3, name = $4, access_token = $5,
			refresh_token = $6, token_type = $7, expires_at = $8, scopes = $9,
			calendar_id = $10, enabled = $11, updated_at = $12
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query,
		integration.ID,
		integration.GoogleUserID,
		integration.Email,
		integration.Name,
		integration.AccessToken,
		integration.RefreshToken,
		integration.TokenType,
		integration.ExpiresAt,
		scopesJSON,
		integration.CalendarID,
		integration.Enabled,
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

func (r *googleIntegrationRepository) UpdateTokens(ctx context.Context, id entities.GoogleIntegrationID, accessToken, refreshToken string, expiresAt time.Time) error {
	query := `
		UPDATE google_integrations
		SET access_token = $2, refresh_token = $3, expires_at = $4, updated_at = $5
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id, accessToken, refreshToken, expiresAt, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *googleIntegrationRepository) Delete(ctx context.Context, id entities.GoogleIntegrationID) error {
	query := `DELETE FROM google_integrations WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *googleIntegrationRepository) GetExpiringSoon(ctx context.Context, beforeTime time.Time) ([]*entities.GoogleIntegration, error) {
	query := `
		SELECT id, user_id, google_user_id, email, name, access_token, 
			   refresh_token, token_type, expires_at, scopes, calendar_id, 
			   enabled, created_at, updated_at
		FROM google_integrations
		WHERE expires_at <= $1 AND enabled = true
		ORDER BY expires_at ASC`

	rows, err := r.db.Query(ctx, query, beforeTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanIntegrations(rows)
}

func (r *googleIntegrationRepository) GetActive(ctx context.Context) ([]*entities.GoogleIntegration, error) {
	query := `
		SELECT id, user_id, google_user_id, email, name, access_token, 
			   refresh_token, token_type, expires_at, scopes, calendar_id, 
			   enabled, created_at, updated_at
		FROM google_integrations
		WHERE enabled = true
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanIntegrations(rows)
}

func (r *googleIntegrationRepository) scanIntegration(row interface{}) (*entities.GoogleIntegration, error) {
	type rowScanner interface {
		Scan(dest ...interface{}) error
	}

	scanner, ok := row.(rowScanner)
	if !ok {
		return nil, fmt.Errorf("invalid row type")
	}

	var integration entities.GoogleIntegration
	var scopesJSON []byte

	err := scanner.Scan(
		&integration.ID,
		&integration.UserID,
		&integration.GoogleUserID,
		&integration.Email,
		&integration.Name,
		&integration.AccessToken,
		&integration.RefreshToken,
		&integration.TokenType,
		&integration.ExpiresAt,
		&scopesJSON,
		&integration.CalendarID,
		&integration.Enabled,
		&integration.CreatedAt,
		&integration.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal scopes
	if err := json.Unmarshal(scopesJSON, &integration.Scopes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scopes: %w", err)
	}

	return &integration, nil
}

func (r *googleIntegrationRepository) scanIntegrations(rows interface{}) ([]*entities.GoogleIntegration, error) {
	type rowScanner interface {
		Next() bool
		Scan(dest ...interface{}) error
	}

	scanner, ok := rows.(rowScanner)
	if !ok {
		return nil, fmt.Errorf("invalid rows type")
	}

	var integrations []*entities.GoogleIntegration
	for scanner.Next() {
		var integration entities.GoogleIntegration
		var scopesJSON []byte

		err := scanner.Scan(
			&integration.ID,
			&integration.UserID,
			&integration.GoogleUserID,
			&integration.Email,
			&integration.Name,
			&integration.AccessToken,
			&integration.RefreshToken,
			&integration.TokenType,
			&integration.ExpiresAt,
			&scopesJSON,
			&integration.CalendarID,
			&integration.Enabled,
			&integration.CreatedAt,
			&integration.UpdatedAt,
		)

		if err != nil {
			continue
		}

		// Unmarshal scopes
		if err := json.Unmarshal(scopesJSON, &integration.Scopes); err != nil {
			continue
		}

		integrations = append(integrations, &integration)
	}

	return integrations, nil
}