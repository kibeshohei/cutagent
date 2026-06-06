// Package repository はデータ永続化を担う。
// 現状はインメモリ実装のみ。将来 Firestore 実装に差し替える
// （その際は Repository インターフェースを切る想定）。
package repository

import (
	"crypto/rand"
	"encoding/hex"
	"sync"

	"github.com/kibeshohei/cutagent/api/internal/schema"
)

// DailyIntakeTarget は残カロリー算出に使う暫定の1日あたり目標摂取量。
// TODO: 目標体重・目標日から動的に算出する（§5.6）。
const DailyIntakeTarget = 1800

// Store はスレッドセーフなインメモリストア。
type Store struct {
	mu sync.RWMutex

	weights  []schema.WeightLog
	meals    []schema.MealLog
	workouts []schema.WorkoutLog
	goal     schema.Goal

	mealMaster    []schema.MealMaster
	workoutMaster []schema.WorkoutMaster
}

// NewStore はプリセットマスターを seed 済みの Store を返す。
func NewStore() *Store {
	return &Store{
		goal: schema.Goal{TargetWeight: 65, TargetDate: "2026-08-01", CurrentWeight: 72.5},
		mealMaster: []schema.MealMaster{
			{ID: "m1", Name: "ゆで卵", Calories: 80},
			{ID: "m2", Name: "鶏むね肉のサラダ", Calories: 350},
			{ID: "m3", Name: "玄米ごはん 150g", Calories: 240},
			{ID: "m4", Name: "プロテイン", Calories: 120},
			{ID: "m5", Name: "納豆", Calories: 100},
		},
		workoutMaster: []schema.WorkoutMaster{
			{ID: "w1", Name: "ウォーキング", CaloriesPerHour: 280},
			{ID: "w2", Name: "ジョギング", CaloriesPerHour: 480},
			{ID: "w3", Name: "筋トレ", CaloriesPerHour: 360},
			{ID: "w4", Name: "サイクリング", CaloriesPerHour: 420},
		},
	}
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// --- 体重 ---

func (s *Store) ListWeights() []schema.WeightLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []schema.WeightLog{}
	return append(out, s.weights...)
}

func (s *Store) AddWeight(w schema.WeightLog) schema.WeightLog {
	s.mu.Lock()
	defer s.mu.Unlock()
	w.ID = newID()
	s.weights = append(s.weights, w)
	return w
}

// --- 食事 ---

func (s *Store) ListMeals(date string) []schema.MealLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []schema.MealLog{}
	for _, m := range s.meals {
		if date == "" || m.Date == date {
			out = append(out, m)
		}
	}
	return out
}

func (s *Store) AddMeal(m schema.MealLog) schema.MealLog {
	s.mu.Lock()
	defer s.mu.Unlock()
	m.ID = newID()
	s.meals = append(s.meals, m)
	return m
}

func (s *Store) DeleteMeal(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, m := range s.meals {
		if m.ID == id {
			s.meals = append(s.meals[:i], s.meals[i+1:]...)
			return true
		}
	}
	return false
}

// --- 運動 ---

func (s *Store) ListWorkouts(date string) []schema.WorkoutLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []schema.WorkoutLog{}
	for _, w := range s.workouts {
		if date == "" || w.Date == date {
			out = append(out, w)
		}
	}
	return out
}

func (s *Store) AddWorkout(w schema.WorkoutLog) schema.WorkoutLog {
	s.mu.Lock()
	defer s.mu.Unlock()
	w.ID = newID()
	s.workouts = append(s.workouts, w)
	return w
}

func (s *Store) DeleteWorkout(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, w := range s.workouts {
		if w.ID == id {
			s.workouts = append(s.workouts[:i], s.workouts[i+1:]...)
			return true
		}
	}
	return false
}

// --- マスター ---

func (s *Store) MealMaster() []schema.MealMaster {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []schema.MealMaster{}
	return append(out, s.mealMaster...)
}

func (s *Store) WorkoutMaster() []schema.WorkoutMaster {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []schema.WorkoutMaster{}
	return append(out, s.workoutMaster...)
}

// --- 目標 ---

func (s *Store) Goal() schema.Goal {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.goal
}

func (s *Store) SetGoal(g schema.Goal) schema.Goal {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.goal = g
	return s.goal
}

// --- サマリー ---

func (s *Store) Summary(date string) schema.Summary {
	s.mu.RLock()
	defer s.mu.RUnlock()
	intake := 0
	for _, m := range s.meals {
		if m.Date == date {
			intake += m.Calories
		}
	}
	burned := 0
	for _, w := range s.workouts {
		if w.Date == date {
			burned += w.CaloriesBurned
		}
	}
	return schema.Summary{
		Date:      date,
		Intake:    intake,
		Burned:    burned,
		Balance:   intake - burned,
		Remaining: DailyIntakeTarget - intake + burned,
	}
}
