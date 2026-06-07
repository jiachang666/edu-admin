package response

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPaginatedUsesRequestedPageAndPageSize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest("GET", "/api/v1/students?page=2&pageSize=2", nil)

	Paginated(context, []int{1, 2, 3, 4, 5})

	var responseBody struct {
		Data struct {
			List     []int `json:"list"`
			Total    int   `json:"total"`
			Page     int   `json:"page"`
			PageSize int   `json:"pageSize"`
		} `json:"data"`
	}

	decodeErr := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	if decodeErr != nil {
		t.Fatalf("failed to decode response body: %v", decodeErr)
	}

	if responseBody.Data.Total != 5 {
		t.Fatalf("expected total 5, got %d", responseBody.Data.Total)
	}
	if responseBody.Data.Page != 2 {
		t.Fatalf("expected page 2, got %d", responseBody.Data.Page)
	}
	if responseBody.Data.PageSize != 2 {
		t.Fatalf("expected page size 2, got %d", responseBody.Data.PageSize)
	}
	if len(responseBody.Data.List) != 2 || responseBody.Data.List[0] != 3 || responseBody.Data.List[1] != 4 {
		t.Fatalf("expected list [3 4], got %#v", responseBody.Data.List)
	}
}

func TestPaginatedFallsBackForInvalidPageParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest("GET", "/api/v1/students?page=0&pageSize=abc", nil)

	Paginated(context, []int{1, 2, 3})

	var responseBody struct {
		Data struct {
			List     []int `json:"list"`
			Total    int   `json:"total"`
			Page     int   `json:"page"`
			PageSize int   `json:"pageSize"`
		} `json:"data"`
	}

	decodeErr := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	if decodeErr != nil {
		t.Fatalf("failed to decode response body: %v", decodeErr)
	}

	if responseBody.Data.Page != defaultPage {
		t.Fatalf("expected default page %d, got %d", defaultPage, responseBody.Data.Page)
	}
	if responseBody.Data.PageSize != defaultPageSize {
		t.Fatalf("expected default page size %d, got %d", defaultPageSize, responseBody.Data.PageSize)
	}
	if len(responseBody.Data.List) != 3 {
		t.Fatalf("expected full list of 3 items, got %#v", responseBody.Data.List)
	}
}

func TestPaginatedReturnsEmptyListWhenPageIsOutOfRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest("GET", "/api/v1/students?page=3&pageSize=2", nil)

	Paginated(context, []int{1, 2, 3})

	var responseBody struct {
		Data struct {
			List     []int `json:"list"`
			Total    int   `json:"total"`
			Page     int   `json:"page"`
			PageSize int   `json:"pageSize"`
		} `json:"data"`
	}

	decodeErr := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	if decodeErr != nil {
		t.Fatalf("failed to decode response body: %v", decodeErr)
	}

	if responseBody.Data.Total != 3 {
		t.Fatalf("expected total 3, got %d", responseBody.Data.Total)
	}
	if len(responseBody.Data.List) != 0 {
		t.Fatalf("expected empty list, got %#v", responseBody.Data.List)
	}
}

func TestPaginatedCapsOversizedPageSize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	items := make([]int, 0, 120)
	for index := 1; index <= 120; index++ {
		items = append(items, index)
	}

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest("GET", "/api/v1/students?page=1&pageSize=500", nil)

	Paginated(context, items)

	var responseBody struct {
		Data struct {
			List     []int `json:"list"`
			Total    int   `json:"total"`
			Page     int   `json:"page"`
			PageSize int   `json:"pageSize"`
		} `json:"data"`
	}

	decodeErr := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	if decodeErr != nil {
		t.Fatalf("failed to decode response body: %v", decodeErr)
	}

	if responseBody.Data.Total != 120 {
		t.Fatalf("expected total 120, got %d", responseBody.Data.Total)
	}
	if responseBody.Data.PageSize != maxPageSize {
		t.Fatalf("expected capped page size %d, got %d", maxPageSize, responseBody.Data.PageSize)
	}
	if len(responseBody.Data.List) != maxPageSize {
		t.Fatalf("expected %d items after cap, got %d", maxPageSize, len(responseBody.Data.List))
	}
}
