import { auth } from "$lib/api/state.svelte";
import { error, redirect } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";
import { client } from "$lib/api";

export const prerender = true;
export const ssr = false;

const authRoutes = ["/register", "/sign-in", "/discord/callback"];
const profileRoutes = ["/me/edit"];

export const load: LayoutLoad = async ({ url }) => {
  const isUnauthedRoute = authRoutes.some((path) => url.pathname.includes(path));
  if (isUnauthedRoute) {
    return;
  }

  if (auth.isSignedIn === "no") {
    redirect(303, "/sign-in");
  }

  const res = await client.GET("/api/profiles/me");
  if (res.error) {
    if (res.error === "Not Found") {
      const isProfileRoute = profileRoutes.some((path) => url.pathname.includes(path));
      if (isProfileRoute) {
        return;
      }
      redirect(303, "/me/edit");
    }

    error(500, res.error);
  }

  if (!res.data) {
    error(500, "Received empty profile.");
  }

  auth.isSignedIn = "yes";
  auth.profile = res.data;
};
