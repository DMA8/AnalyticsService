package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/g6834/team31/analytics/internal/config"
	h "gitlab.com/g6834/team31/analytics/internal/adapters/http"
	models "gitlab.com/g6834/team31/analytics/internal/domain/models"
	mocks "gitlab.com/g6834/team31/analytics/internal/mocks"
	"gitlab.com/g6834/team31/auth/pkg/logging"
)

func TestValidateToken(t *testing.T) {
	cfg := &config.Config{
		HTTP: config.HTTP{
			Port:       ":8080",
			ApiVersion: "/analytics/v1",
		},
	}
	ctr := gomock.NewController(t)
	ctr2 := gomock.NewController(t)
	mockService := mocks.NewMockAnalytics(ctr)
	clientAuth := mocks.NewMockClientAuth(ctr2)
	l := logging.New("info")
	handler := h.New(mockService, clientAuth, cfg, &l)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	handlerToTest := handler.ValidateToken(nextHandler)

	// create a mock request to use
	req := httptest.NewRequest("GET", "https://test.ru", nil)
	rec := httptest.NewRecorder()
	// call the handler using a mock response recorder (we'll not use that anyway)
	handlerToTest.ServeHTTP(rec, req)
	responseBody := rec.Body.Bytes()

	var target h.Message
	err := json.Unmarshal(responseBody, &target)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, target.StatusCode)
	assert.Equal(t, true, target.IsError)
	assert.Equal(t, "cookies doesn't exists in headers", target.Message)

	var target2 h.Message
	req2 := httptest.NewRequest("GET", "https://test.ru", nil)
	req2.AddCookie(&http.Cookie{Name: "accessToken", Value: "1234"})
	rec2 := httptest.NewRecorder()
	handlerToTest.ServeHTTP(rec2, req2)
	responseBody2 := rec2.Body.Bytes()

	err = json.Unmarshal(responseBody2, &target2)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, target2.StatusCode)
	assert.Equal(t, true, target2.IsError)
	assert.Equal(t, "cookies doesn't exists in headers", target2.Message)

	// Validate
	req3 := httptest.NewRequest("GET", "/approved_tasks", nil)
	req3.AddCookie(&http.Cookie{Name: "accessToken", Value: "1234"})
	req3.AddCookie(&http.Cookie{Name: "refreshToken", Value: "1234"})
	rec3 := httptest.NewRecorder()
	clientAuth.EXPECT().Validate(context.Background(), models.JWTTokens{
		Access:  "1234",
		Refresh: "1234",
	})
	handlerToTest.ServeHTTP(rec3, req3)
}
