package handlers_test

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/routes"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

var roleMapCache = make(map[models.RoleName]int64)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	config.LoadConfig("../../.env.test")
	exitCode := m.Run()
	os.Exit(exitCode)
}

func router() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	routes.RegisterUserRoutes(router)
	routes.RegisterPunchCardRoutes(router)
	return router
}

func setup() {
	db.InitDB()
	teardown()
	db.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.UserRole{}, &models.PunchCard{})
	ensureRolesExist()
}

func teardown() {
	migrator := db.DB.Migrator()
	migrator.DropTable(&models.User{}, &models.Role{}, &models.UserRole{}, &models.PunchCard{})
}

func makeRequest(method, url string, body interface{}, authInput *models.AuthInput) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if authInput != nil {
		request.Header.Add("Authorization", "Bearer "+bearerToken(authInput))
	}
	writer := httptest.NewRecorder()
	router().ServeHTTP(writer, request)
	return writer
}

func bearerToken(authInput *models.AuthInput) string {
	writer := makeRequest("POST", "/auth/login", authInput, nil)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	return response["token"]
}

func ensureRolesExist() {
	roles := []models.Role{
		{Name: models.Staff},
		{Name: models.Admin},
	}

	err := db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error
	if err != nil {
		log.Printf("Failed to create roles: %v", err)
	}

	for _, role := range roles {
		roleMapCache[role.Name] = role.ID
	}
}

func createUser(authInput models.AuthInput, isAdmin bool) *models.User {
	var user models.User
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to create hash password: %v", err)
	}

	user = models.User{Username: authInput.Username, Password: string(hashedPassword)}
	err = db.DB.Create(&user).Error
	if err != nil {
		log.Printf("Failed to create users: %v", err)
	}
	var userRoles []models.UserRole
	userRoles = append(userRoles, models.UserRole{UserID: user.ID, RoleID: roleMapCache[models.Staff]})
	if isAdmin {
		userRoles = append(userRoles, models.UserRole{UserID: user.ID, RoleID: roleMapCache[models.Admin]})
	}

	err = db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRoles).Error
	if err != nil {
		log.Printf("Failed to create userroles: %v", err)
	}

	return &user
}
