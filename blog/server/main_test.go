package main

import (
	"blog-server/controllers"
	"blog-server/database"
	"blog-server/middleware"
	"blog-server/models"
	"blog-server/routes"
	"blog-server/utils"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testToken string

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.SetupRoutes(r)
	return r
}

func setupTestDB() {
	os.Remove("test_blog.db")
	database.DB, _ = database.DB.DB()
	var err error
	database.DB, err = database.DB.DB()
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Remove("test_blog.db")
	os.Exit(code)
}

func TestLogin(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	reqBody := `{"username":"admin","password":"123456"}`
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(200), response["code"])
}

func TestGetPublishedArticles(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/articles", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWT(t *testing.T) {
	token, err := utils.GenerateToken(1, "admin")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := utils.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "admin", claims.Username)
}

func TestPasswordHashing(t *testing.T) {
	password := "test123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	assert.NotEqual(t, password, string(hashed))

	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	assert.NoError(t, err)
}
