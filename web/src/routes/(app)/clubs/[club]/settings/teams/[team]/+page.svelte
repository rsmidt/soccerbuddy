<script lang="ts">
  import type { PageData } from "./$types";
  import { Drawer } from "vaul-svelte";
  import { createClient } from "@connectrpc/connect";
  import { TeamService } from "$lib/gen/soccerbuddy/team/v1/team_service_connect";
  import { defaultTransport } from "$lib/client.js";
  import NoBreak from "$lib/components/NoBreak.svelte";
  import Autocomplete from "$lib/components/Autocomplete.svelte";
  import type { AutocompleteOption } from "$lib/components/autocomplete";
  import { debounce } from "$lib/debounce";
  import { pbToRole } from "$lib/protobuf";
  import { runGrpc } from "$lib/runGrpc";
  import DataList from "$lib/components/list/DataList.svelte";
  import { goto, invalidateAll } from "$app/navigation";
  import { Section } from "$lib/components/section";
  import ListAction from "$lib/components/list/ListAction.svelte";
  import type { Component } from "svelte";
  import CarbonSoccer from "virtual:icons/carbon/soccer";
  import MdiClipboardPulseOutline from "virtual:icons/mdi/clipboard-pulse-outline";
  import { invariant } from "$lib/invariant";
  import { ToastControl } from "$lib/toasts/control.svelte";
  import SimpleMessageToast from "$lib/toasts/variants/SimpleMessageToast.svelte";

  type Role = {
    key: string;
    name: string;
    icon: Component;
  };
  const roles: Role[] = [
    { key: "PLAYER", name: "Spieler:in", icon: CarbonSoccer },
    { key: "COACH", name: "Trainer:in", icon: MdiClipboardPulseOutline },
  ] as const;
  type RoleKey = (typeof roles)[number]["key"];

  const { data }: { data: PageData } = $props();

  const toastControl = ToastControl.getGlobal();
  const client = createClient(TeamService, defaultTransport(fetch));

  async function deleteTeam() {
    const confirmed = window.confirm("Möchtest du das Team wirklich löschen?");
    if (!confirmed) return;
    await client.deleteTeam({ teamId: data.team.id });
    await goto(".");
  }

  let options = $state<AutocompleteOption[]>([]);
  const searchPersons = debounce(async (query: string) => {
    const { persons } = await runGrpc(window.location, () =>
      client.searchPersonsNotInTeam({
        query,
        teamId: data.team.id,
      }),
    );
    options = persons.map((person) => ({
      id: person.id,
      name: `${person.firstName} ${person.lastName}`,
    }));
  }, 250);
  const handleFetchOptions = (query: string) => {
    if (query.trim().length < 3) {
      searchPersons.cancel();
      options = [];
      return;
    }
    searchPersons(query);
  };

  let showAddPerson = $state(false);
  let selectedPerson = $state(undefined as AutocompleteOption | undefined);
  let selectedRole = $state(undefined as RoleKey | undefined);

  const handleDrawerClose = () => {
    showAddPerson = false;
    selectedPerson = undefined;
    selectedRole = undefined;
  };
  const selectPerson = (option: AutocompleteOption) => {
    selectedPerson = option;
    selectedRole = undefined;
  };
  const addPlayerToTeam = async () => {
    invariant(selectedPerson, "Person must be selected");
    invariant(selectedRole, "Role must be selected");

    await runGrpc(window.location, () =>
      client.addPersonToTeam({
        personId: selectedPerson!!.id,
        teamId: data.team.id,
        role: selectedRole,
      }),
    );
    handleDrawerClose();
    await invalidateAll();
    toastControl.add({
      component: SimpleMessageToast,
      props: {
        message: "Person erfolgreich hinzugefügt",
        type: "success" as const,
      },
    });
  };

  const membersSortedByName = data.members.sort((a, b) => b.lastName.localeCompare(a.lastName));
</script>

