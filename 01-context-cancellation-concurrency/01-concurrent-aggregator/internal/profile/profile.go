package profile

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

type Profile struct {
	Name string
}

type profileService struct {
	delay      time.Duration
	throwError bool
}

func NewProfileService(delay time.Duration, throwError bool) *profileService {
	return &profileService{delay: delay, throwError: throwError}
}

func (s *profileService) GetProfile(ctx context.Context, id int) (Profile, error) {
	slog.Info("Getting profile", "id", id)

	select {
	case <-time.After(s.delay):
		if s.throwError {
			slog.Error("Failed to get profile", "id", id)
			return Profile{}, errors.New("failed to get profile")
		}
		slog.Info("Profile got", "id", id)
		return Profile{Name: "John Doe"}, nil
	case <-ctx.Done():
		slog.Error("[Profile] Context done", "id", id, "error", ctx.Err())
		return Profile{}, ctx.Err()
	}
}
