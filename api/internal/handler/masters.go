package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kibeshohei/cutagent/api/internal/repository"
	"github.com/kibeshohei/cutagent/api/internal/schema"
)

func RegisterMasters(api huma.API, store *repository.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "list-meal-master",
		Method:      "GET",
		Path:        "/api/meal-master",
		Summary:     "プリセット食品一覧",
		Tags:        []string{"masters"},
	}, func(_ context.Context, _ *struct{}) (*struct {
		Body []schema.MealMaster
	}, error) {
		return &struct{ Body []schema.MealMaster }{Body: store.MealMaster()}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-workout-master",
		Method:      "GET",
		Path:        "/api/workout-master",
		Summary:     "プリセットワークアウト一覧",
		Tags:        []string{"masters"},
	}, func(_ context.Context, _ *struct{}) (*struct {
		Body []schema.WorkoutMaster
	}, error) {
		return &struct{ Body []schema.WorkoutMaster }{Body: store.WorkoutMaster()}, nil
	})
}
