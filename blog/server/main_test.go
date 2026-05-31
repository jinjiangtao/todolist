package main

import (
	"blog-server/database"
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
	"strconv"
	"testing"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.SetupRoutes(r)
	return r
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Remove("blog.db")
	database.InitDB()
	code := m.Run()
	os.Remove("blog.db")
	os.Exit(code)
}

func TestLoginSuccess(t *testing.T) {
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
	assert.NotEmpty(t, response["data"])
	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
}

func TestLoginWrongPassword(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	reqBody := `{"username":"admin","password":"wrongpassword"}`
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(401), response["code"])
}

func TestLoginWrongUsername(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	reqBody := `{"username":"wronguser","password":"123456"}`
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginMissingFields(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	reqBody := `{"username":""}`
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPublishedArticlesEmpty(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/articles", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(200), response["code"])
	data := response["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	assert.Equal(t, 0, len(list))
}

func TestGetArticleNotFound(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/articles/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdminArticlesWithoutToken(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/articles", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminArticlesWithToken(t *testing.T) {
	r := setupTestRouter()

	token, _ := utils.GenerateToken(1, "admin")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/articles", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateArticle(t *testing.T) {
	r := setupTestRouter()

	token, _ := utils.GenerateToken(1, "admin")

	w := httptest.NewRecorder()
	reqBody := `{"title":"Test Article","content":"This is test content","status":1}`
	req, _ := http.NewRequest("POST", "/api/v1/admin/articles", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(200), response["code"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Test Article", data["title"])
}

func TestCreateArticleMissingTitle(t *testing.T) {
	r := setupTestRouter()

	token, _ := utils.GenerateToken(1, "admin")

	w := httptest.NewRecorder()
	reqBody := `{"content":"This is test content","status":1}`
	req, _ := http.NewRequest("POST", "/api/v1/admin/articles", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateArticle(t *testing.T) {
	r := setupTestRouter()

	token, _ := utils.GenerateToken(1, "admin")

	article := models.Article{Title: "Original Title", Content: "Original Content", Status: 0}
	database.DB.Create(&article)

	w := httptest.NewRecorder()
	reqBody := `{"title":"Updated Title","content":"Updated Content","status":1}`
	req, _ := http.NewRequest("PUT", "/api/v1/admin/articles/"+strconv.Itoa(int(article.ID)), bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Updated Title", data["title"])
}

func TestDeleteArticle(t *testing.T) {
	r := setupTestRouter()

	token, _ := utils.GenerateToken(1, "admin")

	article := models.Article{Title: "To Delete", Content: "Content", Status: 0}
	database.DB.Create(&article)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/articles/"+strconv.Itoa(int(article.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var deletedArticle models.Article
	result := database.DB.First(&deletedArticle, article.ID)
	assert.Error(t, result.Error)
}

func TestJWTGenerateAndParse(t *testing.T) {
	token, err := utils.GenerateToken(1, "admin")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := utils.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "admin", claims.Username)
}

func TestJWTParseInvalidToken(t *testing.T) {
	_, err := utils.ParseToken("invalid-token")
	assert.Error(t, err)
}

func TestPasswordHashing(t *testing.T) {
	password := "test123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	assert.NotEqual(t, password, string(hashed))

	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	assert.NoError(t, err)

	err = bcrypt.CompareHashAndPassword(hashed, []byte("wrongpassword"))
	assert.Error(t, err)
}

func TestChangePassword(t *testing.T) {
	r := setupTestRouter()

	token, _ := utils.GenerateToken(1, "admin")

	w := httptest.NewRecorder()
	reqBody := `{"old_password":"123456","new_password":"newpass123"}`
	req, _ := http.NewRequest("PUT", "/api/v1/admin/password", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChangePasswordWrongOld(t *testing.T) {
	r := setupTestRouter()

	token, _ := utils.GenerateToken(1, "admin")

	w := httptest.NewRecorder()
	reqBody := `{"old_password":"wrongpassword","new_password":"newpass123"}`
	req, _ := http.NewRequest("PUT", "/api/v1/admin/password", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestArticleViewCount(t *testing.T) {
	r := setupTestRouter()

	article := models.Article{Title: "Test View", Content: "Content", Status: 1, ViewCount: 0}
	database.DB.Create(&article)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/articles/"+strconv.Itoa(int(article.ID)), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedArticle models.Article
	database.DB.First(&updatedArticle, article.ID)
	assert.Equal(t, 1, updatedArticle.ViewCount)
}

func TestGetPublishedArticlesWithPublishedArticle(t *testing.T) {
	r := setupTestRouter()

	article := models.Article{Title: "Published Article", Content: "Content", Status: 1}
	database.DB.Create(&article)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/articles", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	assert.GreaterOrEqual(t, len(list), 1)
}

func TestGetPublishedArticlesNotIncludeDraft(t *testing.T) {
	r := setupTestRouter()

	draft := models.Article{Title: "Draft Article", Content: "Content", Status: 0}
	database.DB.Create(&draft)

	published := models.Article{Title: "Published Article 2", Content: "Content", Status: 1}
	database.DB.Create(&published)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/articles", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	list := data["list"].([]interface{})

	for _, item := range list {
		article := item.(map[string]interface{})
		assert.Equal(t, float64(1), article["status"])
	}
}

func TestAdminArticlesIncludeDraft(t *testing.T) {
	r := setupTestRouter()

	token, _ := utils.GenerateToken(1, "admin")

	draft := models.Article{Title: "Draft for Admin", Content: "Content", Status: 0}
	database.DB.Create(&draft)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/articles", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	assert.GreaterOrEqual(t, len(list), 1)
}

func TestPagination(t *testing.T) {
	r := setupTestRouter()

	var countBefore int64
	database.DB.Model(&models.Article{}).Where("status = 1").Count(&countBefore)

	for i := 1; i <= 15; i++ {
		article := models.Article{Title: "Pagination Article " + strconv.Itoa(i), Content: "Content", Status: 1}
		database.DB.Create(&article)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/articles?page=1&page_size=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	assert.Equal(t, 10, len(list))
	assert.GreaterOrEqual(t, data["total"], float64(countBefore+15))
}

func TestUpdateArticleNotFound(t *testing.T) {
	r := setupTestRouter()

	token, _ := utils.GenerateToken(1, "admin")

	w := httptest.NewRecorder()
	reqBody := `{"title":"Updated Title","content":"Updated Content","status":1}`
	req, _ := http.NewRequest("PUT", "/api/v1/admin/articles/9999", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteArticleNotFound(t *testing.T) {
	r := setupTestRouter()

	token, _ := utils.GenerateToken(1, "admin")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/admin/articles/9999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}