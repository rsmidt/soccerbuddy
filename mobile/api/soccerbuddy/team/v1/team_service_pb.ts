// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file soccerbuddy/team/v1/team_service.proto (package soccerbuddy.team.v1, syntax proto3)
/* eslint-disable */

import type { GenEnum, GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { enumDesc, fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { DateTime } from "../../../google/type/datetime_pb";
import { file_google_type_datetime } from "../../../google/type/datetime_pb";
import type { RatingPolicy } from "../../shared_pb";
import { file_soccerbuddy_shared } from "../../shared_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file soccerbuddy/team/v1/team_service.proto.
 */
export const file_soccerbuddy_team_v1_team_service: GenFile = /*@__PURE__*/
  fileDesc("CiZzb2NjZXJidWRkeS90ZWFtL3YxL3RlYW1fc2VydmljZS5wcm90bxITc29jY2VyYnVkZHkudGVhbS52MSI5ChFDcmVhdGVUZWFtUmVxdWVzdBIMCgRuYW1lGAEgASgJEhYKDm93bmluZ19jbHViX2lkGAIgASgJIrQBChJDcmVhdGVUZWFtUmVzcG9uc2USCgoCaWQYASABKAkSDAoEbmFtZRgCIAEoCRIMCgRzbHVnGAMgASgJEhYKDm93bmluZ19jbHViX2lkGAQgASgJEi4KCmNyZWF0ZWRfYXQYBSABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wEi4KCnVwZGF0ZWRfYXQYBiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wIioKEExpc3RUZWFtc1JlcXVlc3QSFgoOb3duaW5nX2NsdWJfaWQYASABKAki4AEKEUxpc3RUZWFtc1Jlc3BvbnNlEjoKBXRlYW1zGAEgAygLMisuc29jY2VyYnVkZHkudGVhbS52MS5MaXN0VGVhbXNSZXNwb25zZS5UZWFtGo4BCgRUZWFtEgoKAmlkGAEgASgJEgwKBG5hbWUYAiABKAkSDAoEc2x1ZxgDIAEoCRIuCgpjcmVhdGVkX2F0GAQgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcBIuCgp1cGRhdGVkX2F0GAUgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcCIrChZHZXRUZWFtT3ZlcnZpZXdSZXF1ZXN0EhEKCXRlYW1fc2x1ZxgBIAEoCSK5AQoXR2V0VGVhbU92ZXJ2aWV3UmVzcG9uc2USCgoCaWQYASABKAkSDAoEbmFtZRgCIAEoCRIMCgRzbHVnGAMgASgJEhYKDm93bmluZ19jbHViX2lkGAQgASgJEi4KCmNyZWF0ZWRfYXQYBSABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wEi4KCnVwZGF0ZWRfYXQYBiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wIkoKFkFkZFBlcnNvblRvVGVhbVJlcXVlc3QSDwoHdGVhbV9pZBgBIAEoCRIRCglwZXJzb25faWQYAiABKAkSDAoEcm9sZRgDIAEoCSIZChdBZGRQZXJzb25Ub1RlYW1SZXNwb25zZSIkChFEZWxldGVUZWFtUmVxdWVzdBIPCgd0ZWFtX2lkGAEgASgJIhQKEkRlbGV0ZVRlYW1SZXNwb25zZSI/Ch1TZWFyY2hQZXJzb25zTm90SW5UZWFtUmVxdWVzdBIPCgd0ZWFtX2lkGAEgASgJEg0KBXF1ZXJ5GAIgASgJIqoBCh5TZWFyY2hQZXJzb25zTm90SW5UZWFtUmVzcG9uc2USSwoHcGVyc29ucxgBIAMoCzI6LnNvY2NlcmJ1ZGR5LnRlYW0udjEuU2VhcmNoUGVyc29uc05vdEluVGVhbVJlc3BvbnNlLlBlcnNvbho7CgZQZXJzb24SCgoCaWQYASABKAkSEgoKZmlyc3RfbmFtZRgCIAEoCRIRCglsYXN0X25hbWUYAyABKAkiKQoWTGlzdFRlYW1NZW1iZXJzUmVxdWVzdBIPCgd0ZWFtX2lkGAEgASgJIpUCChdMaXN0VGVhbU1lbWJlcnNSZXNwb25zZRJECgdtZW1iZXJzGAEgAygLMjMuc29jY2VyYnVkZHkudGVhbS52MS5MaXN0VGVhbU1lbWJlcnNSZXNwb25zZS5NZW1iZXIaswEKBk1lbWJlchIKCgJpZBgBIAEoCRISCgpmaXJzdF9uYW1lGAIgASgJEhEKCWxhc3RfbmFtZRgDIAEoCRIRCglwZXJzb25faWQYBCABKAkSFwoKaW52aXRlcl9pZBgFIAEoCUgAiAEBEgwKBHJvbGUYBiABKAkSLQoJam9pbmVkX2F0GAcgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcEINCgtfaW52aXRlcl9pZCLfBAoXU2NoZWR1bGVUcmFpbmluZ1JlcXVlc3QSDwoHdGVhbV9pZBgBIAEoCRIrCgxzY2hlZHVsZWRfYXQYAiABKAsyFS5nb29nbGUudHlwZS5EYXRlVGltZRImCgdlbmRzX2F0GAMgASgLMhUuZ29vZ2xlLnR5cGUuRGF0ZVRpbWUSFQoIbG9jYXRpb24YBCABKAlIAIgBARIXCgpmaWVsZF90eXBlGAUgASgJSAGIAQESQQoPZ2F0aGVyaW5nX3BvaW50GAYgASgLMiMuc29jY2VyYnVkZHkudGVhbS52MS5HYXRoZXJpbmdQb2ludEgCiAEBElIKF2Fja25vd2xlZGdtZW50X3NldHRpbmdzGAcgASgLMiwuc29jY2VyYnVkZHkudGVhbS52MS5BY2tub3dsZWRnZW1lbnRTZXR0aW5nc0gDiAEBEhgKC2Rlc2NyaXB0aW9uGAggASgJSASIAQESQQoPcmF0aW5nX3NldHRpbmdzGAkgASgLMiMuc29jY2VyYnVkZHkudGVhbS52MS5SYXRpbmdTZXR0aW5nc0gFiAEBEjoKC25vbWluYXRpb25zGAogASgLMiAuc29jY2VyYnVkZHkudGVhbS52MS5Ob21pbmF0aW9uc0gGiAEBQgsKCV9sb2NhdGlvbkINCgtfZmllbGRfdHlwZUISChBfZ2F0aGVyaW5nX3BvaW50QhoKGF9hY2tub3dsZWRnbWVudF9zZXR0aW5nc0IOCgxfZGVzY3JpcHRpb25CEgoQX3JhdGluZ19zZXR0aW5nc0IOCgxfbm9taW5hdGlvbnMiUgoOR2F0aGVyaW5nUG9pbnQSEAoIbG9jYXRpb24YBSABKAkSLgoPZ2F0aGVyaW5nX3VudGlsGAYgASgLMhUuZ29vZ2xlLnR5cGUuRGF0ZVRpbWUiQgoXQWNrbm93bGVkZ2VtZW50U2V0dGluZ3MSJwoIZGVhZGxpbmUYASABKAsyFS5nb29nbGUudHlwZS5EYXRlVGltZSJCCg5SYXRpbmdTZXR0aW5ncxIwCgZwb2xpY3kYASABKA4yIC5zb2NjZXJidWRkeS5zaGFyZWQuUmF0aW5nUG9saWN5IoMCCgtOb21pbmF0aW9ucxISCgpwbGF5ZXJfaWRzGAEgAygJEhEKCXN0YWZmX2lkcxgCIAMoCRJQChNub3RpZmljYXRpb25fcG9saWN5GAMgASgOMjMuc29jY2VyYnVkZHkudGVhbS52MS5Ob21pbmF0aW9ucy5Ob3RpZmljYXRpb25Qb2xpY3kiewoSTm90aWZpY2F0aW9uUG9saWN5EiMKH05PVElGSUNBVElPTl9QT0xJQ1lfVU5TUEVDSUZJRUQQABIeChpOT1RJRklDQVRJT05fUE9MSUNZX1NJTEVOVBABEiAKHE5PVElGSUNBVElPTl9QT0xJQ1lfUkVRVUlSRUQQAiIaChhTY2hlZHVsZVRyYWluaW5nUmVzcG9uc2UiJwoUR2V0TXlUZWFtSG9tZVJlcXVlc3QSDwoHdGVhbV9pZBgBIAEoCSKFBQoVR2V0TXlUZWFtSG9tZVJlc3BvbnNlEg8KB3RlYW1faWQYASABKAkSEQoJdGVhbV9uYW1lGAIgASgJEkYKCXRyYWluaW5ncxgDIAMoCzIzLnNvY2NlcmJ1ZGR5LnRlYW0udjEuR2V0TXlUZWFtSG9tZVJlc3BvbnNlLlRyYWluaW5nGv8DCghUcmFpbmluZxIKCgJpZBgBIAEoCRIrCgxzY2hlZHVsZWRfYXQYAiABKAsyFS5nb29nbGUudHlwZS5EYXRlVGltZRImCgdlbmRzX2F0GAMgASgLMhUuZ29vZ2xlLnR5cGUuRGF0ZVRpbWUSFQoIbG9jYXRpb24YBCABKAlIAIgBARIXCgpmaWVsZF90eXBlGAUgASgJSAGIAQESQQoPZ2F0aGVyaW5nX3BvaW50GAYgASgLMiMuc29jY2VyYnVkZHkudGVhbS52MS5HYXRoZXJpbmdQb2ludEgCiAEBElIKF2Fja25vd2xlZGdtZW50X3NldHRpbmdzGAcgASgLMiwuc29jY2VyYnVkZHkudGVhbS52MS5BY2tub3dsZWRnZW1lbnRTZXR0aW5nc0gDiAEBEhgKC2Rlc2NyaXB0aW9uGAggASgJSASIAQESQQoPcmF0aW5nX3NldHRpbmdzGAkgASgLMiMuc29jY2VyYnVkZHkudGVhbS52MS5SYXRpbmdTZXR0aW5nc0gFiAEBQgsKCV9sb2NhdGlvbkINCgtfZmllbGRfdHlwZUISChBfZ2F0aGVyaW5nX3BvaW50QhoKGF9hY2tub3dsZWRnbWVudF9zZXR0aW5nc0IOCgxfZGVzY3JpcHRpb25CEgoQX3JhdGluZ19zZXR0aW5ncyJfCiFOb21pbmF0ZVBlcnNvbnNGb3JUcmFpbmluZ1JlcXVlc3QSEwoLdHJhaW5pbmdfaWQYASABKAkSEgoKcGxheWVyX2lkcxgCIAMoCRIRCglzdGFmZl9pZHMYAyADKAkiJAoiTm9taW5hdGVQZXJzb25zRm9yVHJhaW5pbmdSZXNwb25zZTLyCAoLVGVhbVNlcnZpY2USXwoKQ3JlYXRlVGVhbRImLnNvY2NlcmJ1ZGR5LnRlYW0udjEuQ3JlYXRlVGVhbVJlcXVlc3QaJy5zb2NjZXJidWRkeS50ZWFtLnYxLkNyZWF0ZVRlYW1SZXNwb25zZSIAElwKCUxpc3RUZWFtcxIlLnNvY2NlcmJ1ZGR5LnRlYW0udjEuTGlzdFRlYW1zUmVxdWVzdBomLnNvY2NlcmJ1ZGR5LnRlYW0udjEuTGlzdFRlYW1zUmVzcG9uc2UiABJfCgpEZWxldGVUZWFtEiYuc29jY2VyYnVkZHkudGVhbS52MS5EZWxldGVUZWFtUmVxdWVzdBonLnNvY2NlcmJ1ZGR5LnRlYW0udjEuRGVsZXRlVGVhbVJlc3BvbnNlIgASbgoPR2V0VGVhbU92ZXJ2aWV3Eisuc29jY2VyYnVkZHkudGVhbS52MS5HZXRUZWFtT3ZlcnZpZXdSZXF1ZXN0Giwuc29jY2VyYnVkZHkudGVhbS52MS5HZXRUZWFtT3ZlcnZpZXdSZXNwb25zZSIAEoMBChZTZWFyY2hQZXJzb25zTm90SW5UZWFtEjIuc29jY2VyYnVkZHkudGVhbS52MS5TZWFyY2hQZXJzb25zTm90SW5UZWFtUmVxdWVzdBozLnNvY2NlcmJ1ZGR5LnRlYW0udjEuU2VhcmNoUGVyc29uc05vdEluVGVhbVJlc3BvbnNlIgASbgoPQWRkUGVyc29uVG9UZWFtEisuc29jY2VyYnVkZHkudGVhbS52MS5BZGRQZXJzb25Ub1RlYW1SZXF1ZXN0Giwuc29jY2VyYnVkZHkudGVhbS52MS5BZGRQZXJzb25Ub1RlYW1SZXNwb25zZSIAEm4KD0xpc3RUZWFtTWVtYmVycxIrLnNvY2NlcmJ1ZGR5LnRlYW0udjEuTGlzdFRlYW1NZW1iZXJzUmVxdWVzdBosLnNvY2NlcmJ1ZGR5LnRlYW0udjEuTGlzdFRlYW1NZW1iZXJzUmVzcG9uc2UiABJxChBTY2hlZHVsZVRyYWluaW5nEiwuc29jY2VyYnVkZHkudGVhbS52MS5TY2hlZHVsZVRyYWluaW5nUmVxdWVzdBotLnNvY2NlcmJ1ZGR5LnRlYW0udjEuU2NoZWR1bGVUcmFpbmluZ1Jlc3BvbnNlIgASaAoNR2V0TXlUZWFtSG9tZRIpLnNvY2NlcmJ1ZGR5LnRlYW0udjEuR2V0TXlUZWFtSG9tZVJlcXVlc3QaKi5zb2NjZXJidWRkeS50ZWFtLnYxLkdldE15VGVhbUhvbWVSZXNwb25zZSIAEo8BChpOb21pbmF0ZVBlcnNvbnNGb3JUcmFpbmluZxI2LnNvY2NlcmJ1ZGR5LnRlYW0udjEuTm9taW5hdGVQZXJzb25zRm9yVHJhaW5pbmdSZXF1ZXN0Gjcuc29jY2VyYnVkZHkudGVhbS52MS5Ob21pbmF0ZVBlcnNvbnNGb3JUcmFpbmluZ1Jlc3BvbnNlIgBC2gEKF2NvbS5zb2NjZXJidWRkeS50ZWFtLnYxQhBUZWFtU2VydmljZVByb3RvUAFaP2dpdGh1Yi5jb20vcnNtaWR0L3NvY2NlcmJ1ZGR5L2dlbi9nby9zb2NjZXJidWRkeS90ZWFtL3YxO3RlYW12MaICA1NUWKoCE1NvY2NlcmJ1ZGR5LlRlYW0uVjHKAhNTb2NjZXJidWRkeVxUZWFtXFYx4gIfU29jY2VyYnVkZHlcVGVhbVxWMVxHUEJNZXRhZGF0YeoCFVNvY2NlcmJ1ZGR5OjpUZWFtOjpWMWIGcHJvdG8z", [file_google_protobuf_timestamp, file_google_type_datetime, file_soccerbuddy_shared]);

/**
 * @generated from message soccerbuddy.team.v1.CreateTeamRequest
 */
export type CreateTeamRequest = Message<"soccerbuddy.team.v1.CreateTeamRequest"> & {
  /**
   * @generated from field: string name = 1;
   */
  name: string;

  /**
   * @generated from field: string owning_club_id = 2;
   */
  owningClubId: string;
};

/**
 * Describes the message soccerbuddy.team.v1.CreateTeamRequest.
 * Use `create(CreateTeamRequestSchema)` to create a new message.
 */
export const CreateTeamRequestSchema: GenMessage<CreateTeamRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 0);

/**
 * @generated from message soccerbuddy.team.v1.CreateTeamResponse
 */
export type CreateTeamResponse = Message<"soccerbuddy.team.v1.CreateTeamResponse"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string name = 2;
   */
  name: string;

  /**
   * @generated from field: string slug = 3;
   */
  slug: string;

  /**
   * @generated from field: string owning_club_id = 4;
   */
  owningClubId: string;

  /**
   * @generated from field: google.protobuf.Timestamp created_at = 5;
   */
  createdAt?: Timestamp;

  /**
   * @generated from field: google.protobuf.Timestamp updated_at = 6;
   */
  updatedAt?: Timestamp;
};

