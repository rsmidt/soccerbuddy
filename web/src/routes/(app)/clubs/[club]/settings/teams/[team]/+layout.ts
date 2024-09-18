import type { LayoutLoad } from "./$types";
import { runGrpc } from "$lib/runGrpc";
import { createClient } from "@connectrpc/connect";
import { TeamService } from "$lib/gen/soccerbuddy/team/v1/team_service_connect";
import { defaultTransport } from "$lib/client";

export const load: LayoutLoad = async ({ params, url, fetch }) => {
  const client = createClient(TeamService, defaultTransport(fetch));
  return runGrpc(url, async () => {
    const team = await client.getTeamOverview({ teamSlug: params.team });
    return {
      team,
    };
  });
};
