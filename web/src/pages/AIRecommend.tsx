import { useState } from "react";
import api, { type MealCandidate, type WorkoutCandidate } from "../api";

export default function AIRecommend({ date }: { date: string }) {
  const [mealCandidates, setMealCandidates] = useState<MealCandidate[]>([]);
  const [workoutCandidates, setWorkoutCandidates] = useState<
    WorkoutCandidate[]
  >([]);
  const [loadingMeal, setLoadingMeal] = useState(false);
  const [loadingWorkout, setLoadingWorkout] = useState(false);

  const getMealRec = async () => {
    setLoadingMeal(true);
    try {
      const res = await api.recommendMeal(date);
      setMealCandidates(res.candidates);
    } finally {
      setLoadingMeal(false);
    }
  };

  const getWorkoutRec = async () => {
    setLoadingWorkout(true);
    try {
      const res = await api.recommendWorkout(date);
      setWorkoutCandidates(res.candidates);
    } finally {
      setLoadingWorkout(false);
    }
  };

  return (
    <div className="p-4 space-y-6">
      <h1 className="text-xl font-bold text-emerald-700">AI レコメンド</h1>

      <section className="space-y-3">
        <h2 className="font-semibold text-gray-700">献立提案</h2>
        <button
          type="button"
          onClick={getMealRec}
          disabled={loadingMeal}
          className="w-full py-2 bg-emerald-500 text-white rounded-lg font-medium text-sm disabled:opacity-50"
        >
          {loadingMeal ? "生成中…" : "✨ 献立を提案してもらう"}
        </button>
        {mealCandidates.map((c, i) => (
          <div key={i} className="bg-white rounded-xl p-3 shadow-sm">
            <p className="font-medium">{c.name}</p>
            <p className="text-sm text-emerald-600">{c.calories} kcal</p>
            <p className="text-xs text-gray-500 mt-1">{c.reason}</p>
          </div>
        ))}
      </section>

      <section className="space-y-3">
        <h2 className="font-semibold text-gray-700">ワークアウト提案</h2>
        <button
          type="button"
          onClick={getWorkoutRec}
          disabled={loadingWorkout}
          className="w-full py-2 bg-blue-500 text-white rounded-lg font-medium text-sm disabled:opacity-50"
        >
          {loadingWorkout ? "生成中…" : "✨ 運動を提案してもらう"}
        </button>
        {workoutCandidates.map((c, i) => (
          <div key={i} className="bg-white rounded-xl p-3 shadow-sm">
            <p className="font-medium">{c.name}</p>
            <p className="text-sm text-blue-600">
              {c.durationMinutes}分 · {c.caloriesBurned} kcal
            </p>
            <p className="text-xs text-gray-500 mt-1">{c.reason}</p>
          </div>
        ))}
      </section>
    </div>
  );
}
