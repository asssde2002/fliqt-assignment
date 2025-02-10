package handlers_test

import (
	"backend/internal/db"
	"backend/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClockIn(t *testing.T) {
	setup()
	authInput := models.AuthInput{Username: "testuser", Password: "testpassword"}
	user := createUser(authInput, false)

	writer := makeRequest("POST", "/clock/in", nil, &authInput)
	assert.Equal(t, http.StatusCreated, writer.Code)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, response["message"], "user clocked in for today successfully")

	var punchCard models.PunchCard
	today := time.Now().UTC().Format("2006-01-02")
	err := db.DB.Where("user_id = ? AND DATE(created_at) = ?", user.ID, today).First(&punchCard).Error
	assert.Equal(t, err, nil)
	assert.False(t, punchCard.CreatedAt.IsZero())
	assert.True(t, punchCard.ClockIn.Valid)
	assert.False(t, punchCard.ClockOut.Valid)

	writer = makeRequest("POST", "/clock/in", nil, &authInput)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, response["error"], "user has already clocked in for today")
}

func TestClockOut(t *testing.T) {
	setup()

	authInput := models.AuthInput{Username: "testuser", Password: "testpassword"}
	user := createUser(authInput, false)

	writer := makeRequest("POST", "/clock/out", nil, &authInput)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, response["error"], "user needs to clock in for today first or already clocked out")

	writer = makeRequest("POST", "/clock/in", nil, &authInput)
	assert.Equal(t, http.StatusCreated, writer.Code)
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, response["message"], "user clocked in for today successfully")

	var punchCard models.PunchCard
	today := time.Now().UTC().Format("2006-01-02")
	err := db.DB.Where("user_id = ? AND DATE(created_at) = ?", user.ID, today).First(&punchCard).Error
	assert.Equal(t, err, nil)
	assert.False(t, punchCard.CreatedAt.IsZero())
	assert.True(t, punchCard.ClockIn.Valid)
	assert.False(t, punchCard.ClockOut.Valid)

	writer = makeRequest("POST", "/clock/out", nil, &authInput)
	assert.Equal(t, http.StatusCreated, writer.Code)
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, response["message"], "user clocked out for today successfully")

	err = db.DB.Where("user_id = ? AND DATE(created_at) = ?", user.ID, today).First(&punchCard).Error
	assert.Equal(t, err, nil)
	assert.False(t, punchCard.CreatedAt.IsZero())
	assert.True(t, punchCard.ClockIn.Valid)
	assert.True(t, punchCard.ClockOut.Valid)
	assert.True(t, punchCard.ClockIn.Time.Before(punchCard.ClockOut.Time))
}

func TestGetPunchCard(t *testing.T) {
	setup()
	authInput := models.AuthInput{Username: "testuser", Password: "testpassword"}
	user := createUser(authInput, false)
	currTime := time.Now().UTC().Truncate(time.Millisecond)
	punchCards := []models.PunchCard{
		{
			UserID:    user.ID,
			CreatedAt: currTime,
			ClockIn:   sql.NullTime{Time: currTime, Valid: true},
			ClockOut:  sql.NullTime{Time: currTime.Add(time.Hour), Valid: true},
		},
		{
			UserID:    user.ID,
			CreatedAt: currTime.Add(24 * time.Hour),
			ClockIn:   sql.NullTime{Time: currTime.Add(24 * time.Hour), Valid: true},
			ClockOut:  sql.NullTime{Time: currTime.Add(25 * time.Hour), Valid: true},
		},
	}
	err := db.DB.Create(&punchCards).Error
	assert.Equal(t, err, nil)

	writer := makeRequest("GET", "/clock", nil, &authInput)
	assert.Equal(t, http.StatusOK, writer.Code)
	var response []models.PunchCardResponse
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.True(t, response[0].ClockIn.Equal(punchCards[1].ClockIn.Time))
	assert.True(t, response[0].ClockOut.Equal(punchCards[1].ClockOut.Time))
	assert.True(t, response[1].ClockIn.Equal(punchCards[0].ClockIn.Time))
	assert.True(t, response[1].ClockOut.Equal(punchCards[0].ClockOut.Time))
}
