package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOverviewRequiresDashboardPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Set("current_permissions", []string{"students:view"})

	handler := &Handler{}
	handler.overview(context)

	if recorder.Code != 403 {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}

	var responseBody map[string]any
	decodeErr := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	if decodeErr != nil {
		t.Fatalf("failed to decode response body: %v", decodeErr)
	}
	if responseBody["message"] != "forbidden" {
		t.Fatalf("expected forbidden message, got %#v", responseBody["message"])
	}
}
