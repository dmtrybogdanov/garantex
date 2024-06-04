package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dmtrybogdanov/garantex/internal/converter"
	"github.com/dmtrybogdanov/garantex/internal/repository/rates/modelRepo"
	"github.com/dmtrybogdanov/garantex/pkg/rates_v1"
)

const (
	apiURL       = "https://garantex.org/api/v2/depth"
	currencyPair = "usdtrub"
)

type Order struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

type DepthResponse struct {
	Timestamp int64   `json:"timestamp"`
	Asks      []Order `json:"asks"`
	Bids      []Order `json:"bids"`
}

func getRatesFromGarantex(req *rates_v1.GetRequest) ([]Order, []Order, int64, error) {
	resp, err := http.Get(fmt.Sprintf("%s?market=%s", apiURL, req.Market))
	if err != nil {
		return nil, nil, 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	var depth DepthResponse
	if err := json.NewDecoder(resp.Body).Decode(&depth); err != nil {
		return nil, nil, 0, err
	}

	if len(depth.Asks) == 0 || len(depth.Bids) == 0 {
		return nil, nil, 0, fmt.Errorf("no asks or bids found")
	}
	timestamp := depth.Timestamp

	return depth.Asks, depth.Bids, timestamp, nil
}

func (i *Implementation) Get(ctx context.Context, req *rates_v1.GetRequest) (*rates_v1.GetResponse, error) {
	asks, bids, timestamp, err := getRatesFromGarantex(req)
	if err != nil {
		return nil, err
	}

	var askOrders []*rates_v1.Order
	for _, ask := range asks {
		askOrder := &rates_v1.Order{
			Price:  ask.Price,
			Volume: ask.Volume,
			Amount: ask.Amount,
			Factor: ask.Factor,
			Type:   ask.Type,
		}
		askOrders = append(askOrders, askOrder)
	}

	var bidOrders []*rates_v1.Order
	for _, bid := range bids {
		bidOrder := &rates_v1.Order{
			Price:  bid.Price,
			Volume: bid.Volume,
			Amount: bid.Amount,
			Factor: bid.Factor,
			Type:   bid.Type,
		}
		bidOrders = append(bidOrders, bidOrder)
	}

	askOrdersJSON, err := json.Marshal(askOrders)
	if err != nil {
		return nil, err
	}
	bidOrdersJSON, err := json.Marshal(bidOrders)
	if err != nil {
		return nil, err
	}

	id, err := i.ratesService.Get(ctx, (*modelRepo.RepoGetResponse)(converter.ToRatesFromGarantex(askOrdersJSON, bidOrdersJSON, timestamp)))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted market data with id: %d", id)

	return nil, nil
}
