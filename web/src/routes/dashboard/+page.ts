import { listLevels, UNAUTHORIZED } from "$lib/api.svelte";
import { error, redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { resolve } from "$app/paths";

export const load: PageLoad = async () => {
  const res = await listLevels();
  if (!res.ok) {
    if (res.error.cause === UNAUTHORIZED) {
      redirect(303, resolve("/sign-in"));
    }

    error(500, "Failed to get levels.");
  }

  return { levels: res.levels };
};
