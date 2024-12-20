<script lang="ts">
  import type { PageData } from "./$types";
  import NoBreak from "$lib/components/NoBreak.svelte";
  import { type GetPersonOverviewResponse, PersonService } from "$lib/gen/soccerbuddy/person/v1/person_service_pb";
  import { pbToAccountLink, pbToRole } from "$lib/protobuf";
  import ListAction from "$lib/components/list/ListAction.svelte";
  import { createClient } from "@connectrpc/connect";
  import { defaultTransport } from "$lib/client.js";
  import { runGrpc } from "$lib/runGrpc";
  import { Section } from "$lib/components/section";
  import { type Timestamp, timestampDate } from "@bufbuild/protobuf/wkt";
  import { AccountLink } from "$lib/gen/soccerbuddy/shared_pb";

  const { data }: { data: PageData } = $props();
  const person: GetPersonOverviewResponse = data.person;

  const client = createClient(PersonService, defaultTransport(fetch));
  const handleNewTeamClick = async () => {
    const result = await runGrpc(window.location, () => {
      return client.initiatePersonAccountLink({ personId: person.id, linkAs: AccountLink.LINKED_AS_SELF });
    });
    try {
      await navigator.share({
        title: "Deine Einladung zu SoccerBuddy",
        text: "Du wurdest eingeladen, SoccerBuddy beizutreten und dich zu verknüpfen.",
        url: "https://soccerbuddy.app/signup?code=" + result.linkToken,
      });
    } catch (e) {
      console.log(e);
    }
  };

  function formatTimestamp(timestamp: Timestamp | undefined): string {
    return timestamp ? timestampDate(timestamp).toLocaleDateString() : "";
  }
</script>

<h1 class="default-page-header">Person-Einstellungen</h1>
<p class="description">
  Hier kannst du die Person
  <NoBreak><b>{data.person.firstName} {data.person.lastName}</b></NoBreak>
  bearbeiten.
</p>

{#snippet dataRow(key: string, value: string)}
  <div class="row">
    <div class="item">
      <b>{key}:</b>
      {value}
    </div>
  </div>
{/snippet}

<Section>
  {#snippet title()}
    Details
  {/snippet}
  <div class="list">
    {@render dataRow("Name", `${person.firstName} ${person.lastName}`)}
    {@render dataRow("Geburtstag", formatTimestamp(person.birthdate))}
    {@render dataRow("Erstellt", formatTimestamp(person.createdAt))}
    {@render dataRow("Erstellt von", person.createdBy?.fullName ?? "")}
  </div>
</Section>

<Section>
  {#snippet title()}
    Verknüpfte Accounts
  {/snippet}
  <div class="list">
    <div class="row">
      <div class="item">
        <ListAction onclick={handleNewTeamClick}>Neue Verknüpfung</ListAction>
      </div>
    </div>
    {#each person.linkedAccounts as linkedAccount}
      {@render dataRow(
        pbToAccountLink(linkedAccount.linkedAs),
        `${linkedAccount.fullName} ${linkedAccount.actor.case === "invite" ? `eingeladen von ${linkedAccount.actor.value.invitedBy?.fullName} am ${formatTimestamp(linkedAccount.actor.value.invitedAt)} ` : ""}`,
      )}
    {/each}
  </div>
</Section>

{#if person.pendingAccountLinks.length > 0}
  <Section>
    {#snippet title()}
      Offene Verknüpfungen
    {/snippet}
    <div class="list">
      {#each person.pendingAccountLinks as pendingLink}
        <div class="row">
          <div class="item">
            {`Typ "${pbToAccountLink(pendingLink.linkedAs)}" erstellt von ${pendingLink.invitedBy?.fullName}`}
            <br />
            {`gültig bis ${formatTimestamp(pendingLink.expiresAt)}`}
          </div>
        </div>
      {/each}
    </div>
  </Section>
{/if}

<Section>
  {#snippet title()}Mannschaft{/snippet}
  <div class="list">
    {#each person.teams as team}
      {@render dataRow(pbToRole(team.role), `${team.name}`)}
    {/each}
  </div>
</Section>
