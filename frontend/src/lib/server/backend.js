import { env } from '$env/dynamic/private';

const BACKEND_URL = env.BACKEND_URL || 'http://localhost:8080';
export const TOKEN_COOKIE = 'wc_token';

/**
 * Call the Go backend, attaching the JWT from the request cookie.
 * `cookies` is the SvelteKit cookies object from the event.
 */
export async function backendFetch(cookies, path, init = {}) {
	const token = cookies.get(TOKEN_COOKIE);
	const headers = { ...(init.headers || {}) };
	if (token) headers['Authorization'] = `Bearer ${token}`;
	if (init.body && !headers['Content-Type']) headers['Content-Type'] = 'application/json';
	return fetch(`${BACKEND_URL}${path}`, { ...init, headers });
}

/** Wrap a backend Response for return from a SvelteKit endpoint, preserving
 * status and content type without assuming the body is JSON. */
export function passthrough(res) {
	return new Response(res.body, {
		status: res.status,
		headers: { 'content-type': res.headers.get('content-type') || 'application/json' }
	});
}

/** Returns the logged-in user (via /api/me) or null if not authenticated. */
export async function getCurrentUser(cookies) {
	const token = cookies.get(TOKEN_COOKIE);
	if (!token) return null;
	const res = await backendFetch(cookies, '/api/me');
	if (!res.ok) return null;
	return res.json();
}
