// OpenAPI から型生成する予定。現状は手書き（型の構造は api/internal/schema に準拠）。
// TODO: Huma が生成する openapi.json から openapi-typescript 等で自動生成に切り替える。

export type Goal = {
  targetWeight: number;
  targetDate: string;
  currentWeight: number;
};

export type WeightLog = {
  id?: string;
  date: string;
  weight: number;
};

export type MealLog = {
  id?: string;
  date: string;
  mealType: "breakfast" | "lunch" | "dinner" | "snack";
  name: string;
  calories: number;
};

export type WorkoutLog = {
  id?: string;
  date: string;
  name: string;
  durationMinutes: number;
  caloriesBurned: number;
};

export type MealMaster = { id: string; name: string; calories: number };
export type WorkoutMaster = {
  id: string;
  name: string;
  caloriesPerHour: number;
};

export type Summary = {
  date: string;
  intake: number;
  burned: number;
  balance: number;
  remaining: number;
};

export type MealCandidate = {
  name: string;
  calories: number;
  reason: string;
};

export type WorkoutCandidate = {
  name: string;
  durationMinutes: number;
  caloriesBurned: number;
  reason: string;
};

async function request<T>(
  path: string,
  init?: RequestInit,
): Promise<T> {
  const res = await fetch(path, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

const api = {
  getGoal: () => request<Goal>("/api/goal"),
  updateGoal: (body: Goal) =>
    request<Goal>("/api/goal", { method: "PUT", body: JSON.stringify(body) }),

  listWeights: () => request<WeightLog[]>("/api/weight"),
  addWeight: (body: WeightLog) =>
    request<WeightLog>("/api/weight", {
      method: "POST",
      body: JSON.stringify(body),
    }),

  listMeals: (date: string) => request<MealLog[]>(`/api/meals?date=${date}`),
  addMeal: (body: MealLog) =>
    request<MealLog>("/api/meals", {
      method: "POST",
      body: JSON.stringify(body),
    }),
  deleteMeal: (id: string) =>
    request<void>(`/api/meals/${id}`, { method: "DELETE" }),

  listWorkouts: (date: string) =>
    request<WorkoutLog[]>(`/api/workouts?date=${date}`),
  addWorkout: (body: WorkoutLog) =>
    request<WorkoutLog>("/api/workouts", {
      method: "POST",
      body: JSON.stringify(body),
    }),
  deleteWorkout: (id: string) =>
    request<void>(`/api/workouts/${id}`, { method: "DELETE" }),

  mealMaster: () => request<MealMaster[]>("/api/meal-master"),
  workoutMaster: () => request<WorkoutMaster[]>("/api/workout-master"),

  getSummary: (date: string) =>
    request<Summary>(`/api/summary?date=${date}`),

  recommendMeal: (date: string, preferences = "") =>
    request<{ candidates: MealCandidate[] }>("/api/ai/recommend-meal", {
      method: "POST",
      body: JSON.stringify({ date, preferences }),
    }),
  recommendWorkout: (date: string, preferences = "") =>
    request<{ candidates: WorkoutCandidate[] }>("/api/ai/recommend-workout", {
      method: "POST",
      body: JSON.stringify({ date, preferences }),
    }),
};

export default api;
