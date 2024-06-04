package rates

import (
	"github.com/dmtrybogdanov/garantex/internal/repository"
	"github.com/dmtrybogdanov/garantex/internal/service"
)

type serv struct {
	ratesRepository repository.RatesRepository
}

func NewService(ratesRepository repository.RatesRepository) service.RatesService {
	return &serv{
		ratesRepository: ratesRepository,
	}
}
