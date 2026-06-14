<script>
	import { enhance } from '$app/forms';
	export let form;
</script>

<div class="login-wrap">
	<div class="card">
		<h1>⚽ World Cup 2026</h1>
		<p class="sub">Sign in to plan meetups around the matches.</p>

		<form method="POST" use:enhance>
			<div class="field">
				<label for="username">Username</label>
				<input id="username" name="username" autocomplete="username" value={form?.username ?? ''} required />
			</div>
			<div class="field">
				<label for="password">Password</label>
				<input id="password" name="password" type="password" autocomplete="current-password" required />
			</div>
			{#if form?.error}
				<p class="err">{form.error}</p>
			{/if}
			<button type="submit" style="width:100%">Sign in</button>
		</form>

		<p class="hint">Demo users: <code>alice</code>, <code>bob</code>, <code>carol</code>, <code>dave</code>, <code>erin</code> — password <code>password</code></p>
	</div>
</div>

<style>
	.login-wrap {
		min-height: 60vh;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.card {
		background: var(--panel);
		border: 1px solid var(--border);
		border-radius: 16px;
		padding: 2rem;
		width: 100%;
		max-width: 380px;
	}
	h1 {
		margin: 0 0 0.25rem;
	}
	.sub {
		color: var(--muted);
		margin-top: 0;
	}
	.err {
		color: var(--danger);
		font-size: 0.85rem;
	}
	.hint {
		margin-top: 1.2rem;
		font-size: 0.78rem;
		color: var(--muted);
	}
	code {
		background: var(--bg);
		padding: 0.1rem 0.3rem;
		border-radius: 4px;
	}
</style>
