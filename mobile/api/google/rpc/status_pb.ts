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

// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file google/rpc/status.proto (package google.rpc, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv1";
import type { Any } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_any } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file google/rpc/status.proto.
 */
export const file_google_rpc_status: GenFile =
  /*@__PURE__*/
  fileDesc(
    "Chdnb29nbGUvcnBjL3N0YXR1cy5wcm90bxIKZ29vZ2xlLnJwYyJOCgZTdGF0dXMSDAoEY29kZRgBIAEoBRIPCgdtZXNzYWdlGAIgASgJEiUKB2RldGFpbHMYAyADKAsyFC5nb29nbGUucHJvdG9idWYuQW55QpoBCg5jb20uZ29vZ2xlLnJwY0ILU3RhdHVzUHJvdG9QAVovZ2l0aHViLmNvbS9yc21pZHQvc29jY2VyYnVkZHkvZ2VuL2dvL2dvb2dsZS9ycGP4AQGiAgNHUliqAgpHb29nbGUuUnBjygIKR29vZ2xlXFJwY+ICFkdvb2dsZVxScGNcR1BCTWV0YWRhdGHqAgtHb29nbGU6OlJwY2IGcHJvdG8z",
    [file_google_protobuf_any],
  );

/**
 * The `Status` type defines a logical error model that is suitable for
 * different programming environments, including REST APIs and RPC APIs. It is
 * used by [gRPC](https://github.com/grpc). Each `Status` message contains
 * three pieces of data: error code, error message, and error details.
 *
 * You can find out more about this error model and how to work with it in the
 * [API Design Guide](https://cloud.google.com/apis/design/errors).
 *
 * @generated from message google.rpc.Status
 */
export type Status = Message<"google.rpc.Status"> & {
  /**
   * The status code, which should be an enum value of
   * [google.rpc.Code][google.rpc.Code].
   *
   * @generated from field: int32 code = 1;
   */
  code: number;

  /**
   * A developer-facing error message, which should be in English. Any
   * user-facing error message should be localized and sent in the
   * [google.rpc.Status.details][google.rpc.Status.details] field, or localized
   * by the client.
   *
   * @generated from field: string message = 2;
   */
  message: string;

  /**
   * A list of messages that carry the error details.  There is a common set of
   * message types for APIs to use.
   *
   * @generated from field: repeated google.protobuf.Any details = 3;
   */
  details: Any[];
};

/**
 * Describes the message google.rpc.Status.
 * Use `create(StatusSchema)` to create a new message.
 */
export const StatusSchema: GenMessage<Status> =
  /*@__PURE__*/
  messageDesc(file_google_rpc_status, 0);
