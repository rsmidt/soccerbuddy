<script lang="ts">
  import type { PageProps } from "./$types";
  import ListLink from "$lib/components/list/ListLink.svelte";
  import { configureScreen } from "$lib/components/screen/screen.svelte";

  let { data }: PageProps = $props();

  configureScreen({
    fallbackBackUrl: "/clubs"
  })
</script>

<main>
  <h1 class="default-page-header">Einstellungen für {data.club.name}</h1>

  <div class="list">
    <div class="row">
      <div class="item">
        <ListLink href="settings/teams">Mannschaften</ListLink>
      </div>
    </div>
    <div class="row">
      <div class="item">
        <ListLink href="settings/persons">Personen</ListLink>
      </div>
    </div>
  </div>
</main>
