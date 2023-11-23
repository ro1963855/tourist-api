package apis

import (
	"net/http"
	"strconv"

	database "tourist-api/db"
	"tourist-api/utils"

	"github.com/gin-gonic/gin"
)

func createTrip(c *gin.Context) {
	var createTripPayload struct {
		TripName string `json:"tripName" binding:"required"`
	}

	if err := c.BindJSON(&createTripPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userID, err := c.Cookie(utils.GLOBAL_TOKEN_NAMING)
	userIDAsUint, err1 := strconv.ParseUint(userID, 10, 64)
	if err != nil || err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized failed"})
		return
	}

	newTrip, err := database.CreateTrip(uint(userIDAsUint), createTripPayload.TripName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server side error: 0002"})
		return
	}

	type CreateTripResponse struct {
		TripID uint
	}

	c.JSON(http.StatusOK, CreateTripResponse{TripID: newTrip.ID})
}
