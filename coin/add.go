package coin

import "fmt"

type Coin uint

func (c1 *Coin) Add(c2 *Coin) error {
	if (*c1)+(*c2) <= MaxCoinValue {
		*c1 += (*c2)
		return nil
	} else {
		return fmt.Errorf("Exceeds maximum coin value")
	}
}

const (
	MaxCoinValue = 1000000
)

func ValidCoin(c Coin) bool {
	return 0 <= c && c <= MaxCoinValue
}
