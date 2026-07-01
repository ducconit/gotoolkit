// Package res cung cấp các giải pháp định dạng JSON response chuẩn REST API tối giản,
// hiệu năng cao, zero-dependency và tích hợp tốt với net/http cũng như Gin.
package res

import (
	"encoding/json"
	"net/http"
)

// Envelope đại diện cho vỏ bọc JSON response chuẩn REST API tối giản.
type Envelope struct {
	Code  string `json:"code"`
	Msg   string `json:"msg"`
	Data  any    `json:"data"`
	Extra any    `json:"extra"`
}

// Pagination đại diện cho siêu dữ liệu phân trang chuẩn nằm trong trường extra.
type Pagination struct {
	Total  int64  `json:"total"`
	Cursor string `json:"cursor,omitempty"`
}

// --- CONFIGURATION VARIABLES FOR CODE (Có thể ghi đè lúc ứng dụng khởi chạy) ---
var (
	// DefaultSuccessCode là mã code mặc định cho response thành công.
	DefaultSuccessCode = "0"

	// DefaultValidationCode là mã code mặc định cho lỗi validation.
	DefaultValidationCode = "VALIDATION_FAILED"

	// DefaultBadRequestCode là mã code mặc định cho lỗi 400.
	DefaultBadRequestCode = "BAD_REQUEST"

	// DefaultUnauthorizedCode là mã code mặc định cho lỗi 401.
	DefaultUnauthorizedCode = "UNAUTHORIZED"

	// DefaultForbiddenCode là mã code mặc định cho lỗi 403.
	DefaultForbiddenCode = "FORBIDDEN"

	// DefaultNotFoundCode là mã code mặc định cho lỗi 404.
	DefaultNotFoundCode = "NOT_FOUND"

	// DefaultConflictCode là mã code mặc định cho lỗi 409.
	DefaultConflictCode = "CONFLICT"

	// DefaultPageExpiredCode là mã code mặc định cho lỗi 419.
	DefaultPageExpiredCode = "PAGE_EXPIRED"

	// DefaultRateLimitCode là mã code mặc định cho lỗi 429.
	DefaultRateLimitCode = "RATE_LIMIT_EXCEEDED"

	// DefaultInternalServerCode là mã code mặc định cho lỗi 500.
	DefaultInternalServerCode = "INTERNAL_SERVER_ERROR"

	// DefaultMaintenanceCode là mã code mặc định cho lỗi 503.
	DefaultMaintenanceCode = "MAINTENANCE"
)

// --- CONFIGURATION VARIABLES FOR MESSAGE (Có thể ghi đè lúc ứng dụng khởi chạy) ---
var (
	// DefaultSuccessMsg là thông báo mặc định cho các response thành công.
	DefaultSuccessMsg = "Success"

	// DefaultValidationMsg là thông báo mặc định cho lỗi validation (422).
	DefaultValidationMsg = "Dữ liệu đầu vào không hợp lệ"

	// DefaultBadRequestMsg là thông báo mặc định cho lỗi 400 Bad Request.
	DefaultBadRequestMsg = "Yêu cầu không hợp lệ"

	// DefaultUnauthorizedMsg là thông báo mặc định cho lỗi 401 Unauthorized.
	DefaultUnauthorizedMsg = "Vui lòng đăng nhập"

	// DefaultForbiddenMsg là thông báo mặc định cho lỗi 403 Forbidden.
	DefaultForbiddenMsg = "Không có quyền truy cập"

	// DefaultNotFoundMsg là thông báo mặc định cho lỗi 404 Not Found.
	DefaultNotFoundMsg = "Không tìm thấy tài nguyên"

	// DefaultConflictMsg là thông báo mặc định cho lỗi 409 Conflict.
	DefaultConflictMsg = "Dữ liệu đã tồn tại hoặc xảy ra xung đột"

	// DefaultPageExpiredMsg là thông báo mặc định cho lỗi 419 Page Expired.
	DefaultPageExpiredMsg = "Trang đã hết hạn hoặc phiên làm việc kết thúc"

	// DefaultRateLimitMsg là thông báo mặc định cho lỗi 429 Too Many Requests.
	DefaultRateLimitMsg = "Bạn đã gửi quá nhiều yêu cầu, vui lòng thử lại sau"

	// DefaultInternalServerMsg là thông báo mặc định cho lỗi 500 Internal Server Error.
	DefaultInternalServerMsg = "Lỗi máy chủ nội bộ"

	// DefaultMaintenanceMsg là thông báo mặc định cho lỗi 503 Service Unavailable.
	DefaultMaintenanceMsg = "Hệ thống đang bảo trì, vui lòng quay lại sau"
)

