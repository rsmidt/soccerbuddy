<script lang="ts">
  import "../app.css";
  import ToastsHost from "$lib/toasts/ToastsHost.svelte";

  const { children } = $props();
</script>

<ToastsHost>
  <main>
    {@render children()}
  </main>
</ToastsHost>

<style>
  main {
    padding: 1rem;
    width: 100%;
    margin: 0 auto;
  }
</style>
