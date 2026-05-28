package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kibeshohei/cutagent/api/internal/repository"
	"github.com/kibeshohei/cutagent/api/internal/schema"
)

func RegisterWorkouts(api huma.API, store *repository.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "list-workouts",
		Method:      "GET",
		Path:        "/api/workouts",
		Summary:     "指定日の運動記録",
		Tags:        []string{"workouts"},
	}, func(_ context.Context, in *struct {
		Date string `query:"date" doc:"対象日 (ISO 8601)。未指定なら全件" example:"2026-05-28"`
	}) (*struct {
		Body []schema.WorkoutLog
	}, error) {
		return &struct{ Body []schema.WorkoutLog }{Body: store.ListWorkouts(in.Date)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "add-workout",
		Method:        "POST",
		Path:          "/api/workouts",
		Summary:       "運動記録追加",
		Tags:          []string{"workouts"},
		DefaultStatus: 201,
	}, func(_ context.Context, in *struct {
		Body schema.WorkoutLog
	}) (*struct {
		Body schema.WorkoutLog
	}, error) {
		return &struct{ Body schema.WorkoutLog }{Body: store.AddWorkout(in.Body)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-workout",
		Method:        "DELETE",
		Path:          "/api/workouts/{id}",
		Summary:       "運動記録削除",
		Tags:          []string{"workouts"},
		DefaultStatus: 204,
	}, func(_ context.Context, in *struct {
		ID string `path:"id"`
	}) (*struct{}, error) {
		if !store.DeleteWorkout(in.ID) {
			return nil, huma.Error404NotFound("workout not found")
		}
		return nil, nil
	})
}
