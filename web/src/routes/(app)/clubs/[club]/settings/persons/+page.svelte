<script lang="ts">
  import type { PageData } from "./$types";
  import NoBreak from "$lib/components/NoBreak.svelte";
  import DataList from "$lib/components/list/DataList.svelte";
  import ListLink from "$lib/components/list/ListLink.svelte";

  const { data }: { data: PageData } = $props();
</script>

<h1 class="default-page-header">Alle Personen</h1>
<p class="description">
  Hier findest du alle Personen in deinem Verein
  <NoBreak><b>{data.club.name}</b></NoBreak>
  .
</p>

<DataList data={data.persons}>
  {#snippet addMoreRow()}
    <ListLink href="persons/add">Neue Person hinzufügen</ListLink>
  {/snippet}}
  {#snippet item(itemData: (typeof data.persons)[0])}
    <ListLink href="persons/{itemData.id}">
      {itemData.firstName}
      {itemData.lastName}
    </ListLink>
  {/snippet}
</DataList>

<style>
</style>
