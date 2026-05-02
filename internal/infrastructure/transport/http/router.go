package transporthttp

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewRouter(regHandlers *RegistrationHandlers, queryHandlers *QueryHandlers) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Route("/x-nmos", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`["query/", "registration/"]`))
		})

		r.Route("/registration", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				// List supported versions
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`["v1.3/"]`))
			})

			r.Route("/v1.3", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(`["resource/", "health/"]`))
				})

				r.Post("/resource", regHandlers.RegisterResource)
				r.Delete("/resource/{type}/{id}", regHandlers.UnregisterResource)

				r.Route("/health", func(r chi.Router) {
					r.Post("/nodes/{nodeID}", regHandlers.Heartbeat)
				})
			})
		})

		r.Route("/query", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				// List supported versions
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`["v1.3/"]`))
			})

			r.Route("/v1.3", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(`["nodes/", "devices/", "sources/", "flows/", "senders/", "receivers/", "subscriptions/"]`))
				})

				r.Get("/nodes", queryHandlers.ListNodes)
				r.Get("/nodes/", queryHandlers.ListNodes)
				r.Get("/nodes/{id}", queryHandlers.GetNode)

				r.Get("/devices", queryHandlers.ListDevices)
				r.Get("/devices/", queryHandlers.ListDevices)
				r.Get("/devices/{id}", queryHandlers.GetDevice)

				r.Get("/sources", queryHandlers.ListSources)
				r.Get("/sources/", queryHandlers.ListSources)
				r.Get("/sources/{id}", queryHandlers.GetSource)

				r.Get("/flows", queryHandlers.ListFlows)
				r.Get("/flows/", queryHandlers.ListFlows)
				r.Get("/flows/{id}", queryHandlers.GetFlow)

				r.Get("/senders", queryHandlers.ListSenders)
				r.Get("/senders/", queryHandlers.ListSenders)
				r.Get("/senders/{id}", queryHandlers.GetSender)

				r.Get("/receivers", queryHandlers.ListReceivers)
				r.Get("/receivers/", queryHandlers.ListReceivers)
				r.Get("/receivers/{id}", queryHandlers.GetReceiver)

				r.Get("/subscriptions", queryHandlers.ListSubscriptions)
				r.Post("/subscriptions", queryHandlers.CreateSubscription)
				r.Get("/subscriptions/{id}", queryHandlers.GetSubscription)
				r.Delete("/subscriptions/{id}", queryHandlers.DeleteSubscription)
				r.Get("/subscriptions/{id}/ws", queryHandlers.HandleSubscriptionWebSocket)

				r.Get("/subscriptions/{id}/*", queryHandlers.HandleSubscriptionWebSocket)
			})
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:  http.StatusNotFound,
			Error: "not found",
		})
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:  http.StatusMethodNotAllowed,
			Error: "method not allowed",
		})
	})

	return r
}
