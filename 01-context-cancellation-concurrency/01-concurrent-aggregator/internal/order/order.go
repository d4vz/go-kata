package order

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

type Order struct {
	ID    int
	User  string
	Items []string
}

type orderService struct {
	delay      time.Duration
	throwError bool
}

func NewOrderService(delay time.Duration, throwError bool) *orderService {
	return &orderService{delay: delay, throwError: throwError}
}

func (s *orderService) GetOrders(ctx context.Context, id int) (Order, error) {
	slog.Info("Getting orders", "id", id)

	select {
	case <-time.After(s.delay):
		if s.throwError {
			slog.Error("Failed to get orders", "id", id)
			return Order{}, errors.New("failed to get orders")
		}
		slog.Info("Orders got", "id", id)
		return Order{ID: 1, User: "John Doe", Items: []string{"Item 1", "Item 2"}}, nil
	case <-ctx.Done():
		slog.Error("[Order] Context done", "id", id, "error", ctx.Err())
		return Order{}, ctx.Err()
	}
}
