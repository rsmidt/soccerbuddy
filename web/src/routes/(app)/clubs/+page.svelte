<script lang="ts">
  import type { PageProps } from "./$types";
  import DataList from "$lib/components/list/DataList.svelte";
  import ListLink from "$lib/components/list/ListLink.svelte";
  import type { ListClubsResponse_Club } from "$lib/gen/soccerbuddy/club/v1/club_service_pb";

  const { data }: PageProps = $props();

</script>

<h1 class="default-page-header">Alle Clubs</h1>

<DataList data={data.clubs}>
  {#snippet addMoreRow()}
    <ListLink href="clubs/add">Neuen Club hinzufügen</ListLink>
  {/snippet}}
  {#snippet item(itemData: ListClubsResponse_Club)}
    <ListLink href="clubs/{itemData.slug}/settings">
      {itemData.name}
    </ListLink>
  {/snippet}
</DataList>
