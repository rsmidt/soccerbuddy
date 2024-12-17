import { createApi } from "@reduxjs/toolkit/query/react";
import { AccountService, } from "@/api/soccerbuddy/account/v1/account_service_pb";
import { connectBaseQuery } from "../connect-base-query";
export const accountApi = createApi({
    reducerPath: "accountApi",
    tagTypes: ["account"],
    baseQuery: connectBaseQuery(AccountService),
    refetchOnMountOrArgChange: true,
    refetchOnFocus: true,
    refetchOnReconnect: true,
    endpoints: (builder) => ({
        getMe: builder.query({
            providesTags: [{ id: "me", type: "account" }],
            query: (req) => ({
                method: "getMe",
                req,
            }),
        }),
    }),
});
export const { useGetMeQuery } = accountApi;
