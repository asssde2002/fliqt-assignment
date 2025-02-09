package handlers_test

import (
	"backend/internal/db"
	"backend/internal/handlers"
	"backend/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getMockDB() (sqlmock.Sqlmock, func(), error) {
	dbMock, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      dbMock,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	originalDB := db.DB
	db.DB = gormDB

	restoreDB := func() {
		db.DB = originalDB
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()
	}

	return mock, restoreDB, nil
}

func setupRolesMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`SELECT (.+) FROM roles WHERE name = 'staff'`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(2, "staff"))
}

func TestLogin(t *testing.T) {
	// mock, restoreDB, err := getMockDB()
	// if err != nil {
	// 	t.Fatalf("failed to create mock db: %v", err)
	// }
	// defer restoreDB()

	// gin.SetMode(gin.TestMode)
	// w := httptest.NewRecorder()
	// c, _ := gin.CreateTestContext(w)

	// authInput := models.AuthInput{
	// 	Username: "testuser",
	// 	Password: "testpassword",
	// }
	// body, _ := json.Marshal(authInput)
	// c.Request = httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	// c.Request.Header.Set("Content-Type", "application/json")
	// handlers.SignUp(c)

	// assert.Equal(t, http.StatusCreated, w.Code)
	// assert.Contains(t, w.Body.String(), "user created successfully")

	// if err := mock.ExpectationsWereMet(); err != nil {
	// 	t.Errorf("未滿足的 SQL mock 條件: %v", err)
	// }
}

func TestSignUp(t *testing.T) {
	mock, restoreDB, err := getMockDB()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer restoreDB()

	setupRolesMock(mock)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	authInput := models.AuthInput{
		Username: "testuser",
		Password: "testpassword",
	}
	body, _ := json.Marshal(authInput)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	handlers.SignUp(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "user created successfully")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未滿足的 SQL mock 條件: %v", err)
	}
}

func TestGetUser(t *testing.T) {

}

func TestPutUserRoles(t *testing.T) {

}
