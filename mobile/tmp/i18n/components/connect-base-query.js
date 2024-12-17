import { ConnectError, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
export function connectBaseQuery(clientDef) {
    return async (args, api, extraOptions) => {
        const client = createClient(clientDef, createConnectTransport({
            baseUrl: process.env.EXPO_PUBLIC_API_URL,
            interceptors: [
                (next) => async (req) => {
                    // TODO: How to type state in here without cyclic imports?
                    const token = api.getState().auth.token;
                    if (token) {
                        req.header.set("Authorization", `Bearer ${token}`);
                    }
                    return next(req);
                },
            ],
        }));
        try {
            // TODO: Idk how to resolve this type assertion...
            const result = await client[args.method](args.req);
            return { data: result };
        }
        catch (e) {
            const cErr = ConnectError.from(e);
            return {
                error: {
                    message: cErr.message,
                    status: cErr.code,
                    name: cErr.name,
                },
            };
        }
    };
}
