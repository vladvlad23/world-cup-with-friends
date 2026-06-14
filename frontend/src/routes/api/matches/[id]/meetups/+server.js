import { backendFetch, passthrough } from '$lib/server/backend';

/** @type {import('./$types').RequestHandler} */
export async function POST({ params, cookies, request }) {
	const body = await request.text();
	return passthrough(
		await backendFetch(cookies, `/api/matches/${params.id}/meetups`, { method: 'POST', body })
	);
}
