import type { PageLoad } from "./$types";
import { runGrpc } from "$lib/runGrpc";
import { createClient } from "@connectrpc/connect";
import { TeamService } from "$lib/gen/soccerbuddy/team/v1/team_service_connect";
import { defaultTransport } from "$lib/client";
import { pbToRole } from "$lib/protobuf";

export const load: PageLoad = async ({ url, parent, fetch }) => {
  const { club, team } = await parent();
  const client = createClient(TeamService, defaultTransport(fetch));
  return runGrpc(url, async () => {
    const members = (await client.listTeamMembers({ teamId: team.id })).members.map((member) => ({
      ...member,
      role: pbToRole(member.role),
    }));

    return {
      club,
      team,
      members,
    };
  });
};
