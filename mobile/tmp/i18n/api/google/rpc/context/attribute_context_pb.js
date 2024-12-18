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
import { file_google_protobuf_any, file_google_protobuf_duration, file_google_protobuf_struct, file_google_protobuf_timestamp, } from "@bufbuild/protobuf/wkt";
/**
 * Describes the file google/rpc/context/attribute_context.proto.
 */
export const file_google_rpc_context_attribute_context = 
/*@__PURE__*/
fileDesc("Cipnb29nbGUvcnBjL2NvbnRleHQvYXR0cmlidXRlX2NvbnRleHQucHJvdG8SEmdvb2dsZS5ycGMuY29udGV4dCKDEAoQQXR0cmlidXRlQ29udGV4dBI5CgZvcmlnaW4YByABKAsyKS5nb29nbGUucnBjLmNvbnRleHQuQXR0cmlidXRlQ29udGV4dC5QZWVyEjkKBnNvdXJjZRgBIAEoCzIpLmdvb2dsZS5ycGMuY29udGV4dC5BdHRyaWJ1dGVDb250ZXh0LlBlZXISPgoLZGVzdGluYXRpb24YAiABKAsyKS5nb29nbGUucnBjLmNvbnRleHQuQXR0cmlidXRlQ29udGV4dC5QZWVyEj0KB3JlcXVlc3QYAyABKAsyLC5nb29nbGUucnBjLmNvbnRleHQuQXR0cmlidXRlQ29udGV4dC5SZXF1ZXN0Ej8KCHJlc3BvbnNlGAQgASgLMi0uZ29vZ2xlLnJwYy5jb250ZXh0LkF0dHJpYnV0ZUNvbnRleHQuUmVzcG9uc2USPwoIcmVzb3VyY2UYBSABKAsyLS5nb29nbGUucnBjLmNvbnRleHQuQXR0cmlidXRlQ29udGV4dC5SZXNvdXJjZRI1CgNhcGkYBiABKAsyKC5nb29nbGUucnBjLmNvbnRleHQuQXR0cmlidXRlQ29udGV4dC5BcGkSKAoKZXh0ZW5zaW9ucxgIIAMoCzIULmdvb2dsZS5wcm90b2J1Zi5BbnkavgEKBFBlZXISCgoCaXAYASABKAkSDAoEcG9ydBgCIAEoAxJFCgZsYWJlbHMYBiADKAsyNS5nb29nbGUucnBjLmNvbnRleHQuQXR0cmlidXRlQ29udGV4dC5QZWVyLkxhYmVsc0VudHJ5EhEKCXByaW5jaXBhbBgHIAEoCRITCgtyZWdpb25fY29kZRgIIAEoCRotCgtMYWJlbHNFbnRyeRILCgNrZXkYASABKAkSDQoFdmFsdWUYAiABKAk6AjgBGkwKA0FwaRIPCgdzZXJ2aWNlGAEgASgJEhEKCW9wZXJhdGlvbhgCIAEoCRIQCghwcm90b2NvbBgDIAEoCRIPCgd2ZXJzaW9uGAQgASgJGn8KBEF1dGgSEQoJcHJpbmNpcGFsGAEgASgJEhEKCWF1ZGllbmNlcxgCIAMoCRIRCglwcmVzZW50ZXIYAyABKAkSJwoGY2xhaW1zGAQgASgLMhcuZ29vZ2xlLnByb3RvYnVmLlN0cnVjdBIVCg1hY2Nlc3NfbGV2ZWxzGAUgAygJGu8CCgdSZXF1ZXN0EgoKAmlkGAEgASgJEg4KBm1ldGhvZBgCIAEoCRJKCgdoZWFkZXJzGAMgAygLMjkuZ29vZ2xlLnJwYy5jb250ZXh0LkF0dHJpYnV0ZUNvbnRleHQuUmVxdWVzdC5IZWFkZXJzRW50cnkSDAoEcGF0aBgEIAEoCRIMCgRob3N0GAUgASgJEg4KBnNjaGVtZRgGIAEoCRINCgVxdWVyeRgHIAEoCRIoCgR0aW1lGAkgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcBIMCgRzaXplGAogASgDEhAKCHByb3RvY29sGAsgASgJEg4KBnJlYXNvbhgMIAEoCRI3CgRhdXRoGA0gASgLMikuZ29vZ2xlLnJwYy5jb250ZXh0LkF0dHJpYnV0ZUNvbnRleHQuQXV0aBouCgxIZWFkZXJzRW50cnkSCwoDa2V5GAEgASgJEg0KBXZhbHVlGAIgASgJOgI4ARqBAgoIUmVzcG9uc2USDAoEY29kZRgBIAEoAxIMCgRzaXplGAIgASgDEksKB2hlYWRlcnMYAyADKAsyOi5nb29nbGUucnBjLmNvbnRleHQuQXR0cmlidXRlQ29udGV4dC5SZXNwb25zZS5IZWFkZXJzRW50cnkSKAoEdGltZRgEIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXASMgoPYmFja2VuZF9sYXRlbmN5GAUgASgLMhkuZ29vZ2xlLnByb3RvYnVmLkR1cmF0aW9uGi4KDEhlYWRlcnNFbnRyeRILCgNrZXkYASABKAkSDQoFdmFsdWUYAiABKAk6AjgBGpAECghSZXNvdXJjZRIPCgdzZXJ2aWNlGAEgASgJEgwKBG5hbWUYAiABKAkSDAoEdHlwZRgDIAEoCRJJCgZsYWJlbHMYBCADKAsyOS5nb29nbGUucnBjLmNvbnRleHQuQXR0cmlidXRlQ29udGV4dC5SZXNvdXJjZS5MYWJlbHNFbnRyeRILCgN1aWQYBSABKAkSUwoLYW5ub3RhdGlvbnMYBiADKAsyPi5nb29nbGUucnBjLmNvbnRleHQuQXR0cmlidXRlQ29udGV4dC5SZXNvdXJjZS5Bbm5vdGF0aW9uc0VudHJ5EhQKDGRpc3BsYXlfbmFtZRgHIAEoCRIvCgtjcmVhdGVfdGltZRgIIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXASLwoLdXBkYXRlX3RpbWUYCSABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wEi8KC2RlbGV0ZV90aW1lGAogASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcBIMCgRldGFnGAsgASgJEhAKCGxvY2F0aW9uGAwgASgJGi0KC0xhYmVsc0VudHJ5EgsKA2tleRgBIAEoCRINCgV2YWx1ZRgCIAEoCToCOAEaMgoQQW5ub3RhdGlvbnNFbnRyeRILCgNrZXkYASABKAkSDQoFdmFsdWUYAiABKAk6AjgBQtUBChZjb20uZ29vZ2xlLnJwYy5jb250ZXh0QhVBdHRyaWJ1dGVDb250ZXh0UHJvdG9QAVo3Z2l0aHViLmNvbS9yc21pZHQvc29jY2VyYnVkZHkvZ2VuL2dvL2dvb2dsZS9ycGMvY29udGV4dPgBAaICA0dSQ6oCEkdvb2dsZS5ScGMuQ29udGV4dMoCEkdvb2dsZVxScGNcQ29udGV4dOICHkdvb2dsZVxScGNcQ29udGV4dFxHUEJNZXRhZGF0YeoCFEdvb2dsZTo6UnBjOjpDb250ZXh0YgZwcm90bzM", [
    file_google_protobuf_any,
    file_google_protobuf_duration,
    file_google_protobuf_struct,
    file_google_protobuf_timestamp,
]);
/**
 * Describes the message google.rpc.context.AttributeContext.
 * Use `create(AttributeContextSchema)` to create a new message.
 */
