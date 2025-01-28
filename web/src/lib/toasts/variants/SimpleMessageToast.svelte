<script lang="ts">
  import MaterialSymbolsCheckCircleRounded from "virtual:icons/material-symbols/check-circle-rounded";
  import MaterialSymbolsCloseRounded from "virtual:icons/material-symbols/close-rounded";
  import MaterialSymbolsErrorOutline from "virtual:icons/material-symbols/error-outline";
  import MaterialSymbolsWarningOutline from "virtual:icons/material-symbols/warning-outline";

  const {
    dismiss,
    message,
    type = "success",
  }: {
    message: string;
    dismiss?: () => void;
    type?: "success" | "error" | "warn";
  } = $props();
</script>

<div class="toast">
  <div class="toast-indicator-icon">
    {#if type === "error"}
      <MaterialSymbolsErrorOutline style="color: var(--error)" />
    {:else if type === "warn"}
        <MaterialSymbolsWarningOutline style="color: var(--warning)" />
    {:else}
      <MaterialSymbolsCheckCircleRounded style="color: var(--success)" />
    {/if}
  </div>
  <div class="toast-content">
    <p>{message}</p>
  </div>
  <div class="toast-action">
    <MaterialSymbolsCloseRounded onclick={dismiss} />
  </div>
</div>

<style>
  .toast-indicator-icon {
    margin-top: 3px;
  }

  .toast-action {
    margin-left: auto;
    color: var(--text-400);

    &:hover {
      color: var(--text-100);
    }
  }

  .toast {
    background-color: var(--bg-200);
    color: var(--text-100);
    padding: 1rem;
    border: solid 1px var(--bg-300);
    border-radius: 0.5rem;
    display: flex;
    column-gap: 0.5rem;
  }

  .toast-content {
    display: flex;
    flex-direction: column;
    row-gap: 0.25rem;

    p {
      margin: 0;
    }
  }
</style>
