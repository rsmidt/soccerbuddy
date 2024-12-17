import { createApi } from "@reduxjs/toolkit/query/react";
import {
  AccountService,
  GetMeRequestSchema,
  GetMeResponseSchema,
} from "@/api/soccerbuddy/account/v1/account_service_pb";
import { MessageInitShape, MessageShape } from "@bufbuild/protobuf";
import { connectBaseQuery } from "../connect-base-query";

export const accountApi = createApi({
  reducerPath: "accountApi",
  tagTypes: ["account"],
  baseQuery: connectBaseQuery<typeof AccountService>(AccountService),
  refetchOnMountOrArgChange: true,
  refetchOnFocus: true,
  refetchOnReconnect: true,
  endpoints: (builder) => ({
    getMe: builder.query<
      MessageShape<typeof GetMeResponseSchema>,
      MessageInitShape<typeof GetMeRequestSchema>
    >({
      providesTags: [{ id: "me", type: "account" }],
      query: (req) => ({
        method: "getMe",
        req,
      }),
    }),
  }),
});

export const { useGetMeQuery } = accountApi;
