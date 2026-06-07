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

func TestCreateRejectsOutOfScopeTeacherNoticeTarget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"title":"越权通知","content":"尝试发给别班","category":"课程通知","targetScope":"英语阅读进阶班家长群","relatedClassId":2,"relatedScheduleId":0,"studentIds":[],"status":"草稿","author":"周老师"}`
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/v1/notices", strings.NewReader(body))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Set("current_permissions", []string{"notices:manage"})
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
	if responseBody["message"] != "notice target not found" {
		t.Fatalf("expected notice target not found message, got %#v", responseBody["message"])
	}
}
