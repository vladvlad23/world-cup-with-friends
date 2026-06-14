<script>
	export let title = '';
	export let people = [];
</script>

<div>
	<div class="label">{title}</div>
	<div class="chips">
		{#each people as u (u.id)}
			<span class="chip"><span class="dot" style="background:{u.color}"></span>{u.display_name}</span>
		{/each}
	</div>
</div>

<style>
	.label {
		font-size: 0.8rem;
		color: var(--muted);
		margin-bottom: 0.35rem;
	}
	.chips {
		display: flex;
		flex-wrap: wrap;
		gap: 0.4rem;
	}
	.chip {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		background: var(--panel-2);
		border: 1px solid var(--border);
		border-radius: 999px;
		padding: 0.2rem 0.6rem;
		font-size: 0.82rem;
	}
	.dot {
		width: 10px;
		height: 10px;
		border-radius: 50%;
		display: inline-block;
	}
</style>
