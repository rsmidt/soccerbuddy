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
  baseQuery: connectBaseQuery<typeof AccountService>(AccountService),
  endpoints: (builder) => ({
    getMe: builder.query<
      MessageShape<typeof GetMeResponseSchema>,
      MessageInitShape<typeof GetMeRequestSchema>
    >({
      query: (req) => ({
        method: "getMe",
        req,
      }),
    }),
  }),
});

export const { useGetMeQuery } = accountApi;
