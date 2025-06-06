package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	// Setup Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Test case: Invalid input error
	HandleError(c, http.StatusBadRequest, ErrInvalidInput, "Invalid input format")

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := `{"code":400,"error":"INVALID_INPUT","message":"Invalid input format"}`
	assert.JSONEq(t, expected, w.Body.String())
}
