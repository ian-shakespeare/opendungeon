import { callAPI, NOT_FOUND, UNAUTHORIZED, type APIProfile } from "$lib/api";
import { error, isRedirect, redirect } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";

const profileRoute = "/me/edit";

export const load: LayoutLoad = async ({ url, fetch }) => {
  // TODO: call an admin endpoint to validate
  const profileRes = await callAPI(fetch, "GET", "/profiles/me").catch(
    (error) =>
      ({
        ok: false,
        error: isRedirect(error)
          ? new Error("Unauthorized", { cause: UNAUTHORIZED })
          : (error as Error),
      }) as const,
  );
  if (!profileRes.ok) {
    if (profileRes.error.cause === UNAUTHORIZED) {
      redirect(303, "/sign-in");
    }

    const isEditingProfile = url.pathname.includes(profileRoute);
    if (profileRes.error.cause === NOT_FOUND) {
      if (isEditingProfile) {
        return { isSignedIn: true };
      }

      redirect(303, "/me/edit");
    }

    error(500, profileRes.error.message);
  }

  return {
    isSignedIn: true,
    profile: (await profileRes.data.json()) as APIProfile,
  };
};