/**
 * Describes the message soccerbuddy.team.v1.CreateTeamResponse.
 * Use `create(CreateTeamResponseSchema)` to create a new message.
 */
export const CreateTeamResponseSchema: GenMessage<CreateTeamResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 1);

/**
 * @generated from message soccerbuddy.team.v1.ListTeamsRequest
 */
export type ListTeamsRequest = Message<"soccerbuddy.team.v1.ListTeamsRequest"> & {
  /**
   * @generated from field: string owning_club_id = 1;
   */
  owningClubId: string;
};

/**
 * Describes the message soccerbuddy.team.v1.ListTeamsRequest.
 * Use `create(ListTeamsRequestSchema)` to create a new message.
 */
export const ListTeamsRequestSchema: GenMessage<ListTeamsRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 2);

/**
 * @generated from message soccerbuddy.team.v1.ListTeamsResponse
 */
export type ListTeamsResponse = Message<"soccerbuddy.team.v1.ListTeamsResponse"> & {
  /**
   * @generated from field: repeated soccerbuddy.team.v1.ListTeamsResponse.Team teams = 1;
   */
  teams: ListTeamsResponse_Team[];
};

/**
 * Describes the message soccerbuddy.team.v1.ListTeamsResponse.
 * Use `create(ListTeamsResponseSchema)` to create a new message.
 */
