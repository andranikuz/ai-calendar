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

type moodRepository struct {
	db *pgxpool.Pool
}

func NewMoodRepository(db *pgxpool.Pool) repositories.MoodRepository {
	return &moodRepository{db: db}
}

func (r *moodRepository) Create(ctx context.Context, mood *entities.Mood) error {
	query := `
		INSERT INTO moods (id, user_id, date, level, notes, tags, recorded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	tagsJSON, err := json.Marshal(mood.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	_, err = r.db.Exec(ctx, query,
		mood.ID,
		mood.UserID,
		mood.Date.Format("2006-01-02"),
		int(mood.Level),
		mood.Notes,
		tagsJSON,
		mood.RecordedAt,
	)

	return err
}

func (r *moodRepository) GetByID(ctx context.Context, id entities.MoodID) (*entities.Mood, error) {
	query := `
		SELECT id, user_id, date, level, notes, tags, recorded_at
		FROM moods
		WHERE id = $1`

	var mood entities.Mood
	var tagsJSON []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&mood.ID,
		&mood.UserID,
		&mood.Date,
		&mood.Level,
		&mood.Notes,
		&tagsJSON,
		&mood.RecordedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal tags
	if err := json.Unmarshal(tagsJSON, &mood.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return &mood, nil
}

func (r *moodRepository) GetByUserIDAndDate(ctx context.Context, userID entities.UserID, date time.Time) (*entities.Mood, error) {
	query := `
		SELECT id, user_id, date, level, notes, tags, recorded_at
		FROM moods
		WHERE user_id = $1 AND date = $2`

	var mood entities.Mood
	var tagsJSON []byte

	err := r.db.QueryRow(ctx, query, userID, date.Format("2006-01-02")).Scan(
		&mood.ID,
		&mood.UserID,
		&mood.Date,
		&mood.Level,
		&mood.Notes,
		&tagsJSON,
		&mood.RecordedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal tags
	if err := json.Unmarshal(tagsJSON, &mood.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return &mood, nil
}

func (r *moodRepository) GetByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Mood, error) {
	query := `
		SELECT id, user_id, date, level, notes, tags, recorded_at
		FROM moods
		WHERE user_id = $1
		ORDER BY date DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMoods(rows)
}

