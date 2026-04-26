package transporthttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Woord-En-Lewe/nmos_registry/internal/registry"
	"github.com/go-chi/chi/v5"
)

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Debug string `json:"debug"`
}

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
		h.writeError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	ctx := r.Context()
	var err error
	var resourceID string

	switch req.Type {
	case registry.ResourceTypeNode:
		var node registry.Node
		if err = json.Unmarshal(req.Data, &node); err == nil {
			resourceID = node.ID
			exists, _ := h.resourceManager.NodeExists(ctx, node.ID)
			err = h.resourceManager.RegisterNode(ctx, node)
			if err == nil && !exists {
				w.Header().Set("Location", h.resourcePath(req.Type, resourceID))
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	case registry.ResourceTypeDevice:
		var device registry.Device
		if err = json.Unmarshal(req.Data, &device); err == nil {
			resourceID = device.ID
			exists, _ := h.resourceManager.DeviceExists(ctx, device.ID)
			err = h.resourceManager.RegisterDevice(ctx, device)
			if err == nil && !exists {
				w.Header().Set("Location", h.resourcePath(req.Type, resourceID))
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	case registry.ResourceTypeSource:
		var source registry.Source
		if err = json.Unmarshal(req.Data, &source); err == nil {
			resourceID = source.ID
			exists, _ := h.resourceManager.SourceExists(ctx, source.ID)
			err = h.resourceManager.RegisterSource(ctx, source)
			if err == nil && !exists {
				w.Header().Set("Location", h.resourcePath(req.Type, resourceID))
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	case registry.ResourceTypeFlow:
		var flow registry.Flow
		if err = json.Unmarshal(req.Data, &flow); err == nil {
			resourceID = flow.ID
			exists, _ := h.resourceManager.FlowExists(ctx, flow.ID)
			err = h.resourceManager.RegisterFlow(ctx, flow)
			if err == nil && !exists {
				w.Header().Set("Location", h.resourcePath(req.Type, resourceID))
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	case registry.ResourceTypeSender:
		var sender registry.Sender
		if err = json.Unmarshal(req.Data, &sender); err == nil {
			resourceID = sender.ID
			exists, _ := h.resourceManager.SenderExists(ctx, sender.ID)
			err = h.resourceManager.RegisterSender(ctx, sender)
			if err == nil && !exists {
				w.Header().Set("Location", h.resourcePath(req.Type, resourceID))
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	case registry.ResourceTypeReceiver:
		var receiver registry.Receiver
		if err = json.Unmarshal(req.Data, &receiver); err == nil {
			resourceID = receiver.ID
			exists, _ := h.resourceManager.ReceiverExists(ctx, receiver.ID)
			err = h.resourceManager.RegisterReceiver(ctx, receiver)
			if err == nil && !exists {
				w.Header().Set("Location", h.resourcePath(req.Type, resourceID))
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	default:
		h.writeError(w, http.StatusBadRequest, "invalid resource type", "")
		return
	}

	if err != nil {
		var validationErr *registry.ValidationError
		if errors.As(err, &validationErr) {
			h.writeError(w, http.StatusBadRequest, err.Error(), "")
		} else {
			h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(req.Data)
}

func (h *RegistrationHandlers) resourcePath(resourceType registry.ResourceType, id string) string {
	return fmt.Sprintf("/x-nmos/query/v1.3/%s/%s", strings.TrimSuffix(string(resourceType), "s")+"s", id)
}

func (h *RegistrationHandlers) UnregisterResource(w http.ResponseWriter, r *http.Request) {
	resourceType := registry.ResourceType(chi.URLParam(r, "type"))
	id := chi.URLParam(r, "id")

	ctx := r.Context()
	var err error

	switch resourceType {
	case registry.ResourceTypeNode, "nodes":
		err = h.resourceManager.UnregisterNode(ctx, id)
	case registry.ResourceTypeDevice, "devices":
		err = h.resourceManager.UnregisterDevice(ctx, id)
	case registry.ResourceTypeSource, "sources":
		err = h.resourceManager.UnregisterSource(ctx, id)
	case registry.ResourceTypeFlow, "flows":
		err = h.resourceManager.UnregisterFlow(ctx, id)
	case registry.ResourceTypeSender, "senders":
		err = h.resourceManager.UnregisterSender(ctx, id)
	case registry.ResourceTypeReceiver, "receivers":
		err = h.resourceManager.UnregisterReceiver(ctx, id)
	default:
		h.writeError(w, http.StatusBadRequest, "invalid resource type", "")
		return
	}

	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RegistrationHandlers) Heartbeat(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "nodeID")
	if nodeID == "" {
		h.writeError(w, http.StatusBadRequest, "missing node ID", "")
		return
	}

	if err := h.heartbeatEngine.Heartbeat(r.Context(), nodeID); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"health": nodeID})
}

func (h *RegistrationHandlers) writeError(w http.ResponseWriter, status int, errorMsg string, debug string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:  status,
		Error: errorMsg,
		Debug: debug,
	})
}
