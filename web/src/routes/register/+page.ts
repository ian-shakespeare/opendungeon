import { client } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
  const res = await client.GET("/api/providers");
  if (res.error) {
    error(500, res.error);
  } else if (!res.data) {
    error(500, "Received empty auth providers list.");
  }

  return { providers: res.data.providers, state: res.data.state };
};
