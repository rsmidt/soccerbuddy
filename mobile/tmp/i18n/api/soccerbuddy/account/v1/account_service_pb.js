// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file soccerbuddy/account/v1/account_service.proto (package soccerbuddy.account.v1, syntax proto3)
/* eslint-disable */
import { fileDesc, messageDesc, serviceDesc, } from "@bufbuild/protobuf/codegenv1";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import { file_soccerbuddy_shared } from "../../shared_pb";
/**
 * Describes the file soccerbuddy/account/v1/account_service.proto.
 */
export const file_soccerbuddy_account_v1_account_service = 
/*@__PURE__*/
fileDesc("Cixzb2NjZXJidWRkeS9hY2NvdW50L3YxL2FjY291bnRfc2VydmljZS5wcm90bxIWc29jY2VyYnVkZHkuYWNjb3VudC52MSIOCgxHZXRNZVJlcXVlc3QisQUKDUdldE1lUmVzcG9uc2USCgoCaWQYASABKAkSDQoFZW1haWwYAiABKAkSEgoKZmlyc3RfbmFtZRgDIAEoCRIRCglsYXN0X25hbWUYBCABKAkSSgoObGlua2VkX3BlcnNvbnMYBSADKAsyMi5zb2NjZXJidWRkeS5hY2NvdW50LnYxLkdldE1lUmVzcG9uc2UuTGlua2VkUGVyc29uGiwKCE9wZXJhdG9yEhEKCWZ1bGxfbmFtZRgBIAEoCRINCgVpc19tZRgCIAEoCBriAgoMTGlua2VkUGVyc29uEgoKAmlkGAEgASgJEjIKCWxpbmtlZF9hcxgCIAEoDjIfLnNvY2NlcmJ1ZGR5LnNoYXJlZC5BY2NvdW50TGluaxISCgpmaXJzdF9uYW1lGAMgASgJEhEKCWxhc3RfbmFtZRgEIAEoCRItCglsaW5rZWRfYXQYBSABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wEkYKCWxpbmtlZF9ieRgGIAEoCzIuLnNvY2NlcmJ1ZGR5LmFjY291bnQudjEuR2V0TWVSZXNwb25zZS5PcGVyYXRvckgAiAEBEk4KEHRlYW1fbWVtYmVyc2hpcHMYByADKAsyNC5zb2NjZXJidWRkeS5hY2NvdW50LnYxLkdldE1lUmVzcG9uc2UuVGVhbU1lbWJlcnNoaXASFgoOb3duaW5nX2NsdWJfaWQYCCABKAlCDAoKX2xpbmtlZF9ieRp/Cg5UZWFtTWVtYmVyc2hpcBIKCgJpZBgBIAEoCRIMCgRuYW1lGAIgASgJEgwKBHJvbGUYAyABKAkSLQoJam9pbmVkX2F0GAQgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcBIWCg5vd25pbmdfY2x1Yl9pZBgFIAEoCSJeChRDcmVhdGVBY2NvdW50UmVxdWVzdBINCgVlbWFpbBgBIAEoCRIQCghwYXNzd29yZBgCIAEoCRISCgpmaXJzdF9uYW1lGAMgASgJEhEKCWxhc3RfbmFtZRgEIAEoCSJZChVDcmVhdGVBY2NvdW50UmVzcG9uc2USCgoCaWQYASABKAkSDQoFZW1haWwYAiABKAkSEgoKZmlyc3RfbmFtZRgDIAEoCRIRCglsYXN0X25hbWUYBCABKAkiQwoMTG9naW5SZXF1ZXN0Eg0KBWVtYWlsGAEgASgJEhAKCHBhc3N3b3JkGAIgASgJEhIKCnVzZXJfYWdlbnQYAyABKAkiIwoNTG9naW5SZXNwb25zZRISCgpzZXNzaW9uX2lkGAEgASgJInQKFlJlZ2lzdGVyQWNjb3VudFJlcXVlc3QSEgoKZmlyc3RfbmFtZRgBIAEoCRIRCglsYXN0X25hbWUYAiABKAkSDQoFZW1haWwYAyABKAkSEAoIcGFzc3dvcmQYBCABKAkSEgoKbGlua190b2tlbhgFIAEoCSIlChdSZWdpc3RlckFjY291bnRSZXNwb25zZRIKCgJpZBgBIAEoCTKmAwoOQWNjb3VudFNlcnZpY2USVgoFR2V0TWUSJC5zb2NjZXJidWRkeS5hY2NvdW50LnYxLkdldE1lUmVxdWVzdBolLnNvY2NlcmJ1ZGR5LmFjY291bnQudjEuR2V0TWVSZXNwb25zZSIAEm4KDUNyZWF0ZUFjY291bnQSLC5zb2NjZXJidWRkeS5hY2NvdW50LnYxLkNyZWF0ZUFjY291bnRSZXF1ZXN0Gi0uc29jY2VyYnVkZHkuYWNjb3VudC52MS5DcmVhdGVBY2NvdW50UmVzcG9uc2UiABJWCgVMb2dpbhIkLnNvY2NlcmJ1ZGR5LmFjY291bnQudjEuTG9naW5SZXF1ZXN0GiUuc29jY2VyYnVkZHkuYWNjb3VudC52MS5Mb2dpblJlc3BvbnNlIgASdAoPUmVnaXN0ZXJBY2NvdW50Ei4uc29jY2VyYnVkZHkuYWNjb3VudC52MS5SZWdpc3RlckFjY291bnRSZXF1ZXN0Gi8uc29jY2VyYnVkZHkuYWNjb3VudC52MS5SZWdpc3RlckFjY291bnRSZXNwb25zZSIAQvIBChpjb20uc29jY2VyYnVkZHkuYWNjb3VudC52MUITQWNjb3VudFNlcnZpY2VQcm90b1ABWkVnaXRodWIuY29tL3JzbWlkdC9zb2NjZXJidWRkeS9nZW4vZ28vc29jY2VyYnVkZHkvYWNjb3VudC92MTthY2NvdW50djGiAgNTQViqAhZTb2NjZXJidWRkeS5BY2NvdW50LlYxygIWU29jY2VyYnVkZHlcQWNjb3VudFxWMeICIlNvY2NlcmJ1ZGR5XEFjY291bnRcVjFcR1BCTWV0YWRhdGHqAhhTb2NjZXJidWRkeTo6QWNjb3VudDo6VjFiBnByb3RvMw", [file_google_protobuf_timestamp, file_soccerbuddy_shared]);
/**
 * Describes the message soccerbuddy.account.v1.GetMeRequest.
 * Use `create(GetMeRequestSchema)` to create a new message.
 */
