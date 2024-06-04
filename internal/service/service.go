package service

import (
	"context"

	"github.com/dmtrybogdanov/garantex/internal/repository/rates/modelRepo"
)

type RatesService interface {
	Get(ctx context.Context, rates *modelRepo.RepoGetResponse) (int64, error)
}
