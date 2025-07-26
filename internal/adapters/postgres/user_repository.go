package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repositories.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	profileJSON, err := json.Marshal(user.Profile)
	if err != nil {
		return err
	}
	
	settingsJSON, err := json.Marshal(user.Settings)
	if err != nil {
		return err
	}
	
	query := `
		INSERT INTO users (id, email, name, profile, settings, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err = r.db.Exec(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		profileJSON,
		settingsJSON,
		user.CreatedAt,
		user.UpdatedAt,
	)
	
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	query := `
		SELECT id, email, name, profile, settings, created_at, updated_at
		FROM users
		WHERE id = $1`
	
	row := r.db.QueryRow(ctx, query, id)
	
	var user entities.User
	var profileJSON, settingsJSON []byte
	
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&profileJSON,
		&settingsJSON,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	if err := json.Unmarshal(profileJSON, &user.Profile); err != nil {
		return nil, err
	}
	
	if err := json.Unmarshal(settingsJSON, &user.Settings); err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, email, name, profile, settings, created_at, updated_at
		FROM users
		WHERE email = $1`
	
	row := r.db.QueryRow(ctx, query, email)
	
	var user entities.User
	var profileJSON, settingsJSON []byte
	
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&profileJSON,
		&settingsJSON,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	if err := json.Unmarshal(profileJSON, &user.Profile); err != nil {
		return nil, err
	}
	
	if err := json.Unmarshal(settingsJSON, &user.Settings); err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	profileJSON, err := json.Marshal(user.Profile)
	if err != nil {
		return err
	}
	
	settingsJSON, err := json.Marshal(user.Settings)
	if err != nil {
		return err
	}
	
	query := `
		UPDATE users
		SET email = $2, name = $3, profile = $4, settings = $5, updated_at = $6
		WHERE id = $1`
	
	user.UpdatedAt = time.Now()
	
	result, err := r.db.Exec(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		profileJSON,
		settingsJSON,
		user.UpdatedAt,
	)
	
	if err != nil {
		return err
	}
	
	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id entities.UserID) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	
	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	
	return exists, err
}