package model

type GetResponse struct {
	Asks      []byte
	Bids      []byte
	Timestamp int64
}
