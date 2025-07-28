package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/andranikuz/smart-goal-calendar/internal/domain/repositories"
)

type SyncConflictService struct {
	conflictRepo repositories.SyncConflictRepository
	eventRepo    repositories.EventRepository
}

func NewSyncConflictService(
	conflictRepo repositories.SyncConflictRepository,
	eventRepo repositories.EventRepository,
) *SyncConflictService {
	return &SyncConflictService{
		conflictRepo: conflictRepo,
		eventRepo:    eventRepo,
	}
}

// DetectConflicts analyzes local and Google events to detect synchronization conflicts
func (s *SyncConflictService) DetectConflicts(
	ctx context.Context,
	userID entities.UserID,
	calendarSyncID string,
	localEvents []*entities.Event,
	googleEvents []*entities.Event,
) ([]*entities.SyncConflict, error) {
	var conflicts []*entities.SyncConflict

	// Check for content differences in events with same Google ID
	for _, localEvent := range localEvents {
		if localEvent.GoogleEventID == nil || *localEvent.GoogleEventID == "" {
			continue
		}

		for _, googleEvent := range googleEvents {
			if googleEvent.GoogleEventID != nil && *googleEvent.GoogleEventID == *localEvent.GoogleEventID {
				if conflict := s.detectContentConflict(userID, calendarSyncID, localEvent, googleEvent); conflict != nil {
					conflicts = append(conflicts, conflict)
				}
				break
			}
		}
	}

	// Check for time overlaps
	timeConflicts := s.detectTimeOverlaps(userID, calendarSyncID, localEvents, googleEvents)
	conflicts = append(conflicts, timeConflicts...)

	// Check for duplicate events
	duplicateConflicts := s.detectDuplicateEvents(userID, calendarSyncID, localEvents, googleEvents)
	conflicts = append(conflicts, duplicateConflicts...)

	return conflicts, nil
}

// detectContentConflict checks if the same event has different content in local and Google
func (s *SyncConflictService) detectContentConflict(
	userID entities.UserID,
	calendarSyncID string,
	localEvent, googleEvent *entities.Event,
) *entities.SyncConflict {
	var differences []string

	if localEvent.Title != googleEvent.Title {
		differences = append(differences, "title")
	}
	if localEvent.Description != googleEvent.Description {
		differences = append(differences, "description")
	}
	if localEvent.Location != googleEvent.Location {
		differences = append(differences, "location")
	}
	if !localEvent.StartTime.Equal(googleEvent.StartTime) {
		differences = append(differences, "start_time")
	}
	if !localEvent.EndTime.Equal(googleEvent.EndTime) {
		differences = append(differences, "end_time")
	}

	if len(differences) == 0 {
		return nil // No conflict
	}

	description := fmt.Sprintf("Content differences detected in fields: %s", strings.Join(differences, ", "))

	return &entities.SyncConflict{
		UserID:         userID,
		CalendarSyncID: calendarSyncID,
		ConflictType:   entities.ConflictTypeContentDiff,
		LocalEvent:     localEvent,
		GoogleEvent:    googleEvent,
		Description:    description,
		Status:         "pending",
	}
}

// detectTimeOverlaps checks for time overlapping events
func (s *SyncConflictService) detectTimeOverlaps(
	userID entities.UserID,
	calendarSyncID string,
	localEvents, googleEvents []*entities.Event,
) []*entities.SyncConflict {
	var conflicts []*entities.SyncConflict

	for _, localEvent := range localEvents {
		for _, googleEvent := range googleEvents {
			// Skip if same event
			if localEvent.GoogleEventID != nil && googleEvent.GoogleEventID != nil && 
			   *localEvent.GoogleEventID == *googleEvent.GoogleEventID {
				continue
			}

			// Check for time overlap
			if s.eventsOverlap(localEvent, googleEvent) {
				description := fmt.Sprintf("Events overlap in time: %s - %s", 
					localEvent.StartTime.Format("15:04"), localEvent.EndTime.Format("15:04"))

				conflict := &entities.SyncConflict{
					UserID:         userID,
					CalendarSyncID: calendarSyncID,
					ConflictType:   entities.ConflictTypeTimeOverlap,
					LocalEvent:     localEvent,
					GoogleEvent:    googleEvent,
					Description:    description,
					Status:         "pending",
				}

				conflicts = append(conflicts, conflict)
			}
		}
	}

	return conflicts
}

// detectDuplicateEvents checks for duplicate events (same title, time, but different IDs)
func (s *SyncConflictService) detectDuplicateEvents(
	userID entities.UserID,
	calendarSyncID string,
	localEvents, googleEvents []*entities.Event,
) []*entities.SyncConflict {
	var conflicts []*entities.SyncConflict

	for _, localEvent := range localEvents {
		for _, googleEvent := range googleEvents {
			// Skip if same event
			if localEvent.GoogleEventID != nil && googleEvent.GoogleEventID != nil && 
			   *localEvent.GoogleEventID == *googleEvent.GoogleEventID {
				continue
			}

			// Check if events are duplicates (same title and time)
			if s.areEventsDuplicates(localEvent, googleEvent) {
				description := fmt.Sprintf("Duplicate events detected: '%s' at %s", 
					localEvent.Title, localEvent.StartTime.Format("2006-01-02 15:04"))

				conflict := &entities.SyncConflict{
					UserID:         userID,
					CalendarSyncID: calendarSyncID,
					ConflictType:   entities.ConflictTypeDuplicateEvent,
					LocalEvent:     localEvent,
					GoogleEvent:    googleEvent,
					Description:    description,
					Status:         "pending",
				}

				conflicts = append(conflicts, conflict)
			}
		}
	}

	return conflicts
}