export const ListTeamsResponseSchema: GenMessage<ListTeamsResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 3);

/**
 * @generated from message soccerbuddy.team.v1.ListTeamsResponse.Team
 */
export type ListTeamsResponse_Team = Message<"soccerbuddy.team.v1.ListTeamsResponse.Team"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string name = 2;
   */
  name: string;

  /**
   * @generated from field: string slug = 3;
   */
  slug: string;

  /**
   * @generated from field: google.protobuf.Timestamp created_at = 4;
   */
  createdAt?: Timestamp;

  /**
   * @generated from field: google.protobuf.Timestamp updated_at = 5;
   */
  updatedAt?: Timestamp;
};

/**
 * Describes the message soccerbuddy.team.v1.ListTeamsResponse.Team.
 * Use `create(ListTeamsResponse_TeamSchema)` to create a new message.
 */
export const ListTeamsResponse_TeamSchema: GenMessage<ListTeamsResponse_Team> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 3, 0);

/**
 * @generated from message soccerbuddy.team.v1.GetTeamOverviewRequest
 */
export type GetTeamOverviewRequest = Message<"soccerbuddy.team.v1.GetTeamOverviewRequest"> & {
  /**
   * @generated from field: string team_slug = 1;
   */
  teamSlug: string;
};

/**
 * Describes the message soccerbuddy.team.v1.GetTeamOverviewRequest.
 * Use `create(GetTeamOverviewRequestSchema)` to create a new message.
 */
