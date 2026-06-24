package secureapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/ducconit/gotoolkit/encrypt"
)

type contextKey string

var (
	// SessionKeyContextKey là key được sử dụng để lưu SessionKey trong http.Request Context.
	SessionKeyContextKey = contextKey("secureapi_session_key")

	// SessionKeyGinContextKey là key được sử dụng để lưu SessionKey trong gin.Context.
	SessionKeyGinContextKey = "secureapi_session_key"

	// SessionIDHeader là tên HTTP Header chứa Session ID để xác định Session Key tương ứng.
	SessionIDHeader = "X-Secure-Session-ID"

	// Buffer pool dùng chung để tối ưu hóa bộ nhớ cho responseWrapper.
	bufferPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 4096))
		},
	}
)

// Các lỗi thường gặp trong quá trình xử lý bảo mật API.
var (
	ErrMissingSessionID = errors.New("missing session ID header")
)

// HandshakeRequest định nghĩa cấu trúc yêu cầu bắt tay từ Client.
type HandshakeRequest struct {
	ClientPublicKey string `json:"client_public_key"` // Dạng hex string của raw public key (65 bytes)
}

// HandshakeResponse định nghĩa cấu trúc phản hồi bắt tay từ Server.
type HandshakeResponse struct {
	ServerPublicKey string `json:"server_public_key"` // Dạng hex string
	SessionID       string `json:"session_id"`
}

// APIResponse định nghĩa cấu trúc phản hồi chuẩn của hệ thống.
type APIResponse struct {
	Data any    `json:"data"` // Dữ liệu trả về (chứa Base64 ciphertext nếu API được mã hóa thành công)
	Code string `json:"code"` // Mã lỗi (rỗng hoặc "success" nếu thành công)
	Msg  string `json:"msg"`  // Tin nhắn hiển thị hoặc mô tả lỗi
}

// PayloadExtractor định nghĩa hàm trích xuất ciphertext (raw bytes) từ request body nhận được.
type PayloadExtractor func(body []byte) ([]byte, error)

// PayloadWrapper định nghĩa hàm đóng gói ciphertext (raw bytes) thành response body để gửi về.
type PayloadWrapper func(ciphertext []byte) ([]byte, error)

// ErrorHandler định nghĩa hàm xử lý khi xảy ra lỗi trong Middleware.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error, statusCode int)

// SecureConfig chứa các tùy chỉnh về định dạng dữ liệu và xử lý lỗi cho Middleware.
type SecureConfig struct {
	Extractor   PayloadExtractor
	Wrapper     PayloadWrapper
	OnError     ErrorHandler
	MaxBodySize int64 // Giới hạn kích thước request body (bytes) chống DDoS OOM
}

// DefaultOnError là hàm xử lý lỗi mặc định, trả về cấu trúc APIResponse dưới dạng JSON.
func DefaultOnError(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	code := "ERROR"
	switch statusCode {
	case http.StatusBadRequest:
		code = "BAD_REQUEST"
	case http.StatusUnauthorized:
		code = "UNAUTHORIZED"
	case http.StatusInternalServerError:
		code = "INTERNAL_ERROR"
	}

	resp := APIResponse{
		Data: nil,
		Code: code,
		Msg:  err.Error(),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// DefaultConfig cấu hình mặc định sử dụng định dạng JSON APIResponse chuẩn và DefaultOnError.
var DefaultConfig = SecureConfig{
	Extractor: func(body []byte) ([]byte, error) {
		var resp APIResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}
		dataStr, ok := resp.Data.(string)
		if !ok {
			return nil, errors.New("invalid data field type in payload, string expected")
		}
		return base64.StdEncoding.DecodeString(dataStr)
	},
	Wrapper: func(ciphertext []byte) ([]byte, error) {
		resp := APIResponse{
			Data: base64.StdEncoding.EncodeToString(ciphertext),
			Code: "",
			Msg:  "success",
		}
		return json.Marshal(resp)
	},
	OnError:     DefaultOnError,
	MaxBodySize: 10 * 1024 * 1024, // Mặc định 10MB
}

// GetConfig là helper để lấy cấu hình thực tế từ variadic parameter (nếu rỗng sẽ dùng mặc định).
func GetConfig(configs []SecureConfig) SecureConfig {
	if len(configs) > 0 {
		cfg := configs[0]
		if cfg.Extractor == nil {
			cfg.Extractor = DefaultConfig.Extractor
		}
		if cfg.Wrapper == nil {
			cfg.Wrapper = DefaultConfig.Wrapper
		}
		if cfg.OnError == nil {
			cfg.OnError = DefaultConfig.OnError
		}
		if cfg.MaxBodySize <= 0 {
			cfg.MaxBodySize = DefaultConfig.MaxBodySize
		}
		return cfg
	}
	return DefaultConfig
}

// HandshakeHandler trả về một http.HandlerFunc để xử lý yêu cầu bắt tay từ Client.
func HandshakeHandler(store *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Giới hạn request body cho handshake (chỉ tầm vài KB)
		r.Body = http.MaxBytesReader(w, r.Body, 64*1024)

		var req HandshakeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		clientPubKeyBytes, err := hex.DecodeString(req.ClientPublicKey)
		if err != nil {
			http.Error(w, "Invalid client public key format", http.StatusBadRequest)
			return
		}

		// Sinh khóa tạm thời cho Server
		serverPrivKey, serverPubKeyBytes, err := GenerateEphemeralKeyPair()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Tính Session Key
		sessionKey, err := DeriveSessionKey(serverPrivKey, clientPubKeyBytes)
		if err != nil {
			http.Error(w, "Failed to derive session key", http.StatusBadRequest)
			return
		}

		// Lưu vào Session Store
		sessionID, err := store.CreateSession(sessionKey)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(HandshakeResponse{
			ServerPublicKey: hex.EncodeToString(serverPubKeyBytes),
			SessionID:       sessionID,
		})
	}
}

