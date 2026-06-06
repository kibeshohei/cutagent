import {
	BrowserRouter,
	Navigate,
	NavLink,
	Route,
	Routes,
} from "react-router-dom";
import AIRecommend from "./pages/AIRecommend";
import Dashboard from "./pages/Dashboard";
import Meals from "./pages/Meals";
import Settings from "./pages/Settings";
import WeightPage from "./pages/Weight";
import Workouts from "./pages/Workouts";

const NAV: { path: string; label: string; icon: string }[] = [
	{ path: "/", label: "ホーム", icon: "🏠" },
	{ path: "/meals", label: "食事", icon: "🍱" },
	{ path: "/workouts", label: "運動", icon: "💪" },
	{ path: "/weight", label: "体重", icon: "⚖️" },
	{ path: "/ai", label: "AI", icon: "✨" },
	{ path: "/settings", label: "設定", icon: "⚙️" },
];

export default function App() {
	const today = new Date().toISOString().slice(0, 10);

	return (
		<BrowserRouter>
			<div className="min-h-screen bg-gray-50 md:flex">
				{/* サイドナビ (md+) */}
				<aside className="hidden md:flex md:flex-col md:w-64 md:shrink-0 md:border-r md:border-gray-200 md:bg-white md:min-h-screen md:sticky md:top-0">
					<div className="px-6 py-5 text-xl font-bold text-emerald-700">
						CUTAGENT
					</div>
					<nav className="flex flex-col px-2 gap-1">
						{NAV.map((n) => (
							<NavLink
								key={n.path}
								to={n.path}
								end={n.path === "/"}
								className={({ isActive }) =>
									`flex items-center gap-3 px-4 py-3 rounded-lg text-sm text-left ${
										isActive
											? "bg-emerald-50 text-emerald-700 font-semibold"
											: "text-gray-600 hover:bg-gray-100"
									}`
								}
							>
								<span className="text-lg">{n.icon}</span>
								{n.label}
							</NavLink>
						))}
					</nav>
				</aside>

				{/* メインコンテンツ */}
				<main className="flex-1 overflow-y-auto pb-20 md:pb-0">
					<div className="max-w-md mx-auto md:max-w-3xl">
						<Routes>
							<Route path="/" element={<Dashboard date={today} />} />
							<Route path="/meals" element={<Meals date={today} />} />
							<Route path="/workouts" element={<Workouts date={today} />} />
							<Route path="/weight" element={<WeightPage />} />
							<Route path="/ai" element={<AIRecommend date={today} />} />
							<Route path="/settings" element={<Settings />} />
							<Route path="*" element={<Navigate to="/" replace />} />
						</Routes>
					</div>
				</main>

				{/* ボトムナビ (sm のみ) */}
				<nav className="fixed bottom-0 left-1/2 -translate-x-1/2 w-full max-w-md bg-white border-t border-gray-200 flex md:hidden">
					{NAV.map((n) => (
						<NavLink
							key={n.path}
							to={n.path}
							end={n.path === "/"}
							className={({ isActive }) =>
								`flex-1 flex flex-col items-center py-2 text-xs gap-0.5 ${
									isActive ? "text-emerald-600 font-semibold" : "text-gray-500"
								}`
							}
						>
							<span className="text-lg">{n.icon}</span>
							{n.label}
						</NavLink>
					))}
				</nav>
			</div>
		</BrowserRouter>
	);
}
