package transporthttp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/schema"
	"github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/transport/websocket"
	"github.com/Woord-En-Lewe/nmos_registry/internal/registry"
	"github.com/go-chi/chi/v5"
)

type QueryHandlers struct {
	repo               registry.IRepository
	subscriptionEngine *registry.SubscriptionEngine
	wsManager          *websocket.Manager
	wsBaseURL          string
	validator          *schema.Validator
}

func NewQueryHandlers(repo registry.IRepository, subscriptionEngine *registry.SubscriptionEngine, wsManager *websocket.Manager, wsBaseURL string) *QueryHandlers {
	v, err := schema.NewValidator()
	if err != nil {
		log.Printf("Warning: failed to create schema validator: %v", err)
	}
	return &QueryHandlers{
		repo:               repo,
		subscriptionEngine: subscriptionEngine,
		wsManager:          wsManager,
		wsBaseURL:          wsBaseURL,
		validator:          v,
	}
}

func (h *QueryHandlers) writeError(w http.ResponseWriter, status int, errorMsg string, debug string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:  status,
		Error: errorMsg,
		Debug: debug,
	})
}

func (h *QueryHandlers) checkPaginationSupported(w http.ResponseWriter, r *http.Request) bool {
	query := r.URL.Query()
	if len(query) > 0 {
		h.writeError(w, http.StatusNotImplemented, "query parameters not supported", "")
		return false
	}
	return true
}

