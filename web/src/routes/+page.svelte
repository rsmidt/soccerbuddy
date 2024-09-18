<script lang="ts">
  import { invalidateAll } from "$app/navigation";

  function rerun() {
    invalidateAll();
  }
</script>

<svelte:head>
  <title>Home</title>
  <meta name="description" content="Svelte demo app" />
</svelte:head>

<section>Hello</section>

<button onclick={rerun}> Rerun </button>

<style>
  section {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    flex: 0.6;
  }
</style>