// DecryptMiddleware giải mã Request Body được mã hóa từ Client.
func DecryptMiddleware(store *SessionStore, configs ...SecureConfig) func(http.Handler) http.Handler {
	cfg := GetConfig(configs)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := r.Header.Get(SessionIDHeader)
			if sessionID == "" {
				cfg.OnError(w, r, ErrMissingSessionID, http.StatusBadRequest)
				return
			}

			sessionKey, err := store.GetSession(sessionID)
			if err != nil {
				cfg.OnError(w, r, err, http.StatusUnauthorized)
				return
			}

			// Giới hạn kích thước request body để tránh lỗ hổng DDoS OOM
			r.Body = http.MaxBytesReader(w, r.Body, cfg.MaxBodySize)

			// Đọc Encrypted Request Body
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				cfg.OnError(w, r, err, http.StatusBadRequest)
				return
			}
			r.Body.Close()

			// Nếu body rỗng thì bỏ qua giải mã
			if len(bodyBytes) > 0 {
				ciphertext, err := cfg.Extractor(bodyBytes)
				if err != nil {
					cfg.OnError(w, r, err, http.StatusBadRequest)
					return
				}

				plaintext, err := encrypt.DecryptAESGCM(ciphertext, sessionKey)
				if err != nil {
					cfg.OnError(w, r, err, http.StatusBadRequest)
					return
				}

				// Ghi đè Request Body bằng dữ liệu cleartext đã giải mã
				r.Body = io.NopCloser(bytes.NewBuffer(plaintext))
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(nil))
			}

			// Lưu SessionKey vào Context để các middleware sau (như Encrypt) tái sử dụng
			ctx := context.WithValue(r.Context(), SessionKeyContextKey, sessionKey)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// responseWrapper dùng để chặn dữ liệu ghi ra response và mã hóa nó trước khi gửi đi.
type responseWrapper struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (rw *responseWrapper) Write(b []byte) (int, error) {
	return rw.body.Write(b)
}

func (rw *responseWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

// EncryptMiddleware mã hóa Response Body trước khi gửi về cho Client.
func EncryptMiddleware(store *SessionStore, configs ...SecureConfig) func(http.Handler) http.Handler {
	cfg := GetConfig(configs)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Lấy SessionKey từ Context trước
			var sessionKey []byte
			if val := r.Context().Value(SessionKeyContextKey); val != nil {
				sessionKey = val.([]byte)
			} else {
				// Nếu không có trong Context (do không chạy DecryptMiddleware trước đó), lấy từ Store
				sessionID := r.Header.Get(SessionIDHeader)
				if sessionID == "" {
					cfg.OnError(w, r, ErrMissingSessionID, http.StatusBadRequest)
					return
				}
				var err error
				sessionKey, err = store.GetSession(sessionID)
				if err != nil {
					cfg.OnError(w, r, err, http.StatusUnauthorized)
					return
				}
			}

			// Tái sử dụng bytes.Buffer từ pool để tránh allocation heap
			buf := bufferPool.Get().(*bytes.Buffer)
			buf.Reset()
			defer func() {
				if buf.Cap() <= 65536 {
					bufferPool.Put(buf)
				}
			}()

			wrapper := &responseWrapper{
				ResponseWriter: w,
				body:           buf,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapper, r)

			// Chỉ mã hóa các response thành công (2xx) và có dữ liệu trả về
			if wrapper.statusCode >= 200 && wrapper.statusCode < 300 && wrapper.body.Len() > 0 {
				ciphertext, err := encrypt.AESGCM(wrapper.body.Bytes(), sessionKey)
				if err != nil {
					cfg.OnError(w, r, err, http.StatusInternalServerError)
					return
				}

				respBytes, err := cfg.Wrapper(ciphertext)
				if err != nil {
					cfg.OnError(w, r, err, http.StatusInternalServerError)
					return
				}

				// Xóa Content-Length cũ vì kích thước payload đã thay đổi
				w.Header().Del("Content-Length")
				// Ghi đè Content-Type vì kiểu dữ liệu trả về đã được wrap thành định dạng JSON chuẩn của APIResponse
				w.Header().Set("Content-Type", "application/json")

				w.WriteHeader(wrapper.statusCode)
				_, _ = w.Write(respBytes)
			} else {
				// Nếu lỗi hoặc không có body, trả về như bình thường
				w.WriteHeader(wrapper.statusCode)
				if wrapper.body.Len() > 0 {
					_, _ = w.Write(wrapper.body.Bytes())
				}
			}
		})
	}
}

// CryptoMiddleware gộp cả DecryptMiddleware và EncryptMiddleware.
func CryptoMiddleware(store *SessionStore, configs ...SecureConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		decrypt := DecryptMiddleware(store, configs...)
		encrypt := EncryptMiddleware(store, configs...)
		return decrypt(encrypt(next))
	}
}
