<script>
	import { createEventDispatcher } from 'svelte';
	import MentionInput from './MentionInput.svelte';
	import { createMeetup } from '$lib/api';
	import { parseMentions } from '$lib/mentions';

	export let matchId;
	export let users = [];
	export let currentUser;

	const dispatch = createEventDispatcher();

	let locationName = '';
	let locationUrl = '';
	let isBar = false;
	let note = '';
	let invited = new Set();
	let submitting = false;
	let error = '';

	// You can't invite yourself — you're the host.
	$: invitable = users.filter((u) => u.id !== currentUser.id);

	function toggleInvite(id) {
		if (invited.has(id)) invited.delete(id);
		else invited.add(id);
		invited = invited; // trigger reactivity
	}

	async function submit() {
		if (!locationName.trim()) {
			error = 'Please name the location.';
			return;
		}
		submitting = true;
		error = '';
		try {
			const meetup = await createMeetup(matchId, {
				location_name: locationName.trim(),
				location_url: locationUrl.trim(),
				location_is_bar: isBar,
				note,
				invite_user_ids: [...invited],
				mention_user_ids: parseMentions(note, users)
			});
			dispatch('created', meetup);
		} catch (e) {
			error = e.message;
		} finally {
			submitting = false;
		}
	}
</script>

<h2>Propose a meetup</h2>

<div class="field">
	<label for="loc">Location name</label>
	<input id="loc" bind:value={locationName} placeholder="e.g. The Anchor Pub" />
	<label class="inline">
		<input type="checkbox" bind:checked={isBar} class="cb" /> This is a bar / restaurant
	</label>
	{#if isBar}
		<div class="warning">
			⚠️ Since this is a bar/restaurant, <strong>you (the host) are responsible for making the
			reservation.</strong>
		</div>
	{/if}
</div>

<div class="field">
	<label for="url">Google Maps link (optional)</label>
	<input id="url" bind:value={locationUrl} placeholder="https://maps.google.com/…" />
</div>

<div class="field">
	<label>Invite people</label>
	<div class="invite-grid">
		{#each invitable as u (u.id)}
			<button
				type="button"
				class:active={invited.has(u.id)}
				class="chip"
				on:click={() => toggleInvite(u.id)}
			>
				<span class="dot" style="background:{u.color}"></span>{u.display_name}
			</button>
		{/each}
	</div>
</div>

<div class="field">
	<label>Note &amp; mentions</label>
	<MentionInput bind:value={note} {users} />
</div>

{#if error}<p class="err">{error}</p>{/if}

<div class="actions">
	<button class="ghost" type="button" on:click={() => dispatch('back')}>Back</button>
	<button type="button" on:click={submit} disabled={submitting}>
		{submitting ? 'Creating…' : 'Create meetup'}
	</button>
</div>

<style>
	.inline {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin-top: 0.6rem;
		color: var(--text);
	}
	.cb {
		width: auto;
	}
	.invite-grid {
		display: flex;
		flex-wrap: wrap;
		gap: 0.4rem;
	}
	.chip {
		background: var(--panel-2);
		color: var(--text);
		border: 1px solid var(--border);
		display: flex;
		align-items: center;
		gap: 0.4rem;
		font-weight: 500;
	}
	.chip.active {
		background: var(--accent);
		color: #042;
		border-color: var(--accent);
	}
	.dot {
		width: 10px;
		height: 10px;
		border-radius: 50%;
		display: inline-block;
	}
	.actions {
		display: flex;
		justify-content: space-between;
		margin-top: 1.2rem;
	}
	.err {
		color: var(--danger);
		font-size: 0.85rem;
	}
</style>
