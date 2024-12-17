// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
import { fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv1";
import { file_google_protobuf_duration } from "@bufbuild/protobuf/wkt";
/**
 * Describes the file google/rpc/error_details.proto.
 */
export const file_google_rpc_error_details = 
/*@__PURE__*/
fileDesc("Ch5nb29nbGUvcnBjL2Vycm9yX2RldGFpbHMucHJvdG8SCmdvb2dsZS5ycGMikwEKCUVycm9ySW5mbxIOCgZyZWFzb24YASABKAkSDgoGZG9tYWluGAIgASgJEjUKCG1ldGFkYXRhGAMgAygLMiMuZ29vZ2xlLnJwYy5FcnJvckluZm8uTWV0YWRhdGFFbnRyeRovCg1NZXRhZGF0YUVudHJ5EgsKA2tleRgBIAEoCRINCgV2YWx1ZRgCIAEoCToCOAEiOwoJUmV0cnlJbmZvEi4KC3JldHJ5X2RlbGF5GAEgASgLMhkuZ29vZ2xlLnByb3RvYnVmLkR1cmF0aW9uIjIKCURlYnVnSW5mbxIVCg1zdGFja19lbnRyaWVzGAEgAygJEg4KBmRldGFpbBgCIAEoCSJ5CgxRdW90YUZhaWx1cmUSNgoKdmlvbGF0aW9ucxgBIAMoCzIiLmdvb2dsZS5ycGMuUXVvdGFGYWlsdXJlLlZpb2xhdGlvbhoxCglWaW9sYXRpb24SDwoHc3ViamVjdBgBIAEoCRITCgtkZXNjcmlwdGlvbhgCIAEoCSKVAQoTUHJlY29uZGl0aW9uRmFpbHVyZRI9Cgp2aW9sYXRpb25zGAEgAygLMikuZ29vZ2xlLnJwYy5QcmVjb25kaXRpb25GYWlsdXJlLlZpb2xhdGlvbho/CglWaW9sYXRpb24SDAoEdHlwZRgBIAEoCRIPCgdzdWJqZWN0GAIgASgJEhMKC2Rlc2NyaXB0aW9uGAMgASgJIoMBCgpCYWRSZXF1ZXN0Ej8KEGZpZWxkX3Zpb2xhdGlvbnMYASADKAsyJS5nb29nbGUucnBjLkJhZFJlcXVlc3QuRmllbGRWaW9sYXRpb24aNAoORmllbGRWaW9sYXRpb24SDQoFZmllbGQYASABKAkSEwoLZGVzY3JpcHRpb24YAiABKAkiNwoLUmVxdWVzdEluZm8SEgoKcmVxdWVzdF9pZBgBIAEoCRIUCgxzZXJ2aW5nX2RhdGEYAiABKAkiYAoMUmVzb3VyY2VJbmZvEhUKDXJlc291cmNlX3R5cGUYASABKAkSFQoNcmVzb3VyY2VfbmFtZRgCIAEoCRINCgVvd25lchgDIAEoCRITCgtkZXNjcmlwdGlvbhgEIAEoCSJWCgRIZWxwEiQKBWxpbmtzGAEgAygLMhUuZ29vZ2xlLnJwYy5IZWxwLkxpbmsaKAoETGluaxITCgtkZXNjcmlwdGlvbhgBIAEoCRILCgN1cmwYAiABKAkiMwoQTG9jYWxpemVkTWVzc2FnZRIOCgZsb2NhbGUYASABKAkSDwoHbWVzc2FnZRgCIAEoCUKdAQoOY29tLmdvb2dsZS5ycGNCEUVycm9yRGV0YWlsc1Byb3RvUAFaL2dpdGh1Yi5jb20vcnNtaWR0L3NvY2NlcmJ1ZGR5L2dlbi9nby9nb29nbGUvcnBjogIDR1JYqgIKR29vZ2xlLlJwY8oCCkdvb2dsZVxScGPiAhZHb29nbGVcUnBjXEdQQk1ldGFkYXRh6gILR29vZ2xlOjpScGNiBnByb3RvMw", [file_google_protobuf_duration]);
/**
 * Describes the message google.rpc.ErrorInfo.
 * Use `create(ErrorInfoSchema)` to create a new message.
 */
export const ErrorInfoSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 0);
/**
 * Describes the message google.rpc.RetryInfo.
 * Use `create(RetryInfoSchema)` to create a new message.
 */
export const RetryInfoSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 1);
/**
 * Describes the message google.rpc.DebugInfo.
 * Use `create(DebugInfoSchema)` to create a new message.
 */
export const DebugInfoSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 2);
/**
 * Describes the message google.rpc.QuotaFailure.
 * Use `create(QuotaFailureSchema)` to create a new message.
 */
export const QuotaFailureSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 3);
/**
 * Describes the message google.rpc.QuotaFailure.Violation.
 * Use `create(QuotaFailure_ViolationSchema)` to create a new message.
 */
export const QuotaFailure_ViolationSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 3, 0);
/**
 * Describes the message google.rpc.PreconditionFailure.
 * Use `create(PreconditionFailureSchema)` to create a new message.
 */
export const PreconditionFailureSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 4);
/**
 * Describes the message google.rpc.PreconditionFailure.Violation.
 * Use `create(PreconditionFailure_ViolationSchema)` to create a new message.
 */
export const PreconditionFailure_ViolationSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 4, 0);
/**
 * Describes the message google.rpc.BadRequest.
 * Use `create(BadRequestSchema)` to create a new message.
 */
export const BadRequestSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 5);
/**
 * Describes the message google.rpc.BadRequest.FieldViolation.
 * Use `create(BadRequest_FieldViolationSchema)` to create a new message.
 */
export const BadRequest_FieldViolationSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 5, 0);
/**
 * Describes the message google.rpc.RequestInfo.
 * Use `create(RequestInfoSchema)` to create a new message.
 */
export const RequestInfoSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 6);
/**
 * Describes the message google.rpc.ResourceInfo.
 * Use `create(ResourceInfoSchema)` to create a new message.
 */
export const ResourceInfoSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 7);
/**
 * Describes the message google.rpc.Help.
 * Use `create(HelpSchema)` to create a new message.
 */
export const HelpSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 8);
/**
 * Describes the message google.rpc.Help.Link.
 * Use `create(Help_LinkSchema)` to create a new message.
 */
export const Help_LinkSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 8, 0);
/**
 * Describes the message google.rpc.LocalizedMessage.
 * Use `create(LocalizedMessageSchema)` to create a new message.
 */
export const LocalizedMessageSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_error_details, 9);
