// Package handler は HTTP エンドポイントを Huma に登録する。
package handler

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/kibeshohei/cutagent/api/internal/repository"
)

// RegisterAll は全エンドポイントを api に登録する。
func RegisterAll(api huma.API, store *repository.Store) {
	RegisterHealth(api)
	RegisterWeight(api, store)
	RegisterMeals(api, store)
	RegisterWorkouts(api, store)
	RegisterMasters(api, store)
	RegisterSummary(api, store)
	RegisterGoal(api, store)
	RegisterAI(api, store)
}
