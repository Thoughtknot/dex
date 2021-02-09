package calc

import (
	"dex/pkg/data"
	"fmt"
	"sort"
)

type Portfolio struct {
	Id string
}

type Margin interface {
	GetPortfolioId() string
	GetInitialMargin() float64
}

type InitialMargin struct {
	PortfolioId string
	IM          float64
}

func (i *InitialMargin) GetPortfolioId() string {
	return i.PortfolioId
}

func (i *InitialMargin) GetInitialMargin() float64 {
	return i.IM
}

type Matrix struct {
	Percentile float64
	Scenarios  []Scenario
}

type Scenario struct {
	Id                 string
	ShiftsByInstrument map[string]float64
}

func (p *Portfolio) CalculateTotalIM(m *Matrix, positions map[string]*data.Position, md map[string]data.MarketData) Margin {
	pnlVec := make([]float64, len(m.Scenarios))
	for i, scenario := range m.Scenarios {
		pnl := 0.0
		for k, pos := range positions {
			if mdVal, prs := md[k]; prs {
				pnl += float64(pos.Quantity) * mdVal.GetValue(pos.TimeToExpiry) * scenario.ShiftsByInstrument[k]
			}
		}
		pnlVec[i] = pnl
	}
	sort.Slice(pnlVec, func(i, j int) bool {
		return pnlVec[i] < pnlVec[j]
	})
	percentile := (m.Percentile / 100.0) * float64(len(pnlVec))
	return &InitialMargin{PortfolioId: p.Id, IM: pnlVec[int(percentile)]}
}

func (p *Portfolio) Calculate(pos <-chan *data.Position, md <-chan interface{}, im chan<- Margin) {
	fmt.Println("Calculating portfolio: ", p.Id)
	marketData := map[string]data.MarketData{}
	positions := map[string]*data.Position{}
	var matrix *Matrix
	for {
		select {
		case msg := <-md:
			switch msg.(type) {
			case data.MarketData:
				mdVal := msg.(data.MarketData)
				marketData[mdVal.GetInstrumentId()] = mdVal
			case *Matrix:
				matrix = msg.(*Matrix)
			}
			if matrix != nil {
				val := p.CalculateTotalIM(matrix, positions, marketData)
				im <- val
			}
		case msg := <-pos:
			positions[msg.InstrumentId] = msg
			if matrix != nil {
				val := p.CalculateTotalIM(matrix, positions, marketData)
				im <- val
			}
		}
	}
}
