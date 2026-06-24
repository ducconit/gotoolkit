package gin

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"
	"sync"

	"github.com/ducconit/gotoolkit/encrypt"
	"github.com/ducconit/gotoolkit/secureapi"
	"github.com/gin-gonic/gin"
)

// GinErrorHandler định nghĩa hàm xử lý lỗi khi xảy ra sự cố trong Gin Middleware.
type GinErrorHandler func(c *gin.Context, err error, statusCode int)

// GinSecureConfig chứa các tùy chỉnh về định dạng dữ liệu và xử lý lỗi dành riêng cho Gin.
type GinSecureConfig struct {
	secureapi.SecureConfig                  // Tái sử dụng Extractor, Wrapper và OnError chuẩn
	OnGinError             GinErrorHandler // Handler xử lý lỗi của ứng dụng
	MaxBodySize            int64           // Giới hạn kích thước request body (bytes) chống DDoS OOM
}

// DefaultOnGinError là handler mặc định cho Gin, trả về JSON dạng APIResponse chuẩn.
func DefaultOnGinError(c *gin.Context, err error, statusCode int) {
	code := "ERROR"
	switch statusCode {
	case http.StatusBadRequest:
		code = "BAD_REQUEST"
	case http.StatusUnauthorized:
		code = "UNAUTHORIZED"
	case http.StatusInternalServerError:
		code = "INTERNAL_ERROR"
	}

	c.AbortWithStatusJSON(statusCode, secureapi.APIResponse{
		Data: nil,
		Code: code,
		Msg:  err.Error(),
	})
}

var (
	// Buffer pool dùng chung để tối ưu hóa bộ nhớ cho response.
	bufferPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 4096))
		},
	}
	sessionKeyCtx   = secureapi.SessionKeyGinContextKey
	sessionIDHeader = secureapi.SessionIDHeader
)

// helper để lấy cấu hình Gin thực tế
func getGinConfig(configs []GinSecureConfig) GinSecureConfig {
	if len(configs) > 0 {
		cfg := configs[0]
		// Điền các giá trị mặc định cho SecureConfig nếu bị thiếu
		if cfg.Extractor == nil {
			cfg.Extractor = secureapi.DefaultConfig.Extractor
		}
		if cfg.Wrapper == nil {
			cfg.Wrapper = secureapi.DefaultConfig.Wrapper
		}
		if cfg.OnError == nil {
			cfg.OnError = secureapi.DefaultConfig.OnError
		}
		if cfg.MaxBodySize <= 0 {
			cfg.MaxBodySize = secureapi.DefaultConfig.MaxBodySize
		}

		// Điền giá trị mặc định cho OnGinError
		if cfg.OnGinError == nil {
			if cfg.OnError != nil && &cfg.OnError != &secureapi.DefaultConfig.OnError {
				// Nếu lập trình viên cấu hình OnError stdlib, ta bọc nó lại thành Gin context
				cfg.OnGinError = func(ctx *gin.Context, err error, statusCode int) {
					ctx.Abort()
					cfg.OnError(ctx.Writer, ctx.Request, err, statusCode)
				}
			} else {
				cfg.OnGinError = DefaultOnGinError
			}
		}
		return cfg
	}

	return GinSecureConfig{
		SecureConfig: secureapi.DefaultConfig,
		OnGinError:   DefaultOnGinError,
		MaxBodySize:  secureapi.DefaultConfig.MaxBodySize,
	}
}

// ginBodyWriter dùng để chặn dữ liệu ghi ra response của Gin.
type ginBodyWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *ginBodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *ginBodyWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func (w *ginBodyWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *ginBodyWriter) WriteHeaderNow() {
	// Không cho phép ghi đè header xuống socket sớm
}

func (w *ginBodyWriter) Status() int {
	if w.statusCode == 0 {
		return w.ResponseWriter.Status()
	}
	return w.statusCode
}

func (w *ginBodyWriter) Size() int {
	return w.body.Len()
}

func (w *ginBodyWriter) Written() bool {
	return w.statusCode != 0 || w.body.Len() > 0
}

// HandshakeHandler xử lý endpoint bắt tay trao đổi khóa cho Gin.
func HandshakeHandler(store *secureapi.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Giới hạn body cho handshake
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 64*1024)

		var req secureapi.HandshakeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		clientPubKeyBytes, err := hex.DecodeString(req.ClientPublicKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client public key format"})
			return
		}

		// Sinh khóa tạm thời cho Server
		serverPrivKey, serverPubKeyBytes, err := secureapi.GenerateEphemeralKeyPair()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Tính Session Key
		sessionKey, err := secureapi.DeriveSessionKey(serverPrivKey, clientPubKeyBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to derive session key"})
			return
		}

		// Lưu vào Session Store
		sessionID, err := store.CreateSession(sessionKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, secureapi.HandshakeResponse{
			ServerPublicKey: hex.EncodeToString(serverPubKeyBytes),
			SessionID:       sessionID,
		})
	}
}

