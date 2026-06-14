<script>
	import { createEventDispatcher, onMount } from 'svelte';
	import { getMeetup, joinMeetup } from '$lib/api';
	import PeopleRow from './PeopleRow.svelte';

	export let meetupId;
	export let currentUser;

	const dispatch = createEventDispatcher();

	let meetup = null;
	let loading = true;
	let joining = false;
	let error = '';

	onMount(load);

	async function load() {
		loading = true;
		try {
			meetup = await getMeetup(meetupId);
		} catch (e) {
			error = e.message;
		} finally {
			loading = false;
		}
	}

	async function join() {
		joining = true;
		error = '';
		try {
			meetup = await joinMeetup(meetupId);
			dispatch('joined');
		} catch (e) {
			error = e.message;
		} finally {
			joining = false;
		}
	}

	$: isMember = meetup?.members?.some((u) => u.id === currentUser.id);
</script>

{#if loading}
	<p class="muted">Loading…</p>
{:else if !meetup}
	<p class="err">{error || 'Could not load meetup.'}</p>
{:else}
	<h2>
		{meetup.location_name}
		{#if meetup.location_is_bar}<span class="badge">BAR</span>{/if}
	</h2>
	<p class="muted">Hosted by {meetup.created_by.display_name}</p>

	{#if meetup.location_is_bar}
		<div class="warning">⚠️ The host is responsible for the bar reservation.</div>
	{/if}

	{#if meetup.location_url}
		<p><a href={meetup.location_url} target="_blank" rel="noreferrer">📍 Open in Maps</a></p>
	{/if}

	{#if meetup.note}
		<div class="note">{meetup.note}</div>
	{/if}

	<div class="people">
		<PeopleRow title={`Going (${meetup.members.length})`} people={meetup.members} />
		{#if meetup.invites.length}
			<PeopleRow title="Invited" people={meetup.invites} />
		{/if}
		{#if meetup.mentions.length}
			<PeopleRow title="Mentioned" people={meetup.mentions} />
		{/if}
	</div>

	{#if error}<p class="err">{error}</p>{/if}

	<div class="actions">
		<button class="ghost" type="button" on:click={() => dispatch('back')}>Back</button>
		{#if isMember}
			<button disabled>✓ You're going</button>
		{:else}
			<button type="button" on:click={join} disabled={joining}>
				{joining ? 'Joining…' : 'Join meetup'}
			</button>
		{/if}
	</div>
{/if}

<style>
	.muted {
		color: var(--muted);
	}
	.note {
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		padding: 0.7rem 0.8rem;
		white-space: pre-wrap;
		margin: 0.8rem 0;
	}
	.people {
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
		margin: 1rem 0;
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