// --- CONSTRUCTORS (Dùng cho Gin/Fiber/Echo...) ---

// SuccessEnvelope tạo một Envelope thành công với code mặc định (DefaultSuccessCode), HTTP Status 200.
func SuccessEnvelope(data any, msg ...string) Envelope {
	message := DefaultSuccessMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return Envelope{
		Code:  DefaultSuccessCode,
		Msg:   message,
		Data:  data,
		Extra: nil,
	}
}

// PaginateEnvelope tạo một Envelope thành công dạng phân trang thông thường với total.
func PaginateEnvelope(data any, total int64, msg ...string) Envelope {
	message := DefaultSuccessMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return Envelope{
		Code:  DefaultSuccessCode,
		Msg:   message,
		Data:  data,
		Extra: Pagination{
			Total: total,
		},
	}
}

// PaginateCursorEnvelope tạo một Envelope thành công dạng phân trang bằng cursor với total và cursor.
func PaginateCursorEnvelope(data any, total int64, cursor string, msg ...string) Envelope {
	message := DefaultSuccessMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return Envelope{
		Code:  DefaultSuccessCode,
		Msg:   message,
		Data:  data,
		Extra: Pagination{
			Total:  total,
			Cursor: cursor,
		},
	}
}

// ErrorEnvelope tạo một Envelope lỗi với code lỗi và msg.
func ErrorEnvelope(code string, msg string, extra ...any) Envelope {
	var ext any = nil
	if len(extra) > 0 {
		ext = extra[0]
	}
	return Envelope{
		Code:  code,
		Msg:   msg,
		Data:  nil,
		Extra: ext,
	}
}

// ValidationErrorEnvelope tạo một Envelope lỗi validation với code mặc định (DefaultValidationCode), HTTP Status 422.
func ValidationErrorEnvelope(errors any, msg ...string) Envelope {
	message := DefaultValidationMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return Envelope{
		Code:  DefaultValidationCode,
		Msg:   message,
		Data:  nil,
		Extra: errors,
	}
}

// --- HELPERS GHI TRỰC TIẾP VÀO http.ResponseWriter (net/http) ---

// Write ghi trực tiếp một JSON response tùy biến vào ResponseWriter.
func Write(w http.ResponseWriter, statusCode int, code string, msg string, data any, extra ...any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var ext any = nil
	if len(extra) > 0 {
		ext = extra[0]
	}

	env := Envelope{
		Code:  code,
		Msg:   msg,
		Data:  data,
		Extra: ext,
	}

	return json.NewEncoder(w).Encode(env)
}

// WriteSuccess ghi JSON response thành công vào ResponseWriter (HTTP 200, code mặc định DefaultSuccessCode).
func WriteSuccess(w http.ResponseWriter, data any, msg ...string) error {
	message := DefaultSuccessMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return Write(w, http.StatusOK, DefaultSuccessCode, message, data)
}

// WritePaginate ghi JSON response thành công kèm phân trang thông thường vào ResponseWriter.
func WritePaginate(w http.ResponseWriter, data any, total int64, msg ...string) error {
	message := DefaultSuccessMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	env := Envelope{
		Code: DefaultSuccessCode,
		Msg:  message,
		Data: data,
		Extra: Pagination{
			Total: total,
		},
	}
	return json.NewEncoder(w).Encode(env)
}

