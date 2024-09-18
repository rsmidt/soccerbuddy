<script lang="ts">
  import type { PageData } from "./$types";
  import { dateProxy, superForm } from "sveltekit-superforms";
  import { ToastControl } from "$lib/toasts/control.svelte";
  import EntityCreatedToast from "$lib/toasts/variants/EntityCreatedToast.svelte";
  import NoBreak from "$lib/components/NoBreak.svelte";

  const toastControl = ToastControl.getGlobal();
  const { data }: { data: PageData } = $props();
  const { form, enhance, errors, constraints, submitting } = superForm(data.form, {
    resetForm: false,
    delayMs: 250,
    onUpdate: ({ form, result }) => {
      if (form.valid && result.type === "success") {
        const { newPerson } = result.data;

        toastControl.add({
          component: EntityCreatedToast,
          props: {
            entityName: `${newPerson.firstName} ${newPerson.lastName}`,
            link: `/clubs/${data.club.slug}/settings/persons/${newPerson.id}`,
            linkText: "Zur Person",
            message: "Person wurde erfolgreich erstellt.",
          },
        });
      }
    },
    onSubmit: ({ formData }) => {
      formData.set("owningClubId", data.club.id);
    },
  });
  const birthdateProxy = dateProxy(form, "birthdate", { format: "date" });
</script>

<h1 class="default-page-header">Neue Person erstellen</h1>
<p class="description">
  Lege eine neue Person innerhalb deines Vereins
  <NoBreak><b>{data.club.name}</b></NoBreak>
  an.
</p>

<form method="post" action="?/add" use:enhance>
  <div class="input-group">
    <label for="firstName">Vorname</label>
    <input
      type="text"
      id="firstName"
      name="firstName"
      bind:value={$form.firstName}
      {...$constraints.firstName}
    />
    {#if $errors.firstName}
      <p class="error-firstName">{$errors.firstName}</p>
    {/if}
  </div>

  <div class="input-group">
    <label for="lastName">Nachname</label>
    <input
      type="text"
      id="lastName"
      name="lastName"
      bind:value={$form.lastName}
      {...$constraints.lastName}
    />
    {#if $errors.lastName}
      <p class="error-lastName">{$errors.lastName}</p>
    {/if}
  </div>

  <div class="input-group">
    <label for="birthdate">Geburtsdatum</label>
    <input
      type="date"
      id="birthdate"
      name="birthdate"
      bind:value={$birthdateProxy}
      {...$constraints.birthdate}
    />
    {#if $errors.birthdate}
      <p class="error-birthdate">{$errors.birthdate}</p>
    {/if}
  </div>

  <button disabled={$submitting} type="submit">Person erstellen</button>
</form>

<style>
  form {
    display: flex;
    flex-direction: column;
  }
</style>
