package converter

import (
	"github.com/dmtrybogdanov/garantex/internal/api/rates/model"
	"github.com/dmtrybogdanov/garantex/internal/repository/rates/modelRepo"
)

func ToRatesFromGarantex(askOrdersJSON, bidOrdersJSON []byte, timestamp int64) *model.GetResponse {
	return &model.GetResponse{
		Asks:      askOrdersJSON,
		Bids:      bidOrdersJSON,
		Timestamp: timestamp,
	}
}

func ToRepoFromService(rates *model.GetResponse) *modelRepo.RepoGetResponse {
	return &modelRepo.RepoGetResponse{
		Asks:      rates.Asks,
		Bids:      rates.Bids,
		Timestamp: rates.Timestamp,
	}
}
