import { useEffect, useState } from "react";
import api, { type MealLog, type MealMaster } from "../api";

const MEAL_TYPES = [
  { value: "breakfast", label: "朝食" },
  { value: "lunch", label: "昼食" },
  { value: "dinner", label: "夕食" },
  { value: "snack", label: "間食" },
] as const;

export default function Meals({ date }: { date: string }) {
  const [meals, setMeals] = useState<MealLog[]>([]);
  const [master, setMaster] = useState<MealMaster[]>([]);
  const [mealType, setMealType] =
    useState<MealLog["mealType"]>("breakfast");
  const [selectedId, setSelectedId] = useState("");

  useEffect(() => {
    Promise.all([api.listMeals(date), api.mealMaster()]).then(([m, mm]) => {
      setMeals(m);
      setMaster(mm);
      if (mm.length) setSelectedId(mm[0].id);
    });
  }, [date]);

  const add = async () => {
    const item = master.find((m) => m.id === selectedId);
    if (!item) return;
    const created = await api.addMeal({
      date,
      mealType,
      name: item.name,
      calories: item.calories,
    });
    setMeals((prev) => [...prev, created]);
  };

  const remove = async (id: string) => {
    await api.deleteMeal(id);
    setMeals((prev) => prev.filter((m) => m.id !== id));
  };

  const total = meals.reduce((s, m) => s + m.calories, 0);

  return (
    <div className="p-4 space-y-4">
      <h1 className="text-xl font-bold text-emerald-700">食事記録</h1>
      <p className="text-sm text-gray-500">{date}　合計 {total} kcal</p>

      <div className="bg-white rounded-xl p-4 shadow-sm space-y-3">
        <div className="flex gap-2">
          {MEAL_TYPES.map((t) => (
            <button
              key={t.value}
              type="button"
              onClick={() => setMealType(t.value)}
              className={`flex-1 py-1 rounded-lg text-xs font-medium border ${
                mealType === t.value
                  ? "bg-emerald-500 text-white border-emerald-500"
                  : "border-gray-200 text-gray-600"
              }`}
            >
              {t.label}
            </button>
          ))}
        </div>
        <select
          value={selectedId}
          onChange={(e) => setSelectedId(e.target.value)}
          className="w-full border border-gray-200 rounded-lg p-2 text-sm"
        >
          {master.map((m) => (
            <option key={m.id} value={m.id}>
              {m.name}（{m.calories} kcal）
            </option>
          ))}
        </select>
        <button
          type="button"
          onClick={add}
          className="w-full py-2 bg-emerald-500 text-white rounded-lg font-medium text-sm"
        >
          追加
        </button>
      </div>

      <ul className="space-y-2">
        {meals.map((m) => (
          <li
            key={m.id}
            className="bg-white rounded-xl p-3 shadow-sm flex justify-between items-center"
          >
            <div>
              <p className="font-medium text-sm">{m.name}</p>
              <p className="text-xs text-gray-400">
                {MEAL_TYPES.find((t) => t.value === m.mealType)?.label} ·{" "}
                {m.calories} kcal
              </p>
            </div>
            <button
              type="button"
              onClick={() => remove(m.id!)}
              className="text-gray-300 hover:text-red-400 text-lg"
            >
              ×
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
