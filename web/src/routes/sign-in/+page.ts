import { client } from "$lib/api";
import { error, redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { auth } from "$lib/api/state.svelte";
import { resolve } from "$app/paths";

export const load: PageLoad = async () => {
  if (auth.isSignedIn === "yes") {
    redirect(303, resolve("/dashboard"));
  }

  const res = await client.GET("/api/providers");
  if (res.error) {
    error(500, res.error);
  } else if (!res.data) {
    error(500, "Received empty auth providers list.");
  }

  return { providers: res.data };
};
