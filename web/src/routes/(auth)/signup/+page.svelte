<script lang="ts">
  import type { PageData } from "./$types";
  import { ToastControl } from "$lib/toasts/control.svelte";
  import { superForm } from "sveltekit-superforms";
  import SimpleMessageToast from "$lib/toasts/variants/SimpleMessageToast.svelte";
  import { goto } from "$app/navigation";

  const { data }: { data: PageData } = $props();
  const toastControl = ToastControl.getGlobal();

  let connectRedirect = $derived(encodeURIComponent(`/connect/${data.inviteCode}`));

  const { form, enhance, errors, constraints, submitting } = superForm(data.form, {
    resetForm: false,
    delayMs: 250,
    onSubmit: ({ formData }) => {
      formData.set("inviteCode", data.inviteCode);
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
            message: "Account erfolgreich erstellt. Du wirst gleich weitergeleitet.",
          },
        });
        setTimeout(async () => {
          await goto(`/connect/${data.inviteCode}`);
        }, 5000)
      }
    },
  });
</script>

<h1 class="default-page-header">Registrierung</h1>

<p>Bist du schon registriert? Dann melde dich <a href="/login?redirect={connectRedirect}">hier</a> an.</p>

<form method="post" action="?/register" use:enhance>
  <div class="input-group">
    <label for="firstName">Vorname</label>
    <input {...$constraints.firstName} type="text" id="firstName" name="firstName" bind:value={$form.firstName} />
    {#if $errors.firstName}
      <p class="error">{$errors.firstName}</p>
    {/if}
  </div>
  <div class="input-group">
    <label for="lastName">Nachname</label>
    <input {...$constraints.lastName} type="text" id="lastName" name="lastName" bind:value={$form.lastName} />
    {#if $errors.lastName}
      <p class="error">{$errors.lastName}</p>
    {/if}
  </div>

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
  <div class="input-group">
    <label for="passwordConfirmation">Passwort (Bestätigung)</label>
    <input
      {...$constraints.passwordConfirmation}
      type="password"
      id="passwordConfirmation"
      name="passwordConfirmation"
      bind:value={$form.passwordConfirmation}
    />
    {#if $errors.passwordConfirmation}
      <p class="error">{$errors.passwordConfirmation}</p>
    {/if}
  </div>

  <button type="submit" disabled={$submitting}>
    {$submitting ? "Wird eingeloggt..." : "Einloggen"}
  </button>
</form>
