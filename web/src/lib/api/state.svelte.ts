import type { APIProfile } from "$lib/api";

export const auth = $state<{ isSignedIn: "unknown" | "yes" | "no"; profile: APIProfile | null }>({
  isSignedIn: "unknown",
  profile: null,
});
