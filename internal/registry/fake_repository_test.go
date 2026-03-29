package registry

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type InMemoryRepository struct {
	mu        sync.RWMutex
	nodes     map[string]Node
	devices   map[string]Device
	sources   map[string]Source
	flows     map[string]Flow
	senders   map[string]Sender
	receivers map[string]Receiver
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		nodes:     make(map[string]Node),
		devices:   make(map[string]Device),
		sources:   make(map[string]Source),
		flows:     make(map[string]Flow),
		senders:   make(map[string]Sender),
		receivers: make(map[string]Receiver),
	}
}

// Nodes
func (r *InMemoryRepository) UpsertNode(ctx context.Context, node Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if node.LastSeen.IsZero() {
		node.LastSeen = time.Now()
	}
	r.nodes[node.ID] = node
	return nil
}

func (r *InMemoryRepository) DeleteNode(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.nodes, id)
	return nil
}

func (r *InMemoryRepository) GetNode(ctx context.Context, id string) (Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	node, ok := r.nodes[id]
	if !ok {
		return Node{}, fmt.Errorf("node not found: %s", id)
	}
	return node, nil
}

func (r *InMemoryRepository) ListNodes(ctx context.Context) ([]Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	nodes := make([]Node, 0, len(r.nodes))
	for _, node := range r.nodes {
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (r *InMemoryRepository) UpdateNodeHealth(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	node, ok := r.nodes[id]
	if !ok {
		return fmt.Errorf("node not found: %s", id)
	}
	node.LastSeen = time.Now()
	r.nodes[id] = node
	return nil
}

func (r *InMemoryRepository) GetExpiredNodes(ctx context.Context, expirationTime time.Time) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var expired []string
	for id, node := range r.nodes {
		if node.LastSeen.Before(expirationTime) {
			expired = append(expired, id)
		}
	}
	return expired, nil
}

// Devices
func (r *InMemoryRepository) UpsertDevice(ctx context.Context, device Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.devices[device.ID] = device
	return nil
}

func (r *InMemoryRepository) DeleteDevice(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.devices, id)
	return nil
}

func (r *InMemoryRepository) GetDevice(ctx context.Context, id string) (Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	device, ok := r.devices[id]
	if !ok {
		return Device{}, fmt.Errorf("device not found: %s", id)
	}
	return device, nil
}

func (r *InMemoryRepository) ListDevices(ctx context.Context) ([]Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	devices := make([]Device, 0, len(r.devices))
	for _, device := range r.devices {
		devices = append(devices, device)
	}
	return devices, nil
}

// Sources
func (r *InMemoryRepository) UpsertSource(ctx context.Context, source Source) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sources[source.ID] = source
	return nil
}

func (r *InMemoryRepository) DeleteSource(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sources, id)
	return nil
}

func (r *InMemoryRepository) GetSource(ctx context.Context, id string) (Source, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	source, ok := r.sources[id]
	if !ok {
		return Source{}, fmt.Errorf("source not found: %s", id)
	}
	return source, nil
}

func (r *InMemoryRepository) ListSources(ctx context.Context) ([]Source, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sources := make([]Source, 0, len(r.sources))
	for _, source := range r.sources {
		sources = append(sources, source)
	}
	return sources, nil
}

// Flows
func (r *InMemoryRepository) UpsertFlow(ctx context.Context, flow Flow) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.flows[flow.ID] = flow
	return nil
}

func (r *InMemoryRepository) DeleteFlow(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.flows, id)
	return nil
}

func (r *InMemoryRepository) GetFlow(ctx context.Context, id string) (Flow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	flow, ok := r.flows[id]
	if !ok {
		return Flow{}, fmt.Errorf("flow not found: %s", id)
	}
	return flow, nil
}

func (r *InMemoryRepository) ListFlows(ctx context.Context) ([]Flow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	flows := make([]Flow, 0, len(r.flows))
	for _, flow := range r.flows {
		flows = append(flows, flow)
	}
	return flows, nil
}

// Senders
func (r *InMemoryRepository) UpsertSender(ctx context.Context, sender Sender) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.senders[sender.ID] = sender
	return nil
}

func (r *InMemoryRepository) DeleteSender(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.senders, id)
	return nil
}

func (r *InMemoryRepository) GetSender(ctx context.Context, id string) (Sender, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sender, ok := r.senders[id]
	if !ok {
		return Sender{}, fmt.Errorf("sender not found: %s", id)
	}
	return sender, nil
}

func (r *InMemoryRepository) ListSenders(ctx context.Context) ([]Sender, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	senders := make([]Sender, 0, len(r.senders))
	for _, sender := range r.senders {
		senders = append(senders, sender)
	}
	return senders, nil
}

// Receivers
func (r *InMemoryRepository) UpsertReceiver(ctx context.Context, receiver Receiver) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.receivers[receiver.ID] = receiver
	return nil
}

func (r *InMemoryRepository) DeleteReceiver(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.receivers, id)
	return nil
}

func (r *InMemoryRepository) GetReceiver(ctx context.Context, id string) (Receiver, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	receiver, ok := r.receivers[id]
	if !ok {
		return Receiver{}, fmt.Errorf("receiver not found: %s", id)
	}
	return receiver, nil
}

func (r *InMemoryRepository) ListReceivers(ctx context.Context) ([]Receiver, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	receivers := make([]Receiver, 0, len(r.receivers))
	for _, receiver := range r.receivers {
		receivers = append(receivers, receiver)
	}
	return receivers, nil
}
