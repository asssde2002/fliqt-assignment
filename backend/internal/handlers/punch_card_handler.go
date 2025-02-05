package handlers

import (
	"backend/internal/services"
	"fmt"
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
	fmt.Println(punchCard, exists)
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

}
