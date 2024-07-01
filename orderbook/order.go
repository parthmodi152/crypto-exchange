package orderbook

import (
	"fmt"
	"math/rand"
	"time"
)

type Order struct {
	ID        int64
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

func (o *Order) String() string {
	return fmt.Sprintf("Order{ID: %d, Size: %.2f, Bid: %t, Timestamp: %d}", o.ID, o.Size, o.Bid, o.Timestamp)
}

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		ID:        int64(rand.Intn(1000000)),
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) isFilled() bool {
	return o.Size == 0
}
