import type { PageServerLoad } from "./$types";
import { createClient } from "@connectrpc/connect";
import { defaultTransport } from "$lib/client";
import { create } from "@bufbuild/protobuf";
import { type Actions, fail, redirect } from "@sveltejs/kit";
import { z } from "zod";
import { message, superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { GrpcMutationHandler } from "$lib/grpcMutationHandler";
import {
  AccountService,
  RegisterAccountRequestSchema,
} from "$lib/gen/soccerbuddy/account/v1/account_service_pb";

const registerSchema = z
  .object({
    email: z.string().email(),
    password: z.string().min(8),
    passwordConfirmation: z.string().min(8),
    firstName: z.string().min(1),
    lastName: z.string().min(1),
    inviteCode: z.string(),
  })
  .superRefine((data, ctx) => {
    if (data.password && data.passwordConfirmation && data.password !== data.passwordConfirmation) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Passwörter stimmen nicht überein",
        path: ["password"],
      });
    }
  });

export const load: PageServerLoad = async ({ url }) => {
  const inviteCode = url.searchParams.get("invite");
  if (!inviteCode) {
    redirect(302, "/");
  }

  return {
    inviteCode,
    form: await superValidate(zod(registerSchema)),
  };
};

export const actions = {
  register: async ({ request, fetch }) => {
    const form = await superValidate(request, zod(registerSchema));
    if (!form.valid) {
      return fail(400, { form });
    }

    const client = createClient(AccountService, defaultTransport(fetch));
    const req = create(RegisterAccountRequestSchema, {
      email: form.data.email,
      password: form.data.password,
      firstName: form.data.firstName,
      lastName: form.data.lastName,
      linkToken: form.data.inviteCode,
    });
    return await GrpcMutationHandler.from(async () => {
      await client.registerAccount(req);
      return () => ({ form });
    })
      .onFailedValidation(GrpcMutationHandler.overwriteFormHandler(form))
      .onUnauthenticated(() => message(form, "Ungültige Anmeldedaten", { status: 401 }))
      .run();
  },
} satisfies Actions;
