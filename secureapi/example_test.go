package secureapi_test

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ducconit/gotoolkit/encrypt"
	"github.com/ducconit/gotoolkit/secureapi"
	securegin "github.com/ducconit/gotoolkit/secureapi/gin"
	"github.com/gin-gonic/gin"
)

// Helper để giả lập phía Client thực hiện bắt tay ECDH
type clientSession struct {
	privKey   *ecdh.PrivateKey
	pubKeyHex string
}

func newClientSession() (*clientSession, error) {
	privKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	pubKeyHex := hex.EncodeToString(privKey.PublicKey().Bytes())
	return &clientSession{
		privKey:   privKey,
		pubKeyHex: pubKeyHex,
	}, nil
}

func (c *clientSession) deriveSessionKey(serverPubKeyHex string) ([]byte, error) {
	serverPubKeyBytes, err := hex.DecodeString(serverPubKeyHex)
	if err != nil {
		return nil, err
	}
	serverPubKey, err := ecdh.P256().NewPublicKey(serverPubKeyBytes)
	if err != nil {
		return nil, err
	}
	sharedSecret, err := c.privKey.ECDH(serverPubKey)
	if err != nil {
		return nil, err
	}
	hasher := sha256.New()
	hasher.Write(sharedSecret)
	return hasher.Sum(nil), nil
}

