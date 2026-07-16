<script lang="ts">
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { register, type APIAuthProvider } from "$lib/api.svelte";
  import StyledAnchor from "$lib/components/StyledAnchor.svelte";
  import StyledButton from "$lib/components/StyledButton.svelte";
  import StyledCard from "$lib/components/StyledCard.svelte";
  import StyledInput from "$lib/components/StyledInput.svelte";
  import discordLogo from "$lib/assets/discord-logo.svg";
  import type { PageProps } from "./$types";
  import { addToast } from "$lib/components/Toaster.svelte";
  import StyledMain from "$lib/components/StyledMain.svelte";
  import StyledSeparator from "$lib/components/StyledSeparator.svelte";

  let { data }: PageProps = $props();

  let email = $state("");
  let password = $state("");

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    const res = await register(email, password);
    if (res.ok) {
      await goto(resolve("/dashboard"));
      return;
    }

    addToast({
      data: { title: "Registration Failed", description: res.error.message, level: "danger" },
    });
  }

  function isDiscord(provider: APIAuthProvider): boolean {
    return provider.name.toLowerCase() === "discord";
  }
</script>

<svelte:head>
  <title>Register - OpenDungeon</title>
</svelte:head>

<StyledMain>
  <StyledCard class="px-4 py-6 max-w-96 w-full">
    <ul>
      {#each data.providers as provider, i (i)}
        <li>
          <StyledAnchor
            mode="none"
            rel="external"
            href={provider.authUrl}
            label={"Sign in with " + provider.name}
            icon={!isDiscord(provider) ? undefined : discordLogo}
            class={!isDiscord(provider)
              ? ""
              : "bg-discord-blurple duration-300 hover:bg-white hover:text-discord-blurple"}
          />
        </li>
      {/each}
    </ul>
    <StyledSeparator class="my-6" />
    <form onsubmit={handleSubmit} class="grid gap-4 mb-2">
      <div class="grid gap-2">
        <StyledInput bind:value={email} type="email" placeholder="Email" />
        <StyledInput bind:value={password} type="password" placeholder="Password" />
        <!-- TODO: confirm password -->
        <StyledInput type="password" placeholder="Confirm Password" />
      </div>
      <StyledButton label="Register" />
    </form>
    <p class="text-aurora-gray-700 text-center">
      Already have an account?
      <a href={resolve("/sign-in")} class="text-aurora-magenta-300 underline">Sign in here.</a>
    </p>
  </StyledCard>
</StyledMain>