// eventsOverlap checks if two events overlap in time
func (s *SyncConflictService) eventsOverlap(event1, event2 *entities.Event) bool {
	return event1.StartTime.Before(event2.EndTime) && event2.StartTime.Before(event1.EndTime)
}

// areEventsDuplicates checks if two events are likely duplicates
func (s *SyncConflictService) areEventsDuplicates(event1, event2 *entities.Event) bool {
	// Same title and start time (within 5 minutes tolerance)
	titleMatch := strings.EqualFold(strings.TrimSpace(event1.Title), strings.TrimSpace(event2.Title))
	timeTolerance := 5 * time.Minute
	timeMatch := event1.StartTime.Sub(event2.StartTime).Abs() <= timeTolerance

	return titleMatch && timeMatch
}

// ResolveConflict applies a resolution strategy to a conflict
func (s *SyncConflictService) ResolveConflict(
	ctx context.Context,
	conflictID string,
	action *entities.ConflictResolutionAction,
	resolvedBy string,
) error {
	conflict, err := s.conflictRepo.GetByID(ctx, conflictID)
	if err != nil {
		return fmt.Errorf("failed to get conflict: %w", err)
	}

	switch action.Action {
	case "use_local":
		err = s.applyLocalVersion(ctx, conflict)
	case "use_google":
		err = s.applyGoogleVersion(ctx, conflict)
	case "merge":
		err = s.applyMergedVersion(ctx, conflict, action.EventData)
	case "ignore":
		// Just mark as resolved without making changes
		err = nil
	default:
		return fmt.Errorf("unknown resolution action: %s", action.Action)
	}

	if err != nil {
		return fmt.Errorf("failed to apply resolution: %w", err)
	}

	// Mark conflict as resolved
	return s.conflictRepo.UpdateResolution(ctx, conflictID, action.Resolution, resolvedBy, time.Now())
}

// applyLocalVersion keeps the local version of the event
func (s *SyncConflictService) applyLocalVersion(ctx context.Context, conflict *entities.SyncConflict) error {
	if conflict.LocalEvent == nil {
		return fmt.Errorf("no local event to apply")
	}

	// Update the local event to mark it as the preferred version
	// This might involve updating sync metadata or timestamps
	return s.eventRepo.Update(ctx, conflict.LocalEvent)
}

// applyGoogleVersion uses the Google version of the event
func (s *SyncConflictService) applyGoogleVersion(ctx context.Context, conflict *entities.SyncConflict) error {
	if conflict.GoogleEvent == nil {
		return fmt.Errorf("no Google event to apply")
	}

	// Update local event with Google event data
	if conflict.LocalEvent != nil {
		// Preserve local ID but update content
		conflict.GoogleEvent.ID = conflict.LocalEvent.ID
		return s.eventRepo.Update(ctx, conflict.GoogleEvent)
	} else {
		// Create new local event from Google event
		return s.eventRepo.Create(ctx, conflict.GoogleEvent)
	}
}

// applyMergedVersion merges data from both versions
func (s *SyncConflictService) applyMergedVersion(ctx context.Context, conflict *entities.SyncConflict, eventData map[string]interface{}) error {
	if conflict.LocalEvent == nil {
		return fmt.Errorf("no local event to merge")
	}

	// Apply merged data to the local event
	if title, ok := eventData["title"].(string); ok {
		conflict.LocalEvent.Title = title
	}
	if description, ok := eventData["description"].(string); ok {
		conflict.LocalEvent.Description = description
	}
	if location, ok := eventData["location"].(string); ok {
		conflict.LocalEvent.Location = location
	}

	return s.eventRepo.Update(ctx, conflict.LocalEvent)
}

// AutoResolveConflicts automatically resolves conflicts based on the sync settings
func (s *SyncConflictService) AutoResolveConflicts(
	ctx context.Context,
	conflicts []*entities.SyncConflict,
	strategy entities.ConflictResolutionStrategy,
) error {
	for _, conflict := range conflicts {
		var action *entities.ConflictResolutionAction

		switch strategy {
		case entities.ConflictResolutionGoogleWins:
			action = &entities.ConflictResolutionAction{
				Action:     "use_google",
				Resolution: "Automatically resolved: Google Calendar version preferred",
			}
		case entities.ConflictResolutionLocalWins:
			action = &entities.ConflictResolutionAction{
				Action:     "use_local",
				Resolution: "Automatically resolved: Local version preferred",
			}
		case entities.ConflictResolutionManual:
			// Skip auto-resolution for manual strategy
			continue
		}

		if action != nil {
			if err := s.ResolveConflict(ctx, conflict.ID, action, "auto"); err != nil {
				// Log error but continue with other conflicts
				fmt.Printf("Failed to auto-resolve conflict %s: %v\n", conflict.ID, err)
			}
		}
	}

	return nil
}

// GetPendingConflicts returns all pending conflicts for a user
func (s *SyncConflictService) GetPendingConflicts(ctx context.Context, userID entities.UserID) ([]*entities.SyncConflict, error) {
	return s.conflictRepo.GetPendingByUserID(ctx, userID)
}

// GetConflictStats returns conflict statistics for a user
func (s *SyncConflictService) GetConflictStats(ctx context.Context, userID entities.UserID, days int) (map[string]int, error) {
	since := time.Now().AddDate(0, 0, -days)
	return s.conflictRepo.GetConflictStats(ctx, userID, since)
}