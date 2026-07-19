import { callAPI } from "$lib/api";
import { isRedirect, redirect } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";
import { resolve } from "$app/paths";

export const load: LayoutLoad = async ({ fetch }) => {
  const { ok: isSignedIn } = await callAPI(fetch, "GET", "/profiles/me").catch((error) =>
    !isRedirect(error) ? error : { ok: false },
  );
  if (isSignedIn) {
    redirect(303, resolve("/dashboard"));
  }
};