export const GetTeamOverviewRequestSchema: GenMessage<GetTeamOverviewRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 4);

/**
 * @generated from message soccerbuddy.team.v1.GetTeamOverviewResponse
 */
export type GetTeamOverviewResponse = Message<"soccerbuddy.team.v1.GetTeamOverviewResponse"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string name = 2;
   */
  name: string;

  /**
   * @generated from field: string slug = 3;
   */
  slug: string;

  /**
   * @generated from field: string owning_club_id = 4;
   */
  owningClubId: string;

  /**
   * @generated from field: google.protobuf.Timestamp created_at = 5;
   */
  createdAt?: Timestamp;

  /**
   * @generated from field: google.protobuf.Timestamp updated_at = 6;
   */
  updatedAt?: Timestamp;
};

/**
 * Describes the message soccerbuddy.team.v1.GetTeamOverviewResponse.
 * Use `create(GetTeamOverviewResponseSchema)` to create a new message.
 */
export const GetTeamOverviewResponseSchema: GenMessage<GetTeamOverviewResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 5);

/**
 * @generated from message soccerbuddy.team.v1.AddPersonToTeamRequest
 */
export type AddPersonToTeamRequest = Message<"soccerbuddy.team.v1.AddPersonToTeamRequest"> & {
  /**
   * @generated from field: string team_id = 1;
   */
  teamId: string;

  /**
   * @generated from field: string person_id = 2;
   */
  personId: string;

  /**
   * @generated from field: string role = 3;
   */
  role: string;
};

