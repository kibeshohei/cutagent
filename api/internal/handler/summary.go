package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/kibeshohei/cutagent/api/internal/repository"
	"github.com/kibeshohei/cutagent/api/internal/schema"
)

func RegisterSummary(api huma.API, store *repository.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "get-summary",
		Method:      "GET",
		Path:        "/api/summary",
		Summary:     "指定日の摂取/消費/収支/残kcal",
		Tags:        []string{"summary"},
	}, func(_ context.Context, in *struct {
		Date string `query:"date" required:"true" doc:"対象日 (ISO 8601)" example:"2026-05-28"`
	}) (*struct {
		Body schema.Summary
	}, error) {
		return &struct{ Body schema.Summary }{Body: store.Summary(in.Date)}, nil
	})
}
