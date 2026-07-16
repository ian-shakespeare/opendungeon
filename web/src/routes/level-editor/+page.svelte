<script lang="ts">
  import { goto } from "$app/navigation";
  import { getCellTextureUrl, createLevel } from "$lib/api.svelte";
  import GameWindow from "$lib/components/GameWindow.svelte";
  import LevelEditorToolMenu from "$lib/components/LevelEditorToolMenu.svelte";
  import StyledButton from "$lib/components/StyledButton.svelte";
  import StyledInput from "$lib/components/StyledInput.svelte";
  import { addToast } from "$lib/components/Toaster.svelte";
  import LevelEditor from "$lib/game/level-editor";
  import { resolve } from "$app/paths";

  let editor = new LevelEditor();
  let tool = $state(editor.tool);
  let viewMode = $state(editor.viewMode);
  let levelName = $state("");

  $effect(() => {
    if (tool.type === "texturebrush" || tool.type === "texturepaintbucket") {
      if (tool.texture && !editor.hasTexture(tool.texture)) {
        editor.pause();
        editor
          .loadTexture(tool.texture, getCellTextureUrl(tool.texture).toString())
          .then(() => editor.unpause())
          .catch(() => alert("Failed to load texture. Please reload the page."));
      }
    }
    editor.tool = tool;
  });

  $effect(() => {
    editor.viewMode = viewMode;
  });

  async function handleSaveLevel() {
    if (levelName.length < 3) {
      addToast({
        data: {
          title: "Invalid Name",
          description: "Name must be at least 3 characters.",
          level: "warn",
        },
      });
      return;
    }

    const res = await createLevel(levelName, editor.grid);
    if (!res.ok) {
      addToast({ data: { title: "Save Failed", description: res.error.message, level: "danger" } });
      return;
    }

    addToast({
      data: { title: "Success", description: "Level was successfully saved.", level: "success" },
    });
  }
</script>

<svelte:head>
  <title>Level Editor - OpenDungeon</title>
</svelte:head>

<main class="grid h-full relative grid-rows-[auto_1fr]">
  <div class="h-16 z-10 flex items-center px-6 gap-2 pointer-events-none">
    <StyledButton
      mode="secondary"
      label="Exit"
      onclick={() => goto(resolve("/dashboard"))}
      class="px-6 pointer-events-auto"
    />
    <StyledInput bind:value={levelName} placeholder="Name" class="pointer-events-auto" />
    <StyledButton
      mode="outline"
      label="Save"
      onclick={handleSaveLevel}
      class="w-20 pointer-events-auto"
    />
  </div>
  <LevelEditorToolMenu bind:tool bind:viewMode />
  <GameWindow game={editor} />
</main>
