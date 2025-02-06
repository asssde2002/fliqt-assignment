package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ClockIn(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	today := time.Now().UTC().Format("2006-01-02")
	_, exists := services.CheckPunchCardExist(today, userID)
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user has already clocked in for today"})
		return
	}

	if err := services.CreatePunchCard(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user clocked in for today failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user clocked in for today successfully"})

}

func ClockOut(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	today := time.Now().UTC().Format("2006-01-02")
	punchCard, exists := services.CheckPunchCardExist(today, userID)
	if !exists || punchCard.ClockOut.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user needs to clock in for today first or already clocked out"})
		return
	}
	if err := services.UpdatePunchCard(punchCard.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user clocked out for today failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user clocked out for today successfully"})

}

func GetPunchCard(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	punchCards, err := services.GetPunchCard(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get user's punch cards"})
		return
	}

	var responses []models.PunchCardResponse
	for _, pc := range punchCards {
		var clockIn, clockOut *time.Time
		if pc.ClockIn.Valid {
			clockIn = &pc.ClockIn.Time
		}
		if pc.ClockOut.Valid {
			clockOut = &pc.ClockOut.Time
		}
		pcr := models.PunchCardResponse{
			ClockIn:   clockIn,
			ClockOut:  clockOut,
			CreatedAt: pc.CreatedAt,
		}
		responses = append(responses, pcr)
	}

	c.JSON(http.StatusOK, responses)
}
