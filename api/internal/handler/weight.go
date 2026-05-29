package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kibeshohei/cutagent/api/internal/repository"
	"github.com/kibeshohei/cutagent/api/internal/schema"
)

func RegisterWeight(api huma.API, store *repository.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "list-weights",
		Method:      "GET",
		Path:        "/api/weight",
		Summary:     "体重記録一覧",
		Tags:        []string{"weight"},
	}, func(_ context.Context, _ *struct{}) (*struct {
		Body []schema.WeightLog
	}, error) {
		return &struct{ Body []schema.WeightLog }{Body: store.ListWeights()}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "add-weight",
		Method:        "POST",
		Path:          "/api/weight",
		Summary:       "体重記録追加",
		Tags:          []string{"weight"},
		DefaultStatus: 201,
	}, func(_ context.Context, in *struct {
		Body schema.WeightLog
	}) (*struct {
		Body schema.WeightLog
	}, error) {
		return &struct{ Body schema.WeightLog }{Body: store.AddWeight(in.Body)}, nil
	})
}
