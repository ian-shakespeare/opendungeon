<script lang="ts">
  import { onMount } from "svelte";
  import type { HTMLCanvasAttributes } from "svelte/elements";

  type Props = {
    draw: () => void;
    tick?: (time: number) => void;
    ondbltap?: (event: KeyboardEvent) => void;
    onkeydown?: (event: KeyboardEvent) => void;
    onkeyup?: (event: KeyboardEvent) => void;
    canvas?: HTMLCanvasElement;
  } & Omit<HTMLCanvasAttributes, "onkeydown" | "onkeyup">;

  let {
    draw,
    tick,
    ondbltap,
    onkeydown,
    onkeyup,
    canvas = $bindable(),
    ...props
  }: Props = $props();

  let frameHandle = -1;
  let lastPressKey: string | null = null;
  let lastPressTime = 0;

  onMount(() => {
    loop();

    return () => {
      window.cancelAnimationFrame(frameHandle);
    };
  });

  function loop() {
    frameHandle = window.requestAnimationFrame((ms) => {
      const time = ms / 1000;
      tick?.(time);
      draw();
      loop();
    });
  }

  function handleKeyDown(event: KeyboardEvent) {
    const target = event.target as HTMLElement | null;
    if (
      target?.tagName === "INPUT" ||
      target?.tagName === "TEXTAREA" ||
      target?.isContentEditable
    ) {
      return;
    }

    const currentTime = Date.now();
    const delay = lastPressTime - currentTime;

    if (event.key === lastPressKey && delay < 250) {
      ondbltap?.(event);

      lastPressKey = null;
      lastPressTime = 0;
    } else {
      onkeydown?.(event);

      lastPressKey = event.key;
      lastPressTime = currentTime;
    }
  }
</script>

<svelte:window onkeydown={handleKeyDown} {onkeyup} />

<canvas bind:this={canvas} {...props}></canvas>
