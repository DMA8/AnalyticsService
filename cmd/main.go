package main

import (
	"gitlab.com/g6834/team31/analytics/internal/config"
	"gitlab.com/g6834/team31/analytics/internal/app"
	"context"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()
	app.Run(ctx, cfg)
}
