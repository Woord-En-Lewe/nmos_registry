package transporthttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Woord-En-Lewe/nmos_registry/internal/registry"
)

type mockRepository struct {
	registry.IRepository
	nodes map[string]registry.Node
}

func (m *mockRepository) UpsertNode(ctx context.Context, node registry.Node) error {
	m.nodes[node.ID] = node
	return nil
}

func (m *mockRepository) GetNode(ctx context.Context, id string) (registry.Node, error) {
	node, ok := m.nodes[id]
	if !ok {
		return registry.Node{}, nil
	}
	return node, nil
}

func (m *mockRepository) UpdateNodeHealth(ctx context.Context, id string) error {
	return nil
}

func (m *mockRepository) ListNodes(ctx context.Context) ([]registry.Node, error) {
	nodes := make([]registry.Node, 0, len(m.nodes))
	for _, n := range m.nodes {
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func TestRegisterNode(t *testing.T) {
	repo := &mockRepository{nodes: make(map[string]registry.Node)}
	rm := registry.NewResourceManager(repo)
	he := registry.NewHeartbeatEngine(repo, 0, 0)
	handlers := NewRegistrationHandlers(rm, he)

	node := registry.Node{
		ID:          "test-node-id",
		Version:     "1.0",
		Label:       "Test Node",
		Description: "A test node",
		Api:         json.RawMessage(`{}`),
	}
	data, _ := json.Marshal(node)
	reqBody, _ := json.Marshal(registrationRequest{
		Type: registry.ResourceTypeNode,
		Data: data,
	})

	req := httptest.NewRequest("POST", "/x-nmos/registration/v1.3/resource", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	handlers.RegisterResource(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if len(repo.nodes) != 1 {
		t.Errorf("expected 1 node in repo, got %v", len(repo.nodes))
	}
}

func TestHeartbeat(t *testing.T) {
	repo := &mockRepository{nodes: make(map[string]registry.Node)}
	rm := registry.NewResourceManager(repo)
	he := registry.NewHeartbeatEngine(repo, 0, 0)
	regHandlers := NewRegistrationHandlers(rm, he)
	queryHandlers := NewQueryHandlers(repo, nil, nil, "8080")

	// We need a router to parse the nodeID param
	router := NewRouter(regHandlers, queryHandlers)

	req := httptest.NewRequest("POST", "/x-nmos/registration/v1.3/health/nodes/test-node-id", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
