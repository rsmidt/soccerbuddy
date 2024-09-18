<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLAttributes } from "svelte/elements";
  import MaterialSymbolsNewWindowRounded from "virtual:icons/material-symbols/new-window-rounded";

  const {
    children,
    ...rest
  }: {
    children: Snippet;
  } & HTMLAttributes<any> = $props();
</script>

<div {...rest}>
  {@render children()}
  <MaterialSymbolsNewWindowRounded
    style="color: var(--primary-500); margin-right: 4px"
    height="20px"
    width="20px"
  />
</div>

<style>
  div {
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: var(--text-100);
    text-decoration: none;
  }
</style>
