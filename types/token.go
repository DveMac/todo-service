package types

import (
	"time"
)

type Token struct {
	Expiry time.Time `json:"expiry"`
	Value  string    `json:"value"`
}
