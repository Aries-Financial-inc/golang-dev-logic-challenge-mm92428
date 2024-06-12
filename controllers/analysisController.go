// controllers/analysisController.go
package controllers

import (
	"encoding/json"
	"golang-test/model"
	"math"
	"net/http"
)

// AnalysisResponse represents the data structure of the analysis result
type AnalysisResponse struct {
	XYValues        []XYValue `json:"xy_values"`
	MaxProfit       float64   `json:"max_profit"`
	MaxLoss         float64   `json:"max_loss"`
	BreakEvenPoints []float64 `json:"break_even_points"`
}

// XYValue represents a pair of X and Y values
type XYValue struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func AnalysisHandler(w http.ResponseWriter, r *http.Request) {
	var contracts []model.OptionsContract

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&contracts); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate analysis
	xyValues := CalculateXYValues(contracts)
	maxProfit := CalculateMaxProfit(contracts)
	maxLoss := CalculateMaxLoss(contracts)
	breakEvenPoints := CalculateBreakEvenPoints(contracts)

	response := AnalysisResponse{
		XYValues:        xyValues,
		MaxProfit:       maxProfit,
		MaxLoss:         maxLoss,
		BreakEvenPoints: breakEvenPoints,
	}

	// Encode response to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CalculateXYValues(contracts []model.OptionsContract) []XYValue {
	// Calculate XY values for risk & reward graph
	var xyValues []XYValue
	minPrice, maxPrice := GetMinMaxStrikePrices(contracts)

	for price := minPrice - 20; price <= maxPrice+20; price += 1 {
		profitLoss := 0.0
		for _, contract := range contracts {
			if contract.Type == "call" {
				if contract.LongShort == "long" {
					profitLoss += math.Max(0, price-contract.StrikePrice) - contract.Ask
				} else {
					profitLoss += contract.Bid - math.Max(0, price-contract.StrikePrice)
				}
			} else {
				if contract.LongShort == "long" {
					profitLoss += math.Max(0, contract.StrikePrice-price) - contract.Ask
				} else {
					profitLoss += contract.Bid - math.Max(0, contract.StrikePrice-price)
				}
			}
		}
		xyValues = append(xyValues, XYValue{X: price, Y: profitLoss})
	}
	return xyValues
}

func CalculateMaxProfit(contracts []model.OptionsContract) float64 {
	// Calculate maximum possible profit
	maxProfit := math.Inf(-1)
	for _, xyValue := range CalculateXYValues(contracts) {
		if xyValue.Y > maxProfit {
			maxProfit = xyValue.Y
		}
	}
	return maxProfit
}

func CalculateMaxLoss(contracts []model.OptionsContract) float64 {
	// Calculate maximum possible loss
	maxLoss := math.Inf(1)
	for _, xyValue := range CalculateXYValues(contracts) {
		if xyValue.Y < maxLoss {
			maxLoss = xyValue.Y
		}
	}
	return maxLoss
}

func CalculateBreakEvenPoints(contracts []model.OptionsContract) []float64 {
	// Calculate break even points
	var breakEvenPoints []float64
	xyValues := CalculateXYValues(contracts)
	for i := 1; i < len(xyValues); i++ {
		if (xyValues[i-1].Y < 0 && xyValues[i].Y > 0) || (xyValues[i-1].Y > 0 && xyValues[i].Y < 0) {
			breakEvenPoints = append(breakEvenPoints, xyValues[i].X)
		}
	}
	return breakEvenPoints
}

func GetMinMaxStrikePrices(contracts []model.OptionsContract) (float64, float64) {
	minPrice := math.Inf(1)
	maxPrice := math.Inf(-1)
	for _, contract := range contracts {
		if contract.StrikePrice < minPrice {
			minPrice = contract.StrikePrice
		}
		if contract.StrikePrice > maxPrice {
			maxPrice = contract.StrikePrice
		}
	}
	return minPrice, maxPrice
}
