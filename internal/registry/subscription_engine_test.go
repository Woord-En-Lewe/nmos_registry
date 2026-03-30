package registry

import (
	"context"
	"encoding/json"
	"testing"
)

type mockNotifier struct {
	notifications []Notification
}

type Notification struct {
	ResourceType ResourceType
	Action       string
	Data         interface{}
}

func (m *mockNotifier) Notify(ctx context.Context, resourceType ResourceType, action string, data interface{}) error {
	m.notifications = append(m.notifications, Notification{
		ResourceType: resourceType,
		Action:       action,
		Data:         data,
	})
	return nil
}

func TestSubscriptionEngine_CreateSubscription(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewSubscriptionEngine(repo)
	ctx := context.Background()

	params := json.RawMessage(`{"key": "value"}`)
	maxUpdateRate := 100
	persist := true
	secureWs := false

	sub, err := engine.CreateSubscription(ctx, "/nodes", params, &maxUpdateRate, &persist, &secureWs)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	if sub.ID == "" {
		t.Errorf("Expected subscription ID to be set")
	}
	if sub.ResourcePath != "/nodes" {
		t.Errorf("Expected resource_path to be /nodes, got %s", sub.ResourcePath)
	}
	if string(sub.Params) != string(params) {
		t.Errorf("Expected params to match")
	}
}

func TestSubscriptionEngine_GetSubscription(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewSubscriptionEngine(repo)
	ctx := context.Background()

	sub, err := engine.CreateSubscription(ctx, "/devices", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	retrieved, err := engine.GetSubscription(ctx, sub.ID)
	if err != nil {
		t.Fatalf("GetSubscription failed: %v", err)
	}

	if retrieved.ID != sub.ID {
		t.Errorf("Expected ID to match")
	}
	if retrieved.ResourcePath != "/devices" {
		t.Errorf("Expected resource_path to be /devices")
	}
}

func TestSubscriptionEngine_ListSubscriptions(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewSubscriptionEngine(repo)
	ctx := context.Background()

	_, err := engine.CreateSubscription(ctx, "/nodes", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}
	_, err = engine.CreateSubscription(ctx, "/devices", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	subs, err := engine.ListSubscriptions(ctx)
	if err != nil {
		t.Fatalf("ListSubscriptions failed: %v", err)
	}

	if len(subs) != 2 {
		t.Errorf("Expected 2 subscriptions, got %d", len(subs))
	}
}

func TestSubscriptionEngine_DeleteSubscription(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewSubscriptionEngine(repo)
	ctx := context.Background()

	sub, err := engine.CreateSubscription(ctx, "/nodes", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	err = engine.DeleteSubscription(ctx, sub.ID)
	if err != nil {
		t.Fatalf("DeleteSubscription failed: %v", err)
	}

	_, err = engine.GetSubscription(ctx, sub.ID)
	if err == nil {
		t.Errorf("Expected error when getting deleted subscription")
	}
}

func TestSubscriptionEngine_Notify(t *testing.T) {
	repo := NewInMemoryRepository()
	engine := NewSubscriptionEngine(repo)
	notifier := &mockNotifier{}
	engine.SetNotifier(notifier)
	ctx := context.Background()

	_, err := engine.CreateSubscription(ctx, "/nodes", nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}

	err = engine.Notify(ctx, ResourceTypeNode, "create", map[string]string{"id": "node-1"})
	if err != nil {
		t.Fatalf("Notify failed: %v", err)
	}

	if len(notifier.notifications) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(notifier.notifications))
	}

	if notifier.notifications[0].ResourceType != ResourceTypeNode {
		t.Errorf("Expected resource type to be node")
	}
	if notifier.notifications[0].Action != "create" {
		t.Errorf("Expected action to be create")
	}
}
