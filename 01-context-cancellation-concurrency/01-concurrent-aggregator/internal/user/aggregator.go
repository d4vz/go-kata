package user

import (
	"concurrent-aggregator/internal/order"
	"concurrent-aggregator/internal/profile"
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

type (
	ProfileService interface {
		GetProfile(ctx context.Context, id int) (profile.Profile, error)
	}
	OrderService interface {
		GetOrders(ctx context.Context, id int) (order.Order, error)
	}
	Aggregator struct {
		timeout        time.Duration
		logger         *slog.Logger
		profileService ProfileService
		orderService   OrderService
	}
	Option func(*Aggregator)
)

func WithTimeout(timeout time.Duration) Option {
	return func(a *Aggregator) {
		a.timeout = timeout
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(a *Aggregator) {
		a.logger = logger
	}
}

func WithProfileService(profileService ProfileService) Option {
	return func(a *Aggregator) {
		a.profileService = profileService
	}
}

func WithOrderService(orderService OrderService) Option {
	return func(a *Aggregator) {
		a.orderService = orderService
	}
}

func NewAggregator(opts ...Option) *Aggregator {
	a := &Aggregator{
		timeout: 10 * time.Second,
		logger:  slog.Default(),
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *Aggregator) Aggregate(ctx context.Context, id int) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	var profile profile.Profile
	var orders order.Order

	eg.Go(func() error {
		var err error
		profile, err = a.profileService.GetProfile(ctx, id)
		return err
	})

	eg.Go(func() error {
		var err error
		orders, err = a.orderService.GetOrders(ctx, id)
		return err
	})

	if err := eg.Wait(); err != nil {
		return "", err
	}

	return fmt.Sprintf("User: %s | Orders: %d", profile.Name, len(orders.Items)), nil
}