/**
 * Describes the message soccerbuddy.team.v1.AddPersonToTeamRequest.
 * Use `create(AddPersonToTeamRequestSchema)` to create a new message.
 */
export const AddPersonToTeamRequestSchema: GenMessage<AddPersonToTeamRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 6);

/**
 * @generated from message soccerbuddy.team.v1.AddPersonToTeamResponse
 */
export type AddPersonToTeamResponse = Message<"soccerbuddy.team.v1.AddPersonToTeamResponse"> & {
};

/**
 * Describes the message soccerbuddy.team.v1.AddPersonToTeamResponse.
 * Use `create(AddPersonToTeamResponseSchema)` to create a new message.
 */
export const AddPersonToTeamResponseSchema: GenMessage<AddPersonToTeamResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 7);

/**
 * @generated from message soccerbuddy.team.v1.DeleteTeamRequest
 */
export type DeleteTeamRequest = Message<"soccerbuddy.team.v1.DeleteTeamRequest"> & {
  /**
   * @generated from field: string team_id = 1;
   */
  teamId: string;
};

/**
 * Describes the message soccerbuddy.team.v1.DeleteTeamRequest.
 * Use `create(DeleteTeamRequestSchema)` to create a new message.
 */
export const DeleteTeamRequestSchema: GenMessage<DeleteTeamRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 8);

/**
 * @generated from message soccerbuddy.team.v1.DeleteTeamResponse
 */
export type DeleteTeamResponse = Message<"soccerbuddy.team.v1.DeleteTeamResponse"> & {
};

/**
 * Describes the message soccerbuddy.team.v1.DeleteTeamResponse.
 * Use `create(DeleteTeamResponseSchema)` to create a new message.
 */
export const DeleteTeamResponseSchema: GenMessage<DeleteTeamResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 9);

/**
 * @generated from message soccerbuddy.team.v1.SearchPersonsNotInTeamRequest
 */
export type SearchPersonsNotInTeamRequest = Message<"soccerbuddy.team.v1.SearchPersonsNotInTeamRequest"> & {
  /**
   * @generated from field: string team_id = 1;
   */
  teamId: string;

  /**
   * @generated from field: string query = 2;
   */
  query: string;
};

/**
 * Describes the message soccerbuddy.team.v1.SearchPersonsNotInTeamRequest.
 * Use `create(SearchPersonsNotInTeamRequestSchema)` to create a new message.
 */
export const SearchPersonsNotInTeamRequestSchema: GenMessage<SearchPersonsNotInTeamRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 10);

/**
 * @generated from message soccerbuddy.team.v1.SearchPersonsNotInTeamResponse
 */
export type SearchPersonsNotInTeamResponse = Message<"soccerbuddy.team.v1.SearchPersonsNotInTeamResponse"> & {
  /**
   * @generated from field: repeated soccerbuddy.team.v1.SearchPersonsNotInTeamResponse.Person persons = 1;
   */
  persons: SearchPersonsNotInTeamResponse_Person[];
};

/**
 * Describes the message soccerbuddy.team.v1.SearchPersonsNotInTeamResponse.
 * Use `create(SearchPersonsNotInTeamResponseSchema)` to create a new message.
 */
export const SearchPersonsNotInTeamResponseSchema: GenMessage<SearchPersonsNotInTeamResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 11);

/**
 * @generated from message soccerbuddy.team.v1.SearchPersonsNotInTeamResponse.Person
 */
export type SearchPersonsNotInTeamResponse_Person = Message<"soccerbuddy.team.v1.SearchPersonsNotInTeamResponse.Person"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string first_name = 2;
   */
  firstName: string;

  /**
   * @generated from field: string last_name = 3;
   */
  lastName: string;
};

/**
 * Describes the message soccerbuddy.team.v1.SearchPersonsNotInTeamResponse.Person.
 * Use `create(SearchPersonsNotInTeamResponse_PersonSchema)` to create a new message.
 */
export const SearchPersonsNotInTeamResponse_PersonSchema: GenMessage<SearchPersonsNotInTeamResponse_Person> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 11, 0);

/**
 * @generated from message soccerbuddy.team.v1.ListTeamMembersRequest
 */
