package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/andranikuz/smart-goal-calendar/internal/application/commands"
	"github.com/andranikuz/smart-goal-calendar/internal/application/queries"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/services"
)

type MoodHandler struct {
	moodRepo    repositories.MoodRepository
	moodService *services.MoodService
}

func NewMoodHandler(
	moodRepo repositories.MoodRepository,
	moodService *services.MoodService,
) *MoodHandler {
	return &MoodHandler{
		moodRepo:    moodRepo,
		moodService: moodService,
	}
}

// Command handlers
func (h *MoodHandler) CreateMood(ctx context.Context, cmd commands.CreateMoodCommand) (*entities.Mood, error) {
	mood := &entities.Mood{
		ID:         entities.MoodID(uuid.New().String()),
		UserID:     cmd.UserID,
		Date:       cmd.Date,
		Level:      cmd.Level,
		Notes:      cmd.Notes,
		Tags:       cmd.Tags,
		RecordedAt: time.Now(),
	}

	if err := h.moodService.ValidateMood(mood); err != nil {
		return nil, err
	}

	if err := h.moodRepo.Create(ctx, mood); err != nil {
		return nil, err
	}

	return mood, nil
}

func (h *MoodHandler) UpdateMood(ctx context.Context, cmd commands.UpdateMoodCommand) (*entities.Mood, error) {
	existingMood, err := h.moodRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if existingMood == nil {
		return nil, fmt.Errorf("mood not found")
	}

	if existingMood.UserID != cmd.UserID {
		return nil, fmt.Errorf("unauthorized")
	}

	updatedMood := &entities.Mood{
		ID:         cmd.ID,
		UserID:     cmd.UserID,
		Date:       existingMood.Date, // Date cannot be changed via update
		Level:      cmd.Level,
		Notes:      cmd.Notes,
		Tags:       cmd.Tags,
		RecordedAt: existingMood.RecordedAt,
	}

	if err := h.moodService.ValidateMood(updatedMood); err != nil {
		return nil, err
	}

	if err := h.moodRepo.Update(ctx, updatedMood); err != nil {
		return nil, err
	}

	return updatedMood, nil
}

func (h *MoodHandler) DeleteMood(ctx context.Context, cmd commands.DeleteMoodCommand) error {
	existingMood, err := h.moodRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if existingMood == nil {
		return fmt.Errorf("mood not found")
	}

	if existingMood.UserID != cmd.UserID {
		return fmt.Errorf("unauthorized")
	}

	return h.moodRepo.Delete(ctx, cmd.ID)
}

func (h *MoodHandler) UpsertMoodByDate(ctx context.Context, cmd commands.UpsertMoodByDateCommand) (*entities.Mood, error) {
	mood := &entities.Mood{
		ID:         entities.MoodID(uuid.New().String()),
		UserID:     cmd.UserID,
		Date:       cmd.Date,
		Level:      cmd.Level,
		Notes:      cmd.Notes,
		Tags:       cmd.Tags,
		RecordedAt: time.Now(),
	}

	if err := h.moodService.ValidateMood(mood); err != nil {
		return nil, err
	}

	if err := h.moodRepo.UpsertByDate(ctx, mood); err != nil {
		return nil, err
	}

	// Get the actual mood that was inserted/updated
	return h.moodRepo.GetByUserIDAndDate(ctx, cmd.UserID, cmd.Date)
}

// Query handlers
func (h *MoodHandler) GetMoodByID(ctx context.Context, query queries.GetMoodByIDQuery) (*entities.Mood, error) {
	return h.moodRepo.GetByID(ctx, query.ID)
}

func (h *MoodHandler) GetMoodByUserIDAndDate(ctx context.Context, query queries.GetMoodByUserIDAndDateQuery) (*entities.Mood, error) {
	return h.moodRepo.GetByUserIDAndDate(ctx, query.UserID, query.Date)
}

func (h *MoodHandler) GetMoodsByUserID(ctx context.Context, query queries.GetMoodsByUserIDQuery) ([]*entities.Mood, error) {
	allMoods, err := h.moodRepo.GetByUserID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	// Apply pagination
	start := query.Offset
	if start >= len(allMoods) {
		return []*entities.Mood{}, nil
	}

	end := start + query.Limit
	if end > len(allMoods) {
		end = len(allMoods)
	}

	return allMoods[start:end], nil
}

func (h *MoodHandler) GetMoodsByDateRange(ctx context.Context, query queries.GetMoodsByDateRangeQuery) ([]*entities.Mood, error) {
	return h.moodRepo.GetByUserIDAndDateRange(ctx, query.UserID, query.Start, query.End)
}

func (h *MoodHandler) GetMoodsByLevel(ctx context.Context, query queries.GetMoodsByLevelQuery) ([]*entities.Mood, error) {
	return h.moodRepo.GetByUserIDAndLevel(ctx, query.UserID, query.Level)
}

func (h *MoodHandler) GetMoodsByTags(ctx context.Context, query queries.GetMoodsByTagsQuery) ([]*entities.Mood, error) {
	return h.moodRepo.GetByUserIDAndTags(ctx, query.UserID, query.Tags)
}

func (h *MoodHandler) GetLatestMood(ctx context.Context, query queries.GetLatestMoodQuery) (*entities.Mood, error) {
	return h.moodRepo.GetLatestByUserID(ctx, query.UserID)
}

func (h *MoodHandler) GetMoodStats(ctx context.Context, query queries.GetMoodStatsQuery) (*repositories.MoodStats, error) {
	return h.moodRepo.GetStatsByUserID(ctx, query.UserID, query.Start, query.End)
}

func (h *MoodHandler) GetMoodTrends(ctx context.Context, query queries.GetMoodTrendsQuery) ([]*repositories.MoodTrend, error) {
	return h.moodRepo.GetTrendsByUserID(ctx, query.UserID, query.Days)
}

func (h *MoodHandler) CheckMoodExists(ctx context.Context, query queries.CheckMoodExistsQuery) (bool, error) {
	return h.moodRepo.ExistsByUserIDAndDate(ctx, query.UserID, query.Date)
}