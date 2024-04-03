package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Plot struct {
	Id        string
	Info      string
	CreatedAt time.Time
}

type PlotService interface {
	AddPlot(ctx context.Context, plotId string, info string) (*Plot, error)
	EditPlot(ctx context.Context, plotId string, info string) (*Plot, error)
}

type plotService struct {
	Db *pgxpool.Pool
}

func CreatePlotService(db *pgxpool.Pool) PlotService {
	return &plotService{
		Db: db,
	}
}

func (ps *plotService) AddPlot(ctx context.Context, plotId string, info string) (*Plot, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		query := `INSERT INTO plots (plot_id, info) VALUES (@plotId, @info) RETURNING created_at;`
		args := pgx.NamedArgs{
			"plotId": plotId,
			"info":   info,
		}

		var created_at string
		err := ps.Db.QueryRow(ctx, query, args).Scan(&created_at)
		if err != nil {
			return nil, err
		}

		timeNow, err := time.Parse("2006-01-02", created_at)
		if err != nil {
			return nil, err
		}

		plot := &Plot{
			Id:        plotId,
			Info:      info,
			CreatedAt: timeNow,
		}

		return plot, nil
	}
}

func (ps *plotService) EditPlot(ctx context.Context, plotId string, info string) (*Plot, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		query := `UPDATE plots SET info = @info WHERE plot_id = @plotId RETURNING created_at;`
		args := pgx.NamedArgs{
			"plotId": plotId,
			"info":   info,
		}

		var created_at string
		err := ps.Db.QueryRow(ctx, query, args).Scan(&created_at)
		if err != nil {
			return nil, err
		}

		timeNow, err := time.Parse("2006-01-02", created_at)
		if err != nil {
			return nil, err
		}

		plot := &Plot{
			Id:        plotId,
			Info:      info,
			CreatedAt: timeNow,
		}

		return plot, nil
	}
}
