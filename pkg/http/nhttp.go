package nhttp

import (
	"dex/pkg/calc"
	"dex/pkg/data"
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpState struct {
	pch   chan<- *calc.Portfolio
	posch chan<- *data.Position
	mdch  chan<- interface{}
}

func (h *HttpState) PutYieldCurve(w http.ResponseWriter, r *http.Request) {
	var yield data.YieldCurve
	json.NewDecoder(r.Body).Decode(&yield)
	fmt.Println("Yield: ", yield)
	h.mdch <- &yield
}

func (h *HttpState) PutPrice(w http.ResponseWriter, r *http.Request) {
	var price data.SpotPrice
	json.NewDecoder(r.Body).Decode(&price)
	fmt.Println("Price: ", price)
	h.mdch <- &price
}

func (h *HttpState) PutPosition(w http.ResponseWriter, r *http.Request) {
	var position data.Position
	json.NewDecoder(r.Body).Decode(&position)
	fmt.Println("Position: ", position)
	h.posch <- &position
}

func (h *HttpState) PutMatrix(w http.ResponseWriter, r *http.Request) {
	var matrix calc.Matrix
	json.NewDecoder(r.Body).Decode(&matrix)
	fmt.Println("Matrix: ", matrix)
	h.mdch <- &matrix
}
func (h *HttpState) PutPortfolio(w http.ResponseWriter, r *http.Request) {
	var portfolio calc.Portfolio
	json.NewDecoder(r.Body).Decode(&portfolio)
	fmt.Println("Portfolio: ", portfolio)
	h.pch <- &portfolio
}

func AcceptHttp(pch chan<- *calc.Portfolio, posch chan<- *data.Position, mdch chan<- interface{}) {
	h := HttpState{pch, posch, mdch}
	http.HandleFunc("/portfolio", h.PutPortfolio)
	http.HandleFunc("/position", h.PutPosition)
	http.HandleFunc("/price", h.PutPrice)
	http.HandleFunc("/yieldCurve", h.PutYieldCurve)
	http.HandleFunc("/matrix", h.PutMatrix)
	http.ListenAndServe(":8080", nil)
}
