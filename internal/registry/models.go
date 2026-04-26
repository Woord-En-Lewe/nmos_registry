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
	Tags        json.RawMessage `json:"tags"`
	Href        string          `json:"href"`
	Hostname    *string         `json:"hostname"`
	Caps        json.RawMessage `json:"caps"`
	Api         json.RawMessage `json:"api"`
	Services    json.RawMessage `json:"services"`
	Clocks      json.RawMessage `json:"clocks"`
	Interfaces  json.RawMessage `json:"interfaces"`
	LastSeen    time.Time       `json:"last_seen"`
}

type Device struct {
	ID          string          `json:"id"`
	ApiVersion  string          `json:"api_version"`
	NodeID      string          `json:"node_id"`
	Version     string          `json:"version"`
	Label       string          `json:"label"`
	Description string          `json:"description"`
	Tags        json.RawMessage `json:"tags"`
	Type        string          `json:"type"`
	Senders     json.RawMessage `json:"senders"`
	Receivers   json.RawMessage `json:"receivers"`
	Controls    json.RawMessage `json:"controls"`
}

type Source struct {
	ID          string          `json:"id"`
	ApiVersion  string          `json:"api_version"`
	DeviceID    string          `json:"device_id"`
	Version     string          `json:"version"`
	Label       string          `json:"label"`
	Description string          `json:"description"`
	Tags        json.RawMessage `json:"tags"`
	GrainRate   json.RawMessage `json:"grain_rate"`
	Format      string          `json:"format"`
	Caps        json.RawMessage `json:"caps"`
	Parents     json.RawMessage `json:"parents"`
	ClockName   *string         `json:"clock_name"`
}

type Flow struct {
	ID                     string          `json:"id"`
	ApiVersion             string          `json:"api_version"`
	SourceID               string          `json:"source_id"`
	DeviceID               string          `json:"device_id"`
	Version                string          `json:"version"`
	Label                  string          `json:"label"`
	Description            string          `json:"description"`
	Tags                   json.RawMessage `json:"tags"`
	Format                 string          `json:"format"`
	MediaType              *string         `json:"media_type"`
	SampleRate             json.RawMessage `json:"sample_rate"`
	BitDepth               *int            `json:"bit_depth"`
	DidSdid                json.RawMessage `json:"did_sdid"`
	GrainRate              json.RawMessage `json:"grain_rate"`
	FrameWidth             *int            `json:"frame_width"`
	FrameHeight            *int            `json:"frame_height"`
	InterlaceMode          *string         `json:"interlace_mode"`
	Colorspace             *string         `json:"colorspace"`
	Components             json.RawMessage `json:"components"`
	TransferCharacteristic *string         `json:"transfer_characteristic"`
}

type Sender struct {
	ID                     string          `json:"id"`
	ApiVersion             string          `json:"api_version"`
	DeviceID               string          `json:"device_id"`
	FlowID                 *string         `json:"flow_id"`
	Version                string          `json:"version"`
	Label                  string          `json:"label"`
	Description            string          `json:"description"`
	Tags                   json.RawMessage `json:"tags"`
	Transport              string          `json:"transport"`
	ManifestHref           *string         `json:"manifest_href"`
	InterfaceBindings      json.RawMessage `json:"interface_bindings"`
	SubscriptionReceiverID *string         `json:"subscription_receiver_id"`
	SubscriptionActive     *bool           `json:"subscription_active"`
}

type Receiver struct {
	ID                   string          `json:"id"`
	ApiVersion           string          `json:"api_version"`
	DeviceID             string          `json:"device_id"`
	Version              string          `json:"version"`
	Label                string          `json:"label"`
	Description          string          `json:"description"`
	Tags                 json.RawMessage `json:"tags"`
	Transport            string          `json:"transport"`
	Format               string          `json:"format"`
	Caps                 json.RawMessage `json:"caps"`
	InterfaceBindings    json.RawMessage `json:"interface_bindings"`
	SubscriptionSenderID *string         `json:"subscription_sender_id"`
	SubscriptionActive   *bool           `json:"subscription_active"`
}

type Subscription struct {
	ID              string          `json:"id"`
	ResourcePath    string          `json:"resource_path"`
	Params          json.RawMessage `json:"params"`
	MaxUpdateRateMs *int            `json:"max_update_rate_ms"`
	Persist         *bool           `json:"persist"`
	SecureWebsocket *bool           `json:"secure_websocket"`
	WsHref          *string         `json:"ws_href"`
}
