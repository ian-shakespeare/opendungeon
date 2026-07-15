import createClient, { type Middleware } from "openapi-fetch";
import type { components, paths } from "$lib/api/schema";

export type APICellTexture = components["schemas"]["models.CellTexture"];
export type APIProfile = components["schemas"]["models.Profile"];

const baseUrl = import.meta.env.DEV ? "http://localhost:8000" : window.location.href;
const authMiddleware: Middleware = {
  async onResponse({ response }) {
    if (response.status === 401) {
      // `window.location` to fully reset the router
      window.location.href = "/sign-in";
    }
  },
};

export const client = createClient<paths>({
  credentials: import.meta.env.DEV ? "include" : "same-origin",
  baseUrl: baseUrl,
});
client.use(authMiddleware);

export function getCellTextureUri(key: string): string {
  const url = new URL(baseUrl);
  url.pathname = "/api/cell-textures/" + key;
  return url.toString();
}
