package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kibeshohei/cutagent/api/internal/repository"
)

// MealCandidate は献立レコメンドの1候補。
type MealCandidate struct {
	Name     string `json:"name" doc:"メニュー名" example:"豆腐とわかめの味噌汁セット"`
	Calories int    `json:"calories" doc:"概算カロリー (kcal)" example:"320"`
	Reason   string `json:"reason" doc:"おすすめ理由" example:"残カロリーに収まり高タンパク"`
}

// WorkoutCandidate はワークアウトレコメンドの1候補。
type WorkoutCandidate struct {
	Name            string `json:"name" doc:"種目名" example:"早歩き"`
	DurationMinutes int    `json:"durationMinutes" doc:"推奨時間 (分)" example:"30"`
	CaloriesBurned  int    `json:"caloriesBurned" doc:"概算消費カロリー (kcal)" example:"150"`
	Reason          string `json:"reason" doc:"おすすめ理由" example:"今日の収支を埋めるのに最適"`
}

func RegisterAI(api huma.API, store *repository.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "recommend-meal",
		Method:      "POST",
		Path:        "/api/ai/recommend-meal",
		Summary:     "献立レコメンド",
		Description: "目標と当日の残カロリーを踏まえた献立候補を返す。" +
			"現状はスタブ実装。プロンプトは prompts/recommend-meal.md を参照。",
		Tags: []string{"ai"},
	}, func(_ context.Context, in *struct {
		Body struct {
			Date        string `json:"date" doc:"対象日 (ISO 8601)" example:"2026-05-28"`
			Preferences string `json:"preferences,omitempty" doc:"簡易な好み" example:"和食、低脂質"`
		}
	}) (*struct {
		Body struct {
			Candidates []MealCandidate `json:"candidates"`
		}
	}, error) {
		// TODO: Gemini API (google.golang.org/genai) を呼び出す。
		// API キーは Secret Manager 経由でバックエンドのみが保持する。
		out := &struct {
			Body struct {
				Candidates []MealCandidate `json:"candidates"`
			}
		}{}
		out.Body.Candidates = stubMealCandidates(store)
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "recommend-workout",
		Method:      "POST",
		Path:        "/api/ai/recommend-workout",
		Summary:     "ワークアウトレコメンド",
		Description: "目標と当日の収支を踏まえた運動候補を返す。" +
			"現状はスタブ実装。プロンプトは prompts/recommend-workout.md を参照。",
		Tags: []string{"ai"},
	}, func(_ context.Context, in *struct {
		Body struct {
			Date        string `json:"date" doc:"対象日 (ISO 8601)" example:"2026-05-28"`
			Preferences string `json:"preferences,omitempty" doc:"簡易な好み" example:"室内、膝に優しい"`
		}
	}) (*struct {
		Body struct {
			Candidates []WorkoutCandidate `json:"candidates"`
		}
	}, error) {
		// TODO: Gemini API 呼び出しに差し替える。
		out := &struct {
			Body struct {
				Candidates []WorkoutCandidate `json:"candidates"`
			}
		}{}
		out.Body.Candidates = stubWorkoutCandidates(store)
		return out, nil
	})
}

func stubMealCandidates(_ *repository.Store) []MealCandidate {
	return []MealCandidate{
		{Name: "鶏むね肉のサラダボウル", Calories: 350, Reason: "高タンパク・低脂質で満腹感が高い"},
		{Name: "豆腐とわかめの味噌汁セット", Calories: 320, Reason: "残カロリーに収まり塩分も控えめ"},
		{Name: "玄米おにぎりとゆで卵", Calories: 300, Reason: "腹持ちが良く間食を防げる"},
	}
}

func stubWorkoutCandidates(_ *repository.Store) []WorkoutCandidate {
	return []WorkoutCandidate{
		{Name: "早歩き", DurationMinutes: 30, CaloriesBurned: 150, Reason: "膝に優しく続けやすい"},
		{Name: "自重スクワット", DurationMinutes: 15, CaloriesBurned: 90, Reason: "大きな筋肉を使い消費効率が高い"},
		{Name: "サイクリング", DurationMinutes: 30, CaloriesBurned: 210, Reason: "今日の収支を一気に埋められる"},
	}
}
