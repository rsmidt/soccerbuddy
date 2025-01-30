import { createApi } from "@reduxjs/toolkit/query/react";
import { connectBaseQuery } from "../connect-base-query";
import { MessageInitShape, MessageShape } from "@bufbuild/protobuf";
import {
  GetMyTeamHomeRequestSchema,
  GetMyTeamHomeResponseSchema,
  ListTeamMembersRequestSchema,
  ListTeamMembersResponse,
  ListTeamMembersResponse_Member,
  ListTeamMembersResponseSchema,
  ScheduleTrainingRequestSchema,
  ScheduleTrainingResponseSchema,
  TeamService,
} from "@/api/soccerbuddy/team/v1/team_service_pb";
import {
  GetMeResponse,
  GetMeResponse_LinkedPerson,
} from "@/api/soccerbuddy/account/v1/account_service_pb";
import { AccountLink } from "@/api/soccerbuddy/shared_pb";
import { PLAYER_ROLE_NAME } from "@/components/constants";
import { createSelector } from "reselect";
import { NominationClass } from "@/components/training/schedule-training-slice";

export const teamApi = createApi({
  reducerPath: "teamApi",
  tagTypes: ["team", "team-members"],
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
    listTeamMembers: builder.query<
      MessageShape<typeof ListTeamMembersResponseSchema>,
      MessageInitShape<typeof ListTeamMembersRequestSchema>
    >({
      providesTags: (result) =>
        result?.members.map((member) => ({
          type: "team-members" as const,
          id: member.id,
        })) ?? [],
      query: (req) => ({
        method: "listTeamMembers",
        req,
      }),
    }),
  }),
});

export const {
  useScheduleTrainingMutation,
  useGetMyTeamHomeQuery,
  useListTeamMembersQuery,
} = teamApi;

/**
 * Selects all linked persons that have a team membership for the given team.
 */
export function selectPersonsInTeam(
  data: GetMeResponse | undefined,
  teamId: string,
): GetMeResponse_LinkedPerson[] {
  if (!data) return [];

  return data.linkedPersons.filter((person) =>
    person.teamMemberships.some((team) => team.id === teamId),
  );
}

/**
 * Selects any person wih a parent link ONLY when there's no person with a self link.
 * We do this because we assume that parents often do not really know about the names of their children teams.
 */
export function selectPersonsWithParentLink(
  data: GetMeResponse | undefined,
  teamId: string,
): GetMeResponse_LinkedPerson | undefined {
  if (!data) return undefined;

  const personsInTeam = selectPersonsInTeam(data, teamId);
  const personsLinkedWithParent = personsInTeam.filter(
    (person) => person.linkedAs === AccountLink.LINKED_AS_PARENT,
  );
  const hasPersonsWithSelfLink = personsInTeam.some(
    (person) => person.linkedAs === AccountLink.LINKED_AS_SELF,
  );
  if (hasPersonsWithSelfLink || personsLinkedWithParent.length === 0) {
    return undefined;
  }
  return personsLinkedWithParent[0];
}

/**
 * Selects if any linked person with a team membership for the given team has permission to edit.
 */
export function selectHasEditAllowance(
  data: GetMeResponse | undefined,
  teamId: string,
): boolean {
  if (!data) return false;

  return selectPersonsInTeam(data, teamId).some((person) =>
    person.teamMemberships.some(
      (team) => team.id === teamId && team.role === "COACH",
    ),
  );
}

/**
 * Selects the team members for a given mode.
 */
export const selectTeamMembersByMode = createSelector(
  [
    (data: ListTeamMembersResponse | undefined) => data?.members,
    (data: ListTeamMembersResponse | undefined, mode: NominationClass) => mode,
  ],
  (members, mode) =>
    members
      ?.filter(getFilterCriteria(mode))
      .sort((a, b) => (a.firstName > b.firstName ? 1 : -1)),
);

/**
 * Selects the initials of a team member.
 */
export function selectTeamMemberInitials(
  member: ListTeamMembersResponse_Member | undefined,
): string {
  if (!member) return "";
  return member.firstName[0] + member.lastName[0];
}

/**
 * Selects the full name of a team member.
 */
export function selectTeamMemberName(
  member: ListTeamMembersResponse_Member | undefined,
): string {
  if (!member) return "";
  return `${member.firstName} ${member.lastName}`;
}

/**
 * Returns the filter criteria for a given mode.
 */
export function getFilterCriteria(
  mode: NominationClass,
): (member: ListTeamMembersResponse_Member) => boolean {
  return (member) => {
    if (mode === "player") {
      return member.role === PLAYER_ROLE_NAME;
    } else {
      return member.role !== PLAYER_ROLE_NAME;
    }
  };
}
