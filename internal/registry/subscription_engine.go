package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type SubscriptionEngine struct {
	repo        ISubscriptionRepository
	subscribers map[string][]string
	mu          sync.RWMutex
	notifier    INotifier
}

func NewSubscriptionEngine(repo ISubscriptionRepository) *SubscriptionEngine {
	return &SubscriptionEngine{
		repo:        repo,
		subscribers: make(map[string][]string),
	}
}

func (e *SubscriptionEngine) SetNotifier(notifier INotifier) {
	e.notifier = notifier
}

func (e *SubscriptionEngine) CreateSubscription(ctx context.Context, resourcePath string, params json.RawMessage, maxUpdateRateMs *int, persist *bool, secureWebsocket *bool) (*Subscription, error) {
	sub := &Subscription{
		ID:              uuid.New().String(),
		ResourcePath:    resourcePath,
		Params:          params,
		MaxUpdateRateMs: maxUpdateRateMs,
		Persist:         persist,
		SecureWebsocket: secureWebsocket,
	}

	if err := e.repo.UpsertSubscription(ctx, *sub); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	e.mu.Lock()
	e.subscribers[sub.ID] = []string{sub.ResourcePath}
	e.mu.Unlock()

	return sub, nil
}

func (e *SubscriptionEngine) GetSubscription(ctx context.Context, id string) (*Subscription, error) {
	sub, err := e.repo.GetSubscription(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}
	return &sub, nil
}

func (e *SubscriptionEngine) ListSubscriptions(ctx context.Context) ([]Subscription, error) {
	return e.repo.ListSubscriptions(ctx)
}

func (e *SubscriptionEngine) DeleteSubscription(ctx context.Context, id string) error {
	if err := e.repo.DeleteSubscription(ctx, id); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	e.mu.Lock()
	delete(e.subscribers, id)
	e.mu.Unlock()

	return nil
}

func (e *SubscriptionEngine) Notify(ctx context.Context, resourceType ResourceType, action string, data interface{}) error {
	if e.notifier == nil {
		return nil
	}

	subscriptions, err := e.repo.ListSubscriptions(ctx)
	if err != nil {
		return fmt.Errorf("failed to list subscriptions: %w", err)
	}

	resourcePath := "/" + string(resourceType) + "s"

	for _, sub := range subscriptions {
		if e.matchesSubscription(resourcePath, sub.ResourcePath) {
			if err := e.notifier.Notify(ctx, resourceType, action, data); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *SubscriptionEngine) matchesSubscription(resourcePath, subscriptionPath string) bool {
	return subscriptionPath == resourcePath
}
