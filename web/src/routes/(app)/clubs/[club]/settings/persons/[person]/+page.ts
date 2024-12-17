import type { PageLoad } from "./$types";
import { runGrpc } from "$lib/runGrpc";
import { createClient } from "@connectrpc/connect";
import { defaultTransport } from "$lib/client";
import { PersonService } from "$lib/gen/soccerbuddy/person/v1/person_service_pb";

export const load: PageLoad = async ({ url, parent, fetch, params }) => {
  const { club } = await parent();
  const client = createClient(PersonService, defaultTransport(fetch));
  return runGrpc(url, async () => {
    const person = await client.getPersonOverview({ id: params.person });

    return {
      club,
      person,
    };
  });
};
