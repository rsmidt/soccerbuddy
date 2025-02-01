<script lang="ts">
  import type { PageProps } from "./$types";
  import DataList from "$lib/components/list/DataList.svelte";
  import ListLink from "$lib/components/list/ListLink.svelte";
  import MaterialSymbolsNewWindowRounded from "virtual:icons/material-symbols/new-window-rounded";
  import type { ListClubsResponse_Club } from "$lib/gen/soccerbuddy/club/v1/club_service_pb";
  import { configureScreen } from "$lib/components/screen/screen.svelte";

  const { data }: PageProps = $props();

  configureScreen({ backButtonShown: false });
</script>

<h1 class="default-page-header">Alle Clubs</h1>

{#if data.me.isSuper}
  <DataList data={data.clubs}>
    {#snippet addMoreRow()}
      <ListLink href="clubs/add">
        {#snippet icon()}
          <MaterialSymbolsNewWindowRounded
            style="color: var(--primary-500); margin-right: 4px"
            height="20px"
            width="20px"
          />
        {/snippet}
        Neuen Club hinzufügen
      </ListLink>
    {/snippet}}
    {#snippet item(itemData: ListClubsResponse_Club)}
      <ListLink href="clubs/{itemData.slug}/settings">
        {itemData.name}
      </ListLink>
    {/snippet}
  </DataList>
{:else}
  <DataList data={data.clubs}>
    {#snippet item(itemData: ListClubsResponse_Club)}
      <ListLink href="clubs/{itemData.slug}/settings">
        {itemData.name}
      </ListLink>
    {/snippet}
  </DataList>
{/if}
