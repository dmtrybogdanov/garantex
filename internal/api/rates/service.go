package rates

import (
	"github.com/dmtrybogdanov/garantex/internal/service"
	"github.com/dmtrybogdanov/garantex/pkg/rates_v1"
)

type Implementation struct {
	rates_v1.UnimplementedRatesV1Server
	ratesService service.RatesService
}

func NewImplementation(ratesService service.RatesService) *Implementation {
	return &Implementation{
		ratesService: ratesService,
	}
}
