package services

import (
	"context"
	"testing"
	"time"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSyncConflictRepository is a mock implementation of the SyncConflictRepository
type MockSyncConflictRepository struct {
	mock.Mock
}

func (m *MockSyncConflictRepository) Create(ctx context.Context, conflict *entities.SyncConflict) error {
	args := m.Called(ctx, conflict)
	return args.Error(0)
}

func (m *MockSyncConflictRepository) GetByID(ctx context.Context, id string) (*entities.SyncConflict, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.SyncConflict), args.Error(1)
}

func (m *MockSyncConflictRepository) GetPendingByUserID(ctx context.Context, userID entities.UserID) ([]*entities.SyncConflict, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entities.SyncConflict), args.Error(1)
}

func (m *MockSyncConflictRepository) GetByCalendarSyncID(ctx context.Context, calendarSyncID string) ([]*entities.SyncConflict, error) {
	args := m.Called(ctx, calendarSyncID)
	return args.Get(0).([]*entities.SyncConflict), args.Error(1)
}

func (m *MockSyncConflictRepository) UpdateResolution(ctx context.Context, id string, resolution string, resolvedBy string, resolvedAt time.Time) error {
	args := m.Called(ctx, id, resolution, resolvedBy, resolvedAt)
	return args.Error(0)
}

func (m *MockSyncConflictRepository) MarkAsResolved(ctx context.Context, id string, resolvedBy string) error {
	args := m.Called(ctx, id, resolvedBy)
	return args.Error(0)
}

func (m *MockSyncConflictRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSyncConflictRepository) GetByStatus(ctx context.Context, userID entities.UserID, status string) ([]*entities.SyncConflict, error) {
	args := m.Called(ctx, userID, status)
	return args.Get(0).([]*entities.SyncConflict), args.Error(1)
}

func (m *MockSyncConflictRepository) FindSimilarConflict(ctx context.Context, userID entities.UserID, localEventID, googleEventID string) (*entities.SyncConflict, error) {
	args := m.Called(ctx, userID, localEventID, googleEventID)
	return args.Get(0).(*entities.SyncConflict), args.Error(1)
}

func (m *MockSyncConflictRepository) GetConflictStats(ctx context.Context, userID entities.UserID, since time.Time) (map[string]int, error) {
	args := m.Called(ctx, userID, since)
	return args.Get(0).(map[string]int), args.Error(1)
}

// MockEventRepository is a mock implementation of the EventRepository
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Create(ctx context.Context, event *entities.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) Update(ctx context.Context, event *entities.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) GetByID(ctx context.Context, id entities.EventID) (*entities.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Event), args.Error(1)
}

func (m *MockEventRepository) Delete(ctx context.Context, id entities.EventID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventRepository) GetByUserID(ctx context.Context, userID entities.UserID) ([]*entities.Event, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *MockEventRepository) GetByDateRange(ctx context.Context, userID entities.UserID, start, end time.Time) ([]*entities.Event, error) {
	args := m.Called(ctx, userID, start, end)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *MockEventRepository) BulkCreate(ctx context.Context, events []*entities.Event) error {
	args := m.Called(ctx, events)
	return args.Error(0)
}

// Add other missing methods to satisfy EventRepository interface
func (m *MockEventRepository) GetByUserIDAndTimeRange(ctx context.Context, userID entities.UserID, start, end time.Time) ([]*entities.Event, error) {
	args := m.Called(ctx, userID, start, end)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *MockEventRepository) GetByTimeRange(ctx context.Context, userID entities.UserID, start, end time.Time) ([]*entities.Event, error) {
	args := m.Called(ctx, userID, start, end)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *MockEventRepository) GetByGoalID(ctx context.Context, goalID entities.GoalID) ([]*entities.Event, error) {
	args := m.Called(ctx, goalID)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *MockEventRepository) GetByExternalID(ctx context.Context, userID entities.UserID, externalID string) (*entities.Event, error) {
	args := m.Called(ctx, userID, externalID)
	return args.Get(0).(*entities.Event), args.Error(1)
}

func (m *MockEventRepository) GetByExternalSource(ctx context.Context, userID entities.UserID, source string) ([]*entities.Event, error) {
	args := m.Called(ctx, userID, source)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *MockEventRepository) GetByGoogleEventID(ctx context.Context, googleEventID string) (*entities.Event, error) {
	args := m.Called(ctx, googleEventID)
	return args.Get(0).(*entities.Event), args.Error(1)
}

