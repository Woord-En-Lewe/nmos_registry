package registry

import (
	"context"
	"errors"
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
func (m *ResourceManager) RegisterNode(ctx context.Context, node Node) Result[Node] {
	return Bind(m.checkIDCollision(ctx, node.ID, ResourceTypeNode), func(_ struct{}) Result[Node] {
		return m.persistNode(ctx, node)
	})
}

func (m *ResourceManager) persistNode(ctx context.Context, node Node) Result[Node] {
	_, err := m.repo.GetNode(ctx, node.ID)
	if err != nil && !errors.Is(err, ErrResourceNotFound) {
		return Failuref[Node](500, "database error: %v", err)
	}
	isUpdate := err == nil
	if err := m.repo.UpsertNode(ctx, node); err != nil {
		return Failuref[Node](500, "database error: %v", err)
	}
	if isUpdate {
		return Success(node)
	}
	return Created(node)
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
func (m *ResourceManager) RegisterDevice(ctx context.Context, device Device) Result[Device] {
	return Bind(m.checkIDCollision(ctx, device.ID, ResourceTypeDevice), func(_ struct{}) Result[Device] {
		return Bind(m.checkDeviceParentNotModified(ctx, device), func(_ struct{}) Result[Device] {
			return Bind(m.checkDeviceParentExists(ctx, device), func(_ struct{}) Result[Device] {
				return m.persistDevice(ctx, device)
			})
		})
	})
}

func (m *ResourceManager) checkDeviceParentNotModified(ctx context.Context, device Device) Result[struct{}] {
	existingDevice, err := m.repo.GetDevice(ctx, device.ID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return Success(struct{}{})
		}
		return Failuref[struct{}](500, "database error: %v", err)
	}
	if existingDevice.NodeID != device.NodeID {
		return Failure[struct{}](400, "parent node_id cannot be modified on update", nil)
	}
	return Success(struct{}{})
}

func (m *ResourceManager) checkDeviceParentExists(ctx context.Context, device Device) Result[struct{}] {
	_, err := m.repo.GetNode(ctx, device.NodeID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return Failuref[struct{}](400, "parent node %s not found", device.NodeID)
		}
		return Failuref[struct{}](500, "database error: %v", err)
	}
	return Success(struct{}{})
}

func (m *ResourceManager) persistDevice(ctx context.Context, device Device) Result[Device] {
	_, err := m.repo.GetDevice(ctx, device.ID)
	if err != nil && !errors.Is(err, ErrResourceNotFound) {
		return Failuref[Device](500, "database error: %v", err)
	}
	isUpdate := err == nil
	if err := m.repo.UpsertDevice(ctx, device); err != nil {
		return Failuref[Device](500, "database error: %v", err)
	}
	if isUpdate {
		return Success(device)
	}
	return Created(device)
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
func (m *ResourceManager) RegisterSource(ctx context.Context, source Source) Result[Source] {
	return Bind(m.checkIDCollision(ctx, source.ID, ResourceTypeSource), func(_ struct{}) Result[Source] {
		return Bind(m.checkSourceParentNotModified(ctx, source), func(_ struct{}) Result[Source] {
			return Bind(m.checkSourceParentExists(ctx, source), func(_ struct{}) Result[Source] {
				return m.persistSource(ctx, source)
			})
		})
	})
}

func (m *ResourceManager) checkSourceParentNotModified(ctx context.Context, source Source) Result[struct{}] {
	existingSource, err := m.repo.GetSource(ctx, source.ID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return Success(struct{}{})
		}
		return Failuref[struct{}](500, "database error: %v", err)
	}
	if existingSource.DeviceID != source.DeviceID {
		return Failure[struct{}](400, "parent device_id cannot be modified on update", nil)
	}
	return Success(struct{}{})
}

func (m *ResourceManager) checkSourceParentExists(ctx context.Context, source Source) Result[struct{}] {
	_, err := m.repo.GetDevice(ctx, source.DeviceID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return Failuref[struct{}](400, "parent device %s not found", source.DeviceID)
		}
		return Failuref[struct{}](500, "database error: %v", err)
	}
	return Success(struct{}{})
}

