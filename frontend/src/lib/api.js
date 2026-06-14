// Client-side helpers that talk to the SvelteKit proxy endpoints (which attach
// the auth cookie and forward to the Go backend).

async function asJson(res) {
	const data = await res.json().catch(() => ({}));
	if (!res.ok) throw new Error(data.error || `Request failed (${res.status})`);
	return data;
}

export async function getMatch(id) {
	return asJson(await fetch(`/api/matches/${id}`));
}

export async function createMeetup(matchId, payload) {
	return asJson(
		await fetch(`/api/matches/${matchId}/meetups`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(payload)
		})
	);
}

export async function getMeetup(id) {
	return asJson(await fetch(`/api/meetups/${id}`));
}

export async function joinMeetup(id) {
	return asJson(await fetch(`/api/meetups/${id}/join`, { method: 'POST' }));
}

export async function changePassword(currentPassword, newPassword) {
	return asJson(
		await fetch('/api/me/password', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ current_password: currentPassword, new_password: newPassword })
		})
	);
}
