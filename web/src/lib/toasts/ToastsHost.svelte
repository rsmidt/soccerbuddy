<script lang="ts">
  import { fade, fly } from "svelte/transition";
  import { quintOut } from "svelte/easing";
  import { ToastControl } from "$lib/toasts/control.svelte";
  import type { Snippet } from "svelte";

  const { children }: { children: Snippet } = $props();

  const control = ToastControl.createGlobal();

  function createDismissHandler(id: string) {
    return () => control.remove(id);
  }

  $effect(() => {
    if (control.toasts.length > 0) {
      const cb: EventListener = (event) => {
        if ((event.target as HTMLElement | null)?.closest(".toasts-container")) return;
        control.clear();
      };
      window.addEventListener("click", cb);
      return () => window.removeEventListener("click", cb);
    }
  });
</script>

{@render children()}

<div class="toasts-container">
  <div class="toasts-list">
    {#each control.toasts as toast (toast.id)}
      <div
        in:fly={{ duration: 300, y: 500, opacity: 0.5, easing: quintOut }}
        out:fade={{ duration: 200 }}
      >
        <toast.component {...toast.props} dismiss={createDismissHandler(toast.id)} />
      </div>
    {/each}
  </div>
</div>

<style>
  .toasts-container {
    position: fixed;
    bottom: 1rem;
    left: 1rem;
    right: 1rem;
    z-index: 1000;
  }

  .toasts-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
</style>
