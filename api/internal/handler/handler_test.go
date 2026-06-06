package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/kibeshohei/cutagent/api/internal/handler"
	"github.com/kibeshohei/cutagent/api/internal/repository"
)

func newTestAPI(t *testing.T) humatest.TestAPI {
	t.Helper()
	_, api := humatest.New(t)
	handler.RegisterAll(api, repository.NewStore())
	return api
}

func TestHealth(t *testing.T) {
	api := newTestAPI(t)
	resp := api.Get("/api/health")
	if resp.Code != http.StatusOK {
		t.Fatalf("health: got %d, want 200", resp.Code)
	}
}

func TestWeightCreateAndList(t *testing.T) {
	api := newTestAPI(t)

	resp := api.Post("/api/weight", map[string]any{"date": "2026-05-28", "weight": 71.3})
	if resp.Code != http.StatusCreated {
		t.Fatalf("add weight: got %d, want 201 (%s)", resp.Code, resp.Body.String())
	}

	resp = api.Get("/api/weight")
	if resp.Code != http.StatusOK {
		t.Fatalf("list weight: got %d, want 200", resp.Code)
	}
}

func TestMealLifecycle(t *testing.T) {
	api := newTestAPI(t)

	resp := api.Post("/api/meals", map[string]any{
		"date": "2026-05-28", "mealType": "lunch", "name": "サラダ", "calories": 350,
	})
	if resp.Code != http.StatusCreated {
		t.Fatalf("add meal: got %d, want 201 (%s)", resp.Code, resp.Body.String())
	}

	resp = api.Get("/api/summary?date=2026-05-28")
	if resp.Code != http.StatusOK {
		t.Fatalf("summary: got %d, want 200", resp.Code)
	}
}

func TestMasters(t *testing.T) {
	api := newTestAPI(t)
	if resp := api.Get("/api/meal-master"); resp.Code != http.StatusOK {
		t.Fatalf("meal-master: got %d, want 200", resp.Code)
	}
	if resp := api.Get("/api/workout-master"); resp.Code != http.StatusOK {
		t.Fatalf("workout-master: got %d, want 200", resp.Code)
	}
}

// 体重が 1 件もない状態で /api/weight が JSON null ではなく [] を返すこと。
// null だとフロント側で `logs.at(-1)` が TypeError になり Weight ページが真っ白になる。
func TestWeightListEmptyReturnsArray(t *testing.T) {
	api := newTestAPI(t)
	resp := api.Get("/api/weight")
	if resp.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", resp.Code)
	}
	body := resp.Body.Bytes()
	var logs []map[string]any
	if err := json.Unmarshal(body, &logs); err != nil {
		t.Fatalf("unmarshal: %v body=%q", err, body)
	}
	if logs == nil {
		t.Fatalf("must be [] not null, got body=%q", body)
	}
	if len(logs) != 0 {
		t.Fatalf("must be empty, got %d items", len(logs))
	}
}
