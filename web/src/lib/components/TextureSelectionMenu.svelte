<script lang="ts">
  import { getCellTextureUrl, listCellTextures, type APICellTexture } from "$lib/api.svelte";
  import { type BrushTextureTool, type PaintBucketTextureTool } from "$lib/game/level-editor";
  import { onMount } from "svelte";
  import { addToast } from "./Toaster.svelte";
  import StyledCard from "./StyledCard.svelte";
  import StyledButton from "./StyledButton.svelte";

  let {
    tool = $bindable({ type: "texturebrush", texture: null } as
      BrushTextureTool | PaintBucketTextureTool),
  } = $props();
  let textures = $state<APICellTexture[]>([]);

  onMount(async () => {
    const res = await listCellTextures();
    if (!res.ok) {
      addToast({
        data: { title: "Texture Load Failed", description: res.error.message, level: "danger" },
      });
      return;
    }

    textures = res.textures;
  });
</script>

<StyledCard class="p-4 grid gap-3 pointer-events-auto max-h-[33vh]">
  <ul class="flex flex-wrap gap-4">
    <li>
      <StyledButton
        mode={tool.texture ? "secondary" : "primary"}
        label="Eraser"
        class="px-4 h-full w-[128px]"
      />
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
