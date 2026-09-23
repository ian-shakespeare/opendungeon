<script lang="ts">
  import * as GLM from "gl-matrix";
  import Renderer from "$lib/renderer";
  import { PerspectiveCamera, type Camera } from "$lib/renderer/camera";
  import { onMount } from "svelte";
  import Character from "$lib/character";

  let canvas = $state<HTMLCanvasElement>()!;
  let loading = $state(true);
  let frameHandle: number;
  let characterId: number;
  let renderer: Renderer;
  let camera: Camera;

  onMount(() => {
    renderer = new Renderer(canvas, {
      resizeToWindow: true,
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
</script>

<div class="relative">
  <canvas bind:this={canvas}> </canvas>
</div>
