package registry

import (
	"encoding/json"
	"time"
)

type ResourceType string

const (
	ResourceTypeNode     ResourceType = "node"
	ResourceTypeDevice   ResourceType = "device"
	ResourceTypeSource   ResourceType = "source"
	ResourceTypeFlow     ResourceType = "flow"
	ResourceTypeSender   ResourceType = "sender"
	ResourceTypeReceiver ResourceType = "receiver"
)

type Node struct {
	ID          string          `json:"id"`
	ApiVersion  string          `json:"api_version"`
	Version     string          `json:"version"`
	Label       string          `json:"label"`
	Description string          `json:"description"`
	Tags        json.RawMessage `json:"tags,omitempty"`
	Href        string          `json:"href"`
	Hostname    *string         `json:"hostname,omitempty"`
	Caps        json.RawMessage `json:"caps,omitempty"`
	Api         json.RawMessage `json:"api"`
	Services    json.RawMessage `json:"services,omitempty"`
	Clocks      json.RawMessage `json:"clocks,omitempty"`
	Interfaces  json.RawMessage `json:"interfaces,omitempty"`
	LastSeen    time.Time       `json:"last_seen,omitempty"`
}

type Device struct {
	ID          string          `json:"id"`
	ApiVersion  string          `json:"api_version"`
	NodeID      string          `json:"node_id"`
	Version     string          `json:"version"`
	Label       string          `json:"label"`
	Description string          `json:"description"`
	Tags        json.RawMessage `json:"tags,omitempty"`
	Type        string          `json:"type"`
	Senders     json.RawMessage `json:"senders,omitempty"`
	Receivers   json.RawMessage `json:"receivers,omitempty"`
	Controls    json.RawMessage `json:"controls,omitempty"`
}

type Source struct {
	ID          string          `json:"id"`
	ApiVersion  string          `json:"api_version"`
	DeviceID    string          `json:"device_id"`
	Version     string          `json:"version"`
	Label       string          `json:"label"`
	Description string          `json:"description"`
	Tags        json.RawMessage `json:"tags,omitempty"`
	GrainRate   json.RawMessage `json:"grain_rate,omitempty"`
	Format      string          `json:"format"`
	Caps        json.RawMessage `json:"caps,omitempty"`
	Parents     json.RawMessage `json:"parents,omitempty"`
	ClockName   *string         `json:"clock_name,omitempty"`
}

type Flow struct {
	ID                     string          `json:"id"`
	ApiVersion             string          `json:"api_version"`
	SourceID               string          `json:"source_id"`
	DeviceID               string          `json:"device_id"`
	Version                string          `json:"version"`
	Label                  string          `json:"label"`
	Description            string          `json:"description"`
	Tags                   json.RawMessage `json:"tags,omitempty"`
	Format                 string          `json:"format"`
	MediaType              *string         `json:"media_type,omitempty"`
	SampleRate             json.RawMessage `json:"sample_rate,omitempty"`
	BitDepth               *int            `json:"bit_depth,omitempty"`
	DidSdid                json.RawMessage `json:"did_sdid,omitempty"`
	GrainRate              json.RawMessage `json:"grain_rate,omitempty"`
	FrameWidth             *int            `json:"frame_width,omitempty"`
	FrameHeight            *int            `json:"frame_height,omitempty"`
	InterlaceMode          *string         `json:"interlace_mode,omitempty"`
	Colorspace             *string         `json:"colorspace,omitempty"`
	Components             json.RawMessage `json:"components,omitempty"`
	TransferCharacteristic *string         `json:"transfer_characteristic,omitempty"`
}

type Sender struct {
	ID                     string                `json:"id"`
	ApiVersion             string                `json:"api_version"`
	DeviceID               string                `json:"device_id"`
	FlowID                 *string               `json:"flow_id,omitempty"`
	Version                string                `json:"version"`
	Label                  string                `json:"label"`
	Description            string                `json:"description"`
	Tags                   json.RawMessage       `json:"tags,omitempty"`
	Transport              string                `json:"transport"`
	ManifestHref           *string               `json:"manifest_href,omitempty"`
	InterfaceBindings      json.RawMessage       `json:"interface_bindings,omitempty"`
	Caps                   json.RawMessage       `json:"caps,omitempty"`
	Subscription           *SenderSubscription   `json:"subscription,omitempty"`
}

type SenderSubscription struct {
	ReceiverID *string `json:"receiver_id,omitempty"`
	Active     *bool   `json:"active,omitempty"`
}

type Receiver struct {
	ID                string          `json:"id"`
	ApiVersion        string          `json:"api_version"`
	DeviceID          string          `json:"device_id"`
	Version           string          `json:"version"`
	Label             string          `json:"label"`
	Description       string          `json:"description"`
	Tags              json.RawMessage `json:"tags,omitempty"`
	Transport         string          `json:"transport"`
	Format            string          `json:"format"`
	Caps              json.RawMessage `json:"caps,omitempty"`
	InterfaceBindings json.RawMessage `json:"interface_bindings,omitempty"`
	Subscription      *ReceiverSubscription `json:"subscription,omitempty"`
}

type ReceiverSubscription struct {
	SenderID *string `json:"sender_id,omitempty"`
	Active   *bool   `json:"active,omitempty"`
}

type Subscription struct {
	ID              string          `json:"id"`
	ResourcePath    string          `json:"resource_path"`
	Params          json.RawMessage `json:"params,omitempty"`
	MaxUpdateRateMs int             `json:"max_update_rate_ms"`
	Persist         bool            `json:"persist"`
	Secure          *bool           `json:"secure,omitempty"`
	WsHref          string          `json:"ws_href"`
}
