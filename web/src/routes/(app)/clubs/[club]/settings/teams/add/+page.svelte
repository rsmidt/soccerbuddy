<script lang="ts">
  import { superForm } from "sveltekit-superforms";
  import type { PageData } from "./$types";
  import { ToastControl } from "$lib/toasts/control.svelte";
  import EntityCreatedToast from "$lib/toasts/variants/EntityCreatedToast.svelte";
  import NoBreak from "$lib/components/NoBreak.svelte";

  const { data }: { data: PageData } = $props();

  const toastControl = ToastControl.getGlobal();
  const { form, constraints, enhance, errors, submitting } = superForm(data.form, {
    delayMs: 250,
    onUpdate: ({ form, result }) => {
      if (form.valid && result.type === "success") {
        toastControl.add({
          component: EntityCreatedToast,
          props: {
            entityName: form.data.name,
            link: `/clubs/${data.club.slug}/settings/teams/${result.data.slug}`,
            linkText: "Zur Mannschaft",
            message: "Mannschaft wurde erfolgreich angelegt.",
          },
        });
      }
    },
    onSubmit: ({ formData }) => {
      formData.set("owningClubId", data.club.id);
    },
  });
</script>

<h1 class="default-page-header">Neue Mannschaft anlegen</h1>
<p class="description">
  Lege eine neue Mannschaft innerhalb deines Vereins
  <NoBreak><b>{data.club.name}</b></NoBreak>
  an.
</p>
<form method="post" action="?/add" use:enhance>
  <div class="input-group">
    <label for="name">Name</label>
    <input
      type="text"
      id="name"
      name="name"
      placeholder="z. B. U17 Mädchen"
      bind:value={$form.name}
      {...$constraints.name}
    />
    {#if $errors.name}
      <p class="error-name">{$errors.name}</p>
    {/if}
  </div>
  <button disabled={$submitting} type="submit">Mannschaft anlegen</button>
</form>

<style>
  form {
    display: flex;
    flex-direction: column;
  }
</style>
