package orchestrator

import (
	"sync"
	"time"

	"go.keploy.io/server/v2/pkg/models"
	"go.uber.org/zap"
)

// DynamicMockQueue provides a non-dropping, elastic queue with backpressure signaling
type DynamicMockQueue struct {
	items  []*models.Mock
	mu     sync.Mutex
	cond   *sync.Cond
	closed bool
	logger *zap.Logger
}

func NewDynamicMockQueue(logger *zap.Logger) *DynamicMockQueue {
	q := &DynamicMockQueue{
		items:  make([]*models.Mock, 0, 1024),
		logger: logger,
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *DynamicMockQueue) Push(mock *models.Mock) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return false
	}
	q.items = append(q.items, mock)
	if len(q.items) > 5000 && len(q.items)%1000 == 0 {
		q.logger.Warn("high mock backlog in correlation queue", zap.Int("backlog_size", len(q.items)))
	}
	q.cond.Signal()
	return true
}

func (q *DynamicMockQueue) Pop(timeout time.Duration) (*models.Mock, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 {
		if q.closed {
			return nil, false
		}
		q.cond.Wait()
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *DynamicMockQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.cond.Broadcast()
}
