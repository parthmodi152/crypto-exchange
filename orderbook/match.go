package orderbook

import (
	"fmt"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

func (m *Match) String() string {
	return fmt.Sprintf("Match{Ask: %s, Bid: %s, SizeFilled: %.2f, Price: %.2f}", m.Ask, m.Bid, m.SizeFilled, m.Price)
}
