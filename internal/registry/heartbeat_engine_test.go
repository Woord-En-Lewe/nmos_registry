package registry

import (
	"context"
	"testing"
	"time"
)

func TestHeartbeatEngine_Heartbeat(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewHeartbeatEngine(repo, 0, 0)
	ctx := context.Background()

	nodeID := "node-1"
	node := Node{ID: nodeID}
	repo.UpsertNode(ctx, node)

	// Get initial LastSeen
	n, _ := repo.GetNode(ctx, nodeID)
	initialLastSeen := n.LastSeen

	// Small sleep to ensure time moves
	time.Sleep(10 * time.Millisecond)

	err := engine.Heartbeat(ctx, nodeID)
	if err != nil {
		t.Fatalf("Heartbeat failed: %v", err)
	}

	n, _ = repo.GetNode(ctx, nodeID)
	if !n.LastSeen.After(initialLastSeen) {
		t.Errorf("Expected LastSeen to be updated, but it wasn't")
	}
}

func TestHeartbeatEngine_PerformCleanup(t *testing.T) {
	repo := NewInMemoryRepository()
	// Set a very short expiry for testing
	expiry := 50 * time.Millisecond
	engine := NewHeartbeatEngine(repo, expiry, 0)
	ctx := context.Background()

	node1 := "node-1"
	node2 := "node-2"
	repo.UpsertNode(ctx, Node{ID: node1})
	repo.UpsertNode(ctx, Node{ID: node2})

	// Heartbeat node 2 so it stays fresh
	time.Sleep(30 * time.Millisecond)
	engine.Heartbeat(ctx, node2)

	// Wait for node 1 to expire but node 2 to remain fresh
	time.Sleep(30 * time.Millisecond)

	engine.PerformCleanup(ctx)

	// Node 1 should be gone
	_, err := repo.GetNode(ctx, node1)
	if err == nil {
		t.Errorf("Expected node 1 to be deleted, but it still exists")
	}

	// Node 2 should still exist
	_, err = repo.GetNode(ctx, node2)
	if err != nil {
		t.Errorf("Expected node 2 to still exist, but it was deleted: %v", err)
	}
}
