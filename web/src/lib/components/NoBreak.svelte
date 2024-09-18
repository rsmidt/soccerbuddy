<script lang="ts">
  import type { Snippet } from "svelte";

  const { children }: { children: Snippet } = $props();
</script>

<span>{@render children()}</span>

<style>
  span {
    white-space: nowrap;
  }
</style>
