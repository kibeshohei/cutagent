// Package schema はアプリ全体で使うデータモデルを定義する。
// Huma がこれらの構造体から OpenAPI スキーマを自動生成し、web 側は
// その OpenAPI から TS 型を生成する（型の二重管理を避ける）。
package schema

// Goal はダイエットの目標設定。
type Goal struct {
	TargetWeight  float64 `json:"targetWeight" doc:"目標体重 (kg)" example:"65"`
	TargetDate    string  `json:"targetDate" doc:"目標日 (ISO 8601)" example:"2026-08-01"`
	CurrentWeight float64 `json:"currentWeight" doc:"現在体重 (kg)" example:"72.5"`
}

// WeightLog は1日の体重記録。
type WeightLog struct {
	ID     string  `json:"id,omitempty" doc:"ID（サーバー生成）"`
	Date   string  `json:"date" doc:"記録日 (ISO 8601)" example:"2026-05-28"`
	Weight float64 `json:"weight" doc:"体重 (kg, 小数点1桁)" example:"71.3"`
}

// MealLog は食事記録。
type MealLog struct {
	ID       string `json:"id,omitempty" doc:"ID（サーバー生成）"`
	Date     string `json:"date" doc:"記録日 (ISO 8601)" example:"2026-05-28"`
	MealType string `json:"mealType" enum:"breakfast,lunch,dinner,snack" doc:"食事区分" example:"lunch"`
	Name     string `json:"name" doc:"メニュー名" example:"鶏むね肉のサラダ"`
	Calories int    `json:"calories" doc:"摂取カロリー (kcal)" example:"350"`
}

// WorkoutLog は運動記録。
type WorkoutLog struct {
	ID              string `json:"id,omitempty" doc:"ID（サーバー生成）"`
	Date            string `json:"date" doc:"記録日 (ISO 8601)" example:"2026-05-28"`
	Name            string `json:"name" doc:"種目名" example:"ジョギング"`
	DurationMinutes int    `json:"durationMinutes" doc:"運動時間 (分)" example:"30"`
	CaloriesBurned  int    `json:"caloriesBurned" doc:"消費カロリー (kcal)" example:"240"`
}

// MealMaster はプリセット食品。
type MealMaster struct {
	ID       string `json:"id" doc:"ID"`
	Name     string `json:"name" doc:"食品名" example:"ゆで卵"`
	Calories int    `json:"calories" doc:"カロリー (kcal)" example:"80"`
}

// WorkoutMaster はプリセットワークアウト。
type WorkoutMaster struct {
	ID              string `json:"id" doc:"ID"`
	Name            string `json:"name" doc:"種目名" example:"ウォーキング"`
	CaloriesPerHour int    `json:"caloriesPerHour" doc:"1時間あたりの消費カロリー (kcal)" example:"280"`
}

// Summary は指定日の摂取・消費・収支のまとめ。
type Summary struct {
	Date      string `json:"date" doc:"対象日 (ISO 8601)" example:"2026-05-28"`
	Intake    int    `json:"intake" doc:"摂取カロリー合計 (kcal)" example:"1600"`
	Burned    int    `json:"burned" doc:"消費カロリー合計 (kcal)" example:"400"`
	Balance   int    `json:"balance" doc:"収支 = 摂取 - 消費 (kcal)" example:"1200"`
	Remaining int    `json:"remaining" doc:"目標摂取量に対する残カロリー (kcal)" example:"200"`
}
