<script lang="ts">
  import { callAPI, getCellTextureUrl, type APICellTexture } from "$lib/api";
  import { type BrushTextureTool, type PaintBucketTextureTool } from "$lib/game/level-editor";
  import { onMount } from "svelte";
  import { addToast } from "./Toaster.svelte";
  import StyledCard from "./StyledCard.svelte";
  import Icon from "@iconify/svelte";

  let {
    tool = $bindable({ type: "texturebrush", texture: null } as
      BrushTextureTool | PaintBucketTextureTool),
  } = $props();
  let textures = $state<APICellTexture[]>([]);

  onMount(async () => {
    const res = await callAPI(fetch, "GET", "/cell-textures");
    if (!res.ok) {
      addToast({
        data: { title: "Texture Load Failed", description: res.error.message, level: "danger" },
      });
      return;
    }

    textures = await res.data.json();
  });
</script>

<StyledCard class="p-4 grid gap-3 pointer-events-auto max-h-[33vh]">
  <ul class="flex flex-wrap gap-4">
    <li class="h-full">
      <button
        data-selected={!tool.texture}
        onclick={() => {
          tool = { ...tool, texture: null };
        }}
        class="grid grid-rows-[1fr_auto] bg-aurora-gray-1300 h-full py-2 px-4 rounded grid gap-1 border-2 border-aurora-gray-1300 data-[selected=true]:bg-aurora-gray-1100 data-[selected=true]:border-white"
      >
        <Icon
          icon="material-symbols:ink-eraser-outline-rounded"
          font-size={64}
          class="w-32 self-center"
        />
        <span class="">Eraser</span>
      </button>
    </li>
    {#each textures as { key, displayName }, i (i)}
      <li>
        <button
          data-selected={key === tool.texture}
          onclick={() => {
            tool = { ...tool, texture: key };
          }}
          class="bg-aurora-gray-1300 py-2 px-4 rounded grid gap-1 border-2 border-aurora-gray-1300 data-[selected=true]:bg-aurora-gray-1100 data-[selected=true]:border-white"
        >
          <img
            alt={`${displayName} cell texture`}
            src={getCellTextureUrl(key).toString()}
            width={128}
            height={64}
            aria-hidden="true"
            class="texture pointer-events-none"
          />
          <span>{displayName}</span>
        </button>
      </li>
    {/each}
  </ul>
</StyledCard>

<style>
  .texture {
    image-rendering: pixelated;
  }
</style>
