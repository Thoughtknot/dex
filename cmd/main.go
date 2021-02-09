package main

import (
	"dex/pkg/calc"
	"dex/pkg/data"
	nhttp "dex/pkg/http"
	"dex/pkg/output"
	"dex/pkg/pubsub"
	"fmt"
)

func main() {
	portfolioCh := make(chan *calc.Portfolio, 1)
	positionCh := make(chan *data.Position, 1)
	marketDataCh := make(chan interface{}, 1)
	go nhttp.AcceptHttp(portfolioCh, positionCh, marketDataCh)
	broker := pubsub.NewBroker()
	go broker.Start()
	initialMarginCh := make(chan calc.Margin, 1)
	go output.Start(initialMarginCh)
	riskCalculators := map[string]chan *data.Position{}

	marketDataCache := []interface{}{}
	for {
		select {
		case msg := <-portfolioCh:
			ch := make(chan *data.Position, 1)
			riskCalculators[msg.Id] = ch
			mdCh := broker.Subscribe()
			go msg.Calculate(ch, mdCh, initialMarginCh)
			for _, v := range marketDataCache {
				mdCh <- v
			}
		case msg := <-positionCh:
			posCh := riskCalculators[msg.PortfolioId]
			posCh <- msg
		case msg := <-marketDataCh:
			marketDataCache = append(marketDataCache, msg)
			broker.Publish(msg)
	}
}
