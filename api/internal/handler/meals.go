package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kibeshohei/cutagent/api/internal/repository"
	"github.com/kibeshohei/cutagent/api/internal/schema"
)

func RegisterMeals(api huma.API, store *repository.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "list-meals",
		Method:      "GET",
		Path:        "/api/meals",
		Summary:     "指定日の食事記録",
		Tags:        []string{"meals"},
	}, func(_ context.Context, in *struct {
		Date string `query:"date" doc:"対象日 (ISO 8601)。未指定なら全件" example:"2026-05-28"`
	}) (*struct {
		Body []schema.MealLog
	}, error) {
		return &struct{ Body []schema.MealLog }{Body: store.ListMeals(in.Date)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "add-meal",
		Method:        "POST",
		Path:          "/api/meals",
		Summary:       "食事記録追加",
		Tags:          []string{"meals"},
		DefaultStatus: 201,
	}, func(_ context.Context, in *struct {
		Body schema.MealLog
	}) (*struct {
		Body schema.MealLog
	}, error) {
		return &struct{ Body schema.MealLog }{Body: store.AddMeal(in.Body)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-meal",
		Method:        "DELETE",
		Path:          "/api/meals/{id}",
		Summary:       "食事記録削除",
		Tags:          []string{"meals"},
		DefaultStatus: 204,
	}, func(_ context.Context, in *struct {
		ID string `path:"id"`
	}) (*struct{}, error) {
		if !store.DeleteMeal(in.ID) {
			return nil, huma.Error404NotFound("meal not found")
		}
		return nil, nil
	})
}
