<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLAnchorAttributes } from "svelte/elements";
  import MaterialSymbolsChevronRightRounded from "virtual:icons/material-symbols/chevron-right-rounded";

  const {
    children,
    icon,
    ...rest
  }: {
    children: Snippet;
    icon?: Snippet;
  } & HTMLAnchorAttributes = $props();
</script>

<a data-sveltekit-preload-data="tap" {...rest}>
  {@render children()}
  {#if icon}
    {@render icon()}
  {:else}
    <MaterialSymbolsChevronRightRounded style="color: var(--bg-500)" height="28px" width="28px" />
  {/if}
</a>

<style>
  a {
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: var(--text-100);
    text-decoration: none;
  }
</style>
