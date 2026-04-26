package registry

import (
	"context"
	"errors"
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
	exists, err := m.repo.IDExistsAsOtherType(ctx, node.ID, ResourceTypeNode)
	if err != nil {
		return err
	}
	if exists {
		return NewValidationErrorf("id %s already exists as a different resource type", node.ID)
	}
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
	exists, err := m.repo.IDExistsAsOtherType(ctx, device.ID, ResourceTypeDevice)
	if err != nil {
		return err
	}
	if exists {
		return NewValidationErrorf("id %s already exists as a different resource type", device.ID)
	}

	existingDevice, err := m.repo.GetDevice(ctx, device.ID)
	if err == nil {
		if existingDevice.NodeID != device.NodeID {
			return NewValidationErrorf("parent node_id cannot be modified on update")
		}
	}

	_, err = m.repo.GetNode(ctx, device.NodeID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return NewValidationErrorf("parent node %s not found", device.NodeID)
		}
		return err
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
	exists, err := m.repo.IDExistsAsOtherType(ctx, source.ID, ResourceTypeSource)
	if err != nil {
		return err
	}
	if exists {
		return NewValidationErrorf("id %s already exists as a different resource type", source.ID)
	}

	existingSource, err := m.repo.GetSource(ctx, source.ID)
	if err == nil {
		if existingSource.DeviceID != source.DeviceID {
			return NewValidationErrorf("parent device_id cannot be modified on update")
		}
	}

	_, err = m.repo.GetDevice(ctx, source.DeviceID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return NewValidationErrorf("parent device %s not found", source.DeviceID)
		}
		return err
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
	exists, err := m.repo.IDExistsAsOtherType(ctx, flow.ID, ResourceTypeFlow)
	if err != nil {
		return err
	}
	if exists {
		return NewValidationErrorf("id %s already exists as a different resource type", flow.ID)
	}

	existingFlow, err := m.repo.GetFlow(ctx, flow.ID)
	if err == nil {
		if existingFlow.SourceID != flow.SourceID {
			return NewValidationErrorf("parent source_id cannot be modified on update")
		}
		if existingFlow.DeviceID != flow.DeviceID {
			return NewValidationErrorf("parent device_id cannot be modified on update")
		}
	}

	_, err = m.repo.GetSource(ctx, flow.SourceID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return NewValidationErrorf("parent source %s not found", flow.SourceID)
		}
		return err
	}
	_, err = m.repo.GetDevice(ctx, flow.DeviceID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return NewValidationErrorf("parent device %s not found", flow.DeviceID)
		}
		return err
	}
	return m.repo.UpsertFlow(ctx, flow)
}

func (m *ResourceManager) UnregisterFlow(ctx context.Context, id string) error {
	return m.repo.DeleteFlow(ctx, id)
}

// Sender operations
func (m *ResourceManager) RegisterSender(ctx context.Context, sender Sender) error {
	exists, err := m.repo.IDExistsAsOtherType(ctx, sender.ID, ResourceTypeSender)
	if err != nil {
		return err
	}
	if exists {
		return NewValidationErrorf("id %s already exists as a different resource type", sender.ID)
	}

	existingSender, err := m.repo.GetSender(ctx, sender.ID)
	if err == nil {
		if existingSender.DeviceID != sender.DeviceID {
			return NewValidationErrorf("parent device_id cannot be modified on update")
		}
	}

	_, err = m.repo.GetDevice(ctx, sender.DeviceID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return NewValidationErrorf("parent device %s not found", sender.DeviceID)
		}
		return err
	}
	if sender.FlowID != nil {
		_, err = m.repo.GetFlow(ctx, *sender.FlowID)
		if err != nil {
			if errors.Is(err, ErrResourceNotFound) {
				return NewValidationErrorf("parent flow %s not found", *sender.FlowID)
			}
			return err
		}
	}
	return m.repo.UpsertSender(ctx, sender)
}

func (m *ResourceManager) UnregisterSender(ctx context.Context, id string) error {
	return m.repo.DeleteSender(ctx, id)
}

// Receiver operations
func (m *ResourceManager) RegisterReceiver(ctx context.Context, receiver Receiver) error {
	exists, err := m.repo.IDExistsAsOtherType(ctx, receiver.ID, ResourceTypeReceiver)
	if err != nil {
		return err
	}
	if exists {
		return NewValidationErrorf("id %s already exists as a different resource type", receiver.ID)
	}

	existingReceiver, err := m.repo.GetReceiver(ctx, receiver.ID)
	if err == nil {
		if existingReceiver.DeviceID != receiver.DeviceID {
			return NewValidationErrorf("parent device_id cannot be modified on update")
		}
	}

	_, err = m.repo.GetDevice(ctx, receiver.DeviceID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return NewValidationErrorf("parent device %s not found", receiver.DeviceID)
		}
		return err
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
