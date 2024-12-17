import type { Actions, PageServerLoad } from "./$types";
import { defaultTransport } from "$lib/client";
import {
  CreateTeamRequestSchema,
  type CreateTeamResponse,
  TeamService,
} from "$lib/gen/soccerbuddy/team/v1/team_service_pb";
import { createClient } from "@connectrpc/connect";
import { fail } from "@sveltejs/kit";
import { superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { z } from "zod";
import { GrpcMutationHandler } from "$lib/grpcMutationHandler";
import { create } from "@bufbuild/protobuf";

const schema = z.object({
  name: z.string(),
  owningClubId: z.string(),
});

export const load: PageServerLoad = async () => {
  const form = await superValidate(zod(schema));

  return { form };
};

export const actions = {
  add: async (stuff) => {
    const { request, fetch } = stuff;
    const form = await superValidate(request, zod(schema));
    if (!form.valid) {
      return fail(400, { form });
    }

    const client = createClient(TeamService, defaultTransport(fetch));
    const createTeamRequest = create(CreateTeamRequestSchema, {
      name: form.data.name,
      owningClubId: form.data.owningClubId,
    });

    return await GrpcMutationHandler.from(async () => {
      const resp: CreateTeamResponse = await client.createTeam(createTeamRequest);
      return () => ({
        form,
        id: resp.id,
        name: resp.name,
        slug: resp.slug,
      });
    })
      .onFailedValidation(GrpcMutationHandler.overwriteFormHandler(form))
      .run();
  },
} satisfies Actions;
