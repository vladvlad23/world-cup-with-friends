import { fail, redirect } from '@sveltejs/kit';
import { backendFetch, TOKEN_COOKIE } from '$lib/server/backend';

/** @type {import('./$types').PageServerLoad} */
export function load({ locals }) {
	if (locals.user) throw redirect(303, '/');
	return {};
}

/** @type {import('./$types').Actions} */
export const actions = {
	default: async ({ request, cookies }) => {
		const form = await request.formData();
		const username = String(form.get('username') || '');
		const password = String(form.get('password') || '');

		const res = await backendFetch(cookies, '/api/auth/login', {
			method: 'POST',
			body: JSON.stringify({ username, password })
		});

		if (!res.ok) {
			return fail(401, { username, error: 'Invalid username or password' });
		}
		const { token } = await res.json();
		cookies.set(TOKEN_COOKIE, token, {
			path: '/',
			httpOnly: true,
			sameSite: 'lax',
			maxAge: 60 * 60 * 24 * 30
		});
		throw redirect(303, '/');
	}
};
