package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kibeshohei/cutagent/api/internal/repository"
	"github.com/kibeshohei/cutagent/api/internal/schema"
)

func RegisterGoal(api huma.API, store *repository.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "get-goal",
		Method:      "GET",
		Path:        "/api/goal",
		Summary:     "目標取得",
		Tags:        []string{"goal"},
	}, func(_ context.Context, _ *struct{}) (*struct {
		Body schema.Goal
	}, error) {
		return &struct{ Body schema.Goal }{Body: store.Goal()}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-goal",
		Method:      "PUT",
		Path:        "/api/goal",
		Summary:     "目標更新",
		Tags:        []string{"goal"},
	}, func(_ context.Context, in *struct {
		Body schema.Goal
	}) (*struct {
		Body schema.Goal
	}, error) {
		return &struct{ Body schema.Goal }{Body: store.SetGoal(in.Body)}, nil
	})
}
