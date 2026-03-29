package registry

import (
	"context"
	"log"
	"time"
)

type HeartbeatEngine struct {
	repo            IRepository
	expiryDuration  time.Duration
	cleanupInterval time.Duration
}

func NewHeartbeatEngine(repo IRepository, expiryDuration time.Duration, cleanupInterval time.Duration) *HeartbeatEngine {
	if expiryDuration == 0 {
		expiryDuration = 12 * time.Second // IS-04 recommendation is often > 5s interval, e.g. 12s
	}
	if cleanupInterval == 0 {
		cleanupInterval = 5 * time.Second
	}
	return &HeartbeatEngine{
		repo:            repo,
		expiryDuration:  expiryDuration,
		cleanupInterval: cleanupInterval,
	}
}

func (e *HeartbeatEngine) Heartbeat(ctx context.Context, nodeID string) error {
	return e.repo.UpdateNodeHealth(ctx, nodeID)
}

func (e *HeartbeatEngine) StartGarbageCollector(ctx context.Context) {
	ticker := time.NewTicker(e.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.PerformCleanup(ctx)
		}
	}
}

func (e *HeartbeatEngine) PerformCleanup(ctx context.Context) {
	expirationTime := time.Now().Add(-e.expiryDuration)
	expiredNodes, err := e.repo.GetExpiredNodes(ctx, expirationTime)
	if err != nil {
		log.Printf("Error getting expired nodes: %v", err)
		return
	}

	for _, nodeID := range expiredNodes {
		err := e.repo.DeleteNode(ctx, nodeID)
		if err != nil {
			log.Printf("Error deleting expired node %s: %v", nodeID, err)
		}
	}
}
