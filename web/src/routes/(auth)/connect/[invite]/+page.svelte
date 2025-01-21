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
  const { personDescriptor: pdRaw, linkToken, linked } = data;
  const personDescriptor = pdRaw ? fromJson(DescribePendingPersonLinkResponse_PersonSchema, pdRaw!) : undefined;

  let success = $state(false);
  const toastControl = ToastControl.getGlobal();
  const client = createClient(PersonService, defaultTransport(fetch));
  const handleInviteClick: MouseEventHandler<HTMLButtonElement> = async () => {
    await client.claimPersonLink({ linkToken });
    toastControl.add({
      component: SimpleMessageToast,
      props: {
        message: "Person erfolgreich verknüpft.",
      },
    });
    success = true;
  };
</script>

{#if linked}
  <p>Du bist bereits mit dieser Person verknüpft. Melde dich jetzt einfach in der App an.</p>
{:else if success}
  <p>Verknüpfung erfolgreich, du kannst dich jetzt in der App anmelden!</p>
{:else}
  <h1 class="default-page-header">Neue Verknüpfung</h1>
  <p class="description">Du wurdest eingeladen, dich mit einem neuen Profil zu verknüpfen.</p>

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
