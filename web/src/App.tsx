import { useState } from "react";
import AIRecommend from "./pages/AIRecommend";
import Dashboard from "./pages/Dashboard";
import Meals from "./pages/Meals";
import Settings from "./pages/Settings";
import WeightPage from "./pages/Weight";
import Workouts from "./pages/Workouts";

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
		<div className="min-h-screen bg-gray-50 md:flex">
			{/* サイドナビ (md+) */}
			<aside className="hidden md:flex md:flex-col md:w-64 md:shrink-0 md:border-r md:border-gray-200 md:bg-white md:min-h-screen md:sticky md:top-0">
				<div className="px-6 py-5 text-xl font-bold text-emerald-700">
					CUTAGENT
				</div>
				<nav className="flex flex-col px-2 gap-1">
					{NAV.map((n) => (
						<button
							key={n.id}
							type="button"
							onClick={() => setPage(n.id)}
							className={`flex items-center gap-3 px-4 py-3 rounded-lg text-sm text-left ${
								page === n.id
									? "bg-emerald-50 text-emerald-700 font-semibold"
									: "text-gray-600 hover:bg-gray-100"
							}`}
						>
							<span className="text-lg">{n.icon}</span>
							{n.label}
						</button>
					))}
				</nav>
			</aside>

			{/* メインコンテンツ */}
			<main className="flex-1 overflow-y-auto pb-20 md:pb-0">
				<div className="max-w-md mx-auto md:max-w-3xl">
					{page === "dashboard" && <Dashboard date={today} />}
					{page === "meals" && <Meals date={today} />}
					{page === "workouts" && <Workouts date={today} />}
					{page === "weight" && <WeightPage />}
					{page === "ai" && <AIRecommend date={today} />}
					{page === "settings" && <Settings />}
				</div>
			</main>

			{/* ボトムナビ (sm のみ) */}
			<nav className="fixed bottom-0 left-1/2 -translate-x-1/2 w-full max-w-md bg-white border-t border-gray-200 flex md:hidden">
				{NAV.map((n) => (
					<button
						key={n.id}
						type="button"
						onClick={() => setPage(n.id)}
						className={`flex-1 flex flex-col items-center py-2 text-xs gap-0.5 ${
							page === n.id ? "text-emerald-600 font-semibold" : "text-gray-500"
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
