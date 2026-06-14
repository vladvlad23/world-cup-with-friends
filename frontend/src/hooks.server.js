import { getCurrentUser } from '$lib/server/backend';

/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {
	event.locals.user = await getCurrentUser(event.cookies);
	return resolve(event);
}
