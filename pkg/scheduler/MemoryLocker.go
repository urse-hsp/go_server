// MemoryLocker.go 内存锁
// 进程内锁（单服务锁）

package scheduler

import (
	"context"
	"sync"
	"time"
)

type MemoryLocker struct {
	mu   sync.Mutex
	lock map[string]time.Time
}

func NewMemoryLocker() *MemoryLocker {
	return &MemoryLocker{
		lock: make(map[string]time.Time),
	}
}

func (m *MemoryLocker) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	// 判断是否存在且未过期
	if expire, ok := m.lock[key]; ok && expire.After(now) {
		return false, nil
	}

	m.lock[key] = now.Add(ttl)
	return true, nil
}

func (m *MemoryLocker) Unlock(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.lock, key)
	return nil
}