export type ListTeamMembersRequest = Message<"soccerbuddy.team.v1.ListTeamMembersRequest"> & {
  /**
   * @generated from field: string team_id = 1;
   */
  teamId: string;
};

/**
 * Describes the message soccerbuddy.team.v1.ListTeamMembersRequest.
 * Use `create(ListTeamMembersRequestSchema)` to create a new message.
 */
export const ListTeamMembersRequestSchema: GenMessage<ListTeamMembersRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 12);

/**
 * @generated from message soccerbuddy.team.v1.ListTeamMembersResponse
 */
export type ListTeamMembersResponse = Message<"soccerbuddy.team.v1.ListTeamMembersResponse"> & {
  /**
   * @generated from field: repeated soccerbuddy.team.v1.ListTeamMembersResponse.Member members = 1;
   */
  members: ListTeamMembersResponse_Member[];
};

/**
 * Describes the message soccerbuddy.team.v1.ListTeamMembersResponse.
 * Use `create(ListTeamMembersResponseSchema)` to create a new message.
 */
export const ListTeamMembersResponseSchema: GenMessage<ListTeamMembersResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 13);

/**
 * @generated from message soccerbuddy.team.v1.ListTeamMembersResponse.Member
 */
export type ListTeamMembersResponse_Member = Message<"soccerbuddy.team.v1.ListTeamMembersResponse.Member"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string first_name = 2;
   */
  firstName: string;

  /**
   * @generated from field: string last_name = 3;
   */
  lastName: string;

  /**
   * @generated from field: string person_id = 4;
   */
  personId: string;

  /**
   * @generated from field: optional string inviter_id = 5;
   */
  inviterId?: string;

  /**
   * @generated from field: string role = 6;
   */
  role: string;

  /**
   * @generated from field: google.protobuf.Timestamp joined_at = 7;
   */
  joinedAt?: Timestamp;
};

/**
 * Describes the message soccerbuddy.team.v1.ListTeamMembersResponse.Member.
 * Use `create(ListTeamMembersResponse_MemberSchema)` to create a new message.
 */
export const ListTeamMembersResponse_MemberSchema: GenMessage<ListTeamMembersResponse_Member> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 13, 0);

/**
 * @generated from message soccerbuddy.team.v1.ScheduleTrainingRequest
 */
export type ScheduleTrainingRequest = Message<"soccerbuddy.team.v1.ScheduleTrainingRequest"> & {
  /**
   * @generated from field: string team_id = 1;
   */
  teamId: string;

  /**
   * @generated from field: google.type.DateTime scheduled_at = 2;
   */
  scheduledAt?: DateTime;

  /**
   * @generated from field: google.type.DateTime ends_at = 3;
   */
  endsAt?: DateTime;

  /**
   * @generated from field: optional string location = 4;
   */
  location?: string;

  /**
   * The type of field used: e.g. hard floor, lawn, etc.
   *
   * @generated from field: optional string field_type = 5;
   */
  fieldType?: string;

  /**
   * @generated from field: optional soccerbuddy.team.v1.GatheringPoint gathering_point = 6;
   */
  gatheringPoint?: GatheringPoint;

  /**
   * @generated from field: optional soccerbuddy.team.v1.AcknowledgementSettings acknowledgment_settings = 7;
   */
  acknowledgmentSettings?: AcknowledgementSettings;

  /**
   * @generated from field: optional string description = 8;
   */
  description?: string;

  /**
   * @generated from field: optional soccerbuddy.team.v1.RatingSettings rating_settings = 9;
   */
  ratingSettings?: RatingSettings;

  /**
   * @generated from field: optional soccerbuddy.team.v1.Nominations nominations = 10;
   */
  nominations?: Nominations;
};

/**
 * Describes the message soccerbuddy.team.v1.ScheduleTrainingRequest.
 * Use `create(ScheduleTrainingRequestSchema)` to create a new message.
 */
export const ScheduleTrainingRequestSchema: GenMessage<ScheduleTrainingRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 14);

/**
 * @generated from message soccerbuddy.team.v1.GatheringPoint
 */
export type GatheringPoint = Message<"soccerbuddy.team.v1.GatheringPoint"> & {
  /**
   * @generated from field: string location = 5;
   */
  location: string;

  /**
   * @generated from field: google.type.DateTime gathering_until = 6;
   */
  gatheringUntil?: DateTime;
};

/**
 * Describes the message soccerbuddy.team.v1.GatheringPoint.
 * Use `create(GatheringPointSchema)` to create a new message.
 */
export const GatheringPointSchema: GenMessage<GatheringPoint> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 15);

/**
 * @generated from message soccerbuddy.team.v1.AcknowledgementSettings
 */
export type AcknowledgementSettings = Message<"soccerbuddy.team.v1.AcknowledgementSettings"> & {
  /**
   * @generated from field: google.type.DateTime deadline = 1;
   */
  deadline?: DateTime;
};

/**
 * Describes the message soccerbuddy.team.v1.AcknowledgementSettings.
 * Use `create(AcknowledgementSettingsSchema)` to create a new message.
 */
