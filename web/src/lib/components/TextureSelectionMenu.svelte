<script lang="ts">
  import { client, getCellTextureUri, type APICellTexture } from "$lib/api";
  import { type BrushTextureTool, type PaintBucketTextureTool } from "$lib/game/level-editor";
  import { onMount } from "svelte";

  let {
    tool = $bindable({ type: "texturebrush", texture: null } as
      BrushTextureTool | PaintBucketTextureTool),
  } = $props();
  let textures = $state<APICellTexture[]>([]);

  onMount(() => {
    client.GET("/api/cell-textures").then(({ data, error }) => {
      if (error) {
        console.error("failed to list cell textures");
        return;
      }

      textures = data ?? [];
    });
  });
</script>

<div
  class="text-white grid bg-aurora-gray-1200 rounded px-4 py-3 grid gap-3 pointer-events-auto max-h-[33vh]"
>
  <ul class="flex flex-wrap gap-4">
    <li>
      <button
        data-selected={tool.texture === null}
        onclick={() => {
          tool = { ...tool, texture: null };
        }}
        class="bg-aurora-gray-1300 py-2 px-4 h-full rounded grid gap-1 data-[selected=true]:bg-aurora-gray-1100 data-[selected=true]:border-2 data-[selected=true]:border-white"
      >
        <span class="w-[128px]">Eraser</span>
      </button>
    </li>
    {#each textures as { key, displayName }, i (i)}
      <li>
        <button
          data-selected={key === tool.texture}
          onclick={() => {
            tool = { ...tool, texture: key };
          }}
          class="bg-aurora-gray-1300 py-2 px-4 rounded grid gap-1 data-[selected=true]:bg-aurora-gray-1100 data-[selected=true]:border-2 data-[selected=true]:border-white"
        >
          <img
            alt={`${displayName} cell texture`}
            src={getCellTextureUri(key)}
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
</div>

<style>
  .texture {
    image-rendering: pixelated;
  }
</style>