export const AttributeContextSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_context_attribute_context, 0);
/**
 * Describes the message google.rpc.context.AttributeContext.Peer.
 * Use `create(AttributeContext_PeerSchema)` to create a new message.
 */
export const AttributeContext_PeerSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_context_attribute_context, 0, 0);
/**
 * Describes the message google.rpc.context.AttributeContext.Api.
 * Use `create(AttributeContext_ApiSchema)` to create a new message.
 */
export const AttributeContext_ApiSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_context_attribute_context, 0, 1);
/**
 * Describes the message google.rpc.context.AttributeContext.Auth.
 * Use `create(AttributeContext_AuthSchema)` to create a new message.
 */
export const AttributeContext_AuthSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_context_attribute_context, 0, 2);
/**
 * Describes the message google.rpc.context.AttributeContext.Request.
 * Use `create(AttributeContext_RequestSchema)` to create a new message.
 */
export const AttributeContext_RequestSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_context_attribute_context, 0, 3);
/**
 * Describes the message google.rpc.context.AttributeContext.Response.
 * Use `create(AttributeContext_ResponseSchema)` to create a new message.
 */
export const AttributeContext_ResponseSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_context_attribute_context, 0, 4);
/**
 * Describes the message google.rpc.context.AttributeContext.Resource.
 * Use `create(AttributeContext_ResourceSchema)` to create a new message.
 */
export const AttributeContext_ResourceSchema = 
/*@__PURE__*/
messageDesc(file_google_rpc_context_attribute_context, 0, 5);