export const AcknowledgementSettingsSchema: GenMessage<AcknowledgementSettings> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 16);

/**
 * @generated from message soccerbuddy.team.v1.RatingSettings
 */
export type RatingSettings = Message<"soccerbuddy.team.v1.RatingSettings"> & {
  /**
   * @generated from field: soccerbuddy.shared.RatingPolicy policy = 1;
   */
  policy: RatingPolicy;
};

/**
 * Describes the message soccerbuddy.team.v1.RatingSettings.
 * Use `create(RatingSettingsSchema)` to create a new message.
 */
export const RatingSettingsSchema: GenMessage<RatingSettings> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 17);

/**
 * @generated from message soccerbuddy.team.v1.Nominations
 */
export type Nominations = Message<"soccerbuddy.team.v1.Nominations"> & {
  /**
   * @generated from field: repeated string player_ids = 1;
   */
  playerIds: string[];

  /**
   * @generated from field: repeated string staff_ids = 2;
   */
  staffIds: string[];

  /**
   * @generated from field: soccerbuddy.team.v1.Nominations.NotificationPolicy notification_policy = 3;
   */
  notificationPolicy: Nominations_NotificationPolicy;
};

/**
 * Describes the message soccerbuddy.team.v1.Nominations.
 * Use `create(NominationsSchema)` to create a new message.
 */
export const NominationsSchema: GenMessage<Nominations> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 18);

/**
 * @generated from enum soccerbuddy.team.v1.Nominations.NotificationPolicy
 */
export enum Nominations_NotificationPolicy {
  /**
   * @generated from enum value: NOTIFICATION_POLICY_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: NOTIFICATION_POLICY_SILENT = 1;
   */
  SILENT = 1,

  /**
   * @generated from enum value: NOTIFICATION_POLICY_REQUIRED = 2;
   */
  REQUIRED = 2,
}

/**
 * Describes the enum soccerbuddy.team.v1.Nominations.NotificationPolicy.
 */
export const Nominations_NotificationPolicySchema: GenEnum<Nominations_NotificationPolicy> = /*@__PURE__*/
  enumDesc(file_soccerbuddy_team_v1_team_service, 18, 0);

/**
 * @generated from message soccerbuddy.team.v1.ScheduleTrainingResponse
 */
export type ScheduleTrainingResponse = Message<"soccerbuddy.team.v1.ScheduleTrainingResponse"> & {
};

/**
 * Describes the message soccerbuddy.team.v1.ScheduleTrainingResponse.
 * Use `create(ScheduleTrainingResponseSchema)` to create a new message.
 */
export const ScheduleTrainingResponseSchema: GenMessage<ScheduleTrainingResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 19);

/**
 * @generated from message soccerbuddy.team.v1.GetMyTeamHomeRequest
 */
export type GetMyTeamHomeRequest = Message<"soccerbuddy.team.v1.GetMyTeamHomeRequest"> & {
  /**
   * @generated from field: string team_id = 1;
   */
  teamId: string;
};

/**
 * Describes the message soccerbuddy.team.v1.GetMyTeamHomeRequest.
 * Use `create(GetMyTeamHomeRequestSchema)` to create a new message.
 */
export const GetMyTeamHomeRequestSchema: GenMessage<GetMyTeamHomeRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 20);

/**
 * @generated from message soccerbuddy.team.v1.GetMyTeamHomeResponse
 */
export type GetMyTeamHomeResponse = Message<"soccerbuddy.team.v1.GetMyTeamHomeResponse"> & {
  /**
   * @generated from field: string team_id = 1;
   */
  teamId: string;

  /**
   * @generated from field: string team_name = 2;
   */
  teamName: string;

  /**
   * @generated from field: repeated soccerbuddy.team.v1.GetMyTeamHomeResponse.Training trainings = 3;
   */
  trainings: GetMyTeamHomeResponse_Training[];
};

/**
 * Describes the message soccerbuddy.team.v1.GetMyTeamHomeResponse.
 * Use `create(GetMyTeamHomeResponseSchema)` to create a new message.
 */
export const GetMyTeamHomeResponseSchema: GenMessage<GetMyTeamHomeResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 21);

/**
 * @generated from message soccerbuddy.team.v1.GetMyTeamHomeResponse.Training
 */
export type GetMyTeamHomeResponse_Training = Message<"soccerbuddy.team.v1.GetMyTeamHomeResponse.Training"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: google.type.DateTime scheduled_at = 2;
   */
  scheduledAt?: DateTime;

  /**
   * @generated from field: google.type.DateTime ends_at = 3;
   */
  endsAt?: DateTime;

  /**
   * @generated from field: optional string location = 4;
   */
  location?: string;

  /**
   * The type of field used: e.g. hard floor, lawn, etc.
   *
   * @generated from field: optional string field_type = 5;
   */
  fieldType?: string;

  /**
   * @generated from field: optional soccerbuddy.team.v1.GatheringPoint gathering_point = 6;
   */
  gatheringPoint?: GatheringPoint;

  /**
   * @generated from field: optional soccerbuddy.team.v1.AcknowledgementSettings acknowledgment_settings = 7;
   */
  acknowledgmentSettings?: AcknowledgementSettings;

  /**
   * @generated from field: optional string description = 8;
   */
  description?: string;

  /**
   * @generated from field: optional soccerbuddy.team.v1.RatingSettings rating_settings = 9;
   */
  ratingSettings?: RatingSettings;
};