export const GetMeRequestSchema = 
/*@__PURE__*/
messageDesc(file_soccerbuddy_account_v1_account_service, 0);
/**
 * Describes the message soccerbuddy.account.v1.GetMeResponse.
 * Use `create(GetMeResponseSchema)` to create a new message.
 */
export const GetMeResponseSchema = 
/*@__PURE__*/
messageDesc(file_soccerbuddy_account_v1_account_service, 1);
/**
 * Describes the message soccerbuddy.account.v1.GetMeResponse.Operator.
 * Use `create(GetMeResponse_OperatorSchema)` to create a new message.
 */
export const GetMeResponse_OperatorSchema = 
/*@__PURE__*/
messageDesc(file_soccerbuddy_account_v1_account_service, 1, 0);
/**
 * Describes the message soccerbuddy.account.v1.GetMeResponse.LinkedPerson.
 * Use `create(GetMeResponse_LinkedPersonSchema)` to create a new message.
 */
export const GetMeResponse_LinkedPersonSchema = 
/*@__PURE__*/
messageDesc(file_soccerbuddy_account_v1_account_service, 1, 1);
/**
 * Describes the message soccerbuddy.account.v1.GetMeResponse.TeamMembership.
 * Use `create(GetMeResponse_TeamMembershipSchema)` to create a new message.
 */
export const GetMeResponse_TeamMembershipSchema = 
/*@__PURE__*/
messageDesc(file_soccerbuddy_account_v1_account_service, 1, 2);
/**
 * Describes the message soccerbuddy.account.v1.CreateAccountRequest.
 * Use `create(CreateAccountRequestSchema)` to create a new message.
 */
export const CreateAccountRequestSchema = 
/*@__PURE__*/
messageDesc(file_soccerbuddy_account_v1_account_service, 2);
/**
 * Describes the message soccerbuddy.account.v1.CreateAccountResponse.
 * Use `create(CreateAccountResponseSchema)` to create a new message.
 */
export const CreateAccountResponseSchema = 
/*@__PURE__*/
messageDesc(file_soccerbuddy_account_v1_account_service, 3);
/**
 * Describes the message soccerbuddy.account.v1.LoginRequest.
 * Use `create(LoginRequestSchema)` to create a new message.
 */
export const LoginRequestSchema = 
/*@__PURE__*/
messageDesc(file_soccerbuddy_account_v1_account_service, 4);
/**
 * Describes the message soccerbuddy.account.v1.LoginResponse.
 * Use `create(LoginResponseSchema)` to create a new message.
 */
export const LoginResponseSchema = 
/*@__PURE__*/
messageDesc(file_soccerbuddy_account_v1_account_service, 5);
/**
 * Describes the message soccerbuddy.account.v1.RegisterAccountRequest.
 * Use `create(RegisterAccountRequestSchema)` to create a new message.
 */
export const RegisterAccountRequestSchema = 
/*@__PURE__*/
messageDesc(file_soccerbuddy_account_v1_account_service, 6);
/**
 * Describes the message soccerbuddy.account.v1.RegisterAccountResponse.
 * Use `create(RegisterAccountResponseSchema)` to create a new message.
 */
export const RegisterAccountResponseSchema = 
/*@__PURE__*/
messageDesc(file_soccerbuddy_account_v1_account_service, 7);
/**
 * @generated from service soccerbuddy.account.v1.AccountService
 */
export const AccountService = /*@__PURE__*/ serviceDesc(file_soccerbuddy_account_v1_account_service, 0);