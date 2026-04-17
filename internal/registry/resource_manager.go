package registry

import (
	"context"
	"fmt"
)

type ResourceManager struct {
	repo IRepository
}

func NewResourceManager(repo IRepository) *ResourceManager {
	return &ResourceManager{
		repo: repo,
	}
}

// Node operations
func (m *ResourceManager) RegisterNode(ctx context.Context, node Node) error {
	// IS-04: When a node registers or updates, last_seen should be updated.
	// The repository implementation should handle this (e.g., CURRENT_TIMESTAMP in SQLite).
	return m.repo.UpsertNode(ctx, node)
}

func (m *ResourceManager) UnregisterNode(ctx context.Context, id string) error {
	devices, _ := m.repo.ListDevices(ctx)
	for _, d := range devices {
		if d.NodeID == id {
			m.UnregisterDevice(ctx, d.ID)
		}
	}
	return m.repo.DeleteNode(ctx, id)
}

// Device operations
func (m *ResourceManager) RegisterDevice(ctx context.Context, device Device) error {
	// IS-04: Ensure parent Node exists
	_, err := m.repo.GetNode(ctx, device.NodeID)
	if err != nil {
		return fmt.Errorf("parent node %s not found: %w", device.NodeID, err)
	}
	return m.repo.UpsertDevice(ctx, device)
}

func (m *ResourceManager) UnregisterDevice(ctx context.Context, id string) error {
	sources, _ := m.repo.ListSources(ctx)
	for _, s := range sources {
		if s.DeviceID == id {
			m.UnregisterSource(ctx, s.ID)
		}
	}
	senders, _ := m.repo.ListSenders(ctx)
	for _, s := range senders {
		if s.DeviceID == id {
			m.UnregisterSender(ctx, s.ID)
		}
	}
	receivers, _ := m.repo.ListReceivers(ctx)
	for _, r := range receivers {
		if r.DeviceID == id {
			m.UnregisterReceiver(ctx, r.ID)
		}
	}
	flows, _ := m.repo.ListFlows(ctx)
	for _, f := range flows {
		if f.DeviceID == id {
			m.UnregisterFlow(ctx, f.ID)
		}
	}
	return m.repo.DeleteDevice(ctx, id)
}

// Source operations
func (m *ResourceManager) RegisterSource(ctx context.Context, source Source) error {
	// IS-04: Ensure parent Device exists
	_, err := m.repo.GetDevice(ctx, source.DeviceID)
	if err != nil {
		return fmt.Errorf("parent device %s not found: %w", source.DeviceID, err)
	}
	return m.repo.UpsertSource(ctx, source)
}

func (m *ResourceManager) UnregisterSource(ctx context.Context, id string) error {
	flows, _ := m.repo.ListFlows(ctx)
	for _, f := range flows {
		if f.SourceID == id {
			m.UnregisterFlow(ctx, f.ID)
		}
	}
	return m.repo.DeleteSource(ctx, id)
}

// Flow operations
func (m *ResourceManager) RegisterFlow(ctx context.Context, flow Flow) error {
	// IS-04: Ensure parent Source exists
	_, err := m.repo.GetSource(ctx, flow.SourceID)
	if err != nil {
		return fmt.Errorf("parent source %s not found: %w", flow.SourceID, err)
	}
	// IS-04: Ensure parent Device exists
	_, err = m.repo.GetDevice(ctx, flow.DeviceID)
	if err != nil {
		return fmt.Errorf("parent device %s not found: %w", flow.DeviceID, err)
	}
	return m.repo.UpsertFlow(ctx, flow)
}

func (m *ResourceManager) UnregisterFlow(ctx context.Context, id string) error {
	return m.repo.DeleteFlow(ctx, id)
}

// Sender operations
func (m *ResourceManager) RegisterSender(ctx context.Context, sender Sender) error {
	// IS-04: Ensure parent Device exists
	_, err := m.repo.GetDevice(ctx, sender.DeviceID)
	if err != nil {
		return fmt.Errorf("parent device %s not found: %w", sender.DeviceID, err)
	}
	// Note: FlowID is optional for Senders in some versions/cases, but usually points to a Flow.
	if sender.FlowID != nil {
		_, err = m.repo.GetFlow(ctx, *sender.FlowID)
		if err != nil {
			return fmt.Errorf("parent flow %s not found: %w", *sender.FlowID, err)
		}
	}
	return m.repo.UpsertSender(ctx, sender)
}

func (m *ResourceManager) UnregisterSender(ctx context.Context, id string) error {
	return m.repo.DeleteSender(ctx, id)
}

// Receiver operations
func (m *ResourceManager) RegisterReceiver(ctx context.Context, receiver Receiver) error {
	// IS-04: Ensure parent Device exists
	_, err := m.repo.GetDevice(ctx, receiver.DeviceID)
	if err != nil {
		return fmt.Errorf("parent device %s not found: %w", receiver.DeviceID, err)
	}
	return m.repo.UpsertReceiver(ctx, receiver)
}

func (m *ResourceManager) UnregisterReceiver(ctx context.Context, id string) error {
	return m.repo.DeleteReceiver(ctx, id)
}

func (m *ResourceManager) NodeExists(ctx context.Context, id string) (bool, error) {
	_, err := m.repo.GetNode(ctx, id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (m *ResourceManager) DeviceExists(ctx context.Context, id string) (bool, error) {
	_, err := m.repo.GetDevice(ctx, id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (m *ResourceManager) SourceExists(ctx context.Context, id string) (bool, error) {
	_, err := m.repo.GetSource(ctx, id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (m *ResourceManager) FlowExists(ctx context.Context, id string) (bool, error) {
	_, err := m.repo.GetFlow(ctx, id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (m *ResourceManager) SenderExists(ctx context.Context, id string) (bool, error) {
	_, err := m.repo.GetSender(ctx, id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (m *ResourceManager) ReceiverExists(ctx context.Context, id string) (bool, error) {
	_, err := m.repo.GetReceiver(ctx, id)
	if err != nil {
		return false, nil
	}
	return true, nil
}
