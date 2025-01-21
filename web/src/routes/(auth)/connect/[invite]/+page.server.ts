import type { PageServerLoad } from "./$types";
import { Code, ConnectError, createClient } from "@connectrpc/connect";
import {
  DescribePendingPersonLinkResponse_PersonSchema,
  PersonService,
} from "$lib/gen/soccerbuddy/person/v1/person_service_pb";
import { defaultTransport } from "$lib/client";
import { invariant } from "$lib/invariant";
import { toJson } from "@bufbuild/protobuf";
import {
  AccountAlreadyLinkedToPersonSchema,
  LoginOrAccountCreationRequiredResponseSchema,
} from "$lib/gen/soccerbuddy/shared_pb";
import { redirect } from "@sveltejs/kit";

export const load: PageServerLoad = async ({ params, fetch }) => {
  const client = createClient(PersonService, defaultTransport(fetch));

  try {
    const { person } = await client.describePendingPersonLink({ linkToken: params.invite });
    invariant(person, "person must be defined");

    return {
      linked: false,
      personDescriptor: toJson(DescribePendingPersonLinkResponse_PersonSchema, person),
      linkToken: params.invite,
    };
  } catch (e) {
    const cErr = ConnectError.from(e);
    switch (cErr.code) {
      case Code.FailedPrecondition: {
        const isLoginOrAccountCreationRequired =
          cErr.findDetails(LoginOrAccountCreationRequiredResponseSchema).length > 0;
        if (isLoginOrAccountCreationRequired) {
          redirect(307, `/signup?invite=${params.invite}`);
        }
        const isAccountAlreadyLinked =
          cErr.findDetails(AccountAlreadyLinkedToPersonSchema).length > 0;
        if (isAccountAlreadyLinked) {
          return {
            linked: true,
          };
        }
      }
    }
    throw cErr;
  }
};
