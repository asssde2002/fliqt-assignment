package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if user_id, exists := c.Get("user_id"); !exists || user_id != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	reqUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	currUserID, ok := reqUserID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	currUserRoles, err := services.GetRoleByUserID(currUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if !containsRole(currUserRoles, models.Admin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can update roles"})
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

	if err := services.PutUserRoles(userID, req.Roles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User roles updated successfully"})
}

func containsRole(roles []models.Role, target models.RoleName) bool {
	for _, role := range roles {
		if role.Name == target {
			return true
		}
	}
	return false
}
