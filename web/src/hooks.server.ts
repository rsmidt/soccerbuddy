import type { Handle, HandleFetch } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
  const initToken = event.url.searchParams.get("initToken");

  const response = await resolve(event, {
    filterSerializedResponseHeaders: (name) => name.startsWith("content-type"),
  });

  // When opening from the app, we cannot set the token beforehand.
  // So what we do is to set the token as a cookie on the response.
  // This should be safe as the cookie value is always validated backend side.
  if (initToken && !event.url.pathname.startsWith("/api")) {
    console.log("initToken", initToken, event.url.pathname);
    // Set token as a cookie on the response.
    response.headers.append("set-cookie", `ID=${initToken}; Path=/; HttpOnly; SameSite=Strict`);
  }

  return response;
};

export const handleFetch: HandleFetch = async ({ request, fetch, event }) => {
  const initToken = event.url.searchParams.get("initToken");

  if (initToken) {
    request.headers.set("Authorization", `Bearer ${initToken}`);
  }

  return fetch(request);
};
