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

type EventHandler struct {
	eventRepo repositories.EventRepository
	goalRepo  repositories.GoalRepository
	eventService *services.EventService
}

func NewEventHandler(
	eventRepo repositories.EventRepository,
	goalRepo repositories.GoalRepository,
	eventService *services.EventService,
) *EventHandler {
	return &EventHandler{
		eventRepo:    eventRepo,
		goalRepo:     goalRepo,
		eventService: eventService,
	}
}

// Command Handlers

func (h *EventHandler) HandleCreateEvent(ctx context.Context, cmd commands.CreateEventCommand) (*commands.CreateEventResult, error) {
	// Validate time range
	if cmd.EndTime.Before(cmd.StartTime) || cmd.EndTime.Equal(cmd.StartTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	// Check for conflicts if needed
	hasConflict, err := h.eventRepo.HasConflict(ctx, cmd.UserID, cmd.StartTime, cmd.EndTime, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check for conflicts: %w", err)
	}
	if hasConflict {
		return nil, fmt.Errorf("event conflicts with existing event")
	}

	// Validate goal ownership if GoalID is provided
	if cmd.GoalID != nil {
		goal, err := h.goalRepo.GetByID(ctx, *cmd.GoalID)
		if err != nil {
			return nil, fmt.Errorf("failed to get goal: %w", err)
		}
		if goal == nil {
			return nil, fmt.Errorf("goal not found")
		}
		if goal.UserID != cmd.UserID {
			return nil, fmt.Errorf("access denied: goal belongs to different user")
		}
	}

	// Create event entity
	now := time.Now()
	event := &entities.Event{
		ID:             entities.EventID(uuid.New().String()),
		UserID:         cmd.UserID,
		GoalID:         cmd.GoalID,
		Title:          h.eventService.SanitizeEventTitle(cmd.Title),
		Description:    h.eventService.SanitizeEventDescription(cmd.Description),
		StartTime:      cmd.StartTime,
		EndTime:        cmd.EndTime,
		Timezone:       cmd.Timezone,
		Recurrence:     cmd.Recurrence,
		Location:       cmd.Location,
		Attendees:      cmd.Attendees,
		Status:         cmd.Status,
		ExternalID:     cmd.ExternalID,
		ExternalSource: cmd.ExternalSource,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Validate event
	if err := h.eventService.ValidateEventCreation(event); err != nil {
		return nil, fmt.Errorf("event validation failed: %w", err)
	}

	// Save event
	if err := h.eventRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return &commands.CreateEventResult{
		EventID:   event.ID,
		CreatedAt: event.CreatedAt,
	}, nil
}

func (h *EventHandler) HandleUpdateEvent(ctx context.Context, cmd commands.UpdateEventCommand) (*commands.UpdateEventResult, error) {
	// Get existing event
	event, err := h.eventRepo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	// Check ownership
	if event.UserID != cmd.UserID {
		return nil, fmt.Errorf("access denied: event belongs to different user")
	}

	// Update fields if provided
	if cmd.Title != nil {
		event.Title = h.eventService.SanitizeEventTitle(*cmd.Title)
	}

	if cmd.Description != nil {
		event.Description = h.eventService.SanitizeEventDescription(*cmd.Description)
	}

	if cmd.StartTime != nil {
		event.StartTime = *cmd.StartTime
	}

	if cmd.EndTime != nil {
		event.EndTime = *cmd.EndTime
	}

	// Validate time range after updates
	if event.EndTime.Before(event.StartTime) || event.EndTime.Equal(event.StartTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	// Check for conflicts if time changed
	if cmd.StartTime != nil || cmd.EndTime != nil {
		hasConflict, err := h.eventRepo.HasConflict(ctx, cmd.UserID, event.StartTime, event.EndTime, &cmd.EventID)
		if err != nil {
			return nil, fmt.Errorf("failed to check for conflicts: %w", err)
		}
		if hasConflict {
			return nil, fmt.Errorf("event conflicts with existing event")
		}
	}

	if cmd.Timezone != nil {
		event.Timezone = *cmd.Timezone
	}

	if cmd.Recurrence != nil {
		event.Recurrence = cmd.Recurrence
	}

	if cmd.Location != nil {
		event.Location = *cmd.Location
	}

	if cmd.Attendees != nil {
		event.Attendees = *cmd.Attendees
	}

	if cmd.Status != nil {
		event.Status = *cmd.Status
	}

	if cmd.GoalID != nil {
		// Validate goal ownership
		goal, err := h.goalRepo.GetByID(ctx, *cmd.GoalID)
		if err != nil {
			return nil, fmt.Errorf("failed to get goal: %w", err)
		}
		if goal == nil {
			return nil, fmt.Errorf("goal not found")
		}
		if goal.UserID != cmd.UserID {
			return nil, fmt.Errorf("access denied: goal belongs to different user")
		}
		event.GoalID = cmd.GoalID
	}

	if cmd.ExternalID != nil {
		event.ExternalID = *cmd.ExternalID
	}

	if cmd.ExternalSource != nil {
		event.ExternalSource = *cmd.ExternalSource
	}

	// Update timestamp
	event.UpdatedAt = time.Now()

	// Validate updated event
	if err := h.eventService.ValidateEventCreation(event); err != nil {
		return nil, fmt.Errorf("event validation failed: %w", err)
	}

	// Save updated event
	if err := h.eventRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return &commands.UpdateEventResult{
		UpdatedAt: event.UpdatedAt,
	}, nil
}

func (h *EventHandler) HandleDeleteEvent(ctx context.Context, cmd commands.DeleteEventCommand) error {
	// Get event to check ownership
	event, err := h.eventRepo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	if event == nil {
		return fmt.Errorf("event not found")
	}

	// Check ownership
	if event.UserID != cmd.UserID {
		return fmt.Errorf("access denied: event belongs to different user")
	}

	// Delete event
	if err := h.eventRepo.Delete(ctx, cmd.EventID); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func (h *EventHandler) HandleMoveEvent(ctx context.Context, cmd commands.MoveEventCommand) (*commands.MoveEventResult, error) {
	// Get existing event
	event, err := h.eventRepo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	// Check ownership
	if event.UserID != cmd.UserID {
		return nil, fmt.Errorf("access denied: event belongs to different user")
	}

	// Validate time range
	if cmd.EndTime.Before(cmd.StartTime) || cmd.EndTime.Equal(cmd.StartTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	// Check for conflicts
	hasConflict, err := h.eventRepo.HasConflict(ctx, cmd.UserID, cmd.StartTime, cmd.EndTime, &cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for conflicts: %w", err)
	}
	if hasConflict {
		return nil, fmt.Errorf("event conflicts with existing event")
	}

	// Update times
	event.StartTime = cmd.StartTime
	event.EndTime = cmd.EndTime
	event.UpdatedAt = time.Now()

	// Save updated event
	if err := h.eventRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to move event: %w", err)
	}

	return &commands.MoveEventResult{
		UpdatedAt: event.UpdatedAt,
	}, nil
}

func (h *EventHandler) HandleDuplicateEvent(ctx context.Context, cmd commands.DuplicateEventCommand) (*commands.DuplicateEventResult, error) {
	// Get existing event
	originalEvent, err := h.eventRepo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	if originalEvent == nil {
		return nil, fmt.Errorf("event not found")
	}

	// Check ownership
	if originalEvent.UserID != cmd.UserID {
		return nil, fmt.Errorf("access denied: event belongs to different user")
	}

	// Validate time range
	if cmd.EndTime.Before(cmd.StartTime) || cmd.EndTime.Equal(cmd.StartTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	// Check for conflicts
	hasConflict, err := h.eventRepo.HasConflict(ctx, cmd.UserID, cmd.StartTime, cmd.EndTime, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check for conflicts: %w", err)
	}
	if hasConflict {
		return nil, fmt.Errorf("event conflicts with existing event")
	}

	// Create duplicated event
	now := time.Now()
	newEvent := &entities.Event{
		ID:             entities.EventID(uuid.New().String()),
		UserID:         originalEvent.UserID,
		GoalID:         originalEvent.GoalID,
		Title:          originalEvent.Title + " (Copy)",
		Description:    originalEvent.Description,
		StartTime:      cmd.StartTime,
		EndTime:        cmd.EndTime,
		Timezone:       originalEvent.Timezone,
		Recurrence:     originalEvent.Recurrence,
		Location:       originalEvent.Location,
		Attendees:      originalEvent.Attendees,
		Status:         originalEvent.Status,
		ExternalID:     "", // Clear external references
		ExternalSource: "",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Save new event
	if err := h.eventRepo.Create(ctx, newEvent); err != nil {
		return nil, fmt.Errorf("failed to duplicate event: %w", err)
	}

	return &commands.DuplicateEventResult{
		EventID:   newEvent.ID,
		CreatedAt: newEvent.CreatedAt,
	}, nil
}

func (h *EventHandler) HandleChangeEventStatus(ctx context.Context, cmd commands.ChangeEventStatusCommand) (*commands.ChangeEventStatusResult, error) {
	// Get existing event
	event, err := h.eventRepo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	// Check ownership
	if event.UserID != cmd.UserID {
		return nil, fmt.Errorf("access denied: event belongs to different user")
	}

	// Update status
	event.Status = cmd.Status
	event.UpdatedAt = time.Now()

	// Save updated event
	if err := h.eventRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to update event status: %w", err)
	}

	return &commands.ChangeEventStatusResult{
		UpdatedAt: event.UpdatedAt,
	}, nil
}

func (h *EventHandler) HandleLinkEventToGoal(ctx context.Context, cmd commands.LinkEventToGoalCommand) (*commands.LinkEventToGoalResult, error) {
	// Get existing event
	event, err := h.eventRepo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	// Check event ownership
	if event.UserID != cmd.UserID {
		return nil, fmt.Errorf("access denied: event belongs to different user")
	}

	// Validate goal ownership
	goal, err := h.goalRepo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}
	if goal.UserID != cmd.UserID {
		return nil, fmt.Errorf("access denied: goal belongs to different user")
	}

	// Link event to goal
	event.GoalID = &cmd.GoalID
	event.UpdatedAt = time.Now()

	// Save updated event
	if err := h.eventRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to link event to goal: %w", err)
	}

	return &commands.LinkEventToGoalResult{
		UpdatedAt: event.UpdatedAt,
	}, nil
}

func (h *EventHandler) HandleUnlinkEventFromGoal(ctx context.Context, cmd commands.UnlinkEventFromGoalCommand) (*commands.UnlinkEventFromGoalResult, error) {
	// Get existing event
	event, err := h.eventRepo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	// Check ownership
	if event.UserID != cmd.UserID {
		return nil, fmt.Errorf("access denied: event belongs to different user")
	}

	// Unlink event from goal
	event.GoalID = nil
	event.UpdatedAt = time.Now()

	// Save updated event
	if err := h.eventRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to unlink event from goal: %w", err)
	}

	return &commands.UnlinkEventFromGoalResult{
		UpdatedAt: event.UpdatedAt,
	}, nil
}

// Query Handlers

func (h *EventHandler) HandleGetEventByID(ctx context.Context, query queries.GetEventByIDQuery) (*queries.GetEventResult, error) {
	event, err := h.eventRepo.GetByID(ctx, query.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	// Check ownership
	if event.UserID != query.UserID {
		return nil, fmt.Errorf("access denied: event belongs to different user")
	}

	return &queries.GetEventResult{
		Event: event,
	}, nil
}

func (h *EventHandler) HandleGetEventsByUserID(ctx context.Context, query queries.GetEventsByUserIDQuery) (*queries.GetEventsResult, error) {
	events, totalCount, err := h.eventRepo.GetByUserIDPaginated(ctx, query.UserID, query.Offset, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user events: %w", err)
	}

	return &queries.GetEventsResult{
		Events:     events,
		TotalCount: totalCount,
		Offset:     query.Offset,
		Limit:      query.Limit,
	}, nil
}

func (h *EventHandler) HandleGetEventsByTimeRange(ctx context.Context, query queries.GetEventsByTimeRangeQuery) (*queries.GetEventsByTimeRangeResult, error) {
	events, err := h.eventRepo.GetByUserIDAndTimeRange(ctx, query.UserID, query.StartTime, query.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by time range: %w", err)
	}

	return &queries.GetEventsByTimeRangeResult{
		Events: events,
	}, nil
}

func (h *EventHandler) HandleGetEventsByGoalID(ctx context.Context, query queries.GetEventsByGoalIDQuery) (*queries.GetEventsByGoalIDResult, error) {
	// Check goal ownership
	goal, err := h.goalRepo.GetByID(ctx, query.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}

	if goal.UserID != query.UserID {
		return nil, fmt.Errorf("access denied: goal belongs to different user")
	}

	events, err := h.eventRepo.GetByGoalID(ctx, query.GoalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal events: %w", err)
	}

	return &queries.GetEventsByGoalIDResult{
		Events: events,
	}, nil
}

func (h *EventHandler) HandleGetUpcomingEvents(ctx context.Context, query queries.GetUpcomingEventsQuery) (*queries.GetUpcomingEventsResult, error) {
	events, err := h.eventRepo.GetUpcoming(ctx, query.UserID, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}

	return &queries.GetUpcomingEventsResult{
		Events: events,
	}, nil
}

func (h *EventHandler) HandleGetTodayEvents(ctx context.Context, query queries.GetTodayEventsQuery) (*queries.GetTodayEventsResult, error) {
	events, err := h.eventRepo.GetForToday(ctx, query.UserID, query.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's events: %w", err)
	}

	return &queries.GetTodayEventsResult{
		Events: events,
	}, nil
}

func (h *EventHandler) HandleGetRecurringEvents(ctx context.Context, query queries.GetRecurringEventsQuery) (*queries.GetRecurringEventsResult, error) {
	events, err := h.eventRepo.GetRecurring(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recurring events: %w", err)
	}

	return &queries.GetRecurringEventsResult{
		Events: events,
	}, nil
}

func (h *EventHandler) HandleSearchEvents(ctx context.Context, query queries.SearchEventsQuery) (*queries.SearchEventsResult, error) {
	events, err := h.eventRepo.Search(ctx, query.UserID, query.Query, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search events: %w", err)
	}

	return &queries.SearchEventsResult{
		Events: events,
	}, nil
}

func (h *EventHandler) HandleCheckEventConflict(ctx context.Context, query queries.CheckEventConflictQuery) (*queries.CheckEventConflictResult, error) {
	hasConflict, err := h.eventRepo.HasConflict(ctx, query.UserID, query.StartTime, query.EndTime, query.ExcludeEventID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for conflicts: %w", err)
	}

	return &queries.CheckEventConflictResult{
		HasConflict: hasConflict,
	}, nil
}