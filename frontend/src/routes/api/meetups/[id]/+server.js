import { backendFetch, passthrough } from '$lib/server/backend';

/** @type {import('./$types').RequestHandler} */
export async function GET({ params, cookies }) {
	return passthrough(await backendFetch(cookies, `/api/meetups/${params.id}`));
}
