import { callAPI, type APILevelData } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
  const levelId = url.searchParams.get("levelId");
  if (!levelId) {
    return { level: undefined };
  }

  const res = await callAPI(fetch, "GET", "/levels/" + levelId);
  if (!res.ok) {
    error(500, res.error.message);
  }

  return {
    level: (await res.data.json()) as APILevelData,
  };
};
