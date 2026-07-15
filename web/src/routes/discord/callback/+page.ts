import { client } from "$lib/api";
import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { resolve } from "$app/paths";

export const load: PageLoad = async ({ url }) => {
  const authCode = url.searchParams.get("code");
  if (!authCode) {
    const query = new URLSearchParams({ error: "Missing auth code." });
    redirect(303, resolve(`/sign-in?${query.toString()}`));
  }

  const { error } = await client.POST("/api/sessions/discord", { body: { authCode } });
  if (error) {
    const query = new URLSearchParams({ error });
    redirect(303, resolve(`/sign-in?${query.toString()}`));
  }

  redirect(303, resolve("/dashboard"));
};
