<script lang="ts">
  import type { PageData } from "./$types";
  import { superForm } from "sveltekit-superforms";
  import { ToastControl } from "$lib/toasts/control.svelte";
  import SimpleMessageToast from "$lib/toasts/variants/SimpleMessageToast.svelte";

  const { data }: { data: PageData } = $props();
  const toastControl = ToastControl.getGlobal();

  const { form, enhance, errors, constraints, submitting } = superForm(data.form, {
    resetForm: false,
    delayMs: 250,
    onSubmit: ({ formData }) => {
      formData.set("redirectPath", data.redirectPath);
    },
    onUpdate: ({ form, result }) => {
      if (result.status === 401) {
        toastControl.add({
          component: SimpleMessageToast,
          props: {
            message: form.message,
            type: "error" as const,
          },
        });
      }

      if (form.valid && result.type === "success") {
        toastControl.add({
          component: SimpleMessageToast,
          props: {
            message: "Erfolgreich eingeloggt.",
          },
        });
      }
    },
  });
</script>

<h1 class="default-page-header">Login</h1>

<form method="post" action="?/login" use:enhance>
  <div class="input-group">
    <label for="email">E-Mail</label>
    <input {...$constraints.email} type="email" id="email" name="email" bind:value={$form.email} />
    {#if $errors.email}
      <p class="error">{$errors.email}</p>
    {/if}
  </div>

  <div class="input-group">
    <label for="password">Passwort</label>
    <input
      {...$constraints.password}
      type="password"
      id="password"
      name="password"
      bind:value={$form.password}
    />
    {#if $errors.password}
      <p class="error">{$errors.password}</p>
    {/if}
  </div>

  <button type="submit" disabled={$submitting}>
    {$submitting ? "Wird eingeloggt..." : "Einloggen"}
  </button>
</form>