// decryptHandler là helper xử lý logic giải mã dùng chung để tránh trùng lặp code.
func decryptHandler(c *gin.Context, store *secureapi.SessionStore, cfg GinSecureConfig) ([]byte, bool) {
	sessionID := c.GetHeader(sessionIDHeader)
	if sessionID == "" {
		cfg.OnGinError(c, secureapi.ErrMissingSessionID, http.StatusBadRequest)
		return nil, false
	}

	sessionKey, err := store.GetSession(sessionID)
	if err != nil {
		cfg.OnGinError(c, err, http.StatusUnauthorized)
		return nil, false
	}

	// Giới hạn kích thước request body chống DDoS OOM
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, cfg.MaxBodySize)

	// Đọc payload mã hóa từ Request Body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		cfg.OnGinError(c, err, http.StatusBadRequest)
		return nil, false
	}
	c.Request.Body.Close()

	if len(bodyBytes) > 0 {
		ciphertext, err := cfg.Extractor(bodyBytes)
		if err != nil {
			cfg.OnGinError(c, err, http.StatusBadRequest)
			return nil, false
		}

		plaintext, err := encrypt.DecryptAESGCM(ciphertext, sessionKey)
		if err != nil {
			cfg.OnGinError(c, err, http.StatusBadRequest)
			return nil, false
		}

		// Ghi đè Request Body bằng dữ liệu cleartext đã giải mã
		c.Request.Body = io.NopCloser(bytes.NewBuffer(plaintext))
	} else {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(nil))
	}

	// Lưu SessionKey vào Context để các middleware sau (như Encrypt) tái sử dụng
	c.Set(sessionKeyCtx, sessionKey)
	return sessionKey, true
}

// encryptHandler là helper xử lý logic mã hóa dùng chung để tránh trùng lặp code.
func encryptHandler(c *gin.Context, store *secureapi.SessionStore, cfg GinSecureConfig, sessionKey []byte) {
	if len(sessionKey) == 0 {
		if val, exists := c.Get(sessionKeyCtx); exists {
			sessionKey = val.([]byte)
		} else {
			sessionID := c.GetHeader(sessionIDHeader)
			if sessionID == "" {
				cfg.OnGinError(c, secureapi.ErrMissingSessionID, http.StatusBadRequest)
				return
			}
			var err error
			sessionKey, err = store.GetSession(sessionID)
			if err != nil {
				cfg.OnGinError(c, err, http.StatusUnauthorized)
				return
			}
		}
	}

	// Tái sử dụng bytes.Buffer từ pool để tối ưu hiệu năng
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer func() {
		if buf.Cap() <= 65536 {
			bufferPool.Put(buf)
		}
	}()

	// Thay thế ResponseWriter của Gin để bắt dữ liệu trả về
	w := &ginBodyWriter{body: buf, ResponseWriter: c.Writer, statusCode: http.StatusOK}
	c.Writer = w

	c.Next()

	// Chỉ mã hóa response thành công (2xx) và có dữ liệu trả về
	statusCode := w.Status()
	if statusCode >= 200 && statusCode < 300 && w.body.Len() > 0 {
		ciphertext, err := encrypt.AESGCM(w.body.Bytes(), sessionKey)
		if err != nil {
			w.ResponseWriter.WriteHeader(http.StatusInternalServerError)
			_, _ = w.ResponseWriter.Write([]byte(`{"error":"Failed to encrypt response payload"}`))
			return
		}

		respBytes, err := cfg.Wrapper(ciphertext)
		if err != nil {
			w.ResponseWriter.WriteHeader(http.StatusInternalServerError)
			_, _ = w.ResponseWriter.Write([]byte(`{"error":"Failed to wrap response payload: ` + err.Error() + `"}`))
			return
		}

		// Xóa Content-Length cũ vì kích thước payload đã thay đổi
		w.ResponseWriter.Header().Del("Content-Length")
		// Ghi đè Content-Type vì kiểu dữ liệu trả về đã được wrap thành định dạng JSON chuẩn của APIResponse
		w.ResponseWriter.Header().Set("Content-Type", "application/json")

		w.ResponseWriter.WriteHeader(statusCode)
		_, _ = w.ResponseWriter.Write(respBytes)
	} else {
		w.ResponseWriter.WriteHeader(statusCode)
		if w.body.Len() > 0 {
			_, _ = w.ResponseWriter.Write(w.body.Bytes())
		}
	}
}

// DecryptMiddleware giải mã Request Body được mã hóa từ Client gửi lên cho Gin.
func DecryptMiddleware(store *secureapi.SessionStore, configs ...GinSecureConfig) gin.HandlerFunc {
	cfg := getGinConfig(configs)
	return func(c *gin.Context) {
		if _, ok := decryptHandler(c, store, cfg); ok {
			c.Next()
		}
	}
}

// EncryptMiddleware mã hóa Response Body trước khi gửi về cho Client trong Gin.
func EncryptMiddleware(store *secureapi.SessionStore, configs ...GinSecureConfig) gin.HandlerFunc {
	cfg := getGinConfig(configs)
	return func(c *gin.Context) {
		encryptHandler(c, store, cfg, nil)
	}
}

// CryptoMiddleware gộp cả DecryptMiddleware và EncryptMiddleware dành cho Gin.
func CryptoMiddleware(store *secureapi.SessionStore, configs ...GinSecureConfig) gin.HandlerFunc {
	cfg := getGinConfig(configs)
	return func(c *gin.Context) {
		if sessionKey, ok := decryptHandler(c, store, cfg); ok {
			encryptHandler(c, store, cfg, sessionKey)
		}
	}
}
