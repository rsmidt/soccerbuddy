import {
  DescMethodUnary,
  DescService,
  MessageInitShape,
} from "@bufbuild/protobuf";
import { BaseQueryFn } from "@reduxjs/toolkit/dist/query/react";
import { ConnectError, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

type ConnectBaseQueryArgs<
  TClient extends DescService,
  Key extends keyof TClient["method"] = keyof TClient["method"],
> = {
  method: Key;
  req: TClient["method"][Key] extends DescMethodUnary<infer I, any>
    ? MessageInitShape<I>
    : never;
};

type ConnectBaseQueryError = {
  status: number;
  name: string;
  message: string;
};

export function connectBaseQuery<TService extends DescService>(
  clientDef: TService,
): BaseQueryFn<ConnectBaseQueryArgs<TService>, unknown, ConnectBaseQueryError> {
  return async (args, api, extraOptions) => {
    const client = createClient(
      clientDef,
      createConnectTransport({
        baseUrl: process.env.EXPO_PUBLIC_API_URL,
        interceptors: [
          (next) => async (req) => {
            // TODO: How to type state in here without cyclic imports?
            const token = (api.getState() as any).auth.token;
            if (token) {
              req.header.set("Authorization", `Bearer ${token}`);
            }
            return next(req);
          },
        ],
      }),
    );
    try {
      // TODO: Idk how to resolve this type assertion...
      const result = await client[args.method](args.req as any);
      return { data: result };
    } catch (e) {
      const cErr = ConnectError.from(e);

      return {
        error: {
          message: cErr.message,
          status: cErr.code,
          name: cErr.name,
        } satisfies ConnectBaseQueryError,
      };
    }
  };
}
