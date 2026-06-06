package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	eduservice "edu-admin/internal/modules/edu/service"

	"github.com/gin-gonic/gin"
)

func TestCreateRejectsOutOfScopeTeacherClass(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"classId":2,"scheduleType":"常规课","lessonDate":"2026-06-20","startTime":"10:00","endTime":"11:30","classroom":"B201","remark":"越权测试"}`
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/v1/schedules", strings.NewReader(body))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Set("current_permissions", []string{"schedules:manage"})
	context.Set("current_user_id", uint64(4))
	context.Set("current_role", "teacher")
	context.Set("current_user_name", "周老师")

	handler := New(eduservice.New(nil))
	handler.create(context)

	if recorder.Code != 404 {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}

	var responseBody map[string]any
	decodeErr := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	if decodeErr != nil {
		t.Fatalf("failed to decode response body: %v", decodeErr)
	}
	if responseBody["message"] != "class not found" {
		t.Fatalf("expected class not found message, got %#v", responseBody["message"])
	}
}
