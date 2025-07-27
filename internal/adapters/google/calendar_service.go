package google

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/calendar/v3"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type CalendarService struct {
	oauth2Service *OAuth2Service
}

func NewCalendarService(oauth2Service *OAuth2Service) *CalendarService {
	return &CalendarService{
		oauth2Service: oauth2Service,
	}
}

type GoogleCalendarEvent struct {
	ID          string    `json:"id"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	AllDay      bool      `json:"all_day"`
	CalendarID  string    `json:"calendar_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
	Attendees   []EventAttendee `json:"attendees"`
}

// CalendarEvent represents a simplified Google Calendar event for webhook processing
type CalendarEvent struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Status      string `json:"status"`
	Start       struct {
		DateTime string `json:"dateTime"`
		Date     string `json:"date"`
		TimeZone string `json:"timeZone"`
	} `json:"start"`
	End struct {
		DateTime string `json:"dateTime"`
		Date     string `json:"date"`
		TimeZone string `json:"timeZone"`
	} `json:"end"`
	Attendees []EventAttendee `json:"attendees"`
}

type EventAttendee struct {
	Email string `json:"email"`
	Name  string `json:"displayName"`
}

type CalendarListItem struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Primary     bool   `json:"primary"`
	AccessRole  string `json:"access_role"`
}

func (s *CalendarService) GetCalendars(ctx context.Context, accessToken string) ([]*CalendarListItem, error) {
	service, err := s.oauth2Service.CreateCalendarService(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	calendarList, err := service.CalendarList.List().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list calendars: %w", err)
	}

	var calendars []*CalendarListItem
	for _, item := range calendarList.Items {
		calendars = append(calendars, &CalendarListItem{
			ID:          item.Id,
			Summary:     item.Summary,
			Description: item.Description,
			Primary:     item.Primary,
			AccessRole:  item.AccessRole,
		})
	}

	return calendars, nil
}

func (s *CalendarService) GetEvents(ctx context.Context, accessToken, calendarID string, timeMin, timeMax time.Time) ([]*GoogleCalendarEvent, error) {
	service, err := s.oauth2Service.CreateCalendarService(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	eventsCall := service.Events.List(calendarID).
		TimeMin(timeMin.Format(time.RFC3339)).
		TimeMax(timeMax.Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime")

	events, err := eventsCall.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	var googleEvents []*GoogleCalendarEvent
	for _, event := range events.Items {
		googleEvent := s.convertCalendarEvent(event, calendarID)
		googleEvents = append(googleEvents, googleEvent)
	}

	return googleEvents, nil
}

func (s *CalendarService) CreateEvent(ctx context.Context, accessToken, calendarID string, event *entities.Event) (*GoogleCalendarEvent, error) {
	service, err := s.oauth2Service.CreateCalendarService(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	googleEvent := &calendar.Event{
		Summary:     event.Title,
		Description: event.Description,
		Location:    event.Location,
		Start: &calendar.EventDateTime{
			DateTime: event.StartTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: event.EndTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
	}

	createdEvent, err := service.Events.Insert(calendarID, googleEvent).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return s.convertCalendarEvent(createdEvent, calendarID), nil
}

func (s *CalendarService) UpdateEvent(ctx context.Context, accessToken, calendarID, eventID string, event *entities.Event) (*GoogleCalendarEvent, error) {
	service, err := s.oauth2Service.CreateCalendarService(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	googleEvent := &calendar.Event{
		Summary:     event.Title,
		Description: event.Description,
		Location:    event.Location,
		Start: &calendar.EventDateTime{
			DateTime: event.StartTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: event.EndTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
	}

	updatedEvent, err := service.Events.Update(calendarID, eventID, googleEvent).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return s.convertCalendarEvent(updatedEvent, calendarID), nil
}

func (s *CalendarService) DeleteEvent(ctx context.Context, accessToken, calendarID, eventID string) error {
	service, err := s.oauth2Service.CreateCalendarService(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("failed to create calendar service: %w", err)
	}

	err = service.Events.Delete(calendarID, eventID).Do()
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

func (s *CalendarService) SyncEvents(ctx context.Context, accessToken, calendarID string, localEvents []*entities.Event) ([]*GoogleCalendarEvent, error) {
	// Get events from Google Calendar
	now := time.Now()
	pastMonth := now.AddDate(0, -1, 0)
	futureMonth := now.AddDate(0, 1, 0)

	googleEvents, err := s.GetEvents(ctx, accessToken, calendarID, pastMonth, futureMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get Google events: %w", err)
	}

	// TODO: Implement intelligent sync logic
	// For now, just return Google events
	return googleEvents, nil
}

func (s *CalendarService) convertCalendarEvent(event *calendar.Event, calendarID string) *GoogleCalendarEvent {
	googleEvent := &GoogleCalendarEvent{
		ID:          event.Id,
		Summary:     event.Summary,
		Description: event.Description,
		Location:    event.Location,
		CalendarID:  calendarID,
		AllDay:      false,
	}

	// Parse start time
	if event.Start != nil {
		if event.Start.DateTime != "" {
			if startTime, err := time.Parse(time.RFC3339, event.Start.DateTime); err == nil {
				googleEvent.StartTime = startTime
			}
		} else if event.Start.Date != "" {
			if startTime, err := time.Parse("2006-01-02", event.Start.Date); err == nil {
				googleEvent.StartTime = startTime
				googleEvent.AllDay = true
			}
		}
	}

	// Parse end time
	if event.End != nil {
		if event.End.DateTime != "" {
			if endTime, err := time.Parse(time.RFC3339, event.End.DateTime); err == nil {
				googleEvent.EndTime = endTime
			}
		} else if event.End.Date != "" {
			if endTime, err := time.Parse("2006-01-02", event.End.Date); err == nil {
				googleEvent.EndTime = endTime
				googleEvent.AllDay = true
			}
		}
	}

	// Parse created/updated times
	if event.Created != "" {
		if createdTime, err := time.Parse(time.RFC3339, event.Created); err == nil {
			googleEvent.CreatedAt = createdTime
		}
	}

	if event.Updated != "" {
		if updatedTime, err := time.Parse(time.RFC3339, event.Updated); err == nil {
			googleEvent.UpdatedAt = updatedTime
		}
	}

	return googleEvent
}

// SetupWebhook sets up a webhook for a Google Calendar
func (s *CalendarService) SetupWebhook(ctx context.Context, accessToken, calendarID, channelID, webhookURL string) error {
	service, err := s.oauth2Service.CreateCalendarService(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("failed to create calendar service: %w", err)
	}

	// Create a watch request
	channel := &calendar.Channel{
		Id:      channelID,
		Type:    "web_hook",
		Address: webhookURL,
		Payload: true,
	}

	// Set up the webhook for the calendar events
	_, err = service.Events.Watch(calendarID, channel).Do()
	if err != nil {
		return fmt.Errorf("failed to setup webhook: %w", err)
	}

	return nil
}

// StopWebhook stops a webhook for a Google Calendar
func (s *CalendarService) StopWebhook(ctx context.Context, accessToken, channelID, resourceID string) error {
	service, err := s.oauth2Service.CreateCalendarService(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("failed to create calendar service: %w", err)
	}

	// Create stop request
	channel := &calendar.Channel{
		Id:         channelID,
		ResourceId: resourceID,
	}

	// Stop the webhook
	err = service.Channels.Stop(channel).Do()
	if err != nil {
		return fmt.Errorf("failed to stop webhook: %w", err)
	}

	return nil
}