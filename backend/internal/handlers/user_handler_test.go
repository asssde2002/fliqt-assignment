package handlers_test

import (
	"backend/internal/db"
	"backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	setup()
	username := "testuser"
	password := "testpassword"
	authInput := models.AuthInput{Username: "testuser", Password: "testpassword"}

	user := createUser(authInput, false)

	err := db.DB.Where("username = ?", user.Username).First(&user).Error
	assert.NoError(t, err, "user should exist in the database")

	writer := makeRequest("POST", "/auth/login", models.AuthInput{Username: username, Password: password}, nil)
	assert.Equal(t, http.StatusOK, writer.Code)

	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NotEmpty(t, response["token"], "Token should not be empty")
}

func TestSignUp(t *testing.T) {
	setup()
	newUser := models.AuthInput{
		Username: "testuser",
		Password: "testpassword",
	}

	var user models.User
	err := db.DB.Where("username = ?", newUser.Username).First(&user).Error
	assert.Error(t, err, "user should not exist in the database")

	writer := makeRequest("POST", "/auth/signup", newUser, nil)
	assert.Equal(t, http.StatusCreated, writer.Code)

	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, response["message"], "user created successfully")

	err = db.DB.Where("username = ?", newUser.Username).First(&user).Error
	assert.NoError(t, err, "user should exist in the database")
	assert.Equal(t, newUser.Username, user.Username, "stored username should match")
}

func TestGetUser(t *testing.T) {
	setup()
	authInput := models.AuthInput{Username: "testuser", Password: "testpassword"}
	user1 := createUser(authInput, true)

	writer := makeRequest("GET", fmt.Sprintf("/user/%d", user1.ID), nil, &authInput)
	assert.Equal(t, http.StatusOK, writer.Code)
	var response models.UserResponse
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, response.ID, user1.ID)
	assert.Equal(t, response.Username, user1.Username)
	assert.Equal(t, response.IsActive, true)
	assert.Equal(t, response.Roles, []models.RoleName{models.Admin, models.Staff})

	authInput2 := models.AuthInput{Username: "testuser2", Password: "testpassword2"}
	user2 := createUser(authInput2, false)

	writer2 := makeRequest("GET", fmt.Sprintf("/user/%d", user1.ID), nil, &authInput2)
	assert.Equal(t, http.StatusForbidden, writer2.Code)

	writer3 := makeRequest("GET", fmt.Sprintf("/user/%d", user2.ID), nil, &authInput2)
	assert.Equal(t, http.StatusOK, writer3.Code)
	var response2 models.UserResponse
	json.Unmarshal(writer3.Body.Bytes(), &response2)
	assert.Equal(t, response2.ID, user2.ID)
	assert.Equal(t, response2.Username, user2.Username)
	assert.Equal(t, response2.IsActive, true)
	assert.Equal(t, response2.Roles, []models.RoleName{models.Staff})
}

func TestPutUserRoles(t *testing.T) {
	setup()
	authInput := models.AuthInput{Username: "testuser", Password: "testpassword"}
	user1 := createUser(authInput, true)

	authInput2 := models.AuthInput{Username: "testuser2", Password: "testpassword2"}
	user2 := createUser(authInput2, false)

	roles := struct {
		Roles []models.RoleName `json:"roles" binding:"required"`
	}{
		Roles: []models.RoleName{models.Admin, models.Staff},
	}
	writer := makeRequest("PUT", fmt.Sprintf("/user/%d/roles", user1.ID), roles, &authInput2)
	assert.Equal(t, http.StatusForbidden, writer.Code)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, response["error"], "only admin can update roles")

	db.DB.Model(&user2).Preload("Roles").First(&user2)
	roleNames := []models.RoleName{}
	for _, role := range user2.Roles {
		roleNames = append(roleNames, role.Name)
	}
	assert.NotContains(t, roleNames, models.Admin)
	assert.Contains(t, roleNames, models.Staff)

	writer2 := makeRequest("PUT", fmt.Sprintf("/user/%d/roles", user2.ID), roles, &authInput)
	assert.Equal(t, http.StatusForbidden, writer.Code)
	var response2 map[string]string
	json.Unmarshal(writer2.Body.Bytes(), &response2)
	assert.Equal(t, response2["message"], "user roles updated successfully")

	db.DB.Model(&user2).Preload("Roles").First(&user2)
	roleNames = []models.RoleName{}
	for _, role := range user2.Roles {
		roleNames = append(roleNames, role.Name)
	}
	assert.Contains(t, roleNames, models.Admin)
	assert.Contains(t, roleNames, models.Staff)
}