<h1 class="default-page-header">Team-Einstellungen</h1>
<p class="description">
  Hier kannst du die Einstellungen für das Team
  <NoBreak><b>{data.team.name}</b></NoBreak>
  bearbeiten.
</p>

<Section>
  {#snippet title()}
    Personen
  {/snippet}
  <DataList data={membersSortedByName}>
    {#snippet addMoreRow()}
      <Drawer.Root open={showAddPerson} onClose={handleDrawerClose} shouldScaleBackground>
        <ListAction onclick={() => (showAddPerson = !showAddPerson)}>Person hinzufügen</ListAction>
        <Drawer.Portal>
          <Drawer.Overlay class="drawer-overlay" />
          <Drawer.Content class="drawer-content">
            <div class="person-selector">
              <div class="drawer-handle"></div>
              <div class="input-group">
                <label for="person-search">Durchsuche Personen in deinem Verein</label>
                <Autocomplete
                  id="person-search"
                  placeholder="Maxine Musterfrau..."
                  fetchOptions={handleFetchOptions}
                  {options}
                  select={selectPerson}
                />
              </div>
              {#if selectedPerson}
                <div>
                  <p><b>Ausgewählt: {selectedPerson.name}</b></p>
                  <p class="label">Rolle im Team</p>
                  <div class="roles">
                    {#each roles as role (role.key)}
                      <label class="role">
                        <input
                          type="radio"
                          name="role"
                          value={role.key}
                          bind:group={selectedRole}
                        />
                        <role.icon style="width: 3rem; height: 4rem" />
                        <span class="role-name">{role.name}</span>
                      </label>
                    {/each}
                  </div>
                  <button disabled={!selectedRole} onclick={addPlayerToTeam}>Hinzufügen</button>
                </div>
              {/if}
            </div>
          </Drawer.Content>
        </Drawer.Portal>
      </Drawer.Root>
    {/snippet}}
    {#snippet item(itemData: (typeof data.members)[0])}
      [{itemData.role}] {itemData.firstName} {itemData.lastName}
    {/snippet}
  </DataList>
</Section>

<Section class="delete">
  {#snippet title()}
    Gefahrenzone
  {/snippet}
  {#snippet body()}
    <p>Das Löschen eines Teams kann nicht rückgängig gemacht werden.</p>
    <p>Alle zugehörigen Daten werden unwiderruflich gelöscht:</p>
    <ul>
      <li>Alle Trainings</li>
      <li>Alle Bewertungen</li>
    </ul>
    <button class="delete danger" onclick={deleteTeam}>Team löschen</button>
  {/snippet}
</Section>

<style>
  .delete {
    width: 100%;
    font-size: var(--text-sm);
  }

  :global(.drawer-overlay) {
    position: fixed;
    inset: 0;
    background-color: var(--bg-100);
  }

  :global(.drawer-content) {
    background-color: var(--bg-200);
    border-top-left-radius: 1rem;
    border-top-right-radius: 1rem;
    height: 100%;
    max-height: 96%;
    position: fixed;
    bottom: 0;
    right: 0;
    left: 0;
    transition: all;
    transition-duration: 0.2s;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  }

  .person-selector {
    padding: 1rem;
  }

  .drawer-handle {
    margin: 0 auto 1rem auto;
    height: 0.5rem;
    width: 3rem;
    flex-shrink: 0;
    border-radius: 1rem;
    background-color: var(--bg-500);
  }

  input[type="radio"] {
    position: absolute;
    opacity: 0;
  }

  .roles {
    display: flex;
    column-gap: 1rem;
    justify-content: center;
    margin-bottom: 1rem;
  }

  .role {
    display: flex;
    flex-direction: column;
    align-items: center;
    border: solid 1px var(--bg-300);
    border-radius: 0.5rem;
    aspect-ratio: 1/1;
    padding: 1rem;

    .role-name {
      font-size: var(--text-sm);
    }
  }

  .role:has(input[type="radio"]:checked) {
    border-color: var(--primary-300);
  }

  .label {
    margin-top: 0;
  }
</style>
