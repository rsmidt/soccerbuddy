<script lang="ts">
  import type { PageData } from "./$types";
  import { Section } from "$lib/components/section";
  import type { MouseEventHandler } from "svelte/elements";
  import {
    DescribePendingPersonLinkResponse_PersonSchema,
    PersonService,
  } from "$lib/gen/soccerbuddy/person/v1/person_service_pb";
  import { defaultTransport } from "$lib/client";
  import { createClient } from "@connectrpc/connect";
  import { ToastControl } from "$lib/toasts/control.svelte";
  import SimpleMessageToast from "$lib/toasts/variants/SimpleMessageToast.svelte";
  import { AccountLink } from "$lib/gen/soccerbuddy/shared_pb";
  import { fromJson } from "@bufbuild/protobuf";

  const { data }: { data: PageData } = $props();
  const { type, personDescriptor: pdRaw, linkToken } = data;
  const personDescriptor = fromJson(DescribePendingPersonLinkResponse_PersonSchema, pdRaw!);

  const toastControl = ToastControl.getGlobal();
  const client = createClient(PersonService, defaultTransport(fetch));
  const handleInviteClick: MouseEventHandler<HTMLButtonElement> = async () => {
    await client.claimPersonLink({ linkToken });
    toastControl.add({
      component: SimpleMessageToast,
      props: {
        message: "Person erfolgreich verknüpft.",
      }
    });
  };
</script>

<h1 class="default-page-header">Neue Verknüpfung</h1>
<p class="description">Du wurdest eingeladen, dich mit einem neuen Profil zu verknüpfen.</p>

{#if type === "unauthenticated"}
  <Section>
    {#snippet body()}
      <p>
        Du hast noch keinen Account bei SoccerBuddy? <a href="/signup">Registrieren</a>
      </p>
      <p>
        Du hast bereits einen Account, dann melde dich jetzt an. <a
        href="/login?redirect={data.redirect}">Anmelden</a
      >
      </p>
    {/snippet}
  </Section>
{:else}
  {#snippet dataRow(key: string, value: string)}
    <div class="row">
      <div class="item">
        <b>{key}:</b>
        {value}
      </div>
    </div>
  {/snippet}

  <Section>
    {#snippet title()}
      Profil
    {/snippet}
    <div class="list">
      {@render dataRow("Name", personDescriptor.fullName)}
      {#if personDescriptor.linkAs === AccountLink.LINKED_AS_PARENT}
        {@render dataRow("Verknüpfen als", "Betreuer")}
      {/if}
      {@render dataRow("Verein", personDescriptor.clubName)}
      {@render dataRow("Eingeladen von", personDescriptor.invitedBy)}
    </div>
    {#snippet body()}
      <div class="actions">
        <button class="success" onclick={handleInviteClick}>Verknüpfen</button>
        <button class="neutral">Zurück</button>
      </div>
    {/snippet}
  </Section>
{/if}

<style>
    .actions {
        display: flex;
        justify-content: space-between;
        gap: 1rem;

        button {
            flex: 1 1 0;
        }
    }
</style>
