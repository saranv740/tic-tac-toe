<script lang="ts">
</script>

<svelte:head>
	<title>Global Leaderboard - Tic-Tac-Toe</title>
</svelte:head>

<h1 class="mb-8 text-on-surface">Leaderboard</h1>
<div class="leaderboardContainer w-full py-8 md:py-16">
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
			<tr class="text-on-surface">
				<td class="p-6">1</td>
				<td class="p-6">NeoShinku</td>
				<td class="p-6">412 / 12 / 6</td>
				<td class="p-6">12450</td>
			</tr>
			<tr class="text-on-surface">
				<td class="p-6">1</td>
				<td class="p-6">NeoShinku</td>
				<td class="p-6">412 / 12 / 6</td>
				<td class="p-6">12450</td>
			</tr>
			<tr class="text-on-surface">
				<td class="p-6">1</td>
				<td class="p-6">NeoShinku</td>
				<td class="p-6">412 / 12 / 6</td>
				<td class="p-6">12450</td>
			</tr>
			<tr class="text-on-surface">
				<td class="p-6">1</td>
				<td class="p-6">NeoShinku</td>
				<td class="p-6">412 / 12 / 6</td>
				<td class="p-6">12450</td>
			</tr>
			<tr class="text-on-surface">
				<td class="p-6">1</td>
				<td class="p-6">NeoShinku</td>
				<td class="p-6">412 / 12 / 6</td>
				<td class="p-6">12450</td>
			</tr>
		</tbody>
	</table>
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
