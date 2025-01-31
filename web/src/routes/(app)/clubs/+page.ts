import type { PageLoad } from "./$types";
import { createClient } from "@connectrpc/connect";
import { defaultTransport } from "$lib/client";
import { ClubService } from "$lib/gen/soccerbuddy/club/v1/club_service_pb";
import { runGrpc } from "$lib/runGrpc";

export const load: PageLoad = async ({ fetch, url }) => {
  const client = createClient(ClubService, defaultTransport(fetch));
  return runGrpc(url, async () => {
    const { clubs } = await client.listClubs({});
    return {
      clubs,
    };
  });
};
