// routes/routes.go
package routes

import (
	"net/http"

	"golang-test/controllers"
	"golang-test/model"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/analyze", func(c *gin.Context) {
		var contracts []model.OptionsContract

		if err := c.ShouldBindJSON(&contracts); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		xyValues := controllers.CalculateXYValues(contracts)
		maxProfit := controllers.CalculateMaxProfit(contracts)
		maxLoss := controllers.CalculateMaxLoss(contracts)
		breakEvenPoints := controllers.CalculateBreakEvenPoints(contracts)

		response := controllers.AnalysisResponse{
			XYValues:        xyValues,
			MaxProfit:       maxProfit,
			MaxLoss:         maxLoss,
			BreakEvenPoints: breakEvenPoints,
		}

		c.JSON(http.StatusOK, response)
	})

	return router
}
