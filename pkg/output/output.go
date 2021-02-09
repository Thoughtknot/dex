package output

import (
	"dex/pkg/calc"
	"fmt"
)

func Start(im <-chan calc.Margin) {
	for {
		select {
		case msg := <-im:
			// No op, publish to kafka or something
			fmt.Println("Received result: ", msg.GetPortfolioId(), msg.GetPortfolioId())
		}
	}
}
