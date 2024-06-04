package rates

import (
	"context"

	"github.com/dmtrybogdanov/garantex/internal/repository/rates/modelRepo"
)

func (s *serv) Get(ctx context.Context, rates *modelRepo.RepoGetResponse) (int64, error) {
	id, err := s.ratesRepository.Get(ctx, rates)
	if err != nil {
		return 0, err
	}

	return id, nil
}