func (m *ResourceManager) persistSource(ctx context.Context, source Source) Result[Source] {
	_, err := m.repo.GetSource(ctx, source.ID)
	if err != nil && !errors.Is(err, ErrResourceNotFound) {
		return Failuref[Source](500, "database error: %v", err)
	}
	isUpdate := err == nil
	if err := m.repo.UpsertSource(ctx, source); err != nil {
		return Failuref[Source](500, "database error: %v", err)
	}
	if isUpdate {
		return Success(source)
	}
	return Created(source)
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
func (m *ResourceManager) RegisterFlow(ctx context.Context, flow Flow) Result[Flow] {
	return Bind(m.checkIDCollision(ctx, flow.ID, ResourceTypeFlow), func(_ struct{}) Result[Flow] {
		return Bind(m.checkFlowParentNotModified(ctx, flow), func(_ struct{}) Result[Flow] {
			return Bind(m.checkFlowParentsExist(ctx, flow), func(_ struct{}) Result[Flow] {
				return m.persistFlow(ctx, flow)
			})
		})
	})
}

func (m *ResourceManager) checkFlowParentNotModified(ctx context.Context, flow Flow) Result[struct{}] {
	existingFlow, err := m.repo.GetFlow(ctx, flow.ID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return Success(struct{}{})
		}
		return Failuref[struct{}](500, "database error: %v", err)
	}
	if existingFlow.SourceID != flow.SourceID {
		return Failure[struct{}](400, "parent source_id cannot be modified on update", nil)
	}
	if existingFlow.DeviceID != flow.DeviceID {
		return Failure[struct{}](400, "parent device_id cannot be modified on update", nil)
	}
	return Success(struct{}{})
}

func (m *ResourceManager) checkFlowParentsExist(ctx context.Context, flow Flow) Result[struct{}] {
	_, err := m.repo.GetSource(ctx, flow.SourceID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return Failuref[struct{}](400, "parent source %s not found", flow.SourceID)
		}
		return Failuref[struct{}](500, "database error: %v", err)
	}
	_, err = m.repo.GetDevice(ctx, flow.DeviceID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return Failuref[struct{}](400, "parent device %s not found", flow.DeviceID)
		}
		return Failuref[struct{}](500, "database error: %v", err)
	}
	return Success(struct{}{})
}

func (m *ResourceManager) persistFlow(ctx context.Context, flow Flow) Result[Flow] {
	_, err := m.repo.GetFlow(ctx, flow.ID)
	if err != nil && !errors.Is(err, ErrResourceNotFound) {
		return Failuref[Flow](500, "database error: %v", err)
	}
	isUpdate := err == nil
	if err := m.repo.UpsertFlow(ctx, flow); err != nil {
		return Failuref[Flow](500, "database error: %v", err)
	}
	if isUpdate {
		return Success(flow)
	}
	return Created(flow)
}

func (m *ResourceManager) UnregisterFlow(ctx context.Context, id string) error {
	return m.repo.DeleteFlow(ctx, id)
}

// Sender operations
func (m *ResourceManager) RegisterSender(ctx context.Context, sender Sender) Result[Sender] {
	return Bind(m.checkIDCollision(ctx, sender.ID, ResourceTypeSender), func(_ struct{}) Result[Sender] {
		return Bind(m.checkSenderParentNotModified(ctx, sender), func(_ struct{}) Result[Sender] {
			return Bind(m.checkSenderParentsExist(ctx, sender), func(_ struct{}) Result[Sender] {
				return m.persistSender(ctx, sender)
			})
		})
	})
}

func (m *ResourceManager) checkSenderParentNotModified(ctx context.Context, sender Sender) Result[struct{}] {
	existingSender, err := m.repo.GetSender(ctx, sender.ID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return Success(struct{}{})
		}
		return Failuref[struct{}](500, "database error: %v", err)
	}
	if existingSender.DeviceID != sender.DeviceID {
		return Failure[struct{}](400, "parent device_id cannot be modified on update", nil)
	}
	return Success(struct{}{})
}

func (m *ResourceManager) checkSenderParentsExist(ctx context.Context, sender Sender) Result[struct{}] {
	_, err := m.repo.GetDevice(ctx, sender.DeviceID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return Failuref[struct{}](400, "parent device %s not found", sender.DeviceID)
		}
		return Failuref[struct{}](500, "database error: %v", err)
	}
	if sender.FlowID != nil {
		_, err = m.repo.GetFlow(ctx, *sender.FlowID)
		if err != nil {
			if errors.Is(err, ErrResourceNotFound) {
				return Failuref[struct{}](400, "parent flow %s not found", *sender.FlowID)
			}
			return Failuref[struct{}](500, "database error: %v", err)
		}
	}
	return Success(struct{}{})
}

