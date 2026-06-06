import { useEffect, useState } from "react";
import api, { type Goal, type Summary } from "../api";

export default function Dashboard({ date }: { date: string }) {
	const [summary, setSummary] = useState<Summary | null>(null);
	const [goal, setGoal] = useState<Goal | null>(null);
	const [error, setError] = useState("");

	useEffect(() => {
		Promise.all([api.getSummary(date), api.getGoal()])
			.then(([s, g]) => {
				setSummary(s);
				setGoal(g);
			})
			.catch((e: unknown) =>
				setError(e instanceof Error ? e.message : "エラーが発生しました"),
			);
	}, [date]);

	if (error) return <ErrorMsg msg={error} />;
	if (!summary || !goal) return <Loading />;

	const targetDays = Math.ceil(
		(new Date(goal.targetDate).getTime() - Date.now()) / 86400000,
	);
	const diffKg = +(goal.currentWeight - goal.targetWeight).toFixed(1);
	const remainText =
		diffKg === 0
			? "目標達成"
			: diffKg > 0
				? `あと ${diffKg} kg 減量`
				: `あと ${Math.abs(diffKg)} kg 増量`;

	return (
		<div className="p-4 space-y-4">
			<h1 className="text-xl font-bold text-emerald-700">今日の状況</h1>
			<p className="text-sm text-gray-500">{date}</p>

			<div className="grid grid-cols-2 gap-3">
				<Card label="摂取" value={`${summary.intake} kcal`} color="orange" />
				<Card label="消費" value={`${summary.burned} kcal`} color="blue" />
				<Card label="収支" value={`${summary.balance} kcal`} color="gray" />
				<Card
					label="残カロリー"
					value={`${summary.remaining} kcal`}
					color={summary.remaining < 0 ? "red" : "emerald"}
				/>
			</div>

			<div className="bg-white rounded-xl p-4 shadow-sm space-y-2">
				<h2 className="font-semibold text-gray-700">目標まで</h2>
				<p className="text-3xl font-bold text-emerald-600">{targetDays} 日</p>
				<p className="text-sm text-gray-500">
					{goal.currentWeight} kg → {goal.targetWeight} kg
					{diffKg !== 0 && (
						<span className="text-gray-400">（{remainText}）</span>
					)}
				</p>
				<p className="text-sm text-gray-500">期日 {goal.targetDate}</p>

				{targetDays < 0 && (
					<p className="text-sm text-red-500 font-medium">
						⚠ 目標日を過ぎています
					</p>
				)}
			</div>
		</div>
	);
}

function Card({
	label,
	value,
	color,
}: {
	label: string;
	value: string;
	color: string;
}) {
	return (
		<div className="bg-white rounded-xl p-3 shadow-sm">
			<p className={`text-xs text-${color}-500 font-medium`}>{label}</p>
			<p className="text-lg font-bold text-gray-800">{value}</p>
		</div>
	);
}

function Loading() {
	return <div className="p-4 text-gray-400 text-sm">読み込み中…</div>;
}

function ErrorMsg({ msg }: { msg: string }) {
	return <div className="p-4 text-red-500 text-sm">{msg}</div>;
}
