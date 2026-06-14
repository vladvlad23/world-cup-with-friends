import { redirect } from '@sveltejs/kit';
import { backendFetch } from '$lib/server/backend';

/** @type {import('./$types').PageServerLoad} */
export async function load({ locals, cookies }) {
	if (!locals.user) throw redirect(303, '/login');

	const [matchesRes, usersRes] = await Promise.all([
		backendFetch(cookies, '/api/matches'),
		backendFetch(cookies, '/api/users')
	]);

	return {
		matches: matchesRes.ok ? await matchesRes.json() : [],
		users: usersRes.ok ? await usersRes.json() : []
	};
}
