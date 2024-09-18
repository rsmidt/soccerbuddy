<script lang="ts">
  import MaterialSymbolsCheckCircleRounded from "virtual:icons/material-symbols/check-circle-rounded";
  import MaterialSymbolsCloseRounded from "virtual:icons/material-symbols/close-rounded";

  const {
    dismiss,
    link,
    linkText,
    entityName,
    message,
  }: {
    message: string;
    link: string;
    linkText: string;
    entityName: string;
    dismiss?: () => void;
  } = $props();
</script>

<div class="toast">
  <div class="toast-indicator-icon">
    <MaterialSymbolsCheckCircleRounded style="color: var(--success)" />
  </div>
  <div class="toast-content">
    <p>{message}</p>
    <p class="toast-new-entity"><b>{entityName}</b></p>
    <a href={link} onclick={dismiss}>{linkText}</a>
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

    a {
      color: var(--primary-500);
      text-decoration: underline;
    }
  }
</style>
