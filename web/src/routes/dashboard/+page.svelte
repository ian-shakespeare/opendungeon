<script lang="ts">
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { client } from "$lib/api";
  import { auth } from "$lib/api/state.svelte";

  async function handleSignOut() {
    await client.DELETE("/api/sessions");
    auth.isSignedIn = "no";
    auth.profile = null;
    goto(resolve("/sign-in"));
  }
</script>

<h1>Dashboard</h1>
<p>Welcome back, {auth.profile?.username ?? "[username]"}.</p>
<ul>
  <li>
    <a href={resolve("/level-editor")}>Level Editor</a>
  </li>
</ul>
<button onclick={handleSignOut}>Sign Out</button>
