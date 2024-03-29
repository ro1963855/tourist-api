package apis

import (
	"net/http"
	"net/url"
	"strconv"

	database "tourist-api/db"
	"tourist-api/utils"

	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	var loginPayload struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.BindJSON(&loginPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := database.AuthenticateUser(loginPayload.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	parsedURL, err := url.Parse(utils.ENV_FRONTEND_DOMAIN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server side error: 0001"})
		return
	}

	type LoginResponse struct {
		Success bool
	}

	hostname := parsedURL.Hostname()
	id := strconv.FormatInt(int64(user.ID), 10)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(utils.GLOBAL_TOKEN_NAMING, id, 2147483647, "/", hostname, true, true)
	c.JSON(http.StatusOK, LoginResponse{Success: true})
}
