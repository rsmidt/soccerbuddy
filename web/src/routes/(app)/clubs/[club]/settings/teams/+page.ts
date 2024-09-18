import type { PageLoad } from "./$types";
import { runGrpc } from "$lib/runGrpc";
import { createClient } from "@connectrpc/connect";
import { TeamService } from "$lib/gen/soccerbuddy/team/v1/team_service_connect";
import { defaultTransport } from "$lib/client";

export const load: PageLoad = async ({ parent, url, fetch }) => {
  const { club } = await parent();

  const client = createClient(TeamService, defaultTransport(fetch));
  return runGrpc(url, async () => {
    const { teams } = await client.listTeams({ owningClubId: club.id });
    return {
      club,
      teams: teams.sort((a, b) => a.name.localeCompare(b.name)),
    };
  });
};
