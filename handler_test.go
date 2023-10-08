package main_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"demerzel-badges/internal/handlers"
	"demerzel-badges/internal/models"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) FindSkillById(skillID uint) (*models.Skill, error) {
	args := m.Called(skillID)
	return args.Get(0).(*models.Skill), args.Error(1)
}

func (m *MockDatabase) BadgeExists(skillID uint, name models.Badge) bool {
	args := m.Called(skillID, name)
	return args.Bool(0)
}

func (m *MockDatabase) CreateBadge(badge models.SkillBadge) (*models.SkillBadge, error) {
	args := m.Called(badge)
	return args.Get(0).(*models.SkillBadge), args.Error(1)
}

var mockDB = new(MockDatabase)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/api/badges", handlers.CreateBadgeHandler(mockDB))
	return router
}

var testPayload = []byte(`{"name": "beginner", "min_score": 0, "max_score": 50}`)

func TestCreateBadgeEndpoint(t *testing.T) {
	router := setupRouter()
	defer TestCleanup(t)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/api/badges", bytes.NewBuffer(testPayload))
	
    c, _ := gin.CreateTestContext(w)
    c.Request = req

	router.ServeHTTP(w, req)

	handlers.CreateBadgeHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateBadgeHandler_Success(t *testing.T) {
	mockDB := new(MockDatabase)

	router := setupRouter()

	req, err := http.NewRequest("POST", "/api/badges", bytes.NewBuffer(testPayload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockDB.On("FindSkillById", uint(1)).Return(&models.Skill{ID: 1}, nil)
	mockDB.On("BadgeExists", uint(1), models.Badge("beginner")).Return(false)
	mockDB.On("CreateBadge", mock.AnythingOfType("*models.SkillBadge")).Return(nil, nil)

	handlers.CreateBadgeHandler(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	assert.JSONEq(t, `{"status": "success", "message": "Badge created successfully"}`, w.Body.String())

	mockDB.AssertExpectations(t)
}

func TestCreateBadgeHandler_CreateBadgeError(t *testing.T) {
	mockDB := new(MockDatabase)

	router := setupRouter()

	req, err := http.NewRequest("POST", "/api/badges", bytes.NewBuffer(testPayload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	mockDB.On("FindSkillById", uint(1)).Return(&models.Skill{ID: 1}, nil)
	mockDB.On("BadgeExists", uint(1), models.Badge("beginner")).Return(false)
	mockDB.On("CreateBadge", mock.AnythingOfType("*models.SkillBadge")).Return(nil, errors.New("database error"))

	handlers.CreateBadgeHandler(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	assert.Contains(t, w.Body.String(), "Unable to create badge")
}

func TestCreateBadgeHandler_InvalidInput(t *testing.T) {
	mockDB := new(MockDatabase)

	router := setupRouter()

	req, err := http.NewRequest("POST", "/api/badges", bytes.NewBuffer(testPayload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	mockDB.On("FindSkillById", uint(1)).Return(&models.Skill{ID: 1}, nil)
	mockDB.On("BadgeExists", uint(1), models.Badge("beginner")).Return(false)
	mockDB.On("CreateBadge", mock.AnythingOfType("*models.SkillBadge")).Return(nil, errors.New("Error! Invalid Input"))

	handlers.CreateBadgeHandler(c)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Body.String(), "min_score should be at least 0")
}

func TestCreateBadgeHandler_DatabaseError(t *testing.T) {
	mockDB := new(MockDatabase)

	router := setupRouter()

	req, err := http.NewRequest("POST", "/api/badges", bytes.NewBuffer(testPayload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	

	mockDB.On("FindSkillById", uint(1)).Return(nil, errors.New("database error"))
	mockDB.On("BadgeExists", uint(1), models.Badge("beginner")).Return(false)
	mockDB.On("CreateBadge", mock.AnythingOfType("*models.SkillBadge")).Return(nil, errors.New("database error"))

	handlers.CreateBadgeHandler(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}

func TestCreateBadgeHandler_FindSkillByIdError(t *testing.T) {
	mockDB := new(MockDatabase)

	router := setupRouter()
	
	req, err := http.NewRequest("POST", "/api/badges", bytes.NewBuffer(testPayload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	mockDB.On("FindSkillById", uint(1)).Return(nil, nil)
	mockDB.On("BadgeExists", uint(1), models.Badge("beginner")).Return(false)
	mockDB.On("CreateBadge", mock.AnythingOfType("*models.SkillBadge")).Return(nil, nil)

	handlers.CreateBadgeHandler(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	assert.Contains(t, w.Body.String(), "No skill found matching provided ID")
}

func TestCleanup(t *testing.T) {
	defer mockDB.AssertExpectations(t)
	mockDB.Calls = nil
}
