package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/andranikuz/smart-goal-calendar/internal/domain/entities"
)

type OAuth2Service struct {
	config *oauth2.Config
}

type GoogleTokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	Expiry       time.Time `json:"expiry"`
}

func NewOAuth2Service(clientID, clientSecret, redirectURL string) *OAuth2Service {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			calendar.CalendarReadonlyScope,
			calendar.CalendarEventsScope,
		},
		Endpoint: google.Endpoint,
	}

	return &OAuth2Service{
		config: config,
	}
}

func (s *OAuth2Service) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (s *OAuth2Service) ExchangeCode(ctx context.Context, code string) (*GoogleTokens, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return &GoogleTokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	}, nil
}

func (s *OAuth2Service) RefreshToken(ctx context.Context, refreshToken string) (*GoogleTokens, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := s.config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &GoogleTokens{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		TokenType:    newToken.TokenType,
		Expiry:       newToken.Expiry,
	}, nil
}

func (s *OAuth2Service) CreateCalendarService(ctx context.Context, accessToken string) (*calendar.Service, error) {
	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	client := s.config.Client(ctx, token)
	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	return service, nil
}

func (s *OAuth2Service) ValidateToken(ctx context.Context, accessToken string) error {
	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	client := s.config.Client(ctx, token)
	
	// Test the token by making a simple API call
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=" + accessToken)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token: status %d", resp.StatusCode)
	}

	return nil
}

type TokenInfo struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

func (s *OAuth2Service) GetUserInfo(ctx context.Context, accessToken string) (*TokenInfo, error) {
	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	client := s.config.Client(ctx, token)
	
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	var userInfo TokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}