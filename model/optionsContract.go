// model/optionsContract.go
package model

import (
	"time"
)

// OptionsContract represents the data structure of an options contract
type OptionsContract struct {
	Type           string    `json:"type"`
	StrikePrice    float64   `json:"strike_price"`
	Bid            float64   `json:"bid"`
	Ask            float64   `json:"ask"`
	ExpirationDate time.Time `json:"expiration_date"`
	LongShort      string    `json:"long_short"`
}
