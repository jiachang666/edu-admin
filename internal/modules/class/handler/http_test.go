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

func TestCreateRejectsOutOfScopeTeacherAssignment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"name":"越权班级","courseId":1,"teacherId":2,"campus":"明发校区","capacity":12,"weeklySchedule":"周六 09:00-10:30","startDate":"2026-06-20","endDate":"2026-09-20","status":"开班中","remark":"越权测试"}`
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/v1/classes", strings.NewReader(body))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Set("current_permissions", []string{"classes:manage"})
	context.Set("current_user_id", uint64(4))
	context.Set("current_role", "teacher")
	context.Set("current_user_name", "周老师")

	handler := New(eduservice.New(nil))
	handler.create(context)

	assertClassErrorResponse(t, recorder, 404, "teacher not found")
}

func TestUpdateRejectsOutOfScopeTeacherClass(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"name":"英语阅读进阶班","courseId":2,"teacherId":2,"campus":"百汇校区","capacity":12,"weeklySchedule":"周六 14:00-15:30","startDate":"2026-06-01","endDate":"2026-09-01","status":"开班中","remark":"越权测试"}`
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/classes/2", strings.NewReader(body))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Params = gin.Params{{Key: "id", Value: "2"}}
	context.Set("current_permissions", []string{"classes:manage"})
	context.Set("current_user_id", uint64(4))
	context.Set("current_role", "teacher")
	context.Set("current_user_name", "周老师")

	handler := New(eduservice.New(nil))
	handler.update(context)

	assertClassErrorResponse(t, recorder, 404, "class not found")
}

func TestUpdateRejectsOutOfScopeTeacherAssignment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"name":"周末奥数提高班","courseId":1,"teacherId":2,"campus":"明发校区","capacity":16,"weeklySchedule":"周六 09:00-10:30","startDate":"2026-06-01","endDate":"2026-09-01","status":"开班中","remark":"越权测试"}`
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/classes/1", strings.NewReader(body))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Params = gin.Params{{Key: "id", Value: "1"}}
	context.Set("current_permissions", []string{"classes:manage"})
	context.Set("current_user_id", uint64(4))
	context.Set("current_role", "teacher")
	context.Set("current_user_name", "周老师")

	handler := New(eduservice.New(nil))
	handler.update(context)

	assertClassErrorResponse(t, recorder, 404, "teacher not found")
}

func TestAddStudentsRejectsOutOfScopeTeacherClass(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"studentIds":[3]}`
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/v1/classes/2/students", strings.NewReader(body))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Params = gin.Params{{Key: "id", Value: "2"}}
	context.Set("current_permissions", []string{"classes:manage"})
	context.Set("current_user_id", uint64(4))
	context.Set("current_role", "teacher")
	context.Set("current_user_name", "周老师")

	handler := New(eduservice.New(nil))
	handler.addStudents(context)

	assertClassErrorResponse(t, recorder, 404, "class not found")
}

func TestAddStudentsRejectsOutOfScopeTeacherStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	body := `{"studentIds":[3]}`
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/v1/classes/1/students", strings.NewReader(body))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Params = gin.Params{{Key: "id", Value: "1"}}
	context.Set("current_permissions", []string{"classes:manage"})
	context.Set("current_user_id", uint64(4))
	context.Set("current_role", "teacher")
	context.Set("current_user_name", "周老师")

	handler := New(eduservice.New(nil))
	handler.addStudents(context)

	assertClassErrorResponse(t, recorder, 404, "student not found")
}

func TestRemoveStudentRejectsOutOfScopeTeacherClass(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/classes/2/students/3", nil)
	context.Params = gin.Params{
		{Key: "id", Value: "2"},
		{Key: "studentId", Value: "3"},
	}
	context.Set("current_permissions", []string{"classes:manage"})
	context.Set("current_user_id", uint64(4))
	context.Set("current_role", "teacher")
	context.Set("current_user_name", "周老师")

	handler := New(eduservice.New(nil))
	handler.removeStudent(context)

	assertClassErrorResponse(t, recorder, 404, "class not found")
}

func assertClassErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder, statusCode int, message string) {
	t.Helper()

	if recorder.Code != statusCode {
		t.Fatalf("expected status %d, got %d", statusCode, recorder.Code)
	}

	var responseBody map[string]any
	decodeErr := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	if decodeErr != nil {
		t.Fatalf("failed to decode response body: %v", decodeErr)
	}
	if responseBody["message"] != message {
		t.Fatalf("expected %q message, got %#v", message, responseBody["message"])
	}
}
