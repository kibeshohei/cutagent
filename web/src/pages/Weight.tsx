import { useEffect, useState } from "react";
import api, { type WeightLog } from "../api";

export default function WeightPage() {
  const [logs, setLogs] = useState<WeightLog[]>([]);
  const [date, setDate] = useState(new Date().toISOString().slice(0, 10));
  const [weight, setWeight] = useState("");

  useEffect(() => {
    api.listWeights().then(setLogs);
  }, []);

  const add = async () => {
    const w = parseFloat(weight);
    if (Number.isNaN(w)) return;
    const created = await api.addWeight({ date, weight: w });
    setLogs((prev) =>
      [...prev, created].sort((a, b) => a.date.localeCompare(b.date)),
    );
    setWeight("");
  };

  const latest = logs.at(-1);

  return (
    <div className="p-4 space-y-4">
      <h1 className="text-xl font-bold text-emerald-700">体重記録</h1>
      {latest && (
        <p className="text-sm text-gray-500">
          最新: {latest.date} — {latest.weight} kg
        </p>
      )}

      <div className="bg-white rounded-xl p-4 shadow-sm space-y-3">
        <div className="flex gap-2">
          <input
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            className="flex-1 border border-gray-200 rounded-lg p-2 text-sm"
          />
          <input
            type="number"
            step="0.1"
            placeholder="kg"
            value={weight}
            onChange={(e) => setWeight(e.target.value)}
            className="w-24 border border-gray-200 rounded-lg p-2 text-sm text-center"
          />
        </div>
        <button
          type="button"
          onClick={add}
          className="w-full py-2 bg-emerald-500 text-white rounded-lg font-medium text-sm"
        >
          記録する
        </button>
      </div>

      <ul className="space-y-2">
        {[...logs].reverse().map((l) => (
          <li
            key={l.id}
            className="bg-white rounded-xl p-3 shadow-sm flex justify-between"
          >
            <span className="text-sm text-gray-600">{l.date}</span>
            <span className="font-bold text-emerald-700">{l.weight} kg</span>
          </li>
        ))}
      </ul>
    </div>
  );
}