func (m *ResourceManager) persistSender(ctx context.Context, sender Sender) Result[Sender] {
	_, err := m.repo.GetSender(ctx, sender.ID)
	if err != nil && !errors.Is(err, ErrResourceNotFound) {
		return Failuref[Sender](500, "database error: %v", err)
	}
	isUpdate := err == nil
	if err := m.repo.UpsertSender(ctx, sender); err != nil {
		return Failuref[Sender](500, "database error: %v", err)
	}
	if isUpdate {
		return Success(sender)
	}
	return Created(sender)
}

func (m *ResourceManager) UnregisterSender(ctx context.Context, id string) error {
	return m.repo.DeleteSender(ctx, id)
}

// Receiver operations
func (m *ResourceManager) RegisterReceiver(ctx context.Context, receiver Receiver) Result[Receiver] {
	return Bind(m.checkIDCollision(ctx, receiver.ID, ResourceTypeReceiver), func(_ struct{}) Result[Receiver] {
		return Bind(m.checkReceiverParentNotModified(ctx, receiver), func(_ struct{}) Result[Receiver] {
			return Bind(m.checkReceiverParentExists(ctx, receiver), func(_ struct{}) Result[Receiver] {
				return m.persistReceiver(ctx, receiver)
			})
		})
	})
}

func (m *ResourceManager) checkReceiverParentNotModified(ctx context.Context, receiver Receiver) Result[struct{}] {
	existingReceiver, err := m.repo.GetReceiver(ctx, receiver.ID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return Success(struct{}{})
		}
		return Failuref[struct{}](500, "database error: %v", err)
	}
	if existingReceiver.DeviceID != receiver.DeviceID {
		return Failure[struct{}](400, "parent device_id cannot be modified on update", nil)
	}
	return Success(struct{}{})
}

func (m *ResourceManager) checkReceiverParentExists(ctx context.Context, receiver Receiver) Result[struct{}] {
	_, err := m.repo.GetDevice(ctx, receiver.DeviceID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return Failuref[struct{}](400, "parent device %s not found", receiver.DeviceID)
		}
		return Failuref[struct{}](500, "database error: %v", err)
	}
	return Success(struct{}{})
}

func (m *ResourceManager) persistReceiver(ctx context.Context, receiver Receiver) Result[Receiver] {
	_, err := m.repo.GetReceiver(ctx, receiver.ID)
	if err != nil && !errors.Is(err, ErrResourceNotFound) {
		return Failuref[Receiver](500, "database error: %v", err)
	}
	isUpdate := err == nil
	if err := m.repo.UpsertReceiver(ctx, receiver); err != nil {
		return Failuref[Receiver](500, "database error: %v", err)
	}
	if isUpdate {
		return Success(receiver)
	}
	return Created(receiver)
}

func (m *ResourceManager) UnregisterReceiver(ctx context.Context, id string) error {
	return m.repo.DeleteReceiver(ctx, id)
}

// Helper methods
func (m *ResourceManager) checkIDCollision(ctx context.Context, id string, resourceType ResourceType) Result[struct{}] {
	exists, err := m.repo.IDExistsAsOtherType(ctx, id, resourceType)
	if err != nil {
		return Failuref[struct{}](500, "database error: %v", err)
	}
	if exists {
		return Failuref[struct{}](400, "id %s already exists as a different resource type", id)
	}
	return Success(struct{}{})
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

// GetResource is a helper for query handlers
func (m *ResourceManager) GetNode(ctx context.Context, id string) (Node, error) {
	return m.repo.GetNode(ctx, id)
}

func (m *ResourceManager) GetDevice(ctx context.Context, id string) (Device, error) {
	return m.repo.GetDevice(ctx, id)
}

func (m *ResourceManager) GetSource(ctx context.Context, id string) (Source, error) {
	return m.repo.GetSource(ctx, id)
}

func (m *ResourceManager) GetFlow(ctx context.Context, id string) (Flow, error) {
	return m.repo.GetFlow(ctx, id)
}

func (m *ResourceManager) GetSender(ctx context.Context, id string) (Sender, error) {
	return m.repo.GetSender(ctx, id)
}

func (m *ResourceManager) GetReceiver(ctx context.Context, id string) (Receiver, error) {
	return m.repo.GetReceiver(ctx, id)
}

// Unused import check
var _ = fmt.Sprintf