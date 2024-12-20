// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file soccerbuddy/club/v1/club_service.proto (package soccerbuddy.club.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file soccerbuddy/club/v1/club_service.proto.
 */
export const file_soccerbuddy_club_v1_club_service: GenFile = /*@__PURE__*/
  fileDesc("CiZzb2NjZXJidWRkeS9jbHViL3YxL2NsdWJfc2VydmljZS5wcm90bxITc29jY2VyYnVkZHkuY2x1Yi52MSIhChFDcmVhdGVDbHViUmVxdWVzdBIMCgRuYW1lGAEgASgJIjwKEkNyZWF0ZUNsdWJSZXNwb25zZRIKCgJpZBgBIAEoCRIMCgRuYW1lGAIgASgJEgwKBHNsdWcYAyABKAkiJAoUR2V0Q2x1YkJ5U2x1Z1JlcXVlc3QSDAoEc2x1ZxgBIAEoCSKfAQoVR2V0Q2x1YkJ5U2x1Z1Jlc3BvbnNlEgoKAmlkGAEgASgJEgwKBG5hbWUYAiABKAkSDAoEc2x1ZxgDIAEoCRIuCgpjcmVhdGVkX2F0GAQgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcBIuCgp1cGRhdGVkX2F0GAUgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcDLYAQoLQ2x1YlNlcnZpY2USXwoKQ3JlYXRlQ2x1YhImLnNvY2NlcmJ1ZGR5LmNsdWIudjEuQ3JlYXRlQ2x1YlJlcXVlc3QaJy5zb2NjZXJidWRkeS5jbHViLnYxLkNyZWF0ZUNsdWJSZXNwb25zZSIAEmgKDUdldENsdWJCeVNsdWcSKS5zb2NjZXJidWRkeS5jbHViLnYxLkdldENsdWJCeVNsdWdSZXF1ZXN0Giouc29jY2VyYnVkZHkuY2x1Yi52MS5HZXRDbHViQnlTbHVnUmVzcG9uc2UiAELaAQoXY29tLnNvY2NlcmJ1ZGR5LmNsdWIudjFCEENsdWJTZXJ2aWNlUHJvdG9QAVo/Z2l0aHViLmNvbS9yc21pZHQvc29jY2VyYnVkZHkvZ2VuL2dvL3NvY2NlcmJ1ZGR5L2NsdWIvdjE7Y2x1YnYxogIDU0NYqgITU29jY2VyYnVkZHkuQ2x1Yi5WMcoCE1NvY2NlcmJ1ZGR5XENsdWJcVjHiAh9Tb2NjZXJidWRkeVxDbHViXFYxXEdQQk1ldGFkYXRh6gIVU29jY2VyYnVkZHk6OkNsdWI6OlYxYgZwcm90bzM", [file_google_protobuf_timestamp]);

/**
 * @generated from message soccerbuddy.club.v1.CreateClubRequest
 */
export type CreateClubRequest = Message<"soccerbuddy.club.v1.CreateClubRequest"> & {
  /**
   * @generated from field: string name = 1;
   */
  name: string;
};

/**
 * Describes the message soccerbuddy.club.v1.CreateClubRequest.
 * Use `create(CreateClubRequestSchema)` to create a new message.
 */
export const CreateClubRequestSchema: GenMessage<CreateClubRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_club_v1_club_service, 0);

/**
 * @generated from message soccerbuddy.club.v1.CreateClubResponse
 */
export type CreateClubResponse = Message<"soccerbuddy.club.v1.CreateClubResponse"> & {
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
};

/**
 * Describes the message soccerbuddy.club.v1.CreateClubResponse.
 * Use `create(CreateClubResponseSchema)` to create a new message.
 */
export const CreateClubResponseSchema: GenMessage<CreateClubResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_club_v1_club_service, 1);

/**
 * @generated from message soccerbuddy.club.v1.GetClubBySlugRequest
 */
export type GetClubBySlugRequest = Message<"soccerbuddy.club.v1.GetClubBySlugRequest"> & {
  /**
   * @generated from field: string slug = 1;
   */
  slug: string;
};

/**
 * Describes the message soccerbuddy.club.v1.GetClubBySlugRequest.
 * Use `create(GetClubBySlugRequestSchema)` to create a new message.
 */
export const GetClubBySlugRequestSchema: GenMessage<GetClubBySlugRequest> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_club_v1_club_service, 2);

/**
 * @generated from message soccerbuddy.club.v1.GetClubBySlugResponse
 */
export type GetClubBySlugResponse = Message<"soccerbuddy.club.v1.GetClubBySlugResponse"> & {
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
 * Describes the message soccerbuddy.club.v1.GetClubBySlugResponse.
 * Use `create(GetClubBySlugResponseSchema)` to create a new message.
 */
export const GetClubBySlugResponseSchema: GenMessage<GetClubBySlugResponse> = /*@__PURE__*/
  messageDesc(file_soccerbuddy_club_v1_club_service, 3);

/**
 * @generated from service soccerbuddy.club.v1.ClubService
 */
export const ClubService: GenService<{
  /**
   * @generated from rpc soccerbuddy.club.v1.ClubService.CreateClub
   */
  createClub: {
    methodKind: "unary";
    input: typeof CreateClubRequestSchema;
    output: typeof CreateClubResponseSchema;
  },
  /**
   * @generated from rpc soccerbuddy.club.v1.ClubService.GetClubBySlug
   */
  getClubBySlug: {
    methodKind: "unary";
    input: typeof GetClubBySlugRequestSchema;
    output: typeof GetClubBySlugResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_soccerbuddy_club_v1_club_service, 0);

