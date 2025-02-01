<script lang="ts">
  import type { LayoutProps } from "./$types";
  import "../app.css";
  import ToastsHost from "$lib/toasts/ToastsHost.svelte";
  import IconMaterialArrowBack from "virtual:icons/material-symbols/arrow-back";
  import { initializeScreenContext } from "$lib/components/screen/screen.svelte";
  import { goto } from "$app/navigation";

  const { children }: LayoutProps = $props();
  const context = initializeScreenContext();

  function handleBack(event: MouseEvent, backUrl: URL | string) {
    event.preventDefault();
    // Determine if it's an external target.
    if (typeof backUrl === "string" || backUrl.origin === location.origin) {
      context.state.popping = true;
      goto(backUrl, { replaceState: context.state.replaceBack }).then(() => {
        context.state.stack.pop();
        context.state.popping = false;
      });
    } else {
      window.location.href = backUrl.href;
    }
  }
</script>

<ToastsHost>
  <main>
    {#if context.state.backUrl}
      {@const backUrl = context.state.backUrl}
      <a class="back-link" href={context.backUrlHref} onclick={(event) => handleBack(event, backUrl)}>
        <IconMaterialArrowBack width={24} height={24} />
      </a>
    {/if}
    {@render children()}
  </main>
</ToastsHost>

<style>
  main {
    padding: 1rem;
    width: 100%;
    max-width: 500px;
    margin: 0 auto;
  }

   .back-link {
       text-decoration: none;
       color: white;
   }
</style>
