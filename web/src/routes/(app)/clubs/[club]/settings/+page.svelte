<script lang="ts">
  import type { PageData } from "./$types";

  interface Props {
    data: PageData;
  }

  let { data }: Props = $props();
</script>

<main>
  <h1>Einstellungen für {data.club.name}</h1>

  <a href="settings/teams/add">Team anlegen</a>
</main>
