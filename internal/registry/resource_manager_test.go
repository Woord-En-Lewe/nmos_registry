package registry

import (
	"context"
	"testing"
)

func TestResourceManager_Nodes(t *testing.T) {
	repo := NewInMemoryRepository()
	manager := NewResourceManager(repo)
	ctx := context.Background()

	node := Node{
		ID:          "node-1",
		Version:     "v1.0",
		Label:       "Test Node",
		Description: "A test node",
		Href:        "http://localhost:8080",
	}

	// Register
	err := manager.RegisterNode(ctx, node)
	if err != nil {
		t.Fatalf("failed to register node: %v", err)
	}

	// Verify
	savedNode, err := repo.GetNode(ctx, node.ID)
	if err != nil {
		t.Fatalf("failed to get node from repo: %v", err)
	}
	if savedNode.ID != node.ID {
		t.Errorf("expected node ID %s, got %s", node.ID, savedNode.ID)
	}

	// Unregister
	err = manager.UnregisterNode(ctx, node.ID)
	if err != nil {
		t.Fatalf("failed to unregister node: %v", err)
	}

	// Verify deletion
	_, err = repo.GetNode(ctx, node.ID)
	if err == nil {
		t.Error("expected error getting deleted node, got nil")
	}
}

func TestResourceManager_Devices(t *testing.T) {
	repo := NewInMemoryRepository()
	manager := NewResourceManager(repo)
	ctx := context.Background()

	node := Node{ID: "node-1"}
	repo.UpsertNode(ctx, node)

	device := Device{
		ID:          "device-1",
		NodeID:      "node-1",
		Version:     "v1.0",
		Label:       "Test Device",
		Description: "A test device",
	}

	// Register
	err := manager.RegisterDevice(ctx, device)
	if err != nil {
		t.Fatalf("failed to register device: %v", err)
	}

	// Verify
	savedDevice, err := repo.GetDevice(ctx, device.ID)
	if err != nil {
		t.Fatalf("failed to get device from repo: %v", err)
	}
	if savedDevice.ID != device.ID {
		t.Errorf("expected device ID %s, got %s", device.ID, savedDevice.ID)
	}

	// Register without node should fail
	device2 := Device{ID: "device-2", NodeID: "non-existent"}
	err = manager.RegisterDevice(ctx, device2)
	if err == nil {
		t.Error("expected error registering device with non-existent node, got nil")
	}
}

func TestResourceManager_Sources(t *testing.T) {
	repo := NewInMemoryRepository()
	manager := NewResourceManager(repo)
	ctx := context.Background()

	node := Node{ID: "node-1"}
	repo.UpsertNode(ctx, node)
	device := Device{ID: "device-1", NodeID: "node-1"}
	repo.UpsertDevice(ctx, device)

	source := Source{
		ID:          "source-1",
		DeviceID:    "device-1",
		Version:     "v1.0",
		Label:       "Test Source",
		Description: "A test source",
		Format:      "urn:x-nmos:format:video",
	}

	// Register
	err := manager.RegisterSource(ctx, source)
	if err != nil {
		t.Fatalf("failed to register source: %v", err)
	}

	// Verify
	savedSource, err := repo.GetSource(ctx, source.ID)
	if err != nil {
		t.Fatalf("failed to get source from repo: %v", err)
	}
	if savedSource.ID != source.ID {
		t.Errorf("expected source ID %s, got %s", source.ID, savedSource.ID)
	}
}

func TestResourceManager_Flows(t *testing.T) {
	repo := NewInMemoryRepository()
	manager := NewResourceManager(repo)
	ctx := context.Background()

	node := Node{ID: "node-1"}
	repo.UpsertNode(ctx, node)
	device := Device{ID: "device-1", NodeID: "node-1"}
	repo.UpsertDevice(ctx, device)
	source := Source{ID: "source-1", DeviceID: "device-1"}
	repo.UpsertSource(ctx, source)

	flow := Flow{
		ID:          "flow-1",
		SourceID:    "source-1",
		DeviceID:    "device-1",
		Version:     "v1.0",
		Label:       "Test Flow",
		Description: "A test flow",
		Format:      "urn:x-nmos:format:video",
	}

	// Register
	err := manager.RegisterFlow(ctx, flow)
	if err != nil {
		t.Fatalf("failed to register flow: %v", err)
	}

	// Verify
	savedFlow, err := repo.GetFlow(ctx, flow.ID)
	if err != nil {
		t.Fatalf("failed to get flow from repo: %v", err)
	}
	if savedFlow.ID != flow.ID {
		t.Errorf("expected flow ID %s, got %s", flow.ID, savedFlow.ID)
	}
}

func TestResourceManager_Senders(t *testing.T) {
	repo := NewInMemoryRepository()
	manager := NewResourceManager(repo)
	ctx := context.Background()

	node := Node{ID: "node-1"}
	repo.UpsertNode(ctx, node)
	device := Device{ID: "device-1", NodeID: "node-1"}
	repo.UpsertDevice(ctx, device)

	sender := Sender{
		ID:          "sender-1",
		DeviceID:    "device-1",
		Version:     "v1.0",
		Label:       "Test Sender",
		Description: "A test sender",
		Transport:   "urn:x-nmos:transport:rtp",
	}

	// Register
	err := manager.RegisterSender(ctx, sender)
	if err != nil {
		t.Fatalf("failed to register sender: %v", err)
	}

	// Verify
	savedSender, err := repo.GetSender(ctx, sender.ID)
	if err != nil {
		t.Fatalf("failed to get sender from repo: %v", err)
	}
	if savedSender.ID != sender.ID {
		t.Errorf("expected sender ID %s, got %s", sender.ID, savedSender.ID)
	}
}

func TestResourceManager_Receivers(t *testing.T) {
	repo := NewInMemoryRepository()
	manager := NewResourceManager(repo)
	ctx := context.Background()

	node := Node{ID: "node-1"}
	repo.UpsertNode(ctx, node)
	device := Device{ID: "device-1", NodeID: "node-1"}
	repo.UpsertDevice(ctx, device)

	receiver := Receiver{
		ID:          "receiver-1",
		DeviceID:    "device-1",
		Version:     "v1.0",
		Label:       "Test Receiver",
		Description: "A test receiver",
		Transport:   "urn:x-nmos:transport:rtp",
		Format:      "urn:x-nmos:format:video",
	}

	// Register
	err := manager.RegisterReceiver(ctx, receiver)
	if err != nil {
		t.Fatalf("failed to register receiver: %v", err)
	}

	// Verify
	savedReceiver, err := repo.GetReceiver(ctx, receiver.ID)
	if err != nil {
		t.Fatalf("failed to get receiver from repo: %v", err)
	}
	if savedReceiver.ID != receiver.ID {
		t.Errorf("expected receiver ID %s, got %s", receiver.ID, savedReceiver.ID)
	}
}
