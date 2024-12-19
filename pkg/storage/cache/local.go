package cache

import (
	"ginboot/pkg/logger"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	// LC local cache client
	LC     *cache.Cache
	lcLock sync.Mutex
)

func InitLC() {
	lcLock.Lock()
	if LC == nil {
		LC = cache.New(5*time.Minute, 10*time.Minute)
	}
	logger.Out.Debug("local cache init finished")
	lcLock.Unlock()
}
