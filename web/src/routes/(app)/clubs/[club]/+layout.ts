import type { LayoutLoad } from "./$types";
import { Code, ConnectError, createClient } from "@connectrpc/connect";
import { ClubService } from "$lib/gen/soccerbuddy/club/v1/club_service_connect";
import { error, redirect } from "@sveltejs/kit";
import { defaultTransport } from "$lib/client";
import { runGrpc } from "$lib/runGrpc";

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