/**
 * Describes the message soccerbuddy.team.v1.GetMyTeamHomeResponse.Training.
 * Use `create(GetMyTeamHomeResponse_TrainingSchema)` to create a new message.
 */
export const GetMyTeamHomeResponse_TrainingSchema: GenMessage<GetMyTeamHomeResponse_Training> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 21, 0);

/**
 * @generated from message soccerbuddy.team.v1.NominatePersonsForTrainingRequest
 */
export type NominatePersonsForTrainingRequest = Message<"soccerbuddy.team.v1.NominatePersonsForTrainingRequest"> & {
  /**
   * @generated from field: string training_id = 1;
   */
  trainingId: string;

  /**
   * @generated from field: repeated string player_ids = 2;
   */
  playerIds: string[];

  /**
   * @generated from field: repeated string staff_ids = 3;
   */
  staffIds: string[];
};

/**
 * Describes the message soccerbuddy.team.v1.NominatePersonsForTrainingRequest.
 * Use `create(NominatePersonsForTrainingRequestSchema)` to create a new message.
 */
export const NominatePersonsForTrainingRequestSchema: GenMessage<NominatePersonsForTrainingRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 22);

/**
 * @generated from message soccerbuddy.team.v1.NominatePersonsForTrainingResponse
 */
export type NominatePersonsForTrainingResponse = Message<"soccerbuddy.team.v1.NominatePersonsForTrainingResponse"> & {
};

/**
 * Describes the message soccerbuddy.team.v1.NominatePersonsForTrainingResponse.
 * Use `create(NominatePersonsForTrainingResponseSchema)` to create a new message.
 */
export const NominatePersonsForTrainingResponseSchema: GenMessage<NominatePersonsForTrainingResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_team_v1_team_service, 23);

/**
 * @generated from service soccerbuddy.team.v1.TeamService
 */
export const TeamService: GenService<{
  /**
   * @generated from rpc soccerbuddy.team.v1.TeamService.CreateTeam
   */
  createTeam: {
    methodKind: "unary";
    input: typeof CreateTeamRequestSchema;
    output: typeof CreateTeamResponseSchema;
  },
  /**
   * @generated from rpc soccerbuddy.team.v1.TeamService.ListTeams
   */
  listTeams: {
    methodKind: "unary";
    input: typeof ListTeamsRequestSchema;
    output: typeof ListTeamsResponseSchema;
  },
  /**
   * @generated from rpc soccerbuddy.team.v1.TeamService.DeleteTeam
   */
  deleteTeam: {
    methodKind: "unary";
    input: typeof DeleteTeamRequestSchema;
    output: typeof DeleteTeamResponseSchema;
  },
  /**
   * @generated from rpc soccerbuddy.team.v1.TeamService.GetTeamOverview
   */
  getTeamOverview: {
    methodKind: "unary";
    input: typeof GetTeamOverviewRequestSchema;
    output: typeof GetTeamOverviewResponseSchema;
  },
  /**
   * @generated from rpc soccerbuddy.team.v1.TeamService.SearchPersonsNotInTeam
   */
  searchPersonsNotInTeam: {
    methodKind: "unary";
    input: typeof SearchPersonsNotInTeamRequestSchema;
    output: typeof SearchPersonsNotInTeamResponseSchema;
  },
  /**
   * @generated from rpc soccerbuddy.team.v1.TeamService.AddPersonToTeam
   */
  addPersonToTeam: {
    methodKind: "unary";
    input: typeof AddPersonToTeamRequestSchema;
    output: typeof AddPersonToTeamResponseSchema;
  },
  /**
   * @generated from rpc soccerbuddy.team.v1.TeamService.ListTeamMembers
   */
  listTeamMembers: {
    methodKind: "unary";
    input: typeof ListTeamMembersRequestSchema;
    output: typeof ListTeamMembersResponseSchema;
  },
  /**
   * @generated from rpc soccerbuddy.team.v1.TeamService.ScheduleTraining
   */
  scheduleTraining: {
    methodKind: "unary";
    input: typeof ScheduleTrainingRequestSchema;
    output: typeof ScheduleTrainingResponseSchema;
  },
  /**
   * @generated from rpc soccerbuddy.team.v1.TeamService.GetMyTeamHome
   */
  getMyTeamHome: {
    methodKind: "unary";
    input: typeof GetMyTeamHomeRequestSchema;
    output: typeof GetMyTeamHomeResponseSchema;
  },
  /**
   * @generated from rpc soccerbuddy.team.v1.TeamService.NominatePersonsForTraining
   */
  nominatePersonsForTraining: {
    methodKind: "unary";
    input: typeof NominatePersonsForTrainingRequestSchema;
    output: typeof NominatePersonsForTrainingResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_soccerbuddy_team_v1_team_service, 0);

