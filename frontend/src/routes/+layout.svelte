<script>
	import '../app.css';
	import ChangePasswordModal from '$lib/components/ChangePasswordModal.svelte';
	export let data;

	let showChangePassword = false;
</script>

<div class="shell">
	<header>
		<a class="brand" href="/">⚽ World Cup 2026 <span>Meetups</span></a>
		{#if data.user}
			<div class="user">
				<span class="dot" style="background:{data.user.color}"></span>
				<span>{data.user.display_name}</span>
				<button class="ghost" type="button" on:click={() => (showChangePassword = true)}>
					Change password
				</button>
				<form method="POST" action="/logout">
					<button class="ghost" type="submit">Log out</button>
				</form>
			</div>
		{/if}
	</header>
	<main>
		<slot />
	</main>
</div>

{#if showChangePassword}
	<ChangePasswordModal on:close={() => (showChangePassword = false)} />
{/if}

<style>
	.shell {
		max-width: 1100px;
		margin: 0 auto;
		padding: 0 1rem 3rem;
	}
	header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1.2rem 0;
		border-bottom: 1px solid var(--border);
		margin-bottom: 1.5rem;
	}
	.brand {
		font-size: 1.25rem;
		font-weight: 800;
		text-decoration: none;
		color: var(--text);
	}
	.brand span {
		color: var(--accent);
	}
	.user {
		display: flex;
		align-items: center;
		gap: 0.6rem;
		font-size: 0.9rem;
	}
	.user form {
		margin: 0;
	}
	.dot {
		width: 14px;
		height: 14px;
		border-radius: 50%;
		display: inline-block;
	}
</style>
