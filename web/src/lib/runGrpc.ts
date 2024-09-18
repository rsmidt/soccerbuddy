import { Code, ConnectError } from "@connectrpc/connect";
import { error, redirect } from "@sveltejs/kit";

export async function runGrpc<T>(url: URL, runnable: () => Promise<T>): Promise<T>;
export async function runGrpc<T>(location: Location, runnable: () => Promise<T>): Promise<T>;
export async function runGrpc<T>(
  urlOrLocation: URL | Location,
  runnable: () => Promise<T>,
): Promise<T> {
  try {
    return await runnable();
  } catch (e) {
    if (e instanceof ConnectError) {
      switch (e.code) {
        case Code.Unauthenticated:
          const url = urlOrLocation instanceof URL ? urlOrLocation : new URL(urlOrLocation.href);
          redirect(302, "/login?redirect=" + encodeURIComponent(url.pathname));
        case Code.PermissionDenied:
          error(403, "Permission denied");
        case Code.NotFound:
          error(404, "Not found");
      }
    }
    throw e;
  }
}
