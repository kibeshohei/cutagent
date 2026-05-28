import { useEffect, useState } from "react";
import api, { type Goal } from "../api";

export default function Settings() {
  const [goal, setGoal] = useState<Goal>({
    targetWeight: 65,
    targetDate: "",
    currentWeight: 72.5,
  });
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    api.getGoal().then(setGoal);
  }, []);

  const save = async () => {
    const updated = await api.updateGoal(goal);
    setGoal(updated);
    setSaved(true);
    setTimeout(() => setSaved(false), 2000);
  };

  return (
    <div className="p-4 space-y-4">
      <h1 className="text-xl font-bold text-emerald-700">設定</h1>

      <div className="bg-white rounded-xl p-4 shadow-sm space-y-4">
        <h2 className="font-semibold text-gray-700">目標設定</h2>

        <Field label="現在体重 (kg)">
          <input
            type="number"
            step="0.1"
            value={goal.currentWeight}
            onChange={(e) =>
              setGoal((g) => ({ ...g, currentWeight: parseFloat(e.target.value) }))
            }
            className="w-full border border-gray-200 rounded-lg p-2 text-sm"
          />
        </Field>

        <Field label="目標体重 (kg)">
          <input
            type="number"
            step="0.1"
            value={goal.targetWeight}
            onChange={(e) =>
              setGoal((g) => ({ ...g, targetWeight: parseFloat(e.target.value) }))
            }
            className="w-full border border-gray-200 rounded-lg p-2 text-sm"
          />
        </Field>

        <Field label="目標日">
          <input
            type="date"
            value={goal.targetDate}
            onChange={(e) =>
              setGoal((g) => ({ ...g, targetDate: e.target.value }))
            }
            className="w-full border border-gray-200 rounded-lg p-2 text-sm"
          />
        </Field>

        <button
          type="button"
          onClick={save}
          className="w-full py-2 bg-emerald-500 text-white rounded-lg font-medium text-sm"
        >
          {saved ? "✓ 保存しました" : "保存する"}
        </button>
      </div>
    </div>
  );
}

function Field({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="space-y-1">
      <label className="text-sm text-gray-600">{label}</label>
      {children}
    </div>
  );
}
