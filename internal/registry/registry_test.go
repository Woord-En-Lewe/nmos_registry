package registry

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"
)

type testNotifier struct {
	mu            sync.Mutex
	notifications []testNotification
}

type testNotification struct {
	ResourceType ResourceType
	Action       string
	Data         interface{}
}

func (m *testNotifier) Notify(ctx context.Context, resourceType ResourceType, action string, data interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notifications = append(m.notifications, testNotification{
		ResourceType: resourceType,
		Action:       action,
		Data:         data,
	})
	return nil
}

func (m *testNotifier) GetNotifications() []testNotification {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.notifications
}

func (m *testNotifier) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notifications = nil
}

func TestResourceManager_ReferentialIntegrity(t *testing.T) {
	repo := NewInMemoryRepository()
	rm := NewResourceManager(repo)
	ctx := context.Background()

	node := Node{
		ID:          "node-1",
		Version:     "1.0",
		Label:       "Test Node",
		Description: "A test node",
		Api:         json.RawMessage(`{}`),
	}

	result := rm.RegisterNode(ctx, node)
	if result.IsFailure() {
		err, _ := result.Error()
		t.Fatalf("RegisterNode failed: %v", err)
	}

	device := Device{
		ID:        "device-1",
		NodeID:    "node-1",
		Version:   "1.0",
		Label:     "Test Device",
		Type:      "urn:x-nmos:device:generic",
		Senders:   json.RawMessage(`[]`),
		Receivers: json.RawMessage(`[]`),
	}

	deviceResult := rm.RegisterDevice(ctx, device)
	if deviceResult.IsFailure() {
		err, _ := deviceResult.Error()
		t.Fatalf("RegisterDevice failed: %v", err)
	}

	device2 := Device{
		ID:        "device-2",
		NodeID:    "nonexistent-node",
		Version:   "1.0",
		Label:     "Orphan Device",
		Type:      "urn:x-nmos:device:generic",
		Senders:   json.RawMessage(`[]`),
		Receivers: json.RawMessage(`[]`),
	}

	device2Result := rm.RegisterDevice(ctx, device2)
	if device2Result.IsSuccess() {
		t.Errorf("RegisterDevice should fail for non-existent parent node")
	}
}

func TestResourceManager_ResourceOrdering(t *testing.T) {
	repo := NewInMemoryRepository()
	rm := NewResourceManager(repo)
	ctx := context.Background()

	node := Node{
		ID:          "node-1",
		Version:     "1.0",
		Label:       "Test Node",
		Description: "A test node",
		Api:         json.RawMessage(`{}`),
	}

	device := Device{
		ID:        "device-1",
		NodeID:    "node-1",
		Version:   "1.0",
		Label:     "Test Device",
		Type:      "urn:x-nmos:device:generic",
		Senders:   json.RawMessage(`[]`),
		Receivers: json.RawMessage(`[]`),
	}

	source := Source{
		ID:          "source-1",
		DeviceID:    "device-1",
		Version:     "1.0",
		Label:       "Test Source",
		Description: "A test source",
		Format:      "urn:x-nmos:format:video",
		GrainRate:   json.RawMessage(`{"numerator": 25, "denominator": 1}`),
	}

	flow := Flow{
		ID:          "flow-1",
		SourceID:    "source-1",
		DeviceID:    "device-1",
		Version:     "1.0",
		Label:       "Test Flow",
		Description: "A test flow",
		Format:      "urn:x-nmos:format:video",
	}

	if nodeResult := rm.RegisterNode(ctx, node); nodeResult.IsFailure() {
		err, _ := nodeResult.Error()
		t.Fatalf("RegisterNode failed: %v", err)
	}
	if deviceResult := rm.RegisterDevice(ctx, device); deviceResult.IsFailure() {
		err, _ := deviceResult.Error()
		t.Fatalf("RegisterDevice failed: %v", err)
	}
	if sourceResult := rm.RegisterSource(ctx, source); sourceResult.IsFailure() {
		err, _ := sourceResult.Error()
		t.Fatalf("RegisterSource failed: %v", err)
	}
	if flowResult := rm.RegisterFlow(ctx, flow); flowResult.IsFailure() {
		err, _ := flowResult.Error()
		t.Fatalf("RegisterFlow failed: %v", err)
	}

	nodes, _ := repo.ListNodes(ctx)
	devices, _ := repo.ListDevices(ctx)
	sources, _ := repo.ListSources(ctx)
	flows, _ := repo.ListFlows(ctx)

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}
	if len(devices) != 1 {
		t.Errorf("Expected 1 device, got %d", len(devices))
	}
	if len(sources) != 1 {
		t.Errorf("Expected 1 source, got %d", len(sources))
	}
	if len(flows) != 1 {
		t.Errorf("Expected 1 flow, got %d", len(flows))
	}
}

