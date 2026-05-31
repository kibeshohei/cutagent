package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/kibeshohei/cutagent/api/internal/handler"
	"github.com/kibeshohei/cutagent/api/internal/repository"
)

type mealResp struct {
	Candidates []handler.MealCandidate `json:"candidates"`
}

type workoutResp struct {
	Candidates []handler.WorkoutCandidate `json:"candidates"`
}

// GEMINI_API_KEY 未設定なら、スタブ候補が必ず3件返ることを検証する。
func TestRecommendMealFallback(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	_, api := humatest.New(t)
	handler.RegisterAll(api, repository.NewStore())

	resp := api.Post("/api/ai/recommend-meal", map[string]any{"date": "2026-05-28"})
	if resp.Code != http.StatusOK {
		t.Fatalf("got %d want 200 (%s)", resp.Code, resp.Body.String())
	}
	var body mealResp
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(body.Candidates) != 3 {
		t.Fatalf("want 3 fallback meals, got %d", len(body.Candidates))
	}
}

func TestRecommendWorkoutFallback(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	_, api := humatest.New(t)
	handler.RegisterAll(api, repository.NewStore())

	resp := api.Post("/api/ai/recommend-workout", map[string]any{"date": "2026-05-28"})
	if resp.Code != http.StatusOK {
		t.Fatalf("got %d want 200 (%s)", resp.Code, resp.Body.String())
	}
	var body workoutResp
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(body.Candidates) != 3 {
		t.Fatalf("want 3 fallback workouts, got %d", len(body.Candidates))
	}
}
