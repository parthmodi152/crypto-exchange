package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

type Orders []*Order

func (o Orders) Len() int      { return len(o) }
func (o Orders) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool {
	return o[i].Timestamp < o[j].Timestamp
}

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("Order{Size: %.2f, Bid: %t, Timestamp: %d}", o.Size, o.Bid, o.Timestamp)
}

type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

type ByBestAsk struct{ Limits }

func (a ByBestAsk) Len() int      { return len(a.Limits) }
func (a ByBestAsk) Swap(i, j int) { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsk) Less(i, j int) bool {
	return a.Limits[i].Price < a.Limits[j].Price
}

type ByBestBid struct{ Limits }

func (b ByBestBid) Len() int      { return len(b.Limits) }
func (b ByBestBid) Swap(i, j int) { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestBid) Less(i, j int) bool {
	return b.Limits[i].Price > b.Limits[j].Price
}

func (l *Limit) String() string {
	return fmt.Sprintf("Limit{Price: %.2f, Orders: %v, TotalVolume: %.2f}", l.Price, l.Orders, l.TotalVolume)
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	for i, order := range l.Orders {
		if order == o {
			l.Orders = append(l.Orders[:i], l.Orders[i+1:]...)
			l.TotalVolume -= o.Size
			break
		}
	}

	o.Limit = nil
	l.TotalVolume -= o.Size

	sort.Sort(l.Orders)
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

type OrderBook struct {
	Asks []*Limit
	Bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Asks: []*Limit{},
		Bids: []*Limit{},

		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *OrderBook) PlaceOrder(price float64, o *Order) []Match {
	// 1. Todo: Try to match to order

	// 2. add the rest of the order to the orderbook
	if o.Size > 0.0 {
		ob.add(price, o)
	}

	return []Match{}
}

func (ob *OrderBook) add(price float64, o *Order) {
	if o.Bid {
		if _, ok := ob.BidLimits[price]; !ok {
			l := NewLimit(price)
			ob.BidLimits[price] = l
			ob.Bids = append(ob.Bids, l)
		}

		ob.BidLimits[price].AddOrder(o)
	} else {
		if _, ok := ob.AskLimits[price]; !ok {
			l := NewLimit(price)
			ob.AskLimits[price] = l
			ob.Asks = append(ob.Asks, l)
		}

		ob.AskLimits[price].AddOrder(o)
	}
}
