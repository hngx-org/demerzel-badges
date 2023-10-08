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

    "demerzel-badges/internal/db"
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

func setupTestRouter() *gin.Engine {
    router := gin.Default()
    router.POST("/badges", handlers.CreateBadgeHandler)
    return router
}

func TestIntegrationCreateBadgeHandler_Success(t *testing.T) {
    router := setupTestRouter()
    payload := []byte(`{"skill_id": 1, "name": "Beginner", "min_score": 0, "max_score": 50}`)
    req, _ := http.NewRequest("POST", "/badges", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    testDB := db.SetupTestDB()
    t.Cleanup(func() {
        db.CloseTestDB(testDB)
    })

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req.WithContext(db.SetDBContext(req.Context(), testDB)))

    assert.Equal(t, http.StatusCreated, recorder.Code)
    assert.JSONEq(t, `{"message": "Badge Created Successfully", "badge": {"id": 1, "skill_id": 1, "name": "beginner", "min_score": 0, "max_score": 50}}`, recorder.Body.String())

    mockDB.AssertExpectations(t)
}

func TestIntegrationCreateBadgeHandler_BadgeExists(t *testing.T) {
    testDB := db.SetupTestDB()
    t.Cleanup(func() {
        db.CloseTestDB(testDB)
    })

    router := setupTestRouter()
    payload := []byte(`{"skill_id": 1, "name": "Beginner", "min_score": 0, "max_score": 50}`)
    req, _ := http.NewRequest("POST", "/badges", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req.WithContext(db.SetDBContext(req.Context(), testDB)))

    assert.Equal(t, http.StatusBadRequest, recorder.Code)
    assert.Contains(t, recorder.Body.String(), "Badge already exists")

    mockDB.AssertExpectations(t)
}

func TestIntegrationCreateBadgeHandler_CreateBadgeError(t *testing.T) {
    testDB := db.SetupTestDB()
    t.Cleanup(func() {
        db.CloseTestDB(testDB)
    })

    router := setupTestRouter()
    payload := []byte(`{"skill_id": 1, "name": "Beginner", "min_score": 0, "max_score": 50}`)
    req, _ := http.NewRequest("POST", "/badges", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    mockDB.On("FindSkillById", uint(1)).Return(&models.Skill{ID: 1}, nil)
    mockDB.On("BadgeExists", uint(1), models.Badge("beginner")).Return(false)
    mockDB.On("CreateBadge", mock.AnythingOfType("*models.SkillBadge")).Return(nil, errors.New("database error"))

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req.WithContext(db.SetDBContext(req.Context(), testDB)))

    assert.Equal(t, http.StatusInternalServerError, recorder.Code)
    assert.Contains(t, recorder.Body.String(), "Unable to create badge")

    mockDB.AssertExpectations(t)
}

func TestIntegrationCreateBadgeHandler_InvalidInput(t *testing.T) {
    testDB := db.SetupTestDB()
    t.Cleanup(func() {
        db.CloseTestDB(testDB)
    })

    router := setupTestRouter()
    payload := []byte(`{"skill_id": 1, "name": "Invalid", "min_score": -5, "max_score": 20}`)
    req, _ := http.NewRequest("POST", "/badges", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    mockDB.On("FindSkillById", uint(1)).Return(&models.Skill{ID: 1}, nil)
    mockDB.On("BadgeExists", uint(1), models.Badge("invalid")).Return(false)
    mockDB.On("CreateBadge", mock.AnythingOfType("*models.SkillBadge")).Return(nil, errors.New("Error! Invalid Input"))

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req.WithContext(db.SetDBContext(req.Context(), testDB)))

    assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
    assert.Contains(t, recorder.Body.String(), "min_score should be at least 0")

    mockDB.AssertExpectations(t)
}

func TestIntegrationCreateBadgeHandler_DatabaseError(t *testing.T) {
    testDB := db.SetupTestDB()
    t.Cleanup(func() {
        db.CloseTestDB(testDB)
    })

    router := setupTestRouter()
    payload := []byte(`{"skill_id": 1, "name": "Beginner", "min_score": 0, "max_score": 50}`)
    req, _ := http.NewRequest("POST", "/badges", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    mockDB.On("FindSkillById", uint(1)).Return(nil, errors.New("database error"))
    mockDB.On("BadgeExists", uint(1), models.Badge("beginner")).Return(false)
    mockDB.On("CreateBadge", mock.AnythingOfType("*models.SkillBadge")).Return(nil, errors.New("database error"))

    recorder := httptest.NewRecorder()
    router.ServeHTTP(httptest.NewRecorder(), req.WithContext(db.SetDBContext(req.Context(), testDB)))

    assert.Equal(t, http.StatusInternalServerError, recorder.Code)
    assert.Contains(t, recorder.Body.String(), "Internal Server Error")

    mockDB.AssertExpectations(t)
}

func TestIntegrationCreateBadgeHandler_FindSkillByIdError(t *testing.T) {
    testDB := db.SetupTestDB()
    t.Cleanup(func() {
        db.CloseTestDB(testDB)
    })

    router := setupTestRouter()
    payload := []byte(`{"skill_id": 1, "name": "Beginner", "min_score": 0, "max_score": 50}`)
    req, _ := http.NewRequest("POST", "/badges", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")

    mockDB.On("FindSkillById", uint(1)).Return(nil, nil)
    mockDB.On("BadgeExists", uint(1), models.Badge("beginner")).Return(false)
    mockDB.On("CreateBadge", mock.AnythingOfType("*models.SkillBadge")).Return(nil, nil)

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req.WithContext(db.SetDBContext(req.Context(), testDB)))

    assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
    assert.Contains(t, recorder.Body.String(), "No skill found matching provided ID")
    mockDB.AssertExpectations(t)
}

func TestCleanup(t *testing.T) {
    mockDB.AssertExpectations(t)
    mockDB.Calls = nil
}
