package handlers

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var authInput models.AuthInput
	if err := c.ShouldBindBodyWithJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verifiedUser, err := services.VerifyUser(&authInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := services.GenerateJWT(verifiedUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func SignUp(c *gin.Context) {
	var authInput models.AuthInput
	if err := c.ShouldBindBodyWithJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateUser(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
	})
}

func GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	userID := c.MustGet("user_id").(int64)
	if userID != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := services.FetchUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func PutUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	currUserID := c.MustGet("user_id").(int64)
	currUserRoles, err := services.GetRoleByUserID(currUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if !containsRole(currUserRoles, models.Admin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can update roles"})
		return
	}

	var req struct {
		Roles []models.RoleName `json:"roles" binding:"required"`
	}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roles field is necessary"})
		return
	}

	for _, role := range req.Roles {
		if !role.Valid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "roles field must contain only valid roles (admin, staff)"})
			return
		}
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("user with ID %d does not exist", userID)})
		return
	}

	if err := services.PutUserRoles(userID, req.Roles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user roles updated successfully"})
}

func containsRole(roles []models.Role, target models.RoleName) bool {
	for _, role := range roles {
		if role.Name == target {
			return true
		}
	}
	return false
}
