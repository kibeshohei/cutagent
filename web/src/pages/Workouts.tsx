import { useEffect, useState } from "react";
import api, { type WorkoutLog, type WorkoutMaster } from "../api";

export default function Workouts({ date }: { date: string }) {
  const [workouts, setWorkouts] = useState<WorkoutLog[]>([]);
  const [master, setMaster] = useState<WorkoutMaster[]>([]);
  const [selectedId, setSelectedId] = useState("");
  const [duration, setDuration] = useState(30);

  useEffect(() => {
    Promise.all([api.listWorkouts(date), api.workoutMaster()]).then(
      ([w, wm]) => {
        setWorkouts(w);
        setMaster(wm);
        if (wm.length) setSelectedId(wm[0].id);
      },
    );
  }, [date]);

  const add = async () => {
    const item = master.find((m) => m.id === selectedId);
    if (!item) return;
    const burned = Math.round((item.caloriesPerHour * duration) / 60);
    const created = await api.addWorkout({
      date,
      name: item.name,
      durationMinutes: duration,
      caloriesBurned: burned,
    });
    setWorkouts((prev) => [...prev, created]);
  };

  const remove = async (id: string) => {
    await api.deleteWorkout(id);
    setWorkouts((prev) => prev.filter((w) => w.id !== id));
  };

  const total = workouts.reduce((s, w) => s + w.caloriesBurned, 0);

  return (
    <div className="p-4 space-y-4">
      <h1 className="text-xl font-bold text-emerald-700">運動記録</h1>
      <p className="text-sm text-gray-500">{date}　消費 {total} kcal</p>

      <div className="bg-white rounded-xl p-4 shadow-sm space-y-3">
        <select
          value={selectedId}
          onChange={(e) => setSelectedId(e.target.value)}
          className="w-full border border-gray-200 rounded-lg p-2 text-sm"
        >
          {master.map((m) => (
            <option key={m.id} value={m.id}>
              {m.name}（{m.caloriesPerHour} kcal/h）
            </option>
          ))}
        </select>
        <div className="flex items-center gap-2">
          <label className="text-sm text-gray-600 whitespace-nowrap">
            時間
          </label>
          <input
            type="number"
            min={1}
            max={240}
            value={duration}
            onChange={(e) => setDuration(Number(e.target.value))}
            className="w-20 border border-gray-200 rounded-lg p-2 text-sm text-center"
          />
          <span className="text-sm text-gray-500">分</span>
        </div>
        <button
          type="button"
          onClick={add}
          className="w-full py-2 bg-emerald-500 text-white rounded-lg font-medium text-sm"
        >
          追加
        </button>
      </div>

      <ul className="space-y-2">
        {workouts.map((w) => (
          <li
            key={w.id}
            className="bg-white rounded-xl p-3 shadow-sm flex justify-between items-center"
          >
            <div>
              <p className="font-medium text-sm">{w.name}</p>
              <p className="text-xs text-gray-400">
                {w.durationMinutes}分 · {w.caloriesBurned} kcal
              </p>
            </div>
            <button
              type="button"
              onClick={() => remove(w.id!)}
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