func TestHeartbeatEngine_GarbageCollection(t *testing.T) {
	repo := NewInMemoryRepository()
	expiry := 50 * time.Millisecond
	engine := NewHeartbeatEngine(repo, expiry, 0)
	ctx := context.Background()

	node1 := Node{ID: "node-1", Api: json.RawMessage(`{}`)}
	node2 := Node{ID: "node-2", Api: json.RawMessage(`{}`)}
	repo.UpsertNode(ctx, node1)
	repo.UpsertNode(ctx, node2)

	time.Sleep(30 * time.Millisecond)
	engine.Heartbeat(ctx, "node-2")

	time.Sleep(30 * time.Millisecond)
	engine.PerformCleanup(ctx)

	_, err := repo.GetNode(ctx, "node-1")
	if err == nil {
		t.Errorf("Expected node-1 to be garbage collected")
	}

	_, err = repo.GetNode(ctx, "node-2")
	if err != nil {
		t.Errorf("Expected node-2 to still exist")
	}
}

func TestSubscriptionEngine_Notifications(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewSubscriptionEngine(repo)
	notifier := &testNotifier{}
	engine.SetNotifier(notifier)
	ctx := context.Background()

	_, err := engine.CreateSubscription(ctx, "/nodes", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	notifier.Clear()
	err = engine.Notify(ctx, ResourceTypeNode, "create", map[string]string{"id": "node-1"})
	if err != nil {
		t.Fatalf("Notify failed: %v", err)
	}

	notifications := notifier.GetNotifications()
	if len(notifications) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(notifications))
	}

	if notifications[0].ResourceType != ResourceTypeNode {
		t.Errorf("Expected resource type node, got %s", notifications[0].ResourceType)
	}
	if notifications[0].Action != "create" {
		t.Errorf("Expected action create, got %s", notifications[0].Action)
	}
}

func TestSubscriptionEngine_WildcardSubscriptions(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewSubscriptionEngine(repo)
	notifier := &testNotifier{}
	engine.SetNotifier(notifier)
	ctx := context.Background()

	_, err := engine.CreateSubscription(ctx, "/nodes", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	notifier.Clear()
	engine.Notify(ctx, ResourceTypeDevice, "create", map[string]string{"id": "device-1"})

	notifications := notifier.GetNotifications()
	if len(notifications) != 0 {
		t.Errorf("Expected 0 notifications for /nodes subscription receiving /devices update, got %d", len(notifications))
	}
}

func TestSubscriptionEngine_PersistFlag(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewSubscriptionEngine(repo)
	ctx := context.Background()

	persist := true
	sub, err := engine.CreateSubscription(ctx, "/nodes", nil, nil, &persist, nil)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	retrieved, err := engine.GetSubscription(ctx, sub.ID)
	if err != nil {
		t.Fatalf("GetSubscription failed: %v", err)
	}

	if retrieved.Persist == nil || !*retrieved.Persist {
		t.Errorf("Expected persist to be true")
	}
}

func TestSubscriptionEngine_MultipleSubscriptions(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewSubscriptionEngine(repo)
	notifier := &testNotifier{}
	engine.SetNotifier(notifier)
	ctx := context.Background()

	engine.CreateSubscription(ctx, "/nodes", nil, nil, nil, nil)
	engine.CreateSubscription(ctx, "/devices", nil, nil, nil, nil)

	notifier.Clear()
	engine.Notify(ctx, ResourceTypeNode, "create", map[string]string{"id": "node-1"})
	engine.Notify(ctx, ResourceTypeDevice, "create", map[string]string{"id": "device-1"})

	notifications := notifier.GetNotifications()
	if len(notifications) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(notifications))
	}
}

