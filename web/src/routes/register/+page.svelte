<script lang="ts">
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { client } from "$lib/api";
  import { auth } from "$lib/api/state.svelte";
  import type { PageProps } from "./$types";

  let { data }: PageProps = $props();

  let email = $state("");
  let password = $state("");

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    const userRes = await client.POST("/api/users", { body: { email } });
    if (userRes.error) {
      // TODO: better surfacing
      console.error("failed to create user");
      return;
    }

    const sessionRes = await client.POST("/api/sessions/email", {
      body: { email: userRes.data!.email, password },
    });
    if (sessionRes.error) {
      // TODO: better surfacing
      console.error(sessionRes.error);
      return;
    }

    // TODO: toast or something
    auth.isSignedIn = "yes";
    await goto(resolve("/dashboard"));
  }
</script>

<svelte:head>
  <title>Register - OpenDungeon</title>
</svelte:head>

<h1>Register</h1>
<ul>
  {#each data.providers as provider, i (i)}
    <li><a rel="external" href={provider.authUrl}>{provider.name}</a></li>
  {/each}
</ul>
<form onsubmit={handleSubmit} class="grid">
  <label>
    Email
    <input bind:value={email} type="email" />
  </label>
  <label>
    Password
    <input bind:value={password} type="password" />
  </label>
  <!-- TODO: actually do password confirmation -->
  <label>
    Confirm Password
    <input type="password" />
  </label>
  <input type="submit" value="Register" class="cursor-pointer" />
</form>
