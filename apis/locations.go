package apis

import (
	"fmt"
	"net/http"
	"strconv"

	database "tourist-api/db"
	"tourist-api/utils"

	"github.com/gin-gonic/gin"
)

func bindLocationToTrip(c *gin.Context) {
	var createLocationPayload struct {
		PlaceID       string  `json:"placeId" binding:"required"`
		LocationName  string  `json:"locationName" binding:"required"`
		Longitude     float64 `json:"longitude" binding:"required"`
		Latitude      float64 `json:"latitude" binding:"required"`
		Rating        float64 `json:"rating"`
		CoverImageUrl string  `json:"coverImageUrl"`
		TotalReviews  int     `json:"totalReviews"`
	}

	if err := c.BindJSON(&createLocationPayload); err != nil {
		fmt.Printf("%+v\n", createLocationPayload.CoverImageUrl)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	tripID, err := utils.StringToUnit(c.Param("tripID"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid tripID format"})
		return
	}

	userID, err := c.Cookie(utils.GLOBAL_TOKEN_NAMING)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		return
	}

	userIDAsUint, err := utils.StringToUnit(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorized failed"})
		return
	}

	var location = database.Location{}
	location, err = database.GetLocationByPlaceID(createLocationPayload.PlaceID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server side error: 0004"})
		return
	}

	if location.ID == 0 {
		longitude, _ := strconv.ParseFloat(fmt.Sprintf("%.7f", createLocationPayload.Longitude), 64)
		latitude, _ := strconv.ParseFloat(fmt.Sprintf("%.7f", createLocationPayload.Latitude), 64)
		rating, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", createLocationPayload.Rating), 64)

		newLocation := database.Location{
			PlaceID:       createLocationPayload.PlaceID,
			LocationName:  createLocationPayload.LocationName,
			Longitude:     longitude,
			Latitude:      latitude,
			Rating:        rating,
			CoverImageUrl: createLocationPayload.CoverImageUrl,
			TotalReviews:  createLocationPayload.TotalReviews,
		}

		location, err = database.CreateLocation(newLocation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server side error: 0005"})
			return
		}
	}

	err = database.BindLocationToTrip(userIDAsUint, tripID, location)
	switch err.(type) {
	case *utils.NotFoundError:
		c.JSON(http.StatusNotFound, gin.H{"error": "tripId cannot found trip"})
		return
	case *utils.ConflictError:
		c.JSON(http.StatusConflict, gin.H{"error": "location already bind to into trip"})
		return
	default:
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server side error: 0006"})
			return
		}
	}

	type AddLocationToTripResponse struct {
		Success bool
	}
	c.JSON(http.StatusOK, AddLocationToTripResponse{Success: true})
}

func unbindLocationToTrip(c *gin.Context) {
	tripID, err := utils.StringToUnit(c.Param("tripID"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid tripID format"})
		return
	}

	locationID, err := utils.StringToUnit(c.Param("locationID"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid locationID format"})
		return
	}

	userID, err := c.Cookie(utils.GLOBAL_TOKEN_NAMING)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
		return
	}

	userIDAsUint, err := utils.StringToUnit(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorized failed"})
		return
	}

	err = database.UnbindLocationFromTrip(userIDAsUint, tripID, locationID)
	switch err.(type) {
	case *utils.NotFoundError:
		c.JSON(http.StatusNotFound, gin.H{"error": "tripId cannot found trip"})
		return
	case *utils.ConflictError:
		c.JSON(http.StatusConflict, gin.H{"error": "location not bind into trip before"})
		return
	default:
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server side error: 0007"})
			return
		}
	}

	type AddLocationToTripResponse struct {
		Success bool
	}
	c.JSON(http.StatusOK, AddLocationToTripResponse{Success: true})
}
