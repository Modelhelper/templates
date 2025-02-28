import { error, redirect } from "@sveltejs/kit"
import type { RequestHandler } from "./$types"
import { COOKIE_AUTH_NAME } from "$lib/constants"
//  import { currentUser } from "$lib/stores/currentUser";


export const POST: RequestHandler = async ({ locals, cookies }) => {
	cookies.delete(COOKIE_AUTH_NAME, { path: "/" })
	// currentUser.set(null)
	const { error: err } = await locals.supabase.auth.signOut();

	if (err) {
		throw error(500, 'Something went wrong logging you out.')
	}

	throw redirect(303, "/login")
}
export const GET: RequestHandler = async ({ cookies }) => {
	cookies.delete(COOKIE_AUTH_NAME, { path: "/" })
	// currentUser.set(null)
	throw redirect(303, "/")
}