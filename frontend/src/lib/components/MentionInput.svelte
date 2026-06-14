<script>
	export let value = '';
	export let users = [];

	let textarea;
	let showMenu = false;
	let query = '';
	let menuStart = -1;

	$: suggestions = showMenu
		? users.filter((u) => u.username.toLowerCase().startsWith(query.toLowerCase())).slice(0, 6)
		: [];

	function onInput(e) {
		value = e.target.value;
		const caret = e.target.selectionStart;
		const upto = value.slice(0, caret);
		const match = upto.match(/@([a-zA-Z0-9_]*)$/);
		if (match) {
			showMenu = true;
			query = match[1];
			menuStart = caret - match[0].length;
		} else {
			showMenu = false;
		}
	}

	function pick(user) {
		const before = value.slice(0, menuStart);
		const caret = textarea.selectionStart;
		const after = value.slice(caret);
		value = `${before}@${user.username} ${after}`;
		showMenu = false;
		textarea.focus();
	}
</script>

<div class="mention-wrap">
	<textarea
		bind:this={textarea}
		rows="3"
		placeholder="Add a note. Type @ to mention someone…"
		on:input={onInput}
		on:blur={() => setTimeout(() => (showMenu = false), 150)}>{value}</textarea
	>
	{#if showMenu && suggestions.length}
		<ul class="menu">
			{#each suggestions as u (u.id)}
				<li>
					<button type="button" class="ghost" on:click={() => pick(u)}>
						<span class="dot" style="background:{u.color}"></span>@{u.username}
						<span class="muted">{u.display_name}</span>
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	.mention-wrap {
		position: relative;
	}
	.menu {
		position: absolute;
		left: 0;
		right: 0;
		top: 100%;
		margin: 0.2rem 0 0;
		padding: 0.25rem;
		list-style: none;
		background: var(--panel-2);
		border: 1px solid var(--border);
		border-radius: var(--radius);
		z-index: 10;
	}
	.menu button {
		width: 100%;
		text-align: left;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		border: none;
		padding: 0.4rem 0.5rem;
	}
	.muted {
		color: var(--muted);
		font-size: 0.8rem;
	}
	.dot {
		width: 10px;
		height: 10px;
		border-radius: 50%;
		display: inline-block;
	}
</style>