// WritePaginateCursor ghi JSON response thành công kèm phân trang dạng cursor vào ResponseWriter.
func WritePaginateCursor(w http.ResponseWriter, data any, total int64, cursor string, msg ...string) error {
	message := DefaultSuccessMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	env := Envelope{
		Code: DefaultSuccessCode,
		Msg:  message,
		Data: data,
		Extra: Pagination{
			Total:  total,
			Cursor: cursor,
		},
	}
	return json.NewEncoder(w).Encode(env)
}

// WriteError ghi JSON response lỗi vào ResponseWriter.
func WriteError(w http.ResponseWriter, statusCode int, code string, msg string, extra ...any) error {
	return Write(w, statusCode, code, msg, nil, extra...)
}

// WriteValidationError ghi JSON response lỗi validation vào ResponseWriter (HTTP 422, code mặc định DefaultValidationCode).
func WriteValidationError(w http.ResponseWriter, errors any, msg ...string) error {
	message := DefaultValidationMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return Write(w, http.StatusUnprocessableEntity, DefaultValidationCode, message, nil, errors)
}

// WriteBadRequest ghi JSON response lỗi 400 Bad Request.
func WriteBadRequest(w http.ResponseWriter, msg string, extra ...any) error {
	if msg == "" {
		msg = DefaultBadRequestMsg
	}
	return Write(w, http.StatusBadRequest, DefaultBadRequestCode, msg, nil, extra...)
}

// WriteUnauthorized ghi JSON response lỗi 401 Unauthorized.
func WriteUnauthorized(w http.ResponseWriter, msg string, extra ...any) error {
	if msg == "" {
		msg = DefaultUnauthorizedMsg
	}
	return Write(w, http.StatusUnauthorized, DefaultUnauthorizedCode, msg, nil, extra...)
}

// WriteForbidden ghi JSON response lỗi 403 Forbidden.
func WriteForbidden(w http.ResponseWriter, msg string, extra ...any) error {
	if msg == "" {
		msg = DefaultForbiddenMsg
	}
	return Write(w, http.StatusForbidden, DefaultForbiddenCode, msg, nil, extra...)
}

// WriteNotFound ghi JSON response lỗi 404 Not Found.
func WriteNotFound(w http.ResponseWriter, msg string, extra ...any) error {
	if msg == "" {
		msg = DefaultNotFoundMsg
	}
	return Write(w, http.StatusNotFound, DefaultNotFoundCode, msg, nil, extra...)
}

// WriteConflict ghi JSON response lỗi 409 Conflict.
func WriteConflict(w http.ResponseWriter, msg string, extra ...any) error {
	if msg == "" {
		msg = DefaultConflictMsg
	}
	return Write(w, http.StatusConflict, DefaultConflictCode, msg, nil, extra...)
}

// WritePageExpired ghi JSON response lỗi 419 Page Expired.
func WritePageExpired(w http.ResponseWriter, msg string, extra ...any) error {
	if msg == "" {
		msg = DefaultPageExpiredMsg
	}
	return Write(w, 419, DefaultPageExpiredCode, msg, nil, extra...)
}

// WriteRateLimit ghi JSON response lỗi 429 Too Many Requests.
func WriteRateLimit(w http.ResponseWriter, msg string, extra ...any) error {
	if msg == "" {
		msg = DefaultRateLimitMsg
	}
	return Write(w, http.StatusTooManyRequests, DefaultRateLimitCode, msg, nil, extra...)
}

// WriteInternalServerError ghi JSON response lỗi 500 Internal Server Error.
func WriteInternalServerError(w http.ResponseWriter, msg string, extra ...any) error {
	if msg == "" {
		msg = DefaultInternalServerMsg
	}
	return Write(w, http.StatusInternalServerError, DefaultInternalServerCode, msg, nil, extra...)
}

// WriteMaintenance ghi JSON response lỗi 503 Service Unavailable (Bảo trì).
func WriteMaintenance(w http.ResponseWriter, msg string, extra ...any) error {
	if msg == "" {
		msg = DefaultMaintenanceMsg
	}
	return Write(w, http.StatusServiceUnavailable, DefaultMaintenanceCode, msg, nil, extra...)
}
