<script lang="ts">
  import { callAPI } from "$lib/api";
  import StyledButton from "$lib/components/StyledButton.svelte";
  import StyledCard from "$lib/components/StyledCard.svelte";
  import StyledInput from "$lib/components/StyledInput.svelte";
  import StyledMain from "$lib/components/StyledMain.svelte";
  import StyledSeparator from "$lib/components/StyledSeparator.svelte";
  import { addToast } from "$lib/components/Toaster.svelte";
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import type { PageProps } from "./$types";
  import StyledFileUpload from "$lib/components/StyledFileUpload.svelte";

  let { data }: PageProps = $props();

  let username = $state("");
  let file: File | null = $state(null);

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    const body = new FormData();
    body.append("username", username);

    if (file) {
      body.append("avatar", file);
    }

    const res = await callAPI(fetch, "PUT", "/profiles/me", { body });
    if (!res.ok) {
      addToast({ data: { title: "Edit Failed", description: res.error.message, level: "danger" } });
      return;
    }

    if (!data.profile) {
      await goto(resolve("/dashboard"), { invalidate: [resolve("/dashboard")] });
    }
  }
</script>

<svelte:head>
  <title>Edit Profile - OpenDungeon</title>
</svelte:head>

<StyledMain>
  <StyledCard class="max-w-96 w-full px-4 py-6 grid gap-6 md:px-8">
    {#if data.profile}
      <a
        href={resolve("/dashboard")}
        class="text-aurora-gray-700 underline duration-300 hover:text-aurora-gray-500">Exit</a
      >
      <StyledSeparator />
    {/if}
    <h1 class="text-2xl text-center font-semibold text-aurora-gray-600">
      {data.profile ? "Edit" : "Create"} Profile
    </h1>
    <StyledSeparator />
    <form onsubmit={handleSubmit} class="grid gap-6">
      <StyledFileUpload bind:value={file} label="Avatar" icon="material-symbols:person-rounded" />
      <StyledSeparator />
      <StyledInput bind:value={username} type="text" placeholder="Username" />
      <StyledButton label="Save" />
    </form>
  </StyledCard>
</StyledMain>
