import { createConnectTransport } from "@connectrpc/connect-web";

export function defaultTransport(globalFetch: typeof fetch) {
  return createConnectTransport({
    baseUrl: "/api",
    fetch: globalFetch,
  });
}
