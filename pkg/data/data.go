package data

type MarketData interface {
	GetInstrumentId() string
	GetValue(time float64) float64
}

type SpotPrice struct {
	InstrumentId string
	Price        float64
}

type YieldCurve struct {
	InstrumentId string
	Times        []float64
	Price        []float64
}

func (p *YieldCurve) GetInstrumentId() string {
	return p.InstrumentId
}

func (p *YieldCurve) GetValue(time float64) float64 {
	if time < p.Times[0] {
		return p.Price[0]
	}
	for i := 0; i < len(p.Times)-1; i++ {
		tv := p.Times[i]
		if tv == time {
			return p.Price[i]
		}
		tn := p.Times[i+1]
		if time <= tn {
			p1 := p.Price[i+1]
			p0 := p.Price[i]
			return p0 + (time-tv)*(p1-p0)/(tn-tv)
		}
	}
	return p.Price[len(p.Price)-1]
}

func (p *SpotPrice) GetInstrumentId() string {
	return p.InstrumentId
}

func (p *SpotPrice) GetValue(time float64) float64 {
	return p.Price
}

type Position struct {
	PortfolioId  string
	InstrumentId string
	Quantity     int
	TimeToExpiry float64 
}
