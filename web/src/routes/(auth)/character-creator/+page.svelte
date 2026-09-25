<script lang="ts">
  import * as GLM from "gl-matrix";
  import Renderer from "$lib/renderer";
  import { PerspectiveCamera, type Camera } from "$lib/renderer/camera";
  import { onMount } from "svelte";
  import Character from "$lib/character";

  let canvas = $state<HTMLCanvasElement>()!;
  let loading = $state(true);
  let isDragging = $state(false);
  let shoulderWidth = $state(1);
  let frameHandle: number;
  let characterId: number;
  let renderer: Renderer;
  let camera: Camera;

  onMount(() => {
    renderer = new Renderer(canvas, {
      backgroundColor: new Float32Array([0, 0, 0, 1]),
    });
    camera = new PerspectiveCamera(canvas.width / canvas.height);
    camera.zoom = 3;

    Character.create(renderer).then((character) => {
      characterId = renderer.loadElement(character);
      character.translate(GLM.vec3.fromValues(0, -1.25, -0.5));
      character.setSkinTone(new Float32Array([0.882, 0.607, 0.49, 1.0]));

      loading = false;
    });

    loop();

    return () => {
      window.cancelAnimationFrame(frameHandle);
    };
  });

  $effect(() => {
    shoulderWidth;

    console.log("$effect");
    if (characterId === undefined) {
      return;
    }

    const character = renderer.getElement<Character>(characterId);
    const translation = 1 - shoulderWidth;
    character.setJointTranslation("DEF-shoulder.R", GLM.vec3.fromValues(translation, 0, 0));
    character.setJointTranslation("DEF-shoulder.L", GLM.vec3.fromValues(-translation, 0, 0));
    character.updateTransforms();
    character.computeSkinningMatrix();
  });

  function tick(time: number) {}

  function draw() {
    if (loading) {
      return;
    }

    renderer.clear();

    const character = renderer.getAndUseElement<Character>(characterId);
    character.setCamera(camera);
    character.draw();
  }

  function loop() {
    frameHandle = window.requestAnimationFrame((ms) => {
      const time = ms / 1000;
      tick(time);
      draw();
      loop();
    });
  }

  function handleSpin(event: MouseEvent) {
    if (!isDragging) {
      return;
    }

    const target = event.target as HTMLCanvasElement;
    const width = target.clientWidth;
    const delta = event.movementX / width;
    const theta = 2 * delta * Math.PI;
    const character = renderer.getElement<Character>(characterId);
    character.rotateY(theta);
  }
</script>

<main class="relative w-full h-full">
  <div class="relative z-10">
    <label for="shoulder-width">Shoulder Width ({shoulderWidth})</label>
    <input
      bind:value={shoulderWidth}
      type="range"
      id="shoulder-width"
      name="shoulder-width"
      step={0.1}
      min={0.7}
      max={1.5}
    />
  </div>
  <canvas
    width={640}
    height={480}
    bind:this={canvas}
    onpointerdown={() => (isDragging = true)}
    onpointerup={() => (isDragging = false)}
    onpointerout={() => (isDragging = false)}
    onpointermove={handleSpin}
    class="absolute inset-0 w-full h-full"
  >
  </canvas>
</main>

<style>
  canvas {
    image-rendering: crisp-edges; /* for firefox */
    image-rendering: pixelated;
  }
</style>