func TestSubscriptionEngine_Delete(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewSubscriptionEngine(repo)
	ctx := context.Background()

	sub, _ := engine.CreateSubscription(ctx, "/nodes", nil, nil, nil, nil)

	err := engine.DeleteSubscription(ctx, sub.ID)
	if err != nil {
		t.Fatalf("DeleteSubscription failed: %v", err)
	}

	_, err = engine.GetSubscription(ctx, sub.ID)
	if err == nil {
		t.Errorf("Expected error getting deleted subscription")
	}
}

func TestIntegratedScenario_HeartbeatKeepsNodeAlive(t *testing.T) {
	repo := NewInMemoryRepository()
	rm := NewResourceManager(repo)
	engine := NewHeartbeatEngine(repo, 100*time.Millisecond, 0)
	ctx := context.Background()

	node := Node{
		ID:          "node-1",
		Version:     "1.0",
		Label:       "Test Node",
		Description: "A test node",
		Api:         json.RawMessage(`{}`),
	}
	rm.RegisterNode(ctx, node)

	time.Sleep(50 * time.Millisecond)
	engine.Heartbeat(ctx, "node-1")

	time.Sleep(50 * time.Millisecond)
	engine.PerformCleanup(ctx)

	_, err := repo.GetNode(ctx, "node-1")
	if err != nil {
		t.Errorf("Expected node to still exist after heartbeat")
	}
}

func TestIntegratedScenario_ExpiredNodeGarbageCollected(t *testing.T) {
	repo := NewInMemoryRepository()
	rm := NewResourceManager(repo)
	engine := NewHeartbeatEngine(repo, 50*time.Millisecond, 0)
	ctx := context.Background()

	node := Node{
		ID:          "node-1",
		Version:     "1.0",
		Label:       "Test Node",
		Description: "A test node",
		Api:         json.RawMessage(`{}`),
	}
	rm.RegisterNode(ctx, node)

	time.Sleep(120 * time.Millisecond)
	engine.PerformCleanup(ctx)

	_, err := repo.GetNode(ctx, "node-1")
	if err == nil {
		t.Errorf("Expected node to be garbage collected after expiration")
	}
}

func TestIntegratedScenario_ChildResourcesDeletedWithParent(t *testing.T) {
	repo := NewInMemoryRepository()
	rm := NewResourceManager(repo)
	ctx := context.Background()

	node := Node{
		ID:          "node-1",
		Version:     "1.0",
		Label:       "Test Node",
		Description: "A test node",
		Api:         json.RawMessage(`{}`),
	}
	rm.RegisterNode(ctx, node)

	device := Device{
		ID:        "device-1",
		NodeID:    "node-1",
		Version:   "1.0",
		Label:     "Test Device",
		Type:      "urn:x-nmos:device:generic",
		Senders:   json.RawMessage(`[]`),
		Receivers: json.RawMessage(`[]`),
	}
	rm.RegisterDevice(ctx, device)

	rm.UnregisterNode(ctx, "node-1")

	devices, _ := repo.ListDevices(ctx)
	if len(devices) != 0 {
		t.Errorf("Expected 0 devices after parent node deletion, got %d", len(devices))
	}
}

