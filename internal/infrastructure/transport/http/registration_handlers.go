package transporthttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	raw_data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	byte_reader := bytes.NewReader(raw_data)
	if err := json.NewDecoder(byte_reader).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()

	switch req.Type {
	case registry.ResourceTypeNode:
		var node registry.Node
		if err := json.Unmarshal(req.Data, &node); err != nil {
			h.writeError(w, http.StatusBadRequest, "invalid node data")
			return
		}
		h.registerNode(ctx, w, req.Type, node)
	case registry.ResourceTypeDevice:
		var device registry.Device
		if err := json.Unmarshal(req.Data, &device); err != nil {
			h.writeError(w, http.StatusBadRequest, "invalid device data")
			return
		}
		h.registerDevice(ctx, w, req.Type, device)
	case registry.ResourceTypeSource:
		var source registry.Source
		if err := json.Unmarshal(req.Data, &source); err != nil {
			h.writeError(w, http.StatusBadRequest, "invalid source data")
			return
		}
		h.registerSource(ctx, w, req.Type, source)
	case registry.ResourceTypeFlow:
		var flow registry.Flow
		if err := json.Unmarshal(req.Data, &flow); err != nil {
			h.writeError(w, http.StatusBadRequest, "invalid flow data")
			return
		}
		h.registerFlow(ctx, w, req.Type, flow)
	case registry.ResourceTypeSender:
		var sender registry.Sender
		if err := json.Unmarshal(req.Data, &sender); err != nil {
			h.writeError(w, http.StatusBadRequest, "invalid sender data")
			return
		}
		h.registerSender(ctx, w, req.Type, sender)
	case registry.ResourceTypeReceiver:
		var receiver registry.Receiver
		if err := json.Unmarshal(req.Data, &receiver); err != nil {
			h.writeError(w, http.StatusBadRequest, "invalid receiver data")
			return
		}
		h.registerReceiver(ctx, w, req.Type, receiver)
	default:
		h.writeError(w, http.StatusBadRequest, "invalid resource type")
	}
}

func (h *RegistrationHandlers) registerNode(ctx context.Context, w http.ResponseWriter, resourceType registry.ResourceType, node registry.Node) {
	if err := validateNode(node); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result := h.resourceManager.RegisterNode(ctx, node)
	h.respondWithResource(w, resourceType, node.ID, result.ToVoid(), node)
}

func (h *RegistrationHandlers) registerDevice(ctx context.Context, w http.ResponseWriter, resourceType registry.ResourceType, device registry.Device) {
	result := h.resourceManager.RegisterDevice(ctx, device)
	h.respondWithResource(w, resourceType, device.ID, result.ToVoid(), device)
}

func (h *RegistrationHandlers) registerSource(ctx context.Context, w http.ResponseWriter, resourceType registry.ResourceType, source registry.Source) {
	result := h.resourceManager.RegisterSource(ctx, source)
	h.respondWithResource(w, resourceType, source.ID, result.ToVoid(), source)
}

func (h *RegistrationHandlers) registerFlow(ctx context.Context, w http.ResponseWriter, resourceType registry.ResourceType, flow registry.Flow) {
	result := h.resourceManager.RegisterFlow(ctx, flow)
	h.respondWithResource(w, resourceType, flow.ID, result.ToVoid(), flow)
}

func (h *RegistrationHandlers) registerSender(ctx context.Context, w http.ResponseWriter, resourceType registry.ResourceType, sender registry.Sender) {
	result := h.resourceManager.RegisterSender(ctx, sender)
	h.respondWithResource(w, resourceType, sender.ID, result.ToVoid(), sender)
}

func (h *RegistrationHandlers) registerReceiver(ctx context.Context, w http.ResponseWriter, resourceType registry.ResourceType, receiver registry.Receiver) {
	result := h.resourceManager.RegisterReceiver(ctx, receiver)
	h.respondWithResource(w, resourceType, receiver.ID, result.ToVoid(), receiver)
}

func validateNode(node registry.Node) error {
	if node.ID == "" {
		return errors.New("id is required")
	}
	if !isValidUUID(node.ID) {
		return errors.New("id must be a valid UUID")
	}
	if node.Version == "" {
		return errors.New("version is required")
	}
	if !isValidTAI(node.Version) {
		return errors.New("version must be a valid TAI timestamp (<seconds>:<nanoseconds>)")
	}
	if node.Label == "" {
		return errors.New("label is required")
	}
	if node.Description == "" {
		return errors.New("description is required")
	}
	return nil
}

func isValidUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	parts := strings.Split(s, "-")
	if len(parts) != 5 {
		return false
	}
	if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
		return false
	}
	for _, c := range strings.Join(parts, "") {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

func isValidTAI(s string) bool {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return false
	}
	for _, c := range parts[0] {
		if c < '0' || c > '9' {
			return false
		}
	}
	for _, c := range parts[1] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (h *RegistrationHandlers) respond(w http.ResponseWriter, resourceType registry.ResourceType, resourceID string, result registry.Result[struct{}]) {
	if result.IsFailure() {
		err, _ := result.Error()
		h.writeError(w, err.Code, err.Message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if result.WasCreated() {
		w.Header().Set("Location", h.resourcePath(resourceType, resourceID))
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *RegistrationHandlers) respondWithResource(w http.ResponseWriter, resourceType registry.ResourceType, resourceID string, result registry.Result[struct{}], resource interface{}) {
	if result.IsFailure() {
		err, _ := result.Error()
		h.writeError(w, err.Code, err.Message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", h.resourcePath(resourceType, resourceID))
	if result.WasCreated() {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resource)
}

func (h *RegistrationHandlers) resourcePath(resourceType registry.ResourceType, id string) string {
	return fmt.Sprintf("/x-nmos/registration/v1.3/resource/%s/%s", strings.TrimSuffix(string(resourceType), "s")+"s", id)
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
		h.writeError(w, http.StatusBadRequest, "invalid resource type")
		return
	}

	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RegistrationHandlers) Heartbeat(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "nodeID")
	if nodeID == "" {
		h.writeError(w, http.StatusBadRequest, "missing node ID")
		return
	}

	if err := h.heartbeatEngine.Heartbeat(r.Context(), nodeID); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (h *RegistrationHandlers) writeError(w http.ResponseWriter, status int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:  status,
		Error: errorMsg,
	})
}