func TestSecureAPI_Integration(t *testing.T) {
	// 1. Khởi tạo Session Store phía Server với TTL 5 giây
	store := secureapi.NewSessionStore(5 * time.Second)
	defer store.Close()

	// 2. Thiết lập HTTP Mux & Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/handshake", secureapi.HandshakeHandler(store))

	// API nhạy cảm yêu cầu cả hai chiều mã hóa (CryptoMiddleware)
	cryptoMid := secureapi.CryptoMiddleware(store)
	mux.Handle("/api/secure-data", cryptoMid(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}

		// Nhận dữ liệu cleartext đã được middleware giải mã
		var data map[string]any
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Phản hồi dữ liệu nhạy cảm (sẽ được middleware tự động mã hóa trước khi gửi đi)
		response := map[string]any{
			"message": "Hello " + data["username"].(string) + ", secure transaction success!",
			"balance": 5000000,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	})))

	server := httptest.NewServer(mux)
	defer server.Close()

	// 3. GIẢ LẬP CLIENT GIAO TIẾP VỚI SERVER

	// A. Client tạo cặp khóa ECDH của mình
	client, err := newClientSession()
	if err != nil {
		t.Fatalf("failed to create client session: %v", err)
	}

	// B. Gửi yêu cầu Handshake lên Server để trao đổi khóa
	handshakeReq := secureapi.HandshakeRequest{
		ClientPublicKey: client.pubKeyHex,
	}
	reqBody, _ := json.Marshal(handshakeReq)
	resp, err := http.Post(server.URL+"/api/handshake", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("handshake request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("handshake returned status: %d", resp.StatusCode)
	}

	var handshakeResp secureapi.HandshakeResponse
	_ = json.NewDecoder(resp.Body).Decode(&handshakeResp)

	// C. Client tính toán Session Key của mình dựa trên Server Public Key nhận được
	clientSessionKey, err := client.deriveSessionKey(handshakeResp.ServerPublicKey)
	if err != nil {
		t.Fatalf("failed to derive client session key: %v", err)
	}
	sessionID := handshakeResp.SessionID

	// D. Client chuẩn bị dữ liệu nhạy cảm cần gửi đi
	sensitivePayload := []byte(`{"username":"Vi_Cute_20","action":"transfer"}`)

	// Client mã hóa dữ liệu bằng AES-GCM với Session Key
	ciphertext, err := encrypt.AESGCM(sensitivePayload, clientSessionKey)
	if err != nil {
		t.Fatalf("client encryption failed: %v", err)
	}

	// Client đóng gói dữ liệu mã hóa thành JSON
	encryptedReqBody, _ := json.Marshal(secureapi.APIResponse{
		Data: base64.StdEncoding.EncodeToString(ciphertext),
	})

	// E. Client gửi API Request lên API nhạy cảm, đính kèm Session ID
	req, _ := http.NewRequest("POST", server.URL+"/api/secure-data", bytes.NewBuffer(encryptedReqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(secureapi.SessionIDHeader, sessionID)

	httpClient := &http.Client{}
	apiResp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("API request failed: %v", err)
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode != http.StatusOK {
		t.Fatalf("API returned status: %d", apiResp.StatusCode)
	}

	// F. Client nhận phản hồi đã mã hóa từ Server và giải mã
	var encResp secureapi.APIResponse
	_ = json.NewDecoder(apiResp.Body).Decode(&encResp)

	respCiphertext, err := base64.StdEncoding.DecodeString(encResp.Data.(string))
	if err != nil {
		t.Fatalf("failed to decode base64 response: %v", err)
	}

	respPlaintext, err := encrypt.DecryptAESGCM(respCiphertext, clientSessionKey)
	if err != nil {
		t.Fatalf("client decryption of response failed: %v", err)
	}

	// G. Xác nhận kết quả
	var result map[string]any
	_ = json.Unmarshal(respPlaintext, &result)

	expectedMsg := "Hello Vi_Cute_20, secure transaction success!"
	if result["message"] != expectedMsg {
		t.Errorf("expected message '%s', got '%s'", expectedMsg, result["message"])
	}
	if result["balance"].(float64) != 5000000 {
		t.Errorf("expected balance 5000000, got %v", result["balance"])
	}
}

func TestSecureAPI_Gin_Integration(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	// 1. Khởi tạo Session Store phía Server với TTL 5 giây
	store := secureapi.NewSessionStore(5 * time.Second)
	defer store.Close()

	// 2. Thiết lập Gin Router & Routes
	r := gin.New()
	r.POST("/api/handshake", securegin.HandshakeHandler(store))

	// API nhạy cảm yêu cầu cả hai chiều mã hóa (CryptoMiddleware)
	r.POST("/api/secure-data", securegin.CryptoMiddleware(store), func(c *gin.Context) {
		var data map[string]any
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Phản hồi dữ liệu nhạy cảm (sẽ được middleware tự động mã hóa)
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello " + data["username"].(string) + ", secure transaction success via Gin!",
			"balance": 9999999,
		})
	})

	server := httptest.NewServer(r)
	defer server.Close()

	// 3. GIẢ LẬP CLIENT GIAO TIẾP VỚI SERVER GIN

	// A. Client tạo cặp khóa ECDH của mình
	client, err := newClientSession()
	if err != nil {
		t.Fatalf("failed to create client session: %v", err)
	}

	// B. Gửi yêu cầu Handshake lên Server để trao đổi khóa
	handshakeReq := secureapi.HandshakeRequest{
		ClientPublicKey: client.pubKeyHex,
	}
	reqBody, _ := json.Marshal(handshakeReq)
	resp, err := http.Post(server.URL+"/api/handshake", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("handshake request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("handshake returned status: %d", resp.StatusCode)
	}

	var handshakeResp secureapi.HandshakeResponse
	_ = json.NewDecoder(resp.Body).Decode(&handshakeResp)

	// C. Client tính toán Session Key của mình dựa trên Server Public Key nhận được
	clientSessionKey, err := client.deriveSessionKey(handshakeResp.ServerPublicKey)
	if err != nil {
		t.Fatalf("failed to derive client session key: %v", err)
	}
	sessionID := handshakeResp.SessionID

	// D. Client chuẩn bị dữ liệu nhạy cảm cần gửi đi
	sensitivePayload := []byte(`{"username":"Vi_Cute_20","action":"transfer"}`)

	// Client mã hóa dữ liệu
	ciphertext, err := encrypt.AESGCM(sensitivePayload, clientSessionKey)
	if err != nil {
		t.Fatalf("client encryption failed: %v", err)
	}

	// Client đóng gói dữ liệu mã hóa thành JSON
	encryptedReqBody, _ := json.Marshal(secureapi.APIResponse{
		Data: base64.StdEncoding.EncodeToString(ciphertext),
	})

	// E. Client gửi API Request lên API nhạy cảm, đính kèm Session ID
	req, _ := http.NewRequest("POST", server.URL+"/api/secure-data", bytes.NewBuffer(encryptedReqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(secureapi.SessionIDHeader, sessionID)

	httpClient := &http.Client{}
	apiResp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("API request failed: %v", err)
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode != http.StatusOK {
		t.Fatalf("API returned status: %d", apiResp.StatusCode)
	}

	// F. Client nhận phản hồi đã mã hóa từ Server và giải mã
	var encResp secureapi.APIResponse
	_ = json.NewDecoder(apiResp.Body).Decode(&encResp)

	respCiphertext, err := base64.StdEncoding.DecodeString(encResp.Data.(string))
	if err != nil {
		t.Fatalf("failed to decode base64 response: %v", err)
	}

	respPlaintext, err := encrypt.DecryptAESGCM(respCiphertext, clientSessionKey)
	if err != nil {
		t.Fatalf("client decryption of response failed: %v", err)
	}

	// G. Xác nhận kết quả
	var result map[string]any
	_ = json.Unmarshal(respPlaintext, &result)

	expectedMsg := "Hello Vi_Cute_20, secure transaction success via Gin!"
	if result["message"] != expectedMsg {
		t.Errorf("expected message '%s', got '%s'", expectedMsg, result["message"])
	}
	if result["balance"].(float64) != 9999999 {
		t.Errorf("expected balance 9999999, got %v", result["balance"])
	}
}

func TestSecureAPI_ExpiredSession(t *testing.T) {
	// Kiểm tra xem session hết hạn có bị từ chối
	store := secureapi.NewSessionStore(10 * time.Millisecond) // Hết hạn siêu nhanh
	defer store.Close()

	mux := http.NewServeMux()
	cryptoMid := secureapi.CryptoMiddleware(store)
	mux.Handle("/api/secure-data", cryptoMid(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	server := httptest.NewServer(mux)
	defer server.Close()

	// Tạo session
	dummyKey := make([]byte, 32)
	_, _ = rand.Read(dummyKey)
	sessionID, _ := store.CreateSession(dummyKey)

	// Chờ session hết hạn
	time.Sleep(20 * time.Millisecond)

	req, _ := http.NewRequest("POST", server.URL+"/api/secure-data", bytes.NewBuffer([]byte(`{"data":""}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(secureapi.SessionIDHeader, sessionID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Phải trả về 412 hoặc 401 Unauthorized do session đã hết hạn
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d (Unauthorized), got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

// Bảng kiểm thử cho Cipher (Table-driven Tests)
func TestCipher_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		key       []byte
		wantErr   bool
	}{
		{
			name:      "Valid key 32 bytes (AES-256)",
			plaintext: "hello world, this is a secret payload!",
			key:       bytes.Repeat([]byte{0x01}, 32),
			wantErr:   false,
		},
		{
			name:      "Invalid key length 15 bytes",
			plaintext: "secret",
			key:       bytes.Repeat([]byte{0x01}, 15),
			wantErr:   true,
		},
		{
			name:      "Empty plaintext",
			plaintext: "",
			key:       bytes.Repeat([]byte{0x02}, 32),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := encrypt.AESGCM([]byte(tt.plaintext), tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("AESGCM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			plaintext, err := encrypt.DecryptAESGCM(ciphertext, tt.key)
			if err != nil {
				t.Errorf("DecryptAESGCM() failed: %v", err)
				return
			}

			if string(plaintext) != tt.plaintext {
				t.Errorf("DecryptAESGCM() = %s, want %s", string(plaintext), tt.plaintext)
			}
		})
	}
}

// Ví dụ minh họa sử dụng cho tài liệu API
func Example() {
	// Khởi tạo store với thời gian hết hạn 30 phút
	store := secureapi.NewSessionStore(30 * time.Minute)
	defer store.Close()

	// Gắn handshake handler
	http.HandleFunc("/api/handshake", secureapi.HandshakeHandler(store))

	// Gắn API nhạy cảm với middleware
	cryptoMid := secureapi.CryptoMiddleware(store)
	http.Handle("/api/checkout", cryptoMid(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Nhận dữ liệu đã giải mã tự động
		var order map[string]any
		_ = json.NewDecoder(r.Body).Decode(&order)

		// Xử lý thanh toán...

		// Phản hồi (sẽ tự động được mã hóa)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "paid"})
	})))
}
