import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ request, cookies }) => {
	try {
		const { access_token } = await request.json();
		if (access_token) {
			cookies.set('sb-access-token', access_token, {
				path: '/',
				httpOnly: true,
				sameSite: 'lax',
				secure: process.env.NODE_ENV === 'production',
				maxAge: 60 * 60 * 24 * 7 // 1 week
			});
		} else {
			cookies.delete('sb-access-token', { path: '/' });
		}
		return json({ success: true });
	} catch (e: any) {
		return json({ success: false, error: e.message }, { status: 400 });
	}
};
