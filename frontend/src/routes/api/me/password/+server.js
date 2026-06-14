import { backendFetch, passthrough } from '$lib/server/backend';

/** @type {import('./$types').RequestHandler} */
export async function POST({ cookies, request }) {
	const body = await request.text();
	return passthrough(await backendFetch(cookies, '/api/me/password', { method: 'POST', body }));
}
