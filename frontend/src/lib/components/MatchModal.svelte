<script>
	import { createEventDispatcher, onMount } from 'svelte';
	import { getMatch } from '$lib/api';
	import { timeLabel, prettyDate } from '$lib/date';
	import { dayKey } from '$lib/date';
	import ProposeMeetupForm from './ProposeMeetupForm.svelte';
	import MeetupList from './MeetupList.svelte';
	import MeetupDetail from './MeetupDetail.svelte';

	export let match; // summary from the calendar (id, teams, kickoff, stadium…)
	export let users = [];
	export let currentUser;

	const dispatch = createEventDispatcher();

	let view = 'menu'; // menu | propose | list | detail
	let meetups = [];
	let selectedMeetupId = null;
	let loading = true;
	let dirty = false; // whether anything changed (to refresh calendar)

	onMount(refresh);

	async function refresh() {
		loading = true;
		const data = await getMatch(match.id);
		meetups = data.meetups || [];
		loading = false;
	}

	function close() {
		dispatch('close', { changed: dirty });
	}

	function onCreated() {
		dirty = true;
		refresh();
		view = 'list';
	}

	function openDetail(id) {
		selectedMeetupId = id;
		view = 'detail';
	}
</script>

<div class="modal-backdrop" on:click|self={close} role="presentation">
	<div class="modal">
		<div class="match-head">
			<div class="teams">
				{#if match.home_flag}<img src={match.home_flag} alt="" />{/if}
				<strong>{match.home_team}</strong>
				<span class="vs">vs</span>
				<strong>{match.away_team}</strong>
				{#if match.away_flag}<img src={match.away_flag} alt="" />{/if}
			</div>
			<div class="meta">
				{prettyDate(dayKey(match.kickoff))} · {timeLabel(match.kickoff)}
				{#if match.stadium_name}· {match.stadium_name}, {match.city}{/if}
			</div>
		</div>

		{#if view === 'menu'}
			<div class="menu">
				<button on:click={() => (view = 'propose')}>➕ Propose a meetup</button>
				<button class="secondary" on:click={() => (view = 'list')}>
					📋 See existing meetups{#if !loading}&nbsp;({meetups.length}){/if}
				</button>
				<button class="ghost" on:click={close}>Cancel</button>
			</div>
		{:else if view === 'propose'}
			<ProposeMeetupForm
				matchId={match.id}
				{users}
				{currentUser}
				on:created={onCreated}
				on:back={() => (view = 'menu')}
			/>
		{:else if view === 'list'}
			<MeetupList {meetups} on:select={(e) => openDetail(e.detail)} on:back={() => (view = 'menu')} />
		{:else if view === 'detail'}
			<MeetupDetail
				meetupId={selectedMeetupId}
				{currentUser}
				on:joined={() => {
					dirty = true;
					refresh();
				}}
				on:back={() => (view = 'list')}
			/>
		{/if}
	</div>
</div>

<style>
	.match-head {
		border-bottom: 1px solid var(--border);
		padding-bottom: 1rem;
		margin-bottom: 1rem;
	}
	.teams {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 1.15rem;
	}
	.teams img {
		width: 26px;
		height: auto;
		border-radius: 3px;
	}
	.vs {
		color: var(--muted);
		font-weight: 400;
	}
	.meta {
		color: var(--muted);
		font-size: 0.85rem;
		margin-top: 0.4rem;
	}
	.menu {
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
	}
	.menu button {
		width: 100%;
		padding: 0.8rem;
		font-size: 1rem;
	}
</style>
