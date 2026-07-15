<script lang="ts">
  import { client } from "$lib/api";

  let username = $state("");

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    const res = await client.PUT("/api/profiles/me", { body: { username } });
    if (res.error) {
      console.error(res.error);
    }
  }
</script>

<h1>Edit Me</h1>
<form onsubmit={handleSubmit}>
  <label>
    Username
    <input bind:value={username} type="text" />
  </label>
  <input type="submit" value="Save" />
</form>
