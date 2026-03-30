package transporthttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/transport/websocket"
	"github.com/Woord-En-Lewe/nmos_registry/internal/registry"
	"github.com/go-chi/chi/v5"
)

type QueryHandlers struct {
	repo               registry.IRepository
	subscriptionEngine *registry.SubscriptionEngine
	wsManager          *websocket.Manager
	wsBaseURL          string
}

func NewQueryHandlers(repo registry.IRepository, subscriptionEngine *registry.SubscriptionEngine, wsManager *websocket.Manager, wsBaseURL string) *QueryHandlers {
	return &QueryHandlers{
		repo:               repo,
		subscriptionEngine: subscriptionEngine,
		wsManager:          wsManager,
		wsBaseURL:          wsBaseURL,
	}
}

func (h *QueryHandlers) ListNodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	nodes, err := h.repo.ListNodes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func (h *QueryHandlers) GetNode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	node, err := h.repo.GetNode(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}

func (h *QueryHandlers) ListDevices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	devices, err := h.repo.ListDevices(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func (h *QueryHandlers) GetDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	device, err := h.repo.GetDevice(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

func (h *QueryHandlers) ListSources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sources, err := h.repo.ListSources(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sources)
}

func (h *QueryHandlers) GetSource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	source, err := h.repo.GetSource(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(source)
}

func (h *QueryHandlers) ListFlows(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	flows, err := h.repo.ListFlows(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flows)
}

func (h *QueryHandlers) GetFlow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	flow, err := h.repo.GetFlow(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flow)
}

func (h *QueryHandlers) ListSenders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	senders, err := h.repo.ListSenders(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(senders)
}

func (h *QueryHandlers) GetSender(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	sender, err := h.repo.GetSender(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sender)
}

func (h *QueryHandlers) ListReceivers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	receivers, err := h.repo.ListReceivers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receivers)
}

func (h *QueryHandlers) GetReceiver(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	receiver, err := h.repo.GetReceiver(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receiver)
}

func (h *QueryHandlers) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ResourcePath == "" {
		http.Error(w, "resource_path is required", http.StatusBadRequest)
		return
	}

	sub, err := h.subscriptionEngine.CreateSubscription(ctx, req.ResourcePath, req.Params, req.MaxUpdateRateMs, req.Persist, req.SecureWebsocket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wsProtocol := "ws"
	if sub.SecureWebsocket != nil && *sub.SecureWebsocket {
		wsProtocol = "wss"
	}
	wsHref := fmt.Sprintf("%s://%s/x-nmos/query/v1.3/subscriptions/%s/ws", wsProtocol, h.wsBaseURL, sub.ID)
	sub.WsHref = &wsHref

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

func (h *QueryHandlers) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	sub, err := h.subscriptionEngine.GetSubscription(ctx, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *QueryHandlers) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subs, err := h.subscriptionEngine.ListSubscriptions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}

func (h *QueryHandlers) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	if err := h.subscriptionEngine.DeleteSubscription(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	SecureWebsocket *bool           `json:"secure_websocket,omitempty"`
}
