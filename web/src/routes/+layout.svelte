<script lang="ts">
  import "./layout.css";
  import favicon from "$lib/assets/favicon.svg";
  import Toaster from "$lib/components/Toaster.svelte";
  import type { LayoutProps } from "./$types";

  let { children }: LayoutProps = $props();
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
  <title>OpenDungeon</title>
</svelte:head>

<div class="h-full bg-linear-to-br from-aurora-magenta-900 to-aurora-magenta-1000">
  {@render children()}
</div>

<Toaster />
