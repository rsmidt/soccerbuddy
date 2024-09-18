import type { PageServerLoad } from "./$types";
import { z } from "zod";
import { message, superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { type Actions, fail, redirect } from "@sveltejs/kit";
import { createClient } from "@connectrpc/connect";
import { AccountService } from "$lib/gen/soccerbuddy/account/v1/account_service_connect";
import { defaultTransport } from "$lib/client";
import { LoginRequest } from "$lib/gen/soccerbuddy/account/v1/account_service_pb";
import { GrpcMutationHandler } from "$lib/grpcMutationHandler";

const loginSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
  redirectPath: z.string(),
});

export const load: PageServerLoad = async ({ url }) => {
  const maybeRedirectPath = url.searchParams.get("redirect");
  const redirectPath = maybeRedirectPath ? decodeURIComponent(maybeRedirectPath) : "/clubs";
  return {
    redirectPath,
    form: await superValidate(zod(loginSchema)),
  };
};

export const actions = {
  login: async ({ request, fetch }) => {
    const form = await superValidate(request, zod(loginSchema));
    if (!form.valid) {
      return fail(400, { form });
    }

    const client = createClient(AccountService, defaultTransport(fetch));
    const loginRequest = new LoginRequest({
      email: form.data.email,
      password: form.data.password,
      userAgent: request.headers.get("user-agent") ?? undefined,
    });
    return await GrpcMutationHandler.from(async () => {
      await client.login(loginRequest);
      return () => redirect(302, form.data.redirectPath);
    })
      .onFailedValidation(GrpcMutationHandler.overwriteFormHandler(form))
      .onUnauthenticated(() => message(form, "Ungültige Anmeldedaten", { status: 401 }))
      .run();
  },
} satisfies Actions;
