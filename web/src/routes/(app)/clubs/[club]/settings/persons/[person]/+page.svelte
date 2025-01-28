<script lang="ts">
  import type { PageData } from "./$types";
  // noinspection ES6UnusedImports
  import { Drawer } from "vaul-svelte";
  import NoBreak from "$lib/components/NoBreak.svelte";
  import { type GetPersonOverviewResponse, PersonService } from "$lib/gen/soccerbuddy/person/v1/person_service_pb";
  import { pbToAccountLink, pbToRole } from "$lib/protobuf";
  import ListAction from "$lib/components/list/ListAction.svelte";
  import { ConnectError, createClient } from "@connectrpc/connect";
  import { defaultTransport } from "$lib/client.js";
  import { runGrpc } from "$lib/runGrpc";
  import { Section } from "$lib/components/section";
  import { type Timestamp, timestampDate } from "@bufbuild/protobuf/wkt";
  import { AccountLink } from "$lib/gen/soccerbuddy/shared_pb";
  import IconParentalGuidance from "virtual:icons/soccerbuddy/parental-guidance";
  import IconSelfAwareness from "virtual:icons/soccerbuddy/self-awareness";
  import { ToastControl } from "$lib/toasts/control.svelte";
  import SimpleMessageToast from "$lib/toasts/variants/SimpleMessageToast.svelte";
  import { TooManyLinksCreatedErrorSchema } from "$lib/gen/soccerbuddy/person/v1/errors_pb";

  const { data }: { data: PageData } = $props();
  const person: GetPersonOverviewResponse = data.person;

  const toastControl = ToastControl.getGlobal();

  let showLinkCreator = $state(false);
  let selectedRole = $state(undefined as string | undefined);
  const client = createClient(PersonService, defaultTransport(fetch));
  const handleNewConnectionClick = async () => {
    const linkAs =
      selectedRole === "parent" ? AccountLink.LINKED_AS_PARENT : AccountLink.LINKED_AS_SELF;
    let result;
    try {
      result = await runGrpc(window.location, () => {
        return client.initiatePersonAccountLink({ personId: person.id, linkAs });
      });
    } catch (e) {
      const cErr = ConnectError.from(e);
      if (cErr.findDetails(TooManyLinksCreatedErrorSchema).length > 0) {
        toastControl.add({
          component: SimpleMessageToast,
          props: {
            message: "Die Person hat bereits zu viele offene Links.",
            type: "warn" as const,
          },
        })
      }
      handleDrawerClose();
      return;
    }
    const url = `${process.env.ORIGIN}/connect/${result.linkToken}`;
    try {
      await navigator.share({
        url,
        title: "Deine Einladung zu SoccerBuddy",
        text: "Du wurdest eingeladen, SoccerBuddy beizutreten und dich zu verknüpfen.",
      });
    } catch {
      navigator.clipboard
        .writeText(url)
        .then(() =>
          toastControl.add({
            component: SimpleMessageToast,
            props: {
              message: "Link erfolgreich in die Zwischenablage kopiert.",
              type: "success" as const,
            },
          }),
        )
        .catch((e) => alert(e));
    } finally {
      handleDrawerClose()
    }
  };

  function handleDrawerClose() {
    showLinkCreator = false;
    selectedRole = undefined;
  }

  function formatTimestamp(timestamp: Timestamp | undefined): string {
    return timestamp ? timestampDate(timestamp).toLocaleDateString() : "";
  }

  let maybeSelfLinkedAccount = $derived(
    person.linkedAccounts.find((account) => account.linkedAs === AccountLink.LINKED_AS_SELF),
  );
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
        <Drawer.Root open={showLinkCreator} onClose={handleDrawerClose} shouldScaleBackground>
          <ListAction onclick={() => (showLinkCreator = !showLinkCreator)}
            >Neue Verknüpfung
          </ListAction>
          <Drawer.Portal>
            <Drawer.Overlay class="drawer-overlay" />
            <Drawer.Content class="drawer-content">
              <div class="drawer-handle"></div>
              <p class="label">Wie ist das Verhältnis von dem Account zu dieser Person?</p>
              <div class="types">
                <label class="type">
                  <input
                    type="radio"
                    name="type"
                    value="self"
                    disabled={maybeSelfLinkedAccount !== undefined}
                    bind:group={selectedRole}
                  />
                  <IconSelfAwareness fill="currentColor" style="width: 3rem; height: 4rem" />
                  <span class="type-name">Die Person selbst</span>
                </label>
                <label class="type">
                  <input type="radio" name="type" value="parent" bind:group={selectedRole} />
                  <IconParentalGuidance fill="currentColor" style="width: 4rem; height: 4rem" />
                  <span class="type-name">Eine Bezugsperson</span>
                </label>
              </div>
              {#if maybeSelfLinkedAccount !== undefined}
                <p class="error">
                  {maybeSelfLinkedAccount.fullName} ist bereits als Person selbst verknüpft.
                </p>
              {/if}
              <button disabled={!selectedRole} onclick={handleNewConnectionClick}>Erstellen</button>
            </Drawer.Content>
          </Drawer.Portal>
        </Drawer.Root>
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

<style>
  .types {
    display: flex;
    column-gap: 1rem;
    justify-content: center;
    margin-bottom: 1rem;
  }

  .type {
    display: flex;
    flex-direction: column;
    align-items: center;
    border: solid 1px var(--bg-300);
    border-radius: 0.5rem;
    aspect-ratio: 1/1;
    padding: 1rem;

    .type-name {
      font-size: var(--text-sm);
    }
  }

  .type:has(input[type="radio"]:checked) {
    border-color: var(--primary-300);
  }

  .type:has(input[type="radio"]:disabled) {
      color: hsl(var(--gray-500))
  }

  :global(.drawer-content) {
    width: 100%;
    max-width: 500px;
    margin: 0 auto;
    max-height: 70%;
    padding: 1rem;
  }
</style>
