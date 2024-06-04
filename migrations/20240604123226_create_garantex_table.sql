-- +goose Up
CREATE TABLE IF NOT EXISTS market_data (
                                           id SERIAL PRIMARY KEY,
                                           timestamp INT,
                                           asks JSON,
                                           bids JSON
);

-- +goose Down
drop table market_data;

