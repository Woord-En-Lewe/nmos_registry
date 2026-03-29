package transporthttp

import (
	"encoding/json"
	"net/http"
	"github.com/Woord-En-Lewe/nmos_registry/internal/registry"

	"github.com/go-chi/chi/v5"
)

type RegistrationHandlers struct {
	resourceManager *registry.ResourceManager
	heartbeatEngine *registry.HeartbeatEngine
}

func NewRegistrationHandlers(rm *registry.ResourceManager, he *registry.HeartbeatEngine) *RegistrationHandlers {
	return &RegistrationHandlers{
		resourceManager: rm,
		heartbeatEngine: he,
	}
}

type registrationRequest struct {
	Type registry.ResourceType `json:"type"`
	Data json.RawMessage       `json:"data"`
}

func (h *RegistrationHandlers) RegisterResource(w http.ResponseWriter, r *http.Request) {
	var req registrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var err error

	switch req.Type {
	case registry.ResourceTypeNode:
		var node registry.Node
		if err = json.Unmarshal(req.Data, &node); err == nil {
			err = h.resourceManager.RegisterNode(ctx, node)
		}
	case registry.ResourceTypeDevice:
		var device registry.Device
		if err = json.Unmarshal(req.Data, &device); err == nil {
			err = h.resourceManager.RegisterDevice(ctx, device)
		}
	case registry.ResourceTypeSource:
		var source registry.Source
		if err = json.Unmarshal(req.Data, &source); err == nil {
			err = h.resourceManager.RegisterSource(ctx, source)
		}
	case registry.ResourceTypeFlow:
		var flow registry.Flow
		if err = json.Unmarshal(req.Data, &flow); err == nil {
			err = h.resourceManager.RegisterFlow(ctx, flow)
		}
	case registry.ResourceTypeSender:
		var sender registry.Sender
		if err = json.Unmarshal(req.Data, &sender); err == nil {
			err = h.resourceManager.RegisterSender(ctx, sender)
		}
	case registry.ResourceTypeReceiver:
		var receiver registry.Receiver
		if err = json.Unmarshal(req.Data, &receiver); err == nil {
			err = h.resourceManager.RegisterReceiver(ctx, receiver)
		}
	default:
		http.Error(w, "invalid resource type", http.StatusBadRequest)
		return
	}

	if err != nil {
		// In a real implementation, we might want to distinguish between 400 and 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(req.Data)
}

func (h *RegistrationHandlers) UnregisterResource(w http.ResponseWriter, r *http.Request) {
	resourceType := registry.ResourceType(chi.URLParam(r, "type"))
	id := chi.URLParam(r, "id")

	ctx := r.Context()
	var err error

	switch resourceType {
	case registry.ResourceTypeNode:
		err = h.resourceManager.UnregisterNode(ctx, id)
	case registry.ResourceTypeDevice:
		err = h.resourceManager.UnregisterDevice(ctx, id)
	case registry.ResourceTypeSource:
		err = h.resourceManager.UnregisterSource(ctx, id)
	case registry.ResourceTypeFlow:
		err = h.resourceManager.UnregisterFlow(ctx, id)
	case registry.ResourceTypeSender:
		err = h.resourceManager.UnregisterSender(ctx, id)
	case registry.ResourceTypeReceiver:
		err = h.resourceManager.UnregisterReceiver(ctx, id)
	default:
		http.Error(w, "invalid resource type", http.StatusBadRequest)
		return
	}

	if err != nil {
		// Should probably return 404 if not found, but IRepo doesn't currently distinguish
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RegistrationHandlers) Heartbeat(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "nodeID")
	if nodeID == "" {
		http.Error(w, "missing node ID", http.StatusBadRequest)
		return
	}

	if err := h.heartbeatEngine.Heartbeat(r.Context(), nodeID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// IS-04 says to return the health object, which currently just contains the health value (timestamp)
	// For now, we'll just return a simple JSON with the node ID
	json.NewEncoder(w).Encode(map[string]string{"id": nodeID})
}