func (h *QueryHandlers) ListNodes(w http.ResponseWriter, r *http.Request) {
	if !h.checkPaginationSupported(w, r) {
		return
	}
	ctx := r.Context()
	nodes, err := h.repo.ListNodes(ctx)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		for i, node := range nodes {
			data, _ := json.Marshal(node)
			if err := h.validator.ValidateJSON("node", data); err != nil {
				log.Printf("Node validation failed for index %d: %v", i, err)
				h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
				return
			}
		}
	}
	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) GetNode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	node, err := h.repo.GetNode(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "node not found", "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		data, _ := json.Marshal(node)
		if err := h.validator.ValidateJSON("node", data); err != nil {
			log.Printf("Node validation failed: %v", err)
			h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
			return
		}
	}
	if err := json.NewEncoder(w).Encode(node); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) ListDevices(w http.ResponseWriter, r *http.Request) {
	if !h.checkPaginationSupported(w, r) {
		return
	}
	ctx := r.Context()
	devices, err := h.repo.ListDevices(ctx)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		for i, device := range devices {
			data, _ := json.Marshal(device)
			if err := h.validator.ValidateJSON("device", data); err != nil {
				log.Printf("Device validation failed for index %d: %v", i, err)
				h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
				return
			}
		}
	}
	if err := json.NewEncoder(w).Encode(devices); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) GetDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	device, err := h.repo.GetDevice(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "device not found", "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		data, _ := json.Marshal(device)
		if err := h.validator.ValidateJSON("device", data); err != nil {
			log.Printf("Device validation failed: %v", err)
			h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
			return
		}
	}
	if err := json.NewEncoder(w).Encode(device); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) ListSources(w http.ResponseWriter, r *http.Request) {
	if !h.checkPaginationSupported(w, r) {
		return
	}
	ctx := r.Context()
	sources, err := h.repo.ListSources(ctx)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		for i, source := range sources {
			data, _ := json.Marshal(source)
			if err := h.validator.ValidateJSON("source", data); err != nil {
				log.Printf("Source validation failed for index %d: %v", i, err)
				h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
				return
			}
		}
	}
	if err := json.NewEncoder(w).Encode(sources); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) GetSource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	source, err := h.repo.GetSource(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "source not found", "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		data, _ := json.Marshal(source)
		if err := h.validator.ValidateJSON("source", data); err != nil {
			log.Printf("Source validation failed: %v", err)
			h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
			return
		}
	}
	if err := json.NewEncoder(w).Encode(source); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) ListFlows(w http.ResponseWriter, r *http.Request) {
	if !h.checkPaginationSupported(w, r) {
		return
	}
	ctx := r.Context()
	flows, err := h.repo.ListFlows(ctx)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		for i, flow := range flows {
			data, _ := json.Marshal(flow)
			if err := h.validator.ValidateJSON("flow", data); err != nil {
				log.Printf("Flow validation failed for index %d: %v", i, err)
				h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
				return
			}
		}
	}
	if err := json.NewEncoder(w).Encode(flows); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) GetFlow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	flow, err := h.repo.GetFlow(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "flow not found", "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		data, _ := json.Marshal(flow)
		if err := h.validator.ValidateJSON("flow", data); err != nil {
			log.Printf("Flow validation failed: %v", err)
			h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
			return
		}
	}
	if err := json.NewEncoder(w).Encode(flow); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) ListSenders(w http.ResponseWriter, r *http.Request) {
	if !h.checkPaginationSupported(w, r) {
		return
	}
	ctx := r.Context()
	senders, err := h.repo.ListSenders(ctx)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		for i, sender := range senders {
			data, _ := json.Marshal(sender)
			if err := h.validator.ValidateJSON("sender", data); err != nil {
				log.Printf("Sender validation failed for index %d: %v", i, err)
				h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
				return
			}
		}
	}
	if err := json.NewEncoder(w).Encode(senders); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) GetSender(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	sender, err := h.repo.GetSender(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "sender not found", "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		data, _ := json.Marshal(sender)
		if err := h.validator.ValidateJSON("sender", data); err != nil {
			log.Printf("Sender validation failed: %v", err)
			h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
			return
		}
	}
	if err := json.NewEncoder(w).Encode(sender); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) ListReceivers(w http.ResponseWriter, r *http.Request) {
	if !h.checkPaginationSupported(w, r) {
		return
	}
	ctx := r.Context()
	receivers, err := h.repo.ListReceivers(ctx)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		for i, receiver := range receivers {
			data, _ := json.Marshal(receiver)
			if err := h.validator.ValidateJSON("receiver", data); err != nil {
				log.Printf("Receiver validation failed for index %d: %v", i, err)
				h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
				return
			}
		}
	}
	if err := json.NewEncoder(w).Encode(receivers); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) GetReceiver(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	receiver, err := h.repo.GetReceiver(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "receiver not found", "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		data, _ := json.Marshal(receiver)
		if err := h.validator.ValidateJSON("receiver", data); err != nil {
			log.Printf("Receiver validation failed: %v", err)
			h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
			return
		}
	}
	if err := json.NewEncoder(w).Encode(receiver); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if req.ResourcePath == "" {
		h.writeError(w, http.StatusBadRequest, "resource_path is required", "")
		return
	}

	sub, err := h.subscriptionEngine.CreateSubscription(ctx, req.ResourcePath, req.Params, req.MaxUpdateRateMs, req.Persist, req.Secure)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}

	wsProtocol := "ws"
	if sub.Secure != nil && *sub.Secure {
		wsProtocol = "wss"
	}
	wsHref := fmt.Sprintf("%s://%s/x-nmos/query/v1.3/subscriptions/%s/ws", wsProtocol, h.wsBaseURL, sub.ID)
	sub.WsHref = wsHref

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/x-nmos/query/v1.3/subscriptions/%s", sub.ID))
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(sub); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	sub, err := h.subscriptionEngine.GetSubscription(ctx, id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "subscription not found", "")
		return
	}

	wsProtocol := "ws"
	if sub.Secure != nil && *sub.Secure {
		wsProtocol = "wss"
	}
	sub.WsHref = fmt.Sprintf("%s://%s/x-nmos/query/v1.3/subscriptions/%s/ws", wsProtocol, h.wsBaseURL, sub.ID)

	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		data, _ := json.Marshal(sub)
		if err := h.validator.ValidateJSON("subscription", data); err != nil {
			log.Printf("Subscription validation failed: %v", err)
			h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
			return
		}
	}
	if err := json.NewEncoder(w).Encode(sub); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	if !h.checkPaginationSupported(w, r) {
		return
	}
	ctx := r.Context()
	subs, err := h.subscriptionEngine.ListSubscriptions(ctx)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}
	for i := range subs {
		wsProtocol := "ws"
		if subs[i].Secure != nil && *subs[i].Secure {
			wsProtocol = "wss"
		}
		subs[i].WsHref = fmt.Sprintf("%s://%s/x-nmos/query/v1.3/subscriptions/%s/ws", wsProtocol, h.wsBaseURL, subs[i].ID)
	}
	w.Header().Set("Content-Type", "application/json")
	if h.validator != nil {
		for i, sub := range subs {
			data, _ := json.Marshal(sub)
			if err := h.validator.ValidateJSON("subscription", data); err != nil {
				log.Printf("Subscription validation failed for index %d: %v", i, err)
				h.writeError(w, http.StatusInternalServerError, "schema validation error", schema.ExtractValidationError(err))
				return
			}
		}
	}
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to encode response", err.Error())
		return
	}
}

func (h *QueryHandlers) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	if err := h.subscriptionEngine.DeleteSubscription(ctx, id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error(), "")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *QueryHandlers) HandleSubscriptionWebSocket(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.wsManager.HandleWebSocket(w, r, id)
}

type CreateSubscriptionRequest struct {
	ResourcePath    string          `json:"resource_path"`
	Params          json.RawMessage `json:"params,omitempty"`
	MaxUpdateRateMs *int            `json:"max_update_rate_ms,omitempty"`
	Persist         *bool           `json:"persist,omitempty"`
	Secure          *bool           `json:"secure,omitempty"`
}
