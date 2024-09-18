<script>
  import { page } from "$app/stores";
</script>

{#if $page.status === 404}
  <p>Die angeforderte Seite konnte nicht gefunden werden :(</p>
  <a href="/">Geh zurück zum Start</a>
{:else if $page.status === 500}
  <p>Das ist ein kleines Oopsie passiert 🙈</p>
  <p>Im Idealfall weiß der Admin bereits Bescheid</p>
  <a href="/">Geh zurück zum Start</a>
{:else}
  <p>{$page.status}: {$page.error?.message}</p>
{/if}
