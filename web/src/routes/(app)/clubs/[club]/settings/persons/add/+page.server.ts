import type { Actions, PageServerLoad } from "./$types";
import { z } from "zod";
import { superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { fail } from "@sveltejs/kit";
import { createClient } from "@connectrpc/connect";
import { defaultTransport } from "$lib/client";
import { GrpcMutationHandler } from "$lib/grpcMutationHandler";
import type { CreatePersonResponse } from "$lib/gen/soccerbuddy/person/v1/person_service_pb";
import {
  CreatePersonRequestSchema,
  CreatePersonResponseSchema,
  PersonService,
} from "$lib/gen/soccerbuddy/person/v1/person_service_pb";
import { create, toJson } from "@bufbuild/protobuf";
import { timestampFromDate } from "@bufbuild/protobuf/wkt";

const MIN_DATE = new Date(1900, 0, 1);
const MAX_DATE = new Date();

const personSchema = z.object({
  firstName: z.string().min(1),
  lastName: z.string().min(1),
  birthdate: z.date().min(MIN_DATE).max(MAX_DATE),
  owningClubId: z.string().min(1),
});

export const load: PageServerLoad = async ({ parent }) => {
  const form = await superValidate(zod(personSchema));
  const data = await parent();
  return {
    ...data,
    form,
  };
};

export const actions = {
  add: async ({ request, fetch }) => {
    const form = await superValidate(request, zod(personSchema));
    if (!form.valid) {
      return fail(400, { form });
    }

    const client = createClient(PersonService, defaultTransport(fetch));
    const createPersonRequest = create(CreatePersonRequestSchema, {
      firstName: form.data.firstName,
      lastName: form.data.lastName,
      birthdate: timestampFromDate(form.data.birthdate),
      owningClubId: form.data.owningClubId,
    });

    return await GrpcMutationHandler.from(async () => {
      const resp: CreatePersonResponse = await client.createPerson(createPersonRequest);
      return () => ({
        form,
        newPerson: toJson(CreatePersonResponseSchema, resp),
      });
    })
      .onFailedValidation(GrpcMutationHandler.overwriteFormHandler(form))
      .run();
  },
} satisfies Actions;
