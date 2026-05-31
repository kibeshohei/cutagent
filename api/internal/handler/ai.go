package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kibeshohei/cutagent/api/internal/repository"
	"github.com/kibeshohei/cutagent/api/internal/schema"
	"google.golang.org/genai"
)

const (
	defaultPromptsDir = "prompts"
	geminiModel       = "gemini-2.5-flash"
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

type recommendMealInput struct {
	Body struct {
		Date        string `json:"date" doc:"対象日 (ISO 8601)" example:"2026-05-28"`
		Preferences string `json:"preferences,omitempty" doc:"簡易な好み" example:"和食、低脂質"`
	}
}

type recommendMealOutput struct {
	Body struct {
		Candidates []MealCandidate `json:"candidates"`
	}
}

type recommendWorkoutInput struct {
	Body struct {
		Date        string `json:"date" doc:"対象日 (ISO 8601)" example:"2026-05-28"`
		Preferences string `json:"preferences,omitempty" doc:"簡易な好み" example:"室内、膝に優しい"`
	}
}

type recommendWorkoutOutput struct {
	Body struct {
		Candidates []WorkoutCandidate `json:"candidates"`
	}
}

func RegisterAI(api huma.API, store *repository.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "recommend-meal",
		Method:      "POST",
		Path:        "/api/ai/recommend-meal",
		Summary:     "献立レコメンド",
		Description: "目標と当日の残カロリーを踏まえた献立候補を返す。" +
			"GEMINI_API_KEY が設定されていれば Gemini を呼び、未設定または失敗時はスタブを返す。" +
			"プロンプトは prompts/recommend-meal.md を参照。",
		Tags: []string{"ai"},
	}, func(ctx context.Context, in *recommendMealInput) (*recommendMealOutput, error) {
		out := &recommendMealOutput{}
		goal := store.Goal()
		summary := store.Summary(in.Body.Date)
		if candidates, err := generateMealCandidates(ctx, goal, summary, in.Body.Preferences); err == nil {
			out.Body.Candidates = candidates
		} else {
			out.Body.Candidates = stubMealCandidates()
		}
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "recommend-workout",
		Method:      "POST",
		Path:        "/api/ai/recommend-workout",
		Summary:     "ワークアウトレコメンド",
		Description: "目標と当日の収支を踏まえた運動候補を返す。" +
			"GEMINI_API_KEY が設定されていれば Gemini を呼び、未設定または失敗時はスタブを返す。" +
			"プロンプトは prompts/recommend-workout.md を参照。",
		Tags: []string{"ai"},
	}, func(ctx context.Context, in *recommendWorkoutInput) (*recommendWorkoutOutput, error) {
		out := &recommendWorkoutOutput{}
		goal := store.Goal()
		summary := store.Summary(in.Body.Date)
		if candidates, err := generateWorkoutCandidates(ctx, goal, summary, in.Body.Preferences); err == nil {
			out.Body.Candidates = candidates
		} else {
			out.Body.Candidates = stubWorkoutCandidates()
		}
		return out, nil
	})
}

func generateMealCandidates(ctx context.Context, goal schema.Goal, summary schema.Summary, prefs string) ([]MealCandidate, error) {
	prompt, err := buildPrompt("recommend-meal.md", map[string]string{
		"target_weight":  fmt.Sprintf("%g", goal.TargetWeight),
		"target_date":    goal.TargetDate,
		"current_weight": fmt.Sprintf("%g", goal.CurrentWeight),
		"remaining_kcal": fmt.Sprintf("%d", summary.Remaining),
		"preferences":    prefs,
	})
	if err != nil {
		return nil, fmt.Errorf("build prompt: %w", err)
	}
	text, err := callGemini(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("call gemini: %w", err)
	}
	var out []MealCandidate
	if err := json.Unmarshal([]byte(extractJSON(text)), &out); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("empty candidates")
	}
	return out, nil
}

func generateWorkoutCandidates(ctx context.Context, goal schema.Goal, summary schema.Summary, prefs string) ([]WorkoutCandidate, error) {
	prompt, err := buildPrompt("recommend-workout.md", map[string]string{
		"target_weight":  fmt.Sprintf("%g", goal.TargetWeight),
		"target_date":    goal.TargetDate,
		"current_weight": fmt.Sprintf("%g", goal.CurrentWeight),
		"intake_kcal":    fmt.Sprintf("%d", summary.Intake),
		"burned_kcal":    fmt.Sprintf("%d", summary.Burned),
		"balance_kcal":   fmt.Sprintf("%d", summary.Balance),
		"preferences":    prefs,
	})
	if err != nil {
		return nil, fmt.Errorf("build prompt: %w", err)
	}
	text, err := callGemini(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("call gemini: %w", err)
	}
	var out []WorkoutCandidate
	if err := json.Unmarshal([]byte(extractJSON(text)), &out); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("empty candidates")
	}
	return out, nil
}

// buildPrompt は prompts ディレクトリの Markdown を読み {{var}} を置換する。
// PROMPTS_DIR 環境変数があればそのディレクトリを使う。
func buildPrompt(name string, vars map[string]string) (string, error) {
	dir := os.Getenv("PROMPTS_DIR")
	if dir == "" {
		dir = defaultPromptsDir
	}
	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return "", err
	}
	tmpl := string(b)
	for k, v := range vars {
		tmpl = strings.ReplaceAll(tmpl, "{{"+k+"}}", v)
	}
	return tmpl, nil
}

// callGemini は GEMINI_API_KEY を使って Gemini にプロンプトを送る。
// キー未設定ならエラーを返し、ハンドラ側でスタブにフォールバックさせる。
func callGemini(ctx context.Context, prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set")
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("new client: %w", err)
	}
	resp, err := client.Models.GenerateContent(ctx, geminiModel, genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("generate content: %w", err)
	}
	return resp.Text(), nil
}

// extractJSON は Gemini の出力から JSON 部分だけを取り出す。
// ```json ... ``` ブロックで囲まれていても剥がして本文を返す。
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		if i := strings.Index(s, "\n"); i > 0 {
			s = s[i+1:]
		} else {
			s = strings.TrimPrefix(s, "```")
		}
		s = strings.TrimSuffix(strings.TrimSpace(s), "```")
		s = strings.TrimSpace(s)
	}
	return s
}

func stubMealCandidates() []MealCandidate {
	return []MealCandidate{
		{Name: "鶏むね肉のサラダボウル", Calories: 350, Reason: "高タンパク・低脂質で満腹感が高い"},
		{Name: "豆腐とわかめの味噌汁セット", Calories: 320, Reason: "残カロリーに収まり塩分も控えめ"},
		{Name: "玄米おにぎりとゆで卵", Calories: 300, Reason: "腹持ちが良く間食を防げる"},
	}
}

func stubWorkoutCandidates() []WorkoutCandidate {
	return []WorkoutCandidate{
		{Name: "早歩き", DurationMinutes: 30, CaloriesBurned: 150, Reason: "膝に優しく続けやすい"},
		{Name: "自重スクワット", DurationMinutes: 15, CaloriesBurned: 90, Reason: "大きな筋肉を使い消費効率が高い"},
		{Name: "サイクリング", DurationMinutes: 30, CaloriesBurned: 210, Reason: "今日の収支を一気に埋められる"},
	}
}
