package res

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestWriteSuccess kiểm tra việc ghi response thành công
func TestWriteSuccess(t *testing.T) {
	t.Run("Default Success Message", func(t *testing.T) {
		rec := httptest.NewRecorder()
		data := map[string]string{"id": "1"}

		err := WriteSuccess(rec, data)
		if err != nil {
			t.Fatalf("WriteSuccess() error = %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("HTTP Status = %d, want %d", rec.Code, http.StatusOK)
		}

		var resp Envelope
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse body: %v", err)
		}

		expected := Envelope{
			Code:  "0",
			Msg:   "Success",
			Data:  map[string]any{"id": "1"},
			Extra: nil,
		}

		// So sánh gián tiếp
		if resp.Code != expected.Code || resp.Msg != expected.Msg {
			t.Errorf("Response = %+v, want %+v", resp, expected)
		}
	})

	t.Run("Custom Success Message", func(t *testing.T) {
		rec := httptest.NewRecorder()
		data := []int{1, 2, 3}

		err := WriteSuccess(rec, data, "Lấy dữ liệu thành công")
		if err != nil {
			t.Fatalf("WriteSuccess() error = %v", err)
		}

		var resp Envelope
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse body: %v", err)
		}

		if resp.Msg != "Lấy dữ liệu thành công" {
			t.Errorf("Msg = %q, want %q", resp.Msg, "Lấy dữ liệu thành công")
		}
	})
}

// TestWritePaginate kiểm tra việc ghi response phân trang thành công (cả Paginate và PaginateCursor)
func TestWritePaginate(t *testing.T) {
	t.Run("Paginate (No Cursor)", func(t *testing.T) {
		rec := httptest.NewRecorder()
		data := []string{"user1"}

		err := WritePaginate(rec, data, 50)
		if err != nil {
			t.Fatalf("WritePaginate() error = %v", err)
		}

		// Parse sang map để chắc chắn trường cursor biến mất khỏi JSON (không xuất hiện key "cursor")
		var raw map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
			t.Fatalf("Failed to parse body: %v", err)
		}

		extra, ok := raw["extra"].(map[string]any)
		if !ok {
			t.Fatalf("extra is missing or not a map")
		}

		if int(extra["total"].(float64)) != 50 {
			t.Errorf("Total = %v, want 50", extra["total"])
		}

		// Kiểm tra key "cursor" không có mặt trong map
		if _, exists := extra["cursor"]; exists {
			t.Errorf("cursor field should be omitted from JSON when empty")
		}
	})

	t.Run("PaginateCursor (With Cursor)", func(t *testing.T) {
		rec := httptest.NewRecorder()
		data := []string{"user1", "user2"}

		err := WritePaginateCursor(rec, data, 100, "next_cursor_123")
		if err != nil {
			t.Fatalf("WritePaginateCursor() error = %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("HTTP Status = %d, want %d", rec.Code, http.StatusOK)
		}

		// Parse để đối chiếu
		var result struct {
			Code  string `json:"code"`
			Msg   string `json:"msg"`
			Data  any    `json:"data"`
			Extra struct {
				Total  int64  `json:"total"`
				Cursor string `json:"cursor"`
			} `json:"extra"`
		}

		if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
			t.Fatalf("Failed to parse body: %v", err)
		}

		if result.Extra.Total != 100 {
			t.Errorf("Total = %d, want 100", result.Extra.Total)
		}
		if result.Extra.Cursor != "next_cursor_123" {
			t.Errorf("Cursor = %q, want %q", result.Extra.Cursor, "next_cursor_123")
		}
	})
}

// TestWriteValidationError kiểm tra lỗi validation (422)
func TestWriteValidationError(t *testing.T) {
	rec := httptest.NewRecorder()
	errorsMap := map[string]string{
		"email": "Định dạng email sai",
	}

	err := WriteValidationError(rec, errorsMap)
	if err != nil {
		t.Fatalf("WriteValidationError() error = %v", err)
	}

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("HTTP Status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}

	var resp Envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse body: %v", err)
	}

	if resp.Code != "VALIDATION_FAILED" {
		t.Errorf("Code = %q, want %q", resp.Code, "VALIDATION_FAILED")
	}

	if resp.Msg != "Dữ liệu đầu vào không hợp lệ" {
		t.Errorf("Msg = %q, want %q", resp.Msg, "Dữ liệu đầu vào không hợp lệ")
	}
}

