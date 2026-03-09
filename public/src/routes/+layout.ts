import { authStore } from "$features/auth";
import { redirect } from "@sveltejs/kit";

export const ssr = false;
export const prerender = false;

/** @type {import('./$types').LayoutLoad} */
export async function load({ url, fetch, depends }) {
    await authStore.fetchSession();
    const isAuthenticated = authStore.isAuthenticated;
    const authRoutes = ["/sign-in", "/sign-up", "/forgot-password", "/reset-password"];

    if (!isAuthenticated && url.pathname.startsWith("/")) {
        if (!authRoutes.some((route) => url.pathname.startsWith(route))) {
            throw redirect(303, "/sign-in");
        }

        return {
            isAuthenticated: false,
        };
    }

    if (isAuthenticated && authRoutes.some((route) => url.pathname.startsWith(route))) {
        throw redirect(303, "/");
    }

    return { isAuthenticated };
}
