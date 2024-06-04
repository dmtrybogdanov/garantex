package rates

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/dmtrybogdanov/garantex/internal/repository"
	modelRepo "github.com/dmtrybogdanov/garantex/internal/repository/rates/modelRepo"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("fiber-server")

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.RatesRepository {
	return &repo{db: db}
}

func (r *repo) Get(ctx context.Context, rates *modelRepo.RepoGetResponse) (int64, error) {
	_, span := tracer.Start(ctx, "get rates data postgres", oteltrace.WithAttributes())
	defer span.End()

	builderInsert := sq.Insert("market_data").
		PlaceholderFormat(sq.Dollar).
		Columns("timestamp", "asks", "bids").
		Values(rates.Timestamp, rates.Asks, rates.Bids).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %v", err)
	}

	var ID int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&ID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert market data: %v", err)
	}

	log.Printf("inserted market data: %d", ID)

	return ID, nil
}
