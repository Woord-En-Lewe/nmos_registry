package transporthttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates a new chi router with NMOS routes
func NewRouter(regHandlers *RegistrationHandlers, queryHandlers *QueryHandlers) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/x-nmos", func(r chi.Router) {
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
				r.Get("/nodes/{id}", queryHandlers.GetNode)

				r.Get("/devices", queryHandlers.ListDevices)
				r.Get("/devices/{id}", queryHandlers.GetDevice)

				r.Get("/sources", queryHandlers.ListSources)
				r.Get("/sources/{id}", queryHandlers.GetSource)

				r.Get("/flows", queryHandlers.ListFlows)
				r.Get("/flows/{id}", queryHandlers.GetFlow)

				r.Get("/senders", queryHandlers.ListSenders)
				r.Get("/senders/{id}", queryHandlers.GetSender)

				r.Get("/receivers", queryHandlers.ListReceivers)
				r.Get("/receivers/{id}", queryHandlers.GetReceiver)
			})
		})
	})

	return r
}
