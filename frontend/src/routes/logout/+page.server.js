import { redirect } from '@sveltejs/kit';
import { TOKEN_COOKIE } from '$lib/server/backend';

/** @type {import('./$types').Actions} */
export const actions = {
	default: async ({ cookies }) => {
		cookies.delete(TOKEN_COOKIE, { path: '/' });
		throw redirect(303, '/login');
	}
};
