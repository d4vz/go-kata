package main

import (
	"concurrent-aggregator/internal/order"
	"concurrent-aggregator/internal/profile"
	"concurrent-aggregator/internal/user"
	"context"
	"log/slog"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	profileService := profile.NewProfileService(2*time.Second, true)
	orderService := order.NewOrderService(10*time.Second, false)

	aggregator := user.NewAggregator(
		user.WithTimeout(1*time.Second),
		user.WithLogger(logger),
		user.WithProfileService(profileService),
		user.WithOrderService(orderService),
	)

	result, err := aggregator.Aggregate(context.Background(), 1)

	if err != nil {
		logger.Error("Failed to aggregate", "error", err)
		return
	}

	logger.Info("Aggregated", "result", result)
}