func (m *MockEventRepository) GetUpcoming(ctx context.Context, userID entities.UserID, limit int) ([]*entities.Event, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *MockEventRepository) GetForToday(ctx context.Context, userID entities.UserID, timezone string) ([]*entities.Event, error) {
	args := m.Called(ctx, userID, timezone)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *MockEventRepository) GetRecurring(ctx context.Context, userID entities.UserID) ([]*entities.Event, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *MockEventRepository) GetByStatus(ctx context.Context, userID entities.UserID, status entities.EventStatus) ([]*entities.Event, error) {
	args := m.Called(ctx, userID, status)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func (m *MockEventRepository) HasConflict(ctx context.Context, userID entities.UserID, start, end time.Time, excludeEventID *entities.EventID) (bool, error) {
	args := m.Called(ctx, userID, start, end, excludeEventID)
	return args.Bool(0), args.Error(1)
}

func (m *MockEventRepository) GetByUserIDPaginated(ctx context.Context, userID entities.UserID, offset, limit int) ([]*entities.Event, int64, error) {
	args := m.Called(ctx, userID, offset, limit)
	return args.Get(0).([]*entities.Event), args.Get(1).(int64), args.Error(2)
}

func (m *MockEventRepository) Search(ctx context.Context, userID entities.UserID, query string, limit int) ([]*entities.Event, error) {
	args := m.Called(ctx, userID, query, limit)
	return args.Get(0).([]*entities.Event), args.Error(1)
}

func TestSyncConflictService_DetectConflicts(t *testing.T) {
	mockConflictRepo := new(MockSyncConflictRepository)
	mockEventRepo := new(MockEventRepository)
	service := NewSyncConflictService(mockConflictRepo, mockEventRepo)
	
	ctx := context.Background()
	userID := entities.UserID("test-user-id")
	calendarSyncID := "test-calendar-sync-id"

	// Test data setup
	now := time.Now()
	
	// Event that exists in both local and Google but with different times (time overlap)
	googleEventID1 := "shared-google-id-1"
	localEvent1 := &entities.Event{
		ID:             "local-event-1",
		Title:          "Meeting",
		StartTime:      now,
		EndTime:        now.Add(1 * time.Hour),
		Location:       "Office",
		Description:    "Local version",
		GoogleEventID:  &googleEventID1,
	}
	
	googleEvent1 := &entities.Event{
		ID:             "google-event-1", 
		Title:          "Meeting Updated",
		StartTime:      now.Add(30 * time.Minute), // Different time = content diff
		EndTime:        now.Add(90 * time.Minute),
		Location:       "Office Remote", // Different location = content diff
		Description:    "Google version",
		GoogleEventID:  &googleEventID1,
	}

	// Two separate events that overlap in time
	localEvent2 := &entities.Event{
		ID:        "local-event-2",
		Title:     "Lunch",
		StartTime: now.Add(2 * time.Hour),
		EndTime:   now.Add(3 * time.Hour),
		Location:  "Restaurant A",
		Description: "Local lunch",
	}
	
	googleEvent2 := &entities.Event{
		ID:        "google-event-2",
		Title:     "Team Meeting", // Different event that overlaps
		StartTime: now.Add(2*time.Hour + 30*time.Minute), // Overlaps with lunch
		EndTime:   now.Add(4 * time.Hour),
		Location:  "Conference Room",
		Description: "Google meeting",
	}

	// Duplicate events (same title and similar time, but no time overlap)
	localEvent3 := &entities.Event{
		ID:        "local-event-3",
		Title:     "Conference Call",
		StartTime: now.Add(6 * time.Hour),      // 6 hours later
		EndTime:   now.Add(7 * time.Hour),      // 1 hour duration
	}
	
	googleEvent3 := &entities.Event{
		ID:        "google-event-3",
		Title:     "Conference Call", // Same title + time = duplicate
		StartTime: now.Add(6*time.Hour + 2*time.Minute), // Within 5min tolerance for duplicate (2min diff)
		EndTime:   now.Add(7*time.Hour + 5*time.Minute), // Ends 5min after local event ends - no overlap
	}

	localEvents := []*entities.Event{localEvent1, localEvent2, localEvent3}
	googleEvents := []*entities.Event{googleEvent1, googleEvent2, googleEvent3}

	t.Run("DetectConflicts", func(t *testing.T) {
		conflicts, err := service.DetectConflicts(ctx, userID, calendarSyncID, localEvents, googleEvents)
		
		assert.NoError(t, err)
		
		// Debug: print detected conflicts
		for i, conflict := range conflicts {
			t.Logf("Conflict %d: Type=%s, LocalEvent=%s, GoogleEvent=%s, Description=%s", 
				i, conflict.ConflictType, 
				getEventTitle(conflict.LocalEvent), 
				getEventTitle(conflict.GoogleEvent),
				conflict.Description)
		}
		
		assert.Len(t, conflicts, 4) // Should detect 4 conflicts: 1 content diff, 2 time overlaps, 1 duplicate

		// Check for content diff conflict (events with same GoogleEventID but different content)
		var contentDiffConflict *entities.SyncConflict
		for _, conflict := range conflicts {
			if conflict.ConflictType == entities.ConflictTypeContentDiff {
				contentDiffConflict = conflict
				break
			}
		}
		
		assert.NotNil(t, contentDiffConflict)
		assert.Equal(t, entities.ConflictTypeContentDiff, contentDiffConflict.ConflictType)
		assert.Contains(t, contentDiffConflict.Description, "Content differences detected")
		assert.Contains(t, contentDiffConflict.Description, "title")
		assert.Contains(t, contentDiffConflict.Description, "location")

		// Check for time overlap conflict (separate events that overlap in time)
		var timeOverlapConflict *entities.SyncConflict
		for _, conflict := range conflicts {
			if conflict.ConflictType == entities.ConflictTypeTimeOverlap {
				timeOverlapConflict = conflict
				break
			}
		}
		
		assert.NotNil(t, timeOverlapConflict)
		assert.Equal(t, entities.ConflictTypeTimeOverlap, timeOverlapConflict.ConflictType)
		assert.NotNil(t, timeOverlapConflict.LocalEvent)
		assert.NotNil(t, timeOverlapConflict.GoogleEvent)
		assert.Contains(t, timeOverlapConflict.Description, "Events overlap in time")

		// Check for duplicate event conflict (same title and similar time)
		var duplicateConflict *entities.SyncConflict
		for _, conflict := range conflicts {
			if conflict.ConflictType == entities.ConflictTypeDuplicateEvent {
				duplicateConflict = conflict
				break
			}
		}
		
		assert.NotNil(t, duplicateConflict)
		assert.Equal(t, entities.ConflictTypeDuplicateEvent, duplicateConflict.ConflictType)
		assert.Contains(t, duplicateConflict.Description, "Duplicate events detected")
	})
}

func TestSyncConflictService_ResolveConflict(t *testing.T) {
	mockConflictRepo := new(MockSyncConflictRepository)
	mockEventRepo := new(MockEventRepository)
	service := NewSyncConflictService(mockConflictRepo, mockEventRepo)
	
	ctx := context.Background()
	conflictID := "test-conflict-id"
	
	conflict := &entities.SyncConflict{
		ID:              conflictID,
		UserID:          "test-user-id",
		CalendarSyncID:  "test-calendar-sync-id",
		ConflictType:    entities.ConflictTypeContentDiff,
		LocalEvent:      &entities.Event{ID: "local-1", Title: "Local Event"},
		GoogleEvent:     &entities.Event{ID: "google-1", Title: "Google Event"},
		Status:          "pending",
		Description:     "Test conflict",
	}

	// Mock getting the conflict
	mockConflictRepo.On("GetByID", ctx, conflictID).Return(conflict, nil)
	
	// Mock updating the event repository for resolution
	mockEventRepo.On("Update", ctx, mock.AnythingOfType("*entities.Event")).Return(nil)
	
	// Mock updating the conflict resolution
	mockConflictRepo.On("UpdateResolution", ctx, conflictID, "Use local version", "test-user", mock.AnythingOfType("time.Time")).Return(nil)

	t.Run("ResolveWithUseLocal", func(t *testing.T) {
		action := &entities.ConflictResolutionAction{
			Action:     "use_local",
			Resolution: "Use local version",
		}

		err := service.ResolveConflict(ctx, conflictID, action, "test-user")
		
		assert.NoError(t, err)
		
		// Verify the event update was called
		mockEventRepo.AssertCalled(t, "Update", ctx, conflict.LocalEvent)
		
		// Verify the conflict resolution was updated
		mockConflictRepo.AssertCalled(t, "UpdateResolution", ctx, conflictID, "Use local version", "test-user", mock.AnythingOfType("time.Time"))
	})

	// Verify all mock expectations
	mockConflictRepo.AssertExpectations(t)
	mockEventRepo.AssertExpectations(t)
}

func TestSyncConflictService_GetPendingConflicts(t *testing.T) {
	mockConflictRepo := new(MockSyncConflictRepository)
	mockEventRepo := new(MockEventRepository)
	service := NewSyncConflictService(mockConflictRepo, mockEventRepo)
	
	ctx := context.Background()
	userID := entities.UserID("test-user-id")
	
	expectedConflicts := []*entities.SyncConflict{
		{
			ID:           "conflict-1",
			UserID:       userID,
			ConflictType: entities.ConflictTypeTimeOverlap,
			Status:       "pending",
		},
		{
			ID:           "conflict-2", 
			UserID:       userID,
			ConflictType: entities.ConflictTypeContentDiff,
			Status:       "pending",
		},
	}

	// Mock the repository call
	mockConflictRepo.On("GetPendingByUserID", ctx, userID).Return(expectedConflicts, nil)

	conflicts, err := service.GetPendingConflicts(ctx, userID)
	
	assert.NoError(t, err)
	assert.Len(t, conflicts, 2)
	
	// Verify all conflicts are pending
	for _, conflict := range conflicts {
		assert.Equal(t, "pending", conflict.Status)
	}

	// Verify mock expectations
	mockConflictRepo.AssertExpectations(t)
}

func TestSyncConflictService_GetConflictStats(t *testing.T) {
	mockConflictRepo := new(MockSyncConflictRepository)
	mockEventRepo := new(MockEventRepository)
	service := NewSyncConflictService(mockConflictRepo, mockEventRepo)
	
	ctx := context.Background()
	userID := entities.UserID("test-user-id")
	days := 30

	expectedStats := map[string]int{
		"time_overlap":    5,
		"content_diff":    3,
		"duplicate_event": 2,
		"deleted_event":   1,
	}

	// Mock the repository call
	mockConflictRepo.On("GetConflictStats", ctx, userID, mock.AnythingOfType("time.Time")).Return(expectedStats, nil)

	stats, err := service.GetConflictStats(ctx, userID, days)
	
	assert.NoError(t, err)
	assert.Equal(t, 5, stats["time_overlap"])
	assert.Equal(t, 3, stats["content_diff"])
	assert.Equal(t, 2, stats["duplicate_event"])
	assert.Equal(t, 1, stats["deleted_event"])

	// Verify mock expectations
	mockConflictRepo.AssertExpectations(t)
}

// Helper function to test time overlap detection logic
func TestDetectTimeOverlap(t *testing.T) {
	now := time.Now()
	
	testCases := []struct {
		name     string
		event1   *entities.Event
		event2   *entities.Event
		expected bool
	}{
		{
			name: "Complete overlap",
			event1: &entities.Event{
				StartTime: now,
				EndTime:   now.Add(2 * time.Hour),
			},
			event2: &entities.Event{
				StartTime: now.Add(30 * time.Minute),
				EndTime:   now.Add(90 * time.Minute),
			},
			expected: true,
		},
		{
			name: "No overlap - event2 after event1",
			event1: &entities.Event{
				StartTime: now,
				EndTime:   now.Add(1 * time.Hour),
			},
			event2: &entities.Event{
				StartTime: now.Add(2 * time.Hour),
				EndTime:   now.Add(3 * time.Hour),
			},
			expected: false,
		},
		{
			name: "Adjacent events - no overlap",
			event1: &entities.Event{
				StartTime: now,
				EndTime:   now.Add(1 * time.Hour),
			},
			event2: &entities.Event{
				StartTime: now.Add(1 * time.Hour), // Starts when event1 ends
				EndTime:   now.Add(2 * time.Hour),
			},
			expected: false,
		},
		{
			name: "Partial overlap - event2 starts before event1 ends",
			event1: &entities.Event{
				StartTime: now,
				EndTime:   now.Add(2 * time.Hour),
			},
			event2: &entities.Event{
				StartTime: now.Add(90 * time.Minute),
				EndTime:   now.Add(3 * time.Hour),
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := hasTimeOverlap(tc.event1, tc.event2)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// hasTimeOverlap checks if two events have overlapping time periods
func hasTimeOverlap(event1, event2 *entities.Event) bool {
	return event1.StartTime.Before(event2.EndTime) && event2.StartTime.Before(event1.EndTime)
}

// getEventTitle returns the title of an event or "nil" if event is nil
func getEventTitle(event *entities.Event) string {
	if event == nil {
		return "nil"
	}
	return event.Title
}