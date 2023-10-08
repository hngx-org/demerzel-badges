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
	"demerzel-badges/internal/db"
	"demerzel-badges/internal/models"
	"demerzel-badges/pkg/response"
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


func TestCleanupCreateBadge(t *testing.T) {
	defer mockDB.AssertExpectations(t)
	mockDB.Calls = nil
}


func setupAssignBadgeRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/api/user/badges", handlers.AssignBadgeHandler)
	return router
}

func TestAssignBadgeHandler_Success(t *testing.T) {
	mockDB := new(MockDatabase)
	db.DB = mockDB

	router := setupAssignBadgeRouter()

	testPayload := []byte(`{
		"user_id": "sample_user_id",
		"badge_id": 1,
		"assessment_id": 1
	}`)

	mockDB.On("CheckIfBadgeIsValid", uint(1)).Return(true)
	mockDB.On("VerifyAssessment", uint(1)).Return(true)
	mockDB.On("AssignBadge", "sample_user_id", uint(1), uint(1)).Return(&models.UserBadge{}, nil)

	req, err := http.NewRequest("POST", "/api/user/badges", bytes.NewBuffer(testPayload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	router.ServeHTTP(w, req)

	handlers.AssignBadgeHandler(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Badge Assigned Successfully")
	mockDB.AssertExpectations(t)
}

func TestAssignBadgeHandler_AssignBadgeError(t *testing.T) {
	mockDB := new(MockDatabase)
	db.DB = mockDB

	router := setupAssignBadgeRouter()

	testPayload := []byte(`{
		"user_id": "sample_user_id",
		"badge_id": 1,
		"assessment_id": 1
	}`)

	mockDB.On("CheckIfBadgeIsValid", uint(1)).Return(true)
	mockDB.On("VerifyAssessment", uint(1)).Return(true)
	mockDB.On("AssignBadge", "sample_user_id", uint(1), uint(1)).Return(nil, errors.New("assign badge error"))

	req, err := http.NewRequest("POST", "/api/user/badges", bytes.NewBuffer(testPayload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	router.ServeHTTP(w, req)

	handlers.AssignBadgeHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Unable to assign badge")
	mockDB.AssertExpectations(t)
}


func TestCleanupAssignBadge(t *testing.T) {
	defer mockDB.AssertExpectations(t)
	mockDB.Calls = nil
}



func setupGetUserBadgeRouter() *gin.Engine {
	router := gin.Default()
	router.GET("api/user/badges/:userId/skill/:skillId", handlers.GetUserBadgeHandler)
	return router
}

func TestGetUserBadgeHandler_Success(t *testing.T) {
	mockDB := new(MockDatabase)
	db.DB = mockDB

	router := setupGetUserBadgeRouter()

	userId := "sample_user_id"
	skillId := "sample_skill_id"

	mockDB.On("Where", "user_id=? AND skill_id=?", userId, skillId).Return(mockDB)
	mockDB.On("Find", mock.Anything, "user_id=? AND skill_id=?", userId, skillId).Return(mockDB)
	mockDB.On("Len", mock.Anything).Return(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "api/user/badges/"+userId+"/skill/"+skillId, nil)

	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "userId", Value: userId}, gin.Param{Key: "skillId", Value: skillId})
	c.Request = req

	router.ServeHTTP(w, req)

	handlers.GetUserBadgeHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User Badge Retrieved Successfully")
	mockDB.AssertExpectations(t)
}

func TestGetUserBadgeHandler_NotFound(t *testing.T) {
	mockDB := new(MockDatabase)
	db.DB = mockDB

	router := setupGetUserBadgeRouter()

	userId := "non_existent_user_id"
	skillId := "non_existent_skill_id"

	mockDB.On("Where", "user_id=? AND skill_id=?", userId, skillId).Return(mockDB)
	mockDB.On("Find", mock.Anything, "user_id=? AND skill_id=?", userId, skillId).Return(mockDB)
	mockDB.On("Len", mock.Anything).Return(0)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "api/user/badges/"+userId+"/skill/"+skillId, nil)

	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "userId", Value: userId}, gin.Param{Key: "skillId", Value: skillId})
	c.Request = req

	router.ServeHTTP(w, req)

	handlers.GetUserBadgeHandler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User Badge not Found")
	mockDB.AssertExpectations(t)
}

func TestCleanupGetUserBadge(t *testing.T) {
	defer mockDB.AssertExpectations(t)
	mockDB.Calls = nil
}
