<script lang="ts">
	import { onMount } from 'svelte';
	import { client } from '$lib/services/nakama';
	import { sessionStore } from '$lib/stores/session.svelte';
	import type { UserStats } from '$lib/types';

	interface LeaderboardRow {
		rank: number;
		username: string;
		wins: number;
		losses: number;
		draws: number;
		points: number; // = wins (the leaderboard score)
	}

	let rows = $state<LeaderboardRow[]>([]);
	let isLoading = $state(true);
	let error = $state('');

	onMount(async () => {
		if (!sessionStore.session) return;

		try {
			// 1. Fetch top 50 leaderboard records (score = wins, operator = incr)
			const lb = await client.listLeaderboardRecords(
				sessionStore.session,
				'tic_tac_toe_global',
				[], // ownerIds — empty = global list
				50 // limit: top 50 players
			);

			const records = lb.records ?? [];
			if (records.length === 0) {
				isLoading = false;
				return;
			}

			// 2. Batch-read W/L/D stats from Nakama Storage for all 50 players in one call
			//    Backend writes to collection "stats", key "tic-tac-toe", per-user
			const storageResult = await client.readStorageObjects(sessionStore.session, {
				object_ids: records.map((r) => ({
					collection: 'stats',
					key: 'tic-tac-toe',
					user_id: r.owner_id!
				}))
			});

			// Index storage results by user_id for O(1) lookup
			const statsMap = new Map<string, UserStats>();
			for (const obj of storageResult.objects ?? []) {
				if (obj.user_id && obj.value) {
					statsMap.set(obj.user_id, obj.value as UserStats);
				}
			}

			// 3. Merge leaderboard rows with W/L/D stats
			rows = records.map((r) => {
				const stats = statsMap.get(r.owner_id!) ?? { wins: 0, losses: 0, draws: 0 };
				return {
					rank: r.rank || 0,
					username: r.username ?? 'Unknown',
					wins: stats.wins,
					losses: stats.losses,
					draws: stats.draws,
					points: r.score || 0 // score = total wins
				};
			});
		} catch (err) {
			console.error('Failed to load leaderboard:', err);
			error = 'Could not load leaderboard. Is the server running?';
		} finally {
			isLoading = false;
		}
	});
</script>

<svelte:head>
	<title>Global Leaderboard - Tic-Tac-Toe</title>
</svelte:head>

<h1 class="mb-8 text-on-surface">Leaderboard</h1>
<div class="leaderboardContainer w-full py-8 md:py-16">
	{#if isLoading}
		<p class="text-center text-on-surface-variant">Loading…</p>
	{:else if error}
		<p class="text-center text-red-400">{error}</p>
	{:else if rows.length === 0}
		<p class="text-center text-on-surface-variant">
			No players on the leaderboard yet. Play a game!
		</p>
	{:else}
		<table>
			<thead>
				<tr class="bg-surface-container-high text-xs font-bold text-on-surface-variant uppercase">
					<th scope="col" class="p-6">Rank</th>
					<th scope="col" class="p-6">Player name</th>
					<th scope="col" class="p-6">W / L / D</th>
					<th scope="col" class="p-6">Points</th>
				</tr>
			</thead>
			<tbody>
				{#each rows as row (row.username)}
					<tr class="text-on-surface">
						<td class="p-6">{row.rank}</td>
						<td class="p-6">{row.username}</td>
						<td class="p-6">{row.wins} / {row.losses} / {row.draws}</td>
						<td class="p-6">{row.points}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>

<style>
	h1 {
		font-size: 2.5rem;
		font-weight: bold;
	}

	.leaderboardContainer {
		max-width: 100%;
		overflow: auto;
	}

	table {
		table-layout: fixed;
		width: 100%;
		border-collapse: collapse;
	}

	th {
		white-space: nowrap;
	}

	tbody tr:nth-child(odd) {
		background-color: var(--color-surface-container-low);
	}

	tbody tr:nth-child(even) {
		background-color: #24283b;
	}

	tr :nth-child(1) {
		width: 15%;
		text-align: right;
	}

	tr :nth-child(2) {
		width: 35%;
		text-align: center;
	}

	tr :nth-child(3) {
		width: 25%;
		text-align: left;

		@media screen and (width > 40rem) {
			text-align: right;
		}
	}

	tr :nth-child(4) {
		width: 25%;
		text-align: right;
	}
</style>
