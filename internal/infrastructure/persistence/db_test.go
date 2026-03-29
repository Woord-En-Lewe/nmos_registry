package persistence_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	_ "modernc.org/sqlite"
	"github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/persistence"
	"github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/persistence/db"
)

func TestDatabasePersistence(t *testing.T) {
	// 1. Setup in-memory SQLite
	sqlDB, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	defer sqlDB.Close()

	// 2. Initialize schema
	if err := persistence.InitDB(sqlDB); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	queries := db.New(sqlDB)
	ctx := context.Background()

	// 3. Test Node Persistence
	nodeID := "3b84b7a9-9164-4762-a318-3b84b7a99164"
	node := db.CreateNodeParams{
		ID:          nodeID,
		Version:     "1441704316:592391700",
		Label:       "Test Node",
		Description: "A test node for verification",
		Tags:        json.RawMessage(`{"location":["Studio A"]}`),
		Href:        "http://192.168.1.50:8080/",
		Hostname:    sql.NullString{String: "test-node", Valid: true},
		Api:         json.RawMessage(`{"versions":["v1.0","v1.1","v1.2","v1.3"]}`),
		Caps:        json.RawMessage(`{}`),
		Services:    json.RawMessage(`[]`),
		Clocks:      json.RawMessage(`[]`),
		Interfaces:  json.RawMessage(`[]`),
		LastSeen:    sql.NullTime{},
	}

	if err := queries.CreateNode(ctx, node); err != nil {
		t.Fatalf("failed to create node: %v", err)
	}

	retrievedNode, err := queries.GetNode(ctx, nodeID)
	if err != nil {
		t.Fatalf("failed to get node: %v", err)
	}

	if retrievedNode.Label != node.Label {
		t.Errorf("expected label %s, got %s", node.Label, retrievedNode.Label)
	}

	// 4. Test Device Persistence (Foreign Key Constraint)
	deviceID := "7d64b7a9-9164-4762-a318-3b84b7a99164"
	device := db.CreateDeviceParams{
		ID:          deviceID,
		NodeID:      nodeID,
		Version:     "1441704316:592391700",
		Label:       "Test Device",
		Description: "A test device",
		Tags:        json.RawMessage(`{}`),
		Type:        "urn:x-nmos:device:generic",
		Senders:     json.RawMessage(`[]`),
		Receivers:   json.RawMessage(`[]`),
		Controls:    json.RawMessage(`[]`),
	}

	if err := queries.CreateDevice(ctx, device); err != nil {
		t.Fatalf("failed to create device: %v", err)
	}

	retrievedDevice, err := queries.GetDevice(ctx, deviceID)
	if err != nil {
		t.Fatalf("failed to get device: %v", err)
	}

	if retrievedDevice.NodeID != nodeID {
		t.Errorf("expected node_id %s, got %s", nodeID, retrievedDevice.NodeID)
	}

    // 5. Test Source Persistence
    sourceID := "a1b2c3d4-e5f6-4762-a318-3b84b7a99164"
    source := db.CreateSourceParams{
        ID:          sourceID,
        DeviceID:    deviceID,
        Version:     "1441704316:592391700",
        Label:       "Test Source",
        Description: "A test source",
        Tags:        json.RawMessage(`{}`),
        GrainRate:   json.RawMessage(`{"numerator": 25, "denominator": 1}`),
        Format:      "urn:x-nmos:format:video",
        Caps:        json.RawMessage(`{}`),
        Parents:     json.RawMessage(`[]`),
        ClockName:   sql.NullString{String: "clk0", Valid: true},
    }

    if err := queries.CreateSource(ctx, source); err != nil {
        t.Fatalf("failed to create source: %v", err)
    }

    retrievedSource, err := queries.GetSource(ctx, sourceID)
    if err != nil {
        t.Fatalf("failed to get source: %v", err)
    }

    if retrievedSource.DeviceID != deviceID {
        t.Errorf("expected device_id %s, got %s", deviceID, retrievedSource.DeviceID)
    }

	// 6. Test Cascade Delete
	if err := queries.DeleteNode(ctx, nodeID); err != nil {
		t.Fatalf("failed to delete node: %v", err)
	}

	_, err = queries.GetDevice(ctx, deviceID)
	if err != sql.ErrNoRows {
		t.Errorf("expected device to be deleted via cascade, but it still exists or error: %v", err)
	}
}
