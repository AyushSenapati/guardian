package limiter

import (
	"log"
	"time"
)

// Limit is the rate limit DS
type Limit struct {
	Deadline int64
	Left     int64
}

func getDeadline(limiterType string) (deadline int64) {
	switch limiterType {
	case "s":
		deadline = time.Now().Add(time.Second).UnixNano()
	case "m":
		deadline = time.Now().Add(time.Minute).UnixNano()
	case "h":
		deadline = time.Now().Add(time.Hour).UnixNano()
	default:
		log.Fatal("unsupported rate limit type. should be of (per s/m/h)")
	}
	return deadline
}

// New returns new instance of Limit with the provided deadline and quota
func (l *Limit) New(limiterType string, quota int64) *Limit {
	return &Limit{Deadline: getDeadline(limiterType), Left: quota}
}

// IsDeadlineExceeded checks if deadline is exeeded
func (l *Limit) IsDeadlineExceeded() bool {
	if l == nil || l.Deadline < time.Now().UnixNano() {
		return true
	}
	return false
}

// IsQuotaExceeded checks if quota is left
func (l *Limit) IsQuotaExceeded() bool {
	if l.Left <= 0 {
		return true
	}
	return false
}

func (l *Limit) resetLimit() bool {
	return true
}
