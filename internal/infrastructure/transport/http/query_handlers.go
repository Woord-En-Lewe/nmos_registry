package transporthttp

import (
	"encoding/json"
	"net/http"

	"github.com/Woord-En-Lewe/nmos_registry/internal/registry"
	"github.com/go-chi/chi/v5"
)

type QueryHandlers struct {
	repo registry.IRepository
}

func NewQueryHandlers(repo registry.IRepository) *QueryHandlers {
	return &QueryHandlers{
		repo: repo,
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
