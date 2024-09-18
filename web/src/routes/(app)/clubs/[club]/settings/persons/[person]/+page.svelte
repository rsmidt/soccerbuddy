<script lang="ts">
  import type { PageData } from "./$types";
  import NoBreak from "$lib/components/NoBreak.svelte";
  import {
    type GetPersonOverviewResponse,
    LinkedAs,
  } from "$lib/gen/soccerbuddy/person/v1/person_service_pb";
  import { pbToLinkedAs, pbToRole } from "$lib/protobuf";
  import ListAction from "$lib/components/list/ListAction.svelte";
  import { createClient } from "@connectrpc/connect";
  import { PersonService } from "$lib/gen/soccerbuddy/person/v1/person_service_connect";
  import { defaultTransport } from "$lib/client.js";
  import { runGrpc } from "$lib/runGrpc";
  import { Section } from "$lib/components/section";

  const { data }: { data: PageData } = $props();
  const person: GetPersonOverviewResponse = data.person;

  const client = createClient(PersonService, defaultTransport(fetch));
  const handleNewTeamClick = async () => {
    const result = await runGrpc(window.location, () => {
      return client.initiatePersonAccountLink({ personId: person.id, linkAs: LinkedAs.SELF });
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
    {@render dataRow("Geburtstag", person.birthdate?.toDate()?.toLocaleDateString() ?? "")}
    {@render dataRow("Erstellt", person.createdAt?.toDate().toLocaleDateString() ?? "")}
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
        pbToLinkedAs(linkedAccount.linkedAs),
        `${linkedAccount.fullName} ${linkedAccount.actor.case === "invite" ? `eingeladen von ${linkedAccount.actor.value.invitedBy?.fullName} am ${linkedAccount.actor.value.invitedAt?.toDate()?.toLocaleDateString()} ` : ""}`,
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
            {`Typ "${pbToLinkedAs(pendingLink.linkedAs)}" erstellt von ${pendingLink.invitedBy?.fullName}`}
            <br />
            {`gültig bis ${pendingLink.expiresAt?.toDate().toLocaleDateString()}`}
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