func TestIntegratedScenario_SubscriptionLifecycle(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewSubscriptionEngine(repo)
	notifier := &testNotifier{}
	engine.SetNotifier(notifier)
	ctx := context.Background()

	sub1, err := engine.CreateSubscription(ctx, "/nodes", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	sub2, err := engine.CreateSubscription(ctx, "/devices", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	subs, _ := engine.ListSubscriptions(ctx)
	if len(subs) != 2 {
		t.Errorf("Expected 2 subscriptions, got %d", len(subs))
	}

	err = engine.DeleteSubscription(ctx, sub1.ID)
	if err != nil {
		t.Fatalf("DeleteSubscription failed: %v", err)
	}

	subs, _ = engine.ListSubscriptions(ctx)
	if len(subs) != 1 {
		t.Errorf("Expected 1 subscription after deletion, got %d", len(subs))
	}

	_, err = engine.GetSubscription(ctx, sub2.ID)
	if err != nil {
		t.Errorf("Expected sub2 to still exist")
	}
}

func TestIntegratedScenario_SenderWithFlow(t *testing.T) {
	repo := NewInMemoryRepository()
	rm := NewResourceManager(repo)
	ctx := context.Background()

	node := Node{
		ID:          "node-1",
		Version:     "1.0",
		Label:       "Test Node",
		Description: "A test node",
		Api:         json.RawMessage(`{}`),
	}
	rm.RegisterNode(ctx, node)

	device := Device{
		ID:        "device-1",
		NodeID:    "node-1",
		Version:   "1.0",
		Label:     "Test Device",
		Type:      "urn:x-nmos:device:generic",
		Senders:   json.RawMessage(`[]`),
		Receivers: json.RawMessage(`[]`),
	}
	rm.RegisterDevice(ctx, device)

	source := Source{
		ID:          "source-1",
		DeviceID:    "device-1",
		Version:     "1.0",
		Label:       "Test Source",
		Description: "A test source",
		Format:      "urn:x-nmos:format:video",
		GrainRate:   json.RawMessage(`{"numerator": 25, "denominator": 1}`),
	}
	rm.RegisterSource(ctx, source)

	flow := Flow{
		ID:          "flow-1",
		SourceID:    "source-1",
		DeviceID:    "device-1",
		Version:     "1.0",
		Label:       "Test Flow",
		Description: "A test flow",
		Format:      "urn:x-nmos:format:video",
	}
	rm.RegisterFlow(ctx, flow)

	sender := Sender{
		ID:        "sender-1",
		DeviceID:  "device-1",
		FlowID:    strPtr("flow-1"),
		Version:   "1.0",
		Label:     "Test Sender",
		Transport: "urn:x-nmos:transport:rtp",
	}
	senderResult := rm.RegisterSender(ctx, sender)
	if senderResult.IsFailure() {
		err, _ := senderResult.Error()
		t.Fatalf("RegisterSender failed: %v", err)
	}

	retrieved, _ := repo.GetSender(ctx, "sender-1")
	if retrieved.FlowID == nil || *retrieved.FlowID != "flow-1" {
		t.Errorf("Expected sender to reference flow-1")
	}
}

func TestIntegratedScenario_ReceiverRequiresDevice(t *testing.T) {
	repo := NewInMemoryRepository()
	rm := NewResourceManager(repo)
	ctx := context.Background()

	receiver := Receiver{
		ID:        "receiver-1",
		DeviceID:  "nonexistent-device",
		Version:   "1.0",
		Label:     "Test Receiver",
		Transport: "urn:x-nmos:transport:rtp",
		Format:    "urn:x-nmos:format:video",
	}
	result := rm.RegisterReceiver(ctx, receiver)
	if result.IsSuccess() {
		t.Errorf("Expected error registering receiver without parent device")
	}
}

func strPtr(s string) *string {
	return &s
}
