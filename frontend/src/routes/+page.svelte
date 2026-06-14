<script>
	import { invalidateAll } from '$app/navigation';
	import { buildMonth, monthName, weekdayLabels, dayKey, timeLabel } from '$lib/date';
	import MatchModal from '$lib/components/MatchModal.svelte';

	export let data;

	let selectedMatch = null;

	// Group matches by day for quick lookup in calendar cells.
	$: matchesByDay = data.matches.reduce((acc, m) => {
		const k = dayKey(m.kickoff);
		(acc[k] ||= []).push(m);
		return acc;
	}, {});

	// The tournament runs across June and July 2026.
	const months = [
		{ year: 2026, month: 5 },
		{ year: 2026, month: 6 }
	];

	async function onModalClose(e) {
		selectedMatch = null;
		if (e.detail?.changed) await invalidateAll();
	}
</script>

<h1>Match Calendar</h1>
<p class="lead">Click a match to propose or join a meetup.</p>

{#each months as { year, month } (month)}
	<section class="month">
		<h2>{monthName(month)} {year}</h2>
		<div class="grid head">
			{#each weekdayLabels as wd}<div class="wd">{wd}</div>{/each}
		</div>
		{#each buildMonth(year, month) as week}
			<div class="grid">
				{#each week as cell}
					<div class="cell" class:empty={!cell}>
						{#if cell}
							<div class="daynum">{cell.day}</div>
							{#each matchesByDay[cell.key] || [] as m (m.id)}
								<button class="match" on:click={() => (selectedMatch = m)}>
									<span class="time">{timeLabel(m.kickoff)}</span>
									<span class="tt">{m.home_team} – {m.away_team}</span>
									{#if m.meetup_count > 0}
										<span class="badge">{m.meetup_count}</span>
									{/if}
								</button>
							{/each}
						{/if}
					</div>
				{/each}
			</div>
		{/each}
	</section>
{/each}

{#if selectedMatch}
	<MatchModal
		match={selectedMatch}
		users={data.users}
		currentUser={data.user}
		on:close={onModalClose}
	/>
{/if}

<style>
	h1 {
		margin-bottom: 0.2rem;
	}
	.lead {
		color: var(--muted);
		margin-top: 0;
	}
	.month {
		margin-bottom: 2.5rem;
	}
	.grid {
		display: grid;
		grid-template-columns: repeat(7, 1fr);
		gap: 6px;
	}
	.head {
		margin-bottom: 6px;
	}
	.wd {
		color: var(--muted);
		font-size: 0.75rem;
		text-align: center;
		padding: 0.2rem 0;
	}
	.cell {
		min-height: 96px;
		background: var(--panel);
		border: 1px solid var(--border);
		border-radius: 8px;
		padding: 4px;
		display: flex;
		flex-direction: column;
		gap: 3px;
	}
	.cell.empty {
		background: transparent;
		border: none;
	}
	.daynum {
		font-size: 0.72rem;
		color: var(--muted);
		text-align: right;
		padding-right: 2px;
	}
	.match {
		display: flex;
		align-items: center;
		gap: 4px;
		width: 100%;
		text-align: left;
		background: var(--panel-2);
		color: var(--text);
		border: 1px solid var(--border);
		border-radius: 6px;
		padding: 3px 5px;
		font-size: 0.7rem;
		line-height: 1.15;
	}
	.match:hover {
		border-color: var(--accent);
	}
	.time {
		color: var(--accent-2);
		font-weight: 700;
		flex-shrink: 0;
	}
	.tt {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		flex: 1;
	}
	.badge {
		flex-shrink: 0;
	}
	@media (max-width: 720px) {
		.tt {
			display: none;
		}
	}
</style>
