import { createApi } from "@reduxjs/toolkit/query/react";
import { connectBaseQuery } from "../connect-base-query";
import { MessageInitShape, MessageShape } from "@bufbuild/protobuf";
import {
  GetMyTeamHomeRequestSchema,
  GetMyTeamHomeResponseSchema,
  ScheduleTrainingRequestSchema,
  ScheduleTrainingResponseSchema,
  TeamService,
} from "@/api/soccerbuddy/team/v1/team_service_pb";

export const teamApi = createApi({
  reducerPath: "teamApi",
  tagTypes: ["team"],
  baseQuery: connectBaseQuery<typeof TeamService>(TeamService),
  refetchOnMountOrArgChange: true,
  refetchOnFocus: true,
  refetchOnReconnect: true,
  endpoints: (builder) => ({
    scheduleTraining: builder.mutation<
      MessageShape<typeof ScheduleTrainingResponseSchema>,
      MessageInitShape<typeof ScheduleTrainingRequestSchema>
    >({
      query: (req) => ({
        method: "scheduleTraining",
        req,
      }),
    }),
    getMyTeamHome: builder.query<
      MessageShape<typeof GetMyTeamHomeResponseSchema>,
      MessageInitShape<typeof GetMyTeamHomeRequestSchema>
    >({
      providesTags: (result) =>
        result ? [{ type: "team", id: result.teamId }] : [],
      query: (req) => ({
        method: "getMyTeamHome",
        req,
      }),
    }),
  }),
});

export const { useScheduleTrainingMutation, useGetMyTeamHomeQuery } = teamApi;
