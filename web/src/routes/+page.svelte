<script lang="ts">
  import { configureScreen } from "$lib/components/screen/screen.svelte";

  configureScreen({ backButtonShown: false })
</script>

<svelte:head>
  <title>Home</title>
  <meta name="description" content="SoccerBuddy Management" />
</svelte:head>

<section>Nothing to see here!</section>

<style>
  section {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    flex: 0.6;
  }
</style>
