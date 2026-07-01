package res

import (
	"net/http"
)

// --- GIN FRAMEWORK COMPATIBILITY (Zero-Dependency) ---

// GinContext định nghĩa interface tương thích với *gin.Context.
type GinContext interface {
	JSON(code int, obj any)
}

// GinJSON ghi trực tiếp JSON response tùy biến chuẩn vào *gin.Context.
func GinJSON(c GinContext, statusCode int, code string, msg string, data any, extra ...any) {
	var ext any = nil
	if len(extra) > 0 {
		ext = extra[0]
	}
	c.JSON(statusCode, Envelope{
		Code:  code,
		Msg:   msg,
		Data:  data,
		Extra: ext,
	})
}

// GinSuccess ghi trực tiếp JSON response thành công (Status 200, code mặc định DefaultSuccessCode) vào *gin.Context.
func GinSuccess(c GinContext, data any, msg ...string) {
	message := DefaultSuccessMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	GinJSON(c, http.StatusOK, DefaultSuccessCode, message, data)
}

// GinPaginate ghi trực tiếp JSON response thành công kèm phân trang thông thường (chỉ chứa total) vào *gin.Context.
func GinPaginate(c GinContext, data any, total int64, msg ...string) {
	message := DefaultSuccessMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	c.JSON(http.StatusOK, Envelope{
		Code: DefaultSuccessCode,
		Msg:  message,
		Data: data,
		Extra: Pagination{
			Total: total,
		},
	})
}

// GinPaginateCursor ghi trực tiếp JSON response thành công kèm phân trang bằng cursor (chứa total và cursor) vào *gin.Context.
func GinPaginateCursor(c GinContext, data any, total int64, cursor string, msg ...string) {
	message := DefaultSuccessMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	c.JSON(http.StatusOK, Envelope{
		Code: DefaultSuccessCode,
		Msg:  message,
		Data: data,
		Extra: Pagination{
			Total:  total,
			Cursor: cursor,
		},
	})
}

// GinError ghi trực tiếp JSON response lỗi vào *gin.Context.
func GinError(c GinContext, statusCode int, code string, msg string, extra ...any) {
	GinJSON(c, statusCode, code, msg, nil, extra...)
}

// GinValidationError ghi trực tiếp JSON response lỗi validation vào *gin.Context (HTTP 422, code mặc định DefaultValidationCode).
func GinValidationError(c GinContext, errors any, msg ...string) {
	message := DefaultValidationMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	GinJSON(c, http.StatusUnprocessableEntity, DefaultValidationCode, message, nil, errors)
}

// GinBadRequest ghi JSON response lỗi 400 Bad Request vào Gin.
func GinBadRequest(c GinContext, msg string, extra ...any) {
	if msg == "" {
		msg = DefaultBadRequestMsg
	}
	GinError(c, http.StatusBadRequest, DefaultBadRequestCode, msg, extra...)
}

// GinUnauthorized ghi JSON response lỗi 401 Unauthorized vào Gin.
func GinUnauthorized(c GinContext, msg string, extra ...any) {
	if msg == "" {
		msg = DefaultUnauthorizedMsg
	}
	GinError(c, http.StatusUnauthorized, DefaultUnauthorizedCode, msg, extra...)
}

// GinForbidden ghi JSON response lỗi 403 Forbidden vào Gin.
func GinForbidden(c GinContext, msg string, extra ...any) {
	if msg == "" {
		msg = DefaultForbiddenMsg
	}
	GinError(c, http.StatusForbidden, DefaultForbiddenCode, msg, extra...)
}

// GinNotFound ghi JSON response lỗi 404 Not Found vào Gin.
func GinNotFound(c GinContext, msg string, extra ...any) {
	if msg == "" {
		msg = DefaultNotFoundMsg
	}
	GinError(c, http.StatusNotFound, DefaultNotFoundCode, msg, extra...)
}

// GinConflict ghi JSON response lỗi 409 Conflict vào Gin.
func GinConflict(c GinContext, msg string, extra ...any) {
	if msg == "" {
		msg = DefaultConflictMsg
	}
	GinError(c, http.StatusConflict, DefaultConflictCode, msg, extra...)
}

// GinPageExpired ghi JSON response lỗi 419 Page Expired vào Gin.
func GinPageExpired(c GinContext, msg string, extra ...any) {
	if msg == "" {
		msg = DefaultPageExpiredMsg
	}
	GinError(c, 419, DefaultPageExpiredCode, msg, extra...)
}

// GinRateLimit ghi JSON response lỗi 429 Too Many Requests vào Gin.
func GinRateLimit(c GinContext, msg string, extra ...any) {
	if msg == "" {
		msg = DefaultRateLimitMsg
	}
	GinError(c, http.StatusTooManyRequests, DefaultRateLimitCode, msg, extra...)
}

// GinInternalServerError ghi JSON response lỗi 500 Internal Server Error vào Gin.
func GinInternalServerError(c GinContext, msg string, extra ...any) {
	if msg == "" {
		msg = DefaultInternalServerMsg
	}
	GinError(c, http.StatusInternalServerError, DefaultInternalServerCode, msg, extra...)
}

// GinMaintenance ghi JSON response lỗi 503 Service Unavailable vào Gin.
func GinMaintenance(c GinContext, msg string, extra ...any) {
	if msg == "" {
		msg = DefaultMaintenanceMsg
	}
	GinError(c, http.StatusServiceUnavailable, DefaultMaintenanceCode, msg, extra...)
}
