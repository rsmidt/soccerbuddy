<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLAttributes } from "svelte/elements";

  type SectionProps = Omit<HTMLAttributes<any>, "title" | "body" | "children">;

  const {
    title,
    body,
    children,
    ...rest
  }: { title?: Snippet; body?: Snippet; children?: Snippet } & SectionProps = $props();
</script>

<section {...rest} class="section {rest.class}">
  {#if title}
    <h2>{@render title()}</h2>
  {/if}
  {@render children?.()}
  {#if body}
    <div class="body">{@render body()}</div>
  {/if}
</section>

<style>
  .section {
    margin-top: 2rem;

    h2 {
      margin-top: 0;
      margin-bottom: 0.5rem;
      font-size: var(--text-lg);
    }

    ul {
      margin-bottom: 1rem;
      list-style: disc;
      padding-inline-start: 1rem;
    }

    :global(p, ul) {
      margin-top: 0;
      margin-bottom: 1rem;
      font-size: var(--text-sm);
      color: var(--text-300);
    }

    .body {
      border-radius: 0.5rem;
      background-color: var(--bg-200);
      padding: 0.75rem;
    }

    * + .body {
      margin-top: 1rem;
    }
  }
</style>
