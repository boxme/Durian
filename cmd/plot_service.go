package main

import (
	"context"
	"time"
)

type Plot struct {
	Id        string
	CreatedAt time.Time
}

type PlotService interface {
	AddPlot(ctx context.Context, plot_id string, info string) (*Plot, error)
	EditPlot(ctx context.Context, plot_id string, info string) (*Plot, error)
}
