<script lang="ts" generics="Item extends { id: string }">
  import type { Snippet } from "svelte";

  interface Props {
    data: Item[];
    item: Snippet<[Item]>;
    addMoreRow?: Snippet;
  }

  const { item, addMoreRow, data }: Props = $props();
</script>

<div class="list">
  {#if addMoreRow}
    <div class="row new">
      <div class="item">
        {@render addMoreRow()}
      </div>
    </div>
  {/if}
  {#each data as dataItem (dataItem.id)}
    <div class="row">
      <div class="item">
        {@render item(dataItem)}
      </div>
    </div>
  {/each}
</div>
