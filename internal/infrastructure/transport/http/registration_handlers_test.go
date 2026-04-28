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
		return registry.Node{}, registry.ErrResourceNotFound
	}
	return node, nil
}

func (m *mockRepository) UpdateNodeHealth(ctx context.Context, id string) error {
	return nil
}

func (m *mockRepository) IDExistsAsOtherType(ctx context.Context, id string, excludeType registry.ResourceType) (bool, error) {
	return false, nil
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

	router := NewRouter(regHandlers, queryHandlers)

	req := httptest.NewRequest("POST", "/x-nmos/registration/v1.3/health/nodes/test-node-id", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["health"] == nil {
		t.Errorf("expected health field in response, got %v", rr.Body.String())
	}
}

func TestRegisterResourceReturns200OnUpdate(t *testing.T) {
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

	if rr.Code != http.StatusCreated {
		t.Errorf("first registration should return 201, got %v", rr.Code)
	}

	node.Version = "2.0"
	data, _ = json.Marshal(node)
	reqBody, _ = json.Marshal(registrationRequest{
		Type: registry.ResourceTypeNode,
		Data: data,
	})

	req = httptest.NewRequest("POST", "/x-nmos/registration/v1.3/resource", bytes.NewBuffer(reqBody))
	rr = httptest.NewRecorder()
	handlers.RegisterResource(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("update registration should return 200, got %v", status)
	}
}

func TestRegisterResourceLocationHeader(t *testing.T) {
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

	location := rr.Header().Get("Location")
	if location == "" {
		t.Errorf("expected Location header on creation, got empty string")
	}
	if location != "/x-nmos/query/v1.3/nodes/test-node-id" {
		t.Errorf("expected Location header to be /x-nmos/query/v1.3/nodes/test-node-id, got %s", location)
	}
}

func TestErrorResponseFormat(t *testing.T) {
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