func (r *moodRepository) GetByUserIDAndDateRange(ctx context.Context, userID entities.UserID, start, end time.Time) ([]*entities.Mood, error) {
	query := `
		SELECT id, user_id, date, level, notes, tags, recorded_at
		FROM moods
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC`

	rows, err := r.db.Query(ctx, query, userID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMoods(rows)
}

func (r *moodRepository) GetByUserIDAndLevel(ctx context.Context, userID entities.UserID, level entities.MoodLevel) ([]*entities.Mood, error) {
	query := `
		SELECT id, user_id, date, level, notes, tags, recorded_at
		FROM moods
		WHERE user_id = $1 AND level = $2
		ORDER BY date DESC`

	rows, err := r.db.Query(ctx, query, userID, int(level))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMoods(rows)
}

func (r *moodRepository) GetByUserIDAndTags(ctx context.Context, userID entities.UserID, tags []entities.MoodTag) ([]*entities.Mood, error) {
	if len(tags) == 0 {
		return []*entities.Mood{}, nil
	}

	// For JSONB tags, we need to use different query approach
	searchTags := make([]string, len(tags))
	for i, tag := range tags {
		searchTags[i] = string(tag)
	}

	// Use JSONB containment operator
	query := `
		SELECT id, user_id, date, level, notes, tags, recorded_at
		FROM moods
		WHERE user_id = $1 AND (tags ?| $2)
		ORDER BY date DESC`

	searchTagsJSON, err := json.Marshal(searchTags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search tags: %w", err)
	}

	rows, err := r.db.Query(ctx, query, userID, searchTagsJSON)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMoods(rows)
}

func (r *moodRepository) GetLatestByUserID(ctx context.Context, userID entities.UserID) (*entities.Mood, error) {
	query := `
		SELECT id, user_id, date, level, notes, tags, recorded_at
		FROM moods
		WHERE user_id = $1
		ORDER BY date DESC, recorded_at DESC
		LIMIT 1`

	var mood entities.Mood
	var tagsJSON []byte

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&mood.ID,
		&mood.UserID,
		&mood.Date,
		&mood.Level,
		&mood.Notes,
		&tagsJSON,
		&mood.RecordedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal tags
	if err := json.Unmarshal(tagsJSON, &mood.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return &mood, nil
}

func (r *moodRepository) Update(ctx context.Context, mood *entities.Mood) error {
	query := `
		UPDATE moods
		SET level = $2, notes = $3, tags = $4, recorded_at = $5
		WHERE id = $1`

	tagsJSON, err := json.Marshal(mood.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	result, err := r.db.Exec(ctx, query,
		mood.ID,
		int(mood.Level),
		mood.Notes,
		tagsJSON,
		time.Now(),
	)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *moodRepository) Delete(ctx context.Context, id entities.MoodID) error {
	query := `DELETE FROM moods WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *moodRepository) UpsertByDate(ctx context.Context, mood *entities.Mood) error {
	tagsJSON, err := json.Marshal(mood.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		INSERT INTO moods (id, user_id, date, level, notes, tags, recorded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, date)
		DO UPDATE SET
			level = EXCLUDED.level,
			notes = EXCLUDED.notes,
			tags = EXCLUDED.tags,
			recorded_at = EXCLUDED.recorded_at`

	_, err = r.db.Exec(ctx, query,
		mood.ID,
		mood.UserID,
		mood.Date.Format("2006-01-02"),
		int(mood.Level),
		mood.Notes,
		tagsJSON,
		mood.RecordedAt,
	)

	return err
}

func (r *moodRepository) GetStatsByUserID(ctx context.Context, userID entities.UserID, start, end time.Time) (*repositories.MoodStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_entries,
			AVG(level) as average_level,
			unnest(tags) as tag,
			level,
			date
		FROM moods
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		GROUP BY level, date, tag`

	rows, err := r.db.Query(ctx, query, userID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := &repositories.MoodStats{
		UserID:      userID,
		LevelCounts: make(map[entities.MoodLevel]int),
		TagCounts:   make(map[entities.MoodTag]int),
	}

	var bestLevel entities.MoodLevel = 0
	var worstLevel entities.MoodLevel = 6
	var bestDate, worstDate time.Time

	for rows.Next() {
		var totalEntries int
		var avgLevel float64
		var tag sql.NullString
		var level int
		var dateStr string

		err := rows.Scan(&totalEntries, &avgLevel, &tag, &level, &dateStr)
		if err != nil {
			continue
		}

		stats.TotalEntries = totalEntries
		stats.AverageLevel = avgLevel

		moodLevel := entities.MoodLevel(level)
		stats.LevelCounts[moodLevel]++

		if tag.Valid {
			stats.TagCounts[entities.MoodTag(tag.String)]++
		}

		date, _ := time.Parse("2006-01-02", dateStr)
		if moodLevel > bestLevel {
			bestLevel = moodLevel
			bestDate = date
		}
		if moodLevel < worstLevel {
			worstLevel = moodLevel
			worstDate = date
		}
	}

	stats.BestDay = bestDate
	stats.WorstDay = worstDate

	// Find most common tag
	var maxCount int
	for tag, count := range stats.TagCounts {
		if count > maxCount {
			maxCount = count
			stats.MostCommonTag = tag
		}
	}

	return stats, nil
}

func (r *moodRepository) GetTrendsByUserID(ctx context.Context, userID entities.UserID, days int) ([]*repositories.MoodTrend, error) {
	query := `
		SELECT date, level, tags, notes
		FROM moods
		WHERE user_id = $1 AND date >= $2
		ORDER BY date DESC
		LIMIT $3`

	startDate := time.Now().AddDate(0, 0, -days)
	rows, err := r.db.Query(ctx, query, userID, startDate.Format("2006-01-02"), days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []*repositories.MoodTrend
	for rows.Next() {
		var trend repositories.MoodTrend
		var tagsJSON []byte

		err := rows.Scan(&trend.Date, &trend.Level, &tagsJSON, &trend.Notes)
		if err != nil {
			continue
		}

		// Unmarshal tags
		if err := json.Unmarshal(tagsJSON, &trend.Tags); err != nil {
			continue
		}

		trends = append(trends, &trend)
	}

	return trends, nil
}

func (r *moodRepository) ExistsByUserIDAndDate(ctx context.Context, userID entities.UserID, date time.Time) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM moods WHERE user_id = $1 AND date = $2)`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, date.Format("2006-01-02")).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *moodRepository) scanMoods(rows interface{}) ([]*entities.Mood, error) {
	type rowScanner interface {
		Next() bool
		Scan(dest ...interface{}) error
	}

	scanner, ok := rows.(rowScanner)
	if !ok {
		return nil, fmt.Errorf("invalid row type")
	}

	var moods []*entities.Mood
	for scanner.Next() {
		var mood entities.Mood
		var tagsJSON []byte

		err := scanner.Scan(
			&mood.ID,
			&mood.UserID,
			&mood.Date,
			&mood.Level,
			&mood.Notes,
			&tagsJSON,
			&mood.RecordedAt,
		)
		if err != nil {
			continue
		}

		// Unmarshal tags
		if err := json.Unmarshal(tagsJSON, &mood.Tags); err != nil {
			continue
		}

		moods = append(moods, &mood)
	}

	return moods, nil
}