package limiter

import (
	"log"
	"sync"
	"time"
)

// depending on rate limit type set default quota if nothing is provided
var defaultQuota = map[string]int64{
	"s": 2,
	"m": 15,
	"h": 100,
}

var defaultDuration = map[string]time.Duration{
	"s": time.Second,
	"m": time.Minute,
	"h": time.Hour,
}

var storeCleaner = cleaner{active: false}

// it is the global limiter store cleaner. at specific interval it scans
// registered stores for invalid keys and removes them
type cleaner struct {
	sync.RWMutex
	active bool
	stores []*MemoryStore
}

// Store interface defines basic functionalities of a store
// type Store interface {
// 	Get(key string) int64
// 	Set(key string, deadline int64) bool
// }

// // NewStore takes limiter config and returns appropriate interface to access store
// func NewStore(conf Config) Store {
// 	switch conf.Store {
// 	case "local":
// 		return NewMemoryStore(conf)
// 	default:
// 		log.Printf("limiter: store %s is not supported. using local store", conf.Store)
// 		return NewMemoryStore(conf)
// 	}
// }

// MemoryStore is the local memory store which holds limiter info
type MemoryStore struct {
	sync.RWMutex
	database    map[string]*Limit
	quota       int64
	limiterType string
}

// NewMemoryStore returns a Memory Store
func NewMemoryStore(conf Config) *MemoryStore {
	s := MemoryStore{
		database:    make(map[string]*Limit),
		quota:       getValidQuota(conf.Quota, conf.Per),
		limiterType: conf.Per,
	}
	return &s
}

func getValidQuota(quota int64, limiterType string) int64 {
	q, found := defaultQuota[limiterType]
	if !found {
		log.Fatal("unsupported rate limit type. should be of (per s/m/h)")
	}
	if quota <= 0 {
		quota = q
	}
	return quota
}

// Get retrieves key from the storage
func (s *MemoryStore) Get(key string) (*Limit, bool) {
	s.RLock()
	defer s.RUnlock()
	limit, found := s.database[key]
	return limit, found
}

// Upsert sets provided limit to the key in the storage
func (s *MemoryStore) Upsert(key string, limit *Limit) {
	s.Lock()
	defer s.Unlock()
	s.database[key] = limit
}

// IsAllowed takes key, returns remaining quota and
// bool stating if the request should be allowed
func (s *MemoryStore) IsAllowed(key string) (int64, bool) {
	limit, _ := s.Get(key)

	// If there is an entry check if deadline is exceeded
	if limit.IsDeadlineExceeded() {
		limit = limit.New(s.limiterType, s.quota)
		limit.Left--
		s.Upsert(key, limit)
		return limit.Left, true
	}

	// Check if there is quota left
	if limit.IsQuotaExceeded() {
		return limit.Left, false
	}

	limit.Left--
	return limit.Left, true
}

func (c *cleaner) addStore(store *MemoryStore) bool {
	c.Lock()
	defer c.Unlock()

	for _, s := range c.stores {
		if s == store {
			log.Println("limiter: (cleanup) we are already tracking this store")
			return false
		}
	}
	c.stores = append(c.stores, store)
	return true
}

func (c *cleaner) getStores() []*MemoryStore {
	c.RLock()
	defer c.RUnlock()

	return c.stores
}

func (c *cleaner) dispatchCleaner() {
	log.Println("limiter: (cleanup) dispatching service")
	if c.active {
		log.Println("limiter: (cleanup) service already running")
		return
	}

	duration, found := defaultDuration["m"]
	if !found {
		log.Fatal("unsupported rate limit type. should be of (per s/m/h)")
	}
	minuteTicker := time.NewTicker(duration)

	go func() {
		for {
			<-minuteTicker.C
			stores := c.getStores()
			for i, s := range stores {
				sTime := time.Now().UnixNano()
				i++
				deletedKeys, totalKeys := doCleanup(s)
				log.Printf(
					"limiter: (cleanup) store[%d/%d] %d/%d keys are deleted in %d ns",
					i, len(stores), deletedKeys, totalKeys, time.Now().UnixNano()-sTime)
			}
		}
	}()

	c.active = true
	log.Println("limiter: (cleanup) service dispatched")
}

func doCleanup(s *MemoryStore) (count, total int) {
	s.Lock()
	defer s.Unlock()

	for key, limit := range s.database {
		if limit.IsDeadlineExceeded() {
			delete(s.database, key)
			count++
		}
		total++
	}

	return
}
