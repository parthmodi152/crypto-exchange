package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/parthmodi152/crypto-exchange/orderbook"
)

func main() {
	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler

	ex := NewExchange()

	e.GET("/book/:market", ex.handleGetOrderBook)
	e.DELETE("/order/:id", ex.handleCancelOrder)
	e.POST("/order", ex.handlePlaceOrder)
	e.Start(":3000")
}

type Market string

const (
	MarketETH Market = "ETH"
)

type Exchange struct {
	OrderBook map[Market]*orderbook.OrderBook
}

func NewExchange() *Exchange {
	return &Exchange{
		OrderBook: map[Market]*orderbook.OrderBook{
			MarketETH: orderbook.NewOrderBook(),
		},
	}
}

type OrderType string

const (
	MarketOrder OrderType = "MARKET"
	LimitOrder  OrderType = "LIMIT"
)

type PlaceOrderRequest struct {
	Market Market    `json:"market"`
	Type   OrderType `json:"type"`
	Bid    bool      `json:"bid"`
	Size   float64   `json:"size"`
	Price  float64   `json:"price"`
}

type Order struct {
	ID        int64   `json:"id"`
	Size      float64 `json:"size"`
	Price     float64 `json:"price"`
	Bid       bool    `json:"bid"`
	Timestamp int64   `json:"timestamp"`
}

type OrderBookData struct {
	TotalBidVolume float64
	TotalAskVolume float64
	Asks           []*Order
	Bids           []*Order
}

func httpErrorHandler(err error, c echo.Context) {
	fmt.Println(err)
}

type MatchedOrder struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
	ID    int64   `json:"id"`
}

func (ex *Exchange) handlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	market := Market(placeOrderData.Market)
	ob := ex.OrderBook[market]
	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size)

	if placeOrderData.Type == LimitOrder {
		ob.PlaceLimitOrder(placeOrderData.Price, order)
		return c.JSON(http.StatusOK, "Limit Order placed")
	}

	if placeOrderData.Type == MarketOrder {
		matches := ob.PlaceMarketOrder(order)
		matchedOrders := make([]*MatchedOrder, len(matches))

		for i, match := range matches {

			var id int64
			if order.Bid {
				id = match.Ask.ID
			} else {
				id = match.Bid.ID
			}

			matchedOrders[i] = &MatchedOrder{
				ID:    id,
				Size:  match.SizeFilled,
				Price: match.Price,
			}
		}

		return c.JSON(http.StatusOK, map[string]any{"matches": matchedOrders})
	}

	return c.JSON(http.StatusBadRequest, "Invalid order type")
}

func (ex *Exchange) handleGetOrderBook(c echo.Context) error {
	market := Market(c.Param("market"))
	ob, ok := ex.OrderBook[market]

	if !ok {
		return c.JSON(http.StatusBadRequest, "Market not found")
	}

	orderbookData := OrderBookData{
		TotalBidVolume: ob.BidTotalVolume(),
		TotalAskVolume: ob.AskTotalVolume(),
		Asks:           []*Order{},
		Bids:           []*Order{},
	}

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {

			o := &Order{
				ID:        order.ID,
				Size:      order.Size,
				Price:     limit.Price,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}

			orderbookData.Asks = append(orderbookData.Asks, o)
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {

			o := &Order{
				ID:        order.ID,
				Size:      order.Size,
				Price:     limit.Price,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}

			orderbookData.Bids = append(orderbookData.Bids, o)
		}
	}

	return c.JSON(http.StatusOK, orderbookData)
}

func (ex *Exchange) handleCancelOrder(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid order ID")
	}

	ob := ex.OrderBook[MarketETH]
	order := ob.Orders[id]

	if order != nil {
		ob.CancelOrder(order)
		return c.JSON(http.StatusOK, "Order cancelled")
	}

	return c.JSON(http.StatusBadRequest, "Order not found")
}
