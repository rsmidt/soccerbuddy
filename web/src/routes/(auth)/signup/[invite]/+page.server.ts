import type { PageServerLoad } from "./$types";
import { Code, ConnectError, createClient } from "@connectrpc/connect";
import {
  DescribePendingPersonLinkResponse_PersonSchema,
  PersonService,
} from "$lib/gen/soccerbuddy/person/v1/person_service_pb";
import { defaultTransport } from "$lib/client";
import { invariant } from "$lib/invariant";
import { toJson } from "@bufbuild/protobuf";

export const load: PageServerLoad = async ({ params, fetch, url }) => {
  const client = createClient(PersonService, defaultTransport(fetch));

  try {
    const { person } = await client.describePendingPersonLink({ linkToken: params.invite });
    invariant(person, "person must be defined");

    return {
      type: "authenticated" as const,
      personDescriptor: toJson(DescribePendingPersonLinkResponse_PersonSchema, person),
      linkToken: params.invite,
    };
  } catch (e) {
    const cErr = ConnectError.from(e);
    switch (cErr.code) {
      case Code.Unauthenticated:
        // User is unauthenticated, give the option to sign in or signup.
        return {
          type: "unauthenticated" as const,
          redirect: encodeURIComponent(url.pathname),
        };
    }
    throw cErr;
  }
};
