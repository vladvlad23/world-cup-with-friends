<script>
	import { createEventDispatcher } from 'svelte';
	import { changePassword } from '$lib/api';

	const dispatch = createEventDispatcher();

	let current = '';
	let next = '';
	let confirm = '';
	let saving = false;
	let error = '';
	let done = false;

	async function submit() {
		error = '';
		if (next.length < 8) {
			error = 'New password must be at least 8 characters.';
			return;
		}
		if (next !== confirm) {
			error = 'New passwords do not match.';
			return;
		}
		saving = true;
		try {
			await changePassword(current, next);
			done = true;
		} catch (e) {
			error = e.message;
		} finally {
			saving = false;
		}
	}
</script>

<div class="modal-backdrop" on:click|self={() => dispatch('close')} role="presentation">
	<div class="modal">
		<h2>Change password</h2>

		{#if done}
			<p class="ok">✓ Your password has been updated.</p>
			<div class="actions">
				<button type="button" on:click={() => dispatch('close')}>Done</button>
			</div>
		{:else}
			<div class="field">
				<label for="cur">Current password</label>
				<input id="cur" type="password" bind:value={current} autocomplete="current-password" />
			</div>
			<div class="field">
				<label for="new">New password</label>
				<input id="new" type="password" bind:value={next} autocomplete="new-password" />
			</div>
			<div class="field">
				<label for="conf">Confirm new password</label>
				<input id="conf" type="password" bind:value={confirm} autocomplete="new-password" />
			</div>

			{#if error}<p class="err">{error}</p>{/if}

			<div class="actions">
				<button class="ghost" type="button" on:click={() => dispatch('close')}>Cancel</button>
				<button type="button" on:click={submit} disabled={saving}>
					{saving ? 'Saving…' : 'Update password'}
				</button>
			</div>
		{/if}
	</div>
</div>

<style>
	.actions {
		display: flex;
		justify-content: space-between;
		margin-top: 1.2rem;
	}
	.err {
		color: var(--danger);
		font-size: 0.85rem;
	}
	.ok {
		color: var(--accent);
	}
</style>
