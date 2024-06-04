package repository

import (
	"context"
	modelRepo "github.com/dmtrybogdanov/garantex/internal/repository/rates/modelRepo"
)

type RatesRepository interface {
	Get(ctx context.Context, rates *modelRepo.RepoGetResponse) (int64, error)
}
