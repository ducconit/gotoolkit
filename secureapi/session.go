package secureapi

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

type sessionItem struct {
	key       []byte
	expiresAt time.Time
}

// SessionStore quản lý các session key một cách an toàn và tự động dọn dẹp khi hết hạn.
type SessionStore struct {
	sessions  sync.Map
	ttl       time.Duration
	stopChan  chan struct{}
	closeOnce sync.Once
}

// NewSessionStore tạo một Store mới với thời gian sống (TTL) cấu hình trước.
func NewSessionStore(ttl time.Duration) *SessionStore {
	store := &SessionStore{
		ttl:      ttl,
		stopChan: make(chan struct{}),
	}
	// Chạy tiến trình dọn dẹp định kỳ mỗi phút
	go store.startCleanup(1 * time.Minute)
	return store
}

// CreateSession tạo một session mới, sinh ngẫu nhiên Session ID và lưu khóa đối xứng.
func (s *SessionStore) CreateSession(key []byte) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	sessionID := hex.EncodeToString(bytes)

	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)

	item := sessionItem{
		key:       keyCopy,
		expiresAt: time.Now().Add(s.ttl),
	}
	s.sessions.Store(sessionID, item)
	return sessionID, nil
}

// GetSession lấy khóa đối xứng tương ứng với Session ID.
// Nếu hết hạn, session sẽ bị xóa và trả về lỗi ErrSessionExpired.
func (s *SessionStore) GetSession(sessionID string) ([]byte, error) {
	val, ok := s.sessions.Load(sessionID)
	if !ok {
		return nil, ErrSessionNotFound
	}

	item := val.(sessionItem)
	if time.Now().After(item.expiresAt) {
		s.sessions.Delete(sessionID)
		return nil, ErrSessionExpired
	}

	keyCopy := make([]byte, len(item.key))
	copy(keyCopy, item.key)
	return keyCopy, nil
}

// DeleteSession xóa session khỏi store một cách chủ động.
func (s *SessionStore) DeleteSession(sessionID string) {
	s.sessions.Delete(sessionID)
}

// Close dừng tiến trình dọn dẹp ngầm (thường dùng khi tắt ứng dụng hoặc trong Unit Test).
func (s *SessionStore) Close() {
	s.closeOnce.Do(func() {
		close(s.stopChan)
	})
}

func (s *SessionStore) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			s.sessions.Range(func(key, val any) bool {
				item := val.(sessionItem)
				if now.After(item.expiresAt) {
					s.sessions.Delete(key)
				}
				return true
			})
		case <-s.stopChan:
			return
		}
	}
}
