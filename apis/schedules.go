package apis

import (
	"net/http"
	"strconv"

	database "tourist-api/db"
	"tourist-api/utils"

	"github.com/gin-gonic/gin"
)

func createSchedule(c *gin.Context) {
	var createSchedulePayload struct {
		ScheduleName string `json:"scheduleName" binding:"required"`
	}

	if err := c.BindJSON(&createSchedulePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userID, err := c.Cookie(utils.GLOBAL_TOKEN_NAMING)
	userIDAsUint, err1 := strconv.ParseUint(userID, 10, 64)
	if err != nil || err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized failed"})
		return
	}

	newUserSchedule, err := database.CreateSchedule(uint(userIDAsUint), createSchedulePayload.ScheduleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server side error: 0002"})
		return
	}

	type CreateScheduleResponse struct {
		ScheduleID uint
	}

	c.JSON(http.StatusOK, CreateScheduleResponse{ScheduleID: newUserSchedule.ID})
}
