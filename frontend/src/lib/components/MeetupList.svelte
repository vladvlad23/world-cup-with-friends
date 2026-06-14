<script>
	import { createEventDispatcher } from 'svelte';
	export let meetups = [];
	const dispatch = createEventDispatcher();
</script>

<h2>Existing meetups</h2>

{#if meetups.length === 0}
	<p class="muted">No meetups yet for this match. Be the first to propose one!</p>
{:else}
	<ul class="list">
		{#each meetups as m (m.id)}
			<li>
				<button type="button" class="row" on:click={() => dispatch('select', m.id)}>
					<div>
						<strong>{m.location_name}</strong>
						{#if m.location_is_bar}<span class="badge">BAR</span>{/if}
						<div class="muted">
							by {m.created_by.display_name} · {m.member_count} going
						</div>
					</div>
					<span class="arrow">›</span>
				</button>
			</li>
		{/each}
	</ul>
{/if}

<div class="actions">
	<button class="ghost" type="button" on:click={() => dispatch('back')}>Back</button>
</div>

<style>
	.list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}
	.row {
		width: 100%;
		text-align: left;
		display: flex;
		align-items: center;
		justify-content: space-between;
		background: var(--panel-2);
		border: 1px solid var(--border);
		color: var(--text);
		font-weight: 500;
		padding: 0.7rem 0.9rem;
	}
	.muted {
		color: var(--muted);
		font-size: 0.8rem;
		margin-top: 0.2rem;
	}
	.arrow {
		font-size: 1.4rem;
		color: var(--muted);
	}
	.actions {
		margin-top: 1.2rem;
	}
</style>
