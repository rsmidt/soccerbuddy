<script lang="ts">
  import { configureScreen } from "$lib/components/screen/screen.svelte";
  import type { LayoutProps } from "./$types";

  configureScreen({ backButtonShown: false });

  const { children }: LayoutProps = $props();
</script>

{@render children()}