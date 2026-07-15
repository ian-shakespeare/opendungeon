<script lang="ts">
  import { getCellTextureUri } from "$lib/api";
  import GameWindow from "$lib/components/GameWindow.svelte";
  import LevelEditorToolMenu from "$lib/components/LevelEditorToolMenu.svelte";
  import LevelEditor from "$lib/game/level-editor";

  let editor = new LevelEditor();
  let tool = $state(editor.tool);
  let viewMode = $state(editor.viewMode);

  $effect(() => {
    if (tool.type === "texturebrush" || tool.type === "texturepaintbucket") {
      if (tool.texture && !editor.hasTexture(tool.texture)) {
        editor.pause();
        editor
          .loadTexture(tool.texture, getCellTextureUri(tool.texture))
          .then(() => editor.unpause())
          .catch(() => alert("Failed to load texture. Please reload the page."));
      }
    }
    editor.tool = tool;
  });

  $effect(() => {
    editor.viewMode = viewMode;
  });
</script>

<svelte:head>
  <title>Level Editor - OpenDungeon</title>
</svelte:head>

<main class="grid justify-start h-full relative">
  <LevelEditorToolMenu bind:tool bind:viewMode />
  <GameWindow game={editor} />
</main>
