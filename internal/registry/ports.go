package registry

import (
	"context"
	"time"
)

type IRepository interface {
	// Nodes
	UpsertNode(ctx context.Context, node Node) error
	DeleteNode(ctx context.Context, id string) error
	GetNode(ctx context.Context, id string) (Node, error)
	ListNodes(ctx context.Context) ([]Node, error)
	UpdateNodeHealth(ctx context.Context, id string) error
	GetExpiredNodes(ctx context.Context, expirationTime time.Time) ([]string, error)

	// Devices
	UpsertDevice(ctx context.Context, device Device) error
	DeleteDevice(ctx context.Context, id string) error
	GetDevice(ctx context.Context, id string) (Device, error)
	ListDevices(ctx context.Context) ([]Device, error)

	// Sources
	UpsertSource(ctx context.Context, source Source) error
	DeleteSource(ctx context.Context, id string) error
	GetSource(ctx context.Context, id string) (Source, error)
	ListSources(ctx context.Context) ([]Source, error)

	// Flows
	UpsertFlow(ctx context.Context, flow Flow) error
	DeleteFlow(ctx context.Context, id string) error
	GetFlow(ctx context.Context, id string) (Flow, error)
	ListFlows(ctx context.Context) ([]Flow, error)

	// Senders
	UpsertSender(ctx context.Context, sender Sender) error
	DeleteSender(ctx context.Context, id string) error
	GetSender(ctx context.Context, id string) (Sender, error)
	ListSenders(ctx context.Context) ([]Sender, error)

	// Receivers
	UpsertReceiver(ctx context.Context, receiver Receiver) error
	DeleteReceiver(ctx context.Context, id string) error
	GetReceiver(ctx context.Context, id string) (Receiver, error)
	ListReceivers(ctx context.Context) ([]Receiver, error)
}

type INotifier interface {
	Notify(ctx context.Context, resourceType ResourceType, action string, data interface{}) error
}
