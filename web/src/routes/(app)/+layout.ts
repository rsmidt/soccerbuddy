import type { LayoutLoad } from "./$types";
import { createClient } from "@connectrpc/connect";
import { defaultTransport } from "$lib/client";
import { runGrpc } from "$lib/runGrpc";
import { AccountService } from "$lib/gen/soccerbuddy/account/v1/account_service_pb";

export const load: LayoutLoad = async ({ fetch, url }) => {
  const client = createClient(AccountService, defaultTransport(fetch));

  return await runGrpc(url, async () => {
    const resp = await client.getMe({});
    return {
      me: resp,
    };
  });
};
