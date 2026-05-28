import { useState } from "react";
import Dashboard from "./pages/Dashboard";
import Meals from "./pages/Meals";
import Workouts from "./pages/Workouts";
import WeightPage from "./pages/Weight";
import AIRecommend from "./pages/AIRecommend";
import Settings from "./pages/Settings";

type Page = "dashboard" | "meals" | "workouts" | "weight" | "ai" | "settings";

const NAV: { id: Page; label: string; icon: string }[] = [
  { id: "dashboard", label: "ホーム", icon: "🏠" },
  { id: "meals", label: "食事", icon: "🍱" },
  { id: "workouts", label: "運動", icon: "💪" },
  { id: "weight", label: "体重", icon: "⚖️" },
  { id: "ai", label: "AI", icon: "✨" },
  { id: "settings", label: "設定", icon: "⚙️" },
];

export default function App() {
  const [page, setPage] = useState<Page>("dashboard");
  const today = new Date().toISOString().slice(0, 10);

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col max-w-md mx-auto">
      <main className="flex-1 overflow-y-auto pb-20">
        {page === "dashboard" && <Dashboard date={today} />}
        {page === "meals" && <Meals date={today} />}
        {page === "workouts" && <Workouts date={today} />}
        {page === "weight" && <WeightPage />}
        {page === "ai" && <AIRecommend date={today} />}
        {page === "settings" && <Settings />}
      </main>

      <nav className="fixed bottom-0 left-1/2 -translate-x-1/2 w-full max-w-md bg-white border-t border-gray-200 flex">
        {NAV.map((n) => (
          <button
            key={n.id}
            type="button"
            onClick={() => setPage(n.id)}
            className={`flex-1 flex flex-col items-center py-2 text-xs gap-0.5 ${
              page === n.id
                ? "text-emerald-600 font-semibold"
                : "text-gray-500"
            }`}
          >
            <span className="text-lg">{n.icon}</span>
            {n.label}
          </button>
        ))}
      </nav>
    </div>
  );
}
