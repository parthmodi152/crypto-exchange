package orderbook

import (
	"fmt"
	"sort"
)

type OrderBook struct {
	asks []*Limit
	bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		asks: []*Limit{},
		bids: []*Limit{},

		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *OrderBook) AskTotalVolume() float64 {
	totalVolume := 0.0

	for _, limit := range ob.Asks() {
		totalVolume += limit.TotalVolume
	}

	return totalVolume
}

func (ob *OrderBook) BidTotalVolume() float64 {
	totalVolume := 0.0

	for _, limit := range ob.Bids() {
		totalVolume += limit.TotalVolume
	}

	return totalVolume
}

func (ob *OrderBook) Asks() Limits {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks
}

func (ob *OrderBook) Bids() Limits {
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids
}

func (ob *OrderBook) PlaceMarketOrder(o *Order) []Match {
	matches := []Match{}
	limitsToClear := []*Limit{}

	if o.Bid {
		if o.Size > ob.AskTotalVolume() {
			panic(fmt.Errorf("not enough [%.2f] volume to fill order [%.2f]", ob.AskTotalVolume(), o.Size))
		}

		for _, limit := range ob.Asks() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				limitsToClear = append(limitsToClear, limit)
			}
		}

	} else {
		if o.Size > ob.BidTotalVolume() {
			panic(fmt.Errorf("not enough [%.2f] volume to fill order [%.2f]", ob.BidTotalVolume(), o.Size))
		}

		for _, limit := range ob.Bids() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				limitsToClear = append(limitsToClear, limit)
			}
		}
	}

	for _, limit := range limitsToClear {
		ob.ClearLimit(true, limit)
	}

	return matches
}

func (ob *OrderBook) PlaceLimitOrder(price float64, o *Order) {
	if o.Bid {
		if _, ok := ob.BidLimits[price]; !ok {
			l := NewLimit(price)
			ob.BidLimits[price] = l
			ob.bids = append(ob.bids, l)
		}

		ob.BidLimits[price].AddOrder(o)
	} else {
		if _, ok := ob.AskLimits[price]; !ok {
			l := NewLimit(price)
			ob.AskLimits[price] = l
			ob.asks = append(ob.asks, l)
		}

		ob.AskLimits[price].AddOrder(o)
	}
}

func (ob *OrderBook) CancelOrder(o *Order) {
	limit := o.Limit

	limit.DeleteOrder(o)

	if len(limit.Orders) == 0 {
		ob.ClearLimit(o.Bid, limit)
	}
}

func (ob *OrderBook) ClearLimit(bid bool, l *Limit) {
	if bid {
		delete(ob.BidLimits, l.Price)
		for i, limit := range ob.bids {
			if limit == l {
				ob.bids = append(ob.bids[:i], ob.bids[i+1:]...)
				break
			}
		}
	} else {
		delete(ob.AskLimits, l.Price)

		for i, limit := range ob.asks {
			if limit == l {
				ob.asks = append(ob.asks[:i], ob.asks[i+1:]...)
				break
			}
		}
	}
}
