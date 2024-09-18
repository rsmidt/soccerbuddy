<script lang="ts">
  import MaterialSymbolsChevronRightRounded from "virtual:icons/material-symbols/chevron-right-rounded";
  import MaterialSymbolsNewWindowRounded from "virtual:icons/material-symbols/new-window-rounded";
  import type { PageData } from "./$types";
  import NoBreak from "$lib/components/NoBreak.svelte";

  const { data }: { data: PageData } = $props();
</script>

<h1 class="default-page-header">Alle Teams</h1>
<p class="description">
  Hier findest du alle Teams deines Vereins <NoBreak><b>{data.club.name}</b></NoBreak>.
</p>

<div class="teams-list">
  <div class="team new">
    <a href="teams/add" class="team-name">
      Neue Mannschaft anlegen
      <MaterialSymbolsNewWindowRounded
        style="color: var(--primary-500); margin-right: 4px"
        height="20px"
        width="20px"
      />
    </a>
  </div>
  {#each data.teams as team (team.id)}
    <div class="team">
      <a data-sveltekit-preload-data="tap" href={`teams/${team.slug}`} class="team-name">
        {team.name}
        <MaterialSymbolsChevronRightRounded
          style="color: var(--bg-500)"
          height="28px"
          width="28px"
        />
      </a>
    </div>
  {/each}
</div>

<style>
  .teams-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(15rem, 1fr));
    background-color: var(--bg-200);
    border-radius: 0.5rem;
  }

  .team {
    transition-property: background-color;
    transition-duration: 0.15s;
    transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);

    &:first-of-type {
      border-top-left-radius: 0.5rem;
      border-top-right-radius: 0.5rem;
    }

    &:last-of-type {
      border-bottom-left-radius: 0.5rem;
      border-bottom-right-radius: 0.5rem;
    }

    & + & {
      border-top: 1px solid hsl(var(--gray-300));
    }

    &:hover {
      background-color: var(--bg-300);
    }
  }

  .team-name {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 1rem;
    border-radius: 0.5rem;
    color: var(--text-100);
    text-decoration: none;
  }
</style>
