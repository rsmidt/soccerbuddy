import type { LayoutLoad } from "./$types";
import { createClient } from "@connectrpc/connect";
import { defaultTransport } from "$lib/client";
import { runGrpc } from "$lib/runGrpc";
import { ClubService } from "$lib/gen/soccerbuddy/club/v1/club_service_pb";

export const load: LayoutLoad = async ({ fetch, params, url }) => {
  const client = createClient(ClubService, defaultTransport(fetch));

  return await runGrpc(url, async () => {
    const resp = await client.getClubBySlug({
      slug: params.club,
    });
    return {
      club: resp,
    };
  });
};