// TestWriteCommonErrors kiểm tra toàn bộ các helper lỗi phổ biến
func TestWriteCommonErrors(t *testing.T) {
	tests := []struct {
		name       string
		helperFunc func(w http.ResponseWriter, msg string) error
		wantStatus int
		wantCode   string
		wantMsg    string
	}{
		{"BadRequest", func(w http.ResponseWriter, msg string) error { return WriteBadRequest(w, msg) }, http.StatusBadRequest, "BAD_REQUEST", "invalid payload"},
		{"Unauthorized", func(w http.ResponseWriter, msg string) error { return WriteUnauthorized(w, msg) }, http.StatusUnauthorized, "UNAUTHORIZED", "please login"},
		{"Forbidden", func(w http.ResponseWriter, msg string) error { return WriteForbidden(w, msg) }, http.StatusForbidden, "FORBIDDEN", "no access"},
		{"NotFound", func(w http.ResponseWriter, msg string) error { return WriteNotFound(w, msg) }, http.StatusNotFound, "NOT_FOUND", "not found"},
		{"Conflict", func(w http.ResponseWriter, msg string) error { return WriteConflict(w, msg) }, http.StatusConflict, "CONFLICT", "already exists"},
		{"PageExpired", func(w http.ResponseWriter, msg string) error { return WritePageExpired(w, msg) }, 419, "PAGE_EXPIRED", "csrf expired"},
		{"RateLimit", func(w http.ResponseWriter, msg string) error { return WriteRateLimit(w, msg) }, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "too fast"},
		{"InternalError", func(w http.ResponseWriter, msg string) error { return WriteInternalServerError(w, msg) }, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "panic"},
		{"Maintenance", func(w http.ResponseWriter, msg string) error { return WriteMaintenance(w, msg) }, http.StatusServiceUnavailable, "MAINTENANCE", "upgrading"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			err := tt.helperFunc(rec, tt.wantMsg)
			if err != nil {
				t.Fatalf("Helper error = %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("HTTP Status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var resp Envelope
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse body: %v", err)
			}

			if resp.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", resp.Code, tt.wantCode)
			}
			if resp.Msg != tt.wantMsg {
				t.Errorf("Msg = %q, want %q", resp.Msg, tt.wantMsg)
			}
		})
	}
}

// MockGinContext giả lập Context của Gin framework
type MockGinContext struct {
	Code int
	Obj  any
}

func (m *MockGinContext) JSON(code int, obj any) {
	m.Code = code
	m.Obj = obj
}

// TestGinCompatibility kiểm tra tính tương thích với Gin
func TestGinCompatibility(t *testing.T) {
	t.Run("Gin Success", func(t *testing.T) {
		ctx := &MockGinContext{}
		data := "my-data"
		GinSuccess(ctx, data, "Custom Ok")

		if ctx.Code != http.StatusOK {
			t.Errorf("Gin Status = %d, want 200", ctx.Code)
		}

		env, ok := ctx.Obj.(Envelope)
		if !ok {
			t.Fatalf("Gin Obj is not Envelope")
		}

		if env.Code != "0" || env.Msg != "Custom Ok" || env.Data != data {
			t.Errorf("Envelope = %+v", env)
		}
	})

	t.Run("Gin ValidationError", func(t *testing.T) {
		ctx := &MockGinContext{}
		errs := map[string]string{"name": "required"}
		GinValidationError(ctx, errs)

		if ctx.Code != http.StatusUnprocessableEntity {
			t.Errorf("Gin Status = %d, want 422", ctx.Code)
		}

		env := ctx.Obj.(Envelope)
		if env.Code != "VALIDATION_FAILED" || !reflect.DeepEqual(env.Extra, errs) {
			t.Errorf("Envelope = %+v", env)
		}
	})
}

// === BENCHMARKS ===

// TestGlobalMessageOverride kiểm tra việc thay đổi các biến global message và code hoạt động chính xác
func TestGlobalMessageOverride(t *testing.T) {
	// Lưu trữ giá trị ban đầu để khôi phục
	oldSuccessCode := DefaultSuccessCode
	oldSuccessMsg := DefaultSuccessMsg
	oldUnauthorizedCode := DefaultUnauthorizedCode
	oldUnauthorizedMsg := DefaultUnauthorizedMsg
	defer func() {
		DefaultSuccessCode = oldSuccessCode
		DefaultSuccessMsg = oldSuccessMsg
		DefaultUnauthorizedCode = oldUnauthorizedCode
		DefaultUnauthorizedMsg = oldUnauthorizedMsg
	}()

	// Thay đổi giá trị biến global
	DefaultSuccessCode = "20000"
	DefaultSuccessMsg = "Thao tác thành công"
	DefaultUnauthorizedCode = "40101"
	DefaultUnauthorizedMsg = "Yêu cầu đăng nhập hệ thống"

	// 1. Kiểm tra success
	recSuccess := httptest.NewRecorder()
	_ = WriteSuccess(recSuccess, nil)

	var respSuccess Envelope
	_ = json.Unmarshal(recSuccess.Body.Bytes(), &respSuccess)
	if respSuccess.Code != "20000" {
		t.Errorf("Success Code = %q, want %q", respSuccess.Code, "20000")
	}
	if respSuccess.Msg != "Thao tác thành công" {
		t.Errorf("Success Msg = %q, want %q", respSuccess.Msg, "Thao tác thành công")
	}

	// 2. Kiểm tra unauthorized (truyền chuỗi rỗng)
	recAuth := httptest.NewRecorder()
	_ = WriteUnauthorized(recAuth, "")

	var respAuth Envelope
	_ = json.Unmarshal(recAuth.Body.Bytes(), &respAuth)
	if respAuth.Code != "40101" {
		t.Errorf("Unauthorized Code = %q, want %q", respAuth.Code, "40101")
	}
	if respAuth.Msg != "Yêu cầu đăng nhập hệ thống" {
		t.Errorf("Unauthorized Msg = %q, want %q", respAuth.Msg, "Yêu cầu đăng nhập hệ thống")
	}
}

func BenchmarkWriteSuccess(b *testing.B) {
	rec := httptest.NewRecorder()
	data := map[string]string{"id": "1", "status": "active"}
	b.ResetTimer()
	for b.Loop() {
		rec.Body.Reset()
		_ = WriteSuccess(rec, data)
	}
}

