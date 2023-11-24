package apis

import (
	"net/http"
	"time"

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
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		return
	}

	userIDAsUint, err := utils.StringToUnit(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized failed"})
		return
	}

	newTrip, err := database.CreateTrip(userIDAsUint, createTripPayload.TripName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server side error: 0002"})
		return
	}

	type CreateTripResponse struct {
		TripID uint
	}

	c.JSON(http.StatusOK, CreateTripResponse{TripID: newTrip.ID})
}

func getTrips(c *gin.Context) {
	userID, err := c.Cookie(utils.GLOBAL_TOKEN_NAMING)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		return
	}

	userIDAsUint, err := utils.StringToUnit(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized failed"})
		return
	}

	trips, err := database.GetTrips(userIDAsUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server side error: 0003"})
		return
	}

	type GetTripsResponse struct {
		ID        uint
		TripName  string
		CreatedAt time.Time
	}

	var responseTrips []GetTripsResponse
	for _, trip := range trips {
		responseTrip := GetTripsResponse{
			ID:        trip.ID,
			TripName:  trip.TripName,
			CreatedAt: trip.CreatedAt,
		}

		responseTrips = append(responseTrips, responseTrip)
	}

	c.JSON(http.StatusOK, responseTrips)
}

func getTripDetail(c *gin.Context) {
	userID, err := c.Cookie(utils.GLOBAL_TOKEN_NAMING)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		return
	}

	userIDAsUint, err := utils.StringToUnit(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized failed"})
		return
	}

	tripID, err := utils.StringToUnit(c.Param("tripID"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid tripID format"})
		return
	}

	tripDetail, err := database.GetTripDetail(userIDAsUint, tripID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trip detail not found"})
		return
	}

	c.JSON(http.StatusOK, tripDetail)
}
